package truetypeparser

import (
	"fmt"
	"log"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/cff"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables/kerning"
	"golang.org/x/text/encoding"
)

// ParseFull reads TrueType font data and returns a fully constructed TrueTypeFont.
func ParseFull(data *TrueTypeDataBytes) (*TrueTypeFont, error) {
	if data == nil || len(data.data) == 0 {
		return nil, fmt.Errorf("font data must not be empty")
	}

	version := float32(data.Read32Fixed())
	numTables := int(data.ReadUnsignedShort())

	_ = data.ReadUnsignedShort() // searchRange
	_ = data.ReadUnsignedShort() // entrySelector
	_ = data.ReadUnsignedShort() // rangeShift

	tablesMap := make(map[string]truetype.TrueTypeHeaderTable)

	for i := 0; i < numTables; i++ {
		if table := readTable(data); table.Tag != "" {
			tablesMap[table.Tag] = table
		}
	}

	return parseTables(version, tablesMap, data)
}

func readTable(data *TrueTypeDataBytes) truetype.TrueTypeHeaderTable {
	tag := data.ReadTag()
	checksum := data.ReadUnsignedInt()
	offset := data.ReadUnsignedInt()
	length := data.ReadUnsignedInt()

	if length == 0 && !strings.EqualFold(tag, truetype.Glyf) {
		return truetype.TrueTypeHeaderTable{}
	}

	return truetype.TrueTypeHeaderTable{
		Tag:      tag,
		CheckSum: checksum,
		Offset:   offset,
		Length:   length,
	}
}

func parseTables(version float32, tablesMap map[string]truetype.TrueTypeHeaderTable, data *TrueTypeDataBytes) (*TrueTypeFont, error) {
	isPostScript := false
	var cffFontCollection *cff.CompactFontFormatFontCollection

	if cffTable, ok := tablesMap[truetype.Cff]; ok {
		isPostScript = true
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("failed to parse CFF table: %v", r)
				}
			}()

			if _, err := data.Seek(int64(cffTable.Offset), 0); err != nil {
				return
			}
			buffer := data.ReadByteArray(int(cffTable.Length))
			cffData := cff.NewCompactFontFormatData(buffer)
			fontCollection, parseErr := cff.Parse(cffData, nil)
			if parseErr != nil {
				log.Printf("failed to parse CFF font collection: %v", parseErr)
				return
			}
			cffFontCollection = fontCollection
		}()
	}

	builder := &TableRegisterBuilder{}

	headTable, ok := tablesMap[truetype.Head]
	if !ok {
		return nil, fonts.NewInvalidFontFormatException(fmt.Sprintf("the %s table is required", truetype.Head))
	}
	builder.HeaderTable = loadHeaderTable(data, headTable)

	hheaTable, ok := tablesMap[truetype.Hhea]
	if !ok {
		return nil, fonts.NewInvalidFontFormatException("the horizontal header table is required")
	}
	hht, err := ParseHorizontalHeader(hheaTable, data, builder)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hhea table: %w", err)
	}
	builder.HorizontalHeaderTable = hht

	maxpTable, ok := tablesMap[truetype.Maxp]
	if !ok {
		return nil, fonts.NewInvalidFontFormatException("the maximum profile table is required")
	}
	builder.MaximumProfileTable = loadMaxpTable(data, maxpTable)

	if postTable, ok := tablesMap[truetype.Post]; ok {
		builder.PostScriptTable = loadPostScriptTable(data, postTable, builder.MaximumProfileTable)
	}

	if nameTableHeader, ok := tablesMap[truetype.Name]; ok {
		nt, err := ParseName(nameTableHeader, data, builder)
		if err != nil {
			return nil, fmt.Errorf("failed to parse name table: %w", err)
		}
		builder.NameTable = nt
	}

	if os2TableHeader, ok := tablesMap[truetype.Os2]; ok {
		os2Result, err := ParseOs2(os2TableHeader, data, builder)
		if err == nil && os2Result != nil {
			switch t := os2Result.(type) {
			case tables.Os2Table:
				builder.Os2Table = t
			case tables.Os2RevisedVersion0Table:
				builder.Os2Table = t
			case tables.Os2Version1Table:
				builder.Os2Table = t
			case tables.Os2Version2To4OpenTypeTable:
				builder.Os2Table = t
			case tables.Os2Version5OpenTypeTable:
				builder.Os2Table = t
			}
		}
	}

	if !isPostScript {
		locaTable, ok := tablesMap[truetype.Loca]
		if !ok {
			return nil, fonts.NewInvalidFontFormatException("the location to index table is required for non-PostScript fonts")
		}
		builder.IndexToLocationTable = loadLocaTable(data, locaTable, builder)

		glyfTable, ok := tablesMap[truetype.Glyf]
		if !ok {
			return nil, fonts.NewInvalidFontFormatException("the glyph table is required for non-PostScript fonts")
		}
		builder.GlyphDataTable = loadGlyphTable(data, glyfTable, builder)
	}

	optionallyParseTables(tablesMap, data, builder)

	register, err := NewTableRegister(builder)
	if err != nil {
		return nil, fmt.Errorf("failed to build table register: %w", err)
	}

	font, err := NewTrueTypeFont(version, tablesMap, register, cffFontCollection)
	if err != nil {
		return nil, fmt.Errorf("failed to create TrueType font: %w", err)
	}

	return font, nil
}

