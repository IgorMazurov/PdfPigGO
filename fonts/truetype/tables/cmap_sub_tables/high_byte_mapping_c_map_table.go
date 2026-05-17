package cmapsubtables

// HighByteMappingCMapTable is a format 2 CMap sub-table for Chinese, Japanese and
// Korean characters containing mixed 8/16 bit encodings.
type HighByteMappingCMapTable struct {
	platformID              TrueTypeCMapPlatform
	encodingID              uint16
	characterCodesToGlyphIndices map[int]int
}

// NewHighByteMappingCMapTable creates a new HighByteMappingCMapTable.
func NewHighByteMappingCMapTable(platformID TrueTypeCMapPlatform, encodingID uint16, characterCodesToGlyphIndices map[int]int) *HighByteMappingCMapTable {
	return &HighByteMappingCMapTable{
		platformID:                 platformID,
		encodingID:                 encodingID,
		characterCodesToGlyphIndices: characterCodesToGlyphIndices,
	}
}

// PlatformId returns the platform identifier.
func (t *HighByteMappingCMapTable) PlatformId() TrueTypeCMapPlatform {
	return t.platformID
}

// EncodingId returns the encoding identifier.
func (t *HighByteMappingCMapTable) EncodingId() uint16 {
	return t.encodingID
}

// CharacterCodeToGlyphIndex maps a character code to a glyph index.
// Returns 0 if the character code is not found.
func (t *HighByteMappingCMapTable) CharacterCodeToGlyphIndex(characterCode int) int {
	if idx, ok := t.characterCodesToGlyphIndices[characterCode]; ok {
		return idx
	}

	return 0
}

// SubHeader holds the mapping parameters for a single high-byte sub-header block.
type SubHeader struct {
	FirstCode     int
	EntryCount    int
	IdDelta       int16
	IdRangeOffset int
}
