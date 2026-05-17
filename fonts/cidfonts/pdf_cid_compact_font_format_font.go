package cidfonts

import (
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/cff"
)

// PdfCidCompactFontFormatFont wraps a CFF font collection and implements CidFontProgram.
type PdfCidCompactFontFormatFont struct {
	fontCollection *cff.CompactFontFormatFontCollection
	details        *fonts.FontDetails
}

// NewPdfCidCompactFontFormatFont creates a new CFF-based CID font program from the given font collection.
func NewPdfCidCompactFontFormatFont(fontCollection *cff.CompactFontFormatFontCollection) *PdfCidCompactFontFormatFont {
	first := fontCollection.FirstFont()
	return &PdfCidCompactFontFormatFont{
		fontCollection: fontCollection,
		details:        getDetails(first),
	}
}

func getDetails(font *cff.CompactFontFormatFont) *fonts.FontDetails {
	if font == nil {
		return fonts.GetDefault("")
	}

	withWeightValues := func(isBold bool, weight int) *fonts.FontDetails {
		return fonts.NewFontDetails("", isBold, weight, font.ItalicAngle() != 0)
	}

	switch strings.ToLower(font.Weight()) {
	case "light":
		return withWeightValues(false, 300)
	case "semibold":
		return withWeightValues(true, 600)
	case "bold":
		return withWeightValues(true, fonts.BoldWeight)
	case "black":
		return withWeightValues(true, 900)
	default:
		return withWeightValues(false, fonts.DefaultWeight)
	}
}

// Details returns the font details for this CFF font.
func (f *PdfCidCompactFontFormatFont) Details() *fonts.FontDetails {
	return f.details
}

// GetFontTransformationMatrix returns the transformation matrix of the first font in the collection.
func (f *PdfCidCompactFontFormatFont) GetFontTransformationMatrix() core.TransformationMatrix {
	return f.fontCollection.GetFirstTransformationMatrix()
}

// GetCharacterBoundingBox returns the bounding box for a character by name, or nil if not found.
func (f *PdfCidCompactFontFormatFont) GetCharacterBoundingBox(characterName string) *core.PdfRectangle {
	return f.fontCollection.GetCharacterBoundingBox(characterName)
}

// GetDescent returns nil because ascent/descent are not currently supported for CFF fonts.
func (f *PdfCidCompactFontFormatFont) GetDescent() *float64 {
	return nil
}

// GetAscent returns nil because ascent/descent are not currently supported for CFF fonts.
func (f *PdfCidCompactFontFormatFont) GetAscent() *float64 {
	return nil
}

// TryGetBoundingBox attempts to get the bounding box for the given character identifier.
func (f *PdfCidCompactFontFormatFont) TryGetBoundingBox(characterIdentifier int) (*core.PdfRectangle, bool) {
	font := f.getFont()
	characterName := font.GetCharacterName(characterIdentifier, false)

	if strings.EqualFold(characterName, fonts.NotDefined) {
		return nil, false
	}

	bbox := font.GetCharacterBoundingBox(characterName)
	if bbox == nil {
		rect := core.NewPdfRectangleFromInt(0, 0, 500, 0)
		return &rect, true
	}
	return bbox, true
}

// TryGetBoundingBoxWithMapper attempts to get the bounding box using a character code to glyph ID mapper.
func (f *PdfCidCompactFontFormatFont) TryGetBoundingBoxWithMapper(characterIdentifier int, characterCodeToGlyphId func(int) (int, bool)) (*core.PdfRectangle, bool) {
	font := f.getFont()

	glyphId, hasGlyphId := characterCodeToGlyphId(characterIdentifier)

	var name string
	if hasGlyphId {
		name = font.GetCharacterName(glyphId, false)
	} else {
		name = font.GetCharacterName(characterIdentifier, false)
	}

	bbox := font.GetCharacterBoundingBox(name)
	if bbox != nil {
		return bbox, true
	}

	return nil, false
}

// TryGetBoundingAdvancedWidthWithMapper delegates to TryGetBoundingAdvancedWidth.
func (f *PdfCidCompactFontFormatFont) TryGetBoundingAdvancedWidthWithMapper(characterIdentifier int, characterCodeToGlyphId func(int) (int, bool)) (float64, bool) {
	return f.TryGetBoundingAdvancedWidth(characterIdentifier)
}

// TryGetBoundingAdvancedWidth always returns false because CFF fonts do not currently
// provide bounding advanced width information.
func (f *PdfCidCompactFontFormatFont) TryGetBoundingAdvancedWidth(characterIdentifier int) (float64, bool) {
	return 0, false
}

// GetFontMatrixMultiplier returns the standard CFF font matrix multiplier of 1000.
func (f *PdfCidCompactFontFormatFont) GetFontMatrixMultiplier() int {
	return 1000
}

// TryGetFontMatrix attempts to get the font matrix for the given character code.
func (f *PdfCidCompactFontFormatFont) TryGetFontMatrix(characterCode int) (*core.TransformationMatrix, bool) {
	font := f.getFont()
	name := font.GetCharacterName(characterCode, true)
	if name == "" {
		return nil, false
	}
	matrix := font.GetFontMatrix(name)
	if matrix != nil {
		return matrix, true
	}
	return nil, false
}

// getFont returns the first font from the collection.
func (f *PdfCidCompactFontFormatFont) getFont() *cff.CompactFontFormatFont {
	return f.fontCollection.FirstFont()
}

// TryGetPath attempts to get the drawing path for the given character code.
func (f *PdfCidCompactFontFormatFont) TryGetPath(characterCode int) ([]*core.PdfSubpath, bool) {
	font := f.getFont()
	characterName := font.GetCharacterName(characterCode, false)

	if strings.EqualFold(characterName, fonts.NotDefined) {
		return nil, false
	}

	path, ok := font.TryGetPath(characterName)
	if ok {
		return path, true
	}

	return nil, false
}

// TryGetPathWithMapper attempts to get the drawing path using a character code to glyph ID mapper.
func (f *PdfCidCompactFontFormatFont) TryGetPathWithMapper(characterCode int, characterCodeToGlyphId func(int) (int, bool)) ([]*core.PdfSubpath, bool) {
	glyphId, hasGlyphId := characterCodeToGlyphId(characterCode)

	var characterName string
	if hasGlyphId {
		characterName = f.getFont().GetCharacterName(glyphId, false)
	} else {
		characterName = f.getFont().GetCharacterName(characterCode, false)
	}

	if strings.EqualFold(characterName, fonts.NotDefined) {
		return nil, false
	}

	path, ok := f.getFont().TryGetPath(characterName)
	if ok {
		return path, true
	}

	return nil, false
}

var _ CidFontProgram = (*PdfCidCompactFontFormatFont)(nil)
