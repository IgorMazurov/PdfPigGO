// Package parser provides font parsing utilities for PDF documents.
package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/cff"
	cffcharset "github.com/uglytoad/pdfpig/go/fonts/cff_charset"
	"github.com/uglytoad/pdfpig/go/fonts/cidfonts"
	"github.com/uglytoad/pdfpig/go/fonts/cff/charstrings"
	truetypeparser "github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/logging"
	pdffonts "github.com/uglytoad/pdfpig/go/pdf_fonts"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CidFontFactory creates CID fonts from their PDF dictionary representations.
type CidFontFactory struct {
	log            logging.Log
	scanner        tokenization.PdfTokenScanner
	filterProvider filters.LookupFilterProvider
}

// NewCidFontFactory creates a new CidFontFactory.
func NewCidFontFactory(log logging.Log, scanner tokenization.PdfTokenScanner, filterProvider filters.LookupFilterProvider) *CidFontFactory {
	return &CidFontFactory{
		log:            log,
		scanner:        scanner,
		filterProvider: filterProvider,
	}
}

// Generate creates a CID font from the given dictionary token.
// Returns nil if the subtype is not recognized (neither CIDFontType0 nor CIDFontType2).
func (f *CidFontFactory) Generate(dictionary *tokens.DictionaryToken) (cidfonts.CidFont, error) {
	typeToken := f.getNameOrDefault(dictionary, tokens.Type)
	if typeToken == nil || typeToken.Data() != tokens.Font.Data() {
		return nil, fonts.NewInvalidFontFormatException(fmt.Sprintf("Expected 'Font' dictionary but found '%v'", typeToken))
	}

	widths := f.readWidths(dictionary)

	var defaultWidth *float64
	if dwToken, ok := tryGetNumeric(dictionary, tokens.Dw, f.scanner); ok {
		dw := dwToken.DoubleVal()
		defaultWidth = &dw
	}

	verticalWritingMetrics := f.readVerticalDisplacements(dictionary)

	var descriptor *fonts.FontDescriptor
	if descDict, ok := f.tryGetFontDescriptor(dictionary); ok {
		desc, err := Generate(descDict, f.scanner)
		if err == nil && desc != nil {
			descriptor = desc
		}
	}

	var fontProgram cidfonts.CidFontProgram
	if descriptor != nil {
		fp, err := f.readDescriptorFile(descriptor)
		if err != nil {
			fontName := ""
			if descriptor.FontName != nil {
				fontName = descriptor.FontName.Data()
			}
			f.log.Error(fmt.Sprintf("Invalid descriptor in CID font named '%s': %v.", fontName, err))
		} else {
			fontProgram = fp
		}
	}

	baseFont := f.getNameOrDefault(dictionary, tokens.BaseFont)

	systemInfo, err := f.getSystemInfo(dictionary)
	if err != nil {
		return nil, err
	}

	subType := f.getNameOrDefault(dictionary, tokens.Subtype)

	if subType != nil && subType.Data() == tokens.CidFontType0.Data() {
		return cidfonts.NewType0CidFont(
			fontProgram,
			typeToken,
			subType,
			baseFont,
			systemInfo,
			descriptor,
			&verticalWritingMetrics,
			widths,
			defaultWidth,
		), nil
	}

	if subType != nil && subType.Data() == tokens.CidFontType2.Data() {
		cidToGid := f.getCharacterIdentifierToGlyphIndexMap(dictionary)
		return cidfonts.NewType2CidFont(
			typeToken,
			subType,
			baseFont,
			systemInfo,
			descriptor,
			fontProgram,
			&verticalWritingMetrics,
			widths,
			defaultWidth,
			&cidToGid,
		), nil
	}

	return nil, nil
}

// tryGetFontDescriptor attempts to get the font descriptor dictionary from the given dictionary.
func (f *CidFontFactory) tryGetFontDescriptor(dictionary *tokens.DictionaryToken) (*tokens.DictionaryToken, bool) {
	return tryGetDictionary(dictionary, tokens.FontDescriptor, f.scanner)
}

