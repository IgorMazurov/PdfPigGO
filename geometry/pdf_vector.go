package geometry

import (
	"fmt"
	"math"

	"github.com/uglytoad/pdfpig/go/core"
)

// PdfVector represents a 2D vector in PDF coordinate space.
type PdfVector struct {
	X float64
	Y float64
}

// NewPdfVector creates a new PdfVector with the given coordinates.
func NewPdfVector(x, y float64) PdfVector {
	return PdfVector{x, y}
}

// Scale returns a new vector scaled by the given factor.
func (v PdfVector) Scale(scale float64) PdfVector {
	return PdfVector{v.X * scale, v.Y * scale}
}

// GetMagnitude returns the magnitude (length) of the vector.
func (v PdfVector) GetMagnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Subtract returns a new vector that is the difference between v and other.
func (v PdfVector) Subtract(other PdfVector) PdfVector {
	return PdfVector{v.X - other.X, v.Y - other.Y}
}

// ToPoint converts the vector to a PdfPoint.
func (v PdfVector) ToPoint() core.PdfPoint {
	return core.NewPdfPoint(v.X, v.Y)
}

// String returns a string representation in the format "(X, Y)".
func (v PdfVector) String() string {
	return fmt.Sprintf("(%g, %g)", v.X, v.Y)
}
