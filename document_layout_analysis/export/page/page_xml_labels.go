package page

// PageXmlLabels represents the Labels element in PAGE XML documents,
// containing semantic labels/tags with optional external model references,
// prefixes, and comments.
type PageXmlLabels struct {
	// Labels is a collection of semantic label/tag elements.
	Labels []PageXmlLabel `xml:"Label"`

	// ExternalModel is a reference to an external model/ontology/schema.
	ExternalModel string `xml:"externalModel,attr,omitempty"`

	// ExternalId is an RDF resource identifier (subject or object of an RDF triple).
	ExternalId string `xml:"externalId,attr,omitempty"`

	// Prefix is a prefix for all labels (e.g., first part of a URI).
	Prefix string `xml:"prefix,attr,omitempty"`

	// Comments contains optional annotation text.
	Comments string `xml:"comments,attr,omitempty"`
}
