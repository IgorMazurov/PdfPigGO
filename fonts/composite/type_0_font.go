// Package composite provides types for handling composite (Type 0) fonts in PDFs.
package composite

import (
	"fmt"
	"math"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/cidfonts"
	"github.com/uglytoad/pdfpig/go/fonts/cmap"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// cidFontRequirements captures the methods Type0Font needs from a CID font.
// This interface is intentionally narrower than cidfonts.CidFont to allow
// adapter types that bridge handler interfaces without import cycles.
type cidFontRequirements interface {
	Details() fonts.FontDetails
	FontMatrix() core.TransformationMatrix
	GetDescent() float64
	GetAscent() float64
	GetWidthFromDictionary(cid int) float64
	GetWidthFromFont(characterIdentifier int) float64
	GetBoundingBox(characterIdentifier int) (core.PdfRectangle, error)
	GetPositionVector(characterIdentifier int) geometry.PdfVector
	GetDisplacementVector(characterIdentifier int) geometry.PdfVector
	GetFontMatrix(characterIdentifier int) core.TransformationMatrix
	TryGetPath(characterCode int) ([]core.PdfSubpath, bool)
	TryGetNormalisedPath(characterCode int) ([]core.PdfSubpath, bool)
}

var _ cidFontRequirements = (*cidfonts.Type0CidFont)(nil)

// Type0Font defines glyphs using a CIDFont.
type Type0Font struct {
	ucs2CMap                *cmap.CMap
	isChineseJapaneseOrKorean bool
	boundingBoxCache        sync.Map
	useLenientParsing       bool
	ascent                  float64
	descent                 float64
	baseFont                *tokens.NameToken
	cidFont                 cidFontRequirements
	cMap                    *cmap.CMap
	toUnicode               *ToUnicodeCMap
	details                 *fonts.FontDetails
}

// NewType0Font creates a new Type0Font instance.
func NewType0Font(
	baseFont *tokens.NameToken,
	cidFont cidFontRequirements,
	cmapVal *cmap.CMap,
	toUnicodeCMap *cmap.CMap,
	ucs2CMap *cmap.CMap,
	useLenientParsing bool,
	isChineseJapaneseOrKorean bool,
) (*Type0Font, error) {
	if baseFont == nil {
		return nil, fmt.Errorf("baseFont must not be nil")
	}
	if cidFont == nil {
		return nil, fmt.Errorf("cidFont must not be nil")
	}
	if cmapVal == nil {
		return nil, fmt.Errorf("cmap must not be nil")
	}

	toUnicode := NewToUnicodeCMap(toUnicodeCMap)

	detailsVal := cidFont.Details()
	var details *fonts.FontDetails
	if detailsVal.Weight != 0 && baseFont.Data() != "" {
		details = detailsVal.WithName(baseFont.Data())
	} else {
		details = fonts.GetDefault(baseFont.Data())
	}

	t := &Type0Font{
		ucs2CMap:                ucs2CMap,
		isChineseJapaneseOrKorean: isChineseJapaneseOrKorean,
		useLenientParsing:       useLenientParsing,
		baseFont:                baseFont,
		cidFont:                 cidFont,
		cMap:                    cmapVal,
		toUnicode:               toUnicode,
		details:                 details,
	}

	t.ascent = t.computeAscent()
	t.descent = t.computeDescent()

	return t, nil
}

func (t *Type0Font) computeDescent() float64 {
	d := t.cidFont.GetDescent()
	if math.Abs(d) > 1e-9 {
		return t.GetFontMatrix().TransformY(d)
	}
	return -0.25
}

func (t *Type0Font) computeAscent() float64 {
	a := t.cidFont.GetAscent()
	if math.Abs(a) > 1e-9 {
		return t.GetFontMatrix().TransformY(a)
	}
	return 0.75
}

// Name returns the name token of the font (same as BaseFont).
func (t *Type0Font) Name() *tokens.NameToken {
	return t.baseFont
}

// BaseFont returns the base font name token.
func (t *Type0Font) BaseFont() *tokens.NameToken {
	return t.baseFont
}

// CidFont returns the underlying CID font if it is a cidfonts.CidFont.
// Returns nil if the internal representation cannot be cast.
func (t *Type0Font) CidFont() cidfonts.CidFont {
	if cf, ok := t.cidFont.(cidfonts.CidFont); ok {
		return cf
	}
	return nil
}

// CMap returns the character map for this font.
func (t *Type0Font) CMap() *cmap.CMap {
	return t.cMap
}

// ToUnicode returns the ToUnicode CMap wrapper.
func (t *Type0Font) ToUnicode() *ToUnicodeCMap {
	return t.toUnicode
}

// IsVertical reports whether the font uses vertical writing mode.
func (t *Type0Font) IsVertical() bool {
	return t.cMap.WritingMode() == cmap.Vertical
}

// Details returns the font details.
func (t *Type0Font) Details() *fonts.FontDetails {
	return t.details
}

// ReadCharacterCode reads a character code from the input bytes and returns it
// along with the number of bytes consumed (codeLength).
func (t *Type0Font) ReadCharacterCode(bytes core.InputBytes) (int, int) {
	if t.cMap == nil {
		return 0, 0
	}
	current := bytes.CurrentOffset()

	code, err := t.cMap.ReadCode(bytes, t.useLenientParsing)
	if err != nil {
		if !t.useLenientParsing {
			panic(err)
		}
		return 0, int(bytes.CurrentOffset() - current)
	}

	codeLength := int(bytes.CurrentOffset() - current)
	return code, codeLength
}

// TryGetUnicode attempts to get the Unicode string for the given character code.
// Returns true and the Unicode value if successful, false and an empty string otherwise.
func (t *Type0Font) TryGetUnicode(characterCode int) (string, bool) {
	if t.toUnicode == nil {
		return "", false
	}
	haveCMap := t.toUnicode.CanMapToUnicode()

	if !haveCMap && t.ucs2CMap != nil {
		cid := t.cMap.ConvertToCid(characterCode)
		if cid == 0 {
			return "", false
		}
		if value, ok := t.ucs2CMap.TryConvertToUnicode(cid); ok {
			return value, true
		}
		if value, ok := t.ucs2CMap.TryConvertToUnicode(characterCode); ok {
			return value, true
		}
	}

	if t.toUnicode.IsUsingIdentityAsUnicodeMap() {
		return string(rune(characterCode)), true
	}

	return t.toUnicode.TryGet(characterCode)
}

// GetBoundingBox returns the bounding box for the given character code.
func (t *Type0Font) GetBoundingBox(characterCode int) (*fonts.CharacterBoundingBox, error) {
	if cached, ok := t.boundingBoxCache.Load(characterCode); ok {
		return cached.(*fonts.CharacterBoundingBox), nil
	}

	characterIdentifier := t.cMap.ConvertToCid(characterCode)
	boundingBox, err := t.cidFont.GetBoundingBox(characterIdentifier)
	if err != nil {
		return nil, err
	}
	boundingBox = t.cidFont.GetFontMatrix(characterIdentifier).TransformRect(boundingBox)

	width := t.cidFont.GetWidthFromFont(characterIdentifier)
	advanceWidth := t.GetFontMatrix().TransformX(width)

	result := fonts.NewCharacterBoundingBox(boundingBox, advanceWidth)
	t.boundingBoxCache.Store(characterCode, result)

	return result, nil
}

// GetFontMatrix returns the transformation matrix for this font.
func (t *Type0Font) GetFontMatrix() core.TransformationMatrix {
	return t.cidFont.FontMatrix()
}

// GetDescent returns the descent of the font adjusted by the font matrix.
func (t *Type0Font) GetDescent() float64 {
	return t.descent
}

// GetAscent returns the ascent of the font adjusted by the font matrix.
func (t *Type0Font) GetAscent() float64 {
	return t.ascent
}

// GetPositionVector returns the position vector for vertical text layout.
func (t *Type0Font) GetPositionVector(characterCode int) core.PdfPoint {
	characterIdentifier := t.cMap.ConvertToCid(characterCode)
	vec := t.cidFont.GetPositionVector(characterIdentifier)
	scaled := vec.Scale(-1 / 1000.0)
	return scaled.ToPoint()
}

// GetDisplacementVector returns the displacement vector for vertical text layout.
func (t *Type0Font) GetDisplacementVector(characterCode int) core.PdfPoint {
	characterIdentifier := t.cMap.ConvertToCid(characterCode)
	vec := t.cidFont.GetDisplacementVector(characterIdentifier)
	scaled := vec.Scale(1 / 1000.0)
	return scaled.ToPoint()
}

// TryGetPath attempts to get the glyph path for the given character code.
// Returns true and the path if successful, false and nil otherwise.
func (t *Type0Font) TryGetPath(characterCode int) ([]*core.PdfSubpath, bool) {
	characterIdentifier := t.cMap.ConvertToCid(characterCode)
	path, ok := t.cidFont.TryGetPath(characterIdentifier)
	if !ok {
		return nil, false
	}
	return valuesToPointerSlice(path), true
}

// TryGetNormalisedPath attempts to get the normalised glyph path for the given
// character code. Returns true and the path if successful, false and nil otherwise.
func (t *Type0Font) TryGetNormalisedPath(characterCode int) ([]*core.PdfSubpath, bool) {
	characterIdentifier := t.cMap.ConvertToCid(characterCode)
	path, ok := t.cidFont.TryGetNormalisedPath(characterIdentifier)
	if !ok {
		return nil, false
	}
	return valuesToPointerSlice(path), true
}

func valuesToPointerSlice(values []core.PdfSubpath) []*core.PdfSubpath {
	result := make([]*core.PdfSubpath, len(values))
	for i := range values {
		result[i] = &values[i]
	}
	return result
}
