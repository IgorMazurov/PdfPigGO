// Package truetypeparser provides CMap sub-table loading functions.
package truetypeparser

import (
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
	cmapsubtables "github.com/uglytoad/pdfpig/go/fonts/truetype/tables/cmap_sub_tables"
)

/* ------------------------------------------------------------------ */
/*  Compile-time assertions for TrueTypeTableParser implementations   */
/* ------------------------------------------------------------------ */

// Compile-time assertion: CMapTableParser implements TrueTypeTableParser[*CMapTable].
var _ TrueTypeTableParser[*CMapTable] = (*CMapTableParser)(nil)

// Compile-time assertion: HorizontalHeaderTableParser implements TrueTypeTableParser[tables.HorizontalHeaderTable].
var _ TrueTypeTableParser[tables.HorizontalHeaderTable] = (*HorizontalHeaderTableParser)(nil)

// Compile-time assertion: HorizontalMetricsTableParser implements TrueTypeTableParser[tables.HorizontalMetricsTable].
var _ TrueTypeTableParser[tables.HorizontalMetricsTable] = (*HorizontalMetricsTableParser)(nil)

// Compile-time assertion: NameTableParser implements TrueTypeTableParser[*tables.NameTable].
var _ TrueTypeTableParser[*tables.NameTable] = (*NameTableParser)(nil)

/* ------------------------------------------------------------------ */
/*  CMap Sub-Table Load functions                                     */
/* ------------------------------------------------------------------ */

// ByteEncodingCMapTableLoad reads a format 0 CMap sub-table from the data stream.
func ByteEncodingCMapTableLoad(data *TrueTypeDataBytes, platformID cmapsubtables.TrueTypeCMapPlatform, encodingID uint16) (cmapsubtables.ICMapSubTable, error) {
	length := data.ReadUnsignedShort()
	language := data.ReadUnsignedShort()

	if length == 0 {
		return cmapsubtables.NewByteEncodingCMapTable(platformID, encodingID, language, nil), nil
	}

	glyphMapping := data.ReadByteArray(int(length) - 2*3)

	return cmapsubtables.NewByteEncodingCMapTable(platformID, encodingID, language, glyphMapping), nil
}

// HighByteMappingCMapTableLoad reads a format 2 CMap sub-table from the data stream.
func HighByteMappingCMapTableLoad(data *TrueTypeDataBytes, numberOfGlyphs int, platformID cmapsubtables.TrueTypeCMapPlatform, encodingID uint16) (cmapsubtables.ICMapSubTable, error) {
	data.ReadUnsignedShort() // length
	data.ReadUnsignedShort() // version

	subHeaderKeys := make([]int, 256)
	maximumSubHeaderIndex := 0

	for i := 0; i < 256; i++ {
		value := int(data.ReadUnsignedShort())
		if value/8 > maximumSubHeaderIndex {
			maximumSubHeaderIndex = value / 8
		}
		subHeaderKeys[i] = value
	}

	subHeaderCount := maximumSubHeaderIndex + 1

	subHeaders := make([]cmapsubtables.SubHeader, subHeaderCount)

	for i := 0; i < subHeaderCount; i++ {
		firstCode := int(data.ReadUnsignedShort())
		entryCount := int(data.ReadUnsignedShort())
		idDelta := int16(data.ReadSignedShort())
		idRangeOffset := int(data.ReadUnsignedShort()) - (subHeaderCount-i-1)*8 - 2

		subHeaders[i] = cmapsubtables.SubHeader{
			FirstCode:     firstCode,
			EntryCount:    entryCount,
			IdDelta:       idDelta,
			IdRangeOffset: idRangeOffset,
		}
	}

	glyphIndexArrayOffset := data.Position()

	characterCodeToGlyphId := make(map[int]int)

	for i := 0; i < subHeaderCount; i++ {
		subHeader := subHeaders[i]

		if _, err := data.Seek(glyphIndexArrayOffset+int64(subHeader.IdRangeOffset), 0); err != nil {
			return nil, err
		}

		for j := 0; j < subHeader.EntryCount; j++ {
			characterCode := (i << 8) + (subHeader.FirstCode + j)

			p := int(data.ReadUnsignedShort())

			if p > 0 {
				p = (p + int(subHeader.IdDelta)) % 65536
			}

			if p >= numberOfGlyphs {
				continue
			}

			characterCodeToGlyphId[characterCode] = p
		}
	}

	return cmapsubtables.NewHighByteMappingCMapTable(platformID, encodingID, characterCodeToGlyphId), nil
}

// Format4CMapTableLoad reads a format 4 CMap sub-table from the data stream.
func Format4CMapTableLoad(data *TrueTypeDataBytes, platformID cmapsubtables.TrueTypeCMapPlatform, encodingID uint16) (cmapsubtables.ICMapSubTable, error) {
	length := data.ReadUnsignedShort()
	language := data.ReadUnsignedShort()

	doubleSegmentCount := data.ReadUnsignedShort()
	segmentCount := int(doubleSegmentCount / 2)

	data.ReadUnsignedShort() // searchRange
	data.ReadUnsignedShort() // entrySelector
	data.ReadUnsignedShort() // rangeShift

	endCounts := data.ReadUnsignedShortArray(segmentCount)
	data.ReadUnsignedShort() // reservedPad
	startCounts := data.ReadUnsignedShortArray(segmentCount)
	idDeltas := data.ReadSignedShortArray(segmentCount)
	idRangeOffsets := data.ReadUnsignedShortArray(segmentCount)

	remainingBytes := int(length) - (8*2 + 3*segmentCount*2)
	remainingInts := remainingBytes / 2

	glyphIndices := data.ReadUnsignedShortArray(remainingInts)

	segments := make([]cmapsubtables.Segment, segmentCount)
	for i := 0; i < segmentCount; i++ {
		segments[i] = cmapsubtables.Segment{
			StartCode:     int(startCounts[i]),
			EndCode:       int(endCounts[i]),
			IdDelta:       int(idDeltas[i]),
			IdRangeOffset: int(idRangeOffsets[i]),
		}
	}

	return cmapsubtables.NewFormat4CMapTable(platformID, encodingID, language, segments, glyphIndices), nil
}

// TrimmedTableMappingCMapTableLoad reads a format 6 CMap sub-table from the data stream.
func TrimmedTableMappingCMapTableLoad(data *TrueTypeDataBytes, platformID cmapsubtables.TrueTypeCMapPlatform, encodingID uint16) (cmapsubtables.ICMapSubTable, error) {
	data.ReadUnsignedShort() // length
	data.ReadUnsignedShort() // language

	firstCode := int(data.ReadUnsignedShort())
	entryCount := int(data.ReadUnsignedShort())

	glyphIndices := data.ReadUnsignedShortArray(entryCount)

	return cmapsubtables.NewTrimmedTableMappingCMapTable(platformID, encodingID, firstCode, entryCount, glyphIndices), nil
}
