package page

// PageXmlRegionRefIndexed represents a numbered region reference in PAGE XML documents,
// containing both an IDREF-style cross-reference identifier and its positional order number.
type PageXmlRegionRefIndexed struct {
	// Index is the position (order number) of this item within the current hierarchy level.
	Index int `xml:"index,attr"`

	// RegionRef is the identifier referencing another region element.
	RegionRef string `xml:"regionRef,attr"`
}
