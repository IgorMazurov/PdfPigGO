package functions

import (
	"math"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// PdfFunctionType0 represents a sampled function (PDF function type 0).
// It performs linear or cubic spline interpolation based on a set of sample points.
type PdfFunctionType0 struct {
	PdfFunctionBase

	size          *tokens.ArrayToken
	bitsPerSample int
	order         int
	encodeValues  *tokens.ArrayToken
	decodeValues  *tokens.ArrayToken
	samples       [][]int
}

// NewPdfFunctionType0FromDict creates a PdfFunctionType0 backed by a dictionary token.
func NewPdfFunctionType0FromDict(dict *tokens.DictionaryToken, domain, rangeVals, size *tokens.ArrayToken, bitsPerSample, order int, encode, decode *tokens.ArrayToken) *PdfFunctionType0 {
	return &PdfFunctionType0{
		PdfFunctionBase: *NewPdfFunctionBaseFromDict(dict, domain, rangeVals),
		size:            size,
		bitsPerSample:   bitsPerSample,
		order:           order,
		encodeValues:    encode,
		decodeValues:    decode,
	}
}

// NewPdfFunctionType0FromStream creates a PdfFunctionType0 backed by a stream token.
func NewPdfFunctionType0FromStream(stream *tokens.StreamToken, domain, rangeVals, size *tokens.ArrayToken, bitsPerSample, order int, encode, decode *tokens.ArrayToken) *PdfFunctionType0 {
	return &PdfFunctionType0{
		PdfFunctionBase: *NewPdfFunctionBaseFromStream(stream, domain, rangeVals),
		size:            size,
		bitsPerSample:   bitsPerSample,
		order:           order,
		encodeValues:    encode,
		decodeValues:    decode,
	}
}

// FunctionType returns Sampled for this function type.
func (f *PdfFunctionType0) FunctionType() FunctionTypes {
	return Sampled
}

// Size returns the number of samples in each input dimension of the sample table.
// An array of m positive integers specifying the number of samples in each input dimension.
func (f *PdfFunctionType0) Size() *tokens.ArrayToken {
	return f.size
}

// BitsPerSample returns the number of bits for each output value.
// Valid values are 1, 2, 4, 8, 12, 16, 24, 32.
func (f *PdfFunctionType0) BitsPerSample() int {
	return f.bitsPerSample
}

// Order returns the interpolation order between samples.
// Valid values are 1 (linear) and 3 (cubic spline). Default is 1.
func (f *PdfFunctionType0) Order() int {
	return f.order
}

// getEncodeForParameter returns the encode range for the given parameter number.
// Returns nil if no encode value exists at that index.
func (f *PdfFunctionType0) getEncodeForParameter(paramNum int) *core.PdfRange {
	if f.encodeValues == nil || f.encodeValues.Length() < paramNum*2+1 {
		return nil
	}
	values := extractDoubles(f.encodeValues.Data())
	r := core.NewPdfRangeAtIndex(values, paramNum)
	return &r
}

// getDecodeForParameter returns the decode range for the given parameter number.
// Returns nil if no decode value exists at that index.
func (f *PdfFunctionType0) getDecodeForParameter(paramNum int) *core.PdfRange {
	if f.decodeValues == nil || f.decodeValues.Length() < paramNum*2+1 {
		return nil
	}
	values := extractDoubles(f.decodeValues.Data())
	r := core.NewPdfRangeAtIndex(values, paramNum)
	return &r
}

// Eval evaluates the sampled function at the given input values using interpolation.
// Input values are mapped through domain and encode ranges, interpolated against sample data,
// then mapped through decode and range outputs.
func (f *PdfFunctionType0) Eval(input ...float64) []float64 {
	sizeValues := extractDoubles(f.size.Data())
	bitsPerSample := f.bitsPerSample
	maxSample := math.Pow(2, float64(bitsPerSample)) - 1.0
	numInput := len(input)
	numOutput := f.NumberOfOutputParameters()

	inputPrev := make([]int, numInput)
	inputNext := make([]int, numInput)
	inputCopy := make([]float64, numInput)
	copy(inputCopy, input)

	for i := 0; i < numInput; i++ {
		domain := f.GetDomainForInput(i)
		encodeRange := f.getEncodeForParameter(i)
		inputCopy[i] = ClipToRange(inputCopy[i], domain.Min(), domain.Max())

		if encodeRange != nil {
			inputCopy[i] = interpolate(inputCopy[i], domain.Min(), domain.Max(), encodeRange.Min(), encodeRange.Max())
		} else {
			inputCopy[i] = interpolate(inputCopy[i], domain.Min(), domain.Max(), 0, sizeValues[i]-1)
		}

		inputCopy[i] = ClipToRange(inputCopy[i], 0, sizeValues[i]-1)
		inputPrev[i] = int(math.Floor(inputCopy[i]))
		inputNext[i] = int(math.Ceil(inputCopy[i]))
	}

	outputValues := newRInterpol(inputCopy, inputPrev, inputNext, numOutput, f.size, f.getSamples()).rInterpolate()

	for i := 0; i < numOutput; i++ {
		rng := f.GetRangeForOutput(i)
		decodeRange := f.getDecodeForParameter(i)
		if decodeRange == nil {
			dr := core.NewPdfRange([]float64{rng.Min(), rng.Max()})
			decodeRange = &dr
		}
		outputValues[i] = interpolate(outputValues[i], 0, maxSample, decodeRange.Min(), decodeRange.Max())
		outputValues[i] = ClipToRange(outputValues[i], rng.Min(), rng.Max())
	}

	return outputValues
}

// rInterpol performs recursive N-dimensional interpolation.
type rInterpol struct {
	input             []float64
	inputPrev         []int
	inputNext         []int
	numOutputValues   int
	size              *tokens.ArrayToken
	samples           [][]int
}

func newRInterpol(input []float64, inputPrev, inputNext []int, numOutputValues int, size *tokens.ArrayToken, samples [][]int) *rInterpol {
	return &rInterpol{
		input:           input,
		inputPrev:       inputPrev,
		inputNext:       inputNext,
		numOutputValues: numOutputValues,
		size:            size,
		samples:         samples,
	}
}

// rInterpolate calculates the interpolated result sample.
func (r *rInterpol) rInterpolate() []float64 {
	return r.internalRInterpol(make([]int, len(r.input)), 0)
}

// internalRInterpol recursively interpolates across dimensions.
func (r *rInterpol) internalRInterpol(coord []int, step int) []float64 {
	resultSample := make([]float64, r.numOutputValues)

	if step == len(r.input)-1 {
		if r.inputPrev[step] == r.inputNext[step] {
			coord[step] = r.inputPrev[step]
			tmpSample := r.samples[r.calcSampleIndex(coord)]
			for i := 0; i < r.numOutputValues; i++ {
				resultSample[i] = float64(tmpSample[i])
			}
			return resultSample
		}

		coord[step] = r.inputPrev[step]
		sample1 := r.samples[r.calcSampleIndex(coord)]
		coord[step] = r.inputNext[step]
		sample2 := r.samples[r.calcSampleIndex(coord)]

		for i := 0; i < r.numOutputValues; i++ {
			resultSample[i] = interpolate(r.input[step], float64(r.inputPrev[step]), float64(r.inputNext[step]), float64(sample1[i]), float64(sample2[i]))
		}
		return resultSample
	}

	if r.inputPrev[step] == r.inputNext[step] {
		coord[step] = r.inputPrev[step]
		return r.internalRInterpol(coord, step+1)
	}

	coord[step] = r.inputPrev[step]
	sample1 := r.internalRInterpol(coord, step+1)
	coord[step] = r.inputNext[step]
	sample2 := r.internalRInterpol(coord, step+1)

	for i := 0; i < r.numOutputValues; i++ {
		resultSample[i] = interpolate(r.input[step], float64(r.inputPrev[step]), float64(r.inputNext[step]), sample1[i], sample2[i])
	}

	return resultSample
}

// calcSampleIndex calculates the flat array index from multi-dimensional coordinates.
func (r *rInterpol) calcSampleIndex(vector []int) int {
	sizeValues := extractDoubles(r.size.Data())
	dimension := len(vector)

	index := 0
	sizeProduct := 1

	for i := dimension - 2; i >= 0; i-- {
		sizeProduct = int(float64(sizeProduct) * sizeValues[i])
	}

	for i := dimension - 1; i >= 0; i-- {
		index += sizeProduct * vector[i]
		if i-1 >= 0 {
			sizeProduct = int(float64(sizeProduct) / sizeValues[i-1])
		}
	}

	return index
}

// getSamples lazily parses the raw bit stream into a 2D sample array.
func (f *PdfFunctionType0) getSamples() [][]int {
	if f.samples != nil {
		return f.samples
	}

	arraySize := 1
	nIn := f.NumberOfInputParameters()
	nOut := f.NumberOfOutputParameters()

	for i := 0; i < nIn; i++ {
		arraySize *= f.size.Get(i).(*tokens.NumericToken).IntVal()
	}

	f.samples = make([][]int, arraySize)
	bitsPerSample := f.bitsPerSample
	streamData := f.FunctionStream().Data()

	totalBits := len(streamData) * 8
	bits := make([]bool, totalBits)
	for i := 0; i < len(streamData); i++ {
		for b := 0; b < 8; b++ {
			bits[i*8+b] = (streamData[i]>>uint(b))&1 == 1
		}
	}

	for i := 0; i < arraySize; i++ {
		f.samples[i] = make([]int, nOut)
		for k := 0; k < nOut; k++ {
			var accum int64 = 0
			for l := bitsPerSample - 1; l >= 0; l-- {
				accum <<= 1
				if bits[i*nOut*bitsPerSample+(k*bitsPerSample)+l] {
					accum |= 1
				}
			}
			f.samples[i][k] = int(accum)
		}
	}

	return f.samples
}
