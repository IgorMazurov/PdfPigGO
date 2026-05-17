package png

import "fmt"

// ImageHeader holds the high-level information about a PNG image.
type ImageHeader struct {
	width             int
	height            int
	bitDepth          byte
	colorType         ColorType
	compressionMethod CompressionMethod
	filterMethod      FilterMethod
	interlaceMethod   InterlaceMethod
}

// HeaderBytes is the signature for the IHDR chunk type.
var HeaderBytes = []byte{73, 72, 68, 82} // "IHDR"

var permittedBitDepths = map[ColorType][]byte{
	ColorTypeNone:                    {1, 2, 4, 8, 16},
	ColorTypeColorUsed:               {8, 16},
	ColorTypePaletteUsed | ColorTypeColorUsed: {1, 2, 4, 8},
	ColorTypeAlphaChannelUsed:        {8, 16},
	ColorTypeAlphaChannelUsed | ColorTypeColorUsed: {8, 16},
}

// Width returns the width of the image in pixels.
func (h *ImageHeader) Width() int {
	return h.width
}

// Height returns the height of the image in pixels.
func (h *ImageHeader) Height() int {
	return h.height
}

// BitDepth returns the bit depth of the image.
func (h *ImageHeader) BitDepth() byte {
	return h.bitDepth
}

// ColorType returns the color type of the image.
func (h *ImageHeader) ColorType() ColorType {
	return h.colorType
}

// CompressionMethod returns the compression method used for the image.
func (h *ImageHeader) CompressionMethod() CompressionMethod {
	return h.compressionMethod
}

// FilterMethod returns the filter method used for the image.
func (h *ImageHeader) FilterMethod() FilterMethod {
	return h.filterMethod
}

// InterlaceMethod returns the interlace method used by the image.
func (h *ImageHeader) InterlaceMethod() InterlaceMethod {
	return h.interlaceMethod
}

// NewImageHeader creates a validated ImageHeader. Returns an error if width or height is 0,
// or if bitDepth is not permitted for the given colorType.
func NewImageHeader(
	width int,
	height int,
	bitDepth byte,
	colorType ColorType,
	compressionMethod CompressionMethod,
	filterMethod FilterMethod,
	interlaceMethod InterlaceMethod,
) (*ImageHeader, error) {
	if width == 0 {
		return nil, fmt.Errorf("invalid width (0) for image")
	}

	if height == 0 {
		return nil, fmt.Errorf("invalid height (0) for image")
	}

permitted, ok := permittedBitDepths[colorType]
	if !ok {
		return nil, fmt.Errorf("the bit depth %d is not permitted for color type %v", bitDepth, colorType)
	}

	found := false
	for _, bd := range permitted {
		if bd == bitDepth {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("the bit depth %d is not permitted for color type %v", bitDepth, colorType)
	}

	return &ImageHeader{
		width:             width,
		height:            height,
		bitDepth:          bitDepth,
		colorType:         colorType,
		compressionMethod: compressionMethod,
		filterMethod:      filterMethod,
		interlaceMethod:   interlaceMethod,
	}, nil
}

func (h *ImageHeader) String() string {
	return fmt.Sprintf("w: %d, h: %d, bitDepth: %d, colorType: %v, compression: %v, filter: %v, interlace: %v.",
		h.width, h.height, h.bitDepth, h.colorType, h.compressionMethod, h.filterMethod, h.interlaceMethod)
}
