package page

// PageXmlTextEquiv represents text equivalence data in PAGE XML documents,
// storing text content in various encodings and formats.
type PageXmlTextEquiv struct {
	// PlainText is text in a simple form (ASCII or extended ASCII as mostly used for typing).
	// No use of special characters for ligatures — they should be stored as two separate characters.
	PlainText string `xml:"PlainText"`

	// Unicode is the correct encoding of the original, always using the corresponding Unicode code point.
	// Ligatures must be represented as one character etc.
	Unicode string `xml:"Unicode"`

	// Index is used for sort order when multiple TextEquivs are defined.
	// The text content with the lowest index should be interpreted as the main text content.
	Index string `xml:"index,attr,omitempty"`

	// Conf is the confidence value of the text equivalence.
	Conf *float32 `xml:"conf,attr,omitempty"`

	// DataType describes the type of text content (free text, number, etc.).
	// This is only a descriptive attribute; the text type is not checked during XML validation.
	DataType *PageXmlTextDataSimpleType `xml:"dataType,attr,omitempty"`

	// DataTypeDetails provides refinement for the dataType attribute, e.g., a regular expression.
	DataTypeDetails string `xml:"dataTypeDetails,attr,omitempty"`

	// Comments contains free-form comments about this text equivalence.
	Comments string `xml:"comments,attr,omitempty"`
}
