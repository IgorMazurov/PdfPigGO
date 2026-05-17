package cff

import "github.com/uglytoad/pdfpig/go/fonts/cfftypes"

// CompactFontFormatIndex is an alias for the shared type.
type CompactFontFormatIndex = cfftypes.CompactFontFormatIndex

// None is an empty CompactFontFormatIndex instance.
var None = cfftypes.NewCompactFontFormatIndex(nil)

// NewCompactFontFormatIndex creates a new CompactFontFormatIndex from the given bytes.
var NewCompactFontFormatIndex = cfftypes.NewCompactFontFormatIndex
