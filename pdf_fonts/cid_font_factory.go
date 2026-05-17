// Package pdffonts provides factory interfaces for creating PDF fonts.
package pdffonts

import (
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CidFontFactory creates CID fonts from their PDF dictionary representations.
// The actual implementation is in fonts/parser/cid_font_factory.go to avoid import cycles.
type CidFontFactory struct {
	log    logging.Log
	scanner any
}

// NewCidFontFactory creates a new CidFontFactory.
func NewCidFontFactory(log logging.Log, scanner any, filterProvider any) *CidFontFactory {
	return &CidFontFactory{log: log, scanner: scanner}
}

// CreateCIDFont creates a CID font from the given dictionary.
func (f *CidFontFactory) CreateCIDFont(dictionary *tokens.DictionaryToken) fonts.Font {
	return nil
}
