package util

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// ParserResourceStore defines the ResourceStore methods needed by pattern/shading/color parsing.
type ParserResourceStore interface {
	GetColorSpaceDetails(name *tokens.NameToken, dictionary *tokens.DictionaryToken) colors.ColorSpaceDetails
	TryGetNamedColorSpace(name *tokens.NameToken) (colors.ResourceColorSpace, bool)
	GetPatterns() map[*tokens.NameToken]colors.PatternColor
}

// PatternParserProvider defines the capabilities needed by CreatePattern.
type PatternParserProvider interface {
	filters.LookupFilterProvider
	DecodeStream(stream *tokens.StreamToken, scanner tokenization.PdfTokenScanner) []byte
}

// CreatePattern creates a PatternColor from a pattern token by resolving the
// underlying dictionary or stream and dispatching to the appropriate parser
// based on the /PatternType entry (1 = tiling, 2 = shading).
func CreatePattern(
	pattern tokens.Token,
	scanner tokenization.PdfTokenScanner,
	resourceStore ParserResourceStore,
	filterProvider PatternParserProvider,
) (colors.PatternColor, error) {
	var patternDictionary *tokens.DictionaryToken
	var patternStream *tokens.StreamToken

	if streamToken, ok := parts.TryGet[*tokens.StreamToken](pattern, scanner); ok {
		patternDictionary = streamToken.StreamDictionary
		decoded := filterProvider.DecodeStream(streamToken, scanner)
		patternStream, _ = tokens.NewStreamToken(streamToken.StreamDictionary, decoded)
	} else if dictToken, ok := parts.TryGet[*tokens.DictionaryToken](pattern, scanner); ok {
		patternDictionary = dictToken
	} else {
		return nil, core.NewPdfDocumentFormatException(
			fmt.Sprintf("invalid Pattern token encountered in page resource dictionary: %v", pattern))
	}

	if _, hasPatternType := patternDictionary.TryGet(tokens.PatternType); !hasPatternType {
		return nil, fmt.Errorf("pattern dictionary missing /PatternType entry")
	}

	patternTypeRaw, _ := patternDictionary.TryGet(tokens.PatternType)
	patternTypeToken, err := parts.GetByToken[*tokens.NumericToken](patternTypeRaw, scanner)
	if err != nil {
		return nil, err
	}
	patternTypeInt := patternTypeToken.IntVal()

	matrix, err := parsePatternMatrix(patternDictionary, scanner)
	if err != nil {
		return nil, err
	}

	var patternExtGState *tokens.DictionaryToken
	if _, hasExtGState := patternDictionary.TryGet(tokens.ExtGState); hasExtGState {
		patternExtGState, _ = TryGetOptionalTokenDirect[*tokens.DictionaryToken](patternDictionary, tokens.ExtGState, scanner)
	}

	switch patternTypeInt {
	case 1:
		return createTilingPattern(patternStream, patternExtGState, matrix, scanner)
	case 2:
		return createShadingPattern(patternDictionary, patternExtGState, matrix, scanner, resourceStore, filterProvider)
	default:
		return nil, core.NewPdfDocumentFormatException(
			fmt.Sprintf("invalid Pattern type encountered in page resource dictionary: %d", patternTypeInt))
	}
}

// parsePatternMatrix extracts the /Matrix entry from the pattern dictionary.
// If absent or unresolvable, returns the identity matrix [1 0 0 1 0 0].
func parsePatternMatrix(dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) (core.TransformationMatrix, error) {
	if _, hasMatrix := dict.TryGet(tokens.Matrix); !hasMatrix {
		return core.NewTransformationMatrixFromSlice([]float64{1, 0, 0, 1, 0, 0}), nil
	}

	matrixArray, ok := TryGetOptionalTokenDirect[*tokens.ArrayToken](dict, tokens.Matrix, scanner)
	if !ok {
		return core.NewTransformationMatrixFromSlice([]float64{1, 0, 0, 1, 0, 0}), nil
	}

	numbers := make([]float64, 0, matrixArray.Length())
	for _, elem := range matrixArray.Data() {
		if num, ok := elem.(*tokens.NumericToken); ok {
			numbers = append(numbers, num.DoubleVal())
		}
	}

	if len(numbers) < 6 {
		return core.TransformationMatrix{}, fmt.Errorf("pattern /Matrix array must contain at least 6 numeric values, got %d", len(numbers))
	}

	return core.NewTransformationMatrixFromSlice(numbers[:6]), nil
}

