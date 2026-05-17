// Package handlers provides interfaces for PDF font parsing handlers.
package handlers

import (
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// FontHandler generates a font from a PDF dictionary token.
type FontHandler interface {
	// Generate creates a font instance from the given dictionary.
	Generate(dictionary *tokens.DictionaryToken) fonts.Font
}
