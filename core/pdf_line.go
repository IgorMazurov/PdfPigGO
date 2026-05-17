package core

import "math"

// PdfLine represents a line in a PDF file.
// PDF coordinates are defined with the origin at the lower left (0, 0).
// The Y-axis extends vertically upwards and the X-axis horizontally to the right.
// Unless otherwise specified on a per-page basis, units in PDF space are
// equivalent to a typographic point (1/72 inch).
type PdfLine struct {
	// Point1 is the first point of the line.
	Point1 PdfPoint
	// Point2 is the second point of the line.
	Point2 PdfPoint
}

// NewPdfLine creates a new PdfLine from two points.
func NewPdfLine(point1, point2 PdfPoint) PdfLine {
	return PdfLine{point1, point2}
}

// NewPdfLineFromCoords creates a new PdfLine from four coordinates.
func NewPdfLineFromCoords(x1, y1, x2, y2 float64) PdfLine {
	return NewPdfLine(NewPdfPoint(x1, y1), NewPdfPoint(x2, y2))
}

// Length returns the length of the line (Euclidean distance between Point1 and Point2).
func (l PdfLine) Length() float64 {
	dx := l.Point1.X - l.Point2.X
	dy := l.Point1.Y - l.Point2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// GetBoundingRectangle returns the rectangle completely containing the line.
func (l PdfLine) GetBoundingRectangle() PdfRectangle {
	return NewPdfRectangleFloat(
		math.Min(l.Point1.X, l.Point2.X),
		math.Min(l.Point1.Y, l.Point2.Y),
		math.Max(l.Point1.X, l.Point2.X),
		math.Max(l.Point1.Y, l.Point2.Y),
	)
}

// Equals reports whether l and other have the same endpoints.
func (l PdfLine) Equals(other PdfLine) bool {
	return l.Point1.Equals(other.Point1) && l.Point2.Equals(other.Point2)
}
