package alto

import "encoding/xml"

// AltoSchemaVersion is the ALTO schema version 4.1.
const AltoSchemaVersion = "4.0"

// AltoNamespace is the XML namespace for ALTO documents.
const AltoNamespace = "http://www.loc.gov/standards/alto/ns-v4#"

// AltoDocument is the root element of an ALTO document (version 4.1).
// See https://github.com/altoxml/schema
type AltoDocument struct {
	XMLName xml.Name `xml:"http://www.loc.gov/standards/alto/ns-v4# alto"`

	// Description contains general settings like measurement units and metadata.
	Description *AltoDescription `xml:"Description"`

	// Styles define properties of layout elements. A style defined in a parent element
	// is used as default style for all related children elements.
	Styles *AltoStyles `xml:"Styles"`

	// Tags define properties of additional characteristics referenced from content
	// elements via TAGREF attribute. Contains LayoutTags, StructureTags, RoleTags,
	// NamedEntityTags and OtherTags.
	Tags *AltoTags `xml:"Tags"`

	// Layout is the root layout element containing pages.
	Layout *AltoLayout `xml:"Layout"`

	// SchemaVersion is the schema version of the ALTO file.
	SchemaVersion string `xml:"SCHEMAVERSION,attr"`
}

// NewAltoDocument creates a new AltoDocument with the default schema version.
func NewAltoDocument() *AltoDocument {
	return &AltoDocument{
		SchemaVersion: AltoSchemaVersion,
		Description:   &AltoDescription{},
		Layout:        &AltoLayout{},
	}
}

// AltoTags is a container for tag elements in ALTO format.
// Available tag types: LayoutTag, StructureTag, RoleTag, NamedEntityTag, OtherTag.
type AltoTags struct {
	// Items holds the individual tag elements.
	Items []AltoTag `xml:",any"`

	// ItemsElementName tracks the element type for each item in Items.
	ItemsElementName []AltoItemsChoice `xml:"ItemsElementName,omitempty"`
}
