package util

import (
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/functions"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/tokenization"
)

// ColorSpaceParserProvider defines the capabilities needed by GetColorSpaceDetails.
type ColorSpaceParserProvider interface {
	filters.LookupFilterProvider
	DecodeStream(stream *tokens.StreamToken, scanner tokenization.PdfTokenScanner) []byte
}

// TryMap attempts to map a NameToken to a ColorSpace value.
// First tries direct mapping, then looks up the name in the resource store.
func TryMap(name *tokens.NameToken, resourceStore ParserResourceStore) (colors.ColorSpace, bool) {
	if cs, ok := colors.TryMapToColorSpace(name); ok {
		return cs, true
	}

	resourceCS, found := resourceStore.TryGetNamedColorSpace(name)
	if !found {
		return 0, false
	}

	cs, ok := colors.TryMapToColorSpace(resourceCS.Name)
	return cs, ok
}

// GetColorSpaceDetails parses color space details from the given parameters.
func GetColorSpaceDetails(
	colorSpace *colors.ColorSpace,
	imageDictionary *tokens.DictionaryToken,
	scanner tokenization.PdfTokenScanner,
	resourceStore ParserResourceStore,
	filterProvider ColorSpaceParserProvider,
	cannotRecurse bool,
) colors.ColorSpaceDetails {
	if isImageMaskOrCcitt(imageDictionary, scanner, filterProvider) {
		if cannotRecurse {
			return colors.DeviceGrayColorSpaceDetails
		}

		csWithoutFilters := imageDictionary.Without(tokens.Filter).Without(tokens.F)
		innerDetails := GetColorSpaceDetails(colorSpace, csWithoutFilters, scanner, resourceStore, filterProvider, true)
		return colors.StencilIndexedColorSpace(innerDetails)
	}

	if colorSpace == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	switch *colorSpace {
	case colors.DeviceGray:
		return colors.DeviceGrayColorSpaceDetails

	case colors.DeviceRGB:
		return colors.DeviceRgbColorSpaceDetails

	case colors.DeviceCMYK:
		return colors.DeviceCmykColorSpaceDetails

	case colors.CalGray:
		return parseCalGray(imageDictionary, scanner, resourceStore)

	case colors.CalRGB:
		return parseCalRGB(imageDictionary, scanner, resourceStore)

	case colors.Lab:
		return parseLab(imageDictionary, scanner, resourceStore)

	case colors.ICCBased:
		return parseICCBased(imageDictionary, scanner, resourceStore, filterProvider)

	case colors.Indexed:
		if cannotRecurse {
			return colors.UnsupportedColorSpaceDetails
		}
		return parseIndexed(imageDictionary, scanner, resourceStore, filterProvider)

	case colors.Pattern:
		if cannotRecurse {
			return colors.UnsupportedColorSpaceDetails
		}
		return parsePattern(imageDictionary, scanner, resourceStore, filterProvider)

	case colors.Separation:
		return parseSeparation(imageDictionary, scanner, resourceStore, filterProvider)

	case colors.DeviceN:
		return parseDeviceN(imageDictionary, scanner, resourceStore, filterProvider)

	default:
		return colors.UnsupportedColorSpaceDetails
	}
}

// isImageMaskOrCcitt checks if the image dictionary represents an imagemask or uses CCITT fax decode.
func isImageMaskOrCcitt(dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner, fp ColorSpaceParserProvider) bool {
	if tok, found := dict.TryGet(tokens.ImageMask); found {
		if isImageMask, ok := parts.TryGet[*tokens.BooleanToken](tok, scanner); ok && isImageMask != nil && isImageMask.Data() {
			return true
		}
	}
	if tok, found := dict.TryGet(tokens.Im); found {
		if isIm, ok := parts.TryGet[*tokens.BooleanToken](tok, scanner); ok && isIm != nil && isIm.Data() {
			return true
		}
	}
	if filtersList, err := fp.GetFilters(dict); err == nil {
		for _, f := range filtersList {
			if _, ok := f.(*filters.CcittFaxDecodeFilter); ok {
				return true
			}
		}
	}
	return false
}

