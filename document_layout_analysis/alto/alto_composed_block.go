package alto

// AltoComposedBlock represents a block that consists of other blocks in ALTO format.
type AltoComposedBlock struct {
	AltoBlock

	TypeComposed string `xml:"TYPE,attr,omitempty"`
	FileId       string `xml:"FILEID,attr,omitempty"`
}
