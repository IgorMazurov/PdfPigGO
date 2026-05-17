package alto

// AltoSourceImageInformation holds information about the source image.
type AltoSourceImageInformation struct {
	FileName            string                 `xml:"fileName"`
	FileIdentifiers     []AltoFileIdentifier   `xml:"fileIdentifier"`
	DocumentIdentifiers []AltoDocumentIdentifier `xml:"documentIdentifier"`
}
