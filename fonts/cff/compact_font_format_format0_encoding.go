// Package cff provides types for parsing and working with Compact Font Format (CFF) data in PDF files.

package cff

import "github.com/uglytoad/pdfpig/go/fonts/encodings"

// CFFFormat0CodeEntry holds a character code, SID, and resolved glyph name
// used when building a format 0 encoding table.
type CFFFormat0CodeEntry struct {
	Code int
	Sid  int
	Str  string
}

// CompactFontFormatFormat0Encoding represents a CFF format 0 encoding table.
// Format 0 maps each character code sequentially from glyph ID 1 to numberOfCodes,
// optionally followed by supplement entries that override specific codes.
type CompactFontFormatFormat0Encoding struct {
	*CFFBuiltInEncoding
}

// NewCompactFontFormatFormat0Encoding creates a new format 0 encoding from the given
// code entries and supplements. It adds .notdef at code 0 first, then registers each
// entry via AddWithName, matching the C# constructor semantics exactly.
func NewCompactFontFormatFormat0Encoding(entries []CFFFormat0CodeEntry, supplements []CFFBuiltInEncodingSupplement) *CompactFontFormatFormat0Encoding {
	if supplements == nil {
		supplements = make([]CFFBuiltInEncodingSupplement, 0)
	}

	base := NewCompactFontFormatBaseEncoding()
	base.AddWithName(0, 0, encodings.NotDefined)

	for _, e := range entries {
		base.AddWithName(e.Code, e.Sid, e.Str)
	}

	builtIn := &CFFBuiltInEncoding{
		CompactFontFormatBaseEncoding: base,
		Supplements:                   supplements,
	}

	return &CompactFontFormatFormat0Encoding{
		CFFBuiltInEncoding: builtIn,
	}
}
