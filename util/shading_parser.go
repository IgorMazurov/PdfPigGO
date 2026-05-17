package util

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/functions"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CreateShading creates a Shading from the given shading token by resolving the
// underlying dictionary or stream and dispatching to the appropriate parser based
// on the /ShadingType entry.
func CreateShading(
	shadingToken tokens.Token,
	scanner tokenization.PdfTokenScanner,
	resourceStore ParserResourceStore,
	filterProvider PatternParserProvider,
) (colors.Shading, error) {
	var shadingDictionary *tokens.DictionaryToken
	var shadingStream *tokens.StreamToken

	if fs, ok := parts.TryGet[*tokens.StreamToken](shadingToken, scanner); ok && fs != nil {
		shadingDictionary = fs.StreamDictionary
		decodedData, err := decodeShadingStream(fs, filterProvider, scanner)
		if err != nil {
			return nil, fmt.Errorf("failed to decode shading stream: %w", err)
		}
		streamToken, err := tokens.NewStreamToken(fs.StreamDictionary, decodedData)
		if err != nil {
			return nil, fmt.Errorf("failed to create decoded stream token: %w", err)
		}
		shadingStream = streamToken
	} else if fd, ok := parts.TryGet[*tokens.DictionaryToken](shadingToken, scanner); ok && fd != nil {
		shadingDictionary = fd
	} else {
		return nil, core.NewPdfDocumentFormatException(
			fmt.Sprintf("invalid Shading token encountered in page resource dictionary: %v", shadingToken))
	}

	if shadingDictionary == nil {
		return nil, core.NewPdfDocumentFormatException("shading dictionary is null")
	}

	shadingTypeVal, ok := shadingDictionary.TryGet(tokens.ShadingType)
	if !ok {
		return nil, fmt.Errorf("'%v' is required for shading", tokens.ShadingType)
	}

	shadingTypeToken, err := parts.GetByToken[*tokens.NumericToken](shadingTypeVal, scanner)
	if err != nil {
		return nil, err
	}

	shadingTypeInt := shadingTypeToken.IntVal()
	shadingType := colors.ShadingType(shadingTypeInt)

	if shadingTypeInt >= 4 && shadingStream == nil {
		return nil, fmt.Errorf(
			"shading type '%v' is not properly defined. Shading types 4 to 7 shall be defined by a stream.",
			shadingType)
	}

	colorSpace := resolveColorSpace(shadingDictionary, scanner, resourceStore)
	if colorSpace == nil {
		return nil, fmt.Errorf("color space resolution failed for shading")
	}

	background := getOptionalNumericArray(shadingDictionary, tokens.Background, scanner)

	var bBox *core.PdfRectangle
	if bboxToken, found := TryGetOptionalTokenDirect[*tokens.ArrayToken](shadingDictionary, tokens.Bbox, scanner); found {
		rect, err := ToRectangle(bboxToken, scanner)
		if err == nil && rect != nil {
			bBox = rect
		}
	}

	antiAlias := false
	if antiAliasToken, found := TryGetOptionalTokenDirect[*tokens.BooleanToken](shadingDictionary, tokens.AntiAlias, scanner); found {
		antiAlias = antiAliasToken.Data()
	}

	switch shadingType {
	case colors.FunctionBased:
		return createFunctionBasedShading(shadingDictionary, colorSpace, background, bBox, antiAlias, scanner, filterProvider)

	case colors.Axial:
		return createAxialShading(shadingDictionary, colorSpace, background, bBox, antiAlias, scanner, filterProvider)

	case colors.Radial:
		return createRadialShading(shadingDictionary, colorSpace, background, bBox, antiAlias, scanner, filterProvider)

	case colors.FreeFormGouraud:
		if shadingStream == nil {
			return nil, core.NewPdfDocumentFormatException("shading type 4 requires a stream")
		}
		return createFreeFormGouraudShading(shadingStream, colorSpace, background, bBox, antiAlias, scanner, filterProvider)

	case colors.LatticeFormGouraud:
		if shadingStream == nil {
			return nil, core.NewPdfDocumentFormatException("shading type 5 requires a stream")
		}
		return createLatticeFormGouraudShading(shadingStream, colorSpace, background, bBox, antiAlias, scanner, filterProvider)

	case colors.CoonsPatch:
		if shadingStream == nil {
			return nil, core.NewPdfDocumentFormatException("shading type 6 requires a stream")
		}
		return createCoonsPatchMeshesShading(shadingStream, colorSpace, background, bBox, antiAlias, scanner, filterProvider)

	case colors.TensorProductPatch:
		if shadingStream == nil {
			return nil, core.NewPdfDocumentFormatException("shading type 7 requires a stream")
		}
		return createTensorProductPatchMeshesShading(shadingStream, colorSpace, background, bBox, antiAlias, scanner, filterProvider)

	default:
		return nil, core.NewPdfDocumentFormatException(
			fmt.Sprintf("invalid Shading type encountered in page resource dictionary: '%v'.", shadingType))
	}
}

