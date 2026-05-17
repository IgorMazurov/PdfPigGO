package alto

// AltoFileIdentifier represents a unique identifier for an image file in ALTO format.
type AltoFileIdentifier struct {
	FileIdentifierLocation string `xml:"fileIdentifierLocation,attr"`
	Value                  string `xml:",chardata"`
}
