package annotations

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
)

// QuadPointsQuadrilateral is four points defining the region for an annotation to use.
// An annotation may cover multiple quadrilaterals.
type QuadPointsQuadrilateral struct {
	points []core.PdfPoint
}

// Points returns the 4 points defining this quadrilateral.
// The PDF specification defines these as being in anti-clockwise order starting from the lower-left corner, however
// Adobe's implementation doesn't obey the specification and points seem to go in the order: top-left, top-right,
// bottom-left, bottom-right. See: https://stackoverflow.com/questions/9855814/pdf-spec-vs-acrobat-creation-quadpoints.
func (q QuadPointsQuadrilateral) Points() []core.PdfPoint {
	return q.points
}

// NewQuadPointsQuadrilateral creates a new QuadPointsQuadrilateral.
func NewQuadPointsQuadrilateral(points []core.PdfPoint) (*QuadPointsQuadrilateral, error) {
	if len(points) != 4 {
		return nil, fmt.Errorf("quadpoints quadrilateral should only contain 4 points, instead got %d points", len(points))
	}

	return &QuadPointsQuadrilateral{points: points}, nil
}

// String returns a string representation of the quadrilateral.
func (q QuadPointsQuadrilateral) String() string {
	p := q.points
	return fmt.Sprintf("[ %s, %s, %s, %s ]", p[0], p[1], p[2], p[3])
}
