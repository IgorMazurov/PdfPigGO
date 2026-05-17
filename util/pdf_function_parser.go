package util

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/functions"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Create creates a PdfFunction from the given token. The token may be a direct
// StreamToken, DictionaryToken, or an IndirectReferenceToken that resolves to one.
func Create(function tokens.Token, scanner tokenization.PdfTokenScanner, filterProvider filters.LookupFilterProvider) (functions.PdfFunction, error) {
	var functionDict *tokens.DictionaryToken
	var functionStream *tokens.StreamToken

	if fs, ok := parts.TryGet[*tokens.StreamToken](function, scanner); ok && fs != nil {
		functionDict = fs.StreamDictionary
		decodedData, err := decodeStream(fs, filterProvider, scanner)
		if err != nil {
			return nil, fmt.Errorf("failed to decode function stream: %w", err)
		}
		streamToken, err := tokens.NewStreamToken(fs.StreamDictionary, decodedData)
		if err != nil {
			return nil, fmt.Errorf("failed to create decoded stream token: %w", err)
		}
		functionStream = streamToken
	} else if fd, ok := parts.TryGet[*tokens.DictionaryToken](function, scanner); ok && fd != nil {
		functionDict = fd
	} else {
		return nil, fmt.Errorf("invalid function token: expected StreamToken or DictionaryToken, got %T", function)
	}

	domain, err := getRequiredArray(functionDict, tokens.Domain, scanner)
	if err != nil {
		return nil, fmt.Errorf("missing required /Domain entry in function dictionary: %w", err)
	}

	rangeVals := getOptionalArray(functionDict, tokens.Range, scanner)

	funcTypeVal, ok := functionDict.TryGet(tokens.FunctionType)
	if !ok {
		return nil, core.NewPdfDocumentFormatException("function dictionary does not contain /FT entry")
	}

	numericFuncType, ok := funcTypeVal.(*tokens.NumericToken)
	if !ok || numericFuncType == nil {
		return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("/FT value is not a number: %v", funcTypeVal))
	}

	functionType := numericFuncType.IntVal()

	switch functionType {
	case 0:
		if functionStream == nil {
			return nil, fmt.Errorf("function type 0 requires a stream")
		}
		return createPdfFunctionType0(functionStream, domain, rangeVals, scanner)

	case 2:
		return createPdfFunctionType2(functionDict, domain, rangeVals, scanner)

	case 3:
		return createPdfFunctionType3(functionDict, domain, rangeVals, scanner, filterProvider)

	case 4:
		if functionStream == nil {
			return nil, fmt.Errorf("function type 4 requires a stream")
		}
		return createPdfFunctionType4(functionStream, domain, rangeVals, scanner)

	default:
		return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("unknown function type %d", functionType))
	}
}

// decodeStream decodes the raw stream data by applying all registered filters.
func decodeStream(stream *tokens.StreamToken, provider filters.LookupFilterProvider, scanner tokenization.PdfTokenScanner) ([]byte, error) {
	fl, err := provider.GetFiltersWithScanner(stream.StreamDictionary, scanner)
	if err != nil {
		return nil, err
	}

	transform := stream.Data()
	for i, filter := range fl {
		transform, err = filter.Decode(transform, stream.StreamDictionary, provider, i)
		if err != nil {
			return nil, err
		}
	}

	return transform, nil
}

// getOptionalArray retrieves an optional ArrayToken from the dictionary by key name,
// resolving indirect references through the scanner. Returns nil if not found.
func getOptionalArray(dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) *tokens.ArrayToken {
	token, found := dict.TryGet(name)
	if !found {
		return nil
	}

	arr, ok := parts.TryGet[*tokens.ArrayToken](token, scanner)
	if !ok {
		return nil
	}

	return arr
}

// createPdfFunctionType0 creates a sampled function (type 0).
func createPdfFunctionType0(functionStream *tokens.StreamToken, domain *tokens.ArrayToken, rangeVals *tokens.ArrayToken, scanner tokenization.PdfTokenScanner) (*functions.PdfFunctionType0, error) {
	if rangeVals == nil {
		return nil, fmt.Errorf("could not retrieve Range in type 0 function")
	}

	size, err := getRequiredArray(functionStream.StreamDictionary, tokens.Size, scanner)
	if err != nil {
		return nil, fmt.Errorf("missing /Size entry: %w", err)
	}

	bps, err := getRequiredNumeric(functionStream.StreamDictionary, tokens.BitsPerSample, scanner)
	if err != nil {
		return nil, fmt.Errorf("missing /BitsPerSample entry: %w", err)
	}

	order := 1
	if orderToken := tryGetNumeric(functionStream.StreamDictionary, tokens.Order, scanner); orderToken != nil {
		order = orderToken.IntVal()
	}

	var encode *tokens.ArrayToken
	if enc := getOptionalArray(functionStream.StreamDictionary, tokens.Encode, scanner); enc != nil {
		encode = enc
	} else {
		values := make([]tokens.Token, 0, size.Length())
		sizeValues := size.Data()
		for i := 0; i < len(sizeValues); i++ {
			values = append(values, tokens.NewNumericToken(0))
			numSize, ok := sizeValues[i].(*tokens.NumericToken)
			if !ok || numSize == nil {
				return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("/Size entry at index %d is not a number", i))
			}
			values = append(values, tokens.NewNumericToken(float64(numSize.IntVal()-1)))
		}
		encode = tokens.NewArrayToken(values)
	}

	var decode *tokens.ArrayToken
	if dec := getOptionalArray(functionStream.StreamDictionary, tokens.Decode, scanner); dec != nil {
		decode = dec
	} else {
		decode = rangeVals
	}

	return functions.NewPdfFunctionType0FromStream(functionStream, domain, rangeVals, size, bps.IntVal(), order, encode, decode), nil
}

