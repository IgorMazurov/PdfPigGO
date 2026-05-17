package subsetting

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts/truetype/glyphs"
)

// TrueTypeSubsetGlyphTable holds details of the new glyph table created when
// subsetting a TrueType font.
type TrueTypeSubsetGlyphTable struct {
	// Bytes is the raw data of the subsetted glyf table.
	Bytes []byte

	// GlyphOffsets contains the offset to each glyph within Bytes.
	GlyphOffsets []uint32

	// HorizontalMetrics holds the horizontal metrics for each glyph in the subset.
	HorizontalMetrics []glyphs.HorizontalMetric
}

// NewTrueTypeSubsetGlyphTable creates a new TrueTypeSubsetGlyphTable.
func NewTrueTypeSubsetGlyphTable(
	bytes []byte,
	glyphOffsets []uint32,
	horizontalMetrics []glyphs.HorizontalMetric,
) *TrueTypeSubsetGlyphTable {
	return &TrueTypeSubsetGlyphTable{
		Bytes:             bytes,
		GlyphOffsets:      glyphOffsets,
		HorizontalMetrics: horizontalMetrics,
	}
}

// GlyphCount returns the number of glyphs in the subsetted table.
func (t *TrueTypeSubsetGlyphTable) GlyphCount() uint16 {
	return uint16(len(t.GlyphOffsets) - 1)
}

// OffsetsAsLongs converts the glyph offsets to int64 values.
func (t *TrueTypeSubsetGlyphTable) OffsetsAsLongs() []int64 {
	data := make([]int64, len(t.GlyphOffsets))
	for i, offset := range t.GlyphOffsets {
		data[i] = int64(offset)
	}
	return data
}

// String returns a human-readable description of the subset glyph table.
func (t *TrueTypeSubsetGlyphTable) String() string {
	return fmt.Sprintf("%d glyphs. Data is %d bytes.", int(t.GlyphCount()), len(t.Bytes))
}
