package glyphs

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
)

// GlyphDescription describes a TrueType glyph, including its contour points,
// instructions, and bounding box. Implementations represent either simple or
// composite glyphs.
type GlyphDescription interface {
	// IsSimple reports whether the glyph is a simple (non-composite) glyph.
	IsSimple() bool

	// Bounds returns the bounding rectangle of the character.
	Bounds() core.PdfRectangle

	// Instructions returns the byte program for rendering this glyph.
	Instructions() []byte

	// EndPointsOfContours returns an array of the last point index of each contour.
	EndPointsOfContours() []uint16

	// Points returns the contour points that make up the glyph outline.
	Points() []GlyphPoint

	// IsEmpty reports whether the glyph has no points.
	IsEmpty() bool

	// TryGetGlyphPath attempts to convert the glyph points into a sequence of
	// PDF subpaths. Returns true if successful, false otherwise.
	TryGetGlyphPath() ([]*core.PdfSubpath, bool)

	// DeepClone creates an independent deep copy of this glyph description.
	DeepClone() GlyphDescription

	// Merge combines this glyph with another into a new glyph description.
	Merge(glyph GlyphDescription) GlyphDescription

	// Transform applies a 3×2 transformation matrix to all points and returns
	// a new transformed glyph description.
	Transform(matrix CompositeTransformMatrix3By2) GlyphDescription
}

// Glyph is the concrete implementation of GlyphDescription for TrueType glyphs.
type Glyph struct {
	bounds            core.PdfRectangle
	instructions      []byte
	endPointsOfContours []uint16
	points            []GlyphPoint
	isSimple          bool
}

// NewGlyph creates a new glyph with the given parameters.
func NewGlyph(isSimple bool, instructions []byte, endPointsOfContours []uint16, points []GlyphPoint, bounds core.PdfRectangle) *Glyph {
	return &Glyph{
		bounds:            bounds,
		instructions:      instructions,
		endPointsOfContours: endPointsOfContours,
		points:            points,
		isSimple:          isSimple,
	}
}

// EmptyGlyph returns a new empty glyph with the given bounding box.
func EmptyGlyph(bounds core.PdfRectangle) GlyphDescription {
	return NewGlyph(true, nil, nil, nil, bounds)
}

// IsSimple reports whether the glyph is a simple (non-composite) glyph.
func (g *Glyph) IsSimple() bool {
	return g.isSimple
}

// Bounds returns the bounding rectangle of the character.
func (g *Glyph) Bounds() core.PdfRectangle {
	return g.bounds
}

// Instructions returns the byte program for rendering this glyph.
func (g *Glyph) Instructions() []byte {
	return g.instructions
}

// EndPointsOfContours returns an array of the last point index of each contour.
func (g *Glyph) EndPointsOfContours() []uint16 {
	return g.endPointsOfContours
}

// Points returns the contour points that make up the glyph outline.
func (g *Glyph) Points() []GlyphPoint {
	return g.points
}

// IsEmpty reports whether the glyph has no points.
func (g *Glyph) IsEmpty() bool {
	return len(g.points) == 0
}

// DeepClone creates an independent deep copy of this glyph description.
func (g *Glyph) DeepClone() GlyphDescription {
	clonedInstructions := make([]byte, len(g.instructions))
	copy(clonedInstructions, g.instructions)

	clonedEndPoints := make([]uint16, len(g.endPointsOfContours))
	copy(clonedEndPoints, g.endPointsOfContours)

	clonedPoints := make([]GlyphPoint, len(g.points))
	copy(clonedPoints, g.points)

	return NewGlyph(false, clonedInstructions, clonedEndPoints, clonedPoints, g.bounds)
}

// Merge combines this glyph with another into a new glyph description.
func (g *Glyph) Merge(other GlyphDescription) GlyphDescription {
	newPoints := g.mergePoints(other)
	newEndpoints := g.mergeContourEndPoints(other)

	return NewGlyph(false, g.instructions, newEndpoints, newPoints, g.bounds)
}

// mergePoints concatenates the points of this glyph with those of another.
func (g *Glyph) mergePoints(other GlyphDescription) []GlyphPoint {
	newPoints := make([]GlyphPoint, len(g.points)+len(other.Points()))

	for i := 0; i < len(g.points); i++ {
		newPoints[i] = g.points[i]
	}

	for i := 0; i < len(other.Points()); i++ {
		newPoints[i+len(g.points)] = other.Points()[i]
	}

	return newPoints
}