// createTilingPattern parses a tiling pattern (type 1) from the stream dictionary.
func createTilingPattern(
	patternStream *tokens.StreamToken,
	patternExtGState *tokens.DictionaryToken,
	matrix core.TransformationMatrix,
	scanner tokenization.PdfTokenScanner,
) (colors.PatternColor, error) {
	dict := patternStream.StreamDictionary

	paintTypeToken, err := getRequiredNumeric(dict, tokens.PaintType, scanner)
	if err != nil {
		return nil, core.NewPdfDocumentFormatException("invalid Pattern token encountered: missing /PaintType")
	}

	tilingTypeToken, err := getRequiredNumeric(dict, tokens.TilingType, scanner)
	if err != nil {
		return nil, core.NewPdfDocumentFormatException("invalid Pattern token encountered: missing /TilingType")
	}

	bboxToken, err := getRequiredArray(dict, tokens.Bbox, scanner)
	if err != nil {
		return nil, core.NewPdfDocumentFormatException("invalid Pattern token encountered: missing /BBox")
	}

	xStepToken, err := getRequiredNumeric(dict, tokens.XStep, scanner)
	if err != nil {
		return nil, core.NewPdfDocumentFormatException("invalid Pattern token encountered: missing /XStep")
	}

	yStepToken, err := getRequiredNumeric(dict, tokens.YStep, scanner)
	if err != nil {
		return nil, core.NewPdfDocumentFormatException("invalid Pattern token encountered: missing /YStep")
	}

	resourcesToken, err := getRequiredDictionary(dict, tokens.Resources, scanner)
	if err != nil {
		return nil, core.NewPdfDocumentFormatException("invalid Pattern token encountered: missing /Resources")
	}

	bBox, err := ToRectangle(bboxToken, scanner)
	if err != nil {
		return nil, err
	}

	bBoxVal := *bBox

	return colors.NewTilingPatternColor(
		matrix,
		patternExtGState,
		patternStream,
		colors.PatternPaintType(paintTypeToken.IntVal()),
		colors.PatternTilingType(tilingTypeToken.IntVal()),
		bBoxVal,
		xStepToken.DoubleVal(),
		yStepToken.DoubleVal(),
		resourcesToken,
		patternStream.Data(),
	), nil
}

// createShadingPattern parses a shading pattern (type 2) from the dictionary.
func createShadingPattern(
	patternDictionary *tokens.DictionaryToken,
	patternExtGState *tokens.DictionaryToken,
	matrix core.TransformationMatrix,
	scanner tokenization.PdfTokenScanner,
	resourceStore ParserResourceStore,
	filterProvider PatternParserProvider,
) (colors.PatternColor, error) {
	shadingToken, ok := patternDictionary.TryGet(tokens.Shading)
	if !ok {
		return nil, core.NewPdfDocumentFormatException("shading pattern missing /Shading entry")
	}

	shading, err := CreateShading(shadingToken, scanner, resourceStore, filterProvider)
	if err != nil {
		return nil, err
	}

	return colors.NewShadingPatternColor(matrix, patternExtGState, patternDictionary, shading), nil
}

// getRequiredNumeric retrieves a required NumericToken by key name, resolving indirect references.
func getRequiredNumeric(dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (*tokens.NumericToken, error) {
	token, ok := dict.TryGet(name)
	if !ok {
		return nil, fmt.Errorf("missing required key %v", name)
	}
	result, found := parts.TryGet[*tokens.NumericToken](token, scanner)
	if !found {
		return nil, fmt.Errorf("key %v is not a numeric token", name)
	}
	return result, nil
}

// getRequiredArray retrieves a required ArrayToken by key name, resolving indirect references.
func getRequiredArray(dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (*tokens.ArrayToken, error) {
	token, ok := dict.TryGet(name)
	if !ok {
		return nil, fmt.Errorf("missing required key %v", name)
	}
	result, found := parts.TryGet[*tokens.ArrayToken](token, scanner)
	if !found {
		return nil, fmt.Errorf("key %v is not an array token", name)
	}
	return result, nil
}

// getRequiredDictionary retrieves a required DictionaryToken by key name, resolving indirect references.
func getRequiredDictionary(dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (*tokens.DictionaryToken, error) {
	token, ok := dict.TryGet(name)
	if !ok {
		return nil, fmt.Errorf("missing required key %v", name)
	}
	result, found := parts.TryGet[*tokens.DictionaryToken](token, scanner)
	if !found {
		return nil, fmt.Errorf("key %v is not a dictionary token", name)
	}
	return result, nil
}
