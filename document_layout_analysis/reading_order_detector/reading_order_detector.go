package reading_order_detector

import (
	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
)

// ReadingOrderDetector determines the reading order of text blocks on a page.
// Implementations must call SetReadingOrder on each TextBlock to assign its
// position in the reading sequence.
type ReadingOrderDetector interface {
	// Get returns the text blocks sorted by their reading order.
	// The implementation is expected to mutate each block's ReadingOrder field
	// via SetReadingOrder before returning the ordered slice.
	Get(textBlocks []*document_layout_analysis.TextBlock) []*document_layout_analysis.TextBlock
}
