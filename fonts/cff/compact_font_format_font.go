package cff

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/cfftypes"
	cffcharset "github.com/uglytoad/pdfpig/go/fonts/cff_charset"
)

// Type2CharStringsProvider is an alias for the shared interface.
type Type2CharStringsProvider = cfftypes.Type2CharStringsProvider

// Type2GlyphResult is an alias for the shared interface.
type Type2GlyphResult = cfftypes.Type2GlyphResult

// CompactFontFormatFont represents a Compact Font Format (CFF) font with its
// dictionaries, charset, charstrings, and encoding.
// encodingSource provides glyph name lookup for a character code.
type encodingSource interface {
	GetName(code int) string
}

type CompactFontFormatFont struct {
	topDictionary     *CompactFontFormatTopLevelDictionary
	privateDictionary CompactFontFormatPrivateDictionary
	charset           cffcharset.CompactFontFormatCharset
	type1CharStrings  any
	type2CharStrings  Type2CharStringsProvider
	encoding          encodingSource
}

// NewCompactFontFormatFont creates a new CompactFontFormatFont with the given
// dictionaries, charset, charstrings, and encoding. The type1CharStrings parameter
// should be nil for Type 2 fonts (the common case). Pass a non-nil value only if
// using unsupported Type 1 CharStrings.
func NewCompactFontFormatFont(
	topDictionary *CompactFontFormatTopLevelDictionary,
	privateDictionary CompactFontFormatPrivateDictionary,
	charset cffcharset.CompactFontFormatCharset,
	type1CharStrings any,
	type2CharStrings Type2CharStringsProvider,
	encoding encodingSource,
) *CompactFontFormatFont {
	return &CompactFontFormatFont{
		topDictionary:     topDictionary,
		privateDictionary: privateDictionary,
		charset:           charset,
		type1CharStrings:  type1CharStrings,
		type2CharStrings:  type2CharStrings,
		encoding:          encoding,
	}
}

// FontMatrix returns the font matrix for this font. If not set in the top dictionary,
// it defaults to a scaling matrix with 0.001 on both axes.
func (f *CompactFontFormatFont) FontMatrix() core.TransformationMatrix {
	if f.topDictionary != nil && f.topDictionary.FontMatrix != nil {
		return *f.topDictionary.FontMatrix
	}
	return core.FromValues(0.001, 0, 0, 0.001, 0, 0)
}

// Weight returns the weight string from the top dictionary, or an empty string if not set.
func (f *CompactFontFormatFont) Weight() string {
	if f.topDictionary == nil {
		return ""
	}
	return f.topDictionary.Weight
}

// ItalicAngle returns the italic angle from the top dictionary, or 0 if not set.
func (f *CompactFontFormatFont) ItalicAngle() float64 {
	if f.topDictionary == nil {
		return 0
	}
	return f.topDictionary.ItalicAngle
}

// Encoding returns the encoding source for this font.
func (f *CompactFontFormatFont) Encoding() encodingSource {
	return f.encoding
}

// GetCharacterName returns the character name for the given character code.
// If isCid is true or the charset is a CID charset, it looks up by string ID.
// Otherwise it uses the encoding or falls back to the Adobe glyph list.
func (f *CompactFontFormatFont) GetCharacterName(characterCode int, isCid bool) string {
	if f.encoding != nil {
		return f.encoding.GetName(characterCode)
	}

	if f.charset.IsCidCharset() || isCid {
		return f.charset.GetNameByStringId(characterCode)
	}

	characterName := fonts.AdobeGlyphListUnicodeCodePointToName(characterCode)

	if characterName == fonts.AdobeGlyphListNotDefined {
		return f.charset.GetNameByStringId(characterCode)
	}

	return characterName
}

