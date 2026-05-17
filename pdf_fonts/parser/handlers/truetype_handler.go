// Package handlers provides interfaces for PDF font parsing handlers.
package handlers

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts"
	enc "github.com/uglytoad/pdfpig/go/fonts/encodings"
	"github.com/uglytoad/pdfpig/go/fonts/adobe_font_metrics"
	standard14fonts "github.com/uglytoad/pdfpig/go/fonts/standard14_fonts"
	"github.com/uglytoad/pdfpig/go/fonts/systemfonts"
	truetypeparser "github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/pdf_fonts/simple"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CMapLocalCacheProvider provides a local cache for CMap objects.
// Defined here to avoid import cycles between pdf_fonts and fonts/cmap packages.
type CMapLocalCacheProvider interface {
	TryGetByName(name string) (fonts.CMapProvider, bool)
	TryGetByStream(streamToken *tokens.StreamToken) (fonts.CMapProvider, bool)
}

// EncodingReader reads font encoding from a PDF dictionary.
// Defined here to avoid import cycles between pdf_fonts and fonts/parser packages.
type EncodingReader interface {
	Read(fontDictionary *tokens.DictionaryToken, descriptor *fonts.FontDescriptor, fontEncoding *enc.Encoding) *enc.Encoding
}

// TrueTypeFontHandler handles TrueType font dictionaries (Subtype = Type42).
type TrueTypeFontHandler struct {
	log              logging.Log
	scanner          tokenization.PdfTokenScanner
	filterProvider   filters.LookupFilterProvider
	encodingReader   EncodingReader
	systemFontFinder systemfonts.SystemFontFinder
	type1Handler     FontHandler
	cmapLocalCache   CMapLocalCacheProvider

	// getWidths is a function that resolves the /Widths array from the font dictionary.
	getWidths func(scanner tokenization.PdfTokenScanner, dictionary *tokens.DictionaryToken) ([]float64, error)

	// getFontDescriptor resolves and generates the FontDescriptor from the font dictionary.
	getFontDescriptor func(scanner tokenization.PdfTokenScanner, dictionary *tokens.DictionaryToken) (*fonts.FontDescriptor, error)

	// getName resolves the font name from either /BaseFont or descriptor's FontName.
	getName func(scanner tokenization.PdfTokenScanner, dictionary *tokens.DictionaryToken, descriptor *fonts.FontDescriptor) (*tokens.NameToken, error)
}

// NewTrueTypeFontHandler creates a new TrueTypeFontHandler.
func NewTrueTypeFontHandler(
	log logging.Log,
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
	encodingReader EncodingReader,
	cmapLocalCache CMapLocalCacheProvider,
	systemFontFinder systemfonts.SystemFontFinder,
	type1Handler FontHandler,
	getWidths func(tokenization.PdfTokenScanner, *tokens.DictionaryToken) ([]float64, error),
	getFontDescriptor func(tokenization.PdfTokenScanner, *tokens.DictionaryToken) (*fonts.FontDescriptor, error),
	getName func(tokenization.PdfTokenScanner, *tokens.DictionaryToken, *fonts.FontDescriptor) (*tokens.NameToken, error),
) *TrueTypeFontHandler {
	return &TrueTypeFontHandler{
		log:              log,
		scanner:          scanner,
		filterProvider:   filterProvider,
		encodingReader:   encodingReader,
		systemFontFinder: systemFontFinder,
		type1Handler:     type1Handler,
		cmapLocalCache:   cmapLocalCache,
		getWidths:        getWidths,
		getFontDescriptor: getFontDescriptor,
		getName:          getName,
	}
}

