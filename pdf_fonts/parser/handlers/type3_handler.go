// Package handlers provides interfaces for PDF font parsing handlers.
package handlers

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	pdfparser "github.com/uglytoad/pdfpig/go/pdf_fonts/parser"
	"github.com/uglytoad/pdfpig/go/pdf_fonts/simple"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Type3FontHandler handles Type 3 font dictionaries.
type Type3FontHandler struct {
	scanner        tokenization.PdfTokenScanner
	encodingReader EncodingReader
	cmapLocalCache CMapLocalCacheProvider
	logger         logging.Log
}

// NewType3FontHandler creates a new Type3FontHandler.
func NewType3FontHandler(
	scanner tokenization.PdfTokenScanner,
	encodingReader EncodingReader,
	cmapLocalCache CMapLocalCacheProvider,
) *Type3FontHandler {
	return &Type3FontHandler{
		scanner:        scanner,
		encodingReader: encodingReader,
		cmapLocalCache: cmapLocalCache,
		logger:         logging.NoopLog,
	}
}

// Generate creates a Type 3 font from the given dictionary.
func (h *Type3FontHandler) Generate(dictionary *tokens.DictionaryToken) fonts.Font {
	boundingBox := h.getBoundingBox(dictionary)
	if boundingBox == nil {
		return nil
	}

	fontMatrix, err := h.getFontMatrix(dictionary)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Could not get font matrix for Type 3 font: %v", err))
		return nil
	}

	bb := *boundingBox
	if bb.Left() == 0 && bb.Bottom() == 0 && bb.Height == 0 && bb.Width == 0 &&
		fontMatrix.A != 0 && fontMatrix.D != 0 {
		bb = core.NewPdfRectangleFloat(0, 0, 1/fontMatrix.A, 1/fontMatrix.D)
	}

	firstCharacter, err := pdfparser.GetFirstCharacter(dictionary)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Could not get first character for Type 3 font: %v", err))
		return nil
	}

	lastCharacter, err := pdfparser.GetLastCharacter(dictionary)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Could not get last character for Type 3 font: %v", err))
		return nil
	}

	widths, err := pdfparser.GetWidths(h.scanner, dictionary)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Could not get widths for Type 3 font: %v", err))
		return nil
	}

	encoding := h.encodingReader.Read(dictionary, nil, nil)
	if encoding == nil {
		h.logger.Error(fmt.Sprintf("No encoding found for Type 3 font: %v.", dictionary))
		return nil
	}

	toUnicodeCMap := h.tryGetToUnicode(dictionary)

	name := h.getFontName(dictionary)

	result, err := simple.NewType3Font(name, bb, fontMatrix, encoding, firstCharacter, lastCharacter, widths, toUnicodeCMap)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Could not create Type 3 font: %v", err))
		return nil
	}

	return result
}

// getFontName resolves the font name from the /Name entry or falls back to "Type3".
func (h *Type3FontHandler) getFontName(dictionary *tokens.DictionaryToken) *tokens.NameToken {
	fontName, ok := parts.TryGet[*tokens.NameToken](getDictEntry(dictionary, tokens.Name), h.scanner)
	if ok && fontName != nil {
		return fontName
	}

	return tokens.Type3
}

// getFontMatrix resolves the /FontMatrix entry as a TransformationMatrix.
func (h *Type3FontHandler) getFontMatrix(dictionary *tokens.DictionaryToken) (core.TransformationMatrix, error) {
	matrixObject, ok := dictionary.TryGet(tokens.FontMatrix)
	if !ok {
		return core.Identity, fonts.NewInvalidFontFormatException(fmt.Sprintf("No font matrix found: %v.", dictionary))
	}

	matrixArray, err := parts.GetByToken[*tokens.ArrayToken](matrixObject, h.scanner)
	if err != nil {
		return core.Identity, fonts.NewInvalidFontFormatExceptionWithInner("Could not resolve font matrix array", err)
	}

	if matrixArray.Length() < 6 {
		return core.Identity, fonts.NewInvalidFontFormatException(fmt.Sprintf("Font matrix array has fewer than 6 elements: %v.", dictionary))
	}

	vals := make([]float64, 6)
	for i := 0; i < 6; i++ {
		num, ok := matrixArray.Get(i).(*tokens.NumericToken)
		if !ok {
			return core.Identity, fonts.NewInvalidFontFormatException(fmt.Sprintf("Non-numeric value in font matrix at index %d.", i))
		}
		vals[i] = num.DoubleVal()
	}

	return core.FromValues(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5]), nil
}

// getBoundingBox resolves the /FontBBox entry as a PdfRectangle.
func (h *Type3FontHandler) getBoundingBox(dictionary *tokens.DictionaryToken) *core.PdfRectangle {
	bboxObject, ok := dictionary.TryGet(tokens.FontBbox)
	if !ok {
		h.logger.Error(fmt.Sprintf("Type 3 font was invalid. No Font Bounding Box: %v.", dictionary))
		return nil
	}

	bboxArray, err := parts.GetByToken[*tokens.ArrayToken](bboxObject, h.scanner)
	if err != nil || bboxArray == nil {
		return &core.PdfRectangle{}
	}

	if bboxArray.Length() < 4 {
		return &core.PdfRectangle{}
	}

	vals := make([]float64, 4)
	for i := 0; i < 4; i++ {
		num, ok := bboxArray.Get(i).(*tokens.NumericToken)
		if !ok {
			return &core.PdfRectangle{}
		}
		vals[i] = num.DoubleVal()
	}

	rect := core.NewPdfRectangleFloat(vals[0], vals[1], vals[2], vals[3])
	return &rect
}

// tryGetToUnicode attempts to get and parse the ToUnicode CMap from the dictionary.
func (h *Type3FontHandler) tryGetToUnicode(dictionary *tokens.DictionaryToken) fonts.CMapProvider {
	toUnicodeObj, ok := dictionary.TryGet(tokens.ToUnicode)
	if !ok {
		return nil
	}

	streamToken, err := parts.GetByToken[*tokens.StreamToken](toUnicodeObj, h.scanner)
	if err != nil || streamToken == nil {
		h.logger.Error("Failed to resolve ToUnicode CMap stream for a Type 3 font.")
		return nil
	}

	cmap, found := h.cmapLocalCache.TryGetByStream(streamToken)
	if !found {
		h.logger.Error("Failed to decode ToUnicode CMap for a Type 3 font in file.")
		return nil
	}

	return cmap
}

var _ FontHandler = (*Type3FontHandler)(nil)
