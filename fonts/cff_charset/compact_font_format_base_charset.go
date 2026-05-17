package cffcharset

// GlyphIdEntry holds the string ID and name associated with a glyph ID.
type GlyphIdEntry struct {
	StringId int
	Name     string
}

// BaseCharset is the base implementation for CFF charsets, providing glyph-to-name
// and name-to-glyph lookups backed by a dictionary. It corresponds to the abstract
// CompactFontFormatCharset class in the C# source. Subtypes embed this struct or
// implement CompactFontFormatCharset directly with custom lookup logic.
type BaseCharset struct {
	glyphIdToStringIdAndName map[int]GlyphIdEntry
}

// NewBaseCharset creates a new BaseCharset from the given data entries.
// The entry with glyph ID 0 and name ".notdef" is always inserted first.
func NewBaseCharset(data []struct {
	GlyphId  int
	StringId int
	Name     string
}) *BaseCharset {
	dict := make(map[int]GlyphIdEntry, len(data)+1)
	dict[0] = GlyphIdEntry{StringId: 0, Name: ".notdef"}

	for _, t := range data {
		dict[t.GlyphId] = GlyphIdEntry{StringId: t.StringId, Name: t.Name}
	}

	return &BaseCharset{glyphIdToStringIdAndName: dict}
}

// IsCidCharset returns false for standard (non-CID) charsets.
func (b *BaseCharset) IsCidCharset() bool {
	return false
}

// GetNameByGlyphId returns the character name for the given glyph ID,
// or an empty string if the glyph ID is not found.
func (b *BaseCharset) GetNameByGlyphId(glyphId int) string {
	entry, ok := b.glyphIdToStringIdAndName[glyphId]
	if !ok {
		return ""
	}
	return entry.Name
}

// GetNameByStringId returns the character name for the given string ID.
// If multiple entries share the same string ID, one arbitrary match is returned.
// Returns an empty string if no matching entry is found.
func (b *BaseCharset) GetNameByStringId(stringId int) string {
	for _, entry := range b.glyphIdToStringIdAndName {
		if entry.StringId == stringId {
			return entry.Name
		}
	}
	return ""
}

// GetStringIdByGlyphId returns the string ID for the given glyph ID,
// or 0 if the glyph ID is not found.
func (b *BaseCharset) GetStringIdByGlyphId(glyphId int) int {
	entry, ok := b.glyphIdToStringIdAndName[glyphId]
	if !ok {
		return 0
	}
	return entry.StringId
}

// GetGlyphIdByName returns the glyph ID for the given character name using
// ordinal comparison, or 0 if no matching entry is found.
func (b *BaseCharset) GetGlyphIdByName(characterName string) int {
	for glyphId, entry := range b.glyphIdToStringIdAndName {
		if entry.Name == characterName {
			return glyphId
		}
	}
	return 0
}

var _ CompactFontFormatCharset = (*BaseCharset)(nil)
