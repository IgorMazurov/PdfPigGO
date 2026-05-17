package cffcharset

// CompactFontFormatFormat1Charset is a charset for CFF fonts with well-ordered
// string IDs. It delegates all lookup operations to the embedded BaseCharset.
type CompactFontFormatFormat1Charset struct {
	BaseCharset
}

// NewCompactFontFormatFormat1Charset creates a new Format 1 charset from the given
// glyph ID entries.
func NewCompactFontFormatFormat1Charset(data []struct {
	GlyphId  int
	StringId int
	Name     string
}) *CompactFontFormatFormat1Charset {
	return &CompactFontFormatFormat1Charset{
		BaseCharset: *NewBaseCharset(data),
	}
}

var _ CompactFontFormatCharset = (*CompactFontFormatFormat1Charset)(nil)
