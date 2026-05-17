package tables

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// HeaderTable contains global information about the font.
type HeaderTable struct {
	directoryTable        truetype.TrueTypeHeaderTable
	version               float32
	fontRevision          float32
	checkSumAdjustment    uint32
	magicNumber           uint32
	flags                 uint16
	unitPerEm             uint16
	created               int64
	modified              int64
	xMin                  int16
	yMin                  int16
	xMax                  int16
	yMax                  int16
	macStyle              HeaderMacStyle
	lowestRecommendedPpem uint16
	fontDirectionHint     FontDirection
	indexToLocFormat      IndexToLocationTableEntryFormat
	glyphDataFormat       int16
}

// NewHeaderTable creates a new HeaderTable.
func NewHeaderTable(
	directoryTable truetype.TrueTypeHeaderTable,
	version float32,
	fontRevision float32,
	checkSumAdjustment uint32,
	magicNumber uint32,
	flags uint16,
	unitPerEm uint16,
	created int64,
	modified int64,
	xMin int16,
	yMin int16,
	xMax int16,
	yMax int16,
	macStyle HeaderMacStyle,
	lowestRecommendedPpem uint16,
	fontDirectionHint FontDirection,
	indexToLocFormat IndexToLocationTableEntryFormat,
	glyphDataFormat int16,
) HeaderTable {
	return HeaderTable{
		directoryTable:        directoryTable,
		version:               version,
		fontRevision:          fontRevision,
		checkSumAdjustment:    checkSumAdjustment,
		magicNumber:           magicNumber,
		flags:                 flags,
		unitPerEm:             unitPerEm,
		created:               created,
		modified:              modified,
		xMin:                  xMin,
		yMin:                  yMin,
		xMax:                  xMax,
		yMax:                  yMax,
		macStyle:              macStyle,
		lowestRecommendedPpem: lowestRecommendedPpem,
		fontDirectionHint:     fontDirectionHint,
		indexToLocFormat:      indexToLocFormat,
		glyphDataFormat:       glyphDataFormat,
	}
}

// Tag returns the 4-letter identifier for this table.
func (h HeaderTable) Tag() string {
	return truetype.Head
}

// DirectoryTable returns the directory entry from the font's offset subtable.
func (h HeaderTable) DirectoryTable() truetype.TrueTypeHeaderTable {
	return h.directoryTable
}

// Version returns the font version number.
func (h HeaderTable) Version() float32 { return h.version }

// FontRevision returns the font revision number.
func (h HeaderTable) FontRevision() float32 { return h.fontRevision }

// CheckSumAdjustment returns the checksum adjustment value used to derive
// the checksum of the entire TrueType file.
func (h HeaderTable) CheckSumAdjustment() uint32 { return h.checkSumAdjustment }

// MagicNumber returns the magic number (should be 0x5F0F3CF5).
func (h HeaderTable) MagicNumber() uint32 { return h.magicNumber }

// Flags returns the font flags bit field.
func (h HeaderTable) Flags() uint16 { return h.flags }

// UnitPerEm returns the number of units per em.
func (h HeaderTable) UnitPerEm() uint16 { return h.unitPerEm }

// Created returns the font creation date-time as seconds since 1904-01-01 UTC.
func (h HeaderTable) Created() int64 { return h.created }

// Modified returns the last modification date-time as seconds since 1904-01-01 UTC.
func (h HeaderTable) Modified() int64 { return h.modified }

// Bounds returns the minimum rectangle which contains all glyphs.
func (h HeaderTable) Bounds() core.PdfRectangle {
	return core.NewPdfRectangleFromInt(int(h.xMin), int(h.yMin), int(h.xMax), int(h.yMax))
}

// XMin returns the minimum x coordinate across all glyphs.
func (h HeaderTable) XMin() int16 { return h.xMin }

// YMin returns the minimum y coordinate across all glyphs.
func (h HeaderTable) YMin() int16 { return h.yMin }

// XMax returns the maximum x coordinate across all glyphs.
func (h HeaderTable) XMax() int16 { return h.xMax }

// YMax returns the maximum y coordinate across all glyphs.
func (h HeaderTable) YMax() int16 { return h.yMax }

// MacStyle returns the MacStyle flags.
func (h HeaderTable) MacStyle() HeaderMacStyle { return h.macStyle }

// LowestRecommendedPpem returns the smallest readable size in pixels.
func (h HeaderTable) LowestRecommendedPpem() uint16 { return h.lowestRecommendedPpem }

// FontDirectionHint returns the font direction hint.
func (h HeaderTable) FontDirectionHint() FontDirection { return h.fontDirectionHint }

// IndexToLocFormat returns the index-to-location format: 0 for short offsets, 1 for long.
func (h HeaderTable) IndexToLocFormat() IndexToLocationTableEntryFormat { return h.indexToLocFormat }

// GlyphDataFormat returns the glyph data format. 0 for current format.
func (h HeaderTable) GlyphDataFormat() int16 { return h.glyphDataFormat }

// FontDirection represents values of the font direction hint.
type FontDirection int16

const (
	// StronglyRightToLeftWithNeutrals means strongly right to left with neutrals.
	StronglyRightToLeftWithNeutrals FontDirection = -2
	// StronglyRightToLeft means strongly right to left.
	StronglyRightToLeft FontDirection = -1
	// FullyMixedDirectional means full mixed directional glyphs.
	FullyMixedDirectional FontDirection = 0
	// StronglyLeftToRight means strongly left to right.
	StronglyLeftToRight FontDirection = 1
	// StronglyLeftToRightWithNeutrals means strongly left to right with neutrals.
	StronglyLeftToRightWithNeutrals FontDirection = 2
)

// HeaderMacStyle represents values of the Mac Style flag in the header table.
type HeaderMacStyle uint16

const (
	// HeaderMacStyleNone means no flags set.
	HeaderMacStyleNone HeaderMacStyle = 0
	// HeaderMacStyleBold means bold.
	HeaderMacStyleBold HeaderMacStyle = 1 << 0
	// HeaderMacStyleItalic means italic.
	HeaderMacStyleItalic HeaderMacStyle = 1 << 1
	// HeaderMacStyleUnderline means underline.
	HeaderMacStyleUnderline HeaderMacStyle = 1 << 2
	// HeaderMacStyleOutline means outline.
	HeaderMacStyleOutline HeaderMacStyle = 1 << 3
	// HeaderMacStyleShadow means shadow.
	HeaderMacStyleShadow HeaderMacStyle = 1 << 4
	// HeaderMacStyleCondensed means condensed (narrow).
	HeaderMacStyleCondensed HeaderMacStyle = 1 << 5
	// HeaderMacStyleExtended means extended.
	HeaderMacStyleExtended HeaderMacStyle = 1 << 6
)
