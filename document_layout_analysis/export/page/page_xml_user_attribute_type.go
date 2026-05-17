package page

// PageXmlUserAttributeType represents the type of a user-defined attribute value.
type PageXmlUserAttributeType byte

const (
	// XsdString indicates the attribute value is an xsd:string.
	XsdString PageXmlUserAttributeType = iota

	// XsdInteger indicates the attribute value is an xsd:integer.
	XsdInteger

	// XsdBoolean indicates the attribute value is an xsd:boolean.
	XsdBoolean

	// XsdFloat indicates the attribute value is an xsd:float.
	XsdFloat
)
