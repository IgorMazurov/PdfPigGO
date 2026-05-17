package tables

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// Os2Version2To4OpenTypeTable represents version 2-4 of the OS/2 table as defined in OpenType 1.5,
// with five additional fields beyond those in version 1: xHeight, capHeight, defaultCharacter,
// breakCharacter, and maximumContext. Although no new fields were added beyond version 2,
// the specification of certain fields was revised in versions 3 and 4.
type Os2Version2To4OpenTypeTable struct {
	Os2Version1Table
	xHeight          int16
	capHeight        int16
	defaultCharacter uint16
	breakCharacter   uint16
	maximumContext   uint16
}

// NewOs2Version2To4OpenTypeTable creates a new Os2Version2To4OpenTypeTable.
func NewOs2Version2To4OpenTypeTable(
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
) Os2Version2To4OpenTypeTable {
	return Os2Version2To4OpenTypeTable{
		Os2Version1Table: NewOs2Version1Table(
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
		),
		xHeight:          xHeight,
		capHeight:        capHeight,
		defaultCharacter: defaultCharacter,
		breakCharacter:   breakCharacter,
		maximumContext:   maximumContext,
	}
}

// XHeight returns the distance between the baseline and the approximate height
// of non-ascending lowercase letters.
func (o Os2Version2To4OpenTypeTable) XHeight() int16 {
	return o.xHeight
}

// CapHeight returns the distance between the baseline and the approximate height
// of uppercase letters.
func (o Os2Version2To4OpenTypeTable) CapHeight() int16 {
	return o.capHeight
}

// DefaultCharacter returns the Unicode code point in UTF-16 encoding of a character
// that can be used for a default glyph if a requested character is not supported.
// A value of zero means glyph ID 0 should be used as the default character.
func (o Os2Version2To4OpenTypeTable) DefaultCharacter() uint16 {
	return o.defaultCharacter
}

// BreakCharacter returns the Unicode code point in UTF-16 encoding of a character
// that can be used as a default break character for word separation and text justification.
// Most fonts specify U+0020 SPACE as the break character.
func (o Os2Version2To4OpenTypeTable) BreakCharacter() uint16 {
	return o.breakCharacter
}

// MaximumContext returns the maximum distance in glyphs that any feature of this font
// is capable of affecting. For example, kerning has a value of 2 (one for each glyph
// in the kerning pair), and fonts with the 'ffi' ligature would have a value of 3.
func (o Os2Version2To4OpenTypeTable) MaximumContext() uint16 {
	return o.maximumContext
}

// Write writes the OS/2 version 2-4 table data to the writer.
func (o Os2Version2To4OpenTypeTable) Write(w io.Writer) error {
	if err := o.Os2Version1Table.Write(w); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.xHeight); err != nil {
		return err
	}

	if _, err := core.WriteShort(w, o.capHeight); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, o.defaultCharacter); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, o.breakCharacter); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, o.maximumContext); err != nil {
		return err
	}

	return nil
}
