// Package handlers provides interfaces for PDF font parsing handlers.
package handlers

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	cff "github.com/uglytoad/pdfpig/go/fonts/cff"
	"github.com/uglytoad/pdfpig/go/fonts/cff/charstrings"
	cffcharset "github.com/uglytoad/pdfpig/go/fonts/cff_charset"
	enc "github.com/uglytoad/pdfpig/go/fonts/encodings"
	"github.com/uglytoad/pdfpig/go/fonts"
	standard14fonts "github.com/uglytoad/pdfpig/go/fonts/standard14_fonts"
	type1parser "github.com/uglytoad/pdfpig/go/fonts/type1/parser"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	parser "github.com/uglytoad/pdfpig/go/pdf_fonts/parser"
	"github.com/uglytoad/pdfpig/go/pdf_fonts/simple"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Type1FontHandler handles Type 1 font dictionaries.
type Type1FontHandler struct {
	scanner          tokenization.PdfTokenScanner
	filterProvider   filters.LookupFilterProvider
	encodingReader   EncodingReader
	cmapLocalCache   CMapLocalCacheProvider
	logger           logging.Log
	stackDepthGuard  *core.StackDepthGuard
	isLenientParsing bool
}

// NewType1FontHandler creates a new Type1FontHandler.
func NewType1FontHandler(
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
	encodingReader EncodingReader,
	cmapLocalCache CMapLocalCacheProvider,
	stackDepthGuard *core.StackDepthGuard,
	isLenientParsing bool,
	logger logging.Log,
) *Type1FontHandler {
	return &Type1FontHandler{
		scanner:          scanner,
		filterProvider:   filterProvider,
		encodingReader:   encodingReader,
		cmapLocalCache:   cmapLocalCache,
		stackDepthGuard:  stackDepthGuard,
		isLenientParsing: isLenientParsing,
		logger:           logger,
	}
}

// Generate creates a font from the given dictionary.
func (h *Type1FontHandler) Generate(dictionary *tokens.DictionaryToken) fonts.Font {
	usingStandard14Only := !dictionary.ContainsKey(tokens.FirstChar) || !dictionary.ContainsKey(tokens.Widths)

	if usingStandard14Only {
		baseFont, ok := parts.TryGet[*tokens.NameToken](getDictEntry(dictionary, tokens.BaseFont), h.scanner)
		if !ok || baseFont == nil {
			h.logger.Error(fmt.Sprintf("The Type 1 font did not contain a first character entry but also did not reference a standard 14 font: %v", dictionary))
			return nil
		}

		metrics, found := standard14fonts.GetAdobeFontMetrics(baseFont.Data())
		if found {
			overrideEncoding := h.encodingReader.Read(dictionary, nil, nil)
			font, err := simple.NewType1Standard14Font(metrics, overrideEncoding)
			if err != nil {
				h.logger.Error(fmt.Sprintf("Could not create Type1Standard14Font: %v", err))
				return nil
			}
			return font
		}
	}

	firstCharacter := 0
	lastCharacter := 0
	var widths []float64

	if !usingStandard14Only {
		firstCharVal, err := parser.GetFirstCharacter(dictionary)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Could not get first character for Type 1 font: %v", err))
			return nil
		}
		lastCharVal, err := parser.GetLastCharacter(dictionary)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Could not get last character for Type 1 font: %v", err))
			return nil
		}
		widthsArr, err := parser.GetWidths(h.scanner, dictionary)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Could not get widths for Type 1 font: %v", err))
			return nil
		}
		firstCharacter = firstCharVal
		lastCharacter = lastCharVal
		widths = widthsArr
	}

	if !dictionary.ContainsKey(tokens.FontDescriptor) {
		baseFont, ok := parts.TryGet[*tokens.NameToken](getDictEntry(dictionary, tokens.BaseFont), h.scanner)
		if ok && baseFont != nil {
			metrics, found := standard14fonts.GetAdobeFontMetrics(baseFont.Data())
			if !found {
				if h.isLenientParsing {
					metrics, found = standard14fonts.GetAdobeFontMetrics("Times-Roman")
					if !found {
						h.logger.Error(fmt.Sprintf("Type 1 Standard 14 font with name %v requested, this is an invalid name.", baseFont))
						return nil
					}
				} else {
					h.logger.Error(fmt.Sprintf("Type 1 Standard 14 font with name %v requested, this is an invalid name.", baseFont))
					return nil
				}
			}

			overrideEncoding := h.encodingReader.Read(dictionary, nil, nil)
			font, err := simple.NewType1Standard14Font(metrics, overrideEncoding)
			if err != nil {
				h.logger.Error(fmt.Sprintf("Could not create Type1Standard14Font: %v", err))
				return nil
			}
			return font
		}
	}

	descriptor, err := parser.GetFontDescriptor(h.scanner, dictionary)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Could not get font descriptor for Type 1 font: %v", err))
		return nil
	}

	fontProgram := h.parseFontProgram(descriptor)

	name, err := parser.GetName(h.scanner, dictionary, descriptor)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Could not get name for Type 1 font: %v", err))
		return nil
	}

	var toUnicodeCMap fonts.CMapProvider
	toUnicodeObj, hasToUnicode := dictionary.TryGet(tokens.ToUnicode)
	if hasToUnicode {
		streamToken, okStream := parts.TryGet[*tokens.StreamToken](toUnicodeObj, h.scanner)
		if okStream && streamToken != nil {
			cmapVal, found := h.cmapLocalCache.TryGetByStream(streamToken)
			if found {
				toUnicodeCMap = cmapVal
			}
		}
	}

	var fromFont *enc.Encoding
	if fontProgram.T1 != nil && len(fontProgram.T1.Encoding) > 0 {
		fromFont = enc.NewBuiltInEncoding(fontProgram.T1.Encoding).Encoding
	} else if fromFont == nil && fontProgram.Cff != nil {
		if cffEnc := encodingFromCFF(fontProgram.Cff); cffEnc != nil {
			fromFont = cffEnc
		}
	}

	encoding := h.encodingReader.Read(dictionary, descriptor, fromFont)

	if encoding == nil && fontProgram.T1 != nil && len(fontProgram.T1.Encoding) > 0 {
		encoding = enc.NewBuiltInEncoding(fontProgram.T1.Encoding).Encoding
	} else if encoding == nil && fontProgram.Cff != nil {
		if cffEnc := encodingFromCFF(fontProgram.Cff); cffEnc != nil {
			encoding = cffEnc
		}
	}

	result, err := simple.NewType1FontSimple(name, firstCharacter, lastCharacter, widths, descriptor, encoding, toUnicodeCMap, fontProgram)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Could not create Type1FontSimple: %v", err))
		return nil
	}

	return result
}

