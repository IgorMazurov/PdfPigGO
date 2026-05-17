// Package parser provides types for parsing PDF font structures.
package parser

import (
	"github.com/uglytoad/pdfpig/go/fonts"
	enc "github.com/uglytoad/pdfpig/go/fonts/encodings"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// EncodingReader reads a font encoding from a PDF font dictionary.
type EncodingReader interface {
	// Read extracts the encoding from the given font dictionary, using the
	// optional font descriptor and fallback encoding when available.
	Read(
		fontDictionary *tokens.DictionaryToken,
		descriptor *fonts.FontDescriptor,
		fontEncoding *enc.Encoding,
	) *enc.Encoding
}
