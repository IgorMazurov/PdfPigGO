package page

import "strings"

// PageXmlTextLine represents a text line element in PAGE XML documents,
// containing words, baseline geometry, text equivalents, and styling information.
type PageXmlTextLine struct {
	// AlternativeImages are alternative text line images (e.g., black-and-white).
	AlternativeImages []AlternativeImage `xml:"AlternativeImage"`

	// Coords defines the polygon outline of the text line.
	Coords *PageXmlCoords `xml:"Coords,omitempty"`

	// Baseline marks the multiple connected points that define the baseline of the glyphs.
	Baseline *Baseline `xml:"Baseline,omitempty"`

	// Words are the word elements contained within this text line.
	Words []PageXmlWord `xml:"Word"`

	// TextEquivs store text content in various encodings and formats.
	TextEquivs []PageXmlTextEquiv `xml:"TextEquiv"`

	// TextStyle contains font family, size, colour, and typographic attributes.
	TextStyle *PageXmlTextStyle `xml:"TextStyle,omitempty"`

	// UserDefined holds structured custom data defined by name, type, and value.
	UserDefined []PageXmlUserAttribute `xml:"UserAttribute"`

	// Labels contain semantic labels/tags with optional external model references.
	Labels []PageXmlLabels `xml:"Labels"`

	// Id is the unique identifier for this text line.
	Id string `xml:"id,attr,omitempty"`

	// PrimaryLanguage indicates the main language of the text content.
	PrimaryLanguage *PageXmlLanguageSimpleType `xml:"primaryLanguage,attr,omitempty"`

	// PrimaryScript indicates the primary writing script used in the text.
	PrimaryScript *PageXmlScriptSimpleType `xml:"primaryScript,attr,omitempty"`

	// SecondaryScript indicates a secondary writing script present in the text.
	SecondaryScript *PageXmlScriptSimpleType `xml:"secondaryScript,attr,omitempty"`

	// ReadingDirection specifies the direction in which the text should be read.
	ReadingDirection *PageXmlReadingDirectionSimpleType `xml:"readingDirection,attr,omitempty"`

	// Production indicates how the text was produced (printed, handwritten, etc.).
	Production *PageXmlProductionSimpleType `xml:"production,attr,omitempty"`

	// Custom is a free-form attribute for custom data.
	Custom string `xml:"custom,attr,omitempty"`

	// Comments contains optional annotation text about this text line.
	Comments string `xml:"comments,attr,omitempty"`

	// Index defines the position of this text line within its parent element.
	Index *int `xml:"index,attr,omitempty"`
}

// String returns the concatenated Unicode text from all text equivalents,
// separated by newlines.
func (t *PageXmlTextLine) String() string {
	var result []string
	for _, te := range t.TextEquivs {
		result = append(result, te.Unicode)
	}
	return strings.Join(result, "\n")
}
