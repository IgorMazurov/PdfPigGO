package cmapsubtables

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// ByteEncodingCMapTable is a format 0 CMap sub-table where character codes and glyph
// indices are restricted to a single byte.
type ByteEncodingCMapTable struct {
	platformID   TrueTypeCMapPlatform
	encodingID   uint16
	languageID   uint16
	glyphMapping []byte
}

const (
	byteFormat              = 0
	byteDefaultLanguageID   = uint16(0)
	byteSizeOfShort         = 2
	byteGlyphMappingLength  = 256
)

// NewByteEncodingCMapTable creates a new ByteEncodingCMapTable.
func NewByteEncodingCMapTable(platformID TrueTypeCMapPlatform, encodingID, languageID uint16, glyphMapping []byte) *ByteEncodingCMapTable {
	return &ByteEncodingCMapTable{
		platformID:   platformID,
		encodingID:   encodingID,
		languageID:   languageID,
		glyphMapping: glyphMapping,
	}
}

// PlatformId returns the platform identifier.
func (t *ByteEncodingCMapTable) PlatformId() TrueTypeCMapPlatform {
	return t.platformID
}

// EncodingId returns the encoding identifier.
func (t *ByteEncodingCMapTable) EncodingId() uint16 {
	return t.encodingID
}

// LanguageId returns the language identifier.
func (t *ByteEncodingCMapTable) LanguageId() uint16 {
	return t.languageID
}

// CharacterCodeToGlyphIndex maps a character code to a glyph index.
func (t *ByteEncodingCMapTable) CharacterCodeToGlyphIndex(characterCode int) int {
	if characterCode < 0 || characterCode >= len(t.glyphMapping) {
		return 0
	}

	return int(t.glyphMapping[characterCode])
}

// Write serializes the CMap sub-table to the given writer.
func (t *ByteEncodingCMapTable) Write(w io.Writer) error {
	if _, err := core.WriteUShort(w, byteFormat); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, byteGlyphMappingLength+byteSizeOfShort*3); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, byteDefaultLanguageID); err != nil {
		return err
	}

	for i := range t.glyphMapping {
		if _, err := w.Write([]byte{t.glyphMapping[i]}); err != nil {
			return err
		}
	}

	return nil
}
