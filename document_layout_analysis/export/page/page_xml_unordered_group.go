package page

// PageXmlUnorderedGroup represents a numbered group containing unordered elements
// in PAGE XML documents, with optional caption, type classification, and continuation flag.
type PageXmlUnorderedGroup struct {
	// UserDefined contains custom user-defined attributes attached to this element.
	UserDefined []PageXmlUserAttribute `xml:"UserAttribute"`

	// Labels holds semantic labels/tags associated with this group.
	Labels []PageXmlLabels `xml:"Labels"`

	// Items contains the child elements: ordered groups, region references, or nested unordered groups.
	Items []any `xml:",any"`

	// Id is the unique identifier of this element.
	Id string `xml:"id,attr,omitempty"`

	// RegionRef is an optional link to a parent region of nested regions.
	RegionRef string `xml:"regionRef,attr,omitempty"`

	// Caption provides a human-readable title for this group.
	Caption string `xml:"caption,attr,omitempty"`

	// Type specifies the kind of group (paragraph, list, figure, etc.).
	Type *PageXmlGroupSimpleType `xml:"type,attr,omitempty"`

	// Continuation indicates whether this group is a continuation of another group
	// from a previous column or page.
	Continuation *bool `xml:"continuation,attr,omitempty"`

	// Custom is a free-form attribute for generic use.
	Custom string `xml:"custom,attr,omitempty"`

	// Comments contains optional annotation text.
	Comments string `xml:"comments,attr,omitempty"`
}