// GetCharacterBoundingBox returns the bounding box for the character with the given name,
// or nil if it cannot be computed. Type 1 CharStrings are currently unsupported.
func (f *CompactFontFormatFont) GetCharacterBoundingBox(characterName string) *core.PdfRectangle {
	defaultWidthX := f.getDefaultWidthX(characterName)
	nominalWidthX := f.getNominalWidthX(characterName)

	if f.type1CharStrings != nil {
		panic("Type 1 CharStrings in a CFF font are currently unsupported")
	}

	if f.type2CharStrings == nil {
		return nil
	}

	glyph, err := f.type2CharStrings.Generate(characterName, defaultWidthX, nominalWidthX)
	if err != nil {
		return nil
	}

	rectangle := core.GetBoundingRectangleForPath(glyph.Path())
	if rectangle != nil {
		return rectangle
	}

	defaultBoundingBox := f.topDictionary.FontBoundingBox
	width := float64(0)
	if glyph.Width() != nil {
		width = *glyph.Width()
	}
	rect := core.NewPdfRectangleFloat(0, 0, width, defaultBoundingBox.Height)
	return &rect
}

// TryGetPath attempts to get the path for the character with the given name.
// Returns true and the path if successful, false and nil otherwise.
// Type 1 CharStrings are currently unsupported.
func (f *CompactFontFormatFont) TryGetPath(characterName string) ([]*core.PdfSubpath, bool) {
	defaultWidthX := f.getDefaultWidthX(characterName)
	nominalWidthX := f.getNominalWidthX(characterName)

	if f.type1CharStrings != nil {
		panic("Type 1 CharStrings in a CFF font are currently unsupported")
	}

	if f.type2CharStrings == nil {
		return nil, false
	}

	glyph, err := f.type2CharStrings.Generate(characterName, defaultWidthX, nominalWidthX)
	if err != nil {
		return nil, false
	}

	return glyph.Path(), true
}

// GetCharacterPath returns the path for the character with the given name,
// or nil if it cannot be generated. Type 1 CharStrings are currently unsupported.
func (f *CompactFontFormatFont) GetCharacterPath(characterName string) []*core.PdfSubpath {
	defaultWidthX := f.getDefaultWidthX(characterName)
	nominalWidthX := f.getNominalWidthX(characterName)

	if f.type1CharStrings != nil {
		panic("Type 1 CharStrings in a CFF font are currently unsupported")
	}

	if f.type2CharStrings == nil {
		return nil
	}

	glyph, err := f.type2CharStrings.Generate(characterName, defaultWidthX, nominalWidthX)
	if err != nil {
		return nil
	}

	return glyph.Path()
}

// GetFontMatrix returns the font matrix for the corresponding character name,
// or nil if not available. Subclasses may override this behavior.
func (f *CompactFontFormatFont) GetFontMatrix(_ string) *core.TransformationMatrix {
	if f.topDictionary == nil || f.topDictionary.FontMatrix == nil {
		return nil
	}
	m := *f.topDictionary.FontMatrix
	return &m
}

// CharacterNames returns the names of all characters that have charstring definitions,
// or nil if this font does not use Type 2 CharStrings with enumeration support.
func (f *CompactFontFormatFont) CharacterNames() []string {
	type namedEnumerator interface {
		CharacterNames() []string
	}
	if e, ok := f.type2CharStrings.(namedEnumerator); ok {
		return e.CharacterNames()
	}
	return nil
}

// getDefaultWidthX returns the default width X from the private dictionary.
func (f *CompactFontFormatFont) getDefaultWidthX(_ string) float64 {
	return f.privateDictionary.DefaultWidthX
}

// getNominalWidthX returns the nominal width X from the private dictionary.
func (f *CompactFontFormatFont) getNominalWidthX(_ string) float64 {
	return f.privateDictionary.NominalWidthX
}

// CompactFontFormatCidFont represents a CFF font that uses CID (Character Identifier)
// encoding with multiple font dictionaries and private dictionaries.
type CompactFontFormatCidFont struct {
	*CompactFontFormatFont
	fontDictionaries    []*CompactFontFormatTopLevelDictionary
	privateDictionaries []CompactFontFormatPrivateDictionary
	fdSelect            FdSelect
}

