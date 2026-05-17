package page

// PageXmlGridPoints represents grid points with x,y coordinates in PAGE XML documents.
type PageXmlGridPoints struct {
	// Index is the grid row index.
	Index int `xml:"index,attr"`

	// Points is a string of coordinate pairs defining the grid point positions.
	Points string `xml:"points,attr"`
}
