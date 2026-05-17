package parser

import (
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/type1"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Type1Font represents a parsed Adobe Type 1 font program.
type Type1Font struct {
	Name              string
	Encoding          map[int]string
	FontMatrix        core.TransformationMatrix
	BoundingBox       core.PdfRectangle
	PrivateDictionary *type1.Type1PrivateDictionary
	CharStrings       *charstrings.Type1CharStrings
}

// NewType1Font creates a new Type1Font instance, transforming the raw font matrix
// array token into a TransformationMatrix value.
func NewType1Font(
	name string,
	encoding map[int]string,
	fontMatrix *tokens.ArrayToken,
	boundingBox core.PdfRectangle,
	privateDict *type1.Type1PrivateDictionary,
	charStrings *charstrings.Type1CharStrings,
) *Type1Font {
	return &Type1Font{
		Name:              name,
		Encoding:          encoding,
		FontMatrix:        getFontTransformationMatrix(fontMatrix),
		BoundingBox:       boundingBox,
		PrivateDictionary: privateDict,
		CharStrings:       charStrings,
	}
}

// GetCharacterBoundingBox returns the bounding box for the character with the given name,
// or nil if the character is .notdef or has no path data.
func (f *Type1Font) GetCharacterBoundingBox(characterName string) *core.PdfRectangle {
	if strings.EqualFold(characterName, fonts.NotDefined) {
		return nil
	}

	glyph := f.GetCharacterPath(characterName)
	return core.GetBoundingRectangleForPath(glyph)
}

// ContainsNamedCharacter reports whether the font contains a character with the given name.
func (f *Type1Font) ContainsNamedCharacter(name string) bool {
	_, ok := f.CharStrings.CharStrings()[name]
	return ok
}

// GetCharacterPath returns the subpath geometry for the character with the given name,
// or nil if the character is .notdef or has no charstring definition.
func (f *Type1Font) GetCharacterPath(characterName string) []*core.PdfSubpath {
	if strings.EqualFold(characterName, fonts.NotDefined) {
		return nil
	}

	glyph, ok := f.CharStrings.TryGenerate(characterName)
	if !ok {
		return nil
	}

	return glyph
}

func getFontTransformationMatrix(array *tokens.ArrayToken) core.TransformationMatrix {
	if array == nil || len(array.Data()) != 6 {
		return core.FromValues(0.001, 0, 0, 0.001, 0, 0)
	}

	a := array.Data()[0].(*tokens.NumericToken).Data()
	b := array.Data()[1].(*tokens.NumericToken).Data()
	c := array.Data()[2].(*tokens.NumericToken).Data()
	d := array.Data()[3].(*tokens.NumericToken).Data()
	e := array.Data()[4].(*tokens.NumericToken).Data()
	f := array.Data()[5].(*tokens.NumericToken).Data()

	return core.FromValues(a, b, c, d, e, f)
}