// parseCalGray parses a CalGray color space from the image dictionary.
func parseCalGray(dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner, rs ParserResourceStore) colors.ColorSpaceDetails {
	csArray, ok := tryGetColorSpaceArray(dict, rs, scanner)
	if !ok || csArray.Length() != 2 {
		return colors.UnsupportedColorSpaceDetails
	}

	first, ok := csArray.Get(0).(*tokens.NameToken)
	if !ok {
		return colors.UnsupportedColorSpaceDetails
	}
	innerCS, mapped := TryMap(first, rs)
	if !mapped || innerCS != colors.CalGray {
		return colors.UnsupportedColorSpaceDetails
	}

	dictToken, ok := parts.TryGet[*tokens.DictionaryToken](csArray.Get(1), scanner)
	if !ok || dictToken == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	wpTok, found := dictToken.TryGet(tokens.WhitePoint)
	if !found {
		return colors.UnsupportedColorSpaceDetails
	}
	wpToken, ok := parts.TryGet[*tokens.ArrayToken](wpTok, scanner)
	if !ok || wpToken == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	whitePoint := extractDoubles(wpToken)

	var blackPoint []float64
	if bpTok, found := dictToken.TryGet(tokens.BlackPoint); found {
		if bpToken, ok := parts.TryGet[*tokens.ArrayToken](bpTok, scanner); ok && bpToken != nil {
			blackPoint = extractDoubles(bpToken)
		}
	}

	var gamma float64
	if gTok, found := dictToken.TryGet(tokens.Gamma); found {
		if gToken, ok := parts.TryGet[*tokens.NumericToken](gTok, scanner); ok && gToken != nil {
			gamma = gToken.DoubleVal()
		}
	}

	details, err := colors.NewCalGrayColorSpaceDetails(whitePoint, blackPoint, gamma)
	if err != nil {
		return colors.UnsupportedColorSpaceDetails
	}
	return details
}

// parseCalRGB parses a CalRGB color space from the image dictionary.
func parseCalRGB(dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner, rs ParserResourceStore) colors.ColorSpaceDetails {
	csArray, ok := tryGetColorSpaceArray(dict, rs, scanner)
	if !ok || csArray.Length() != 2 {
		return colors.UnsupportedColorSpaceDetails
	}

	first, ok := csArray.Get(0).(*tokens.NameToken)
	if !ok {
		return colors.UnsupportedColorSpaceDetails
	}
	innerCS, mapped := TryMap(first, rs)
	if !mapped || innerCS != colors.CalRGB {
		return colors.UnsupportedColorSpaceDetails
	}

	dictToken, ok := parts.TryGet[*tokens.DictionaryToken](csArray.Get(1), scanner)
	if !ok || dictToken == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	wpTok, found := dictToken.TryGet(tokens.WhitePoint)
	if !found {
		return colors.UnsupportedColorSpaceDetails
	}
	wpToken, ok := parts.TryGet[*tokens.ArrayToken](wpTok, scanner)
	if !ok || wpToken == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	whitePoint := extractDoubles(wpToken)

	var blackPoint []float64
	if bpTok, found := dictToken.TryGet(tokens.BlackPoint); found {
		if bpToken, ok := parts.TryGet[*tokens.ArrayToken](bpTok, scanner); ok && bpToken != nil {
			blackPoint = extractDoubles(bpToken)
		}
	}

	var gamma []float64
	if gTok, found := dictToken.TryGet(tokens.Gamma); found {
		if gToken, ok := parts.TryGet[*tokens.ArrayToken](gTok, scanner); ok && gToken != nil {
			gamma = extractDoubles(gToken)
		}
	}

	var matrix []float64
	if mTok, found := dictToken.TryGet(tokens.Matrix); found {
		if mToken, ok := parts.TryGet[*tokens.ArrayToken](mTok, scanner); ok && mToken != nil {
			matrix = extractDoubles(mToken)
		}
	}

	details, err := colors.NewCalRGBColorSpaceDetails(whitePoint, blackPoint, gamma, matrix)
	if err != nil {
		return colors.UnsupportedColorSpaceDetails
	}
	return details
}