// Generate creates a font from the given dictionary.
func (h *TrueTypeFontHandler) Generate(dictionary *tokens.DictionaryToken) fonts.Font {
	firstCharDirect := h.tryGetFirstCharDirect(dictionary)

	hasRequiredEntries := firstCharDirect != nil &&
		dictionary.ContainsKey(tokens.FontDescriptor) &&
		dictionary.ContainsKey(tokens.Widths)

	if !hasRequiredEntries {
		return h.generateStandardOrFallback(dictionary, firstCharDirect)
	}

	firstCharacter := *firstCharDirect
	widths, err := h.getWidths(h.scanner, dictionary)
	if err != nil {
		h.log.Error(fmt.Sprintf("Could not get widths for TrueType font: %v", err))
		return nil
	}

	descriptor, err := h.getFontDescriptor(h.scanner, dictionary)
	if err != nil {
		h.log.Error(fmt.Sprintf("Could not get font descriptor for TrueType font: %v", err))
		return nil
	}

	font, actualHandler := h.parseTrueTypeFont(descriptor)
	if font == nil && actualHandler != nil {
		return actualHandler.Generate(dictionary)
	}

	name, err := h.getName(h.scanner, dictionary, descriptor)
	if err != nil {
		h.log.Error(fmt.Sprintf("Could not get name for TrueType font: %v", err))
		return nil
	}

	toUnicodeCMap := h.tryGetToUnicode(dictionary)
	encoding := h.encodingReader.Read(dictionary, descriptor, nil)

	encoding = h.resolveEncoding(encoding, font)

	result, err := simple.NewTrueTypeSimpleFont(name, descriptor, toUnicodeCMap, encoding, font, firstCharacter, widths)
	if err != nil {
		h.log.Error(fmt.Sprintf("Could not create TrueTypeSimpleFont: %v", err))
		return nil
	}

	return result
}

// tryGetFirstCharDirect attempts to get the /FirstChar entry directly without
// resolving indirect references, matching C# TryGetOptionalTokenDirect behavior.
func (h *TrueTypeFontHandler) tryGetFirstCharDirect(dictionary *tokens.DictionaryToken) *int {
	rawToken, ok := dictionary.TryGet(tokens.FirstChar)
	if !ok || rawToken == nil {
		return nil
	}

	token, ok := rawToken.(*tokens.NumericToken)
	if !ok || token == nil {
		return nil
	}
	v := token.IntVal()
	return &v
}

// tryGetFirstChar attempts to get the /FirstChar entry as a numeric token,
// resolving indirect references if necessary.
func (h *TrueTypeFontHandler) tryGetFirstChar(dictionary *tokens.DictionaryToken) *int {
	token, ok := parts.TryGet[*tokens.NumericToken](getDictEntry(dictionary, tokens.FirstChar), h.scanner)
	if !ok || token == nil {
		return nil
	}
	v := token.IntVal()
	return &v
}

