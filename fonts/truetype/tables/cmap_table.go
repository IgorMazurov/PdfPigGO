package tables

import (
	"errors"
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	cmapsubtables "github.com/uglytoad/pdfpig/go/fonts/truetype/tables/cmap_sub_tables"
)

// CMapTable maps character codes to glyph indices in a TrueType font.
// The choice of encoding for a particular font is dependent on the conventions
// used by the intended platform. The cmap table can contain multiple encoding
// tables for use on different platforms, one for each supported encoding scheme.
type CMapTable struct {
	version        uint16
	directoryTable truetype.TrueTypeHeaderTable
	subTables      []cmapsubtables.ICMapSubTable
}

// NewCMapTable creates a new CMapTable instance.
func NewCMapTable(version uint16, directoryTable truetype.TrueTypeHeaderTable, subTables []cmapsubtables.ICMapSubTable) *CMapTable {
	return &CMapTable{
		version:        version,
		directoryTable: directoryTable,
		subTables:      subTables,
	}
}

// Tag returns the 4-letter identifier for this table.
func (c *CMapTable) Tag() string {
	return truetype.Cmap
}

// DirectoryTable returns the directory entry from the font's offset subtable.
func (c *CMapTable) DirectoryTable() truetype.TrueTypeHeaderTable {
	return c.directoryTable
}

// Version returns the version number (always 0).
func (c *CMapTable) Version() uint16 {
	return c.version
}

// SubTables returns the sub-tables, one for each supported encoding scheme and platform.
func (c *CMapTable) SubTables() []cmapsubtables.ICMapSubTable {
	return c.subTables
}

// TryGetGlyphIndex gets the glyph index for the corresponding character code.
// Returns the glyph index and true if found, or 0 and false otherwise.
func (c *CMapTable) TryGetGlyphIndex(characterCode int) (int, bool) {
	if len(c.subTables) == 0 {
		return 0, false
	}

	var windowsMapping cmapsubtables.ICMapSubTable

	for _, subTable := range c.subTables {
		glyphIndex := subTable.CharacterCodeToGlyphIndex(characterCode)

		if glyphIndex != 0 {
			return glyphIndex, true
		}

		if subTable.EncodingId() == 0 && subTable.PlatformId() == cmapsubtables.Windows {
			windowsMapping = subTable
		}
	}

	if windowsMapping != nil && characterCode >= 0 && characterCode <= 255 {
		glyphIndex := windowsMapping.CharacterCodeToGlyphIndex(characterCode + 0xF000)

		if glyphIndex != 0 {
			return glyphIndex, true
		}

		glyphIndex = windowsMapping.CharacterCodeToGlyphIndex(characterCode + 0xF100)

		if glyphIndex != 0 {
			return glyphIndex, true
		}

		glyphIndex = windowsMapping.CharacterCodeToGlyphIndex(characterCode + 0xF200)

		if glyphIndex != 0 {
			return glyphIndex, true
		}
	}

	return 0, false
}

// Write writes the CMap table data to the writer.
// The writer must also implement io.Seeker for offset patching.
func (c *CMapTable) Write(w io.Writer) error {
	seeker, ok := w.(io.Seeker)
	if !ok {
		return errors.New("writer does not support seeking")
	}

	startPosition, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("failed to get stream position: %w", err)
	}

	if _, err := core.WriteUShort(w, c.version); err != nil {
		return err
	}

	if _, err := core.WriteUShortInt32(w, int32(len(c.subTables))); err != nil {
		return err
	}

	subTableIndexOffsetPositions := make([]int64, len(c.subTables))
	for i, subTable := range c.subTables {
		if _, err := core.WriteUShort(w, uint16(subTable.PlatformId())); err != nil {
			return err
		}

		if _, err := core.WriteUShort(w, subTable.EncodingId()); err != nil {
			return err
		}

		subTableIndexOffsetPositions[i], err = seeker.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("failed to get offset position: %w", err)
		}

		if _, err := core.WriteUInt(w, 0); err != nil {
			return err
		}
	}

	subTableActualPositions := make([]uint32, len(c.subTables))
	for i, subTable := range c.subTables {
		writeableSubTable, ok := subTable.(core.Writeable)
		if !ok {
			return fmt.Errorf("cannot write subtable of type: %T", subTable)
		}

		currentPos, err := seeker.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("failed to get current position: %w", err)
		}

		subTableActualPositions[i] = uint32(currentPos - startPosition)

		if err := writeableSubTable.Write(w); err != nil {
			return err
		}
	}

	endAt, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("failed to get end position: %w", err)
	}

	for i := range subTableIndexOffsetPositions {
		actual := subTableActualPositions[i]

		if _, err := seeker.Seek(subTableIndexOffsetPositions[i], io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek to offset position: %w", err)
		}

		if _, err := core.WriteUInt(w, actual); err != nil {
			return err
		}
	}

	if _, err := seeker.Seek(endAt, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek to end position: %w", err)
	}

	return nil
}
