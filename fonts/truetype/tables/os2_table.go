package tables

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// IOs2Table is the common interface for all OS/2 table versions.
type IOs2Table interface {
	Tag() string
	Write(w io.Writer) error
}

// Os2Table contains the most basic format of the OS/2 table, excluding the fields
// not included in the Apple version of the specification.
type Os2Table struct {
	directoryTable         truetype.TrueTypeHeaderTable
	version                uint16
	xAverageCharacterWidth int16
	weightClass            uint16
	widthClass             uint16
	typeFlags              uint16
	ySubscriptXSize        int16
	ySubscriptYSize        int16
	ySubscriptXOffset      int16
	ySubscriptYOffset      int16
	ysuperscriptXSize      int16
	ysuperscriptYSize      int16
	ysuperscriptXOffset    int16
	ysuperscriptYOffset    int16
	yStrikeoutSize         int16
	yStrikeoutPosition     int16
	familyClass            int16
	panose                 []byte
	unicodeRanges          []uint32
	vendorId               string
	fontSelectionFlags     uint16
	firstCharacterIndex    uint16
	lastCharacterIndex     uint16
}

// NewOs2Table creates a new Os2Table.
func NewOs2Table(
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
) Os2Table {
	return Os2Table{
		directoryTable:         directoryTable,
		version:                version,
		xAverageCharacterWidth: xAverageCharacterWidth,
		weightClass:            weightClass,
		widthClass:             widthClass,
		typeFlags:              typeFlags,
		ySubscriptXSize:        ySubscriptXSize,
		ySubscriptYSize:        ySubscriptYSize,
		ySubscriptXOffset:      ySubscriptXOffset,
		ySubscriptYOffset:      ySubscriptYOffset,
		ysuperscriptXSize:      ysuperscriptXSize,
		ysuperscriptYSize:      ysuperscriptYSize,
		ysuperscriptXOffset:    ysuperscriptXOffset,
		ysuperscriptYOffset:    ysuperscriptYOffset,
		yStrikeoutSize:         yStrikeoutSize,
		yStrikeoutPosition:     yStrikeoutPosition,
		familyClass:            familyClass,
		panose:                 panose,
		unicodeRanges:          unicodeRanges,
		vendorId:               vendorId,
		fontSelectionFlags:     fontSelectionFlags,
		firstCharacterIndex:    firstCharacterIndex,
		lastCharacterIndex:     lastCharacterIndex,
	}
}

// Tag returns the 4-letter identifier for this table.
func (o Os2Table) Tag() string {
	return truetype.Os2
}

// DirectoryTable returns the directory entry from the font's offset subtable.
func (o Os2Table) DirectoryTable() truetype.TrueTypeHeaderTable {
	return o.directoryTable
}

// Version returns the version number 0-5 detailing the layout of the OS/2 table.
func (o Os2Table) Version() uint16 {
	return o.version
}

// XAverageCharacterWidth returns the average width of all non-zero width characters in the font.
func (o Os2Table) XAverageCharacterWidth() int16 {
	return o.xAverageCharacterWidth
}

// WeightClass indicates the visual weight of characters in the font from 1-1000.
func (o Os2Table) WeightClass() uint16 {
	return o.weightClass
}

// WidthClass returns the percentage difference from normal of the aspect ratio for this font.
func (o Os2Table) WidthClass() uint16 {
	return o.widthClass
}

// TypeFlags returns the font embedding licensing rights for this font.
func (o Os2Table) TypeFlags() uint16 {
	return o.typeFlags
}

// YSubscriptXSize returns the recommended horizontal size for subscripts using this font.
func (o Os2Table) YSubscriptXSize() int16 {
	return o.ySubscriptXSize
}

// YSubscriptYSize returns the recommended vertical size for subscripts using this font.
func (o Os2Table) YSubscriptYSize() int16 {
	return o.ySubscriptYSize
}

// YSubscriptXOffset returns the recommended horizontal offset from the previous glyph origin
// to the subscript's origin for subscripts using this font.
func (o Os2Table) YSubscriptXOffset() int16 {
	return o.ySubscriptXOffset
}

