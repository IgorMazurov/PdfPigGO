package cmapsubtables

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// TrimmedTableMappingCMapTable is a format 6 CMap sub-table that uses 2 bytes to map
// a contiguous range of character codes to glyph indices.
type TrimmedTableMappingCMapTable struct {
	platformID       TrueTypeCMapPlatform
	encodingID       uint16
	firstCharacterCode int
	entryCount       int
	glyphIndices     []uint16
}

const (
	trimmedFormat            = uint16(6)
	trimmedDefaultLanguageID = uint16(0)
	sizeOfShort              = 2
)

// NewTrimmedTableMappingCMapTable creates a new TrimmedTableMappingCMapTable.
func NewTrimmedTableMappingCMapTable(platformID TrueTypeCMapPlatform, encodingID uint16, firstCharacterCode int, entryCount int, glyphIndices []uint16) *TrimmedTableMappingCMapTable {
	return &TrimmedTableMappingCMapTable{
		platformID:       platformID,
		encodingID:       encodingID,
		firstCharacterCode: firstCharacterCode,
		entryCount:       entryCount,
		glyphIndices:     glyphIndices,
	}
}

// PlatformId returns the platform identifier.
func (t *TrimmedTableMappingCMapTable) PlatformId() TrueTypeCMapPlatform {
	return t.platformID
}

// EncodingId returns the encoding identifier.
func (t *TrimmedTableMappingCMapTable) EncodingId() uint16 {
	return t.encodingID
}

// FirstCharacterCode returns the first character code in the contiguous range.
func (t *TrimmedTableMappingCMapTable) FirstCharacterCode() int {
	return t.firstCharacterCode
}

// LastCharacterCode returns the last character code in the contiguous range.
func (t *TrimmedTableMappingCMapTable) LastCharacterCode() int {
	return t.firstCharacterCode + t.entryCount - 1
}

// CharacterCodeToGlyphIndex maps a character code to a glyph index.
// Returns 0 if the character code is outside the mapped range.
func (t *TrimmedTableMappingCMapTable) CharacterCodeToGlyphIndex(characterCode int) int {
	if characterCode < t.firstCharacterCode || characterCode > t.firstCharacterCode+t.entryCount {
		return 0
	}

	offset := characterCode - t.firstCharacterCode

	if offset < 0 || offset >= len(t.glyphIndices) {
		return 0
	}

	return int(t.glyphIndices[offset])
}

// Write serializes the CMap sub-table to the given writer.
func (t *TrimmedTableMappingCMapTable) Write(w io.Writer) error {
	if _, err := core.WriteUShort(w, trimmedFormat); err != nil {
		return err
	}

	length := uint16(5*sizeOfShort + sizeOfShort*len(t.glyphIndices))

	if _, err := core.WriteUShort(w, length); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, trimmedDefaultLanguageID); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, uint16(t.firstCharacterCode)); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, uint16(len(t.glyphIndices))); err != nil {
		return err
	}

	for i := range t.glyphIndices {
		if _, err := core.WriteUShort(w, t.glyphIndices[i]); err != nil {
			return err
		}
	}

	return nil
}
