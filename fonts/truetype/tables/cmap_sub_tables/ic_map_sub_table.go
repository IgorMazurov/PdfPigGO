package cmapsubtables

// ICMapSubTable represents a single CMap sub-table in a TrueType font.
// The CMap table maps character codes to glyph indices. A font that runs on
// multiple platforms has multiple encoding tables stored as sub-tables.
type ICMapSubTable interface {
	// PlatformId returns the platform identifier for this CMap sub-table.
	PlatformId() TrueTypeCMapPlatform

	// EncodingId returns the platform-specific encoding identifier.
	// Interpretation depends on the value of PlatformId.
	EncodingId() uint16

	// CharacterCodeToGlyphIndex maps a character code to the array index
	// of the glyph in the font data.
	CharacterCodeToGlyphIndex(characterCode int) int
}
