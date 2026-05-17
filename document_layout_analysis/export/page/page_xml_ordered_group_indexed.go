package page

// PageXmlOrderedGroupIndexed represents an ordered group in PAGE XML documents,
// containing child items (groups or region references) with explicit ordering,
// optional caption, type classification, and continuation flag.
type PageXmlOrderedGroupIndexed struct {
	// UserDefined contains custom user-defined attributes attached to this element.
	UserDefined []PageXmlUserAttribute `xml:"UserAttribute"`

	// Labels holds semantic labels/tags associated with this group.
	Labels []PageXmlLabels `xml:"Labels"`

	// Items contains the child elements: ordered groups, unordered groups, or region references.
	Items []any `xml:",any"`

	// Id is the unique identifier of this element.
	Id string `xml:"id,attr,omitempty"`

	// RegionRef is a reference to another region element by its ID.
	RegionRef string `xml:"regionRef,attr,omitempty"`

	// Index defines the order position within the parent container.
	Index int `xml:"index,attr,omitempty"`

	// Caption provides a human-readable title for this group.
	Caption string `xml:"caption,attr,omitempty"`

	// Type specifies the kind of group (paragraph, list, figure, etc.).
	Type *PageXmlGroupSimpleType `xml:"type,attr,omitempty"`

	// Continuation indicates whether this group is a continuation of a previous one.
	Continuation *bool `xml:"continuation,attr,omitempty"`

	// Custom is a free-form attribute for custom data.
	Custom string `xml:"custom,attr,omitempty"`

	// Comments contains optional annotation text.
	Comments string `xml:"comments,attr,omitempty"`
}
