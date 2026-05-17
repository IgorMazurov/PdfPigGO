package page

import "strings"

// PageXmlGlyph represents a glyph element in PAGE XML documents,
// containing graphemes and text equivalence data.
type PageXmlGlyph struct {
	// AlternativeImages are alternative glyph images (e.g., black-and-white).
	AlternativeImages []AlternativeImage `xml:"AlternativeImage"`

	// Coords defines the polygon outline of the glyph.
	Coords *PageXmlCoords `xml:"Coords,omitempty"`

	// Graphemes are the grapheme elements contained within this glyph.
	Graphemes []PageXmlGraphemeBase `xml:"Grapheme"`

	// TextEquivs store text content in various encodings and formats.
	TextEquivs []PageXmlTextEquiv `xml:"TextEquiv"`

	// TextStyle contains font family, size, colour, and typographic attributes.
	TextStyle *PageXmlTextStyle `xml:"TextStyle,omitempty"`

	// UserDefined holds structured custom data defined by name, type, and value.
	UserDefined []PageXmlUserAttribute `xml:"UserAttribute"`

	// Labels contain semantic labels/tags with optional external model references.
	Labels []PageXmlLabels `xml:"Labels"`

	// Id is the unique identifier for this glyph.
	Id string `xml:"id,attr,omitempty"`

	// Ligature indicates whether this glyph represents a ligature.
	Ligature *bool `xml:"ligature,attr,omitempty"`

	// Symbol indicates whether this glyph represents a symbol.
	Symbol *bool `xml:"symbol,attr,omitempty"`

	// Script indicates the writing script used for the glyph.
	Script *PageXmlScriptSimpleType `xml:"script,attr,omitempty"`

	// Production indicates how the text was produced (printed, handwritten, etc.).
	Production *PageXmlProductionSimpleType `xml:"production,attr,omitempty"`

	// Custom is a free-form attribute for custom data.
	Custom string `xml:"custom,attr,omitempty"`

	// Comments contains optional annotation text about this glyph.
	Comments string `xml:"comments,attr,omitempty"`
}

// String returns the concatenated Unicode text from all text equivalents,
// separated by newlines.
func (g *PageXmlGlyph) String() string {
	var result []string
	for _, te := range g.TextEquivs {
		result = append(result, te.Unicode)
	}
	return strings.Join(result, "\n")
}
