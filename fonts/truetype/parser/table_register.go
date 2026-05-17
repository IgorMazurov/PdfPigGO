package truetypeparser

import (
	"errors"

	cmapsubtables "github.com/uglytoad/pdfpig/go/fonts/truetype/tables/cmap_sub_tables"

	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
)

// TableRegister holds the set of tables in a TrueType font interpreted by the library.
type TableRegister struct {
	// HeaderTable contains global information about the font.
	HeaderTable tables.HeaderTable

	// GlyphTable contains the data that defines the appearance of the glyphs in the font.
	GlyphTable *tables.GlyphDataTable

	// HorizontalHeaderTable contains information needed to layout fonts whose characters
	// are written horizontally.
	HorizontalHeaderTable tables.HorizontalHeaderTable

	// HorizontalMetricsTable contains metric information for the horizontal layout of each
	// glyph in the font.
	HorizontalMetricsTable tables.HorizontalMetricsTable

	// IndexToLocationTable stores the offsets to the locations of the glyphs relative to
	// the glyph table.
	IndexToLocationTable tables.IndexToLocationTable

	// MaximumProfileTable establishes the memory requirements for the font.
	MaximumProfileTable tables.BasicMaximumProfileTable

	// NameTable defines strings used by the font.
	NameTable *tables.NameTable

	// PostScriptTable contains information needed to use a TrueType font on a PostScript
	// printer, including PostScript names for all glyphs.
	PostScriptTable tables.PostScriptTable

	// CMapTable defines mapping of character codes to glyph index values in the font.
	CMapTable *CMapTable

	// KerningTable contains kerning pair data.
	KerningTable *tables.KerningTable

	// Os2Table consists of metrics required by Windows.
	Os2Table tables.IOs2Table
}

// NewTableRegister creates a TableRegister from the builder, validating that all
// required tables are present.
func NewTableRegister(builder *TableRegisterBuilder) (*TableRegister, error) {
	if builder == nil {
		return nil, errors.New("builder must not be nil")
	}

	if builder.HeaderTable.Tag() == "" {
		return nil, errors.New("the builder did not contain the header table")
	}

	if builder.HorizontalHeaderTable.Tag() == "" {
		return nil, errors.New("the builder did not contain the horizontal header table")
	}

	if builder.MaximumProfileTable.Tag() == "" {
		return nil, errors.New("the builder did not contain the maximum profile table")
	}

	return &TableRegister{
		HeaderTable:            builder.HeaderTable,
		GlyphTable:             builder.GlyphDataTable,
		HorizontalHeaderTable:  builder.HorizontalHeaderTable,
		HorizontalMetricsTable: builder.HorizontalMetricsTable,
		IndexToLocationTable:   builder.IndexToLocationTable,
		MaximumProfileTable:    builder.MaximumProfileTable,
		NameTable:              builder.NameTable,
		PostScriptTable:        builder.PostScriptTable,
		CMapTable:              builder.CMapTable,
		KerningTable:           builder.KerningTable,
		Os2Table:               builder.Os2Table,
	}, nil
}

// TableRegisterBuilder collects parsed TrueType tables during font loading.
type TableRegisterBuilder struct {
	HeaderTable            tables.HeaderTable
	GlyphDataTable         *tables.GlyphDataTable
	HorizontalHeaderTable  tables.HorizontalHeaderTable
	HorizontalMetricsTable tables.HorizontalMetricsTable
	IndexToLocationTable   tables.IndexToLocationTable
	MaximumProfileTable    tables.BasicMaximumProfileTable
	NameTable              *tables.NameTable
	PostScriptTable        tables.PostScriptTable
	CMapTable              *CMapTable
	KerningTable           *tables.KerningTable
	Os2Table               tables.IOs2Table
}

// CMapTable is an alias for the real type in the tables package.
type CMapTable = tables.CMapTable

// NewCMapTable creates a new CMapTable instance, delegating to the tables package.
func NewCMapTable(version uint16, directoryTable truetype.TrueTypeHeaderTable, subTables []cmapsubtables.ICMapSubTable) *CMapTable {
	return tables.NewCMapTable(version, directoryTable, subTables)
}
