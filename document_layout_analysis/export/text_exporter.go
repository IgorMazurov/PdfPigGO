package export

import "github.com/uglytoad/pdfpig/go/content"

// TextExporter exports the page's text into a desired format.
type TextExporter interface {
	// Get returns the text representation of a page in a desired format.
	Get(page *content.Page) string
}
