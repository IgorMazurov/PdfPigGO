package png

import "fmt"

// Pixel represents a single pixel with RGBA color components.
type Pixel struct {
	R         byte
	G         byte
	B         byte
	A         byte
	IsGrayscale bool
}

// NewPixel creates a new Pixel with the given RGBA values and grayscale flag.
func NewPixel(r, g, b, a byte, isGrayscale bool) Pixel {
	return Pixel{R: r, G: g, B: b, A: a, IsGrayscale: isGrayscale}
}

// NewRgbPixel creates a fully opaque RGB pixel (alpha = 255).
func NewRgbPixel(r, g, b byte) Pixel {
	return Pixel{R: r, G: g, B: b, A: 255, IsGrayscale: false}
}

// NewGrayPixel creates a grayscale pixel with equal RGB channels and full opacity.
func NewGrayPixel(gray byte) Pixel {
	return Pixel{R: gray, G: gray, B: gray, A: 255, IsGrayscale: true}
}

// Equals reports whether p and other have identical RGBA values and grayscale flag.
func (p Pixel) Equals(other Pixel) bool {
	return p.R == other.R && p.G == other.G && p.B == other.B && p.A == other.A && p.IsGrayscale == other.IsGrayscale
}

// String returns a string representation of the pixel.
func (p Pixel) String() string {
	if p.IsGrayscale {
		return fmt.Sprintf("Pixel(G=%d, A=%d)", p.R, p.A)
	}
	return fmt.Sprintf("Pixel(R=%d, G=%d, B=%d, A=%d)", p.R, p.G, p.B, p.A)
}
