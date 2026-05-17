package page_segmenter

import "github.com/uglytoad/pdfpig/go/document_layout_analysis"

// PageSegmenterOptions holds configuration for page segmentation operations.
type PageSegmenterOptions interface {
	document_layout_analysis.DlaOptions

	// WordSeparator returns the separator used between words when building lines.
	WordSeparator() string

	// SetWordSeparator sets the separator used between words when building lines.
	SetWordSeparator(value string)

	// LineSeparator returns the separator used between lines when building paragraphs.
	LineSeparator() string

	// SetLineSeparator sets the separator used between lines when building paragraphs.
	SetLineSeparator(value string)
}