// YSubscriptYOffset returns the recommended vertical offset from the previous glyph origin
// to the subscript's origin for subscripts using this font.
func (o Os2Table) YSubscriptYOffset() int16 {
	return o.ySubscriptYOffset
}

// YSuperscriptXSize returns the recommended horizontal size for superscripts using this font.
func (o Os2Table) YSuperscriptXSize() int16 {
	return o.ysuperscriptXSize
}

// YSuperscriptYSize returns the recommended vertical size for superscripts using this font.
func (o Os2Table) YSuperscriptYSize() int16 {
	return o.ysuperscriptYSize
}

// YSuperscriptXOffset returns the recommended horizontal offset from the previous glyph origin
// to the superscript's origin for superscripts using this font.
func (o Os2Table) YSuperscriptXOffset() int16 {
	return o.ysuperscriptXOffset
}

// YSuperscriptYOffset returns the recommended vertical offset from the previous glyph origin
// to the superscript's origin for superscripts using this font.
func (o Os2Table) YSuperscriptYOffset() int16 {
	return o.ysuperscriptYOffset
}

// YStrikeoutSize returns the thickness of the strikeout stroke.
func (o Os2Table) YStrikeoutSize() int16 {
	return o.yStrikeoutSize
}

// YStrikeoutPosition returns the position of the top of the strikeout stroke relative
// to the baseline. Positive values are above the baseline, negative values below.
func (o Os2Table) YStrikeoutPosition() int16 {
	return o.yStrikeoutPosition
}

// FamilyClass returns the value registered by IBM for each font family to find substitutes.
// The high byte is the family class, the low byte is the family subclass.
func (o Os2Table) FamilyClass() int16 {
	return o.familyClass
}

// Panose returns the PANOSE definition of 10 bytes that defines various information about
// the font enabling matching fonts based on requirements. The first byte is the family type
// (Latin, Latin Hand Written, etc.).
func (o Os2Table) Panose() []byte {
	return o.panose
}

// UnicodeRanges specifies Unicode blocks supported by the font file for the Microsoft platform.
func (o Os2Table) UnicodeRanges() []uint32 {
	return o.unicodeRanges
}

// VendorId returns the four-character identifier for the vendor of the given type face.
func (o Os2Table) VendorId() string {
	return o.vendorId
}

// FontSelectionFlags contains information concerning the nature of the font patterns.
func (o Os2Table) FontSelectionFlags() uint16 {
	return o.fontSelectionFlags
}

// FirstCharacterIndex returns the minimum Unicode character code in this font.
func (o Os2Table) FirstCharacterIndex() uint16 {
	return o.firstCharacterIndex
}

// LastCharacterIndex returns the maximum Unicode character code in this font.
func (o Os2Table) LastCharacterIndex() uint16 {
	return o.lastCharacterIndex
}

// Write writes the OS/2 table data to the writer.
func (o Os2Table) Write(w io.Writer) error {
	if _, err := core.WriteUShort(w, o.version); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.xAverageCharacterWidth); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, o.weightClass); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, o.widthClass); err != nil {
		return err
	}

	if _, err := core.WriteShortUint16(w, uint16(o.typeFlags)); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.ySubscriptXSize); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.ySubscriptYSize); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.ySubscriptXOffset); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.ySubscriptYOffset); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.ysuperscriptXSize); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.ysuperscriptYSize); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.ysuperscriptXOffset); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.ysuperscriptYOffset); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.yStrikeoutSize); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.yStrikeoutPosition); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.familyClass); err != nil {
		return err
	}

	if _, err := w.Write(o.panose); err != nil {
		return err
	}

	for _, r := range o.unicodeRanges {
		if _, err := core.WriteUInt(w, r); err != nil {
			return err
		}
	}

	for i := 0; i < len(o.vendorId); i++ {
		if _, err := w.Write([]byte{o.vendorId[i]}); err != nil {
			return err
		}
	}

	if _, err := core.WriteUShort(w, o.fontSelectionFlags); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, o.firstCharacterIndex); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, o.lastCharacterIndex); err != nil {
		return err
	}

	return nil
}
