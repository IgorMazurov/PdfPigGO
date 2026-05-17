package adobe_font_metrics

import "fmt"

// AdobeFontMetricsCharacterSize holds the x and y components of the width vector
// of a font's characters. Presence implies that IsFixedPitch is true.
type AdobeFontMetricsCharacterSize struct {
	// X is the horizontal width component.
	X float64

	// Y is the vertical width component.
	Y float64
}

// NewAdobeFontMetricsCharacterSize creates a new character size instance.
func NewAdobeFontMetricsCharacterSize(x, y float64) AdobeFontMetricsCharacterSize {
	return AdobeFontMetricsCharacterSize{
		X: x,
		Y: y,
	}
}

// String returns the string representation of the character size.
func (s AdobeFontMetricsCharacterSize) String() string {
	return fmt.Sprintf("%g, %g", s.X, s.Y)
}