// readDescriptorFile reads and parses the font program from the font descriptor's file stream.
func (f *CidFontFactory) readDescriptorFile(descriptor *fonts.FontDescriptor) (cidfonts.CidFontProgram, error) {
	if descriptor.FontFile == nil {
		return nil, nil
	}

	fontFileStream, err := parts.GetByRef[*tokens.StreamToken](descriptor.FontFile.ObjectKey.Data(), f.scanner)
	if err != nil || fontFileStream == nil {
		return nil, nil
	}

	fontFileBytes := f.decodeStream(fontFileStream)
	if fontFileBytes == nil {
		return nil, fmt.Errorf("failed to decode font file stream")
	}

	switch descriptor.FontFile.FileType {
	case fonts.TrueType:
		return f.parseTrueTypeFont(fontFileBytes, fontFileStream, descriptor)
	case fonts.FromSubtype:
		return f.parseFromSubtype(fontFileStream)
	default:
		return nil, core.NewPdfDocumentFormatException("currently only TrueType fonts are supported")
	}
}

// parseTrueTypeFont parses a TrueType font from the decoded bytes, checking if it's actually CFF.
func (f *CidFontFactory) parseTrueTypeFont(fontFileBytes []byte, streamToken *tokens.StreamToken, descriptor *fonts.FontDescriptor) (cidfonts.CidFontProgram, error) {
	if isTrueTypeCff(fontFileBytes) {
		f.log.Warn("The CID TrueType font has the signature of a CFF font. Using CID CFF instead.")
		cffData := cff.NewCompactFontFormatData(fontFileBytes)
		parseCharStrings := func(charStringBytes [][]byte, subroutinesSelector *cff.CompactFontFormatSubroutinesSelector, charset cffcharset.CompactFontFormatCharset) (cff.Type2CharStringsProvider, error) {
			return charstrings.Parse(charStringBytes, subroutinesSelector, charset)
		}
		font, err := cff.Parse(cffData, parseCharStrings)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CFF font: %w", err)
		}
		return cidfonts.NewPdfCidCompactFontFormatFont(font), nil
	}

	input := truetypeparser.NewTrueTypeDataBytes(fontFileBytes)
	ttf, err := truetypeparser.ParseFull(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TrueType font: %w", err)
	}
	return cidfonts.NewPdfCidTrueTypeFont(ttf), nil
}

// parseFromSubtype parses a font based on the Subtype entry in the stream dictionary.
func (f *CidFontFactory) parseFromSubtype(streamToken *tokens.StreamToken) (cidfonts.CidFontProgram, error) {
	subtypeName, ok := parts.TryGet[*tokens.NameToken](getDictEntry(streamToken.StreamDictionary, tokens.Subtype), f.scanner)
	if !ok {
		return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("the font file stream did not contain a subtype entry: %v", streamToken.StreamDictionary))
	}

	bytes := f.decodeStream(streamToken)
	if bytes == nil {
		return nil, fmt.Errorf("failed to decode font file stream")
	}

	subtypeData := subtypeName.Data()
	if subtypeData == tokens.CidFontType0C.Data() || subtypeData == tokens.Type1C.Data() {
		cffData := cff.NewCompactFontFormatData(bytes)
		parseCharStrings := func(charStringBytes [][]byte, subroutinesSelector *cff.CompactFontFormatSubroutinesSelector, charset cffcharset.CompactFontFormatCharset) (cff.Type2CharStringsProvider, error) {
			return charstrings.Parse(charStringBytes, subroutinesSelector, charset)
		}
		font, parseErr := cff.Parse(cffData, parseCharStrings)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse CFF font: %w", parseErr)
		}
		return cidfonts.NewPdfCidCompactFontFormatFont(font), nil
	}

	if subtypeData == tokens.OpenType.Data() {
		if isTrueTypeCff(bytes) {
			f.log.Warn("The CID OpenType font has the signature of a CFF font. Using CID CFF instead.")
			cffData := cff.NewCompactFontFormatData(bytes)
			parseCharStrings := func(charStringBytes [][]byte, subroutinesSelector *cff.CompactFontFormatSubroutinesSelector, charset cffcharset.CompactFontFormatCharset) (cff.Type2CharStringsProvider, error) {
				return charstrings.Parse(charStringBytes, subroutinesSelector, charset)
			}
			font, parseErr := cff.Parse(cffData, parseCharStrings)
			if parseErr != nil {
				return nil, fmt.Errorf("failed to parse CFF font from OpenType stream: %w", parseErr)
			}
			return cidfonts.NewPdfCidCompactFontFormatFont(font), nil
		}

		input := truetypeparser.NewTrueTypeDataBytes(bytes)
		ttf, parseErr := truetypeparser.ParseFull(input)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse OpenType font as TrueType: %w", parseErr)
		}
		return cidfonts.NewPdfCidTrueTypeFont(ttf), nil
	}

	return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("unexpected subtype for CID font: %v", subtypeName))
}

