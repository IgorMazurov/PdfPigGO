package cffcharset

// CompactFontFormatFormat2Charset is a charset for CFF fonts with a large number
// of well-ordered string IDs. It delegates all lookup operations to the embedded BaseCharset.
type CompactFontFormatFormat2Charset struct {
	BaseCharset
}

// NewCompactFontFormatFormat2Charset creates a new Format 2 charset from the given
// glyph ID entries.
func NewCompactFontFormatFormat2Charset(data []struct {
	GlyphId  int
	StringId int
	Name     string
}) *CompactFontFormatFormat2Charset {
	return &CompactFontFormatFormat2Charset{
		BaseCharset: *NewBaseCharset(data),
	}
}

var _ CompactFontFormatCharset = (*CompactFontFormatFormat2Charset)(nil)
