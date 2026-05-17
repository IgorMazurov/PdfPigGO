package truetypeparser

import (
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
)

// Os2TableParser parses the OS/2 table of a TrueType font.
type Os2TableParser struct{}

// Parse reads and interprets the OS/2 table from raw TrueType data, returning
// the appropriate versioned table type based on the version field:
//   - Version 0 (length == 68): tables.Os2Table (Apple shortened format)
//   - Version 0 (with extra fields): tables.Os2RevisedVersion0Table
//   - Version 1: tables.Os2Version1Table
//   - Version 2-4: tables.Os2Version2To4OpenTypeTable
//   - Version >= 5: tables.Os2Version5OpenTypeTable
func (p *Os2TableParser) Parse(header truetype.TrueTypeHeaderTable, data *TrueTypeDataBytes, register *TableRegisterBuilder) (any, error) {
	if _, err := data.Seek(int64(header.Offset), 0); err != nil {
		return nil, err
	}

	version := data.ReadUnsignedShort()

	xAvgCharWidth := data.ReadSignedShort()
	weightClass := data.ReadUnsignedShort()
	widthClass := data.ReadUnsignedShort()
	typeFlags := data.ReadUnsignedShort()
	ySubscriptXSize := data.ReadSignedShort()
	ySubscriptYSize := data.ReadSignedShort()
	ySubscriptXOffset := data.ReadSignedShort()
	ySubscriptYOffset := data.ReadSignedShort()
	ySuperscriptXSize := data.ReadSignedShort()
	ySuperscriptYSize := data.ReadSignedShort()
	ySuperscriptXOffset := data.ReadSignedShort()
	ySuperscriptYOffset := data.ReadSignedShort()
	yStrikeoutSize := data.ReadSignedShort()
	yStrikeoutPosition := data.ReadSignedShort()
	familyClass := data.ReadSignedShort()
	panose := readByteArray(data, 10)
	ulCharRange1 := data.ReadUnsignedInt()
	ulCharRange2 := data.ReadUnsignedInt()
	ulCharRange3 := data.ReadUnsignedInt()
	ulCharRange4 := data.ReadUnsignedInt()
	vendorIdBytes := readByteArray(data, 4)
	selectionFlags := data.ReadUnsignedShort()
	firstCharacterIndex := data.ReadUnsignedShort()
	lastCharacterIndex := data.ReadUnsignedShort()

	unicodeCharRange := []uint32{ulCharRange1, ulCharRange2, ulCharRange3, ulCharRange4}
	vendorId := bytesToString(vendorIdBytes)

	baseArgs := baseArgsTuple{
		header, version, xAvgCharWidth,
		weightClass, widthClass, typeFlags, ySubscriptXSize,
		ySubscriptYSize, ySubscriptXOffset, ySubscriptYOffset,
		ySuperscriptXSize, ySuperscriptYSize,
		ySuperscriptXOffset, ySuperscriptYOffset,
		yStrikeoutSize, yStrikeoutPosition, familyClass,
		panose, unicodeCharRange, vendorId,
		selectionFlags, firstCharacterIndex, lastCharacterIndex,
	}

	/*
	 * Documentation for OS/2 version 0 in Apple's TrueType Reference Manual stops at the usLastCharIndex field
	 * and does not include the last five fields of the table as it was defined by Microsoft.
	 * Some legacy TrueType fonts may have been built with a shortened version 0 OS/2 table.
	 * Applications should check the table length for a version 0 OS/2 table before reading these fields.
	 */
	if version == 0 && header.Length == 68 {
		return newBaseOs2(baseArgs), nil
	}

	sTypoAscender, sTypoDescender, sTypoLineGap, usWinAscent, usWinDescent, ok := tryReadExtraFields(data)
	if !ok {
		/*
		 * Font may be invalid. Try falling back to shorter version...
		 */
		return newBaseOs2(baseArgs), nil
	}

	if version == 0 {
		return newRevisedV0(baseArgs, sTypoAscender, sTypoDescender, sTypoLineGap, usWinAscent, usWinDescent), nil
	}

	ulCodePageRange1 := data.ReadUnsignedInt()
	ulCodePageRange2 := data.ReadUnsignedInt()

	if version == 1 {
		return newV1(baseArgs, sTypoAscender, sTypoDescender, sTypoLineGap, usWinAscent, usWinDescent, ulCodePageRange1, ulCodePageRange2), nil
	}

	sxHeight := data.ReadSignedShort()
	sCapHeight := data.ReadSignedShort()
	usDefaultChar := data.ReadUnsignedShort()
	usBreakChar := data.ReadUnsignedShort()
	usMaxContext := data.ReadUnsignedShort()

	if version < 5 {
		return newV2To4(baseArgs, sTypoAscender, sTypoDescender, sTypoLineGap, usWinAscent, usWinDescent,
			ulCodePageRange1, ulCodePageRange2, sxHeight, sCapHeight, usDefaultChar, usBreakChar, usMaxContext), nil
	}

	usLowerOpticalPointSize := data.ReadUnsignedShort()
	usUpperOpticalPointSize := data.ReadUnsignedShort()

	return newV5(baseArgs, sTypoAscender, sTypoDescender, sTypoLineGap, usWinAscent, usWinDescent,
		ulCodePageRange1, ulCodePageRange2, sxHeight, sCapHeight, usDefaultChar, usBreakChar, usMaxContext,
		usLowerOpticalPointSize, usUpperOpticalPointSize), nil
}