// decodeStream decodes the stream data using the filter provider.
func (f *CidFontFactory) decodeStream(streamToken *tokens.StreamToken) []byte {
	fl, err := f.filterProvider.GetFiltersWithScanner(streamToken.StreamDictionary, f.scanner)
	if err != nil || len(fl) == 0 {
		return streamToken.Data()
	}

	transform := streamToken.Data()
	for i, filter := range fl {
		result, decodeErr := filter.Decode(transform, streamToken.StreamDictionary, f.filterProvider, i)
		if decodeErr != nil {
			f.log.Error(fmt.Sprintf("Could not decode CID font file stream: %v", decodeErr))
			return nil
		}
		transform = result
	}

	return transform
}

// readWidths reads the width dictionary from the /W entry.
func (f *CidFontFactory) readWidths(dict *tokens.DictionaryToken) map[int]float64 {
	widths := make(map[int]float64)

	widthArray, ok := tryGetArray(dict, tokens.W, f.scanner)
	if !ok {
		return widths
	}

	size := widthArray.Length()
	counter := 0
	for counter < size {
		firstCodeToken := widthArray.Get(counter)
		counter++
		firstCode, err := parts.GetByToken[*tokens.NumericToken](firstCodeToken, f.scanner)
		if err != nil {
			break
		}

		next := widthArray.Get(counter)
		counter++

		if arr, okArr := parts.TryGet[*tokens.ArrayToken](next, f.scanner); okArr {
			startRange := firstCode.IntVal()
			arraySize := arr.Length()
			for i := 0; i < arraySize; i++ {
				widthToken, err := parts.GetByToken[*tokens.NumericToken](arr.Get(i), f.scanner)
				if err != nil {
					break
				}
				widths[startRange+i] = widthToken.DoubleVal()
			}
		} else {
			secondCode, err := parts.GetByToken[*tokens.NumericToken](next, f.scanner)
			if err != nil {
				break
			}
			rangeWidthToken := widthArray.Get(counter)
			counter++
			rangeWidth, err := parts.GetByToken[*tokens.NumericToken](rangeWidthToken, f.scanner)
			if err != nil {
				break
			}
			startRange := firstCode.IntVal()
			endRange := secondCode.IntVal()
			width := rangeWidth.DoubleVal()
			for i := startRange; i <= endRange; i++ {
				widths[i] = width
			}
		}
	}

	return widths
}