// generateStandardOrFallback handles the case where required TrueType entries
// are missing, falling back to Standard 14 fonts or character-width inference.
func (h *TrueTypeFontHandler) generateStandardOrFallback(dictionary *tokens.DictionaryToken, firstCharPtr *int) fonts.Font {
	baseFont := h.tryGetBaseFont(dictionary)

	if baseFont == nil {
		h.log.Error(fmt.Sprintf("The provided TrueType font dictionary did not contain a /FirstChar or a /BaseFont entry: %v.", dictionary))
		return nil
	}

	standard14Font, ok := standard14fonts.GetAdobeFontMetrics(baseFont.Data())
	if !ok {
		lastCharToken, okLC := parts.TryGet[*tokens.NumericToken](getDictEntry(dictionary, tokens.LastChar), h.scanner)
		widthsArrayLoc, okWA := dictionary.TryGet(tokens.Widths)

		if okLC && okWA && lastCharToken != nil {
			widthsArr, err := parts.GetByToken[*tokens.ArrayToken](widthsArrayLoc, h.scanner)
			if err == nil && widthsArr != nil {
				v := (lastCharToken.IntVal() - widthsArr.Length()) + 1
				firstCharPtr = &v
			} else {
				h.log.Error(fmt.Sprintf("The provided TrueType font dictionary did not have a /FirstChar and did not match a Standard 14 font: %v.", dictionary))
				return nil
			}
		} else {
			h.log.Error(fmt.Sprintf("The provided TrueType font dictionary did not have a /FirstChar and did not match a Standard 14 font: %v.", dictionary))
			return nil
		}

		h.log.Error(fmt.Sprintf("TrueType font fallback path without Standard 14 metrics is not fully implemented: %v.", dictionary))
		return nil
	}

	// Second attempt to get FirstChar via indirect reference resolution.
	// In C#, TryGetOptionalTokenDirect failed initially but TryGet may succeed here.
	if resolvedFirstChar := h.tryGetFirstChar(dictionary); resolvedFirstChar != nil {
		firstCharPtr = resolvedFirstChar
	}

	fileSystemFont := h.systemFontFinder.GetTrueTypeFont(baseFont.Data())
	thisEncoding := h.encodingReader.Read(dictionary, nil, nil)

	if thisEncoding == nil {
		afmEnc, err := adobe_font_metrics.NewAdobeFontMetricsEncoding(&standard14Font)
		if err != nil {
			h.log.Error(fmt.Sprintf("Could not create AdobeFontMetricsEncoding: %v", err))
			return nil
		}
		thisEncoding = afmEnc.Encoding
	}

	widthsOverride := h.tryGetWidthsOverride(dictionary)

	result, err := simple.NewTrueTypeStandard14FallbackSimpleFont(
		baseFont,
		&standard14Font,
		thisEncoding,
		fileSystemFont,
		simple.NewMetricOverrides(firstCharPtr, widthsOverride),
	)
	if err != nil {
		h.log.Error(fmt.Sprintf("Could not create TrueTypeStandard14FallbackSimpleFont: %v", err))
		return nil
	}

	return result
}

// tryGetBaseFont attempts to get the /BaseFont entry as a NameToken.
func (h *TrueTypeFontHandler) tryGetBaseFont(dictionary *tokens.DictionaryToken) *tokens.NameToken {
	baseFont, ok := parts.TryGet[*tokens.NameToken](getDictEntry(dictionary, tokens.BaseFont), h.scanner)
	if !ok || baseFont == nil {
		return nil
	}
	return baseFont
}

