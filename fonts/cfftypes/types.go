package cfftypes

import "github.com/uglytoad/pdfpig/go/core"

// Type2GlyphResult is the result of generating a glyph path from a Type 2 CharString.
type Type2GlyphResult interface {
	Path() []*core.PdfSubpath
	Width() *float64
}

// Type2CharStringsProvider generates glyph paths from character names.
type Type2CharStringsProvider interface {
	Generate(name string, defaultWidthX, nominalWidthX float64) (glyph Type2GlyphResult, err error)
}

// FdSelect maps glyph IDs to font dictionary indices in a CID-keyed CFF font.
type FdSelect interface {
	GetFontDictionaryIndex(glyphId int) int
}

// CompactFontFormatIndex holds an index of byte entries in a Compact Font Format file.
type CompactFontFormatIndex struct {
	bytes [][]byte
}

// NewCompactFontFormatIndex creates a new CompactFontFormatIndex from the given bytes.
func NewCompactFontFormatIndex(bytes [][]byte) *CompactFontFormatIndex {
	if bytes == nil {
		bytes = [][]byte{}
	}
	return &CompactFontFormatIndex{bytes: bytes}
}

// Count returns the number of entries in the index.
func (i *CompactFontFormatIndex) Count() int {
	return len(i.bytes)
}

// Get returns the entry at the specified index.
func (i *CompactFontFormatIndex) Get(index int) []byte {
	return i.bytes[index]
}

// GetBytes returns all entries in the index as a slice of byte slices.
func (i *CompactFontFormatIndex) GetBytes() [][]byte {
	return i.bytes
}

// CompactFontFormatSubroutinesSelector provides access to global and local subroutines
// for resolving subroutine calls during charstring interpretation.
type CompactFontFormatSubroutinesSelector struct {
	GlobalSubroutines    *CompactFontFormatIndex
	LocalSubroutines     *CompactFontFormatIndex
	FdSelect             FdSelect
	FontLocalSubroutines []*CompactFontFormatIndex
}

// GetSubroutines returns the global and local subroutine indices for the given glyph ID.
func (s CompactFontFormatSubroutinesSelector) GetSubroutines(glyphId int) (*CompactFontFormatIndex, *CompactFontFormatIndex) {
	if s.FdSelect == nil || len(s.FontLocalSubroutines) == 0 {
		return s.GlobalSubroutines, s.LocalSubroutines
	}

	fdIdx := s.FdSelect.GetFontDictionaryIndex(glyphId)
	if fdIdx < 0 || fdIdx >= len(s.FontLocalSubroutines) {
		return s.GlobalSubroutines, s.LocalSubroutines
	}

	localPerFont := s.FontLocalSubroutines[fdIdx]
	if localPerFont == nil {
		localPerFont = s.LocalSubroutines
	}

	return s.GlobalSubroutines, localPerFont
}
