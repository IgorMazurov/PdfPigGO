// Package simple provides concrete font implementations for PDF simple fonts.
package simple

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	enc "github.com/uglytoad/pdfpig/go/fonts/encodings"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Type3Font represents a Type 3 font in PDF.
// Type 3 fonts define their glyphs using PDF graphics operators rather than
// embedded font programs, making them unique among simple font types.
type Type3Font struct {
	name          *tokens.NameToken
	boundingBox   core.PdfRectangle
	fontMatrix    core.TransformationMatrix
	encoding      *enc.Encoding
	firstChar     int
	lastChar      int
	widths        []float64
	toUnicodeCMap fonts.CMapProvider
	descent       float64
	ascent        float64
}

// NewType3Font creates a new Type3Font instance.
func NewType3Font(
	name *tokens.NameToken,
	boundingBox core.PdfRectangle,
	fontMatrix core.TransformationMatrix,
	encoding *enc.Encoding,
	firstChar int,
	lastChar int,
	widths []float64,
	toUnicodeCMap fonts.CMapProvider,
) (*Type3Font, error) {
	if encoding == nil {
		return nil, fonts.NewInvalidFontFormatException("Type 3 font requires an encoding.")
	}

	ascent := fontMatrix.TransformY(boundingBox.Top())

	f := &Type3Font{
		name:          name,
		boundingBox:   boundingBox,
		fontMatrix:    fontMatrix,
		encoding:      encoding,
		firstChar:     firstChar,
		lastChar:      lastChar,
		widths:        widths,
		toUnicodeCMap: toUnicodeCMap,
		ascent:        ascent,
	}

	return f, nil
}

// Name returns the name of the font.
func (f *Type3Font) Name() *tokens.NameToken {
	return f.name
}

// IsVertical reports whether the font is used for vertical text layout.
// Type 3 fonts are never vertical.
func (f *Type3Font) IsVertical() bool {
	return false
}

// Details returns the default font details for this font.
func (f *Type3Font) Details() *fonts.FontDetails {
	name := ""
	if f.name != nil {
		name = f.name.Data()
	}
	return fonts.GetDefault(name)
}

// ReadCharacterCode reads a single byte as the character code.
func (f *Type3Font) ReadCharacterCode(bytes core.InputBytes) (int, int) {
	return int(bytes.CurrentByte()), 1
}

// TryGetUnicode attempts to get the Unicode string for the given character code.
func (f *Type3Font) TryGetUnicode(characterCode int) (string, bool) {
	if f.toUnicodeCMap != nil {
		if unicode, ok := f.toUnicodeCMap.TryConvertToUnicode(characterCode); ok {
			return unicode, true
		}
	}

	if f.encoding == nil {
		return "", false
	}

	glyphName := f.encoding.GetName(characterCode)
	if glyphName == "" {
		return "", false
	}

	unicode, ok := fonts.GlyphListNameToUnicode(glyphName)
	return unicode, ok
}

// GetBoundingBox returns the bounding box for the given character code.
func (f *Type3Font) GetBoundingBox(characterCode int) (*fonts.CharacterBoundingBox, error) {
	if characterCode < f.firstChar || characterCode > f.lastChar {
		return nil, fmt.Errorf("character code %d out of range [%d, %d]", characterCode, f.firstChar, f.lastChar)
	}

	width := f.widths[characterCode-f.firstChar]
	glyphSpaceRect := core.NewPdfRectangleFloat(0, 0, width, f.boundingBox.Height)
	transformedRect := f.fontMatrix.TransformRect(glyphSpaceRect)
	transformedWidth := f.fontMatrix.TransformX(width)

	return fonts.NewCharacterBoundingBox(transformedRect, transformedWidth), nil
}

// GetFontMatrix returns the transformation matrix for this font.
func (f *Type3Font) GetFontMatrix() core.TransformationMatrix {
	return f.fontMatrix
}

// GetDescent returns the descent of the font adjusted by the font matrix.
// Type 3 fonts have a fixed descent of 0.
func (f *Type3Font) GetDescent() float64 {
	return f.descent
}

// GetAscent returns the ascent of the font adjusted by the font matrix.
func (f *Type3Font) GetAscent() float64 {
	return f.ascent
}

// TryGetPath always returns false since Type 3 fonts do not use vector paths.
func (f *Type3Font) TryGetPath(characterCode int) ([]*core.PdfSubpath, bool) {
	return nil, false
}

// TryGetNormalisedPath always returns false since Type 3 fonts do not use vector paths.
func (f *Type3Font) TryGetNormalisedPath(characterCode int) ([]*core.PdfSubpath, bool) {
	return nil, false
}

var _ fonts.Font = (*Type3Font)(nil)