// tryGetWidthsOverride attempts to get the /Widths array and convert it to []float64.
func (h *TrueTypeFontHandler) tryGetWidthsOverride(dictionary *tokens.DictionaryToken) []float64 {
	widthsArray, ok := dictionary.TryGet(tokens.Widths)
	if !ok {
		return nil
	}

	arr, err := parts.GetByToken[*tokens.ArrayToken](widthsArray, h.scanner)
	if err != nil || arr == nil {
		return nil
	}

	result := make([]float64, 0, arr.Length())
	for i := 0; i < arr.Length(); i++ {
		elem := arr.Get(i)
		if num, ok := elem.(*tokens.NumericToken); ok {
			result = append(result, num.DoubleVal())
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

// parseTrueTypeFont attempts to parse the TrueType font from the descriptor.
// Returns (font, fallbackHandler) where fallbackHandler is non-nil if the
// stream contains a Type 1C font that should be handled by type1Handler.
func (h *TrueTypeFontHandler) parseTrueTypeFont(descriptor *fonts.FontDescriptor) (*truetypeparser.TrueTypeFont, FontHandler) {
	if descriptor == nil || descriptor.FontFile == nil {
		if descriptor != nil && descriptor.FontName != nil {
			ttf := h.systemFontFinder.GetTrueTypeFont(descriptor.FontName.Data())
			return ttf, nil
		}
		h.log.Error("Failed finding system font by name: missing font name in descriptor.")
		return nil, nil
	}

	streamToken, err := parts.GetByToken[*tokens.StreamToken](descriptor.FontFile.ObjectKey, h.scanner)
	if err != nil || streamToken == nil {
		h.log.Error("Could not parse the TrueType font.")
		return nil, nil
	}

	if descriptor.FontFile.FileType == fonts.FromSubtype {
		subtypeName, ok := parts.TryGet[*tokens.NameToken](getDictEntry(streamToken.StreamDictionary, tokens.Subtype), h.scanner)
		if ok && subtypeName != nil {
			if subtypeName.Equals(tokens.Type1C) {
				return nil, h.type1Handler
			}
			if !subtypeName.Equals(tokens.OpenType) {
				h.log.Error(fmt.Sprintf("Expected a TrueType font in the TrueType font descriptor, instead it was %v.", descriptor.FontFile.FileType))
				return nil, nil
			}
		} else {
			h.log.Error(fmt.Sprintf("Expected a TrueType font in the TrueType font descriptor, instead it was %v.", descriptor.FontFile.FileType))
			return nil, nil
		}
	}

	fontFileBytes := h.decodeStream(streamToken)
	if fontFileBytes == nil {
		return nil, nil
	}

	dataBytes := truetypeparser.NewTrueTypeDataBytes(fontFileBytes)
	font := truetypeparser.Parse(dataBytes)
	return font, nil
}

// decodeStream decodes the stream data using the filter provider and scanner.
func (h *TrueTypeFontHandler) decodeStream(streamToken *tokens.StreamToken) []byte {
	fl, err := h.filterProvider.GetFiltersWithScanner(streamToken.StreamDictionary, h.scanner)
	if err != nil || len(fl) == 0 {
		return streamToken.Data()
	}

	transform := streamToken.Data()
	for _, filter := range fl {
		result, err := filter.Decode(transform, streamToken.StreamDictionary, h.filterProvider, 0)
		if err != nil {
			h.log.Error(fmt.Sprintf("Could not decode TrueType font file stream: %v", err))
			return nil
		}
		transform = result
	}

	return transform
}

// tryGetToUnicode attempts to get and parse the ToUnicode CMap from the dictionary.
func (h *TrueTypeFontHandler) tryGetToUnicode(dictionary *tokens.DictionaryToken) fonts.CMapProvider {
	toUnicodeObj, ok := dictionary.TryGet(tokens.ToUnicode)
	if !ok {
		return nil
	}

	streamToken, err := parts.GetByToken[*tokens.StreamToken](toUnicodeObj, h.scanner)
	if err != nil || streamToken == nil {
		h.log.Error("Failed to resolve ToUnicode CMap stream for a TrueType font.")
		return nil
	}

	cmap, found := h.cmapLocalCache.TryGetByStream(streamToken)
	if !found {
		h.log.Error("Failed to decode ToUnicode CMap for a TrueType font in file.")
		return nil
	}

	return cmap
}

// resolveEncoding synthesizes an encoding from the font's CMap and PostScript tables
// if no explicit encoding is found.
func (h *TrueTypeFontHandler) resolveEncoding(encoding *enc.Encoding, font *truetypeparser.TrueTypeFont) *enc.Encoding {
	if encoding != nil || font == nil {
		return encoding
	}

	tableReg := font.TableRegister()
	if tableReg.CMapTable == nil || tableReg.PostScriptTable.GlyphNames() == nil {
		return nil
	}

	postscript := tableReg.PostScriptTable
	fakeEncoding := make(map[int]string)

	for i := 0; i < 256; i++ {
		index, ok := tableReg.CMapTable.TryGetGlyphIndex(i)
		if !ok {
			continue
		}

		var glyphName string
		glyphNames := postscript.GlyphNames()
		if index >= 0 && index < len(glyphNames) {
			glyphName = glyphNames[index]
		} else {
			glyphName = fmt.Sprintf("%d", index)
		}

		fakeEncoding[i] = glyphName
	}

	builtInEnc := enc.NewBuiltInEncoding(fakeEncoding)

	return builtInEnc.Encoding
}

func getDictEntry(dict *tokens.DictionaryToken, name *tokens.NameToken) tokens.Token {
	token, _ := dict.TryGet(name)
	return token
}

var _ FontHandler = (*TrueTypeFontHandler)(nil)