// parseLab parses a Lab color space from the image dictionary.
func parseLab(dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner, rs ParserResourceStore) colors.ColorSpaceDetails {
	csArray, ok := tryGetColorSpaceArray(dict, rs, scanner)
	if !ok || csArray.Length() != 2 {
		return colors.UnsupportedColorSpaceDetails
	}

	first, ok := csArray.Get(0).(*tokens.NameToken)
	if !ok {
		return colors.UnsupportedColorSpaceDetails
	}
	innerCS, mapped := TryMap(first, rs)
	if !mapped || innerCS != colors.Lab {
		return colors.UnsupportedColorSpaceDetails
	}

	dictToken, ok := parts.TryGet[*tokens.DictionaryToken](csArray.Get(1), scanner)
	if !ok || dictToken == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	wpTok, found := dictToken.TryGet(tokens.WhitePoint)
	if !found {
		return colors.UnsupportedColorSpaceDetails
	}
	wpToken, ok := parts.TryGet[*tokens.ArrayToken](wpTok, scanner)
	if !ok || wpToken == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	whitePoint := extractDoubles(wpToken)

	var blackPoint []float64
	if bpTok, found := dictToken.TryGet(tokens.BlackPoint); found {
		if bpToken, ok := parts.TryGet[*tokens.ArrayToken](bpTok, scanner); ok && bpToken != nil {
			blackPoint = extractDoubles(bpToken)
		}
	}

	var matrix []float64
	if mTok, found := dictToken.TryGet(tokens.Matrix); found {
		if mToken, ok := parts.TryGet[*tokens.ArrayToken](mTok, scanner); ok && mToken != nil {
			matrix = extractDoubles(mToken)
		}
	}

	details, err := colors.NewLabColorSpaceDetails(whitePoint, blackPoint, matrix)
	if err != nil {
		return colors.UnsupportedColorSpaceDetails
	}
	return details
}

// parseICCBased parses an ICCBased color space from the image dictionary.
func parseICCBased(dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner, rs ParserResourceStore, fp ColorSpaceParserProvider) colors.ColorSpaceDetails {
	csArray, ok := tryGetColorSpaceArray(dict, rs, scanner)
	if !ok || csArray.Length() != 2 {
		return colors.UnsupportedColorSpaceDetails
	}

	first, ok := csArray.Get(0).(*tokens.NameToken)
	if !ok {
		return colors.UnsupportedColorSpaceDetails
	}
	innerCS, mapped := TryMap(first, rs)
	if !mapped || innerCS != colors.ICCBased {
		return colors.UnsupportedColorSpaceDetails
	}

	streamToken, ok := parts.TryGet[*tokens.StreamToken](csArray.Get(1), scanner)
	if !ok || streamToken == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	nTok, found := streamToken.StreamDictionary.TryGet(tokens.N)
	if !found {
		return colors.UnsupportedColorSpaceDetails
	}
	nToken, ok := parts.TryGet[*tokens.NumericToken](nTok, scanner)
	if !ok || nToken == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	var alternateCS colors.ColorSpaceDetails
	if altTok, found := streamToken.StreamDictionary.TryGet(tokens.Alternate); found {
		if altNameToken, ok := altTok.(*tokens.NameToken); ok {
			if altCS, mapped := TryMap(altNameToken, rs); mapped {
				alternateCS = GetColorSpaceDetails(&altCS, dict, scanner, rs, fp, true)
			}
		}
	}

	var rangeVals []float64
	if rTok, found := streamToken.StreamDictionary.TryGet(tokens.Range); found {
		if rToken, ok := parts.TryGet[*tokens.ArrayToken](rTok, scanner); ok && rToken != nil {
			rangeVals = extractDoubles(rToken)
		}
	}

	details, err := colors.NewICCBasedColorSpaceDetails(nToken.IntVal(), alternateCS, rangeVals)
	if err != nil {
		return colors.UnsupportedColorSpaceDetails
	}

	// TODO: C# also extracts optional Metadata stream (XmpMetadata) from the ICCBased
	// stream dictionary and passes it to ICCBasedColorSpaceDetails. The Go struct
	// iccBasedColorSpaceDetails currently has no metadata field. Adding it would require
	// modifying the colors package (struct + constructor). XMP metadata for ICC profiles
	// is optional and does not affect core color operations, so this gap is acceptable.

	return details
}

