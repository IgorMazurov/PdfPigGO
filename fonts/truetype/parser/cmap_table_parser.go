package truetypeparser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	cmapsubtables "github.com/uglytoad/pdfpig/go/fonts/truetype/tables/cmap_sub_tables"
)

// CMapTableParser parses the cmap (character-to-glyph mapping) table of a TrueType font.
type CMapTableParser struct{}

// Parse reads and interprets the CMap table from the raw TrueType data.
func (p *CMapTableParser) Parse(header truetype.TrueTypeHeaderTable, data *TrueTypeDataBytes, register *TableRegisterBuilder) (*CMapTable, error) {
	if _, err := data.Seek(int64(header.Offset), 0); err != nil {
		return nil, fmt.Errorf("failed to seek to cmap table offset: %w", err)
	}

	tableVersionNumber := data.ReadUnsignedShort()

	numberOfEncodingTables := data.ReadUnsignedShort()

	subTableHeaders := make([]subTableHeaderEntry, numberOfEncodingTables)

	for i := uint16(0); i < numberOfEncodingTables; i++ {
		platformId := cmapsubtables.TrueTypeCMapPlatform(data.ReadUnsignedShort())
		encodingId := data.ReadUnsignedShort()
		offset := data.ReadUnsignedInt()

		subTableHeaders[i] = subTableHeaderEntry{
			PlatformId: platformId,
			EncodingId: encodingId,
			Offset:     offset,
		}
	}

	tables := make([]cmapsubtables.ICMapSubTable, 0, numberOfEncodingTables)

	numberofGlyphs := register.MaximumProfileTable.NumberOfGlyphs()

	for i := range subTableHeaders {
		subTableHeader := subTableHeaders[i]

		if _, err := data.Seek(int64(header.Offset)+int64(subTableHeader.Offset), 0); err != nil {
			return nil, fmt.Errorf("failed to seek to cmap sub-table offset: %w", err)
		}

		format := data.ReadUnsignedShort()

		switch format {
		case 0:
			item, err := ByteEncodingCMapTableLoad(data, subTableHeader.PlatformId, subTableHeader.EncodingId)
			if err != nil {
				return nil, fmt.Errorf("failed to load byte encoding cmap table: %w", err)
			}
			tables = append(tables, item)

		case 2:
			item, err := HighByteMappingCMapTableLoad(data, numberofGlyphs, subTableHeader.PlatformId, subTableHeader.EncodingId)
			if err != nil {
				return nil, fmt.Errorf("failed to load high byte mapping cmap table: %w", err)
			}
			tables = append(tables, item)

		case 4:
			item, err := Format4CMapTableLoad(data, subTableHeader.PlatformId, subTableHeader.EncodingId)
			if err != nil {
				return nil, fmt.Errorf("failed to load format 4 cmap table: %w", err)
			}
			tables = append(tables, item)

		case 6:
			item, err := TrimmedTableMappingCMapTableLoad(data, subTableHeader.PlatformId, subTableHeader.EncodingId)
			if err != nil {
				return nil, fmt.Errorf("failed to load trimmed table mapping cmap table: %w", err)
			}
			tables = append(tables, item)
		}
	}

	return NewCMapTable(tableVersionNumber, header, tables), nil
}

// subTableHeaderEntry holds the header information for a single CMap sub-table.
type subTableHeaderEntry struct {
	PlatformId cmapsubtables.TrueTypeCMapPlatform
	EncodingId uint16
	Offset     uint32
}