// mergeContourEndPoints concatenates the contour end points, offsetting the
// second glyph's endpoints by the last endpoint of this glyph.
func (g *Glyph) mergeContourEndPoints(other GlyphDescription) []uint16 {
	var destinationLastEndPoint uint16 = 0
	if len(g.endPointsOfContours) > 0 {
		destinationLastEndPoint = g.endPointsOfContours[len(g.endPointsOfContours)-1] + 1
	}

	endPoints := make([]uint16, len(g.endPointsOfContours)+len(other.EndPointsOfContours()))

	for i := 0; i < len(g.endPointsOfContours); i++ {
		endPoints[i] = g.endPointsOfContours[i]
	}

	for i := 0; i < len(other.EndPointsOfContours()); i++ {
		endPoints[i+len(g.endPointsOfContours)] = other.EndPointsOfContours()[i] + destinationLastEndPoint
	}

	return endPoints
}

// Transform applies a 3×2 transformation matrix to all points and returns
// a new transformed glyph description.
func (g *Glyph) Transform(matrix CompositeTransformMatrix3By2) GlyphDescription {
	newPoints := make([]GlyphPoint, len(g.points))

	for i := len(g.points) - 1; i >= 0; i-- {
		point := g.points[i]

		scaled := matrix.ScaleAndRotate(core.NewPdfPoint(float64(point.X), float64(point.Y)))
		scaled = matrix.Translate(scaled)

		newPoints[i] = NewGlyphPoint(int16(scaled.X), int16(scaled.Y), point.IsOnCurve, point.IsEndOfContour)
	}

	return NewGlyph(g.isSimple, g.instructions, g.endPointsOfContours, newPoints, g.bounds)
}

// TryGetGlyphPath attempts to convert the glyph points into a sequence of PDF subpaths.
func (g *Glyph) TryGetGlyphPath() ([]*core.PdfSubpath, bool) {
	if g.points == nil {
		return nil, false
	}

	if len(g.points) > 0 {
		return calculatePath(g.points), true
	}

	return []*core.PdfSubpath{}, true
}

// calculatePath converts glyph contour points into PDF subpaths.
// Based on Apache PDFBox GlyphRenderer logic.
func calculatePath(points []GlyphPoint) []*core.PdfSubpath {
	path := make([]*core.PdfSubpath, 0)

	start := 0
	for p := 0; p < len(points); p++ {
		if !points[p].IsEndOfContour {
			continue
		}

		subpath := core.NewPdfSubpath()
		firstPoint := points[start]
		lastPoint := points[p]
		contour := make([]GlyphPoint, 0)

		for q := start; q <= p; q++ {
			contour = append(contour, points[q])
		}

		if points[start].IsOnCurve {
			contour = append(contour, firstPoint)
		} else if points[p].IsOnCurve {
			contour = append([]GlyphPoint{lastPoint}, contour...)
		} else {
			pmid := midValueGlyphPoint(firstPoint, lastPoint)
			contour = append([]GlyphPoint{pmid}, contour...)
			contour = append(contour, pmid)
		}

		subpath.MoveTo(float64(contour[0].X), float64(contour[0].Y))

		j := 1
		for j < len(contour) {
			pNow := contour[j]
			if pNow.IsOnCurve {
				subpath.LineTo(float64(pNow.X), float64(pNow.Y))
				j++
			} else if j+1 < len(contour) && contour[j+1].IsOnCurve {
				pNext := contour[j+1]
				subpath.BezierCurveToQuadratic(
					float64(pNow.X), float64(pNow.Y),
					float64(pNext.X), float64(pNext.Y),
				)
				j += 2
			} else if j+1 < len(contour) {
				pmid := midValueGlyphPoint(pNow, contour[j+1])
				subpath.BezierCurveToQuadratic(
					float64(pNow.X), float64(pNow.Y),
					float64(pmid.X), float64(pmid.Y),
				)
				j++
			} else {
				break
			}
		}

		subpath.CloseSubpath()
		path = append(path, subpath)
		start = p + 1
	}

	return path
}

// midValueShort computes the midpoint between two int16 values.
func midValueShort(a, b int16) int16 {
	return a + (b-a)/2
}

// midValueGlyphPoint creates an on-curve point that is between point1 and point2.
func midValueGlyphPoint(point1, point2 GlyphPoint) GlyphPoint {
	return NewGlyphPoint(
		midValueShort(point1.X, point2.X),
		midValueShort(point1.Y, point2.Y),
		true,
		false,
	)
}

// String returns a string representation of the glyph.
func (g *Glyph) String() string {
	glyphType := "C"
	if g.isSimple {
		glyphType = "S"
	}
	return fmt.Sprintf("%s: Width %g, Height: %g, Points: %d",
		glyphType, g.bounds.Width, g.bounds.Height, len(g.points))
}
