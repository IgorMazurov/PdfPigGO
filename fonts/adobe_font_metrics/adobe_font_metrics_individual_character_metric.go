package adobe_font_metrics

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
)

// NewAdobeFontMetricsIndividualCharacterMetric creates a new individual character metric.
func NewAdobeFontMetricsIndividualCharacterMetric(
	characterCode int,
	name string,
	width AdobeFontMetricsVector,
	widthDirection0 AdobeFontMetricsVector,
	widthDirection1 AdobeFontMetricsVector,
	vVector AdobeFontMetricsVector,
	boundingBox core.PdfRectangle,
	ligature *AdobeFontMetricsLigature,
) AdobeFontMetricsIndividualCharacterMetric {
	return AdobeFontMetricsIndividualCharacterMetric{
		CharacterCode:   characterCode,
		Name:            name,
		Width:           width,
		WidthDirection0: widthDirection0,
		WidthDirection1: widthDirection1,
		VVector:         vVector,
		BoundingBox:     boundingBox,
		Ligature:        ligature,
	}
}

// String returns a string representation of the character metric.
func (m AdobeFontMetricsIndividualCharacterMetric) String() string {
	return fmt.Sprintf("[%d] %s Width: %.4g.", m.CharacterCode, m.Name, m.Width.X)
}
