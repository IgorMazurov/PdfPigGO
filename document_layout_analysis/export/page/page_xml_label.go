package page

// PageXmlLabel represents a semantic label element in PAGE XML documents,
// carrying a tag value (e.g., 'person'), optional type information
// (e.g., 'YYYY-mm-dd' for a date label), and optional comments.
type PageXmlLabel struct {
	// Value is the label or tag (e.g., 'person'). Can be an RDF resource identifier.
	Value string `xml:"value,attr"`

	// Type is additional information on the label (e.g., 'YYYY-mm-dd' for a date).
	Type string `xml:"type,attr,omitempty"`

	// Comments contains optional annotation text.
	Comments string `xml:"comments,attr,omitempty"`
}
