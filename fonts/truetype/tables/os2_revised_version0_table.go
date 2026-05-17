package tables

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// Os2RevisedVersion0Table represents version 0 of the OS/2 table as defined in TrueType revision 1.5,
// including fields not present in the Apple specification.
type Os2RevisedVersion0Table struct {
	Os2Table
	typographicAscender int16
	typographicDescender int16
	typographicLineGap   int16
	windowsAscent       uint16
	windowsDescent      uint16
}

// NewOs2RevisedVersion0Table creates a new Os2RevisedVersion0Table.
func NewOs2RevisedVersion0Table(
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
) Os2RevisedVersion0Table {
	return Os2RevisedVersion0Table{
		Os2Table: NewOs2Table(
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
		),
		typographicAscender: typographicAscender,
		typographicDescender: typographicDescender,
		typographicLineGap:   typographicLineGap,
		windowsAscent:        windowsAscent,
		windowsDescent:       windowsDescent,
	}
}

// TypographicAscender returns the typographic ascender.
func (o Os2RevisedVersion0Table) TypographicAscender() int16 {
	return o.typographicAscender
}

// TypographicDescender returns the typographic descender.
func (o Os2RevisedVersion0Table) TypographicDescender() int16 {
	return o.typographicDescender
}

// TypographicLineGap returns the typographic line gap.
func (o Os2RevisedVersion0Table) TypographicLineGap() int16 {
	return o.typographicLineGap
}

// WindowsAscent returns the Windows ascender metric for specifying
// the height above the baseline for a clipping region.
func (o Os2RevisedVersion0Table) WindowsAscent() uint16 {
	return o.windowsAscent
}

// WindowsDescent returns the Windows descender metric for specifying
// the vertical extent below the baseline for a clipping region.
func (o Os2RevisedVersion0Table) WindowsDescent() uint16 {
	return o.windowsDescent
}

// Write writes the OS/2 version 0 table data to the writer.
func (o Os2RevisedVersion0Table) Write(w io.Writer) error {
	if err := o.Os2Table.Write(w); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.typographicAscender); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.typographicDescender); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.typographicLineGap); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, o.windowsAscent); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, o.windowsDescent); err != nil {
		return err
	}

	return nil
}