// baseArgsTuple holds the common constructor arguments shared by all OS/2 table versions.
type baseArgsTuple struct {
	header              truetype.TrueTypeHeaderTable
	version             uint16
	xAvgCharWidth       int16
	weightClass         uint16
	widthClass          uint16
	typeFlags           uint16
	ySubscriptXSize     int16
	ySubscriptYSize     int16
	ySubscriptXOffset   int16
	ySubscriptYOffset   int16
	ySuperscriptXSize   int16
	ySuperscriptYSize   int16
	ySuperscriptXOffset int16
	ySuperscriptYOffset int16
	yStrikeoutSize      int16
	yStrikeoutPosition  int16
	familyClass         int16
	panose              []byte
	unicodeCharRange    []uint32
	vendorId            string
	selectionFlags      uint16
	firstCharacterIndex uint16
	lastCharacterIndex  uint16
}

// readByteArray reads length bytes from the data at the current offset and advances position.
func readByteArray(data *TrueTypeDataBytes, length int) []byte {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = data.data[data.offset+int64(i)]
	}
	data.offset += int64(length)
	return result
}

// bytesToString converts a byte slice to an ASCII string.
func bytesToString(b []byte) string {
	s := make([]byte, len(b))
	copy(s, b)
	return string(s)
}

// tryReadExtraFields attempts to read the five extra OS/2 fields (typo ascender/descender/gap,
// windows ascent/descent). Returns (values, true) on success or (zero values, false) on failure.
func tryReadExtraFields(data *TrueTypeDataBytes) (sTypoAscender int16, sTypoDescender int16, sTypoLineGap int16, usWinAscent uint16, usWinDescent uint16, ok bool) {
	required := 2 + 2 + 2 + 2 + 2 // five 2-byte fields
	if data.offset+int64(required) > int64(len(data.data)) {
		return 0, 0, 0, 0, 0, false
	}

	sTypoAscender = data.ReadSignedShort()
	sTypoDescender = data.ReadSignedShort()
	sTypoLineGap = data.ReadSignedShort()
	usWinAscent = data.ReadUnsignedShort()
	usWinDescent = data.ReadUnsignedShort()
	return sTypoAscender, sTypoDescender, sTypoLineGap, usWinAscent, usWinDescent, true
}

// newBaseOs2 constructs a base Os2Table from common arguments.
func newBaseOs2(a baseArgsTuple) tables.Os2Table {
	return tables.NewOs2Table(
		a.header, a.version, a.xAvgCharWidth,
		a.weightClass, a.widthClass, a.typeFlags,
		a.ySubscriptXSize, a.ySubscriptYSize,
		a.ySubscriptXOffset, a.ySubscriptYOffset,
		a.ySuperscriptXSize, a.ySuperscriptYSize,
		a.ySuperscriptXOffset, a.ySuperscriptYOffset,
		a.yStrikeoutSize, a.yStrikeoutPosition,
		a.familyClass, a.panose, a.unicodeCharRange,
		a.vendorId, a.selectionFlags,
		a.firstCharacterIndex, a.lastCharacterIndex,
	)
}

// newRevisedV0 constructs an Os2RevisedVersion0Table.
func newRevisedV0(a baseArgsTuple, sTypoAscender int16, sTypoDescender int16, sTypoLineGap int16, usWinAscent uint16, usWinDescent uint16) tables.Os2RevisedVersion0Table {
	return tables.NewOs2RevisedVersion0Table(
		a.header, a.version, a.xAvgCharWidth,
		a.weightClass, a.widthClass, a.typeFlags,
		a.ySubscriptXSize, a.ySubscriptYSize,
		a.ySubscriptXOffset, a.ySubscriptYOffset,
		a.ySuperscriptXSize, a.ySuperscriptYSize,
		a.ySuperscriptXOffset, a.ySuperscriptYOffset,
		a.yStrikeoutSize, a.yStrikeoutPosition,
		a.familyClass, a.panose, a.unicodeCharRange,
		a.vendorId, a.selectionFlags,
		a.firstCharacterIndex, a.lastCharacterIndex,
		sTypoAscender, sTypoDescender, sTypoLineGap,
		usWinAscent, usWinDescent,
	)
}

