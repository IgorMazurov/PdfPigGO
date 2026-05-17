package alto

// AltoPositionedElement is the base element for positioned ALTO elements.
type AltoPositionedElement struct {
	Height               *float32 `xml:"HEIGHT,attr"`
	Width                *float32 `xml:"WIDTH,attr"`
	HorizontalPosition   *float32 `xml:"HPOS,attr"`
	VerticalPosition     *float32 `xml:"VPOS,attr"`
}
