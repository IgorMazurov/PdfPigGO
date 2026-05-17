package adobe_font_metrics

import "github.com/uglytoad/pdfpig/go/core"

// AdobeFontMetricsIndividualCharacterMetricBuilder accumulates fields for constructing an individual character metric.
type AdobeFontMetricsIndividualCharacterMetricBuilder struct {
	CharacterCode         int
	Name                  string
	WidthX                float64
	WidthY                float64
	WidthXDirection0      float64
	WidthYDirection0      float64
	WidthXDirection1      float64
	WidthYDirection1      float64
	VVector               AdobeFontMetricsVector
	BoundingBox           core.PdfRectangle
	Ligature              *AdobeFontMetricsLigature
}

// Build constructs an AdobeFontMetricsIndividualCharacterMetric from the accumulated fields.
func (b AdobeFontMetricsIndividualCharacterMetricBuilder) Build() AdobeFontMetricsIndividualCharacterMetric {
	return NewAdobeFontMetricsIndividualCharacterMetric(
		b.CharacterCode,
		b.Name,
		AdobeFontMetricsVector{X: b.WidthX, Y: b.WidthY},
		AdobeFontMetricsVector{X: b.WidthXDirection0, Y: b.WidthYDirection0},
		AdobeFontMetricsVector{X: b.WidthXDirection1, Y: b.WidthYDirection1},
		b.VVector,
		b.BoundingBox,
		b.Ligature,
	)
}
