// Package handlers provides interfaces for PDF font parsing handlers.
package handlers

import (
	"fmt"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CidFontSystemInfo provides access to the character collection definition
// for a CID font (registry, ordering). Defined locally to avoid import cycles
// between handlers and pdf_fonts packages.
type CidFontSystemInfo interface {
	RegistryName() string
	OrderingName() string
	String() string
}

// CidFontLike represents the subset of cidfonts.CidFont needed by Type0FontHandler,
// defined locally to avoid import cycles between handlers and cidfonts packages.
type CidFontLike interface {
	SystemInfo() CidFontSystemInfo
}

// CombinedCidFont combines CidFontLike with all methods needed by composite.NewType0Font.
// This interface is satisfied by concrete CID font types (e.g., *cidfonts.Type0CidFont)
// but avoids importing cidfonts package directly which would create import cycles.
type CombinedCidFont interface {
	CidFontLike

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

// CidFontFactoryProvider creates CID fonts from dictionary tokens.
// Returns CombinedCidFont which carries all methods needed downstream for Type0 font construction.
type CidFontFactoryProvider interface {
	Generate(dictionary *tokens.DictionaryToken) CombinedCidFont
}

// ParsingOptionsProvider provides access to parsing configuration.
type ParsingOptionsProvider interface {
	UseLenientParsing() bool
	Logger() logging.Log
}

// Type0FontConstructor creates a Type 0 font from its components.
type Type0FontConstructor func(
	baseFont *tokens.NameToken,
	cidFont CombinedCidFont,
	cmapVal fonts.CMapProvider,
	toUnicodeCMap fonts.CMapProvider,
	ucs2CMap fonts.CMapProvider,
	useLenientParsing bool,
	isChineseJapaneseOrKorean bool,
) fonts.Font

// CMapCacheLookup is the callback type for looking up a predefined CMap by name
// from the global cache. The caller provides this function at construction time
// to avoid import cycles between handlers -> cmap -> pdf_fonts -> handlers.
type CMapCacheLookup func(name string) (fonts.CMapProvider, CidFontSystemInfo, bool)

// Type0FontHandler handles Type 0 (CID-keyed) font dictionaries.
type Type0FontHandler struct {
	cidFontFactory   CidFontFactoryProvider
	scanner          tokenization.PdfTokenScanner
	cmapLocalCache   CMapLocalCacheProvider
	cmapCacheLookup  CMapCacheLookup
	logger           logging.Log
	parsingOptions   ParsingOptionsProvider
	createType0Font  Type0FontConstructor
}

// NewType0FontHandler creates a new Type0FontHandler.
func NewType0FontHandler(
	cidFontFactory CidFontFactoryProvider,
	scanner tokenization.PdfTokenScanner,
	cmapLocalCache CMapLocalCacheProvider,
	cmapCacheLookup CMapCacheLookup,
	parsingOptions ParsingOptionsProvider,
	createType0Font Type0FontConstructor,
) *Type0FontHandler {
	return &Type0FontHandler{
		cidFontFactory:  cidFontFactory,
		scanner:         scanner,
		cmapLocalCache:  cmapLocalCache,
		cmapCacheLookup: cmapCacheLookup,
		logger:          parsingOptions.Logger(),
		parsingOptions:  parsingOptions,
		createType0Font: createType0Font,
	}
}

// Generate creates a Type 0 font from the given dictionary.
func (h *Type0FontHandler) Generate(dictionary *tokens.DictionaryToken) fonts.Font {
	baseFont := h.getNameOrDefault(dictionary, tokens.BaseFont)

	cMap, isCMapPredefined := h.readEncoding(dictionary)

	descendantObj, ok := h.tryGetFirstDescendant(dictionary)
	if !ok {
		h.logger.Error(fmt.Sprintf("No descendant font dictionary was declared for this Type 0 font. This dictionary should contain the CIDFont for the Type 0 font. %v", dictionary))
		return nil
	}

	descendantFontDictionary := h.resolveToDictionary(descendantObj)
	if descendantFontDictionary == nil {
		return nil
	}

	cidFont := h.parseDescendant(descendantFontDictionary)
	if cidFont == nil {
		return nil
	}

	ucs2CMap, isCJK := h.getUcs2CMap(dictionary, isCMapPredefined, cidFont)

	var toUnicodeCMap fonts.CMapProvider
	toUnicodeObj, hasToUnicode := dictionary.TryGet(tokens.ToUnicode)
	if hasToUnicode {
		streamToken, okStream := parts.TryGet[*tokens.StreamToken](toUnicodeObj, h.scanner)
		if okStream && streamToken != nil {
			c, found := h.cmapLocalCache.TryGetByStream(streamToken)
			if found {
				toUnicodeCMap = c
			}
		} else {
			nameToken, okName := parts.TryGet[*tokens.NameToken](toUnicodeObj, h.scanner)
			if okName && nameToken != nil {
				c, found := h.cmapLocalCache.TryGetByName(nameToken.Data())
				if found {
					toUnicodeCMap = c
				}
			} else {
				h.logger.Error(fmt.Sprintf("Invalid type of toUnicode CMap encountered for font named %v. Got: %v.", baseFont, toUnicodeObj))
			}
		}
	}

	font := h.createType0Font(baseFont, cidFont, cMap, toUnicodeCMap, ucs2CMap, h.parsingOptions.UseLenientParsing(), isCJK)
	return font
}

func (h *Type0FontHandler) getNameOrDefault(dictionary *tokens.DictionaryToken, name *tokens.NameToken) *tokens.NameToken {
	token, ok := parts.TryGet[*tokens.NameToken](getDictEntry(dictionary, name), h.scanner)
	if !ok || token == nil {
		return nil
	}
	return token
}

func (h *Type0FontHandler) resolveToDictionary(token tokens.Token) *tokens.DictionaryToken {
	if dict, ok := any(token).(*tokens.DictionaryToken); ok {
		return dict
	}

	ref, ok := any(token).(*tokens.IndirectReferenceToken)
	if !ok {
		return nil
	}

	parsed, err := parts.GetByToken[*tokens.DictionaryToken](ref, h.scanner)
	if err != nil || parsed == nil {
		return nil
	}

	return parsed
}

func (h *Type0FontHandler) tryGetFirstDescendant(dictionary *tokens.DictionaryToken) (tokens.Token, bool) {
	value, ok := dictionary.TryGet(tokens.DescendantFonts)
	if !ok {
		return nil, false
	}

	if ref, ok := any(value).(*tokens.IndirectReferenceToken); ok {
		return ref, true
	}

	array, isArray := any(value).(*tokens.ArrayToken)
	if !isArray || array.Length() == 0 {
		return nil, false
	}

	first := array.Get(0)
	if _, ok := any(first).(*tokens.IndirectReferenceToken); ok {
		return first, true
	}
	if _, ok := any(first).(*tokens.DictionaryToken); ok {
		return first, true
	}

	return nil, false
}

func (h *Type0FontHandler) parseDescendant(dictionary *tokens.DictionaryToken) CombinedCidFont {
	typeName := h.getNameOrDefault(dictionary, tokens.Type)
	if typeName == nil || !typeName.Equals(tokens.Font) {
		h.logger.Error(fmt.Sprintf("Expected 'Font' dictionary but found '%v'", typeName))
		return nil
	}

	result := h.cidFontFactory.Generate(dictionary)
	return result
}

func (h *Type0FontHandler) readEncoding(dictionary *tokens.DictionaryToken) (fonts.CMapProvider, bool) {
	isPredefined := false

	encodingName, okName := parts.TryGet[*tokens.NameToken](getDictEntry(dictionary, tokens.Encoding), h.scanner)
	if okName && encodingName != nil {
		cmapVal, found := h.cmapLocalCache.TryGetByName(encodingName.Data())
		if !found || cmapVal == nil {
			h.logger.Error(fmt.Sprintf("Missing CMap named %s.", encodingName.Data()))
			return nil, false
		}
		isPredefined = true
		return cmapVal, isPredefined
	}

	streamToken, okStream := parts.TryGet[*tokens.StreamToken](getDictEntry(dictionary, tokens.Encoding), h.scanner)
	if okStream && streamToken != nil {
		cmapVal, found := h.cmapLocalCache.TryGetByStream(streamToken)
		if !found || cmapVal == nil {
			h.logger.Error(fmt.Sprintf("Could not read CMap from stream in the dictionary: %v", dictionary))
			return nil, false
		}
		return cmapVal, isPredefined
	}

	h.logger.Error(fmt.Sprintf("Could not read the encoding, expected a name or a stream but it was not found in the dictionary: %v", dictionary))
	return nil, false
}

func (h *Type0FontHandler) getUcs2CMap(dictionary *tokens.DictionaryToken, isCMapPredefined bool, cidFont CombinedCidFont) (fonts.CMapProvider, bool) {
	if !isCMapPredefined {
		return nil, false
	}

	encodingName := h.getNameOrDefault(dictionary, tokens.Encoding)
	if encodingName == nil {
		return nil, false
	}

	isCJK := cidFontIsAdobeCJK(cidFont)

	isIdentityMap := encodingName.Equals(tokens.IdentityH) || encodingName.Equals(tokens.IdentityV)
	if isIdentityMap && !isCJK {
		return nil, false
	}

	if !isCJK {
		return nil, false
	}

	if h.cmapCacheLookup == nil {
		return nil, false
	}

	sysInfo := cidFont.SystemInfo()
	fullCmapName := sysInfo.String()
	var registry, ordering string

	nonUnicodeCMap, nonUnicodeSysInfo, found := h.cmapCacheLookup(fullCmapName)
	if found && nonUnicodeCMap != nil && nonUnicodeSysInfo != nil {
		registry = nonUnicodeSysInfo.RegistryName()
		ordering = nonUnicodeSysInfo.OrderingName()
	} else {
		registry = sysInfo.RegistryName()
		ordering = sysInfo.OrderingName()
	}

	unicodeCMapName := fmt.Sprintf("%s-%s-UCS2", registry, ordering)

	unicodeCMap, _, found := h.cmapCacheLookup(unicodeCMapName)
	if !found || unicodeCMap == nil {
		h.logger.Error(fmt.Sprintf("Could not locate CMap by name: %s.", unicodeCMapName))
		return nil, false
	}

	return unicodeCMap, true
}

// cidFontIsAdobeCJK checks whether the CID font uses an Adobe CJK character collection.
func cidFontIsAdobeCJK(cidFont CombinedCidFont) bool {
	sysInfo := cidFont.SystemInfo()
	if sysInfo == nil {
		return false
	}
	registry := sysInfo.RegistryName()
	ordering := sysInfo.OrderingName()
	return strings.EqualFold(registry, "Adobe") && (
		strings.EqualFold(ordering, "GB1") ||
			strings.EqualFold(ordering, "CNS1") ||
			strings.EqualFold(ordering, "Japan1") ||
			strings.EqualFold(ordering, "Korea1"))
}

var _ FontHandler = (*Type0FontHandler)(nil)
