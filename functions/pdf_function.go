package functions

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// FunctionTypes defines the type of a PDF function as specified by the PDF specification.
// Possible values are: Sampled (0), Exponential (2), Stitching (3), PostScript (4).
type FunctionTypes int

const (
	// Sampled represents a sampled function (function type 0).
	Sampled FunctionTypes = iota

	// Exponential represents an exponential interpolation/smoothness function (function type 2).
	Exponential FunctionTypes = 2

	// Stitching represents a stitching function (function type 3).
	Stitching FunctionTypes = 3

	// PostScript represents a PostScript calculator function (function type 4).
	PostScript FunctionTypes = 4
)

// PdfFunction represents a function object in a PDF document.
type PdfFunction interface {
	// FunctionType returns one of Sampled, Exponential, Stitching, or PostScript.
	FunctionType() FunctionTypes

	// Eval evaluates the function at the given input values and returns the output values.
	Eval(input ...float64) []float64
}

// PdfFunctionBase provides shared state and helper methods for all PDF function implementations.
// Concrete function types should embed this struct to reuse domain/range handling logic.
type PdfFunctionBase struct {
	functionDict *tokens.DictionaryToken
	functionStream *tokens.StreamToken
	domainValues   *tokens.ArrayToken
	rangeValues    *tokens.ArrayToken

	numberOfInputValues  int
	numberOfOutputValues int
}

// NewPdfFunctionBaseFromDict creates a PdfFunctionBase backed by a dictionary token.
func NewPdfFunctionBaseFromDict(dict *tokens.DictionaryToken, domain *tokens.ArrayToken, rangeVals *tokens.ArrayToken) *PdfFunctionBase {
	return &PdfFunctionBase{
		functionDict:         dict,
		domainValues:         domain,
		rangeValues:          rangeVals,
		numberOfInputValues:  -1,
		numberOfOutputValues: -1,
	}
}

// NewPdfFunctionBaseFromStream creates a PdfFunctionBase backed by a stream token.
func NewPdfFunctionBaseFromStream(stream *tokens.StreamToken, domain *tokens.ArrayToken, rangeVals *tokens.ArrayToken) *PdfFunctionBase {
	return &PdfFunctionBase{
		functionStream:       stream,
		domainValues:         domain,
		rangeValues:          rangeVals,
		numberOfInputValues:  -1,
		numberOfOutputValues: -1,
	}
}

// FunctionDictionary returns the function dictionary if one was set during construction.
func (f *PdfFunctionBase) FunctionDictionary() *tokens.DictionaryToken {
	return f.functionDict
}

// FunctionStream returns the function stream if one was set during construction.
func (f *PdfFunctionBase) FunctionStream() *tokens.StreamToken {
	return f.functionStream
}

// GetDictionary returns the function's dictionary. If a stream is used, it returns
// the stream's StreamDictionary; otherwise it returns the function dictionary directly.
func (f *PdfFunctionBase) GetDictionary() *tokens.DictionaryToken {
	if f.functionStream != nil {
		return f.functionStream.StreamDictionary
	}
	return f.functionDict
}

// RangeValues returns all ranges for the output values as an ArrayToken.
// Required for type 0 and type 4 functions. May be nil.
func (f *PdfFunctionBase) RangeValues() *tokens.ArrayToken {
	return f.rangeValues
}

// NumberOfOutputParameters returns the number of output parameters that have a range specified.
// A range for output parameters is optional so this may return zero for a function
// that does have output parameters but no range specified.
func (f *PdfFunctionBase) NumberOfOutputParameters() int {
	if f.numberOfOutputValues == -1 {
		if f.rangeValues == nil {
			f.numberOfOutputValues = 0
		} else {
			f.numberOfOutputValues = f.rangeValues.Length() / 2
		}
	}
	return f.numberOfOutputValues
}

// GetRangeForOutput returns the range for the nth output parameter.
func (f *PdfFunctionBase) GetRangeForOutput(n int) core.PdfRange {
	values := extractDoubles(f.rangeValues.Data())
	return core.NewPdfRangeAtIndex(values, n)
}

// NumberOfInputParameters returns the number of input parameters that have a domain specified.
func (f *PdfFunctionBase) NumberOfInputParameters() int {
	if f.numberOfInputValues == -1 {
		f.numberOfInputValues = f.domainValues.Length() / 2
	}
	return f.numberOfInputValues
}

// GetDomainForInput returns the domain range for the nth input parameter.
func (f *PdfFunctionBase) GetDomainForInput(n int) core.PdfRange {
	values := extractDoubles(f.domainValues.Data())
	return core.NewPdfRangeAtIndex(values, n)
}

// clipToRanges clips each value in inputValues to its corresponding range from RangeValues.
func (f *PdfFunctionBase) clipToRanges(inputValues []float64) []float64 {
	if f.rangeValues == nil || f.rangeValues.Length() == 0 {
		return inputValues
	}

	rangeVals := extractDoubles(f.rangeValues.Data())
	numRanges := len(rangeVals) / 2
	result := make([]float64, numRanges)

	for i := 0; i < numRanges; i++ {
		idx := i << 1
		result[i] = ClipToRange(inputValues[i], rangeVals[idx], rangeVals[idx+1])
	}

	return result
}

// ClipToRange clips a single value to the given [minVal, maxVal] range.
func ClipToRange(x, minVal, maxVal float64) float64 {
	if x < minVal {
		return minVal
	}
	if x > maxVal {
		return maxVal
	}
	return x
}

// interpolate calculates the y value on the line defined by two point pairs:
// (xMin, yMin) and (xMax, yMax), for a given x coordinate.
func interpolate(x, xMin, xMax, yMin, yMax float64) float64 {
	return yMin + (x-xMin)*(yMax-yMin)/(xMax-xMin)
}

// extractDoubles extracts float64 values from NumericToken elements in a token slice.
func extractDoubles(data []tokens.Token) []float64 {
	result := make([]float64, 0, len(data))
	for _, t := range data {
		if nt, ok := t.(*tokens.NumericToken); ok {
			result = append(result, nt.DoubleVal())
		}
	}
	return result
}
