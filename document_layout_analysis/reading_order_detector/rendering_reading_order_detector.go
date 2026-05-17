package reading_order_detector

import (
	"sort"

	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
)

// RenderingReadingOrderDetector determines reading order by sorting text blocks
// according to the average TextSequence of their constituent letters.
type RenderingReadingOrderDetector struct{}

var _ ReadingOrderDetector = (*RenderingReadingOrderDetector)(nil)

// Instance is the singleton instance of RenderingReadingOrderDetector.
var Instance = &RenderingReadingOrderDetector{}

// Get returns the text blocks sorted by their average letter TextSequence,
// setting each block's ReadingOrder field accordingly.
func (r *RenderingReadingOrderDetector) Get(textBlocks []*document_layout_analysis.TextBlock) []*document_layout_analysis.TextBlock {
	sorted := make([]*document_layout_analysis.TextBlock, len(textBlocks))
	copy(sorted, textBlocks)

	sort.SliceStable(sorted, func(i, j int) bool {
		return avgTextSequence(sorted[i]) < avgTextSequence(sorted[j])
	})

	for order, block := range sorted {
		block.SetReadingOrder(order) //nolint:errcheck
	}

	return sorted
}

// avgTextSequence returns the average TextSequence across all letters in all
// words of all text lines within the given block. Returns 0 for empty blocks.
func avgTextSequence(block *document_layout_analysis.TextBlock) float64 {
	var sum float64
	count := 0

	for _, line := range block.TextLines {
		for _, word := range line.Words {
			for _, letter := range word.Letters {
				sum += float64(letter.TextSequence)
				count++
			}
		}
	}

	if count == 0 {
		return 0
	}

	return sum / float64(count)
}
