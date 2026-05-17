// Package truetypeparser provides types and helpers for TrueType font parsing.
package truetypeparser

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/cff"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	cmapsubtables "github.com/uglytoad/pdfpig/go/fonts/truetype/tables/cmap_sub_tables"
)

// TrueTypeFont represents a parsed TrueType font.
type TrueTypeFont struct {
	version            float32
	tableHeaders       map[string]truetype.TrueTypeHeaderTable
	tableRegister      *TableRegister
	cffFontCollection  *cff.CompactFontFormatFontCollection
	windowsUnicodeCMap cmapsubtables.ICMapSubTable
	macRomanCMap       cmapsubtables.ICMapSubTable
	windowsSymbolCMap  cmapsubtables.ICMapSubTable
	numberOfTables     int
}

// NewTrueTypeFont creates a new TrueTypeFont from parsed data.
// The tableHeaders and tableRegister must not be nil.
func NewTrueTypeFont(
	version float32,
	tableHeaders map[string]truetype.TrueTypeHeaderTable,
	tableRegister *TableRegister,
	cffFontCollection *cff.CompactFontFormatFontCollection,
) (*TrueTypeFont, error) {
	if tableHeaders == nil {
		return nil, errNilTableHeaders
	}

	if tableRegister == nil {
		return nil, errNilTableRegister
	}

	font := &TrueTypeFont{
		version:           version,
		tableHeaders:      tableHeaders,
		tableRegister:     tableRegister,
		cffFontCollection: cffFontCollection,
		numberOfTables:    len(tableHeaders),
	}

	font.extractCMapSubTables()

	return font, nil
}

// Version returns the font version number.
func (f *TrueTypeFont) Version() float32 {
	return f.version
}

// TableHeaders returns the table directory entries indicating the offset and length
// of the data for each table name.
func (f *TrueTypeFont) TableHeaders() map[string]truetype.TrueTypeHeaderTable {
	return f.tableHeaders
}

// TableRegister returns the actual parsed table data for this TrueType font.
func (f *TrueTypeFont) TableRegister() *TableRegister {
	return f.tableRegister
}

// Name returns the name of the font according to the font's name table,
// or an empty string if the name table is absent.
func (f *TrueTypeFont) Name() string {
	if f.tableRegister.NameTable == nil {
		return ""
	}
	return f.tableRegister.NameTable.FontName()
}

// WindowsUnicodeCMap returns the cmap subtable for Windows Unicode (platform 3, encoding 1).
// Returns nil if no such subtable is present.
func (f *TrueTypeFont) WindowsUnicodeCMap() cmapsubtables.ICMapSubTable {
	return f.windowsUnicodeCMap
}

// MacRomanCMap returns the cmap subtable for Mac Roman (platform 1, encoding 0).
// Returns nil if no such subtable is present.
func (f *TrueTypeFont) MacRomanCMap() cmapsubtables.ICMapSubTable {
	return f.macRomanCMap
}

// WindowsSymbolCMap returns the cmap subtable for Windows Symbol (platform 3, encoding 0).
// Returns nil if no such subtable is present.
func (f *TrueTypeFont) WindowsSymbolCMap() cmapsubtables.ICMapSubTable {
	return f.windowsSymbolCMap
}

// NumberOfTables returns the number of tables in this font.
func (f *TrueTypeFont) NumberOfTables() int {
	return f.numberOfTables
}

// extractCMapSubTables iterates over the CMap table's sub-tables and extracts
// the commonly used ones: Windows Unicode (3,1), Windows Symbol (3,0),
// and Mac Roman (1,0).
func (f *TrueTypeFont) extractCMapSubTables() {
	if f.tableRegister.CMapTable == nil {
		return
	}

	for _, subTable := range f.tableRegister.CMapTable.SubTables() {
		if f.windowsSymbolCMap == nil &&
			subTable.PlatformId() == cmapsubtables.Windows &&
			subTable.EncodingId() == 0 {
			f.windowsSymbolCMap = subTable
		} else if f.windowsUnicodeCMap == nil &&
			subTable.PlatformId() == cmapsubtables.Windows &&
			subTable.EncodingId() == 1 {
			f.windowsUnicodeCMap = subTable
		} else if f.macRomanCMap == nil &&
			subTable.PlatformId() == cmapsubtables.Macintosh &&
			subTable.EncodingId() == 0 {
			f.macRomanCMap = subTable
		}
	}
}

// TryGetBoundingBox attempts to get the bounding box for a glyph representing
// the specified character code. Uses default character-to-glyph mapping.
func (f *TrueTypeFont) TryGetBoundingBox(characterCode int) (core.PdfRectangle, bool) {
	return f.TryGetBoundingBoxWithMapping(characterCode, nil)
}

