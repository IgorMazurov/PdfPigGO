package truetypeparser

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
)

// HorizontalHeaderTableParser parses the horizontal header table of a TrueType font.
type HorizontalHeaderTableParser struct{}

// Parse reads and interprets the horizontal header (hhea) table from raw TrueType data.
func (p *HorizontalHeaderTableParser) Parse(header truetype.TrueTypeHeaderTable, data *TrueTypeDataBytes, register *TableRegisterBuilder) (tables.HorizontalHeaderTable, error) {
	if _, err := data.Seek(int64(header.Offset), 0); err != nil {
		return tables.HorizontalHeaderTable{}, err
	}

	majorVersion := int(data.ReadUnsignedShort())
	minorVersion := int(data.ReadUnsignedShort())

	ascent := data.ReadSignedShort()
	descent := data.ReadSignedShort()
	lineGap := data.ReadSignedShort()

	advanceWidthMaximum := data.ReadUnsignedShort()

	minLeftSideBearing := data.ReadSignedShort()
	minRightSideBearing := data.ReadSignedShort()
	xMaxExtent := data.ReadSignedShort()

	caretSlopeRise := data.ReadSignedShort()
	caretSlopeRun := data.ReadSignedShort()
	caretOffset := data.ReadSignedShort()

	_ = data.ReadSignedShort()
	_ = data.ReadSignedShort()
	_ = data.ReadSignedShort()
	_ = data.ReadSignedShort()

	metricDataFormat := data.ReadSignedShort()

	if metricDataFormat != 0 {
		return tables.HorizontalHeaderTable{}, errors.New("the metric data format for a horizontal header table should be 0")
	}

	numberOfHMetrics := data.ReadUnsignedShort()

	return tables.NewHorizontalHeaderTable(
		header,
		majorVersion, minorVersion,
		ascent, descent, lineGap,
		advanceWidthMaximum,
		minLeftSideBearing, minRightSideBearing,
		xMaxExtent, caretSlopeRise, caretSlopeRun, caretOffset, metricDataFormat,
		numberOfHMetrics,
	), nil
}
