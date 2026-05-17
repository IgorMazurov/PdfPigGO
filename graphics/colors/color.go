package colors

// RGBValues holds red, green, and blue color component values between 0 and 1.
type RGBValues struct {
	R float64
	G float64
	B float64
}

// Color represents a color used for text or paths in a PDF.
type Color interface {
	// ColorSpace returns the color space used for this color.
	ColorSpace() ColorSpace

	// ToRGBValues converts the color to RGB values between 0 and 1.
	ToRGBValues() RGBValues
}
