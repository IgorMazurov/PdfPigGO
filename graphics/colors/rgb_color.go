package colors

import (
	"fmt"
)

// RGBColor represents a color with red, green and blue components between 0 and 1.
type RGBColor struct {
	R float64
	G float64
	B float64
}

// NewRGBColor creates a new RGBColor with the given component values.
func NewRGBColor(r, g, b float64) RGBColor {
	return RGBColor{R: r, G: g, B: b}
}

// RGBBlack is the RGB black color (all components zero).
var RGBBlack = Color(NewRGBColor(0, 0, 0))

// RGBWhite is the RGB white color (all components one).
var RGBWhite = Color(NewRGBColor(1, 1, 1))

// ColorSpace returns DeviceRGB for this color.
func (c RGBColor) ColorSpace() ColorSpace {
	return DeviceRGB
}

// ToRGBValues returns the RGB component values between 0 and 1.
func (c RGBColor) ToRGBValues() RGBValues {
	return RGBValues{R: c.R, G: c.G, B: c.B}
}

// Equals reports whether this RGBColor equals another across all channels.
func (c RGBColor) Equals(other RGBColor) bool {
	return c.R == other.R && c.G == other.G && c.B == other.B
}

// String returns a string representation of the RGB color.
func (c RGBColor) String() string {
	return fmt.Sprintf("RGB: (%g, %g, %g)", c.R, c.G, c.B)
}
