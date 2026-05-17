// Package simple provides simple font implementations for PDF text rendering.
package simple

import (
	"log"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/adobe_font_metrics"
	"github.com/uglytoad/pdfpig/go/fonts/encodings"
	truetypeparser "github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// defaultTransformation is the standard font matrix for Type 42 fonts.
var defaultTransformation = core.FromValues(1/1000.0, 0, 0, 1/1000.0, 0, 0)

// trueTypeStandard14FallbackSimpleFont handles TrueType fonts that use both
// the Standard 14 descriptor and the TrueType font from disk.
type trueTypeStandard14FallbackSimpleFont struct {
	fontMatrix  core.TransformationMatrix
	ascent      float64
	descent     float64
	fontMetrics *adobe_font_metrics.AdobeFontMetrics
	encoding    *encodings.Encoding
	font        *truetypeparser.TrueTypeFont
	overrides   *metricOverrides

	name    *tokens.NameToken
	details *fonts.FontDetails
}

// NewTrueTypeStandard14FallbackSimpleFont creates a new TrueType Standard 14 fallback simple font.
func NewTrueTypeStandard14FallbackSimpleFont(
	name *tokens.NameToken,
	fontMetrics *adobe_font_metrics.AdobeFontMetrics,
	encoding *encodings.Encoding,
	font *truetypeparser.TrueTypeFont,
	overrides *metricOverrides,
) (*trueTypeStandard14FallbackSimpleFont, error) {
	if encoding == nil {
		return nil, errNilEncoding
	}

	f := &trueTypeStandard14FallbackSimpleFont{
		fontMetrics: fontMetrics,
		encoding:    encoding,
		font:        font,
		overrides:   overrides,
		name:        name,
	}

	fontName := ""
	if name != nil {
		fontName = name.Data()
	}

	if fontMetrics == nil {
		f.details = fonts.GetDefault(fontName)
	} else {
		isBold := fontMetrics.Weight == "Bold"
		weight := fonts.DefaultWeight
		if isBold {
			weight = fonts.BoldWeight
		}
		isItalic := fontMetrics.ItalicAngle != 0
		f.details = fonts.NewFontDetails(fontName, isBold, weight, isItalic)
	}

	_ = encoding // ZapfDingbats assertion omitted per Go conventions

	if font != nil && font.TableRegister().HeaderTable.Tag() != "" {
		scale := float64(font.GetUnitsPerEm())
		f.fontMatrix = core.FromValues(1.0/scale, 0, 0, 1.0/scale, 0, 0)
	} else {
		f.fontMatrix = defaultTransformation
	}

	f.descent = f.computeDescent()
	f.ascent = f.computeAscent()

	return f, nil
}

func (f *trueTypeStandard14FallbackSimpleFont) computeDescent() float64 {
	if f.fontMetrics != nil {
		return f.fontMatrix.TransformY(f.fontMetrics.Descender)
	}
	return f.fontMatrix.TransformY(float64(f.font.TableRegister().HorizontalHeaderTable.Descent()))
}

func (f *trueTypeStandard14FallbackSimpleFont) computeAscent() float64 {
	if f.fontMetrics != nil {
		return f.fontMatrix.TransformY(f.fontMetrics.Ascender)
	}
	return f.fontMatrix.TransformY(float64(f.font.TableRegister().HorizontalHeaderTable.Ascent()))
}

// Name returns the name of the font.
func (f *trueTypeStandard14FallbackSimpleFont) Name() *tokens.NameToken {
	return f.name
}

// IsVertical reports whether the font is used for vertical text layout.
func (f *trueTypeStandard14FallbackSimpleFont) IsVertical() bool {
	return false
}

// Details returns the details associated with this font.
func (f *trueTypeStandard14FallbackSimpleFont) Details() *fonts.FontDetails {
	return f.details
}

// SetDetails sets the font details.
func (f *trueTypeStandard14FallbackSimpleFont) SetDetails(details *fonts.FontDetails) {
	f.details = details
}

// ReadCharacterCode reads the next character code from the given bytes.
func (f *trueTypeStandard14FallbackSimpleFont) ReadCharacterCode(bytes core.InputBytes) (int, int) {
	return int(bytes.CurrentByte()), 1
}

// TryGetUnicode attempts to get the Unicode string for the given character code.
func (f *trueTypeStandard14FallbackSimpleFont) TryGetUnicode(characterCode int) (string, bool) {
	encodedCharacterName := f.encoding.GetName(characterCode)

	glyphList, err := fonts.AdobeGlyphList()
	if err != nil || glyphList == nil {
		return "", false
	}

	value, err := glyphList.NameToUnicode(encodedCharacterName)
	if err != nil {
		log.Printf("warning: failed to get Unicode for character name %q: %v", encodedCharacterName, err)
		return "", false
	}

	return value, true
}

// GetBoundingBox returns the bounding box for the given character code.
func (f *trueTypeStandard14FallbackSimpleFont) GetBoundingBox(characterCode int) (*fonts.CharacterBoundingBox, error) {
	fontMatrix := f.fontMatrix

	if f.font != nil {
		bounds, ok := f.font.TryGetBoundingBox(characterCode)
		if ok {
			bounds = fontMatrix.TransformRect(bounds)

			width := 0.0
			hasOverride := false
			if f.overrides != nil {
				var w float64
				hasOverride = f.overrides.tryGetWidth(characterCode, &w)
				if hasOverride {
					width = w
				}
			}

			if !hasOverride {
				encodedName := f.encoding.GetName(characterCode)
				if metric, exists := f.fontMetrics.CharacterMetrics[encodedName]; exists {
					width = defaultTransformation.TransformX(metric.Width.X)
				} else {
					width = bounds.Width
				}
			} else {
				width = defaultTransformation.TransformX(width)
			}

			return fonts.NewCharacterBoundingBox(bounds, width), nil
		}
	}

	name := f.encoding.GetName(characterCode)
	metric, exists := f.fontMetrics.CharacterMetrics[name]
	if !exists {
		return fonts.NewCharacterBoundingBox(core.PdfRectangle{}, 0), nil
	}

	width := 0.0
	hasOverride := false
	if f.overrides != nil {
		var w float64
		hasOverride = f.overrides.tryGetWidth(characterCode, &w)
		if hasOverride {
			width = w
		}
	}

	if !hasOverride {
		width = fontMatrix.TransformX(metric.Width.X)
	} else {
		width = defaultTransformation.TransformX(width)
	}

	bounds := fontMatrix.TransformRect(metric.BoundingBox)

	return fonts.NewCharacterBoundingBox(bounds, width), nil
}

// GetFontMatrix returns the transformation matrix for this font.
func (f *trueTypeStandard14FallbackSimpleFont) GetFontMatrix() core.TransformationMatrix {
	return f.fontMatrix
}

// GetDescent returns the descent of the font adjusted by the font matrix.
func (f *trueTypeStandard14FallbackSimpleFont) GetDescent() float64 {
	return f.descent
}

// GetAscent returns the ascent of the font adjusted by the font matrix.
func (f *trueTypeStandard14FallbackSimpleFont) GetAscent() float64 {
	return f.ascent
}

// TryGetPath attempts to get the glyph path for the given character code.
func (f *trueTypeStandard14FallbackSimpleFont) TryGetPath(characterCode int) ([]*core.PdfSubpath, bool) {
	if f.font == nil {
		return nil, false
	}
	return f.font.TryGetPath(characterCode)
}

// TryGetNormalisedPath attempts to get the normalised glyph path for the given character code.
func (f *trueTypeStandard14FallbackSimpleFont) TryGetNormalisedPath(characterCode int) ([]*core.PdfSubpath, bool) {
	path, ok := f.TryGetPath(characterCode)
	if !ok {
		return nil, false
	}

	transformed, err := f.fontMatrix.TransformPath(path)
	if err != nil {
		log.Printf("warning: failed to transform path: %v", err)
		return nil, false
	}

	return transformed, true
}

// metricOverrides holds optional width overrides for character codes.
type metricOverrides struct {
	firstCharacterCode *int
	widths             []float64
	hasOverridden      bool
}

// NewMetricOverrides creates a new MetricOverrides instance.
func NewMetricOverrides(firstCharacterCode *int, widths []float64) *metricOverrides {
	hasOverridden := firstCharacterCode != nil && widths != nil && len(widths) > 0
	return &metricOverrides{
		firstCharacterCode: firstCharacterCode,
		widths:             widths,
		hasOverridden:      hasOverridden,
	}
}

// tryGetWidth attempts to get an overridden width for the given character code.
func (m *metricOverrides) tryGetWidth(characterCode int, width *float64) bool {
	*width = 0

	if !m.hasOverridden || m.firstCharacterCode == nil {
		return false
	}

	index := characterCode - *m.firstCharacterCode

	if index < 0 || index >= len(m.widths) {
		return false
	}

	*width = m.widths[index]
	return true
}

var _ fonts.Font = (*trueTypeStandard14FallbackSimpleFont)(nil)

var errNilEncoding = err("encoding must not be nil")

type err string

func (e err) Error() string {
	return string(e)
}
