package tables

import (
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// HorizontalHeaderTable contains information needed to layout fonts whose
// characters are written horizontally, that is, either left to right or right to left.
type HorizontalHeaderTable struct {
	directoryTable        truetype.TrueTypeHeaderTable
	majorVersion          int
	minorVersion          int
	ascent                int16
	descent               int16
	lineGap               int16
	advanceWidthMaximum   uint16
	minimumLeftSideBearing int16
	minimumRightSideBearing int16
	xMaxExtent            int16
	caretSlopeRise        int16
	caretSlopeRun         int16
	caretOffset           int16
	metricDataFormat      int16
	numberOfHeaderMetrics uint16
}

// NewHorizontalHeaderTable creates a new HorizontalHeaderTable.
func NewHorizontalHeaderTable(
	directoryTable truetype.TrueTypeHeaderTable,
	majorVersion, minorVersion int,
	ascent, descent, lineGap int16,
	advanceWidthMaximum uint16,
	minimumLeftSideBearing, minimumRightSideBearing int16,
	xMaxExtent, caretSlopeRise, caretSlopeRun, caretOffset, metricDataFormat int16,
	numberOfHeaderMetrics uint16,
) HorizontalHeaderTable {
	return HorizontalHeaderTable{
		directoryTable:          directoryTable,
		majorVersion:            majorVersion,
		minorVersion:            minorVersion,
		ascent:                  ascent,
		descent:                 descent,
		lineGap:                 lineGap,
		advanceWidthMaximum:     advanceWidthMaximum,
		minimumLeftSideBearing:  minimumLeftSideBearing,
		minimumRightSideBearing: minimumRightSideBearing,
		xMaxExtent:              xMaxExtent,
		caretSlopeRise:          caretSlopeRise,
		caretSlopeRun:           caretSlopeRun,
		caretOffset:             caretOffset,
		metricDataFormat:        metricDataFormat,
		numberOfHeaderMetrics:   numberOfHeaderMetrics,
	}
}

// Tag returns the 4-letter identifier for this table.
func (h HorizontalHeaderTable) Tag() string {
	return truetype.Hhea
}

// DirectoryTable returns the directory entry from the font's offset subtable.
func (h HorizontalHeaderTable) DirectoryTable() truetype.TrueTypeHeaderTable {
	return h.directoryTable
}

// MajorVersion returns the major version number of this table (1).
func (h HorizontalHeaderTable) MajorVersion() int {
	return h.majorVersion
}

// MinorVersion returns the minor version number of this table (0).
func (h HorizontalHeaderTable) MinorVersion() int {
	return h.minorVersion
}

// Ascent returns the distance from baseline to highest ascender.
func (h HorizontalHeaderTable) Ascent() int16 {
	return h.ascent
}

// Descent returns the distance from baseline to lower descender.
func (h HorizontalHeaderTable) Descent() int16 {
	return h.descent
}

// LineGap returns the typographic line gap.
func (h HorizontalHeaderTable) LineGap() int16 {
	return h.lineGap
}

// AdvanceWidthMaximum returns the maximum advance width value as given by
// the Horizontal Metrics table.
func (h HorizontalHeaderTable) AdvanceWidthMaximum() uint16 {
	return h.advanceWidthMaximum
}

// MinimumLeftSideBearing returns the minimum left side bearing as given by
// the Horizontal Metrics table.
func (h HorizontalHeaderTable) MinimumLeftSideBearing() int16 {
	return h.minimumLeftSideBearing
}

// MinimumRightSideBearing returns the minimum right sidebearing.
func (h HorizontalHeaderTable) MinimumRightSideBearing() int16 {
	return h.minimumRightSideBearing
}

// XMaxExtent returns the maximum X extent.
func (h HorizontalHeaderTable) XMaxExtent() int16 {
	return h.xMaxExtent
}

// CaretSlopeRise returns the value used to calculate the slope of the cursor.
// 1 is vertical.
func (h HorizontalHeaderTable) CaretSlopeRise() int16 {
	return h.caretSlopeRise
}

// CaretSlopeRun returns the caret slope run value. 0 is vertical.
func (h HorizontalHeaderTable) CaretSlopeRun() int16 {
	return h.caretSlopeRun
}

// CaretOffset returns the amount by which a slanted highlight on a glyph should
// be shifted to provide the best appearance. 0 for non-slanted fonts.
func (h HorizontalHeaderTable) CaretOffset() int16 {
	return h.caretOffset
}

// MetricDataFormat returns the metric data format. 0 for the current format.
func (h HorizontalHeaderTable) MetricDataFormat() int16 {
	return h.metricDataFormat
}

// NumberOfHeaderMetrics returns the number of horizontal metrics in the
// Horizontal Metrics table.
func (h HorizontalHeaderTable) NumberOfHeaderMetrics() uint16 {
	return h.numberOfHeaderMetrics
}
