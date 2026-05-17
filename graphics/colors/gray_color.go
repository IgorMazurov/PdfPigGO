package colors

import (
	"fmt"
)

// GrayColor represents a grayscale color with a single gray component between 0 and 1.
type GrayColor struct {
	Gray float64
}

// NewGrayColor creates a new GrayColor with the given gray value.
func NewGrayColor(gray float64) GrayColor {
	return GrayColor{Gray: gray}
}

// GrayBlack is the grayscale black color (0).
var GrayBlack = Color(NewGrayColor(0))

// GrayWhite is the grayscale white color (1).
var GrayWhite = Color(NewGrayColor(1))

// ColorSpace returns DeviceGray for this color.
func (g GrayColor) ColorSpace() ColorSpace {
	return DeviceGray
}

// ToRGBValues converts the gray color to RGB values between 0 and 1.
func (g GrayColor) ToRGBValues() RGBValues {
	return RGBValues{R: g.Gray, G: g.Gray, B: g.Gray}
}

// Equals reports whether this GrayColor equals another.
func (g GrayColor) Equals(other GrayColor) bool {
	return g.Gray == other.Gray
}

// String returns a string representation of the gray color.
func (g GrayColor) String() string {
	return fmt.Sprintf("Gray: %g", g.Gray)
}
