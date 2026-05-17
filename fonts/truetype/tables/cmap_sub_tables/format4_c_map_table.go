package cmapsubtables

// Format4CMapTable is a format 4 CMap sub-table that defines gappy ranges of
// character code to glyph index mappings. It uses segment coverage tables for
// efficient binary search-based lookup.
type Format4CMapTable struct {
	platformID TrueTypeCMapPlatform
	encodingID uint16
	language   uint16
	segments   []Segment
	glyphIDs   []uint16
}

// NewFormat4CMapTable creates a new Format4CMapTable.
func NewFormat4CMapTable(platformID TrueTypeCMapPlatform, encodingID, language uint16, segments []Segment, glyphIDs []uint16) *Format4CMapTable {
	return &Format4CMapTable{
		platformID: platformID,
		encodingID: encodingID,
		language:   language,
		segments:   segments,
		glyphIDs:   glyphIDs,
	}
}

// PlatformId returns the platform identifier.
func (t *Format4CMapTable) PlatformId() TrueTypeCMapPlatform {
	return t.platformID
}

// EncodingId returns the encoding identifier.
func (t *Format4CMapTable) EncodingId() uint16 {
	return t.encodingID
}

// Language returns the language identifier.
func (t *Format4CMapTable) Language() uint16 {
	return t.language
}

// Segments returns the list of mapping segments.
func (t *Format4CMapTable) Segments() []Segment {
	return t.segments
}

// GlyphIds returns the glyph ID array used for non-zero offset lookups.
func (t *Format4CMapTable) GlyphIds() []uint16 {
	return t.glyphIDs
}

// CharacterCodeToGlyphIndex maps a character code to a glyph index using
// segment coverage tables. Returns 0 if the character code is not found in any
// segment.
func (t *Format4CMapTable) CharacterCodeToGlyphIndex(characterCode int) int {
	for i, seg := range t.segments {
		if characterCode < seg.StartCode || characterCode > seg.EndCode {
			continue
		}

		if seg.IdRangeOffset == 0 {
			return int(uint32(characterCode+seg.IdDelta)&0xFFFF)
		}

		offset := seg.IdRangeOffset/2 + (characterCode - seg.StartCode)
		return int(t.glyphIDs[offset-len(t.segments)+i])
	}

	return 0
}

// Segment represents a contiguous range of character codes mapped to glyph indices
// within a Format4CMapTable.
type Segment struct {
	StartCode     int
	EndCode       int
	IdDelta       int
	IdRangeOffset int
}
