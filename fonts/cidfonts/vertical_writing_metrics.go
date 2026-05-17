package cidfonts

import "github.com/uglytoad/pdfpig/go/geometry"

// VerticalWritingMetrics holds displacement and position vectors for glyphs in
// fonts that support vertical writing mode. The position vector transforms the
// horizontal writing origin into the vertical writing origin, and the
// displacement vector specifies how far to move vertically before drawing the
// next glyph.
type VerticalWritingMetrics struct {
	DefaultVerticalWritingMetrics     VerticalVectorComponents
	IndividualVerticalWritingDisplacements map[int]float64
	IndividualVerticalWritingPositions    map[int]geometry.PdfVector
}

// NewVerticalWritingMetrics creates a new VerticalWritingMetrics instance.
func NewVerticalWritingMetrics(
	defaultMetrics VerticalVectorComponents,
	individualDisplacements map[int]float64,
	individualPositions map[int]geometry.PdfVector,
) VerticalWritingMetrics {
	if individualDisplacements == nil {
		individualDisplacements = make(map[int]float64)
	}
	if individualPositions == nil {
		individualPositions = make(map[int]geometry.PdfVector)
	}
	return VerticalWritingMetrics{
		DefaultVerticalWritingMetrics:     defaultMetrics,
		IndividualVerticalWritingDisplacements: individualDisplacements,
		IndividualVerticalWritingPositions:    individualPositions,
	}
}

// GetPositionVector returns the position vector used to convert horizontal glyph
// origin to vertical origin. If an override exists for the given character
// identifier it is returned; otherwise the default is computed from the glyph width.
func (v VerticalWritingMetrics) GetPositionVector(characterIdentifier int, glyphWidth float64) geometry.PdfVector {
	if vec, ok := v.IndividualVerticalWritingPositions[characterIdentifier]; ok {
		return vec
	}
	return v.DefaultVerticalWritingMetrics.GetPositionVector(glyphWidth)
}

// GetDisplacementVector returns the displacement vector used to move the origin
// to the next glyph location after drawing. If an override exists for the given
// character identifier it is returned; otherwise the default displacement vector is used.
func (v VerticalWritingMetrics) GetDisplacementVector(characterIdentifier int) geometry.PdfVector {
	if dispY, ok := v.IndividualVerticalWritingDisplacements[characterIdentifier]; ok {
		return geometry.NewPdfVector(0, dispY)
	}
	return v.DefaultVerticalWritingMetrics.GetDisplacementVector()
}
