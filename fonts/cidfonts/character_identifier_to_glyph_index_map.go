// Package cidfonts provides types for handling CID (Character Identifier) fonts in PDFs.
package cidfonts

// CharacterIdentifierToGlyphIndexMap specifies mapping from character identifiers
// to glyph indices. Can either be defined as a name in which case it must be
// Identity or a stream which defines the mapping.
type CharacterIdentifierToGlyphIndexMap struct {
	isIdentity bool
	mapData    []int
}

// NewCharacterIdentifierToGlyphIndexMap creates an identity mapping where each
// character identifier maps directly to the same glyph index.
func NewCharacterIdentifierToGlyphIndexMap() *CharacterIdentifierToGlyphIndexMap {
	return &CharacterIdentifierToGlyphIndexMap{
		isIdentity: true,
	}
}

// NewCharacterIdentifierToGlyphIndexMapFromStream creates a mapping from stream bytes.
// Each pair of bytes represents a 16-bit big-endian glyph index.
func NewCharacterIdentifierToGlyphIndexMapFromStream(streamBytes []byte) *CharacterIdentifierToGlyphIndexMap {
	numberOfEntries := len(streamBytes) / 2

	mapData := make([]int, numberOfEntries)

	for i := 0; i < numberOfEntries; i++ {
		glyphIndex := int(streamBytes[i*2])<<8 | int(streamBytes[i*2+1])
		mapData[i] = glyphIndex
	}

	return &CharacterIdentifierToGlyphIndexMap{
		isIdentity: false,
		mapData:    mapData,
	}
}

// GetGlyphIndex returns the glyph index for the given character identifier.
// For identity mappings, the character identifier is returned directly.
// If the character identifier is out of range, 0 is returned.
func (m *CharacterIdentifierToGlyphIndexMap) GetGlyphIndex(characterIdentifier int) int {
	if m.isIdentity {
		return characterIdentifier
	}

	if characterIdentifier >= len(m.mapData) || characterIdentifier < 0 {
		return 0
	}

	return m.mapData[characterIdentifier]
}
