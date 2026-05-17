package cff

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
)

// UnsetOffset indicates that an offset has not yet been assigned.
const UnsetOffset = -1

// CompactFontFormatCharStringType defines the format of the CharString data contained
// within a Compact Font Format font.
type CompactFontFormatCharStringType byte

const (
	// CompactFontFormatCharStringType1 is the Type 1 CharString format as defined by
	// the Adobe Type 1 Font Format.
	CompactFontFormatCharStringType1 CompactFontFormatCharStringType = 1
	// CompactFontFormatCharStringType2 is the Type 2 CharString format as defined by
	// Adobe Technical Note #5177. This is the default type.
	CompactFontFormatCharStringType2 CompactFontFormatCharStringType = 2
)

// CompactFontFormatTopLevelDictionary holds the top-level dictionary for a Compact
// Font Format (CFF) font, containing metadata and offsets to sub-tables such as
// the charset, encoding, charstrings, and private dictionary.
type CompactFontFormatTopLevelDictionary struct {
	Version                   string
	Notice                    string
	Copyright                 string
	FullName                  string
	FamilyName                string
	Weight                    string
	IsFixedPitch              bool
	ItalicAngle               float64
	UnderlinePosition         float64
	UnderlineThickness        float64
	PaintType                 float64
	CharStringType            CompactFontFormatCharStringType
	FontMatrix                *core.TransformationMatrix
	StrokeWidth               float64
	UniqueId                  float64
	FontBoundingBox           core.PdfRectangle
	Xuid                      []float64
	CharSetOffset             int
	EncodingOffset            int
	PrivateDictionaryLocation *SizeAndOffset
	CharStringsOffset         int
	SyntheticBaseFontIndex    int
	PostScript                string
	BaseFontName              string
	BaseFontBlend             []float64
	IsCidFont                 bool
	CidFontOperators          CidFontOperators
}

// NewCompactFontFormatTopLevelDictionary creates a top-level dictionary with the
// default values matching the C# source.
func NewCompactFontFormatTopLevelDictionary() *CompactFontFormatTopLevelDictionary {
	return &CompactFontFormatTopLevelDictionary{
		UnderlinePosition:    -100,
		UnderlineThickness:   50,
		CharStringType:       CompactFontFormatCharStringType2,
		FontBoundingBox:      core.NewPdfRectangleFloat(0, 0, 0, 0),
		CharSetOffset:        UnsetOffset,
		EncodingOffset:       UnsetOffset,
		CharStringsOffset:    -1,
		CidFontOperators:     NewCidFontOperators(),
	}
}

// SizeAndOffset holds a byte size and offset pair for locating a CFF sub-table.
type SizeAndOffset struct {
	Size   int
	Offset int
}

// String returns a human-readable representation of the size and offset.
func (s SizeAndOffset) String() string {
	return fmt.Sprintf("Size: %d, Offset: %d", s.Size, s.Offset)
}

// CidFontOperators holds the CID font-specific operators from the top-level dictionary.
type CidFontOperators struct {
	Ros                  RegistryOrderingSupplement
	Version              int
	Revision             int
	Type                 int
	Count                int
	UidBase              float64
	FontDictionaryArray  int
	FontDictionarySelect int
	FontName             string
}

// NewCidFontOperators creates a CidFontOperators with the default values matching
// the C# source.
func NewCidFontOperators() CidFontOperators {
	return CidFontOperators{
		Version:  0,
		Revision: 0,
		Type:     0,
		Count:    8720,
	}
}

// RegistryOrderingSupplement identifies a character collection by registry string,
// ordering string, and supplement number.
type RegistryOrderingSupplement struct {
	Registry   string
	Ordering   string
	Supplement float64
}
