package page

// AlternativeImage represents an alternative image reference in PAGE XML documents.
type AlternativeImage struct {
	// FileName is the path or URI to the alternative image file.
	FileName string `xml:"filename,attr"`

	// Comments contains optional annotation text for the image.
	Comments string `xml:"comments,attr"`

	// Conf is a confidence value between 0 and 1.
	Conf *float32 `xml:"conf,attr,omitempty"`
}
