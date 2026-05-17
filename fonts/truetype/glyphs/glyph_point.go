package glyphs

import "fmt"

// GlyphPoint represents a single point on a glyph contour in a TrueType font.
type GlyphPoint struct {
	// X is the horizontal coordinate of the point.
	X int16

	// Y is the vertical coordinate of the point.
	Y int16

	// IsOnCurve indicates whether this point lies directly on the curve,
	// as opposed to being a control point for a Bézier segment.
	IsOnCurve bool

	// IsEndOfContour marks the last point in a contour.
	IsEndOfContour bool
}

// NewGlyphPoint creates a new GlyphPoint with the given coordinates and flags.
func NewGlyphPoint(x, y int16, isOnCurve, isEndOfContour bool) GlyphPoint {
	return GlyphPoint{
		X:              x,
		Y:              y,
		IsOnCurve:      isOnCurve,
		IsEndOfContour: isEndOfContour,
	}
}

// String returns a string representation of the glyph point.
func (p GlyphPoint) String() string {
	return fmt.Sprintf("(%d, %d) | %v | %v", p.X, p.Y, p.IsOnCurve, p.IsEndOfContour)
}
