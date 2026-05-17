package glyphs

import "github.com/uglytoad/pdfpig/go/core"

// CompositeTransformMatrix3By2 represents a 3x2 transformation matrix used for
// composite glyph positioning in TrueType fonts. The matrix stores six values
// that define scaling, rotation, and translation transformations.
type CompositeTransformMatrix3By2 struct {
	r0c0 float64
	r0c1 float64
	r1c0 float64
	r1c1 float64
	r2c0 float64
	r2c1 float64
}

// Identity is the identity transformation matrix (no scaling, rotation, or translation).
var Identity = CompositeTransformMatrix3By2{1, 0, 0, 1, 0, 0}

// NewCompositeTransformMatrix3By2 creates a new 3x2 transformation matrix.
func NewCompositeTransformMatrix3By2(r0c0, r0c1, r1c0, r1c1, r2c0, r2c1 float64) CompositeTransformMatrix3By2 {
	return CompositeTransformMatrix3By2{r0c0, r0c1, r1c0, r1c1, r2c0, r2c1}
}

// CreateTranslation returns a translation-only transformation matrix.
func CreateTranslation(x, y float64) CompositeTransformMatrix3By2 {
	return CompositeTransformMatrix3By2{1, 0, 0, 1, x, y}
}

// WithTranslation returns a new matrix with the same scaling/rotation but
// different translation components.
func (m CompositeTransformMatrix3By2) WithTranslation(x, y float64) CompositeTransformMatrix3By2 {
	return CompositeTransformMatrix3By2{m.r0c0, m.r0c1, m.r1c0, m.r1c1, x, y}
}

// ScaleAndRotate applies only the scaling and rotation portion of the matrix
// to the source point, ignoring translation.
func (m CompositeTransformMatrix3By2) ScaleAndRotate(source core.PdfPoint) core.PdfPoint {
	newX := source.X*m.r0c0 + source.Y*m.r1c0
	newY := source.X*m.r0c1 + source.Y*m.r1c1
	return core.NewPdfPoint(newX, newY)
}

// Translate applies only the translation portion of the matrix to the source point.
func (m CompositeTransformMatrix3By2) Translate(source core.PdfPoint) core.PdfPoint {
	return core.NewPdfPoint(source.X+m.r2c0, source.Y+m.r2c1)
}