// createPdfFunctionType2 creates an exponential interpolation function (type 2).
func createPdfFunctionType2(functionDict *tokens.DictionaryToken, domain *tokens.ArrayToken, rangeVals *tokens.ArrayToken, scanner tokenization.PdfTokenScanner) (*functions.PdfFunctionType2, error) {
	var array0 *tokens.ArrayToken
	if arr := getOptionalArray(functionDict, tokens.C0, scanner); arr != nil && arr.Length() > 0 {
		array0 = arr
	} else {
		array0 = tokens.NewArrayToken([]tokens.Token{tokens.NewNumericToken(0)})
	}

	var array1 *tokens.ArrayToken
	if arr := getOptionalArray(functionDict, tokens.C1, scanner); arr != nil && arr.Length() > 0 {
		array1 = arr
	} else {
		array1 = tokens.NewArrayToken([]tokens.Token{tokens.NewNumericToken(1)})
	}

	exp, err := getRequiredNumeric(functionDict, tokens.N, scanner)
	if err != nil {
		return nil, fmt.Errorf("missing /N entry: %w", err)
	}

	return functions.NewPdfFunctionType2FromDict(functionDict, domain, rangeVals, array0, array1, exp.DoubleVal()), nil
}

// createPdfFunctionType3 creates a stitching function (type 3).
func createPdfFunctionType3(functionDict *tokens.DictionaryToken, domain *tokens.ArrayToken, rangeVals *tokens.ArrayToken, scanner tokenization.PdfTokenScanner, filterProvider filters.LookupFilterProvider) (*functions.PdfFunctionType3, error) {
	functionsArray := make([]functions.PdfFunction, 0)

	functionsToken, err := getRequiredArray(functionDict, tokens.Functions, scanner)
	if err != nil {
		return nil, fmt.Errorf("missing /Functions entry: %w", err)
	}

	for _, fnToken := range functionsToken.Data() {
		var pdfFn functions.PdfFunction
		if streamTok, ok := parts.TryGet[*tokens.StreamToken](fnToken, scanner); ok && streamTok != nil {
			pdfFn, err = Create(streamTok, scanner, filterProvider)
			if err != nil {
				return nil, fmt.Errorf("failed to create nested function from stream: %w", err)
			}
		} else if dictTok, ok := parts.TryGet[*tokens.DictionaryToken](fnToken, scanner); ok && dictTok != nil {
			pdfFn, err = Create(dictTok, scanner, filterProvider)
			if err != nil {
				return nil, fmt.Errorf("failed to create nested function from dictionary: %w", err)
			}
		} else {
			return nil, fmt.Errorf("could not find function for token '%v' inside type 3 function", fnToken)
		}
		functionsArray = append(functionsArray, pdfFn)
	}

	bounds, err := getRequiredArray(functionDict, tokens.Bounds, scanner)
	if err != nil {
		return nil, fmt.Errorf("missing /Bounds entry: %w", err)
	}

	encode, err := getRequiredArray(functionDict, tokens.Encode, scanner)
	if err != nil {
		return nil, fmt.Errorf("missing /Encode entry: %w", err)
	}

	return functions.NewPdfFunctionType3FromDict(functionDict, domain, rangeVals, functionsArray, bounds, encode)
}

// createPdfFunctionType4 creates a PostScript calculator function (type 4).
func createPdfFunctionType4(functionStream *tokens.StreamToken, domain *tokens.ArrayToken, rangeVals *tokens.ArrayToken, _ tokenization.PdfTokenScanner) (*functions.PdfFunctionType4, error) {
	if rangeVals == nil {
		return nil, fmt.Errorf("could not retrieve Range in type 4 function")
	}

	return functions.NewPdfFunctionType4(functionStream, domain, rangeVals)
}

// tryGetNumeric retrieves an optional NumericToken from the dictionary by key name,
// resolving indirect references through the scanner. Returns nil if not found.
func tryGetNumeric(dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) *tokens.NumericToken {
	token, found := dict.TryGet(name)
	if !found {
		return nil
	}

	num, ok := parts.TryGet[*tokens.NumericToken](token, scanner)
	if !ok {
		return nil
	}

	return num
}
