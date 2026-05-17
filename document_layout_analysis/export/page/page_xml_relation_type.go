package page

// PageXmlRelationType represents the type of relation between elements in PAGE XML.
type PageXmlRelationType byte

const (
	// Link indicates a link relation between elements.
	PageXmlRelationLink PageXmlRelationType = iota

	// Join indicates a join relation between elements.
	PageXmlRelationJoin
)
