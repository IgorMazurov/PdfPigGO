// Package cidfonts provides types for handling CID fonts within PDF documents.
package cidfonts

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/geometry"
	pdffonts "github.com/uglytoad/pdfpig/go/pdf_fonts"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Type2CidFont represents a CID font containing glyph descriptions based on
// the TrueType font format.
type Type2CidFont struct {
	fontProgram            CidFontProgram
	verticalWritingMetrics *VerticalWritingMetrics
	widths                 map[int]float64
	defaultWidth           *float64
	cidToGid               *CharacterIdentifierToGlyphIndexMap

	typeToken  *tokens.NameToken
	subType    *tokens.NameToken
	baseFont   *tokens.NameToken
	systemInfo pdffonts.CharacterIdentifierSystemInfo
	fontMatrix core.TransformationMatrix
	descriptor *fonts.FontDescriptor
}

// NewType2CidFont creates a new Type 2 CID font.
func NewType2CidFont(
	typeToken *tokens.NameToken,
	subType *tokens.NameToken,
	baseFont *tokens.NameToken,
	systemInfo pdffonts.CharacterIdentifierSystemInfo,
	descriptor *fonts.FontDescriptor,
	fontProgram CidFontProgram,
	verticalWritingMetrics *VerticalWritingMetrics,
	widths map[int]float64,
	defaultWidth *float64,
	cidToGid *CharacterIdentifierToGlyphIndexMap,
) *Type2CidFont {
	multiplier := 1000.0
	if fontProgram != nil {
		multiplier = float64(fontProgram.GetFontMatrixMultiplier())
	}
	scale := 1.0 / multiplier

	return &Type2CidFont{
		fontProgram:            fontProgram,
		verticalWritingMetrics: verticalWritingMetrics,
		widths:                 widths,
		defaultWidth:           defaultWidth,
		cidToGid:               cidToGid,
		typeToken:              typeToken,
		subType:                subType,
		baseFont:               baseFont,
		systemInfo:             systemInfo,
		fontMatrix:             core.FromValues(scale, 0, 0, scale, 0, 0),
		descriptor:             descriptor,
	}
}

// Type returns the font type name token.
func (f *Type2CidFont) Type() *tokens.NameToken {
	return f.typeToken
}

// SubType returns the subtype name token.
func (f *Type2CidFont) SubType() *tokens.NameToken {
	return f.subType
}

// BaseFont returns the PostScript name of the CIDFont.
func (f *Type2CidFont) BaseFont() *tokens.NameToken {
	return f.baseFont
}

// SystemInfo returns the character collection definition for the font.
func (f *Type2CidFont) SystemInfo() pdffonts.CidFontSystemInfo {
	return f.systemInfo
}

// FontMatrix returns the transformation matrix for the font.
func (f *Type2CidFont) FontMatrix() core.TransformationMatrix {
	return f.fontMatrix
}

// CidFontType returns Type2 indicating this is a TrueType-based CID font.
func (f *Type2CidFont) CidFontType() CidFontType {
	return Type2
}

// Descriptor returns the font descriptor.
func (f *Type2CidFont) Descriptor() *fonts.FontDescriptor {
	return f.descriptor
}

// Details returns the font details for this CID font.
func (f *Type2CidFont) Details() fonts.FontDetails {
	if f.fontProgram != nil {
		details := f.fontProgram.Details()
		if details != nil {
			return *details
		}
	}
	if f.descriptor != nil {
		name := ""
		if f.baseFont != nil {
			name = f.baseFont.Data()
		}
		detail := f.descriptor.ToDetails(name)
		return *detail
	}
	name := ""
	if f.baseFont != nil {
		name = f.baseFont.Data()
	}
	defaultDetail := fonts.GetDefault(name)
	return *defaultDetail
}

// GetWidthFromFont gets the glyph width from the underlying font program for the given character identifier.
func (f *Type2CidFont) GetWidthFromFont(characterIdentifier int) float64 {
	if f.fontProgram == nil {
		return f.GetWidthFromDictionary(characterIdentifier)
	}

	glyphMapper := func(code int) (int, bool) {
		idx := f.cidToGid.GetGlyphIndex(code)
		return idx, true
	}

	if width, ok := f.fontProgram.TryGetBoundingAdvancedWidthWithMapper(characterIdentifier, glyphMapper); ok {
		return width
	}

	return f.GetWidthFromDictionary(characterIdentifier)
}

// GetWidthFromDictionary gets the glyph width from the font dictionary for the given CID.
func (f *Type2CidFont) GetWidthFromDictionary(characterIdentifier int) float64 {
	if width, ok := f.widths[characterIdentifier]; ok {
		return width
	}

	if f.defaultWidth != nil {
		return *f.defaultWidth
	}

	if f.descriptor == nil {
		return 1000
	}

	return f.descriptor.MissingWidth
}

