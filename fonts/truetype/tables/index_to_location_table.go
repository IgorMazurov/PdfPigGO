package tables

import (
	"errors"
	"io"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// IndexToLocationTableEntryFormat indicates the format of glyph offset entries
// stored in the raw TrueType data.
type IndexToLocationTableEntryFormat int16

const (
	// IndexToLocationTableShort means the actual local offset divided by 2 is stored.
	IndexToLocationTableShort IndexToLocationTableEntryFormat = 0
	// IndexToLocationTableLong means the actual local offset is stored.
	IndexToLocationTableLong IndexToLocationTableEntryFormat = 1
)

// IndexToLocationTable stores the offsets to the locations of the glyphs relative
// to the start of the glyph data table.
type IndexToLocationTable struct {
	directoryTable truetype.TrueTypeHeaderTable
	format         IndexToLocationTableEntryFormat
	glyphOffsets   []uint32
}

// NewIndexToLocationTable creates a new IndexToLocationTable.
func NewIndexToLocationTable(
	directoryTable truetype.TrueTypeHeaderTable,
	format IndexToLocationTableEntryFormat,
	glyphOffsets []uint32,
) IndexToLocationTable {
	return IndexToLocationTable{
		directoryTable: directoryTable,
		format:         format,
		glyphOffsets:   glyphOffsets,
	}
}

// Tag returns the 4-letter identifier for this table.
func (i IndexToLocationTable) Tag() string {
	return truetype.Loca
}

// DirectoryTable returns the directory entry from the font's offset subtable.
func (i IndexToLocationTable) DirectoryTable() truetype.TrueTypeHeaderTable {
	return i.directoryTable
}

// Format returns the format of glyph offset entries.
func (i IndexToLocationTable) Format() IndexToLocationTableEntryFormat {
	return i.format
}

// GlyphOffsets returns the glyph offsets relative to the start of the glyph data table.
func (i IndexToLocationTable) GlyphOffsets() []uint32 {
	return i.glyphOffsets
}

// Write serializes the index-to-location table to the output stream.
func (i IndexToLocationTable) Write(w io.Writer) error {
	for idx := 0; idx < len(i.glyphOffsets); idx++ {
		offset := i.glyphOffsets[idx]
		switch i.format {
		case IndexToLocationTableShort:
			if _, err := core.WriteUShort(w, uint16(offset/2)); err != nil {
				return err
			}
		case IndexToLocationTableLong:
			if _, err := core.WriteUInt(w, offset); err != nil {
				return err
			}
		default:
			return errors.New("invalid format for index-to-location table")
		}
	}

	return nil
}

// Compile-time interface assertions.
var (
	_ TrueTypeTable = IndexToLocationTable{}
	_ core.Writeable = IndexToLocationTable{}
)
