package core

import (
	"fmt"
	"math"
)

// PdfRectangle represents a rectangle in a PDF file.
// PDF coordinates are defined with the origin at the lower left (0, 0).
// The Y-axis extends vertically upwards and the X-axis horizontally to the right.
// Unless otherwise specified on a per-page basis, units in PDF space are
// equivalent to a typographic point (1/72 inch).
type PdfRectangle struct {
	// TopLeft is the top left point of the rectangle.
	TopLeft PdfPoint
	// TopRight is the top right point of the rectangle.
	TopRight PdfPoint
	// BottomRight is the bottom right point of the rectangle.
	BottomRight PdfPoint
	// BottomLeft is the bottom left point of the rectangle.
	BottomLeft PdfPoint
	// Width is the width of the rectangle (always positive).
	Width float64
	// Height is the height of the rectangle (always positive).
	Height float64
}

// NewPdfRectangle creates a new PdfRectangle from bottom-left and top-right points.
func NewPdfRectangle(bottomLeft, topRight PdfPoint) PdfRectangle {
	return NewPdfRectangleFloat(bottomLeft.X, bottomLeft.Y, topRight.X, topRight.Y)
}

// NewPdfRectangleFromInt creates a new PdfRectangle from integer coordinates.
// x1,y1 is the bottom-left corner; x2,y2 is the top-right corner.
func NewPdfRectangleFromInt(x1, y1, x2, y2 int) PdfRectangle {
	return NewPdfRectangleFloat(float64(x1), float64(y1), float64(x2), float64(y2))
}

// NewPdfRectangleFloat creates a new PdfRectangle from floating-point coordinates.
// x1,y1 is the bottom-left corner; x2,y2 is the top-right corner.
func NewPdfRectangleFloat(x1, y1, x2, y2 float64) PdfRectangle {
	return NewPdfRectangleFromCorners(
		NewPdfPoint(x1, y2), // topLeft
		NewPdfPoint(x2, y2), // topRight
		NewPdfPoint(x1, y1), // bottomLeft
		NewPdfPoint(x2, y1), // bottomRight
	)
}

// NewPdfRectangleFromCorners creates a new PdfRectangle from four corner points.
func NewPdfRectangleFromCorners(topLeft, topRight, bottomLeft, bottomRight PdfPoint) PdfRectangle {
	w := math.Sqrt((bottomLeft.X-bottomRight.X)*(bottomLeft.X-bottomRight.X) + (bottomLeft.Y-bottomRight.Y)*(bottomLeft.Y-bottomRight.Y))
	h := math.Sqrt((bottomLeft.X-topLeft.X)*(bottomLeft.X-topLeft.X) + (bottomLeft.Y-topLeft.Y)*(bottomLeft.Y-topLeft.Y))
	return PdfRectangle{
		TopLeft:     topLeft,
		TopRight:    topRight,
		BottomLeft:  bottomLeft,
		BottomRight: bottomRight,
		Width:       w,
		Height:      h,
	}
}

// Centroid returns the centroid point of the rectangle.
func (r PdfRectangle) Centroid() PdfPoint {
	cx := (r.BottomRight.X + r.TopRight.X + r.TopLeft.X + r.BottomLeft.X) / 4.0
	cy := (r.BottomRight.Y + r.TopRight.Y + r.TopLeft.Y + r.BottomLeft.Y) / 4.0
	return NewPdfPoint(cx, cy)
}

// Rotation returns the rotation angle of the rectangle in degrees, counterclockwise.
// The value is in the range [-180, 180].
func (r PdfRectangle) Rotation() float64 {
	return r.getTheta() * 180 / math.Pi
}

// Area returns the area of the rectangle.
func (r PdfRectangle) Area() float64 {
	return math.Abs(r.Width * r.Height)
}

// Left returns the left coordinate of the rectangle.
// This value is only valid if the rectangle is not rotated; check Rotation().
func (r PdfRectangle) Left() float64 {
	if r.TopLeft.X < r.TopRight.X {
		return r.TopLeft.X
	}
	return r.TopRight.X
}

// Top returns the top coordinate of the rectangle.
// This value is only valid if the rectangle is not rotated; check Rotation().
func (r PdfRectangle) Top() float64 {
	if r.TopLeft.Y > r.BottomLeft.Y {
		return r.TopLeft.Y
	}
	return r.BottomLeft.Y
}