func (h *Type1FontHandler) parseFontProgram(descriptor *fonts.FontDescriptor) simple.Type1FontProgram {
	var result simple.Type1FontProgram

	if descriptor == nil || descriptor.FontFile == nil {
		return result
	}

	if descriptor.FontFile.ObjectKey.Data().ObjectNumber() == 0 {
		return result
	}

	streamToken, ok := parts.TryGet[*tokens.StreamToken](descriptor.FontFile.ObjectKey, h.scanner)
	if !ok || streamToken == nil {
		return result
	}

	bytes := h.decodeStream(streamToken)
	if bytes == nil {
		return result
	}

	subtypeName, hasSubtype := parts.TryGet[*tokens.NameToken](getDictEntry(streamToken.StreamDictionary, tokens.Subtype), h.scanner)
	if hasSubtype && subtypeName != nil && subtypeName.Equals(tokens.Type1C) {
		dataBytes := cff.NewCompactFontFormatData(bytes)
		parseCharStrings := func(
			charStringBytes [][]byte,
			subroutinesSelector *cff.CompactFontFormatSubroutinesSelector,
			charset cffcharset.CompactFontFormatCharset,
		) (cff.Type2CharStringsProvider, error) {
			return charstrings.Parse(charStringBytes, subroutinesSelector, charset)
		}
		cffCollection, err := cff.Parse(dataBytes, parseCharStrings)
		if err != nil || cffCollection == nil {
			h.logger.Error(fmt.Sprintf("Could not parse CFF font: %v", err))
			return result
		}
		result.Cff = cffCollection
		return result
	}

	length1Token, okL1 := streamToken.StreamDictionary.TryGet(tokens.Length1)
	length2Token, okL2 := streamToken.StreamDictionary.TryGet(tokens.Length2)
	if !okL1 || !okL2 {
		h.logger.Error("Type 1 font stream missing /Length1 or /Length2 entries")
		return result
	}

	length1Num, ok1 := length1Token.(*tokens.NumericToken)
	length2Num, ok2 := length2Token.(*tokens.NumericToken)
	if !ok1 || !ok2 {
		// Try resolving indirect references to get the actual numeric values.
		if ir, okIR := any(length1Token).(*tokens.IndirectReferenceToken); okIR && !ok1 {
			resolved, _ := parts.TryGet[*tokens.NumericToken](ir, h.scanner)
			if resolved != nil {
				length1Num = resolved
				ok1 = true
			}
		}
		if ir, okIR := any(length2Token).(*tokens.IndirectReferenceToken); okIR && !ok2 {
			resolved, _ := parts.TryGet[*tokens.NumericToken](ir, h.scanner)
			if resolved != nil {
				length2Num = resolved
				ok2 = true
			}
		}
		if !ok1 || !ok2 {
			h.logger.Error("Type 1 font stream /Length1 or /Length2 are not numeric")
			return result
		}
	}

	inputBytes := core.NewMemoryInputBytes(bytes)
	font, err := type1parser.Parse(inputBytes, int(length1Num.IntVal()), int(length2Num.IntVal()), h.stackDepthGuard)
	if err != nil || font == nil {
		h.logger.Error(fmt.Sprintf("Could not parse Type 1 font program: %v", err))
		return result
	}

	result.T1 = font
	return result
}

func (h *Type1FontHandler) decodeStream(streamToken *tokens.StreamToken) []byte {
	fl, err := h.filterProvider.GetFiltersWithScanner(streamToken.StreamDictionary, h.scanner)
	if err != nil || len(fl) == 0 {
		return streamToken.Data()
	}

	transform := streamToken.Data()
	for _, filter := range fl {
		result, decodeErr := filter.Decode(transform, streamToken.StreamDictionary, h.filterProvider, 0)
		if decodeErr != nil {
			h.logger.Error(fmt.Sprintf("Could not decode Type 1 font file stream: %v", decodeErr))
			return nil
		}
		transform = result
	}

	return transform
}

// encodingFromCFF extracts a *enc.Encoding from a CFF font collection by
// converting the internal encodingSource (CompactFontFormatBaseEncoding)
// to the common encodings.Encoding type.
func encodingFromCFF(collection *cff.CompactFontFormatFontCollection) *enc.Encoding {
	if collection == nil {
		return nil
	}
	firstFont := collection.FirstFont()
	if firstFont == nil {
		return nil
	}
	cffEnc := firstFont.Encoding()
	if cffEnc == nil {
		return nil
	}
	if base, ok := cffEnc.(interface{ ToEncoding() *enc.Encoding }); ok {
		return base.ToEncoding()
	}
	return nil
}

var _ FontHandler = (*Type1FontHandler)(nil)