// NewCompactFontFormatCidFont creates a new CompactFontFormatCidFont with the given
// dictionaries, charset, charstrings, font dictionary select, and collections.
func NewCompactFontFormatCidFont(
	topDictionary *CompactFontFormatTopLevelDictionary,
	privateDictionary CompactFontFormatPrivateDictionary,
	charset cffcharset.CompactFontFormatCharset,
	type1CharStrings any,
	type2CharStrings Type2CharStringsProvider,
	fontDictionaries []*CompactFontFormatTopLevelDictionary,
	privateDictionaries []CompactFontFormatPrivateDictionary,
	fdSelect FdSelect,
) *CompactFontFormatCidFont {
	return &CompactFontFormatCidFont{
		CompactFontFormatFont: NewCompactFontFormatFont(topDictionary, privateDictionary, charset, type1CharStrings, type2CharStrings, nil),
		fontDictionaries:      fontDictionaries,
		privateDictionaries:   privateDictionaries,
		fdSelect:              fdSelect,
	}
}

// FontDictionaries returns the list of top-level font dictionaries for this CID font.
func (f *CompactFontFormatCidFont) FontDictionaries() []*CompactFontFormatTopLevelDictionary {
	return f.fontDictionaries
}

// PrivateDictionaries returns the list of private dictionaries for this CID font.
func (f *CompactFontFormatCidFont) PrivateDictionaries() []CompactFontFormatPrivateDictionary {
	return f.privateDictionaries
}

// getDefaultWidthX overrides the base implementation to look up the correct private
// dictionary based on the character's font dictionary index. Returns 1000 if not found.
func (f *CompactFontFormatCidFont) getDefaultWidthX(characterName string) float64 {
	dict := f.tryGetPrivateDictionaryForCharacter(characterName)
	if dict == nil {
		return 1000
	}
	return dict.DefaultWidthX
}

// getNominalWidthX overrides the base implementation to look up the correct private
// dictionary based on the character's font dictionary index. Returns 0 if not found.
func (f *CompactFontFormatCidFont) getNominalWidthX(characterName string) float64 {
	dict := f.tryGetPrivateDictionaryForCharacter(characterName)
	if dict == nil {
		return 0
	}
	return dict.NominalWidthX
}

// GetFontMatrix overrides the base implementation to compose the top-level font matrix
// with the per-character font dictionary matrix when both are available.
func (f *CompactFontFormatCidFont) GetFontMatrix(characterName string) *core.TransformationMatrix {
	fontDict := f.tryGetFontDictionaryForCharacter(characterName)

	topHasMatrix := f.topDictionary != nil && f.topDictionary.FontMatrix != nil
	dictHasMatrix := fontDict != nil && fontDict.FontMatrix != nil

	if topHasMatrix && dictHasMatrix {
		multiplied := f.topDictionary.FontMatrix.Multiply(*fontDict.FontMatrix)
		return &multiplied
	}

	if topHasMatrix {
		m := *f.topDictionary.FontMatrix
		return &m
	}

	if dictHasMatrix {
		m := *fontDict.FontMatrix
		return &m
	}

	return nil
}

// tryGetPrivateDictionaryForCharacter looks up the private dictionary for the given
// character name via its glyph ID and font dictionary select.
func (f *CompactFontFormatCidFont) tryGetPrivateDictionaryForCharacter(characterName string) *CompactFontFormatPrivateDictionary {
	glyphId := f.charset.GetGlyphIdByName(characterName)
	fd := f.fdSelect.GetFontDictionaryIndex(glyphId)
	if fd == -1 || fd >= len(f.privateDictionaries) {
		return nil
	}
	dict := f.privateDictionaries[fd]
	return &dict
}

// tryGetFontDictionaryForCharacter looks up the top-level font dictionary for the given
// character name via its glyph ID and font dictionary select.
func (f *CompactFontFormatCidFont) tryGetFontDictionaryForCharacter(characterName string) *CompactFontFormatTopLevelDictionary {
	glyphId := f.charset.GetGlyphIdByName(characterName)
	fd := f.fdSelect.GetFontDictionaryIndex(glyphId)
	if fd == -1 || fd >= len(f.fontDictionaries) {
		return nil
	}
	return f.fontDictionaries[fd]
}
