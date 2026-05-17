package page

// PageXmlNonPrintingChar represents a glyph component without visual representation
// but with Unicode code point. Non-visual / non-printing / control character.
// Part of grapheme container (of glyph) or grapheme sub group.
type PageXmlNonPrintingChar struct {
	PageXmlGraphemeBase
}
