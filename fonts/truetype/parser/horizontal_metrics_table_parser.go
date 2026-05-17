package truetypeparser

import (
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/glyphs"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
)

// HorizontalMetricsTableParser parses the horizontal metrics (hmtx) table of a TrueType font.
type HorizontalMetricsTableParser struct{}

// Parse reads and interprets the horizontal metrics table from raw TrueType data.
func (p *HorizontalMetricsTableParser) Parse(header truetype.TrueTypeHeaderTable, data *TrueTypeDataBytes, register *TableRegisterBuilder) (tables.HorizontalMetricsTable, error) {
	glyphCount := register.MaximumProfileTable.NumberOfGlyphs()
	metricCount := int(register.HorizontalHeaderTable.NumberOfHeaderMetrics())

	if _, err := data.Seek(int64(header.Offset), 0); err != nil {
		return tables.HorizontalMetricsTable{}, err
	}

	horizontalMetrics := make([]glyphs.HorizontalMetric, metricCount)

	for i := 0; i < metricCount; i++ {
		width := data.ReadUnsignedShort()
		lsb := data.ReadSignedShort()
		horizontalMetrics[i] = glyphs.NewHorizontalMetric(width, lsb)
	}

	numberNonHorizontal := glyphCount - metricCount

	if numberNonHorizontal < 0 {
		numberNonHorizontal = glyphCount
	}

	additionalLeftSideBearings := make([]int16, numberNonHorizontal)

	for i := 0; i < len(additionalLeftSideBearings); i++ {
		if int64((metricCount*4)+(i*2)) >= int64(header.Length) {
			break
		}
		additionalLeftSideBearings[i] = data.ReadSignedShort()
	}

	return tables.NewHorizontalMetricsTable(header, horizontalMetrics, additionalLeftSideBearings)
}
