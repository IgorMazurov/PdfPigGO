package images

// JpegInformation holds information read from a JPEG image.
type JpegInformation struct {
	// Width of the image in pixels.
	Width int

	// Height of the image in pixels.
	Height int

	// BitsPerComponent is the bits per component.
	BitsPerComponent int

	// NumberOfComponents is 1 for grayscale, 3 for RGB, 4 for CMYK.
	NumberOfComponents int
}

// NewJpegInformation creates a new JpegInformation.
func NewJpegInformation(width, height, bitsPerComponent, numberOfComponents int) *JpegInformation {
	return &JpegInformation{
		Width:              width,
		Height:             height,
		BitsPerComponent:   bitsPerComponent,
		NumberOfComponents: numberOfComponents,
	}
}
