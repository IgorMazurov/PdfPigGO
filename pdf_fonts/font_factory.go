// Package pdffonts provides factory interfaces for creating PDF fonts.
package pdffonts

import (
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// FontFactory creates a font instance from a PDF dictionary token.
type FontFactory interface {
	// Get returns the font described by the given dictionary, or an error if
	// the dictionary does not describe a recognised font type.
	Get(dictionary *tokens.DictionaryToken) (fonts.Font, error)
}
