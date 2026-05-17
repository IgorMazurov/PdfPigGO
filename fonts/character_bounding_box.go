// Package fonts provides types for PDF font handling.
package fonts

import (
	"github.com/uglytoad/pdfpig/go/core"
)

// CharacterBoundingBox represents the bounding box of a character glyph.
type CharacterBoundingBox struct {
	// GlyphBounds is the rectangular bounds of the glyph.
	GlyphBounds core.PdfRectangle
	// Width is the width of the character.
	Width float64
}

// NewCharacterBoundingBox creates a new CharacterBoundingBox with the given
// glyph bounds and width.
func NewCharacterBoundingBox(bounds core.PdfRectangle, width float64) *CharacterBoundingBox {
	return &CharacterBoundingBox{
		GlyphBounds: bounds,
		Width:       width,
	}
}
