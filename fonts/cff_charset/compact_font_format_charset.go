package cffcharset

// CompactFontFormatCharset provides glyph-to-name and name-to-glyph lookups
// for a Compact Font Format (CFF) charset.
type CompactFontFormatCharset interface {
	// IsCidCharset reports whether the charset is a CID (Character Identifier) charset.
	IsCidCharset() bool

	// GetNameByGlyphId returns the character name for the given glyph ID.
	GetNameByGlyphId(glyphId int) string

	// GetNameByStringId returns the character name for the given string ID.
	GetNameByStringId(stringId int) string

	// GetStringIdByGlyphId returns the string ID corresponding to the given glyph ID.
	GetStringIdByGlyphId(glyphId int) int

	// GetGlyphIdByName returns the glyph ID for the given character name.
	GetGlyphIdByName(characterName string) int
}
