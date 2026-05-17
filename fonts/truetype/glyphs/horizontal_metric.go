package glyphs

import "fmt"

// HorizontalMetric holds the pair of horizontal metrics for an individual glyph.
type HorizontalMetric struct {
	// AdvanceWidth is the advance width of the glyph.
	AdvanceWidth uint16

	// LeftSideBearing is the left side bearing of the glyph.
	LeftSideBearing int16
}

// NewHorizontalMetric creates a new HorizontalMetric with the given values.
func NewHorizontalMetric(advanceWidth uint16, leftSideBearing int16) HorizontalMetric {
	return HorizontalMetric{
		AdvanceWidth:    advanceWidth,
		LeftSideBearing: leftSideBearing,
	}
}

// String returns a string representation of the horizontal metric.
func (m HorizontalMetric) String() string {
	return fmt.Sprintf("Width: %d. LSB: %d.", m.AdvanceWidth, m.LeftSideBearing)
}
