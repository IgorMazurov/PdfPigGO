package functions

import (
	"fmt"
	"math"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// PdfFunctionType2 represents an exponential interpolation function (PDF function type 2).
// It computes output values using the formula: c0 + x^N * (c1 - c0) for each component.
type PdfFunctionType2 struct {
	PdfFunctionBase

	// C0 is the array of base values for the exponential interpolation.
	C0 *tokens.ArrayToken

	// C1 is the array of target values for the exponential interpolation.
	C1 *tokens.ArrayToken

	// N is the exponent applied to the normalized input value.
	N float64
}

// NewPdfFunctionType2FromDict creates a PdfFunctionType2 backed by a dictionary token.
func NewPdfFunctionType2FromDict(dict *tokens.DictionaryToken, domain, rangeVals, c0, c1 *tokens.ArrayToken, n float64) *PdfFunctionType2 {
	return &PdfFunctionType2{
		PdfFunctionBase: *NewPdfFunctionBaseFromDict(dict, domain, rangeVals),
		C0:              c0,
		C1:              c1,
		N:               n,
	}
}

// NewPdfFunctionType2FromStream creates a PdfFunctionType2 backed by a stream token.
func NewPdfFunctionType2FromStream(stream *tokens.StreamToken, domain, rangeVals, c0, c1 *tokens.ArrayToken, n float64) *PdfFunctionType2 {
	return &PdfFunctionType2{
		PdfFunctionBase: *NewPdfFunctionBaseFromStream(stream, domain, rangeVals),
		C0:              c0,
		C1:              c1,
		N:               n,
	}
}

// FunctionType returns Exponential for this function type.
func (f *PdfFunctionType2) FunctionType() FunctionTypes {
	return Exponential
}

// Eval evaluates the exponential interpolation function at the given input values.
// The formula applied per component is: c0[j] + x^N * (c1[j] - c0[j]).
// Results are clipped to the function's range if one is specified.
func (f *PdfFunctionType2) Eval(input ...float64) []float64 {
	xToN := math.Pow(input[0], f.N)

	length := f.C0.Length()
	if f.C1.Length() < length {
		length = f.C1.Length()
	}

	result := make([]float64, length)
	c0Data := f.C0.Data()
	c1Data := f.C1.Data()

	for j := 0; j < length; j++ {
		c0j := c0Data[j].(*tokens.NumericToken).DoubleVal()
		c1j := c1Data[j].(*tokens.NumericToken).DoubleVal()
		result[j] = c0j + xToN*(c1j-c0j)
	}

	return f.clipToRanges(result)
}

// String returns a string representation of the function.
func (f *PdfFunctionType2) String() string {
	return fmt.Sprintf("FunctionType2{C0: %v C1: %v N: %g}", f.C0, f.C1, f.N)
}