// parseIndexed parses an Indexed color space from the image dictionary.
func parseIndexed(dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner, rs ParserResourceStore, fp ColorSpaceParserProvider) colors.ColorSpaceDetails {
	csArray, ok := tryGetColorSpaceArray(dict, rs, scanner)
	if !ok || csArray.Length() != 4 {
		return colors.UnsupportedColorSpaceDetails
	}

	first, ok := csArray.Get(0).(*tokens.NameToken)
	if !ok {
		return colors.UnsupportedColorSpaceDetails
	}
	innerCS, mapped := TryMap(first, rs)
	if !mapped || innerCS != colors.Indexed {
		return colors.UnsupportedColorSpaceDetails
	}

	baseDetails := getSecondaryColorSpace(csArray.Get(1), dict, scanner, fp, rs)
	if colors.IsUnsupported(baseDetails) {
		return colors.UnsupportedColorSpaceDetails
	}

	hiValToken, ok := parts.TryGet[*tokens.NumericToken](csArray.Get(2), scanner)
	if !ok || hiValToken == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	hival := byte(hiValToken.IntVal())

	tableBytes := extractTableBytes(csArray.Get(3), scanner, fp)
	if tableBytes == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	details := colors.NewIndexedColorSpaceDetails(baseDetails, hival, tableBytes)

	return details
}

// parsePattern parses a Pattern color space from the image dictionary.
func parsePattern(dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner, rs ParserResourceStore, fp ColorSpaceParserProvider) colors.ColorSpaceDetails {
	var underlyingCS colors.ColorSpaceDetails = colors.UnsupportedColorSpaceDetails

	if len(dict.Data()) > 0 {
		csArray, ok := tryGetColorSpaceArray(dict, rs, scanner)
		if !ok || (csArray.Length() != 1 && csArray.Length() != 2) {
			return colors.UnsupportedColorSpaceDetails
		}

		patternName, ok := parts.TryGet[*tokens.NameToken](csArray.Get(0), scanner)
		if !ok || patternName == nil || !patternName.Equals(tokens.Pattern) {
			return colors.UnsupportedColorSpaceDetails
		}

		if csArray.Length() > 1 {
			underlyingCS = getSecondaryColorSpace(csArray.Get(1), dict, scanner, fp, rs)
		}
	}

	patterns := rs.GetPatterns()
	// Both TilingPatternColor and ShadingPatternColor implement both PatternColor and Color
	// interfaces (they have ColorSpace() and ToRGBValues() methods), so runtime type assertion works.
	patternMap := make(map[*tokens.NameToken]colors.Color, len(patterns))
	for name, pat := range patterns {
		if c, ok := pat.(colors.Color); ok {
			patternMap[name] = c
		} else {
			patternMap[name] = nil
		}
	}

	return colors.NewPatternColorSpaceDetails(patternMap, underlyingCS)
}