func optionallyParseTables(tablesMap map[string]truetype.TrueTypeHeaderTable, data *TrueTypeDataBytes, builder *TableRegisterBuilder) {
	if cmapHeader, ok := tablesMap[truetype.Cmap]; ok {
		cmap, err := ParseCMap(cmapHeader, data, builder)
		if err == nil && cmap != nil {
			builder.CMapTable = cmap
		}
	}

	if hmtxHeader, ok := tablesMap[truetype.Hmtx]; ok {
		hmt, err := ParseHorizontalMetrics(hmtxHeader, data, builder)
		if err == nil {
			builder.HorizontalMetricsTable = hmt
		}
	}

	if kernHeader, ok := tablesMap[truetype.Kern]; ok {
		kerningTable := loadKerningTable(data, kernHeader)
		if kerningTable != nil {
			builder.KerningTable = kerningTable
		}
	}
}

// --- Load functions for tables that don't use the generic parser pattern ---

func loadHeaderTable(data *TrueTypeDataBytes, directoryTable truetype.TrueTypeHeaderTable) tables.HeaderTable {
	if _, err := data.Seek(int64(directoryTable.Offset), 0); err != nil {
		return tables.NewHeaderTable(directoryTable, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, tables.HeaderMacStyleNone, 0, tables.FullyMixedDirectional, tables.IndexToLocationTableShort, 0)
	}

	version := float32(data.Read32Fixed())
	fontRevision := float32(data.Read32Fixed())
	checkSumAdjustment := data.ReadUnsignedInt()
	magicNumber := data.ReadUnsignedInt()
	flags := data.ReadUnsignedShort()
	unitPerEm := data.ReadUnsignedShort()

	createdHi := data.ReadUnsignedInt()
	createdLo := data.ReadUnsignedInt()
	created := int64(createdHi)<<32 | int64(createdLo)

	modifiedHi := data.ReadUnsignedInt()
	modifiedLo := data.ReadUnsignedInt()
	modified := int64(modifiedHi)<<32 | int64(modifiedLo)

	xMin := data.ReadSignedShort()
	yMin := data.ReadSignedShort()
	xMax := data.ReadSignedShort()
	yMax := data.ReadSignedShort()
	macStyle := tables.HeaderMacStyle(data.ReadUnsignedShort())
	lowestRecommendedPpem := data.ReadUnsignedShort()
	fontDirectionHint := tables.FontDirection(data.ReadSignedShort())
	indexToLocFormat := tables.IndexToLocationTableEntryFormat(data.ReadSignedShort())
	glyphDataFormat := data.ReadSignedShort()

	return tables.NewHeaderTable(
		directoryTable, version, fontRevision, checkSumAdjustment, magicNumber,
		flags, unitPerEm, created, modified, xMin, yMin, xMax, yMax,
		macStyle, lowestRecommendedPpem, fontDirectionHint, indexToLocFormat, glyphDataFormat,
	)
}