// resolveColorSpace resolves the color space from either /CS name or array entry.
func resolveColorSpace(
	dict *tokens.DictionaryToken,
	scanner tokenization.PdfTokenScanner,
	resourceStore ParserResourceStore,
) colors.ColorSpaceDetails {
	csVal, hasCS := dict.TryGet(tokens.ColorSpace)
	if !hasCS {
		return nil
	}

	if nameToken, err := parts.GetByToken[*tokens.NameToken](csVal, scanner); err == nil && nameToken != nil {
		return resourceStore.GetColorSpaceDetails(nameToken, dict)
	}

	if arrayToken, ok := parts.TryGet[*tokens.ArrayToken](csVal, scanner); ok && arrayToken != nil {
		data := arrayToken.Data()
		if len(data) > 0 {
			if firstColorSpaceName, ok := data[0].(*tokens.NameToken); ok && firstColorSpaceName != nil {
				return resourceStore.GetColorSpaceDetails(firstColorSpaceName, dict)
			}
		}
		return nil
	}

	return nil
}

// decodeShadingStream decodes the raw stream data by applying all registered filters.
func decodeShadingStream(
	stream *tokens.StreamToken,
	provider PatternParserProvider,
	scanner tokenization.PdfTokenScanner,
) ([]byte, error) {
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

// getOptionalNumericArray extracts numeric values from an optional array entry.
func getOptionalNumericArray(
	dict *tokens.DictionaryToken,
	name *tokens.NameToken,
	scanner tokenization.PdfTokenScanner,
) []float64 {
	arr, ok := TryGetOptionalTokenDirect[*tokens.ArrayToken](dict, name, scanner)
	if !ok || arr == nil {
		return nil
	}

	data := arr.Data()
	result := make([]float64, 0, len(data))
	for _, elem := range data {
		if num, ok := elem.(*tokens.NumericToken); ok && num != nil {
			result = append(result, num.DoubleVal())
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

// getRequiredNumericArray extracts numeric values from a required array entry.
func getRequiredNumericArray(
	dict *tokens.DictionaryToken,
	name *tokens.NameToken,
	scanner tokenization.PdfTokenScanner,
	shadingType colors.ShadingType,
) ([]float64, error) {
	arr, ok := TryGetOptionalTokenDirect[*tokens.ArrayToken](dict, name, scanner)
	if !ok || arr == nil {
		return nil, fmt.Errorf("%v is required for shading type '%v'.", name, shadingType)
	}

	data := arr.Data()
	result := make([]float64, 0, len(data))
	for _, elem := range data {
		if num, ok := elem.(*tokens.NumericToken); ok && num != nil {
			result = append(result, num.DoubleVal())
		}
	}

	return result, nil
}

// getFunctions resolves one or more PdfFunction(s) from a function token.
func getFunctions(
	functionToken tokens.Token,
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
) ([]functions.PdfFunction, error) {
	if fa, ok := parts.TryGet[*tokens.ArrayToken](functionToken, scanner); ok && fa != nil {
		data := fa.Data()
		functionArray := make([]functions.PdfFunction, len(data))
		for i := 0; i < len(data); i++ {
			fn, err := Create(data[i], scanner, filterProvider)
			if err != nil {
				return nil, fmt.Errorf("failed to create function at index %d: %w", i, err)
			}
			functionArray[i] = fn
		}
		return functionArray, nil
	}

	fn, err := Create(functionToken, scanner, filterProvider)
	if err != nil {
		return nil, err
	}
	return []functions.PdfFunction{fn}, nil
}

// createFunctionBasedShading creates a type 1 (function-based) shading.
func createFunctionBasedShading(
	shadingDictionary *tokens.DictionaryToken,
	colorSpace colors.ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
) (colors.Shading, error) {
	var domain []float64
	if domArr, found := TryGetOptionalTokenDirect[*tokens.ArrayToken](shadingDictionary, tokens.Domain, scanner); found && domArr != nil {
		domain = extractNumericSlice(domArr.Data())
	}
	if domain == nil || len(domain) == 0 {
		domain = []float64{0.0, 1.0, 0.0, 1.0}
	}

	matrix := core.NewTransformationMatrixFromSlice([]float64{1, 0, 0, 1, 0, 0})
	if matArr, found := TryGetOptionalTokenDirect[*tokens.ArrayToken](shadingDictionary, tokens.Matrix, scanner); found && matArr != nil {
		nums := extractNumericSlice(matArr.Data())
		if len(nums) >= 6 {
			matrix = core.NewTransformationMatrixFromSlice(nums[:6])
		}
	}

	funcVal, hasFunc := shadingDictionary.TryGet(tokens.Function)
	if !hasFunc {
		return nil, fmt.Errorf("'%v' is required for shading type '%v'.", tokens.Function, colors.FunctionBased)
	}

	functions, err := getFunctions(funcVal, scanner, filterProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve function(s): %w", err)
	}

	return colors.NewFunctionBasedShading(
		shadingDictionary, colorSpace, background, bBox, antiAlias, domain, matrix, functions), nil
}

// createAxialShading creates a type 2 (axial) shading.
func createAxialShading(
	shadingDictionary *tokens.DictionaryToken,
	colorSpace colors.ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
) (colors.Shading, error) {
	coords, err := getRequiredNumericArray(shadingDictionary, tokens.Coords, scanner, colors.Axial)
	if err != nil {
		return nil, err
	}

	var domain []float64
	if domArr, found := TryGetOptionalTokenDirect[*tokens.ArrayToken](shadingDictionary, tokens.Domain, scanner); found && domArr != nil {
		domain = extractNumericSlice(domArr.Data())
	}
	if domain == nil || len(domain) == 0 {
		domain = []float64{0, 1}
	}

	funcVal, hasFunc := shadingDictionary.TryGet(tokens.Function)
	if !hasFunc {
		return nil, fmt.Errorf("%v is required for shading type '%v'.", tokens.Function, colors.Axial)
	}

	functions, err := getFunctions(funcVal, scanner, filterProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve function(s): %w", err)
	}

	extend := []bool{false, false}
	if extArr, found := TryGetOptionalTokenDirect[*tokens.ArrayToken](shadingDictionary, tokens.Extend, scanner); found && extArr != nil {
		ext := extractBoolSlice(extArr.Data())
		if len(ext) > 0 {
			extend = ext
		}
	}

	return colors.NewAxialShading(
		shadingDictionary, colorSpace, background, bBox, antiAlias, coords, domain, functions, extend), nil
}

// createRadialShading creates a type 3 (radial) shading.
func createRadialShading(
	shadingDictionary *tokens.DictionaryToken,
	colorSpace colors.ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
) (colors.Shading, error) {
	coords, err := getRequiredNumericArray(shadingDictionary, tokens.Coords, scanner, colors.Radial)
	if err != nil {
		return nil, err
	}

	var domain []float64
	if domArr, found := TryGetOptionalTokenDirect[*tokens.ArrayToken](shadingDictionary, tokens.Domain, scanner); found && domArr != nil {
		domain = extractNumericSlice(domArr.Data())
	}
	if domain == nil || len(domain) == 0 {
		domain = []float64{0, 1}
	}

	funcVal, hasFunc := shadingDictionary.TryGet(tokens.Function)
	if !hasFunc {
		return nil, fmt.Errorf("%v is required for shading type '%v'.", tokens.Function, colors.Radial)
	}

	functions, err := getFunctions(funcVal, scanner, filterProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve function(s): %w", err)
	}

	extend := []bool{false, false}
	if extArr, found := TryGetOptionalTokenDirect[*tokens.ArrayToken](shadingDictionary, tokens.Extend, scanner); found && extArr != nil {
		ext := extractBoolSlice(extArr.Data())
		if len(ext) > 0 {
			extend = ext
		}
	}

	return colors.NewRadialShading(
		shadingDictionary, colorSpace, background, bBox, antiAlias, coords, domain, functions, extend), nil
}

// createFreeFormGouraudShading creates a type 4 (free-form Gouraud-shaded triangle mesh) shading.
func createFreeFormGouraudShading(
	shadingStream *tokens.StreamToken,
	colorSpace colors.ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
) (colors.Shading, error) {
	dict := shadingStream.StreamDictionary

	bitsPerCoordinate, err := getRequiredInt(dict, tokens.BitsPerCoordinate, scanner, colors.FreeFormGouraud)
	if err != nil {
		return nil, err
	}

	bitsPerComponent, err := getRequiredInt(dict, tokens.BitsPerComponent, scanner, colors.FreeFormGouraud)
	if err != nil {
		return nil, err
	}

	bitsPerFlag, err := getRequiredInt(dict, tokens.BitsPerFlag, scanner, colors.FreeFormGouraud)
	if err != nil {
		return nil, err
	}

	decode, err := getRequiredNumericArray(dict, tokens.Decode, scanner, colors.FreeFormGouraud)
	if err != nil {
		return nil, err
	}

	var fnArray []functions.PdfFunction
	if funcVal, hasFunc := dict.TryGet(tokens.Function); hasFunc {
		fnArray, err = getFunctions(funcVal, scanner, filterProvider)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve function(s): %w", err)
		}
	}

	return colors.NewFreeFormGouraudShading(
		shadingStream, colorSpace, background, bBox, antiAlias,
		bitsPerCoordinate, bitsPerComponent, bitsPerFlag, decode, fnArray), nil
}

// createLatticeFormGouraudShading creates a type 5 (lattice-form Gouraud-shaded triangle mesh) shading.
func createLatticeFormGouraudShading(
	shadingStream *tokens.StreamToken,
	colorSpace colors.ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
) (colors.Shading, error) {
	dict := shadingStream.StreamDictionary

	bitsPerCoordinate, err := getRequiredInt(dict, tokens.BitsPerCoordinate, scanner, colors.LatticeFormGouraud)
	if err != nil {
		return nil, err
	}

	bitsPerComponent, err := getRequiredInt(dict, tokens.BitsPerComponent, scanner, colors.LatticeFormGouraud)
	if err != nil {
		return nil, err
	}

	verticesPerRow, err := getRequiredInt(dict, tokens.VerticesPerRow, scanner, colors.LatticeFormGouraud)
	if err != nil {
		return nil, err
	}

	decode, err := getRequiredNumericArray(dict, tokens.Decode, scanner, colors.LatticeFormGouraud)
	if err != nil {
		return nil, err
	}

	var fnArray []functions.PdfFunction
	if funcVal, hasFunc := dict.TryGet(tokens.Function); hasFunc {
		fnArray, err = getFunctions(funcVal, scanner, filterProvider)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve function(s): %w", err)
		}
	}

	return colors.NewLatticeFormGouraudShading(
		shadingStream, colorSpace, background, bBox, antiAlias,
		bitsPerCoordinate, bitsPerComponent, verticesPerRow, decode, fnArray), nil
}

// createCoonsPatchMeshesShading creates a type 6 (Coons patch mesh) shading.
func createCoonsPatchMeshesShading(
	shadingStream *tokens.StreamToken,
	colorSpace colors.ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
) (colors.Shading, error) {
	dict := shadingStream.StreamDictionary

	bitsPerCoordinate, err := getRequiredInt(dict, tokens.BitsPerCoordinate, scanner, colors.CoonsPatch)
	if err != nil {
		return nil, err
	}

	bitsPerComponent, err := getRequiredInt(dict, tokens.BitsPerComponent, scanner, colors.CoonsPatch)
	if err != nil {
		return nil, err
	}

	bitsPerFlag, err := getRequiredInt(dict, tokens.BitsPerFlag, scanner, colors.CoonsPatch)
	if err != nil {
		return nil, err
	}

	decode, err := getRequiredNumericArray(dict, tokens.Decode, scanner, colors.CoonsPatch)
	if err != nil {
		return nil, err
	}

	var fnArray []functions.PdfFunction
	if funcVal, hasFunc := dict.TryGet(tokens.Function); hasFunc {
		fnArray, err = getFunctions(funcVal, scanner, filterProvider)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve function(s): %w", err)
		}
	}

	return colors.NewCoonsPatchMeshesShading(
		shadingStream, colorSpace, background, bBox, antiAlias,
		bitsPerCoordinate, bitsPerComponent, bitsPerFlag, decode, fnArray), nil
}

// createTensorProductPatchMeshesShading creates a type 7 (tensor-product patch mesh) shading.
func createTensorProductPatchMeshesShading(
	shadingStream *tokens.StreamToken,
	colorSpace colors.ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
) (colors.Shading, error) {
	dict := shadingStream.StreamDictionary

	bitsPerCoordinate, err := getRequiredInt(dict, tokens.BitsPerCoordinate, scanner, colors.TensorProductPatch)
	if err != nil {
		return nil, err
	}

	bitsPerComponent, err := getRequiredInt(dict, tokens.BitsPerComponent, scanner, colors.TensorProductPatch)
	if err != nil {
		return nil, err
	}

	bitsPerFlag, err := getRequiredInt(dict, tokens.BitsPerFlag, scanner, colors.TensorProductPatch)
	if err != nil {
		return nil, err
	}

	decode, err := getRequiredNumericArray(dict, tokens.Decode, scanner, colors.TensorProductPatch)
	if err != nil {
		return nil, err
	}

	var fnArray []functions.PdfFunction
	if funcVal, hasFunc := dict.TryGet(tokens.Function); hasFunc {
		fnArray, err = getFunctions(funcVal, scanner, filterProvider)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve function(s): %w", err)
		}
	}

	return colors.NewTensorProductPatchMeshesShading(
		shadingStream, colorSpace, background, bBox, antiAlias,
		bitsPerCoordinate, bitsPerComponent, bitsPerFlag, decode, fnArray), nil
}

