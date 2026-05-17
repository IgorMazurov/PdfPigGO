package page

// PageXmlGraphemeBase represents the base type for graphemes, grapheme groups,
// and non-printing characters in PAGE XML documents.
type PageXmlGraphemeBase struct {
	// TextEquivs store text content in various encodings and formats.
	TextEquivs []PageXmlTextEquiv `xml:"TextEquiv"`

	// Id is the unique identifier.
	Id string `xml:"id,attr,omitempty"`

	// Index defines the order within the parent container.
	Index *int `xml:"index,attr,omitempty"`

	// Ligature indicates whether this element represents a ligature.
	Ligature *bool `xml:"ligature,attr,omitempty"`

	// CharType specifies the character type (base or combining).
	CharType *PageXmlGraphemeBaseCharType `xml:"charType,attr,omitempty"`

	// Custom is a free-form attribute for custom data.
	Custom string `xml:"custom,attr,omitempty"`

	// Comments contains optional annotation text.
	Comments string `xml:"comments,attr,omitempty"`
}
