package page

// PageXmlBorder represents the border of an actual page in PAGE XML documents,
// used when the scanned image contains parts not belonging to the page.
type PageXmlBorder struct {
	// Coords defines the polygon outline of the page border.
	Coords *PageXmlCoords `xml:"coords,omitempty"`
}