// GetBoundingBox returns the bounding box for the given character identifier.
func (f *Type2CidFont) GetBoundingBox(characterIdentifier int) (core.PdfRectangle, error) {
	if f.fontProgram == nil {
		return f.descriptor.BoundingBox, nil
	}

	var bbox *core.PdfRectangle
	var ok bool
	var recovered interface{}
	fn := func() {
		glyphMapper := func(code int) (int, bool) {
			idx := f.cidToGid.GetGlyphIndex(code)
			return idx, true
		}
		bbox, ok = f.fontProgram.TryGetBoundingBoxWithMapper(characterIdentifier, glyphMapper)
	}
	defer func() {
		if r := recover(); r != nil {
			recovered = r
		}
	}()
	fn()
	if recovered != nil {
		return f.descriptor.BoundingBox, nil
	}

	if ok {
		return *bbox, nil
	}

	return f.descriptor.BoundingBox, nil
}

// GetPositionVector returns the position vector for the given character identifier.
func (f *Type2CidFont) GetPositionVector(characterIdentifier int) geometry.PdfVector {
	width := f.GetWidthFromFont(characterIdentifier)
	return f.verticalWritingMetrics.GetPositionVector(characterIdentifier, width)
}

// GetDisplacementVector returns the displacement vector for the given character identifier.
func (f *Type2CidFont) GetDisplacementVector(characterIdentifier int) geometry.PdfVector {
	return f.verticalWritingMetrics.GetDisplacementVector(characterIdentifier)
}

// GetFontMatrix returns the font matrix for the given character identifier.
func (f *Type2CidFont) GetFontMatrix(characterIdentifier int) core.TransformationMatrix {
	return f.fontMatrix
}

// GetDescent returns the font descent value.
func (f *Type2CidFont) GetDescent() float64 {
	if f.fontProgram == nil {
		return f.descriptor.Descent
	}

	if descent := f.fontProgram.GetDescent(); descent != nil {
		return *descent
	}

	return f.descriptor.Descent
}

// GetAscent returns the font ascent value.
func (f *Type2CidFont) GetAscent() float64 {
	if f.fontProgram == nil {
		return f.descriptor.Ascent
	}

	if ascent := f.fontProgram.GetAscent(); ascent != nil {
		return *ascent
	}

	return f.descriptor.Ascent
}

// TryGetPath attempts to get the glyph path for the given character code.
func (f *Type2CidFont) TryGetPath(characterCode int) ([]core.PdfSubpath, bool) {
	if f.fontProgram == nil {
		return nil, false
	}

	path, ok := f.fontProgram.TryGetPathWithMapper(characterCode, func(code int) (int, bool) {
		idx := f.cidToGid.GetGlyphIndex(code)
		return idx, true
	})
	if !ok {
		return nil, false
	}

	return pointersToValues(path), true
}

// TryGetPathWithMapper attempts to get the glyph path using a custom character-code-to-glyph-ID mapper function.
func (f *Type2CidFont) TryGetPathWithMapper(characterCode int, mapper func(int) *int) ([]core.PdfSubpath, bool) {
	if f.fontProgram == nil {
		return nil, false
	}

	path, ok := f.fontProgram.TryGetPathWithMapper(characterCode, func(code int) (int, bool) {
		result := mapper(code)
		if result == nil {
			return 0, false
		}
		return *result, true
	})
	if !ok {
		return nil, false
	}

	return pointersToValues(path), true
}

// TryGetNormalisedPath attempts to get the normalised glyph path for the given character code.
func (f *Type2CidFont) TryGetNormalisedPath(characterCode int) ([]core.PdfSubpath, bool) {
	if f.fontProgram == nil {
		return nil, false
	}

	path, ok := f.fontProgram.TryGetPathWithMapper(characterCode, func(code int) (int, bool) {
		idx := f.cidToGid.GetGlyphIndex(code)
		return idx, true
	})
	if !ok {
		return nil, false
	}

	matrix := f.GetFontMatrix(characterCode)
	transformed, err := matrix.TransformPath(path)
	if err != nil {
		return nil, false
	}

	return pointersToValues(transformed), true
}

// TryGetNormalisedPathWithMapper attempts to get the normalised glyph path using a custom
// character-code-to-glyph-ID mapper function.
func (f *Type2CidFont) TryGetNormalisedPathWithMapper(characterCode int, mapper func(int) *int) ([]core.PdfSubpath, bool) {
	if f.fontProgram == nil {
		return nil, false
	}

	path, ok := f.fontProgram.TryGetPathWithMapper(characterCode, func(code int) (int, bool) {
		result := mapper(code)
		if result == nil {
			return 0, false
		}
		return *result, true
	})
	if !ok {
		return nil, false
	}

	matrix := f.GetFontMatrix(characterCode)
	transformed, err := matrix.TransformPath(path)
	if err != nil {
		return nil, false
	}

	return pointersToValues(transformed), true
}

var _ CidFont = (*Type2CidFont)(nil)
