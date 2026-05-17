package page

import "strings"

// PageXmlWord represents a word element in PAGE XML documents,
// containing glyphs and text equivalence data.
type PageXmlWord struct {
	// AlternativeImages are alternative word images (e.g., black-and-white).
	AlternativeImages []AlternativeImage `xml:"AlternativeImage"`

	// Coords defines the polygon outline of the word.
	Coords *PageXmlCoords `xml:"Coords,omitempty"`

	// Glyphs are the glyph elements contained within this word.
	Glyphs []PageXmlGlyph `xml:"Glyph"`

	// TextEquivs store text content in various encodings and formats.
	TextEquivs []PageXmlTextEquiv `xml:"TextEquiv"`

	// TextStyle contains font family, size, colour, and typographic attributes.
	TextStyle *PageXmlTextStyle `xml:"TextStyle,omitempty"`

	// UserDefined holds structured custom data defined by name, type, and value.
	UserDefined []PageXmlUserAttribute `xml:"UserAttribute"`

	// Labels contain semantic labels/tags with optional external model references.
	Labels []PageXmlLabels `xml:"Labels"`

	// Id is the unique identifier for this word.
	Id string `xml:"id,attr,omitempty"`

	// Language overrides the primaryLanguage attribute of parent line/text region.
	Language *PageXmlLanguageSimpleType `xml:"language,attr,omitempty"`

	// PrimaryScript indicates the primary writing script used in the word.
	PrimaryScript *PageXmlScriptSimpleType `xml:"primaryScript,attr,omitempty"`

	// SecondaryScript indicates a secondary writing script present in the word.
	SecondaryScript *PageXmlScriptSimpleType `xml:"secondaryScript,attr,omitempty"`

	// ReadingDirection specifies the direction in which the text should be read.
	ReadingDirection *PageXmlReadingDirectionSimpleType `xml:"readingDirection,attr,omitempty"`

	// Production indicates how the text was produced (printed, handwritten, etc.).
	Production *PageXmlProductionSimpleType `xml:"production,attr,omitempty"`

	// Custom is a free-form attribute for custom data.
	Custom string `xml:"custom,attr,omitempty"`

	// Comments contains optional annotation text about this word.
	Comments string `xml:"comments,attr,omitempty"`
}

// String returns the concatenated Unicode text from all text equivalents,
// separated by newlines.
func (w *PageXmlWord) String() string {
	var result []string
	for _, te := range w.TextEquivs {
		result = append(result, te.Unicode)
	}
	return strings.Join(result, "\n")
}