// TryGetBoundingBoxWithMapping attempts to get the bounding box for a glyph
// representing the specified character code using a custom character-to-glyph mapping.
func (f *TrueTypeFont) TryGetBoundingBoxWithMapping(
	characterCode int,
	characterCodeToGlyphId func(int) *int,
) (core.PdfRectangle, bool) {
	if f.tableRegister.GlyphTable.Tag() == "" {
		if f.cffFontCollection != nil {
			first := f.cffFontCollection.FirstFont()
			if first == nil {
				return core.PdfRectangle{}, false
			}
			name := first.GetCharacterName(characterCode, true)
			if name == "" {
				return core.PdfRectangle{}, false
			}

			bbox := first.GetCharacterBoundingBox(name)
			if bbox != nil {
				return *bbox, true
			}
		}

		return core.PdfRectangle{}, false
	}

	glyphIndex, ok := f.tryGetGlyphIndex(characterCode, characterCodeToGlyphId)
	if !ok {
		return core.PdfRectangle{}, false
	}

	boundingBox, ok := f.tableRegister.GlyphTable.TryGetGlyphBounds(glyphIndex)
	if !ok {
		return core.PdfRectangle{}, false
	}

	if boundingBox.Width == 0 {
		advanceWidth, ok := f.tryGetBoundingAdvancedWidthByIndex(glyphIndex)
		if ok {
			boundingBox = core.NewPdfRectangleFloat(0, 0, advanceWidth, 0)
		}
	}

	return boundingBox, true
}

// TryGetPath attempts to get the path for a glyph representing the specified
// character code. Uses default character-to-glyph mapping.
func (f *TrueTypeFont) TryGetPath(characterCode int) ([]*core.PdfSubpath, bool) {
	return f.TryGetPathWithMapping(characterCode, nil)
}

// TryGetPathWithMapping attempts to get the path for a glyph representing the
// specified character code using a custom character-to-glyph mapping.
func (f *TrueTypeFont) TryGetPathWithMapping(
	characterCode int,
	characterCodeToGlyphId func(int) *int,
) ([]*core.PdfSubpath, bool) {
	if f.tableRegister.GlyphTable.Tag() == "" {
		if f.cffFontCollection != nil {
			first := f.cffFontCollection.FirstFont()
			if first == nil {
				return nil, false
			}
			name := first.GetCharacterName(characterCode, true)
			if name == "" {
				return nil, false
			}

			path, ok := first.TryGetPath(name)
			return path, ok
		}

		return nil, false
	}

	glyphIndex, ok := f.tryGetGlyphIndex(characterCode, characterCodeToGlyphId)
	if !ok {
		return nil, false
	}

	return f.tableRegister.GlyphTable.TryGetGlyphPath(glyphIndex)
}

// TryGetAdvanceWidth attempts to get the advance width for a glyph representing
// the specified character code. Uses default character-to-glyph mapping.
func (f *TrueTypeFont) TryGetAdvanceWidth(characterCode int) (float64, bool) {
	return f.TryGetAdvanceWidthWithMapping(characterCode, nil)
}

// TryGetAdvanceWidthWithMapping attempts to get the advance width for a glyph
// representing the specified character code using a custom character-to-glyph mapping.
func (f *TrueTypeFont) TryGetAdvanceWidthWithMapping(
	characterCode int,
	characterCodeToGlyphId func(int) *int,
) (float64, bool) {
	glyphIndex, ok := f.tryGetGlyphIndex(characterCode, characterCodeToGlyphId)
	if !ok {
		return 0, false
	}

	width, ok := f.tryGetBoundingAdvancedWidthByIndex(glyphIndex)
	if !ok {
		return 0, false
	}

	return width, true
}

// GetUnitsPerEm returns the number of units per em for this font.
func (f *TrueTypeFont) GetUnitsPerEm() uint16 {
	return f.tableRegister.HeaderTable.UnitPerEm()
}

// tryGetBoundingAdvancedWidthByIndex gets the advance width from the horizontal
// metrics table for a glyph at the given index. Returns false if the horizontal
// metrics table is absent.
func (f *TrueTypeFont) tryGetBoundingAdvancedWidthByIndex(index int) (float64, bool) {
	if f.tableRegister.HorizontalMetricsTable.Tag() == "" {
		return 0, false
	}

	width, _ := f.tableRegister.HorizontalMetricsTable.GetAdvanceWidth(index)
	return float64(width), true
}

// tryGetGlyphIndex resolves a character code to a glyph index. If an external
// mapping function is provided and returns a value, it is used; otherwise the
// CMap table is consulted. Returns false if no mapping can be found.
func (f *TrueTypeFont) tryGetGlyphIndex(
	characterIdentifier int,
	characterCodeToGlyphId func(int) *int,
) (int, bool) {
	if characterCodeToGlyphId != nil {
		externalGlyphId := characterCodeToGlyphId(characterIdentifier)
		if externalGlyphId != nil {
			return *externalGlyphId, true
		}
	}

	if f.tableRegister.CMapTable == nil {
		return 0, false
	}

	return f.tableRegister.CMapTable.TryGetGlyphIndex(characterIdentifier)
}

var (
	errNilTableHeaders = err("table headers cannot be nil")
	errNilTableRegister = err("table register cannot be nil")
)

type err string

func (e err) Error() string {
	return string(e)
}
