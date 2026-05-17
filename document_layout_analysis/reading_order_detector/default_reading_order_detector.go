package reading_order_detector

import (
	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
)

// DefaultReadingOrderDetector returns text blocks in their original order
// without applying any sorting logic.
type DefaultReadingOrderDetector struct{}

var _ ReadingOrderDetector = (*DefaultReadingOrderDetector)(nil)

// Default is the singleton instance of DefaultReadingOrderDetector.
var Default = &DefaultReadingOrderDetector{}

// Get returns the text blocks unchanged, performing no reordering.
func (d *DefaultReadingOrderDetector) Get(textBlocks []*document_layout_analysis.TextBlock) []*document_layout_analysis.TextBlock {
	return textBlocks
}
