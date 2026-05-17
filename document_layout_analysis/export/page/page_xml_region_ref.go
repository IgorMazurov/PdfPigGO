package page

// PageXmlRegionRef represents a reference to another region in a PAGE XML document,
// using an IDREF-style cross-reference identifier.
type PageXmlRegionRef struct {
	// RegionRef is the identifier referencing another region element.
	RegionRef string `xml:"regionRef,attr"`
}
