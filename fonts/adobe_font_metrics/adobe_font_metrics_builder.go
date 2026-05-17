package adobe_font_metrics

import (
	"github.com/uglytoad/pdfpig/go/core"
)

// AdobeFontMetricsBuilder builds an AdobeFontMetrics instance incrementally.
type AdobeFontMetricsBuilder struct {
	AfmVersion       float64
	Comments         []string
	CharacterMetrics []AdobeFontMetricsIndividualCharacterMetric
	FontName         string
	FullName         string
	FamilyName       string
	Weight           string
	ItalicAngle      float64
	IsFixedPitch     bool
	PdfBoundingBox   core.PdfRectangle
	UnderlinePosition  float64
	UnderlineThickness float64
	Version          string
	Notice           string
	EncodingScheme   string
MappingScheme    int
	CharacterSet     string
	IsBaseFont       bool
	CapHeight        float64
	XHeight          float64
	Ascender         float64
	Descender        float64
	StdHw            float64
	StdVw            float64
	EscapeCharacter  int
	CharacterWidth   AdobeFontMetricsCharacterSize
	Characters       int
	VVector          AdobeFontMetricsVector
	IsFixedV         bool
}

// NewAdobeFontMetricsBuilder creates a new builder initialized with the given AFM version.
func NewAdobeFontMetricsBuilder(afmVersion float64) *AdobeFontMetricsBuilder {
	return &AdobeFontMetricsBuilder{
		AfmVersion:       afmVersion,
		Comments:         []string{},
		CharacterMetrics: []AdobeFontMetricsIndividualCharacterMetric{},
		IsBaseFont:       true,
	}
}

// SetBoundingBox sets the font bounding box dimensions.
func (b *AdobeFontMetricsBuilder) SetBoundingBox(x1, y1, x2, y2 float64) {
	b.PdfBoundingBox = core.NewPdfRectangleFloat(x1, y1, x2, y2)
}

// SetCharacterWidth sets the uniform character width and height.
func (b *AdobeFontMetricsBuilder) SetCharacterWidth(x, y float64) {
	b.CharacterWidth = NewAdobeFontMetricsCharacterSize(x, y)
}

// SetVVector sets the vector from writing direction 0 to direction 1.
func (b *AdobeFontMetricsBuilder) SetVVector(x, y float64) {
	b.VVector = AdobeFontMetricsVector{X: x, Y: y}
}

// Build constructs and returns the final AdobeFontMetrics instance.
func (b *AdobeFontMetricsBuilder) Build() AdobeFontMetrics {
	dictionary := make(map[string]AdobeFontMetricsIndividualCharacterMetric, len(b.CharacterMetrics))
	for _, m := range b.CharacterMetrics {
		dictionary[m.Name] = m
	}

	return NewAdobeFontMetrics(
		b.AfmVersion,
		b.Comments,
		0,
		b.FontName,
		b.FullName,
		b.FamilyName,
		b.Weight,
		b.PdfBoundingBox,
		b.Version,
		b.Notice,
		b.EncodingScheme,
		b.MappingScheme,
		b.EscapeCharacter,
		b.CharacterSet,
		b.Characters,
		b.IsBaseFont,
		b.VVector,
		b.IsFixedV,
		b.CapHeight,
		b.XHeight,
		b.Ascender,
		b.Descender,
		b.UnderlinePosition,
		b.UnderlineThickness,
		b.ItalicAngle,
		b.CharacterWidth,
		b.StdHw,
		b.StdVw,
		dictionary,
	)
}
