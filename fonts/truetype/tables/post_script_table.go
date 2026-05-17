package tables

import (
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// PostScriptTable contains information for TrueType fonts on PostScript printers,
// including data for the FontInfo dictionary and the PostScript glyph names.
type PostScriptTable struct {
	directoryTable     truetype.TrueTypeHeaderTable
	format             float32
	italicAngle        float32
	underlinePosition  int16
	underlineThickness int16
	isFixedPitch       uint32
	minMemoryType42    uint32
	maxMemoryType42    uint32
	minMemoryType1     uint32
	maxMemoryType1     uint32
	glyphNames         []string
}

// NewPostScriptTable creates a new PostScriptTable.
func NewPostScriptTable(
	directoryTable truetype.TrueTypeHeaderTable,
	format float32,
	italicAngle float32,
	underlinePosition int16,
	underlineThickness int16,
	isFixedPitch uint32,
	minMemoryType42 uint32,
	maxMemoryType42 uint32,
	minMemoryType1 uint32,
	maxMemoryType1 uint32,
	glyphNames []string,
) PostScriptTable {
	return PostScriptTable{
		directoryTable:     directoryTable,
		format:             format,
		italicAngle:        italicAngle,
		underlinePosition:  underlinePosition,
		underlineThickness: underlineThickness,
		isFixedPitch:       isFixedPitch,
		minMemoryType42:    minMemoryType42,
		maxMemoryType42:    maxMemoryType42,
		minMemoryType1:     minMemoryType1,
		maxMemoryType1:     maxMemoryType1,
		glyphNames:         glyphNames,
	}
}

// Tag returns the 4-letter identifier for this table.
func (p PostScriptTable) Tag() string {
	return truetype.Post
}

// DirectoryTable returns the directory entry from the font's offset subtable.
func (p PostScriptTable) DirectoryTable() truetype.TrueTypeHeaderTable {
	return p.directoryTable
}

// Format returns the format version. 1 = standard Mac, 2 = Microsoft, 2.5 = subset, 3 = no info.
func (p PostScriptTable) Format() float32 { return p.format }

// ItalicAngle returns the angle in counter-clockwise degrees from vertical.
func (p PostScriptTable) ItalicAngle() float32 { return p.italicAngle }

// UnderlinePosition returns the suggested underline position with negative values below baseline.
func (p PostScriptTable) UnderlinePosition() int16 { return p.underlinePosition }

// UnderlineThickness returns the suggested underline thickness.
func (p PostScriptTable) UnderlineThickness() int16 { return p.underlineThickness }

// IsFixedPitch returns 0 if the font is proportionally spaced, non-zero for monospace.
func (p PostScriptTable) IsFixedPitch() uint32 { return p.isFixedPitch }

// MinMemoryType42 returns minimum memory usage when the TrueType font is downloaded.
func (p PostScriptTable) MinMemoryType42() uint32 { return p.minMemoryType42 }

// MaxMemoryType42 returns maximum memory usage when the TrueType font is downloaded.
func (p PostScriptTable) MaxMemoryType42() uint32 { return p.maxMemoryType42 }

// MinMemoryType1 returns minimum memory usage when downloaded as a Type 1 font.
func (p PostScriptTable) MinMemoryType1() uint32 { return p.minMemoryType1 }

// MaxMemoryType1 returns maximum memory usage when downloaded as a Type 1 font.
func (p PostScriptTable) MaxMemoryType1() uint32 { return p.maxMemoryType1 }

// GlyphNames returns the PostScript names of the glyphs.
func (p PostScriptTable) GlyphNames() []string { return p.glyphNames }
