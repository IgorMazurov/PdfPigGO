package colors

import (
	"fmt"
)

// CMYKColor represents a color with cyan, magenta, yellow and black (K) components.
type CMYKColor struct {
	C float64
	M float64
	Y float64
	K float64
}

// NewCMYKColor creates a new CMYKColor with the given component values.
func NewCMYKColor(c, m, y, k float64) CMYKColor {
	return CMYKColor{C: c, M: m, Y: y, K: k}
}

// CMYKBlack is the CMYK black color (0, 0, 0, 1).
var CMYKBlack = Color(NewCMYKColor(0, 0, 0, 1))

// Black is the CMYK black color (0, 0, 0, 1).
var Black = CMYKBlack

// CMYKWhite is the CMYK white color (all components zero).
var CMYKWhite = Color(NewCMYKColor(0, 0, 0, 0))

// White is the CMYK white color (all components zero).
var White = CMYKWhite

// ColorSpace returns DeviceCMYK for this color.
func (c CMYKColor) ColorSpace() ColorSpace {
	return DeviceCMYK
}

// ToRGBValues converts the CMYK color to RGB values between 0 and 1.
func (c CMYKColor) ToRGBValues() RGBValues {
	k := 1 - c.K
	return RGBValues{
		R: (1 - c.C) * k,
		G: (1 - c.M) * k,
		B: (1 - c.Y) * k,
	}
}

// Equals reports whether this CMYKColor equals another.
func (c CMYKColor) Equals(other CMYKColor) bool {
	return c.C == other.C && c.M == other.M && c.Y == other.Y && c.K == other.K
}

// String returns a string representation of the CMYK color.
func (c CMYKColor) String() string {
	return fmt.Sprintf("CMYK: (%g, %g, %g, %g)", c.C, c.M, c.Y, c.K)
}
