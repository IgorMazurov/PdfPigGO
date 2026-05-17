package cffcharset

// CompactFontFormatFormat0Charset is a charset for CFF fonts with relatively
// unordered string IDs. It delegates all lookup operations to the embedded BaseCharset.
type CompactFontFormatFormat0Charset struct {
	BaseCharset
}

// NewCompactFontFormatFormat0Charset creates a new Format 0 charset from the given
// glyph ID entries.
func NewCompactFontFormatFormat0Charset(data []struct {
	GlyphId  int
	StringId int
	Name     string
}) *CompactFontFormatFormat0Charset {
	return &CompactFontFormatFormat0Charset{
		BaseCharset: *NewBaseCharset(data),
	}
}

var _ CompactFontFormatCharset = (*CompactFontFormatFormat0Charset)(nil)
