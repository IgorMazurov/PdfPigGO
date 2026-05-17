package adobe_font_metrics

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
)

// AdobeFontMetricsVector represents a 2D vector used in AFM metrics.
type AdobeFontMetricsVector struct {
	X float64
	Y float64
}

// NewAdobeFontMetricsVector creates a new vector instance.
func NewAdobeFontMetricsVector(x, y float64) AdobeFontMetricsVector {
	return AdobeFontMetricsVector{
		X: x,
		Y: y,
	}
}

// String returns the string representation of the vector.
func (v AdobeFontMetricsVector) String() string {
	return fmt.Sprintf("%g, %g", v.X, v.Y)
}

// AdobeFontMetricsIndividualCharacterMetric holds metrics for a single character in an AFM font file.
type AdobeFontMetricsIndividualCharacterMetric struct {
	// CharacterCode is the character code.
	CharacterCode int

	// Name is the PostScript language character name.
	Name string

	// Width is the width vector.
	Width AdobeFontMetricsVector

	// WidthDirection0 is the width for writing direction 0 (horizontal).
	WidthDirection0 AdobeFontMetricsVector

	// WidthDirection1 is the width for writing direction 1 (vertical).
	WidthDirection1 AdobeFontMetricsVector

	// VVector is the vector from origin of writing direction 1 to origin of writing direction 0.
	VVector AdobeFontMetricsVector

	// BoundingBox is the character bounding box.
	BoundingBox core.PdfRectangle

	// Ligature is the ligature information (nil if none).
	Ligature *AdobeFontMetricsLigature
}

// AdobeFontMetricsLigature represents a ligature in an Adobe Font Metrics individual character.
type AdobeFontMetricsLigature struct {
	// Successor is the character to join with to form a ligature.
	Successor string

	// Value is the current character.
	Value string
}

// NewAdobeFontMetricsLigature creates a new ligature instance.
func NewAdobeFontMetricsLigature(successor, value string) AdobeFontMetricsLigature {
	return AdobeFontMetricsLigature{
		Successor: successor,
		Value:     value,
	}
}

// String returns the string representation of the ligature.
func (l AdobeFontMetricsLigature) String() string {
	return fmt.Sprintf("Ligature: %s -> Successor: %s", l.Value, l.Successor)
}

// AdobeFontMetrics holds the global metrics for a font program and the metrics of each character.
type AdobeFontMetrics struct {
	// AfmVersion is the version of the Adobe Font Metrics specification used to generate this file.
	AfmVersion float64

	// Comments contains any comments in the file.
	Comments []string

	// MetricSets describes the writing directions described by these metrics.
	MetricSets AdobeFontMetricsWritingDirection

	// FontName is the font name.
	FontName string

	// FullName is the font full name.
	FullName string

	// FamilyName is the font family name.
	FamilyName string

	// Weight is the font weight.
	Weight string

	// BoundingBox is the minimum bounding box for all characters in the font.
	BoundingBox core.PdfRectangle

	// Version is the font program version identifier.
	Version string

	// Notice is the font name trademark or copyright notice.
	Notice string

	// EncodingScheme indicates the default encoding vector for this font program.
	// Common ones are AdobeStandardEncoding and JIS12-88-CFEncoding.
	// Special font programs might state FontSpecific.
	EncodingScheme string

	// MappingScheme describes the mapping scheme.
	MappingScheme int

	// EscapeCharacter is the byte value of the escape character used if this font is escape-mapped.
	EscapeCharacter int

	// CharacterSet describes the character set of this font.
	CharacterSet string

	// Characters is the number of characters in this font.
	Characters int

	// IsBaseFont indicates whether this is a base font.
	IsBaseFont bool

	// VVector is a vector from the origin of writing direction 0 to direction 1.
	VVector AdobeFontMetricsVector

	// IsFixedV indicates whether VVector is the same for every character in this font.
	IsFixedV bool

	// CapHeight is usually the y-value of the top of capital 'H'.
	CapHeight float64

	// XHeight is usually the y-value of the top of lowercase 'x'.
	XHeight float64

	// Ascender is usually the y-value of the top of lowercase 'd'.
	Ascender float64

	// Descender is usually the y-value of the bottom of lowercase 'p'.
	Descender float64

	// UnderlinePosition is the distance from the baseline for underlining.
	UnderlinePosition float64

	// UnderlineThickness is the width of the line for underlining.
	UnderlineThickness float64

	// ItalicAngle is the angle in degrees counter-clockwise from the vertical of the vertical linea.
	// Zero for non-italic fonts.
	ItalicAngle float64

	// CharacterWidth, if present, means all characters have this width and height.
	CharacterWidth AdobeFontMetricsCharacterSize

	// HorizontalStemWidth is the horizontal stem width.
	HorizontalStemWidth float64

	// VerticalStemWidth is the vertical stem width.
	VerticalStemWidth float64

	// CharacterMetrics holds metrics for the individual characters.
	CharacterMetrics map[string]AdobeFontMetricsIndividualCharacterMetric
}

// NewAdobeFontMetrics creates a new AdobeFontMetrics instance.
func NewAdobeFontMetrics(
	afmVersion float64,
	comments []string,
	metricSets AdobeFontMetricsWritingDirection,
	fontName string,
	fullName string,
	familyName string,
	weight string,
	boundingBox core.PdfRectangle,
	version string,
	notice string,
	encodingScheme string,
	mappingScheme int,
	escapeCharacter int,
	characterSet string,
	characters int,
	isBaseFont bool,
	vVector AdobeFontMetricsVector,
	isFixedV bool,
	capHeight float64,
	xHeight float64,
	ascender float64,
	descender float64,
	underlinePosition float64,
	underlineThickness float64,
	italicAngle float64,
	characterWidth AdobeFontMetricsCharacterSize,
	horizontalStemWidth float64,
	verticalStemWidth float64,
	characterMetrics map[string]AdobeFontMetricsIndividualCharacterMetric,
) AdobeFontMetrics {
	return AdobeFontMetrics{
		AfmVersion:         afmVersion,
		Comments:           comments,
		MetricSets:         metricSets,
		FontName:           fontName,
		FullName:           fullName,
		FamilyName:         familyName,
		Weight:             weight,
		BoundingBox:        boundingBox,
		Version:            version,
		Notice:             notice,
		EncodingScheme:     encodingScheme,
		MappingScheme:      mappingScheme,
		EscapeCharacter:    escapeCharacter,
		CharacterSet:       characterSet,
		Characters:         characters,
		IsBaseFont:         isBaseFont,
		VVector:            vVector,
		IsFixedV:           isFixedV,
		CapHeight:          capHeight,
		XHeight:            xHeight,
		Ascender:           ascender,
		Descender:          descender,
		UnderlinePosition:  underlinePosition,
		UnderlineThickness: underlineThickness,
		ItalicAngle:        italicAngle,
		CharacterWidth:     characterWidth,
		HorizontalStemWidth: horizontalStemWidth,
		VerticalStemWidth:  verticalStemWidth,
		CharacterMetrics:   characterMetrics,
	}
}

// String returns a string representation of the font metrics.
func (a AdobeFontMetrics) String() string {
	name := a.FullName
	if a.FontName != "" {
		name = a.FontName
	}
	return fmt.Sprintf("AFM Font %s with %d characters.", name, a.Characters)
}