// parseSeparation parses a Separation color space from the image dictionary.
func parseSeparation(dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner, rs ParserResourceStore, fp ColorSpaceParserProvider) colors.ColorSpaceDetails {
	csArray, ok := tryGetColorSpaceArray(dict, rs, scanner)
	if !ok || csArray.Length() != 4 {
		return colors.UnsupportedColorSpaceDetails
	}

	separationName, ok := parts.TryGet[*tokens.NameToken](csArray.Get(0), scanner)
	if !ok || separationName == nil || !separationName.Equals(tokens.Separation) {
		return colors.UnsupportedColorSpaceDetails
	}

	nameToken, ok := parts.TryGet[*tokens.NameToken](csArray.Get(1), scanner)
	if !ok || nameToken == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	alternateCS := getSecondaryColorSpace(csArray.Get(2), dict, scanner, fp, rs)
	fn := createPdfFunction(csArray.Get(3), scanner, fp)
	if fn == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	return colors.NewSeparationColorSpaceDetails(nameToken, alternateCS, fn)
}

// parseDeviceN parses a DeviceN color space from the image dictionary.
func parseDeviceN(dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner, rs ParserResourceStore, fp ColorSpaceParserProvider) colors.ColorSpaceDetails {
	csArray, ok := tryGetColorSpaceArray(dict, rs, scanner)
	if !ok || (csArray.Length() != 4 && csArray.Length() != 5) {
		return colors.UnsupportedColorSpaceDetails
	}

	devnName, ok := parts.TryGet[*tokens.NameToken](csArray.Get(0), scanner)
	if !ok || devnName == nil || !devnName.Equals(tokens.Devicen) {
		return colors.UnsupportedColorSpaceDetails
	}

	namesToken, ok := parts.TryGet[*tokens.ArrayToken](csArray.Get(1), scanner)
	if !ok || namesToken == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	alternateCS := getSecondaryColorSpace(csArray.Get(2), dict, scanner, fp, rs)
	tintFunc := createPdfFunction(csArray.Get(3), scanner, fp)
	if tintFunc == nil {
		return colors.UnsupportedColorSpaceDetails
	}

	var attrs *colors.DeviceNColorSpaceAttributes
	if csArray.Length() > 4 {
		attrDict, ok := parts.TryGet[*tokens.DictionaryToken](csArray.Get(4), scanner)
		if ok && attrDict != nil {
			subtype := tokens.Devicen
			if s, found := attrDict.TryGet(tokens.Subtype); found {
				if sn, ok := s.(*tokens.NameToken); ok {
					subtype = sn
				}
			}

			var colorants *tokens.DictionaryToken
			if c, found := attrDict.TryGet(tokens.Colorants); found {
				if cd, ok := c.(*tokens.DictionaryToken); ok {
					colorants = cd
				}
			}

			var process *tokens.DictionaryToken
			if p, found := attrDict.TryGet(tokens.Process); found {
				if pd, ok := p.(*tokens.DictionaryToken); ok {
					process = pd
				}
			}

			var mixingHints *tokens.DictionaryToken
			if m, found := attrDict.TryGet(tokens.MixingHints); found {
				if md, ok := m.(*tokens.DictionaryToken); ok {
					mixingHints = md
				}
			}

			attrs = &colors.DeviceNColorSpaceAttributes{
				Subtype:     subtype,
				Colorants:   colorants,
				Process:     process,
				MixingHints: mixingHints,
			}
		}
	}

	nameTokens := make([]*tokens.NameToken, 0, namesToken.Length())
	for i := 0; i < namesToken.Length(); i++ {
		if nt, ok := namesToken.Get(i).(*tokens.NameToken); ok {
			nameTokens = append(nameTokens, nt)
		}
	}

	if attrs != nil {
		return colors.NewDeviceNColorSpaceDetails(nameTokens, alternateCS, tintFunc, attrs)
	}
	return colors.NewDeviceNColorSpaceDetails(nameTokens, alternateCS, tintFunc, nil)
}

