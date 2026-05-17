package page

// PageXmlCoords represents coordinate data in PAGE XML documents,
// describing a polygon outline of an element as a path of points.
type PageXmlCoords struct {
	// Points is a string of coordinate pairs defining the polygon geometry.
	Points string `xml:"points,attr"`

	// Conf is a confidence value between 0 and 1.
	Conf *float32 `xml:"conf,attr,omitempty"`
}
