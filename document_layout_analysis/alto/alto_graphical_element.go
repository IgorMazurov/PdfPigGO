package alto

// AltoGraphicalElement represents a graphic used to separate blocks in ALTO format.
// Usually a line or rectangle.
type AltoGraphicalElement struct {
	AltoBlock

	Title string `xml:"TITLE,attr,omitempty"`
	Type  string `xml:"TYPE,attr,omitempty"`
}
