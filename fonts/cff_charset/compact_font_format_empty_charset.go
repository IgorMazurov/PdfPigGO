package cffcharset

// CompactFontFormatEmptyCharset is an empty charset for CID fonts which map from
// Character ID to Glyph Id without using strings.
type CompactFontFormatEmptyCharset struct {
	numberOfCharstrings int
}

// NewCompactFontFormatEmptyCharset creates a new empty charset with the given
// number of charstrings.
func NewCompactFontFormatEmptyCharset(numberOfCharstrings int) *CompactFontFormatEmptyCharset {
	return &CompactFontFormatEmptyCharset{
		numberOfCharstrings: numberOfCharstrings,
	}
}

// IsCidCharset always returns true for this charset.
func (c *CompactFontFormatEmptyCharset) IsCidCharset() bool {
	return true
}

// GetNameByGlyphId panics because CID charsets do not support named glyphs.
func (c *CompactFontFormatEmptyCharset) GetNameByGlyphId(glyphId int) string {
	panic("CID charsets do not support named glyphs")
}

// GetNameByStringId panics because CID charsets do not support named glyphs.
func (c *CompactFontFormatEmptyCharset) GetNameByStringId(stringId int) string {
	panic("CID charsets do not support named glyphs")
}

// GetStringIdByGlyphId panics because CID charsets do not support named glyphs.
func (c *CompactFontFormatEmptyCharset) GetStringIdByGlyphId(glyphId int) int {
	panic("CID charsets do not support named glyphs")
}

// GetGlyphIdByName returns 0 as CID charsets do not use named glyph lookups.
func (c *CompactFontFormatEmptyCharset) GetGlyphIdByName(characterName string) int {
	return 0
}

var _ CompactFontFormatCharset = (*CompactFontFormatEmptyCharset)(nil)
