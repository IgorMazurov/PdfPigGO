package truetypeparser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
)

var cmapParser = &CMapTableParser{}

var horizontalMetricsParser = &HorizontalMetricsTableParser{}

var nameParser = &NameTableParser{}

var os2Parser = &Os2TableParser{}

var horizontalHeaderParser = &HorizontalHeaderTableParser{}

// ParseTable seeks to the table offset in the raw data and delegates parsing to
// the provided parser. Callers pass a pointer to an already-constructed parser
// instance, e.g.:
//
//	cmap, err := ParseTable[truetypeparser.CMapTable](cmapParser, header, data, register)
func ParseTable[T any](
	parser TrueTypeTableParser[T],
	header truetype.TrueTypeHeaderTable,
	data *TrueTypeDataBytes,
	register *TableRegisterBuilder,
) (T, error) {
	if _, err := data.Seek(int64(header.Offset), 0); err != nil {
		var zero T
		return zero, fmt.Errorf("failed to seek to table offset %d: %w", header.Offset, err)
	}

	result, err := parser.Parse(header, data, register)
	if err != nil {
		var zero T
		return zero, err
	}

	return result, nil
}

// ParseCMap parses the cmap (character-to-glyph mapping) table from raw data.
func ParseCMap(
	header truetype.TrueTypeHeaderTable,
	data *TrueTypeDataBytes,
	register *TableRegisterBuilder,
) (*CMapTable, error) {
	return ParseTable[*CMapTable](cmapParser, header, data, register)
}

// ParseHorizontalMetrics parses the horizontal metrics (hmtx) table from raw data.
func ParseHorizontalMetrics(
	header truetype.TrueTypeHeaderTable,
	data *TrueTypeDataBytes,
	register *TableRegisterBuilder,
) (tables.HorizontalMetricsTable, error) {
	return ParseTable[tables.HorizontalMetricsTable](horizontalMetricsParser, header, data, register)
}

// ParseName parses the name (name) table from raw data.
func ParseName(
	header truetype.TrueTypeHeaderTable,
	data *TrueTypeDataBytes,
	register *TableRegisterBuilder,
) (*tables.NameTable, error) {
	return ParseTable[*tables.NameTable](nameParser, header, data, register)
}

// ParseOs2 parses the OS/2 table from raw data. Returns a versioned Os2Table
// struct depending on the font's OS/2 table version field.
func ParseOs2(
	header truetype.TrueTypeHeaderTable,
	data *TrueTypeDataBytes,
	register *TableRegisterBuilder,
) (any, error) {
	if _, err := data.Seek(int64(header.Offset), 0); err != nil {
		return nil, fmt.Errorf("failed to seek to OS/2 table offset %d: %w", header.Offset, err)
	}

	return os2Parser.Parse(header, data, register)
}

// ParseHorizontalHeader parses the horizontal header (hhea) table from raw data.
func ParseHorizontalHeader(
	header truetype.TrueTypeHeaderTable,
	data *TrueTypeDataBytes,
	register *TableRegisterBuilder,
) (tables.HorizontalHeaderTable, error) {
	return ParseTable[tables.HorizontalHeaderTable](horizontalHeaderParser, header, data, register)
}
