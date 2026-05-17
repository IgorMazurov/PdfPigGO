// Package fonts provides writing font implementations for PDF generation.
package fonts

import (
	"fmt"
	"math"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	subsetting "github.com/uglytoad/pdfpig/go/fonts/truetype/subsetting"
	truetypeparser "github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/writer/interfaces"
)

// TrueTypeWritingFont represents a TrueType font for writing PDF documents.
type TrueTypeWritingFont struct {
	font             *truetypeparser.TrueTypeFont
	fontFileBytes    []byte
	mappingLock      sync.Mutex
	characterMapping map[rune]byte
	counter          int
}

// NewTrueTypeWritingFont creates a new TrueTypeWritingFont from a parsed TrueType font and its raw bytes.
func NewTrueTypeWritingFont(font *truetypeparser.TrueTypeFont, fontFileBytes []byte) *TrueTypeWritingFont {
	return &TrueTypeWritingFont{
		font:             font,
		fontFileBytes:    fontFileBytes,
		characterMapping: make(map[rune]byte),
		counter:          1,
	}
}

// HasWidths reports whether the font has explicit character widths.
func (f *TrueTypeWritingFont) HasWidths() bool {
	return true
}

// Name returns the PostScript name of the font.
func (f *TrueTypeWritingFont) Name() string {
	return f.font.Name()
}

// TryGetBoundingBox attempts to get the bounding box for a character.
func (f *TrueTypeWritingFont) TryGetBoundingBox(character rune) (*core.PdfRectangle, bool) {
	rect, ok := f.font.TryGetBoundingBox(int(character))
	if !ok {
		return nil, false
	}
	return &rect, true
}

// TryGetAdvanceWidth attempts to get the advance width for a character.
func (f *TrueTypeWritingFont) TryGetAdvanceWidth(character rune) (float64, bool) {
	width, ok := f.font.TryGetAdvanceWidth(int(character))
	return width, ok
}

// GetFontMatrix returns the transformation matrix for the font based on unitsPerEm.
func (f *TrueTypeWritingFont) GetFontMatrix() core.TransformationMatrix {
	unitsPerEm := float64(f.font.GetUnitsPerEm())
	return core.FromValues(1.0/unitsPerEm, 0, 0, 1.0/unitsPerEm, 0, 0)
}

