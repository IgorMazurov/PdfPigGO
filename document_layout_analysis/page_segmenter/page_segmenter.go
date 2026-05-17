package page_segmenter

import (
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
)

// PageSegmenter divides a page into areas, each consisting of a layout structure
// (blocks, lines, etc.). See "Performance Comparison of Six Algorithms for Page
// Segmentation" by Faisal Shafait, Daniel Keysers, and Thomas M. Breuel.
type PageSegmenter interface {
	// GetBlocks returns text blocks generated from the given words on a page.
	GetBlocks(words []*content.Word) []*document_layout_analysis.TextBlock
}