// readVerticalDisplacements reads vertical writing metrics from the dictionary.
func (f *CidFontFactory) readVerticalDisplacements(dict *tokens.DictionaryToken) cidfonts.VerticalWritingMetrics {
	verticalDisplacements := make(map[int]float64)
	positionVectors := make(map[int]geometry.PdfVector)

	var dw2 cidfonts.VerticalVectorComponents = cidfonts.DefaultVerticalVectorComponents

	dw2Token, hasDw2 := dict.TryGet(tokens.Dw2)
	if hasDw2 {
		if arr, ok := dw2Token.(*tokens.ArrayToken); ok && arr.Length() >= 2 {
			posToken, err := parts.GetByToken[*tokens.NumericToken](arr.Get(0), f.scanner)
			if err == nil {
				dispToken, err2 := parts.GetByToken[*tokens.NumericToken](arr.Get(1), f.scanner)
				if err2 == nil {
					dw2 = cidfonts.NewVerticalVectorComponents(posToken.DoubleVal(), dispToken.DoubleVal())
				}
			}
		}
	}

	if w2, ok := tryGetArray(dict, tokens.W2, f.scanner); ok {
		i := 0
		for i < w2.Length() {
			cToken, err := parts.GetByToken[*tokens.NumericToken](w2.Get(i), f.scanner)
			if err != nil {
				break
			}
			i++

			next := w2.Get(i)
			i++

			if arr, okArr := parts.TryGet[*tokens.ArrayToken](next, f.scanner); okArr {
				for j := 0; j < arr.Length(); {
					cid := cToken.IntVal() + j
					w1yToken, err := parts.GetByToken[*tokens.NumericToken](arr.Get(j), f.scanner)
					if err != nil {
						break
					}
					j++
					v1xToken, err := parts.GetByToken[*tokens.NumericToken](arr.Get(j), f.scanner)
					if err != nil {
						break
					}
					j++
					v1yToken, err := parts.GetByToken[*tokens.NumericToken](arr.Get(j), f.scanner)
					if err != nil {
						break
					}
					j++

					verticalDisplacements[cid] = w1yToken.DoubleVal()
					positionVectors[cid] = geometry.NewPdfVector(v1xToken.DoubleVal(), v1yToken.DoubleVal())
				}
			} else {
				lastToken, err := parts.GetByToken[*tokens.NumericToken](next, f.scanner)
				if err != nil {
					break
				}
				first := cToken.IntVal()
				last := lastToken.IntVal()

				w1yRaw := w2.Get(i)
				i++
				v1xRaw := w2.Get(i)
				i++
				v1yRaw := w2.Get(i)
				i++

				w1yToken, err := parts.GetByToken[*tokens.NumericToken](w1yRaw, f.scanner)
				if err != nil {
					break
				}
				v1xToken, err := parts.GetByToken[*tokens.NumericToken](v1xRaw, f.scanner)
				if err != nil {
					break
				}
				v1yToken, err := parts.GetByToken[*tokens.NumericToken](v1yRaw, f.scanner)
				if err != nil {
					break
				}

				for cid := first; cid <= last; cid++ {
					verticalDisplacements[cid] = w1yToken.DoubleVal()
					positionVectors[cid] = geometry.NewPdfVector(v1xToken.DoubleVal(), v1yToken.DoubleVal())
				}
			}
		}
	}

	return cidfonts.NewVerticalWritingMetrics(dw2, verticalDisplacements, positionVectors)
}

// getSystemInfo reads the CID system info from the dictionary.
func (f *CidFontFactory) getSystemInfo(dictionary *tokens.DictionaryToken) (pdffonts.CharacterIdentifierSystemInfo, error) {
	cidEntry, ok := dictionary.TryGet(tokens.CidSystemInfo)
	if !ok {
		return pdffonts.CharacterIdentifierSystemInfo{}, fonts.NewInvalidFontFormatException(
			fmt.Sprintf("No CID System Info was found in the CID Font dictionary: %v", dictionary))
	}

	var cidDictionary *tokens.DictionaryToken
	if dict, ok := cidEntry.(*tokens.DictionaryToken); ok {
		cidDictionary = dict
	} else {
		resolved, err := parts.GetByToken[*tokens.DictionaryToken](cidEntry, f.scanner)
		if err != nil {
			return pdffonts.CharacterIdentifierSystemInfo{}, fonts.NewInvalidFontFormatException(
				fmt.Sprintf("Could not resolve CID System Info: %v", dictionary))
		}
		cidDictionary = resolved
	}

	registry := safeKeyAccess(cidDictionary, tokens.Registry, f.scanner)
	ordering := safeKeyAccess(cidDictionary, tokens.Ordering, f.scanner)
	supplement := 0
	if supToken, ok := cidDictionary.TryGet(tokens.Supplement); ok {
		if num, okNum := supToken.(*tokens.NumericToken); okNum {
			supplement = num.IntVal()
		}
	}

	return pdffonts.NewCharacterIdentifierSystemInfo(registry, ordering, supplement), nil
}

