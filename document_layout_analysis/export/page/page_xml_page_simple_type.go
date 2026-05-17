package page

// PageXmlPageSimpleType represents the type of page in a PAGE XML document.
type PageXmlPageSimpleType byte

const (
	// FrontCover indicates a front cover page.
	FrontCover PageXmlPageSimpleType = iota

	// BackCover indicates a back cover page.
	BackCover

	// Title indicates a title page.
	Title

	// TableOfContents indicates a table of contents page.
	TableOfContents

	// Index indicates an index page.
	Index

	// Content indicates a regular content page.
	Content

	// Blank indicates a blank page.
	Blank

	// PageOther indicates an unrecognized page type.
	PageOther
)
