// Package cidfonts provides types for handling CID fonts within PDF documents.
package cidfonts

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/geometry"
	pdffonts "github.com/uglytoad/pdfpig/go/pdf_fonts"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Type0CidFont represents a CID font containing glyph descriptions based on
// the Adobe Type 1 font format.
type Type0CidFont struct {
	fontProgram            CidFontProgram
	verticalWritingMetrics *VerticalWritingMetrics
	defaultWidth           *float64
	scale                  float64

	typeToken *tokens.NameToken
	subType   *tokens.NameToken
	baseFont  *tokens.NameToken
	systemInfo pdffonts.CharacterIdentifierSystemInfo
	fontMatrix core.TransformationMatrix
	descriptor *fonts.FontDescriptor
	widths     map[int]float64
}

// NewType0CidFont creates a new Type 0 CID font.
func NewType0CidFont(
	fontProgram CidFontProgram,
	typeToken *tokens.NameToken,
	subType *tokens.NameToken,
	baseFont *tokens.NameToken,
	systemInfo pdffonts.CharacterIdentifierSystemInfo,
	descriptor *fonts.FontDescriptor,
	verticalWritingMetrics *VerticalWritingMetrics,
	widths map[int]float64,
	defaultWidth *float64,
) *Type0CidFont {
	multiplier := 1000.0
	if fontProgram != nil {
		multiplier = float64(fontProgram.GetFontMatrixMultiplier())
	}
	scale := 1.0 / multiplier

	return &Type0CidFont{
		fontProgram:            fontProgram,
		verticalWritingMetrics: verticalWritingMetrics,
		defaultWidth:           defaultWidth,
		scale:                  scale,
		typeToken:              typeToken,
		subType:                subType,
		baseFont:               baseFont,
		systemInfo:             systemInfo,
		fontMatrix:             core.FromValues(scale, 0, 0, scale, 0, 0),
		descriptor:             descriptor,
		widths:                 widths,
	}
}

// Type returns the font type name token.
func (f *Type0CidFont) Type() *tokens.NameToken {
	return f.typeToken
}

// SubType returns either Type0 or Type2.
func (f *Type0CidFont) SubType() *tokens.NameToken {
	return f.subType
}

// BaseFont returns the PostScript name of the CIDFont.
func (f *Type0CidFont) BaseFont() *tokens.NameToken {
	return f.baseFont
}

// SystemInfo returns the character collection definition for the font.
func (f *Type0CidFont) SystemInfo() pdffonts.CidFontSystemInfo {
	return f.systemInfo
}

// Details returns the font details for this CID font.
func (f *Type0CidFont) Details() fonts.FontDetails {
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

// FontMatrix returns the transformation matrix for the font.
func (f *Type0CidFont) FontMatrix() core.TransformationMatrix {
	return f.fontMatrix
}

// CidFontType returns Type0 indicating this is an Adobe Type 1 CID font.
func (f *Type0CidFont) CidFontType() CidFontType {
	return Type0
}

// Descriptor returns the font descriptor.
func (f *Type0CidFont) Descriptor() *fonts.FontDescriptor {
	return f.descriptor
}

// GetWidthFromFont gets the glyph width from the underlying font program for the given character identifier.
func (f *Type0CidFont) GetWidthFromFont(characterCode int) float64 {
	return f.GetWidthFromDictionary(characterCode)
}

// GetWidthFromDictionary gets the glyph width from the font dictionary for the given CID.
func (f *Type0CidFont) GetWidthFromDictionary(cid int) float64 {
	if cid < 0 {
		panic(fmt.Sprintf("The provided character code was negative: %d.", cid))
	}

	if width, ok := f.widths[cid]; ok {
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
func (f *Type0CidFont) GetBoundingBox(characterIdentifier int) (core.PdfRectangle, error) {
	if characterIdentifier < 0 {
		return core.PdfRectangle{}, fmt.Errorf("the provided character identifier was negative: %d", characterIdentifier)
	}

	if f.fontProgram == nil {
		if f.descriptor != nil {
			return f.descriptor.BoundingBox, nil
		}
		return core.NewPdfRectangleFloat(0, 0, 1000, 1.0/f.scale), nil
	}

	if bbox, ok := f.fontProgram.TryGetBoundingBox(characterIdentifier); ok {
		return *bbox, nil
	}

	if width, ok := f.widths[characterIdentifier]; ok {
		return core.NewPdfRectangleFloat(0, 0, width, 1.0/f.scale), nil
	}

	if f.defaultWidth != nil {
		return core.NewPdfRectangleFloat(0, 0, *f.defaultWidth, 1.0/f.scale), nil
	}

	return core.NewPdfRectangleFloat(0, 0, 1000, 1.0/f.scale), nil
}

// GetPositionVector returns the position vector for the given character identifier.
func (f *Type0CidFont) GetPositionVector(characterIdentifier int) geometry.PdfVector {
	width := f.GetWidthFromFont(characterIdentifier)
	return f.verticalWritingMetrics.GetPositionVector(characterIdentifier, width)
}

// GetDisplacementVector returns the displacement vector for the given character identifier.
func (f *Type0CidFont) GetDisplacementVector(characterIdentifier int) geometry.PdfVector {
	return f.verticalWritingMetrics.GetDisplacementVector(characterIdentifier)
}

// GetFontMatrix returns the font matrix for the given character identifier.
func (f *Type0CidFont) GetFontMatrix(characterIdentifier int) core.TransformationMatrix {
	if f.fontProgram == nil {
		return f.fontMatrix
	}

	if m, ok := f.fontProgram.TryGetFontMatrix(characterIdentifier); ok {
		return *m
	}

	return f.fontMatrix
}

// GetDescent returns the font descent value.
func (f *Type0CidFont) GetDescent() float64 {
	if f.fontProgram == nil {
		return f.descriptor.Descent
	}

	if descent := f.fontProgram.GetDescent(); descent != nil {
		return *descent
	}

	return f.descriptor.Descent
}

// GetAscent returns the font ascent value.
func (f *Type0CidFont) GetAscent() float64 {
	if f.fontProgram == nil {
		return f.descriptor.Ascent
	}

	if ascent := f.fontProgram.GetAscent(); ascent != nil {
		return *ascent
	}

	return f.descriptor.Ascent
}

// TryGetPath attempts to get the glyph path for the given character code.
func (f *Type0CidFont) TryGetPath(characterCode int) ([]core.PdfSubpath, bool) {
	if f.fontProgram == nil {
		return nil, false
	}

	path, ok := f.fontProgram.TryGetPath(characterCode)
	if !ok {
		return nil, false
	}

	return pointersToValues(path), true
}

// TryGetPathWithMapper attempts to get the glyph path using a custom character-code-to-glyph-ID mapper function.
func (f *Type0CidFont) TryGetPathWithMapper(characterCode int, mapper func(int) *int) ([]core.PdfSubpath, bool) {
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
func (f *Type0CidFont) TryGetNormalisedPath(characterCode int) ([]core.PdfSubpath, bool) {
	if f.fontProgram == nil {
		return nil, false
	}

	path, ok := f.fontProgram.TryGetPath(characterCode)
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
func (f *Type0CidFont) TryGetNormalisedPathWithMapper(characterCode int, mapper func(int) *int) ([]core.PdfSubpath, bool) {
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

func pointersToValues(paths []*core.PdfSubpath) []core.PdfSubpath {
	result := make([]core.PdfSubpath, len(paths))
	for i, p := range paths {
		result[i] = *p
	}
	return result
}

var _ CidFont = (*Type0CidFont)(nil)
