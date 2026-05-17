package functions

import (
	"errors"
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// PdfFunctionType3 represents a stitching function (PDF function type 3).
// It selects one of several child functions based on the input value,
// then evaluates that child function with an appropriately mapped input.
type PdfFunctionType3 struct {
	PdfFunctionBase

	functionsArray []PdfFunction
	bounds         *tokens.ArrayToken
	encode         *tokens.ArrayToken
	boundsValues   []float64
}

// NewPdfFunctionType3FromDict creates a PdfFunctionType3 backed by a dictionary token.
func NewPdfFunctionType3FromDict(dict *tokens.DictionaryToken, domain, rangeVals *tokens.ArrayToken, functionsArray []PdfFunction, bounds, encode *tokens.ArrayToken) (*PdfFunctionType3, error) {
	if len(functionsArray) == 0 {
		return nil, errors.New("functions array must not be empty")
	}

	boundsValues := extractDoubles(bounds.Data())

	return &PdfFunctionType3{
		PdfFunctionBase: *NewPdfFunctionBaseFromDict(dict, domain, rangeVals),
		functionsArray:  functionsArray,
		bounds:          bounds,
		encode:          encode,
		boundsValues:    boundsValues,
	}, nil
}

// NewPdfFunctionType3FromStream creates a PdfFunctionType3 backed by a stream token.
func NewPdfFunctionType3FromStream(stream *tokens.StreamToken, domain, rangeVals *tokens.ArrayToken, functionsArray []PdfFunction, bounds, encode *tokens.ArrayToken) (*PdfFunctionType3, error) {
	if len(functionsArray) == 0 {
		return nil, errors.New("functions array must not be empty")
	}

	boundsValues := extractDoubles(bounds.Data())

	return &PdfFunctionType3{
		PdfFunctionBase: *NewPdfFunctionBaseFromStream(stream, domain, rangeVals),
		functionsArray:  functionsArray,
		bounds:          bounds,
		encode:          encode,
		boundsValues:    boundsValues,
	}, nil
}

// FunctionType returns Stitching for this function type.
func (f *PdfFunctionType3) FunctionType() FunctionTypes {
	return Stitching
}

// Eval evaluates the stitching function at the given input values.
// The function selects a child function based on which partition interval
// contains the input value, maps the input through encode ranges, and returns
// the clipped result of the selected child function.
func (f *PdfFunctionType3) Eval(input ...float64) []float64 {
	x := input[0]
	domain := f.GetDomainForInput(0)
	x = ClipToRange(x, domain.Min(), domain.Max())

	var fn PdfFunction
	var xMapped float64

	if len(f.functionsArray) == 1 {
		fn = f.functionsArray[0]
		encRange := f.getEncodeForParameter(0)
		xMapped = interpolate(x, domain.Min(), domain.Max(), encRange.Min(), encRange.Max())
	} else {
		partitionSize := len(f.boundsValues) + 2
		partitionValues := make([]float64, partitionSize)
		partitionValues[0] = domain.Min()
		partitionValues[partitionSize-1] = domain.Max()
		copy(partitionValues[1:], f.boundsValues)

		found := false
		for i := 0; i < partitionSize-1; i++ {
			if x >= partitionValues[i] && (x < partitionValues[i+1] || (i == partitionSize-2 && x == partitionValues[i+1])) {
				fn = f.functionsArray[i]
				encRange := f.getEncodeForParameter(i)
				xMapped = interpolate(x, partitionValues[i], partitionValues[i+1], encRange.Min(), encRange.Max())
				found = true
				break
			}
		}

		if !found {
			panic(errors.New("partition not found in type 3 function"))
		}
	}

	result := fn.Eval(xMapped)
	return f.clipToRanges(result)
}

// FunctionsArray returns all child functions used by this stitching function.
func (f *PdfFunctionType3) FunctionsArray() []PdfFunction {
	return f.functionsArray
}

// Bounds returns the bounds array as an ArrayToken.
func (f *PdfFunctionType3) Bounds() *tokens.ArrayToken {
	return f.bounds
}

// Encode returns the encode array as an ArrayToken.
func (f *PdfFunctionType3) Encode() *tokens.ArrayToken {
	return f.encode
}

// getEncodeForParameter returns the encode range for the given parameter number.
func (f *PdfFunctionType3) getEncodeForParameter(n int) core.PdfRange {
	values := extractDoubles(f.encode.Data())
	return core.NewPdfRangeAtIndex(values, n)
}

// String returns a string representation of the stitching function.
func (f *PdfFunctionType3) String() string {
	return fmt.Sprintf("StitchingFunction{FunctionsCount: %d, Bounds: %v}", len(f.functionsArray), f.bounds)
}
