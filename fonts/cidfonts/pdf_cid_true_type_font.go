// Package cidfonts provides types for handling CID fonts within PDF documents.
package cidfonts

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
)

// PdfCidTrueTypeFont wraps a TrueType font and implements CidFontProgram.
type PdfCidTrueTypeFont struct {
	font    *truetypeparser.TrueTypeFont
	details *fonts.FontDetails
}

// NewPdfCidTrueTypeFont creates a new TrueType-based CID font program from the given font.
func NewPdfCidTrueTypeFont(font *truetypeparser.TrueTypeFont) *PdfCidTrueTypeFont {
	register := font.TableRegister()
	header := register.HeaderTable
	isBold := header.MacStyle()&tables.HeaderMacStyleBold != 0
	isItalic := header.MacStyle()&tables.HeaderMacStyleItalic != 0

	weight := fonts.DefaultWeight
	if isBold {
		weight = fonts.BoldWeight
	}

	return &PdfCidTrueTypeFont{
		font:    font,
		details: fonts.NewFontDetails(font.Name(), isBold, weight, isItalic),
	}
}

// Details returns the font details for this TrueType font.
func (f *PdfCidTrueTypeFont) Details() *fonts.FontDetails {
	return f.details
}

// TryGetBoundingBox attempts to get the bounding box for the given character identifier.
func (f *PdfCidTrueTypeFont) TryGetBoundingBox(characterIdentifier int) (*core.PdfRectangle, bool) {
	rect, ok := f.font.TryGetBoundingBox(characterIdentifier)
	if !ok {
		return nil, false
	}
	return &rect, true
}

// TryGetBoundingBoxWithMapper attempts to get the bounding box using a character code to glyph ID mapper.
func (f *PdfCidTrueTypeFont) TryGetBoundingBoxWithMapper(
	characterIdentifier int,
	characterCodeToGlyphId func(int) (int, bool),
) (*core.PdfRectangle, bool) {
	adapter := func(code int) *int {
		id, ok := characterCodeToGlyphId(code)
		if !ok {
			return nil
		}
		return &id
	}

	rect, ok := f.font.TryGetBoundingBoxWithMapping(characterIdentifier, adapter)
	if !ok {
		return nil, false
	}
	return &rect, true
}

// TryGetBoundingAdvancedWidth attempts to get the bounding advanced width for the given character identifier.
func (f *PdfCidTrueTypeFont) TryGetBoundingAdvancedWidth(characterIdentifier int) (float64, bool) {
	return f.font.TryGetAdvanceWidth(characterIdentifier)
}

// TryGetBoundingAdvancedWidthWithMapper attempts to get the bounding advanced width using a character code to glyph ID mapper.
func (f *PdfCidTrueTypeFont) TryGetBoundingAdvancedWidthWithMapper(
	characterIdentifier int,
	characterCodeToGlyphId func(int) (int, bool),
) (float64, bool) {
	adapter := func(code int) *int {
		id, ok := characterCodeToGlyphId(code)
		if !ok {
			return nil
		}
		return &id
	}

	return f.font.TryGetAdvanceWidthWithMapping(characterIdentifier, adapter)
}

// GetFontMatrixMultiplier returns the units per em value from the font.
func (f *PdfCidTrueTypeFont) GetFontMatrixMultiplier() int {
	return int(f.font.GetUnitsPerEm())
}

// GetDescent returns the font descent as a float64 pointer, or nil if unavailable.
func (f *PdfCidTrueTypeFont) GetDescent() *float64 {
	v := float64(f.font.TableRegister().HorizontalHeaderTable.Descent())
	return &v
}

// GetAscent returns the font ascent as a float64 pointer, or nil if unavailable.
func (f *PdfCidTrueTypeFont) GetAscent() *float64 {
	v := float64(f.font.TableRegister().HorizontalHeaderTable.Ascent())
	return &v
}

// TryGetFontMatrix always returns false because TrueType fonts do not provide per-character matrices.
func (f *PdfCidTrueTypeFont) TryGetFontMatrix(characterCode int) (*core.TransformationMatrix, bool) {
	return nil, false
}

// TryGetPath attempts to get the drawing path for the given character code.
func (f *PdfCidTrueTypeFont) TryGetPath(characterCode int) ([]*core.PdfSubpath, bool) {
	return f.font.TryGetPath(characterCode)
}

// TryGetPathWithMapper attempts to get the drawing path using a character code to glyph ID mapper.
func (f *PdfCidTrueTypeFont) TryGetPathWithMapper(
	characterCode int,
	characterCodeToGlyphId func(int) (int, bool),
) ([]*core.PdfSubpath, bool) {
	adapter := func(code int) *int {
		id, ok := characterCodeToGlyphId(code)
		if !ok {
			return nil
		}
		return &id
	}

	return f.font.TryGetPathWithMapping(characterCode, adapter)
}

var _ CidFontProgram = (*PdfCidTrueTypeFont)(nil)
