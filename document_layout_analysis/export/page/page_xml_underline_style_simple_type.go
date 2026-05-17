package page

// PageXmlUnderlineStyleSimpleType represents underline style values in PAGE XML documents.
type PageXmlUnderlineStyleSimpleType byte

const (
	// SingleLine represents a single-line underline style.
	SingleLine PageXmlUnderlineStyleSimpleType = iota

	// DoubleLine represents a double-line underline style.
	DoubleLine

	// UnderlineStyleOther represents an other/unknown underline style.
	UnderlineStyleOther
)
