package page

// PageXmlUserAttribute represents structured custom data defined by name, type and value.
type PageXmlUserAttribute struct {
	// Name is the attribute identifier.
	Name string `xml:"name,attr,omitempty"`

	// Description provides a human-readable explanation of the attribute.
	Description string `xml:"description,attr,omitempty"`

	// Type indicates the data type of the value (string, integer, boolean, float).
	Type *PageXmlUserAttributeType `xml:"type,attr,omitempty"`

	// Value is the attribute's data as a string.
	Value string `xml:"value,attr,omitempty"`
}