// Right returns the right coordinate of the rectangle.
// This value is only valid if the rectangle is not rotated; check Rotation().
func (r PdfRectangle) Right() float64 {
	if r.BottomRight.X > r.BottomLeft.X {
		return r.BottomRight.X
	}
	return r.BottomLeft.X
}

// Bottom returns the bottom coordinate of the rectangle.
// This value is only valid if the rectangle is not rotated; check Rotation().
func (r PdfRectangle) Bottom() float64 {
	if r.BottomRight.Y < r.TopRight.Y {
		return r.BottomRight.Y
	}
	return r.TopRight.Y
}

// Translate returns a new rectangle shifted by dx on X and dy on Y.
func (r PdfRectangle) Translate(dx, dy float64) PdfRectangle {
	return NewPdfRectangleFromCorners(
		r.TopLeft.Translate(dx, dy),
		r.TopRight.Translate(dx, dy),
		r.BottomLeft.Translate(dx, dy),
		r.BottomRight.Translate(dx, dy),
	)
}

// getTheta computes the rotation angle in radians: -pi <= theta <= pi.
func (r PdfRectangle) getTheta() float64 {
	if !r.BottomRight.Equals(r.BottomLeft) {
		return math.Atan2(r.BottomRight.Y-r.BottomLeft.Y, r.BottomRight.X-r.BottomLeft.X)
	}
	// Handle the case where both bottom points are identical.
	return math.Atan2(r.TopLeft.Y-r.BottomLeft.Y, r.TopLeft.X-r.BottomLeft.X) - math.Pi/2
}

// Equals reports whether r and other have the same corner points.
func (r PdfRectangle) Equals(other PdfRectangle) bool {
	return r.TopLeft.Equals(other.TopLeft) &&
		r.TopRight.Equals(other.TopRight) &&
		r.BottomRight.Equals(other.BottomRight) &&
		r.BottomLeft.Equals(other.BottomLeft)
}

// Contains reports whether the point p lies within the rectangle.
// If inclusive is true, points on the boundary are considered inside.
// For rotated rectangles, uses an area-based approach to correctly determine containment.
func (r PdfRectangle) Contains(p PdfPoint, inclusive bool) bool {
	const eps = 1e-5

	if math.Abs(r.Area()) < eps {
		return false
	}

	if math.Abs(r.Rotation()) < eps {
		if inclusive {
			return p.X >= r.Left() && p.X <= r.Right() && p.Y >= r.Bottom() && p.Y <= r.Top()
		}
		return p.X > r.Left() && p.X < r.Right() && p.Y > r.Bottom() && p.Y < r.Top()
	}

	area3 := func(p1, p2, p3 PdfPoint) float64 {
		return math.Abs((p2.X*p1.Y-p1.X*p2.Y)+(p3.X*p2.Y-p2.X*p3.Y)+(p1.X*p3.Y-p3.X*p1.Y)) / 2.0
	}

	a1 := area3(r.BottomLeft, p, r.TopLeft)
	a2 := area3(r.TopLeft, p, r.TopRight)
	a3 := area3(r.TopRight, p, r.BottomRight)
	a4 := area3(r.BottomRight, p, r.BottomLeft)

	sum := a1 + a2 + a3 + a4

	if sum-r.Area() > eps {
		return false
	}

	if a1 < eps || a2 < eps || a3 < eps || a4 < eps {
		return inclusive
	}

	return true
}

// Rotate returns a new rectangle rotated by degrees counterclockwise around the origin.
func (r PdfRectangle) Rotate(degrees float64) PdfRectangle {
	return GetRotationMatrix(degrees).TransformRect(r)
}

// Scale returns a new rectangle scaled by scaleX on X and scaleY on Y from the origin.
func (r PdfRectangle) Scale(scaleX, scaleY float64) PdfRectangle {
	return GetScaleMatrix(scaleX, scaleY).TransformRect(r)
}

// Corners returns the four corner points of the rectangle in order:
// BottomRight, TopRight, TopLeft, BottomLeft.
func (r PdfRectangle) Corners() [4]PdfPoint {
	return [4]PdfPoint{
		r.BottomRight,
		r.TopRight,
		r.TopLeft,
		r.BottomLeft,
	}
}

// String returns a string representation in the format "[topLeft, width, height]".
func (r PdfRectangle) String() string {
	return fmt.Sprintf("[%s, %g, %g]", r.TopLeft, r.Width, r.Height)
}