// getCharacterIdentifierToGlyphIndexMap reads the CIDToGIDMap from the dictionary.
func (f *CidFontFactory) getCharacterIdentifierToGlyphIndexMap(dictionary *tokens.DictionaryToken) cidfonts.CharacterIdentifierToGlyphIndexMap {
	entry, ok := dictionary.TryGet(tokens.CidToGidMap)
	if !ok {
		return *cidfonts.NewCharacterIdentifierToGlyphIndexMap()
	}

	if _, okName := parts.TryGet[*tokens.NameToken](entry, f.scanner); okName {
		return *cidfonts.NewCharacterIdentifierToGlyphIndexMap()
	}

	stream, err := parts.GetByToken[*tokens.StreamToken](entry, f.scanner)
	if err != nil {
		return *cidfonts.NewCharacterIdentifierToGlyphIndexMap()
	}

	bytes := f.decodeStream(stream)
	if bytes == nil {
		return *cidfonts.NewCharacterIdentifierToGlyphIndexMap()
	}

	return *cidfonts.NewCharacterIdentifierToGlyphIndexMapFromStream(bytes)
}

// getNameOrDefault gets a name token from the dictionary or returns nil if not found.
func (f *CidFontFactory) getNameOrDefault(dictionary *tokens.DictionaryToken, name *tokens.NameToken) *tokens.NameToken {
	token, ok := dictionary.TryGet(name)
	if !ok {
		return nil
	}
	if nameToken, ok := token.(*tokens.NameToken); ok {
		return nameToken
	}
	resolved, _ := parts.GetByToken[*tokens.NameToken](token, f.scanner)
	return resolved
}

// safeKeyAccess safely retrieves a string value from the dictionary for the given key.
func safeKeyAccess(dictionary *tokens.DictionaryToken, keyName *tokens.NameToken, scanner tokenization.PdfTokenScanner) string {
	token, ok := dictionary.TryGet(keyName)
	if !ok {
		return ""
	}

	if str, okStr := token.(*tokens.StringToken); okStr {
		return str.Data()
	}

	if hex, okHex := token.(*tokens.HexToken); okHex {
		return hex.Data()
	}

	if ref, okRef := token.(*tokens.IndirectReferenceToken); okRef {
		if stringToken, ok := parts.TryGet[*tokens.StringToken](ref, scanner); ok {
			return stringToken.Data()
		}
		if hexToken, ok := parts.TryGet[*tokens.HexToken](ref, scanner); ok {
			return hexToken.Data()
		}
		return ""
	}

	return ""
}

// getDictEntry returns the token for a given key from the dictionary, or nil if not found.
func getDictEntry(dictionary *tokens.DictionaryToken, name *tokens.NameToken) tokens.Token {
	token, _ := dictionary.TryGet(name)
	return token
}

// isTrueTypeCff checks whether the given byte slice has a CFF font signature.
// See https://docs.fileformat.com/font/cff/ and
// https://adobe-type-tools.github.io/font-tech-notes/pdfs/5176.CFF.pdf
func isTrueTypeCff(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	major := data[0]   // Major version
	minor := data[1]   // Minor version
	hdrSize := data[2] // Header size
	offSize := data[3] // Absolute offset size

	return major == 0x01 &&
		minor == 0x00 &&
		hdrSize >= 0x04 &&
		offSize >= 0x01 && offSize <= 0x04
}

// tryGetNumeric tries to get a NumericToken from the dictionary by name, resolving indirect refs.
func tryGetNumeric(dictionary *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (*tokens.NumericToken, bool) {
	token, ok := dictionary.TryGet(name)
	if !ok {
		return nil, false
	}
	if num, ok := token.(*tokens.NumericToken); ok {
		return num, true
	}
	resolved, ok := parts.TryGet[*tokens.NumericToken](token, scanner)
	return resolved, ok
}

// tryGetArray tries to get an ArrayToken from the dictionary by name, resolving indirect refs.
func tryGetArray(dictionary *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (*tokens.ArrayToken, bool) {
	token, ok := dictionary.TryGet(name)
	if !ok {
		return nil, false
	}
	if arr, ok := token.(*tokens.ArrayToken); ok {
		return arr, true
	}
	resolved, ok := parts.TryGet[*tokens.ArrayToken](token, scanner)
	return resolved, ok
}

// tryGetDictionary tries to get a DictionaryToken from the dictionary by name, resolving indirect refs.
func tryGetDictionary(dictionary *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (*tokens.DictionaryToken, bool) {
	token, ok := dictionary.TryGet(name)
	if !ok {
		return nil, false
	}
	if dict, ok := token.(*tokens.DictionaryToken); ok {
		return dict, true
	}
	resolved, ok := parts.TryGet[*tokens.DictionaryToken](token, scanner)
	return resolved, ok
}
