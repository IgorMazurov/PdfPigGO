package page

// PageXmlGroupSimpleType represents the type of group in a PAGE XML document.
type PageXmlGroupSimpleType byte

const (
	// Paragraph indicates a paragraph group.
	Paragraph PageXmlGroupSimpleType = iota

	// List indicates a list group.
	List

	// ListItem indicates a list item group.
	ListItem

	// Figure indicates a figure group.
	Figure

	// Article indicates an article group.
	Article

	// Div indicates a division group.
	Div

	// OtherGroup indicates an unrecognized group type.
	OtherGroup
)