// getRequiredInt retrieves a required integer value from the dictionary by key name.
func getRequiredInt(
	dict *tokens.DictionaryToken,
	name *tokens.NameToken,
	scanner tokenization.PdfTokenScanner,
	shadingType colors.ShadingType,
) (int, error) {
	token, found := dict.TryGet(name)
	if !found {
		return 0, fmt.Errorf("%v is required for shading type '%v'.", name, shadingType)
	}

	num, err := parts.GetByToken[*tokens.NumericToken](token, scanner)
	if err != nil || num == nil {
		return 0, fmt.Errorf("%v is not a numeric token in shading type '%v'.", name, shadingType)
	}

	return num.IntVal(), nil
}

// extractNumericSlice extracts float64 values from a slice of tokens.
func extractNumericSlice(data []tokens.Token) []float64 {
	result := make([]float64, 0, len(data))
	for _, elem := range data {
		if num, ok := elem.(*tokens.NumericToken); ok && num != nil {
			result = append(result, num.DoubleVal())
		}
	}
	return result
}

// extractBoolSlice extracts bool values from a slice of tokens.
func extractBoolSlice(data []tokens.Token) []bool {
	result := make([]bool, 0, len(data))
	for _, elem := range data {
		if b, ok := elem.(*tokens.BooleanToken); ok && b != nil {
			result = append(result, b.Data())
		}
	}
	return result
}
