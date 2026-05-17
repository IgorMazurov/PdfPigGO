package tables

import (
	"errors"
	"io"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/glyphs"
)

// HorizontalMetricsTable contains metric information for the horizontal layout
// of each glyph in the font. The 'hmtx' table stores advance widths and left-side bearings.
type HorizontalMetricsTable struct {
	directoryTable             truetype.TrueTypeHeaderTable
	horizontalMetrics          []glyphs.HorizontalMetric
	additionalLeftSideBearings []int16
}

// NewHorizontalMetricsTable creates a new HorizontalMetricsTable.
func NewHorizontalMetricsTable(
	directoryTable truetype.TrueTypeHeaderTable,
	horizontalMetrics []glyphs.HorizontalMetric,
	additionalLeftSideBearings []int16,
) (HorizontalMetricsTable, error) {
	if horizontalMetrics == nil {
		return HorizontalMetricsTable{}, errors.New("horizontalMetrics cannot be null")
	}

	if additionalLeftSideBearings == nil {
		return HorizontalMetricsTable{}, errors.New("additionalLeftSideBearings cannot be null")
	}

	return HorizontalMetricsTable{
		directoryTable:             directoryTable,
		horizontalMetrics:          horizontalMetrics,
		additionalLeftSideBearings: additionalLeftSideBearings,
	}, nil
}

// Tag returns the 4-letter identifier for this table.
func (h HorizontalMetricsTable) Tag() string {
	return truetype.Hmtx
}

// DirectoryTable returns the directory entry from the font's offset subtable.
func (h HorizontalMetricsTable) DirectoryTable() truetype.TrueTypeHeaderTable {
	return h.directoryTable
}

// HorizontalMetrics returns the left-side bearing and advance widths for the glyphs
// in the font. For a monospace font this may only contain a single entry with additional
// left-side bearings defined in AdditionalLeftSideBearings.
func (h HorizontalMetricsTable) HorizontalMetrics() []glyphs.HorizontalMetric {
	return h.horizontalMetrics
}

// AdditionalLeftSideBearings returns the array of left side bearings following the
// horizontal metrics. This is generally used for a run of monospaced glyphs where
// each glyph shares the advance width from the last entry in HorizontalMetrics.
func (h HorizontalMetricsTable) AdditionalLeftSideBearings() []int16 {
	return h.additionalLeftSideBearings
}

// GetAdvanceWidth returns the advance width for a glyph at the given index.
// For monospaced fonts without a width entry per glyph, the last metric's
// advance width is returned for subsequent glyphs.
func (h HorizontalMetricsTable) GetAdvanceWidth(index int) (uint16, error) {
	if index < 0 {
		return 0, errors.New("index cannot be less than zero")
	}

	if index < len(h.horizontalMetrics) {
		return h.horizontalMetrics[index].AdvanceWidth, nil
	}

	return h.horizontalMetrics[len(h.horizontalMetrics)-1].AdvanceWidth, nil
}

// Write serializes the horizontal metrics table to the output stream.
func (h HorizontalMetricsTable) Write(w io.Writer) error {
	for i := 0; i < len(h.horizontalMetrics); i++ {
		metric := h.horizontalMetrics[i]
		if _, err := core.WriteUShort(w, metric.AdvanceWidth); err != nil {
			return err
		}
		if _, err := core.WriteShort(w, metric.LeftSideBearing); err != nil {
			return err
		}
	}

	for i := 0; i < len(h.additionalLeftSideBearings); i++ {
		if _, err := core.WriteShort(w, h.additionalLeftSideBearings[i]); err != nil {
			return err
		}
	}

	return nil
}

// GetHorizontalMetric returns the full horizontal metric (advance width + left side bearing)
// for a given glyph index. For indices beyond the explicit metrics array, the last advance
// width is combined with either an additional left-side bearing or zero.
func (h HorizontalMetricsTable) GetHorizontalMetric(glyphIndex int) glyphs.HorizontalMetric {
	if glyphIndex < len(h.horizontalMetrics) {
		return h.horizontalMetrics[glyphIndex]
	}

	lastWidth := uint16(0)
	if len(h.horizontalMetrics) > 0 {
		lastWidth = h.horizontalMetrics[len(h.horizontalMetrics)-1].AdvanceWidth
	}

	lsbIndex := glyphIndex - len(h.horizontalMetrics)
	if lsbIndex < len(h.additionalLeftSideBearings) {
		return glyphs.NewHorizontalMetric(lastWidth, h.additionalLeftSideBearings[lsbIndex])
	}

	return glyphs.NewHorizontalMetric(lastWidth, 0)
}

// Compile-time interface assertions.
var (
	_ TrueTypeTable = HorizontalMetricsTable{}
	_ core.Writeable = HorizontalMetricsTable{}
)
