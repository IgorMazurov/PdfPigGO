package page

// PageXmlGrapheme represents a sub-element of a glyph — the smallest graphical
// unit that can be assigned a Unicode code point, in PAGE XML documents.
type PageXmlGrapheme struct {
	PageXmlGraphemeBase

	// Coords defines the polygon outline of the grapheme.
	Coords *PageXmlCoords `xml:"Coords,omitempty"`
}
