package alto

// AltoIllustration represents a picture or image in ALTO format.
type AltoIllustration struct {
	AltoBlock

	// IllustrationType is a user defined string to identify the type of illustration like photo, map, drawing, chart, etc.
	IllustrationType string `xml:"TYPE,attr,omitempty"`

	// FileId is a link to an image which contains only the illustration.
	FileId string `xml:"FILEID,attr,omitempty"`
}