// newV1 constructs an Os2Version1Table.
func newV1(a baseArgsTuple, sTypoAscender int16, sTypoDescender int16, sTypoLineGap int16, usWinAscent uint16, usWinDescent uint16, ulCodePageRange1 uint32, ulCodePageRange2 uint32) tables.Os2Version1Table {
	return tables.NewOs2Version1Table(
		a.header, a.version, a.xAvgCharWidth,
		a.weightClass, a.widthClass, a.typeFlags,
		a.ySubscriptXSize, a.ySubscriptYSize,
		a.ySubscriptXOffset, a.ySubscriptYOffset,
		a.ySuperscriptXSize, a.ySuperscriptYSize,
		a.ySuperscriptXOffset, a.ySuperscriptYOffset,
		a.yStrikeoutSize, a.yStrikeoutPosition,
		a.familyClass, a.panose, a.unicodeCharRange,
		a.vendorId, a.selectionFlags,
		a.firstCharacterIndex, a.lastCharacterIndex,
		sTypoAscender, sTypoDescender, sTypoLineGap,
		usWinAscent, usWinDescent,
		ulCodePageRange1, ulCodePageRange2,
	)
}

// newV2To4 constructs an Os2Version2To4OpenTypeTable.
func newV2To4(a baseArgsTuple, sTypoAscender int16, sTypoDescender int16, sTypoLineGap int16, usWinAscent uint16, usWinDescent uint16, ulCodePageRange1 uint32, ulCodePageRange2 uint32, sxHeight int16, sCapHeight int16, usDefaultChar uint16, usBreakChar uint16, usMaxContext uint16) tables.Os2Version2To4OpenTypeTable {
	return tables.NewOs2Version2To4OpenTypeTable(
		a.header, a.version, a.xAvgCharWidth,
		a.weightClass, a.widthClass, a.typeFlags,
		a.ySubscriptXSize, a.ySubscriptYSize,
		a.ySubscriptXOffset, a.ySubscriptYOffset,
		a.ySuperscriptXSize, a.ySuperscriptYSize,
		a.ySuperscriptXOffset, a.ySuperscriptYOffset,
		a.yStrikeoutSize, a.yStrikeoutPosition,
		a.familyClass, a.panose, a.unicodeCharRange,
		a.vendorId, a.selectionFlags,
		a.firstCharacterIndex, a.lastCharacterIndex,
		sTypoAscender, sTypoDescender, sTypoLineGap,
		usWinAscent, usWinDescent,
		ulCodePageRange1, ulCodePageRange2,
		sxHeight, sCapHeight, usDefaultChar, usBreakChar, usMaxContext,
	)
}

// newV5 constructs an Os2Version5OpenTypeTable.
func newV5(a baseArgsTuple, sTypoAscender int16, sTypoDescender int16, sTypoLineGap int16, usWinAscent uint16, usWinDescent uint16, ulCodePageRange1 uint32, ulCodePageRange2 uint32, sxHeight int16, sCapHeight int16, usDefaultChar uint16, usBreakChar uint16, usMaxContext uint16, usLowerOpticalPointSize uint16, usUpperOpticalPointSize uint16) tables.Os2Version5OpenTypeTable {
	return tables.NewOs2Version5OpenTypeTable(
		a.header, a.version, a.xAvgCharWidth,
		a.weightClass, a.widthClass, a.typeFlags,
		a.ySubscriptXSize, a.ySubscriptYSize,
		a.ySubscriptXOffset, a.ySubscriptYOffset,
		a.ySuperscriptXSize, a.ySuperscriptYSize,
		a.ySuperscriptXOffset, a.ySuperscriptYOffset,
		a.yStrikeoutSize, a.yStrikeoutPosition,
		a.familyClass, a.panose, a.unicodeCharRange,
		a.vendorId, a.selectionFlags,
		a.firstCharacterIndex, a.lastCharacterIndex,
		sTypoAscender, sTypoDescender, sTypoLineGap,
		usWinAscent, usWinDescent,
		ulCodePageRange1, ulCodePageRange2,
		sxHeight, sCapHeight, usDefaultChar, usBreakChar, usMaxContext,
		usLowerOpticalPointSize, usUpperOpticalPointSize,
	)
}
