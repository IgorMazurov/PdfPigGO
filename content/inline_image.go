package content

import (
	"fmt"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/images/png"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// InlineImage represents a small image completely defined inline within a page's content stream.
type InlineImage struct {
	boundingBox      core.PdfRectangle
	widthInSamples   int
	heightInSamples  int
	bitsPerComponent int
	isImageMask      bool
	decode           []float64
	renderingIntent  graphiccore.RenderingIntent
	interpolate      bool
	imageDictionary  *tokens.DictionaryToken
	rawMemory        []byte
	colorSpaceDetails *colors.ColorSpaceDetails
	maskImage        PdfImage

	filterProvider filters.FilterProvider
	filterNames    []*tokens.NameToken

	decodedOnce sync.Once
	decoded     []byte
	supported   bool
}

var _ PdfImage = (*InlineImage)(nil)

// NewInlineImage creates a new InlineImage.
func NewInlineImage(
	bounds core.PdfRectangle,
	widthInSamples int,
	heightInSamples int,
	bitsPerComponent int,
	isImageMask bool,
	renderingIntent graphiccore.RenderingIntent,
	interpolate bool,
	decode []float64,
	rawMemory []byte,
	filterProvider filters.FilterProvider,
	filterNames []*tokens.NameToken,
	streamDictionary *tokens.DictionaryToken,
	colorSpaceDetails *colors.ColorSpaceDetails,
	softMaskImage PdfImage,
) *InlineImage {
	var filtersList []filters.Filter

	if filterProvider != nil {
		filtersList, _ = filterProvider.GetNamedFilters(filterNames)
	}

	supported := true
	for _, f := range filtersList {
		if !f.IsSupported() {
			supported = false
			break
		}
	}

	img := &InlineImage{
		boundingBox:       bounds,
		widthInSamples:    widthInSamples,
		heightInSamples:   heightInSamples,
		bitsPerComponent:  bitsPerComponent,
		isImageMask:       isImageMask,
		renderingIntent:   renderingIntent,
		interpolate:       interpolate,
		decode:            decode,
		imageDictionary:   streamDictionary,
		rawMemory:         rawMemory,
		colorSpaceDetails: colorSpaceDetails,
		maskImage:         softMaskImage,
		filterProvider:    filterProvider,
		filterNames:       filterNames,
		supported:         supported,
	}

	return img
}

// BoundingBox returns the bounding rectangle of the image in PDF coordinates.
func (i *InlineImage) BoundingBox() core.PdfRectangle {
	return i.boundingBox
}

// Bounds returns the placement rectangle of the image in PDF coordinates.
// Deprecated: use BoundingBox() instead.
func (i *InlineImage) Bounds() core.PdfRectangle {
	return i.boundingBox
}

// WidthInSamples returns the width of the image in samples.
func (i *InlineImage) WidthInSamples() int {
	return i.widthInSamples
}

// HeightInSamples returns the height of the image in samples.
func (i *InlineImage) HeightInSamples() int {
	return i.heightInSamples
}

// BitsPerComponent returns the number of bits used to represent each color component.
func (i *InlineImage) BitsPerComponent() int {
	return i.bitsPerComponent
}

// IsImageMask reports whether the image is to be treated as an image mask.
func (i *InlineImage) IsImageMask() bool {
	return i.isImageMask
}

// Decode returns the decode array that maps image samples into values
// appropriate for the color space.
func (i *InlineImage) Decode() []float64 {
	return i.decode
}

// IsInlineImage reports whether this image is an inline image.
func (i *InlineImage) IsInlineImage() bool {
	return true
}

// ImageDictionary returns the full dictionary token for this image object.
func (i *InlineImage) ImageDictionary() *tokens.DictionaryToken {
	return i.imageDictionary
}

// RenderingIntent returns the color rendering intent for the image.
func (i *InlineImage) RenderingIntent() graphiccore.RenderingIntent {
	return i.renderingIntent
}

// Interpolate reports whether interpolation should be performed when rendering.
func (i *InlineImage) Interpolate() bool {
	return i.interpolate
}

// RawMemory returns the encoded bytes of the image with all filters still applied.
func (i *InlineImage) RawMemory() []byte {
	return i.rawMemory
}

// RawBytes returns the encoded byte slice of the image with all filters still applied.
func (i *InlineImage) RawBytes() []byte {
	return i.rawMemory
}

// ColorSpaceDetails returns the color space details used to interpret the image.
func (i *InlineImage) ColorSpaceDetails() *colors.ColorSpaceDetails {
	return i.colorSpaceDetails
}

// MaskImage returns the image mask, either a soft-mask or a stencil mask.
func (i *InlineImage) MaskImage() PdfImage {
	return i.maskImage
}

// TryGetBytes attempts to return the decoded bytes of the image.
// Returns false if any filter in the chain is not supported or decoding fails.
func (i *InlineImage) TryGetBytes() ([]byte, bool) {
	if !i.supported || i.filterProvider == nil {
		return nil, false
	}

	var success bool
	i.decodedOnce.Do(func() {
		filtersList, err := i.filterProvider.GetNamedFilters(i.filterNames)
		if err != nil {
			return
		}

		b := i.rawMemory
		for idx, f := range filtersList {
			b, err = f.Decode(b, i.imageDictionary, i.filterProvider, idx)
			if err != nil {
				return
			}
		}

		i.decoded = b
		success = true
	})

	if !success {
		return nil, false
	}

	return i.decoded, true
}

// TryGetPng attempts to convert the image to PNG format.
func (i *InlineImage) TryGetPng() ([]byte, bool) {
	return png.TryGenerate(i, i.maskImage)
}

// String returns a human-readable description of the inline image.
func (i *InlineImage) String() string {
	return fmt.Sprintf("Inline Image (w %g, h %g)", i.boundingBox.Width, i.boundingBox.Height)
}
