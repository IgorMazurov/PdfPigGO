package alto

// AltoCircle represents a circle shape in ALTO format.
// HorizontalPosition and VerticalPosition describe the center of the circle.
type AltoCircle struct {
	HorizontalPosition float32 `xml:"HPOS,attr"`
	VerticalPosition   float32 `xml:"VPOS,attr"`
	Radius             float32 `xml:"RADIUS,attr"`
}