func loadMaxpTable(data *TrueTypeDataBytes, directoryTable truetype.TrueTypeHeaderTable) tables.BasicMaximumProfileTable {
	if _, err := data.Seek(int64(directoryTable.Offset), 0); err != nil {
		return tables.NewBasicMaximumProfileTable(directoryTable, 0, 0)
	}

	version := float32(data.Read32Fixed())
	numberOfGlyphs := int(data.ReadUnsignedShort())
	_ = version

	return tables.NewBasicMaximumProfileTable(directoryTable, version, numberOfGlyphs)
}

func loadPostScriptTable(
	data *TrueTypeDataBytes,
	directoryTable truetype.TrueTypeHeaderTable,
	maxProfile tables.BasicMaximumProfileTable,
) tables.PostScriptTable {
	if _, err := data.Seek(int64(directoryTable.Offset), 0); err != nil {
		return tables.NewPostScriptTable(directoryTable, 0, 0, 0, 0, 0, 0, 0, 0, 0, nil)
	}

	format := float32(data.Read32Fixed())
	italicAngle := float32(data.Read32Fixed())
	underlinePosition := data.ReadSignedShort()
	underlineThickness := data.ReadSignedShort()
	isFixedPitch := data.ReadUnsignedInt()
	minMemoryType42 := data.ReadUnsignedInt()
	maxMemoryType42 := data.ReadUnsignedInt()
	minMemoryType1 := data.ReadUnsignedInt()
	maxMemoryType1 := data.ReadUnsignedInt()

	glyphNames := loadPostScriptGlyphNames(data, maxProfile.NumberOfGlyphs(), format)

	return tables.NewPostScriptTable(
		directoryTable, format, italicAngle, underlinePosition, underlineThickness,
		isFixedPitch, minMemoryType42, maxMemoryType42, minMemoryType1, maxMemoryType1,
		glyphNames,
	)
}

func loadPostScriptGlyphNames(data *TrueTypeDataBytes, numberOfGlyphs int, format float32) []string {
	if floatWithinEpsilon(format, 1.0) {
		names := make([]string, truetype.NumberOfMacGlyphs)
		copy(names, truetype.MacGlyphNames)
		return names
	}

	if floatWithinEpsilon(format, 2.0) {
		return loadFormat2GlyphNames(data)
	}

	if floatWithinEpsilon(format, 2.5) {
		nameIndices := make([]int, numberOfGlyphs)
		for i := 0; i < numberOfGlyphs; i++ {
			offset := data.ReadSignedByte()
			nameIndices[i] = i + 1 + int(offset)
		}

		names := make([]string, numberOfGlyphs)
		for i, idx := range nameIndices {
			if idx >= 0 && idx < truetype.NumberOfMacGlyphs {
				names[i] = truetype.MacGlyphNames[idx]
			}
		}
		return names
	}

	return []string{}
}

func loadFormat2GlyphNames(data *TrueTypeDataBytes) []string {
	const reservedIndexStart = 32768

	numGlyphs := int(data.ReadUnsignedShort())
	nameIndices := make([]uint16, numGlyphs)
	names := make([]string, numGlyphs)

	maxIndex := -1
	for i := 0; i < numGlyphs; i++ {
		idx := data.ReadUnsignedShort()
		nameIndices[i] = idx
		if idx < reservedIndexStart && int(idx) > maxIndex {
			maxIndex = int(idx)
		}
	}

	var customNameArray []string
	if maxIndex >= truetype.NumberOfMacGlyphs {
		customCount := maxIndex - truetype.NumberOfMacGlyphs + 1
		customNameArray = make([]string, customCount)
		for i := 0; i < customCount; i++ {
			nc, _ := data.ReadByte()
			numChars := int(nc)
			if s, ok := data.TryReadString(numChars, encoding.Nop); ok {
				customNameArray[i] = s
			}
		}
	}

	for i := 0; i < numGlyphs; i++ {
		idx := nameIndices[i]
		if idx < truetype.NumberOfMacGlyphs {
			names[i] = truetype.MacGlyphNames[idx]
		} else if idx >= truetype.NumberOfMacGlyphs && idx < reservedIndexStart {
			names[i] = customNameArray[idx-truetype.NumberOfMacGlyphs]
		} else {
			names[i] = ".undefined"
		}
	}

	return names
}

func floatWithinEpsilon(a, b float32) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < 1e-6
}

