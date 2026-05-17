// Package simple provides simple font implementations for PDF text rendering.
package simple

import (
	"fmt"
	"math"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/adobe_font_metrics"
	enc "github.com/uglytoad/pdfpig/go/fonts/encodings"
	cff "github.com/uglytoad/pdfpig/go/fonts/cff"
	parser "github.com/uglytoad/pdfpig/go/fonts/type1/parser"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Type1Standard14Font represents a font using one of the Adobe Standard 14 fonts.
// Can use a custom encoding.
type Type1Standard14Font struct {
	standardFontMetrics adobe_font_metrics.AdobeFontMetrics
	encoding            *enc.Encoding
	isZapfDingbats      bool
	name                *tokens.NameToken
	isVertical          bool
	details             *fonts.FontDetails
	fontMatrix          core.TransformationMatrix
	ascent              float64
	descent             float64
}

// NewType1Standard14Font creates a new Type1Standard14Font from Standard 14 metrics.
func NewType1Standard14Font(
	metrics adobe_font_metrics.AdobeFontMetrics,
	overrideEncoding *enc.Encoding,
) (*Type1Standard14Font, error) {
	if metrics.FontName == "" {
		return nil, fmt.Errorf("standard font metrics cannot be empty")
	}

	var encVal *enc.Encoding
	if overrideEncoding != nil {
		encVal = overrideEncoding
	} else {
		afmEnc, err := adobe_font_metrics.NewAdobeFontMetricsEncoding(&metrics)
		if err != nil {
			return nil, err
		}
		encVal = afmEnc.Encoding
	}

	nameToken := tokens.Create(metrics.FontName)

	isBold := metrics.Weight == "Bold"
	details := fonts.NewFontDetails(
		nameToken.Data(),
		isBold,
		func() int {
			if isBold {
				return 700
			}
			return fonts.DefaultWeight
		}(),
		metrics.ItalicAngle != 0,
	)

	isZapfDingbats := details.Name != "" && strings.Contains(details.Name, "ZapfDingbats")
	if !isZapfDingbats && encVal != nil {
		isZapfDingbats = strings.EqualFold(encVal.EncodingName(), "ZapfDingbatsEncoding")
	}

	fontMatrix := core.FromValues(0.001, 0, 0, 0.001, 0, 0)

	f := &Type1Standard14Font{
		standardFontMetrics: metrics,
		encoding:            encVal,
		isZapfDingbats:      isZapfDingbats,
		name:                nameToken,
		isVertical:          false,
		details:             details,
		fontMatrix:          fontMatrix,
	}

	f.descent = f.computeDescent()
	f.ascent = f.computeAscent()

	return f, nil
}

func (f *Type1Standard14Font) computeDescent() float64 {
	if math.Abs(f.standardFontMetrics.Descender) < 1e-9 {
		return -0.25
	}
	return f.fontMatrix.TransformY(f.standardFontMetrics.Descender)
}

func (f *Type1Standard14Font) computeAscent() float64 {
	if math.Abs(f.standardFontMetrics.Ascender) < 1e-9 {
		return 0.75
	}
	return f.fontMatrix.TransformY(f.standardFontMetrics.Ascender)
}

// Name returns the font name token.
func (f *Type1Standard14Font) Name() *tokens.NameToken {
	return f.name
}

// IsVertical always returns false for Standard 14 fonts.
func (f *Type1Standard14Font) IsVertical() bool {
	return f.isVertical
}

// Details returns the font details.
func (f *Type1Standard14Font) Details() *fonts.FontDetails {
	return f.details
}

// ReadCharacterCode reads a single-byte character code and returns the code
// along with the number of bytes consumed (always 1).
func (f *Type1Standard14Font) ReadCharacterCode(bytes core.InputBytes) (int, int) {
	return int(bytes.CurrentByte()), 1
}

// TryGetUnicode attempts to get the Unicode string for the given character code.
// Returns true and the Unicode value if successful, false otherwise.
func (f *Type1Standard14Font) TryGetUnicode(characterCode int) (string, bool) {
	name := f.encoding.GetName(characterCode)

	if strings.EqualFold(name, fonts.NotDefined) {
		return "", false
	}

	var value string
	var ok bool

	if f.isZapfDingbats {
		zapf, err := fonts.ZapfDingbats()
		if err == nil && zapf != nil {
			value, ok = tryNameToUnicode(zapf, name)
			if ok {
				return value, true
			}
		}
	}

	adobe, err := fonts.AdobeGlyphList()
	if err != nil || adobe == nil {
		return "", false
	}

	value, ok = tryNameToUnicode(adobe, name)
	return value, ok
}

func tryNameToUnicode(gl *fonts.GlyphList, name string) (string, bool) {
	value, err := gl.NameToUnicode(name)
	if err != nil {
		return "", false
	}
	return value, value != ""
}

// GetBoundingBox returns the bounding box for the given character code.
func (f *Type1Standard14Font) GetBoundingBox(characterCode int) (*fonts.CharacterBoundingBox, error) {
	boundingBox, advanceWidth := f.getBoundingBoxInGlyphSpace(characterCode)

	boundingBox = f.fontMatrix.TransformRect(boundingBox)
	advanceWidth = f.fontMatrix.TransformX(advanceWidth)

	return fonts.NewCharacterBoundingBox(boundingBox, advanceWidth), nil
}

func (f *Type1Standard14Font) getBoundingBoxInGlyphSpace(characterCode int) (core.PdfRectangle, float64) {
	name := f.encoding.GetName(characterCode)

	metrics, ok := f.standardFontMetrics.CharacterMetrics[name]
	if !ok {
		return core.NewPdfRectangleFloat(0, 0, 250, 0), 250
	}

	x := metrics.Width.X

	if x == 0 && metrics.BoundingBox.Width > 0 {
		x = metrics.BoundingBox.Width
	}

	return metrics.BoundingBox, x
}

// GetFontMatrix returns the transformation matrix for this font.
func (f *Type1Standard14Font) GetFontMatrix() core.TransformationMatrix {
	return f.fontMatrix
}

// GetDescent returns the descent of the font adjusted by the font matrix.
func (f *Type1Standard14Font) GetDescent() float64 {
	return f.descent
}

// GetAscent returns the ascent of the font adjusted by the font matrix.
func (f *Type1Standard14Font) GetAscent() float64 {
	return f.ascent
}

// TryGetPath attempts to get the glyph path for the given character code.
// Not implemented for Standard 14 fonts.
func (f *Type1Standard14Font) TryGetPath(characterCode int) ([]*core.PdfSubpath, bool) {
	return nil, false
}

// TryGetNormalisedPath attempts to get the normalised glyph path.
// Not implemented for Standard 14 fonts.
func (f *Type1Standard14Font) TryGetNormalisedPath(characterCode int) ([]*core.PdfSubpath, bool) {
	return f.TryGetPath(characterCode)
}

var _ fonts.Font = (*Type1Standard14Font)(nil)

// Type1FontProgram holds either a Type1Font or a CompactFontFormatFontCollection.
type Type1FontProgram struct {
	T1  *parser.Type1Font
	Cff *cff.CompactFontFormatFontCollection
}

// Type1FontSimple is a font based on the Adobe Type 1 font format.
type Type1FontSimple struct {
	name             *tokens.NameToken
	firstChar        int
	lastChar         int
	widths           []float64
	descriptor       *fonts.FontDescriptor
	encoding         *enc.Encoding
	toUnicodeCMap    fonts.CMapProvider
	fontProgram      Type1FontProgram
	fontMatrix       core.TransformationMatrix
	ascent           float64
	descent          float64
	details          *fonts.FontDetails
	isZapfDingbats   bool
	cachedBoundingBoxes map[int]*fonts.CharacterBoundingBox
}

// NewType1FontSimple creates a new Type1FontSimple from the given parameters.
func NewType1FontSimple(
	name *tokens.NameToken,
	firstChar int,
	lastChar int,
	widths []float64,
	descriptor *fonts.FontDescriptor,
	encoding *enc.Encoding,
	toUnicodeCMap fonts.CMapProvider,
	fontProgram Type1FontProgram,
) (*Type1FontSimple, error) {
	matrix := core.FromValues(0.001, 0, 0, 0.001, 0, 0)

	if fontProgram.T1 != nil {
		matrix = fontProgram.T1.FontMatrix
	} else if fontProgram.Cff != nil {
		matrix = fontProgram.Cff.GetFirstTransformationMatrix()
	}

	f := &Type1FontSimple{
		name:             name,
		firstChar:        firstChar,
		lastChar:         lastChar,
		widths:           widths,
		descriptor:       descriptor,
		encoding:         encoding,
		toUnicodeCMap:    toUnicodeCMap,
		fontProgram:      fontProgram,
		fontMatrix:       matrix,
		cachedBoundingBoxes: make(map[int]*fonts.CharacterBoundingBox),
	}

	fontName := ""
	if name != nil {
		fontName = name.Data()
	}

	if descriptor != nil {
		f.details = descriptor.ToDetails(fontName)
	} else {
		f.details = fonts.GetDefault(fontName)
	}

	f.isZapfDingbats = f.checkZapfDingbats(encoding) || (f.details != nil && strings.Contains(f.details.Name, "ZapfDingbats"))

	f.descent = f.computeDescent()
	f.ascent = f.computeAscent()

	return f, nil
}

func (f *Type1FontSimple) checkZapfDingbats(e *enc.Encoding) bool {
	if e == nil {
		return false
	}
	return strings.EqualFold(e.EncodingName(), "ZapfDingbats")
}

func (f *Type1FontSimple) computeDescent() float64 {
	if f.descriptor != nil && math.Abs(f.descriptor.Descent) > 1e-9 {
		return f.fontMatrix.TransformY(f.descriptor.Descent)
	}
	return -0.25
}

func (f *Type1FontSimple) computeAscent() float64 {
	if f.descriptor != nil && math.Abs(f.descriptor.Ascent) > 1e-9 {
		return f.fontMatrix.TransformY(f.descriptor.Ascent)
	}
	return 0.75
}

// Name returns the font name token.
func (f *Type1FontSimple) Name() *tokens.NameToken {
	return f.name
}

// IsVertical always returns false for simple Type 1 fonts.
func (f *Type1FontSimple) IsVertical() bool {
	return false
}

// Details returns the font details.
func (f *Type1FontSimple) Details() *fonts.FontDetails {
	if f.details != nil {
		return f.details
	}
	fontName := ""
	if f.name != nil {
		fontName = f.name.Data()
	}
	f.details = fonts.GetDefault(fontName)
	return f.details
}

// ReadCharacterCode reads a single-byte character code.
func (f *Type1FontSimple) ReadCharacterCode(bytes core.InputBytes) (int, int) {
	return int(bytes.CurrentByte()), 1
}

// TryGetUnicode attempts to get the Unicode string for the given character code.
func (f *Type1FontSimple) TryGetUnicode(characterCode int) (string, bool) {
	if f.toUnicodeCMap != nil {
		if value, ok := f.toUnicodeCMap.TryConvertToUnicode(characterCode); ok {
			return value, true
		}
	}

	if f.encoding == nil {
		value := unicodeFromRune(rune(characterCode))
		if value != "" {
			return value, true
		}

		if f.fontProgram.T1 != nil {
			if name, ok := f.fontProgram.T1.Encoding[characterCode]; ok {
				return name, true
			}
		}

		return "", false
	}

	name := f.encoding.GetName(characterCode)
	value, err := f.nameToUnicode(name)
	if err != nil {
		return "", false
	}

	return value, value != ""
}

func unicodeFromRune(r rune) string {
	if r > 0x10FFFF {
		return ""
	}
	if r >= 0xD800 && r <= 0xDFFF {
		return ""
	}
	return string(r)
}

// nameToUnicode converts a glyph name to its Unicode string, checking ZapfDingbats first if applicable.
func (f *Type1FontSimple) nameToUnicode(name string) (string, error) {
	if f.isZapfDingbats {
		zapf, err := fonts.ZapfDingbats()
		if err == nil && zapf != nil {
			value, err2 := zapf.NameToUnicode(name)
			if err2 == nil && value != "" {
				return value, nil
			}
		}
	}

	adobe, err := fonts.AdobeGlyphList()
	if err != nil || adobe == nil {
		return "", nil
	}

	value, err := adobe.NameToUnicode(name)
	if err != nil {
		return "", err
	}

	return value, nil
}

// UnicodeCodePointToName returns the glyph name for a given Unicode code point value.
func (f *Type1FontSimple) UnicodeCodePointToName(unicodeValue int) string {
	if f.isZapfDingbats {
		zapf, err := fonts.ZapfDingbats()
		if err == nil && zapf != nil {
			value := zapf.UnicodeCodePointToName(unicodeValue)
			if value != fonts.NotDefined {
				return value
			}
		}
	}

	adobe, err := fonts.AdobeGlyphList()
	if err != nil || adobe == nil {
		return fonts.NotDefined
	}

	return adobe.UnicodeCodePointToName(unicodeValue)
}

// GetBoundingBox returns the bounding box for the given character code, using caching.
func (f *Type1FontSimple) GetBoundingBox(characterCode int) (*fonts.CharacterBoundingBox, error) {
	if cached, ok := f.cachedBoundingBoxes[characterCode]; ok {
		return cached, nil
	}

	boundingBox := f.getBoundingBoxInGlyphSpace(characterCode)
	boundingBox = f.fontMatrix.TransformRect(boundingBox)

	width := f.getWidth(characterCode, boundingBox)
	result := fonts.NewCharacterBoundingBox(boundingBox, width/1000.0)

	f.cachedBoundingBoxes[characterCode] = result

	return result, nil
}

// GetFontMatrix returns the transformation matrix for this font.
func (f *Type1FontSimple) GetFontMatrix() core.TransformationMatrix {
	return f.fontMatrix
}

// GetDescent returns the descent adjusted by the font matrix.
func (f *Type1FontSimple) GetDescent() float64 {
	return f.descent
}

// GetAscent returns the ascent adjusted by the font matrix.
func (f *Type1FontSimple) GetAscent() float64 {
	return f.ascent
}

// TryGetPath attempts to get the glyph path for the given character code.
func (f *Type1FontSimple) TryGetPath(characterCode int) ([]*core.PdfSubpath, bool) {
	if characterCode < f.firstChar || characterCode > f.lastChar {
		return nil, false
	}

	if f.fontProgram.T1 == nil && f.fontProgram.Cff == nil {
		return nil, false
	}

	var tempPath []*core.PdfSubpath

	if f.fontProgram.T1 != nil {
		name := f.getGlyphName(characterCode)
		tempPath = f.fontProgram.T1.GetCharacterPath(name)
	} else if f.fontProgram.Cff != nil {
		firstFont := f.fontProgram.Cff.FirstFont()
		characterName := f.getGlyphNameForCFF(characterCode, firstFont)
		if firstFont != nil {
			tempPath = firstFont.GetCharacterPath(characterName)
		}
	}

	if len(tempPath) > 0 {
		return tempPath, true
	}

	return nil, false
}

// TryGetNormalisedPath attempts to get the normalised glyph path.
func (f *Type1FontSimple) TryGetNormalisedPath(characterCode int) ([]*core.PdfSubpath, bool) {
	path, ok := f.TryGetPath(characterCode)
	if !ok {
		return nil, false
	}
	transformed, err := f.GetFontMatrix().TransformPath(path)
	if err != nil {
		return nil, false
	}
	return transformed, true
}

func (f *Type1FontSimple) getWidth(characterCode int, boundingBox core.PdfRectangle) float64 {
	widthIndex := characterCode - f.firstChar

	if widthIndex >= 0 && widthIndex < len(f.widths) {
		return f.widths[widthIndex]
	}

	if f.descriptor != nil {
		return f.descriptor.MissingWidth
	}

	return boundingBox.Width
}

func (f *Type1FontSimple) getBoundingBoxInGlyphSpace(characterCode int) core.PdfRectangle {
	if characterCode < f.firstChar || characterCode > f.lastChar {
		return core.NewPdfRectangleFloat(0, 0, 250, 0)
	}

	if f.fontProgram.T1 == nil && f.fontProgram.Cff == nil {
		return core.NewPdfRectangleFloat(0, 0, f.widths[characterCode-f.firstChar], 0)
	}

	var rect *core.PdfRectangle

	if f.fontProgram.T1 != nil {
		name := f.getGlyphName(characterCode)
		rect = f.fontProgram.T1.GetCharacterBoundingBox(name)
	} else if f.fontProgram.Cff != nil {
		firstFont := f.fontProgram.Cff.FirstFont()
		characterName := f.getGlyphNameForCFF(characterCode, firstFont)
		if firstFont != nil {
			rect = firstFont.GetCharacterBoundingBox(characterName)
		}
	}

	if rect == nil {
		return core.NewPdfRectangleFloat(0, 0, f.widths[characterCode-f.firstChar], 0)
	}

	return *rect
}

func (f *Type1FontSimple) getGlyphName(characterCode int) string {
	if f.encoding != nil {
		return f.encoding.GetName(characterCode)
	}
	return f.UnicodeCodePointToName(characterCode)
}

func (f *Type1FontSimple) getGlyphNameForCFF(characterCode int, firstFont *cff.CompactFontFormatFont) string {
	if f.encoding != nil {
		return f.encoding.GetName(characterCode)
	}
	if firstFont != nil {
		return firstFont.GetCharacterName(characterCode, false)
	}
	return fonts.NotDefined
}

var _ fonts.Font = (*Type1FontSimple)(nil)
