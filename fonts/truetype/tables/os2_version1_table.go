package tables

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// Os2Version1Table represents version 1 of the OS/2 table as defined in TrueType revision 1.66,
// with two additional fields beyond those in version 0: codePage1 and codePage2.
type Os2Version1Table struct {
	Os2RevisedVersion0Table
	codePage1 uint32
	codePage2 uint32
}

// NewOs2Version1Table creates a new Os2Version1Table.
func NewOs2Version1Table(
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
) Os2Version1Table {
	return Os2Version1Table{
		Os2RevisedVersion0Table: NewOs2RevisedVersion0Table(
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
		),
		codePage1: codePage1,
		codePage2: codePage2,
	}
}

// CodePage1 returns the first code page range bit mask for the cmap subtable.
func (o Os2Version1Table) CodePage1() uint32 {
	return o.codePage1
}

// CodePage2 returns the second code page range bit mask for the cmap subtable.
func (o Os2Version1Table) CodePage2() uint32 {
	return o.codePage2
}

// Write writes the OS/2 version 1 table data to the writer.
func (o Os2Version1Table) Write(w io.Writer) error {
	if err := o.Os2RevisedVersion0Table.Write(w); err != nil {
		return err
	}

	if _, err := core.WriteUInt(w, o.codePage1); err != nil {
		return err
	}

	if _, err := core.WriteUInt(w, o.codePage2); err != nil {
		return err
	}

	return nil
}
