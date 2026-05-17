package page

// PageXmlLayers represents the Layers element in PAGE XML documents,
// expressing the z-index ordering of overlapping regions. An element with a
// greater z-index is always in front of another element with lower z-index.
type PageXmlLayers struct {
	// Layers is the collection of layer elements defining region grouping and stacking order.
	Layers []PageXmlLayer `xml:"Layer"`
}
