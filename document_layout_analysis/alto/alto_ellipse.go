package alto

// AltoEllipse represents an ellipse shape in ALTO format.
// HorizontalPosition and VerticalPosition describe the center of the ellipse.
// HorizontalLength and VerticalLength are the width and height.
type AltoEllipse struct {
	HorizontalPosition float32  `xml:"HPOS,attr"`
	VerticalPosition   float32  `xml:"VPOS,attr"`
	HorizontalLength   float32  `xml:"HLENGTH,attr"`
	VerticalLength     float32  `xml:"VLENGTH,attr"`
	Rotation           *float32 `xml:"ROTATION,attr"`
}
