// Package png provides utilities for converting PDF images to PNG format.
package png

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	"github.com/uglytoad/pdfpig/go/images"
)

// PdfImageSource defines the minimal interface required from a PDF image
// for PNG conversion. This avoids an import cycle with the content package
// while remaining compatible with content.PdfImage via structural typing.
type PdfImageSource interface {
	WidthInSamples() int
	HeightInSamples() int
	BitsPerComponent() int
	ColorSpaceDetails() *colors.ColorSpaceDetails
	TryGetBytes() ([]byte, bool)
	Decode() []float64
}

// TryGenerate attempts to convert a PDF image to PNG format.
// The mask parameter may be nil if no soft mask is available; when non-nil
// it must have the same dimensions as the main image.
// It handles color space conversion, CMYK-to-RGB transformation, and soft masks.
// Returns the PNG byte slice and true if successful, or nil and false otherwise.
func TryGenerate(image PdfImageSource, mask PdfImageSource) ([]byte, bool) {
	csdPtr := image.ColorSpaceDetails()

	if csdPtr == nil || colors.IsUnsupported(*csdPtr) {
		return nil, false
	}

	csd := *csdPtr

	if csd.BaseType() == colors.Pattern {
		return nil, false
	}

	rawBytes, ok := image.TryGetBytes()
	if !ok {
		return nil, false
	}

	bytesPure, err := images.Convert(csd, rawBytes, image.BitsPerComponent(), image.WidthInSamples(), image.HeightInSamples())
	if err != nil {
		return nil, false
	}

	numComponents := csd.BaseNumberOfColorComponents()

	softMask, hasMask := generateSoftMask(image, mask)

	getAlphaChannel := func(_ int) byte {
		return 255
	}

	if hasMask {
		smCopy := make([]byte, len(softMask))
		copy(smCopy, softMask)
		if needsReverseDecode(mask) {
			getAlphaChannel = func(i int) byte {
				return 255 - smCopy[i]
			}
		} else {
			getAlphaChannel = func(i int) byte {
				return smCopy[i]
			}
		}
	}

	if !isCorrectlySized(image, bytesPure) {
		return nil, false
	}

	builder := Create(image.WidthInSamples(), image.HeightInSamples(), hasMask)

	if csd.BaseType() == colors.DeviceCMYK || numComponents == 4 {
		i := 0
		sm := 0
		for col := 0; col < image.HeightInSamples(); col++ {
			for row := 0; row < image.WidthInSamples(); row++ {
				a := getAlphaChannel(sm)
				sm++

				c := float64(bytesPure[i]) / 255.0
				i++
				m := float64(bytesPure[i]) / 255.0
				i++
				y := float64(bytesPure[i]) / 255.0
				i++
				k := float64(bytesPure[i]) / 255.0
				i++

				r := byte(255 * (1 - c) * (1 - k))
				g := byte(255 * (1 - m) * (1 - k))
				b := byte(255 * (1 - y) * (1 - k))

				builder.SetPixel(NewPixel(r, g, b, a, false), row, col)
			}
		}
	} else if numComponents == 3 {
		i := 0
		sm := 0
		for col := 0; col < image.HeightInSamples(); col++ {
			for row := 0; row < image.WidthInSamples(); row++ {
				a := getAlphaChannel(sm)
				sm++
				r := bytesPure[i]
				i++
				g := bytesPure[i]
				i++
				b := bytesPure[i]
				i++
				builder.SetPixel(NewPixel(r, g, b, a, false), row, col)
			}
		}
	} else {
		i := 0
		if !needsReverseDecode(image) {
			for col := 0; col < image.HeightInSamples(); col++ {
				for row := 0; row < image.WidthInSamples(); row++ {
					a := getAlphaChannel(i)
					pixel := bytesPure[i]
					i++
					builder.SetPixel(NewPixel(pixel, pixel, pixel, a, true), row, col)
				}
			}
		} else {
			for col := 0; col < image.HeightInSamples(); col++ {
				for row := 0; row < image.WidthInSamples(); row++ {
					a := getAlphaChannel(i)
					pixel := byte(255 - bytesPure[i])
					i++
					builder.SetPixel(NewPixel(pixel, pixel, pixel, a, true), row, col)
				}
			}
		}
	}

	pngBytes := builder.SaveBytes()
	return pngBytes, true
}

// generateSoftMask attempts to generate a soft mask from the image's mask image.
// The mask is only applied if it has the same dimensions as the main image.
// Returns the converted mask bytes and true if successful, nil and false otherwise.
func generateSoftMask(image PdfImageSource, mask PdfImageSource) ([]byte, bool) {
	if mask == nil {
		return nil, false
	}

	if image.HeightInSamples() != mask.HeightInSamples() || image.WidthInSamples() != mask.WidthInSamples() {
		return nil, false
	}

	maskBytes, ok := mask.TryGetBytes()
	if !ok {
		return nil, false
	}

	csdPtr := mask.ColorSpaceDetails()
	if csdPtr == nil {
		return nil, false
	}

	converted, err := images.Convert(*csdPtr, maskBytes, mask.BitsPerComponent(), mask.WidthInSamples(), mask.HeightInSamples())
	if err != nil {
		return nil, false
	}

	if !isCorrectlySized(mask, converted) {
		return nil, false
	}

	return converted, true
}

// isCorrectlySized checks whether the byte slice has the expected size for the given image dimensions.
// Per PDF spec p.37, an extra end-of-line marker (LF, CR, or CRLF) at the end is tolerated.
func isCorrectlySized(image PdfImageSource, bytesPure []byte) bool {
	csdPtr := image.ColorSpaceDetails()
	if csdPtr == nil {
		return false
	}

	numComponents := (*csdPtr).BaseNumberOfColorComponents()
	requiredSize := image.WidthInSamples() * image.HeightInSamples() * numComponents
	actualSize := len(bytesPure)

	if actualSize == requiredSize {
		return true
	}

	if actualSize == requiredSize+1 {
		lastByte := bytesPure[actualSize-1]
		if lastByte == core.AsciiLineFeed || lastByte == core.AsciiCarriageReturn {
			return true
		}
	}

	if actualSize == requiredSize+2 {
		lastTwo := bytesPure[actualSize-2:]
		if len(lastTwo) == 2 && lastTwo[0] == core.AsciiCarriageReturn && lastTwo[1] == core.AsciiLineFeed {
			return true
		}
	}

	return false
}

// needsReverseDecode reports whether the image colors need to be reversed based on the Decode array.
// Returns true when the Decode array contains at least two elements with values [1, 0].
func needsReverseDecode(img PdfImageSource) bool {
	decode := img.Decode()
	return len(decode) >= 2 && decode[0] == 1 && decode[1] == 0
}
