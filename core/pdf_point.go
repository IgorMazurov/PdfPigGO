package core

import "fmt"

// PdfPoint represents a point in a PDF file.
// PDF coordinates are defined with the origin at the lower left (0, 0).
// The Y-axis extends vertically upwards and the X-axis horizontally to the right.
// Unless otherwise specified on a per-page basis, units in PDF space are
// equivalent to a typographic point (1/72 inch).
type PdfPoint struct {
	// X is the horizontal coordinate.
	X float64
	// Y is the vertical coordinate.
	Y float64
}

// Origin is the origin of the coordinate system (0, 0).
var Origin = PdfPoint{0.0, 0.0}

// NewPdfPoint creates a new PdfPoint with floating-point coordinates.
func NewPdfPoint(x, y float64) PdfPoint {
	return PdfPoint{x, y}
}

// PdfPointFromInt creates a new PdfPoint from integer coordinates.
func PdfPointFromInt(x, y int) PdfPoint {
	return PdfPoint{float64(x), float64(y)}
}

// MoveX returns a new point shifted on the X axis by dx.
func (p PdfPoint) MoveX(dx float64) PdfPoint {
	return PdfPoint{p.X + dx, p.Y}
}

// MoveY returns a new point shifted on the Y axis by dy.
func (p PdfPoint) MoveY(dy float64) PdfPoint {
	return PdfPoint{p.X, p.Y + dy}
}

// Translate returns a new point shifted by dx on X and dy on Y.
func (p PdfPoint) Translate(dx, dy float64) PdfPoint {
	return PdfPoint{p.X + dx, p.Y + dy}
}

// Equals reports whether p and other have the same coordinates.
func (p PdfPoint) Equals(other PdfPoint) bool {
	return p.X == other.X && p.Y == other.Y
}

// String returns a string representation of the point in the format "(x:X, y:Y)".
func (p PdfPoint) String() string {
	return fmt.Sprintf("(x:%g, y:%g)", p.X, p.Y)
}
