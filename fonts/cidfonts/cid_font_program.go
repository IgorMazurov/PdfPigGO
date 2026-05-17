package cidfonts

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
)

// CidFontProgram represents either an Adobe Type 1 or TrueType font program for a CIDFont.
type CidFontProgram interface {
	// Details returns the font details associated with this font program.
	Details() *fonts.FontDetails

	// TryGetBoundingBox attempts to get the bounding box for the given character identifier.
	TryGetBoundingBox(characterIdentifier int) (*core.PdfRectangle, bool)

	// TryGetBoundingBoxWithMapper attempts to get the bounding box using a character code to glyph ID mapper.
	TryGetBoundingBoxWithMapper(characterIdentifier int, characterCodeToGlyphId func(int) (int, bool)) (*core.PdfRectangle, bool)

	// TryGetBoundingAdvancedWidthWithMapper attempts to get the bounding advanced width using a character code to glyph ID mapper.
	TryGetBoundingAdvancedWidthWithMapper(characterIdentifier int, characterCodeToGlyphId func(int) (int, bool)) (float64, bool)

	// TryGetBoundingAdvancedWidth attempts to get the bounding advanced width for the given character identifier.
	TryGetBoundingAdvancedWidth(characterIdentifier int) (float64, bool)

	// GetDescent returns the font descent, or nil if not available.
	GetDescent() *float64

	// GetAscent returns the font ascent, or nil if not available.
	GetAscent() *float64

	// TryGetPath attempts to get the drawing path for the given character code.
	TryGetPath(characterCode int) ([]*core.PdfSubpath, bool)

	// TryGetPathWithMapper attempts to get the drawing path using a character code to glyph ID mapper.
	TryGetPathWithMapper(characterCode int, characterCodeToGlyphId func(int) (int, bool)) ([]*core.PdfSubpath, bool)

	// GetFontMatrixMultiplier returns the font matrix multiplier value.
	GetFontMatrixMultiplier() int

	// TryGetFontMatrix attempts to get the font matrix for the given character code.
	TryGetFontMatrix(characterCode int) (*core.TransformationMatrix, bool)
}
