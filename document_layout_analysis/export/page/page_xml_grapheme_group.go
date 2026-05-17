package page

// PageXmlGraphemeGroup represents a group of graphemes within a glyph in PAGE XML documents.
type PageXmlGraphemeGroup struct {
	PageXmlGraphemeBase

	// Items contains the child graphemes and non-printing characters in this group.
	Items []PageXmlGraphemeBase `xml:"Grapheme,NonPrintingChar"`
}
