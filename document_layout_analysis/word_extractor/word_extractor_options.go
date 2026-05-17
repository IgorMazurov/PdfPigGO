package word_extractor

import "github.com/uglytoad/pdfpig/go/document_layout_analysis"

// WordExtractorOptions holds configuration for word extraction operations.
type WordExtractorOptions interface {
	document_layout_analysis.DlaOptions
}