// getSecondaryColorSpace resolves a secondary color space token to ColorSpaceDetails.
func getSecondaryColorSpace(csToken tokens.Token, dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner, fp ColorSpaceParserProvider, rs ParserResourceStore) colors.ColorSpaceDetails {
	if altName, ok := parts.TryGet[*tokens.NameToken](csToken, scanner); ok && altName != nil {
		if baseCS, mapped := TryMap(altName, rs); mapped {
			return GetColorSpaceDetails(&baseCS, dict, scanner, rs, fp, true)
		}
	}

	if altArray, ok := parts.TryGet[*tokens.ArrayToken](csToken, scanner); ok && altArray != nil && altArray.Length() > 0 {
		if altName, ok := altArray.Get(0).(*tokens.NameToken); ok {
			if altCS, mapped := TryMap(altName, rs); mapped {
				pseudoDict, err := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
					tokens.ColorSpace: altArray,
				})
				if err != nil {
					return colors.UnsupportedColorSpaceDetails
				}
				return GetColorSpaceDetails(&altCS, pseudoDict, scanner, rs, fp, true)
			}
		}
	}

	return colors.UnsupportedColorSpaceDetails
}

// tryGetColorSpaceArray extracts the color space array from the image dictionary.
func tryGetColorSpaceArray(dict *tokens.DictionaryToken, rs ParserResourceStore, scanner tokenization.PdfTokenScanner) (*tokens.ArrayToken, bool) {
	csToken := getObjectOrDefault(dict, tokens.ColorSpace, tokens.Cs)

	csArray, ok := parts.TryGet[*tokens.ArrayToken](csToken, scanner)
	if ok && csArray != nil {
		return csArray, true
	}

	if csName, ok := parts.TryGet[*tokens.NameToken](csToken, scanner); ok && csName != nil {
		if namedCS, found := rs.TryGetNamedColorSpace(csName); found {
			if arr, ok := namedCS.Data.(*tokens.ArrayToken); ok {
				return arr, true
			}
		}
	}

	return nil, false
}

// getObjectOrDefault returns the token for the first key that exists in the dictionary.
func getObjectOrDefault(dict *tokens.DictionaryToken, keys ...*tokens.NameToken) tokens.Token {
	for _, key := range keys {
		if tok, found := dict.TryGet(key); found {
			return tok
		}
	}
	return nil
}

// extractDoubles extracts float64 values from an ArrayToken of NumericTokens.
func extractDoubles(arr *tokens.ArrayToken) []float64 {
	data := arr.Data()
	result := make([]float64, 0, len(data))
	for _, tok := range data {
		if nt, ok := tok.(*tokens.NumericToken); ok {
			result = append(result, nt.DoubleVal())
		}
	}
	return result
}

// extractTableBytes extracts the color table bytes from an Indexed color space's fourth element.
func extractTableBytes(token tokens.Token, scanner tokenization.PdfTokenScanner, fp ColorSpaceParserProvider) []byte {
	if hexTok, ok := parts.TryGet[*tokens.HexToken](token, scanner); ok && hexTok != nil {
		return []byte(hexTok.Data())
	}

	if streamTok, ok := parts.TryGet[*tokens.StreamToken](token, scanner); ok && streamTok != nil {
		return fp.DecodeStream(streamTok, scanner)
	}

	if strTok, ok := parts.TryGet[*tokens.StringToken](token, scanner); ok && strTok != nil {
		return strTok.GetBytes()
	}

	return nil
}

// createPdfFunction creates a PdfFunction from the given token.
// Returns nil if the function cannot be created (e.g., unsupported type or missing entries).
func createPdfFunction(token tokens.Token, scanner tokenization.PdfTokenScanner, fp ColorSpaceParserProvider) functions.PdfFunction {
	fn, err := Create(token, scanner, fp)
	if err != nil {
		return nil
	}
	return fn
}
