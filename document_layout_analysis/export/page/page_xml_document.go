package page

import "encoding/xml"

// PageXmlNamespace is the XML namespace for PAGE-XML documents.
const PageXmlNamespace = "http://schema.primaresearch.org/PAGE/gts/pagebased/2019-07-15"

// PageXmlDocument represents the root element (PcGts) of a PAGE XML document
// according to the PAGE scheme version 2019-07-15. It contains metadata and
// exactly one page with its regions, reading order, and other properties.
type PageXmlDocument struct {
	XMLName xml.Name `xml:"http://schema.primaresearch.org/PAGE/gts/pagebased/2019-07-15 PcGts"`

	// Metadata holds creator information, timestamps, comments, and metadata items.
	Metadata *PageXmlMetadata `xml:"Metadata,omitempty"`

	// Page is the single document page containing image info, regions, etc.
	Page *PageXmlPage `xml:"Page,omitempty"`

	// PcGtsId is a unique identifier for this PAGE XML document root element.
	PcGtsId string `xml:"pcGtsId,attr,omitempty"`
}
