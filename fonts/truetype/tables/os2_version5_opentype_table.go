package tables

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// Os2Version5OpenTypeTable represents version 5 of the OS/2 table as defined in OpenType 1.7,
// with two additional fields beyond those in versions 2-4: lowerOpticalPointSize and upperOpticalPointSize.
type Os2Version5OpenTypeTable struct {
	Os2Version2To4OpenTypeTable
	lowerOpticalPointSize uint16
	upperOpticalPointSize uint16
}

// NewOs2Version5OpenTypeTable creates a new Os2Version5OpenTypeTable.
func NewOs2Version5OpenTypeTable(
	directoryTable truetype.TrueTypeHeaderTable,
	version uint16,
	xAverageCharacterWidth int16,
	weightClass uint16,
	widthClass uint16,
	typeFlags uint16,
	ySubscriptXSize int16,
	ySubscriptYSize int16,
	ySubscriptXOffset int16,
	ySubscriptYOffset int16,
	ysuperscriptXSize int16,
	ysuperscriptYSize int16,
	ysuperscriptXOffset int16,
	ysuperscriptYOffset int16,
	yStrikeoutSize int16,
	yStrikeoutPosition int16,
	familyClass int16,
	panose []byte,
	unicodeRanges []uint32,
	vendorId string,
	fontSelectionFlags uint16,
	firstCharacterIndex uint16,
	lastCharacterIndex uint16,
	typographicAscender int16,
	typographicDescender int16,
	typographicLineGap int16,
	windowsAscent uint16,
	windowsDescent uint16,
	codePage1 uint32,
	codePage2 uint32,
	xHeight int16,
	capHeight int16,
	defaultCharacter uint16,
	breakCharacter uint16,
	maximumContext uint16,
	lowerOpticalPointSize uint16,
	upperOpticalPointSize uint16,
) Os2Version5OpenTypeTable {
	return Os2Version5OpenTypeTable{
		Os2Version2To4OpenTypeTable: NewOs2Version2To4OpenTypeTable(
			directoryTable,
			version,
			xAverageCharacterWidth,
			weightClass,
			widthClass,
			typeFlags,
			ySubscriptXSize,
			ySubscriptYSize,
			ySubscriptXOffset,
			ySubscriptYOffset,
			ysuperscriptXSize,
			ysuperscriptYSize,
			ysuperscriptXOffset,
			ysuperscriptYOffset,
			yStrikeoutSize,
			yStrikeoutPosition,
			familyClass,
			panose,
			unicodeRanges,
			vendorId,
			fontSelectionFlags,
			firstCharacterIndex,
			lastCharacterIndex,
			typographicAscender,
			typographicDescender,
			typographicLineGap,
			windowsAscent,
			windowsDescent,
			codePage1,
			codePage2,
			xHeight,
			capHeight,
			defaultCharacter,
			breakCharacter,
			maximumContext,
		),
		lowerOpticalPointSize: lowerOpticalPointSize,
		upperOpticalPointSize: upperOpticalPointSize,
	}
}

// LowerOpticalPointSize returns the lower value of the size range for which this font has been designed.
// The units are TWIPs (one-twentieth of a point, or 1440 per inch). This is the inclusive lower bound.
func (o Os2Version5OpenTypeTable) LowerOpticalPointSize() uint16 {
	return o.lowerOpticalPointSize
}

// UpperOpticalPointSize returns the upper value of the size range for which this font has been designed.
// The units are TWIPs (one-twentieth of a point, or 1440 per inch). This is the exclusive upper bound.
func (o Os2Version5OpenTypeTable) UpperOpticalPointSize() uint16 {
	return o.upperOpticalPointSize
}

// Write writes the OS/2 version 5 table data to the writer.
func (o Os2Version5OpenTypeTable) Write(w io.Writer) error {
	if err := o.Os2Version2To4OpenTypeTable.Write(w); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, o.lowerOpticalPointSize); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, o.upperOpticalPointSize); err != nil {
		return err
	}

	return nil
}
