// Package cff provides types for parsing and working with Compact Font Format (CFF) data in PDF files.
package cff

import (
	"github.com/uglytoad/pdfpig/go/fonts/encodings"
)

// CompactFontFormatBaseEncoding is the base encoding type for CFF fonts,
// mapping character codes to glyph names from a Compact Font Format encoding.
type CompactFontFormatBaseEncoding struct {
	encodings.Encoding
	codeToName map[int]string
}

// NewCompactFontFormatBaseEncoding creates a new CompactFontFormatBaseEncoding.
func NewCompactFontFormatBaseEncoding() *CompactFontFormatBaseEncoding {
	return &CompactFontFormatBaseEncoding{
		Encoding:   *encodings.NewEncoding(),
		codeToName: make(map[int]string, 250),
	}
}

// EncodingName returns the name of this encoding.
func (e *CompactFontFormatBaseEncoding) EncodingName() string {
	return "CFF"
}

// GetName returns the PostScript name of the glyph for the given character code,
// or NotDefined if not found.
func (e *CompactFontFormatBaseEncoding) GetName(code int) string {
	name, ok := e.codeToName[code]
	if !ok {
		return encodings.NotDefined
	}
	return name
}

// AddWithName adds a character code, SID, and explicit name to the encoding.
func (e *CompactFontFormatBaseEncoding) AddWithName(code int, _sid int, name string) {
	e.codeToName[code] = name
	e.Encoding.Add(code, name)
}

// AddWithSid adds a character code resolved from an SID to the encoding.
func (e *CompactFontFormatBaseEncoding) AddWithSid(code int, sid int) {
	name := GetName(sid)
	e.codeToName[code] = name
	e.Encoding.Add(code, name)
}

// ToEncoding returns the underlying encodings.Encoding pointer.
// This bridges the internal CFF encodingSource interface to the
// common *encodings.Encoding type used by font handlers.
func (e *CompactFontFormatBaseEncoding) ToEncoding() *encodings.Encoding {
	return &e.Encoding
}
