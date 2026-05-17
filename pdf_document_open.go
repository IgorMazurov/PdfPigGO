// Package pdfpig provides the main entry points for opening and reading PDF documents.
package pdfpig

import (
	"io"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/parser"
)

// Open creates a PdfDocument for reading from the provided file bytes.
func Open(fileBytes []byte, options *content.ParsingOptions) (*content.PdfDocument, error) {
	return parser.OpenMemory(fileBytes, options)
}

// OpenFile opens a file and creates a PdfDocument for reading from the provided file path.
func OpenFile(filePath string, options *content.ParsingOptions) (*content.PdfDocument, error) {
	return parser.OpenFile(filePath, options)
}

// OpenStream opens a PDF document from the given seekable stream.
// If the stream does not support seeking, it will be copied into memory first.
// The caller must manage disposing the stream; the created PdfDocument will not dispose it.
func OpenStream(stream io.ReadSeeker, options *content.ParsingOptions) (*content.PdfDocument, error) {
	return parser.OpenStream(stream, options)
}
