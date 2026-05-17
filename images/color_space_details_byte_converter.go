package images

import (
	"math"

	"github.com/uglytoad/pdfpig/go/graphics/colors"
)

// Convert transforms decoded image bytes according to color space details,
// handling bit unpacking and stride padding removal.
func Convert(details colors.ColorSpaceDetails, decoded []byte, bitsPerComponent, imageWidth, imageHeight int) ([]byte, error) {
	if len(decoded) == 0 {
		return nil, nil
	}

	if details == nil {
		return decoded, nil
	}

	if bitsPerComponent != 8 {
		var err error
		decoded, err = unpackComponents(decoded, bitsPerComponent, details.Type())
		if err != nil {
			return nil, err
		}
	}

	bytesPerPixel := details.NumberOfColorComponents()
	strideWidth := len(decoded) / imageHeight / bytesPerPixel
	if strideWidth != imageWidth {
		if bytesPerPixel > 1 && imageWidth*imageHeight*bytesPerPixel < len(decoded) {
			decoded = decoded[:imageWidth*imageHeight*bytesPerPixel]
		} else {
			var err error
			decoded, err = removeStridePadding(decoded, strideWidth, imageWidth, imageHeight, bytesPerPixel)
			if err != nil {
				return nil, err
			}
		}
	}

	return details.Transform(decoded), nil
}

// unpackComponents expands sub-byte components to full bytes.
func unpackComponents(input []byte, bitsPerComponent int, colorSpace colors.ColorSpace) ([]byte, error) {
	if bitsPerComponent == 16 {
		size := len(input) / 2
		unpacked16 := make([]byte, size)
		for b := 0; b < size; b++ {
			i := 2 * b
			val := uint16(input[i])<<8 | uint16(input[i+1])
			unpacked16[b] = byte(val / 256)
		}
		return unpacked16, nil
	}

	end := 8 - bitsPerComponent
	componentCount := int(math.Ceil(float64(end+1) / float64(bitsPerComponent)))
	unpacked := make([]byte, len(input)*componentCount)
	right := (1 << bitsPerComponent) - 1
	u := 0

	if bitsPerComponent == 1 && colorSpace != colors.Indexed {
		for _, b := range input {
			for i := end; i >= 0; i-- {
				if (b>>uint(i))&1 == 1 {
					unpacked[u] = math.MaxUint8
				} else {
					unpacked[u] = 0
				}
				u++
			}
		}
		return unpacked, nil
	}

	for _, b := range input {
		for i := end; i >= 0; i -= bitsPerComponent {
			unpacked[u] = byte((b >> uint(i)) & byte(right))
			u++
		}
	}
	return unpacked, nil
}

// removeStridePadding removes extra padding bytes between image rows.
func removeStridePadding(input []byte, strideWidth, imageWidth, imageHeight, multiplier int) ([]byte, error) {
	size := imageWidth * imageHeight * multiplier
	result := make([]byte, size)
	for y := 0; y < imageHeight; y++ {
		sourceIndex := y * strideWidth * multiplier
		targetIndex := y * imageWidth * multiplier
		rowLen := imageWidth * multiplier
		copy(result[targetIndex:targetIndex+rowLen], input[sourceIndex:sourceIndex+rowLen])
	}
	return result, nil
}


