package tokenization

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// PdfTokenScanner tokenizes objects from bytes in a PDF file.
type PdfTokenScanner interface {
	SeekableTokenScanner

	// Close releases any resources held by the scanner.
	Close() error

	// Get returns the tokenized object with the given indirect reference.
	// Returns nil if the reference is undefined.
	Get(reference core.IndirectReference) *tokens.ObjectToken

	// ReplaceToken adds a token to an internal cache that will be returned
	// instead of scanning the source PDF data.
	ReplaceToken(reference core.IndirectReference, token tokens.Token)
}
