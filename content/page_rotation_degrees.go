package content

import (
	"fmt"
	"math"
)

// PageRotationDegrees represents the rotation of a page in a PDF document
// defined by the page dictionary in degrees clockwise.
type PageRotationDegrees struct {
	// Value is the rotation of the page in degrees clockwise.
	Value int
}

// SwapsAxis reports whether the rotation flips the x and y axes.
func (p PageRotationDegrees) SwapsAxis() bool {
	return p.Value == 90 || p.Value == 270
}

// Radians returns the rotation expressed in radians (anti-clockwise).
func (p PageRotationDegrees) Radians() float64 {
	switch p.Value {
	case 0:
		return 0
	case 90:
		return -0.5 * math.Pi
	case 180:
		return -math.Pi
	case 270:
		return -1.5 * math.Pi
	default:
		panic(fmt.Sprintf("invalid value for rotation: %d", p.Value))
	}
}

// NewPageRotationDegrees creates a PageRotationDegrees from the given rotation in degrees clockwise.
// Rotation must be a multiple of 90 (0, 90, 180, or 270). Negative values are normalized by adding 360,
// and values >= 360 are reduced modulo 360.
func NewPageRotationDegrees(rotation int) (PageRotationDegrees, error) {
	if rotation < 0 {
		rotation = 360 + rotation
	}

	for rotation >= 360 {
		rotation -= 360
	}

	if rotation != 0 && rotation != 90 && rotation != 180 && rotation != 270 {
		return PageRotationDegrees{}, fmt.Errorf("rotation must be 0, 90, 180 or 270. Got: %d", rotation)
	}

	return PageRotationDegrees{Value: rotation}, nil
}

// Equals reports whether p and other represent the same rotation.
func (p PageRotationDegrees) Equals(other PageRotationDegrees) bool {
	return p.Value == other.Value
}

// String returns the string representation of the rotation value.
func (p PageRotationDegrees) String() string {
	return fmt.Sprintf("%d", p.Value)
}

var _ fmt.Stringer = PageRotationDegrees{}