// WriteFont writes the font dictionary to the stream and returns an indirect reference.
func (f *TrueTypeWritingFont) WriteFont(w interfaces.PdfStreamWriter, reservedIndirect *tokens.IndirectReferenceToken) *tokens.IndirectReferenceToken {
	f.mappingLock.Lock()
	keys := make([]rune, 0, len(f.characterMapping))
	for k := range f.characterMapping {
		keys = append(keys, k)
	}
	f.mappingLock.Unlock()

	newEncoding := subsetting.NewTrueTypeSubsetEncoding(keys)
	subsetBytes, err := subsetting.Subset(f.fontFileBytes, newEncoding)
	if err != nil {
		panic(fmt.Errorf("failed to subset font: %w", err))
	}

	embeddedFile := interfaces.CompressToStream(subsetBytes)
	fileRef := w.WriteToken(embeddedFile)

	tableRegister := f.font.TableRegister()
	baseFont := tokens.MustCreate(tableRegister.NameTable.GetPostscriptName())

	postscript := tableRegister.PostScriptTable
	hhead := tableRegister.HorizontalHeaderTable
	bbox := tableRegister.HeaderTable.Bounds()
	scaling := 1000.0 / float64(tableRegister.HeaderTable.UnitPerEm())

	descriptorDictionary := map[*tokens.NameToken]tokens.Token{
		tokens.Type:          tokens.FontDescriptor,
		tokens.FontName:      baseFont,
		tokens.Flags:         tokens.NewNumericTokenFromInt(int(fonts.Symbolic)),
		tokens.FontBbox:       getBoundingBox(bbox, scaling),
		tokens.ItalicAngle:   tokens.NewNumericToken(float64(postscript.ItalicAngle())),
		tokens.Ascent:        tokens.NewNumericToken(math.Round(float64(hhead.Ascent())*scaling*100) / 100),
		tokens.Descent:       tokens.NewNumericToken(math.Round(float64(hhead.Descent())*scaling*100) / 100),
		tokens.CapHeight:     tokens.NewNumericTokenFromInt(90),
		tokens.StemV:         tokens.NewNumericTokenFromInt(90),
		tokens.FontFile2:    fileRef,
	}

	os2 := tableRegister.Os2Table
	if os2.Tag() == "" {
		panic(fonts.NewInvalidFontFormatException("Embedding TrueType font requires OS/2 table."))
	}

	descriptorDictionary[tokens.StemV] = tokens.NewNumericToken(math.Round(bbox.Width*scaling*0.13*100) / 100)

	lastCharacter := 0
	widths := []tokens.Token{tokens.Zero}

	f.mappingLock.Lock()
	for kvpChar, kvpVal := range f.characterMapping {
		if int(kvpVal) > lastCharacter {
			lastCharacter = int(kvpVal)
		}

		glyphId := f.font.WindowsUnicodeCMap().CharacterCodeToGlyphIndex(int(kvpChar))
		width, _ := tableRegister.HorizontalMetricsTable.GetAdvanceWidth(glyphId)
		widthVal := math.Round(float64(width)*scaling*100) / 100
		widths = append(widths, tokens.NewNumericToken(widthVal))
	}
	f.mappingLock.Unlock()

	descriptorDict, _ := tokens.NewDictionary(descriptorDictionary)
	descriptor := w.WriteToken(descriptorDict)

	toUnicodeBuilder := &ToUnicodeCMapBuilder{}
	toUnicodeBytes, _ := toUnicodeBuilder.ConvertToCMapStream(f.characterMapping)
	toUnicodeStream := interfaces.CompressToStream(toUnicodeBytes)
	toUnicode := w.WriteToken(toUnicodeStream)

	fontDictionary := map[*tokens.NameToken]tokens.Token{
		tokens.Type:          tokens.Font,
		tokens.Subtype:       tokens.TrueType,
		tokens.BaseFont:      baseFont,
		tokens.FontDescriptor: descriptor,
		tokens.FirstChar:     tokens.Zero,
		tokens.LastChar:      tokens.NewNumericTokenFromInt(lastCharacter),
		tokens.Widths:        tokens.NewArrayToken(widths),
		tokens.ToUnicode:     toUnicode,
	}

	tokenDict, _ := tokens.NewDictionary(fontDictionary)

	if reservedIndirect != nil {
		return w.WriteTokenAt(tokenDict, reservedIndirect)
	}

	return w.WriteToken(tokenDict)
}

// GetValueForCharacter returns the byte value for a character in the font's encoding.
func (f *TrueTypeWritingFont) GetValueForCharacter(character rune) byte {
	f.mappingLock.Lock()
	defer f.mappingLock.Unlock()

	if result, ok := f.characterMapping[character]; ok {
		return result
	}

	if f.counter > 255 {
		panic(fmt.Errorf("cannot support more than 255 separate characters in a simple TrueType font"))
	}

	value := byte(f.counter)
	f.counter++
	f.characterMapping[character] = value
	return value
}

// getBoundingBox creates an ArrayToken from a scaled bounding box.
func getBoundingBox(boundingBox core.PdfRectangle, scaling float64) *tokens.ArrayToken {
	return tokens.NewArrayToken([]tokens.Token{
		tokens.NewNumericToken(math.Round(boundingBox.Left()*scaling*100) / 100),
		tokens.NewNumericToken(math.Round(boundingBox.Bottom()*scaling*100) / 100),
		tokens.NewNumericToken(math.Round(boundingBox.Right()*scaling*100) / 100),
		tokens.NewNumericToken(math.Round(boundingBox.Top()*scaling*100) / 100),
	})
}

var _ interfaces.WritingFont = (*TrueTypeWritingFont)(nil)
