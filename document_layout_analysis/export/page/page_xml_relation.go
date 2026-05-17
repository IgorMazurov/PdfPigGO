package page

// PageXmlRelation represents a one-to-one relation between two layout objects in PAGE XML.
// Use 'link' for loose relations and 'join' for strong relations (e.g., fragmented content).
type PageXmlRelation struct {
	// Id is the unique identifier of this relation element.
	Id string `xml:"id,attr"`

	// Type specifies whether the relation is a link or join.
	Type *PageXmlRelationType `xml:"type,attr,omitempty"`

	// SourceRegionRef is a reference to the source region of the relation.
	SourceRegionRef *PageXmlRegionRef `xml:"sourceRegionRef"`

	// TargetRegionRef is a reference to the target region of the relation.
	TargetRegionRef *PageXmlRegionRef `xml:"targetRegionRef"`

	// Labels contains semantic labels/tags for this relation.
	Labels []PageXmlLabels `xml:"Labels"`

	// Custom is a generic attribute for arbitrary use.
	Custom string `xml:"custom,attr,omitempty"`

	// Comments contains optional annotation text.
	Comments string `xml:"comments,attr,omitempty"`
}
