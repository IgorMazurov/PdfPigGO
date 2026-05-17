package alto

// AltoDocumentIdentifier represents a unique identifier for the document in ALTO format.
type AltoDocumentIdentifier struct {
	DocumentIdentifierLocation string `xml:"documentIdentifierLocation,attr"`
	Value                      string `xml:",chardata"`
}