func loadLocaTable(
	data *TrueTypeDataBytes,
	directoryTable truetype.TrueTypeHeaderTable,
	builder *TableRegisterBuilder,
) tables.IndexToLocationTable {
	if _, err := data.Seek(int64(directoryTable.Offset), 0); err != nil {
		return tables.NewIndexToLocationTable(directoryTable, tables.IndexToLocationTableShort, nil)
	}

	numGlyphs := builder.MaximumProfileTable.NumberOfGlyphs()
	format := builder.HeaderTable.IndexToLocFormat()

	glyphOffsets := make([]uint32, numGlyphs+1)
	lastOffset := uint32(0)

	for i := 0; i <= numGlyphs; i++ {
		var offset uint32
		if format == tables.IndexToLocationTableShort {
			offset = uint32(data.ReadUnsignedShort()) * 2
		} else {
			offset = data.ReadUnsignedInt()
		}
		if offset > lastOffset {
			lastOffset = offset
		}
		glyphOffsets[i] = offset
	}

	return tables.NewIndexToLocationTable(directoryTable, format, glyphOffsets)
}

func loadGlyphTable(
	data *TrueTypeDataBytes,
	directoryTable truetype.TrueTypeHeaderTable,
	builder *TableRegisterBuilder,
) *tables.GlyphDataTable {
	if _, err := data.Seek(int64(directoryTable.Offset), 0); err != nil {
		gt := tables.NewGlyphDataTable(directoryTable, nil, core.PdfRectangle{}, nil)
		return &gt
	}

	bytes := data.ReadByteArray(int(directoryTable.Length))

	maxBounds := core.NewPdfRectangleFromInt(
		int(builder.HeaderTable.XMin()),
		int(builder.HeaderTable.YMin()),
		int(builder.HeaderTable.XMax()),
		int(builder.HeaderTable.YMax()),
	)

	glyphOffsets := builder.IndexToLocationTable.GlyphOffsets()

	gt := tables.NewGlyphDataTable(directoryTable, glyphOffsets, maxBounds, bytes)
	return &gt
}

func loadKerningTable(data *TrueTypeDataBytes, directoryTable truetype.TrueTypeHeaderTable) *tables.KerningTable {
	if _, err := data.Seek(int64(directoryTable.Offset), 0); err != nil {
		return tables.NewKerningTable(nil)
	}

	version := data.ReadUnsignedShort()
	numTablesVal := int(data.ReadUnsignedShort())
	subTables := make([]*kerning.KerningSubTable, 0, numTablesVal)

	for i := 0; i < numTablesVal; i++ {
		unitVersion := data.ReadUnsignedShort()
		length := int(data.ReadUnsignedShort())
		coverage := kerning.KernCoverage(data.ReadUnsignedShort())

		var pairs []kerning.KernPair
		bytesRead := 6

		if unitVersion == 0 {
			pairCount := int(data.ReadUnsignedShort())
			bytesRead += 2
			for j := 0; j < pairCount && bytesRead+6 <= length; j++ {
				left := int(data.ReadUnsignedShort())
				right := int(data.ReadUnsignedShort())
				adjustment := data.ReadSignedShort()
				pairs = append(pairs, kerning.NewKernPair(left, right, adjustment))
				bytesRead += 6
			}
		} else if unitVersion == 1 {
			_ = data.ReadUnsignedShort() // searchRange
			_ = data.ReadUnsignedShort() // entrySelector
			_ = data.ReadUnsignedShort() // rangeShift
			bytesRead += 6
			pairCount := int(data.ReadUnsignedShort())
			bytesRead += 2
			for j := 0; j < pairCount && bytesRead+10 <= length; j++ {
				left := int(data.ReadUnsignedShort())
				right := int(data.ReadUnsignedShort())
				_ = data.ReadUnsignedShort() // valueFormat
				adjustment := data.ReadSignedShort()
				pairs = append(pairs, kerning.NewKernPair(left, right, adjustment))
				_ = data.ReadUnsignedInt() // attributeOutput
				bytesRead += 10
			}
		}

		if len(pairs) > 0 {
			subTables = append(subTables, kerning.NewKerningSubTable(int(version), coverage, pairs))
		}

		remaining := length - bytesRead
		if remaining > 0 {
			data.Seek(int64(remaining), 1) // io.SeekCurrent
		}
	}

	return tables.NewKerningTable(subTables)
}
