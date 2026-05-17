package alto

// AltoSP represents a white space in ALTO format.
type AltoSP struct {
	AltoPositionedElement
	Id string `xml:"ID,attr"`
}
