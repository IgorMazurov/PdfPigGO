package page

// PageXmlLayer represents a layer element in PAGE XML documents,
// containing references to regions and layer metadata such as z-index and caption.
type PageXmlLayer struct {
	// RegionRefs is the list of region references belonging to this layer.
	RegionRefs []PageXmlRegionRef `xml:"RegionRef"`

	// Id is the unique identifier for this layer element.
	Id string `xml:"id,attr"`

	// ZIndex specifies the stacking order of this layer relative to other layers.
	ZIndex int `xml:"zIndex,attr"`

	// Caption provides a human-readable description of the layer.
	Caption string `xml:"caption,attr"`
}
