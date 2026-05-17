package cff

import "github.com/uglytoad/pdfpig/go/fonts"

// CompactFontFormatPrivateDictionary holds the private dictionary for a Compact
// Font Format (CFF) font, extending AdobeStylePrivateDictionary with CFF-specific
// properties.
type CompactFontFormatPrivateDictionary struct {
	fonts.AdobeStylePrivateDictionary

	// InitialRandomSeed is the compatibility entry used as the initial seed value
	// for the random number generator in Type 2 charstrings.
	InitialRandomSeed float64

	// LocalSubroutineOffset is the offset in bytes for the local subroutine index
	// in this font, relative to this private dictionary. Nil if not present.
	LocalSubroutineOffset *int

	// DefaultWidthX is the default glyph width in X direction. If a glyph's width
	// equals DefaultWidthX it can be omitted from the charstring.
	DefaultWidthX float64

	// NominalWidthX is used to compute glyph widths: if not equal to DefaultWidthX,
	// the glyph width is computed by adding the charstring width to the nominal
	// width X value.
	NominalWidthX float64
}

// CompactFontFormatPrivateDictionaryBuilder is a mutable builder for constructing
// a CompactFontFormatPrivateDictionary. It embeds AdobeStylePrivateDictionaryBuilder
// for all inherited properties and adds CFF-specific fields.
type CompactFontFormatPrivateDictionaryBuilder struct {
	fonts.AdobeStylePrivateDictionaryBuilder

	InitialRandomSeed     float64
	LocalSubroutineOffset *int
	DefaultWidthX         float64
	NominalWidthX         float64
}

// NewCompactFontFormatPrivateDictionary creates a new CompactFontFormatPrivateDictionary
// from the given builder.
func NewCompactFontFormatPrivateDictionary(builder *CompactFontFormatPrivateDictionaryBuilder) CompactFontFormatPrivateDictionary {
	if builder == nil {
		builder = &CompactFontFormatPrivateDictionaryBuilder{}
	}

	return CompactFontFormatPrivateDictionary{
		AdobeStylePrivateDictionary: fonts.NewAdobeStylePrivateDictionary(&builder.AdobeStylePrivateDictionaryBuilder),
		InitialRandomSeed:           builder.InitialRandomSeed,
		LocalSubroutineOffset:       builder.LocalSubroutineOffset,
		DefaultWidthX:               builder.DefaultWidthX,
		NominalWidthX:               builder.NominalWidthX,
	}
}

// DefaultCompactFontFormatPrivateDictionary returns a CompactFontFormatPrivateDictionary
// with all default values.
var DefaultCompactFontFormatPrivateDictionary = NewCompactFontFormatPrivateDictionary(&CompactFontFormatPrivateDictionaryBuilder{})
