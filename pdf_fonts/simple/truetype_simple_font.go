// Package simple provides simple font implementations for PDF text rendering.
package simple

import (
	"fmt"
	"strings"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	enc "github.com/uglytoad/pdfpig/go/fonts/encodings"
	truetypeparser "github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// toUnicodeProvider wraps fonts.CMapProvider with CanMapToUnicode check.
// The C# ToUnicodeCMap has both CanMapToUnicode and TryGet methods; the Go
// CMapProvider interface only exposes TryConvertToUnicode which combines both.
type toUnicodeProvider struct {
	cmap fonts.CMapProvider
}

func newToUnicodeProvider(cmap fonts.CMapProvider) *toUnicodeProvider {
	return &toUnicodeProvider{cmap: cmap}
}

func (t *toUnicodeProvider) CanMapToUnicode() bool {
	return t.cmap != nil
}

func (t *toUnicodeProvider) TryGet(code int) (string, bool) {
	if t.cmap == nil {
		return "", false
	}
	return t.cmap.TryConvertToUnicode(code)
}

// trueTypeSimpleFont represents a TrueType-based simple font in PDF.
type trueTypeSimpleFont struct {
	name           *tokens.NameToken
	descriptor     *fonts.FontDescriptor
	toUnicode      *toUnicodeProvider
	encoding       *enc.Encoding
	font           *truetypeparser.TrueTypeFont
	firstCharacter int
	widths         []float64
	isZapfDingbats bool
	fontMatrix     core.TransformationMatrix
	descent        float64
	ascent         float64
	details        *fonts.FontDetails

	boundingBoxCache map[int]*fonts.CharacterBoundingBox
	bboxMu           sync.RWMutex

	unicodeValuesCache map[int]string
	unicodeMu          sync.RWMutex
}

// NewTrueTypeSimpleFont creates a new TrueType simple font.
func NewTrueTypeSimpleFont(
	name *tokens.NameToken,
	descriptor *fonts.FontDescriptor,
	toUnicodeCMap fonts.CMapProvider,
	encoding *enc.Encoding,
	font *truetypeparser.TrueTypeFont,
	firstCharacter int,
	widths []float64,
) (*trueTypeSimpleFont, error) {
	f := &trueTypeSimpleFont{
		name:               name,
		descriptor:         descriptor,
		toUnicode:          newToUnicodeProvider(toUnicodeCMap),
		encoding:           encoding,
		font:               font,
		firstCharacter:     firstCharacter,
		widths:             widths,
		boundingBoxCache:   make(map[int]*fonts.CharacterBoundingBox),
		unicodeValuesCache: make(map[int]string),
	}

	f.init(name, descriptor, encoding)
	return f, nil
}

func (f *trueTypeSimpleFont) init(
	name *tokens.NameToken,
	descriptor *fonts.FontDescriptor,
	encoding *enc.Encoding,
) {
	fontName := ""
	if name != nil {
		fontName = name.Data()
	}

	if descriptor != nil {
		f.details = descriptor.ToDetails(fontName)
	} else {
		f.details = fonts.GetDefault(fontName)
	}

	f.isZapfDingbats = encoding != nil && encoding.EncodingName() == "ZapfDingbatsEncoding" ||
		strings.Contains(f.details.Name, "ZapfDingbats")

	scale := 1000.0
	if f.font != nil {
		headerTable := f.font.TableRegister().HeaderTable
		if headerTable.Tag() != "" {
			scale = float64(f.font.GetUnitsPerEm())
		}
	}

	f.fontMatrix = core.FromValues(1.0/scale, 0, 0, 1.0/scale, 0, 0)
	f.descent = f.computeDescent()
	f.ascent = f.computeAscent()
}

func (f *trueTypeSimpleFont) computeDescent() float64 {
	if f.font == nil {
		return defaultTransformation.TransformY(f.descriptor.Descent)
	}
	return f.fontMatrix.TransformY(float64(f.font.TableRegister().HorizontalHeaderTable.Descent()))
}

func (f *trueTypeSimpleFont) computeAscent() float64 {
	if f.font == nil {
		return defaultTransformation.TransformY(f.descriptor.Ascent)
	}
	return f.fontMatrix.TransformY(float64(f.font.TableRegister().HorizontalHeaderTable.Ascent()))
}

// Name returns the name of the font.
func (f *trueTypeSimpleFont) Name() *tokens.NameToken {
	return f.name
}

// IsVertical reports whether the font is used for vertical text layout.
func (f *trueTypeSimpleFont) IsVertical() bool {
	return false
}

// Details returns the details associated with this font.
func (f *trueTypeSimpleFont) Details() *fonts.FontDetails {
	return f.details
}

// ToUnicode returns the ToUnicode CMap provider for this font.
func (f *trueTypeSimpleFont) ToUnicode() fonts.CMapProvider {
	if f.toUnicode != nil {
		return f.toUnicode.cmap
	}
	return nil
}

// SetToUnicode sets the ToUnicode CMap provider for this font.
func (f *trueTypeSimpleFont) SetToUnicode(cmap fonts.CMapProvider) {
	f.toUnicode = newToUnicodeProvider(cmap)
}

// ReadCharacterCode reads the next character code from the given bytes.
func (f *trueTypeSimpleFont) ReadCharacterCode(bytes core.InputBytes) (int, int) {
	return int(bytes.CurrentByte()), 1
}

// TryGetUnicode attempts to get the Unicode string for the given character code.
func (f *trueTypeSimpleFont) TryGetUnicode(characterCode int) (string, bool) {
	f.unicodeMu.RLock()
	if value, ok := f.unicodeValuesCache[characterCode]; ok {
		f.unicodeMu.RUnlock()
		return value, true
	}
	f.unicodeMu.RUnlock()

	if f.toUnicode.CanMapToUnicode() {
		if value, ok := f.toUnicode.TryGet(characterCode); ok {
			f.unicodeMu.Lock()
			f.unicodeValuesCache[characterCode] = value
			f.unicodeMu.Unlock()
			return value, true
		}
	}

	if f.encoding == nil {
		return "", false
	}

	name := f.encoding.GetName(characterCode)

	if f.isZapfDingbats {
		zg, err := fonts.ZapfDingbats()
		if err == nil && zg != nil {
			value, err := zg.NameToUnicode(name)
			if err == nil && value != "" {
				f.unicodeMu.Lock()
				f.unicodeValuesCache[characterCode] = value
				f.unicodeMu.Unlock()
				return value, true
			}
		}
	}

	ag, err := fonts.AdobeGlyphList()
	if err != nil || ag == nil {
		return "", false
	}

	value, err := ag.NameToUnicode(name)
	if err != nil {
		return "", false
	}

	if value != "" {
		f.unicodeMu.Lock()
		f.unicodeValuesCache[characterCode] = value
		f.unicodeMu.Unlock()
	}

	return value, value != ""
}

// GetBoundingBox returns the bounding box for the given character code.
func (f *trueTypeSimpleFont) GetBoundingBox(characterCode int) (*fonts.CharacterBoundingBox, error) {
	f.bboxMu.RLock()
	if cached, ok := f.boundingBoxCache[characterCode]; ok {
		f.bboxMu.RUnlock()
		return cached, nil
	}
	f.bboxMu.RUnlock()

	boundingRect, fromFont := f.getBoundingBoxInGlyphSpace(characterCode)
	preTransformWidth := boundingRect.Width

	if fromFont {
		boundingRect = f.fontMatrix.TransformRect(boundingRect)
	} else {
		boundingRect = defaultTransformation.TransformRect(boundingRect)
	}

	var width float64

	index := characterCode - f.firstCharacter
	if f.widths != nil && index >= 0 && index < len(f.widths) {
		fromFont = false
		width = f.widths[index]
	} else if f.font != nil {
		w, ok := f.font.TryGetAdvanceWidth(characterCode)
		if !ok {
			width = preTransformWidth
		} else {
			width = w
		}
	} else if len(f.widths) > 0 {
		width = f.widths[0]
	} else {
		fontName := "<unknown>"
		if f.name != nil {
			fontName = f.name.Data()
		}
		return nil, fmt.Errorf("could not retrieve width for character code: %d in font %s", characterCode, fontName)
	}

	if fromFont {
		width = f.fontMatrix.TransformX(width)
	} else {
		width = defaultTransformation.TransformX(width)
	}

	result := fonts.NewCharacterBoundingBox(boundingRect, width)

	f.bboxMu.Lock()
	f.boundingBoxCache[characterCode] = result
	f.bboxMu.Unlock()

	return result, nil
}

// GetFontMatrix returns the transformation matrix for this font.
func (f *trueTypeSimpleFont) GetFontMatrix() core.TransformationMatrix {
	return f.fontMatrix
}

// GetDescent returns the descent of the font adjusted by the font matrix.
func (f *trueTypeSimpleFont) GetDescent() float64 {
	return f.descent
}

// GetAscent returns the ascent of the font adjusted by the font matrix.
func (f *trueTypeSimpleFont) GetAscent() float64 {
	return f.ascent
}

// TryGetPath attempts to get the glyph path for the given character code.
func (f *trueTypeSimpleFont) TryGetPath(characterCode int) ([]*core.PdfSubpath, bool) {
	if f.font == nil {
		return nil, false
	}
	return f.font.TryGetPathWithMapping(characterCode, f.characterCodeToGlyphId)
}

// TryGetNormalisedPath attempts to get the normalised glyph path for the given character code.
func (f *trueTypeSimpleFont) TryGetNormalisedPath(characterCode int) ([]*core.PdfSubpath, bool) {
	path, ok := f.TryGetPath(characterCode)
	if !ok {
		return nil, false
	}
	transformed, err := f.fontMatrix.TransformPath(path)
	if err != nil {
		return nil, false
	}
	return transformed, true
}

func (f *trueTypeSimpleFont) getBoundingBoxInGlyphSpace(characterCode int) (core.PdfRectangle, bool) {
	if f.font == nil {
		return f.descriptor.BoundingBox, true
	}

	if bounds, ok := f.font.TryGetBoundingBoxWithMapping(characterCode, f.characterCodeToGlyphId); ok {
		return bounds, true
	}

	if width, ok := f.font.TryGetAdvanceWidthWithMapping(characterCode, f.characterCodeToGlyphId); ok {
		return core.NewPdfRectangleFloat(0, 0, width, 0), true
	}

	return core.NewPdfRectangleFloat(0, 0, f.getWidth(characterCode), 0), false
}

func (f *trueTypeSimpleFont) characterCodeToGlyphId(characterCode int) *int {
	if f.descriptor == nil || f.font == nil || f.encoding == nil {
		return nil
	}

	f.unicodeMu.RLock()
	unicode, hasUnicode := f.unicodeValuesCache[characterCode]
	f.unicodeMu.RUnlock()

	if !hasUnicode {
		return nil
	}

	tr := f.font.TableRegister()
	if tr.CMapTable == nil {
		return nil
	}

	name, ok := f.encoding.CodeToName[characterCode]
	if !ok || name == "" {
		return nil
	}

	if strings.EqualFold(name, fonts.NotDefined) {
		zero := 0
		return &zero
	}

	glyphId := 0

	if f.descriptor.Flags.HasFlag(fonts.Symbolic) && f.font.WindowsSymbolCMap() != nil {
		symbolCMap := f.font.WindowsSymbolCMap()
		glyphId = symbolCMap.CharacterCodeToGlyphIndex(characterCode)

		if glyphId == 0 && characterCode >= 0 && characterCode <= 0xFF {
			glyphId = symbolCMap.CharacterCodeToGlyphIndex(characterCode + 0xF000)
			if glyphId == 0 {
				glyphId = symbolCMap.CharacterCodeToGlyphIndex(characterCode + 0xF100)
			}
			if glyphId == 0 {
				glyphId = symbolCMap.CharacterCodeToGlyphIndex(characterCode + 0xF200)
			}
		}

		if glyphId == 0 && f.font.WindowsUnicodeCMap() != nil && unicode != "" {
			glyphId = f.font.WindowsUnicodeCMap().CharacterCodeToGlyphIndex(int(unicode[0]))
		}
	} else {
		if f.font.WindowsUnicodeCMap() != nil && unicode != "" {
			glyphId = f.font.WindowsUnicodeCMap().CharacterCodeToGlyphIndex(int(unicode[0]))
		}

		if glyphId == 0 && f.font.MacRomanCMap() != nil {
			macCode := enc.MacOsRomanEncodingValue.GetCode(name)
			if macCode >= 0 {
				glyphId = f.font.MacRomanCMap().CharacterCodeToGlyphIndex(macCode)
			}
		}

		if glyphId == 0 && tr.PostScriptTable.Tag() != "" {
			for i, glyphName := range tr.PostScriptTable.GlyphNames() {
				if strings.EqualFold(glyphName, name) {
					return &i
				}
			}
		}
	}

	if glyphId != 0 {
		return &glyphId
	}

	return nil
}

func (f *trueTypeSimpleFont) getWidth(characterCode int) float64 {
	index := characterCode - f.firstCharacter
	if index < 0 || index >= len(f.widths) {
		return f.descriptor.MissingWidth
	}
	return f.widths[index]
}

var _ fonts.Font = (*trueTypeSimpleFont)(nil)
