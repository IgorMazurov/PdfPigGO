package core

// PdfRange represents a numeric range where Min <= value <= Max.
// The underlying array stores pairs of values, and StartingIndex selects
// which pair to interpret as the current range.
type PdfRange struct {
	rangeArray    []float64
	startingIndex int
}

// NewPdfRange creates a new PdfRange from a slice of doubles,
// assuming a starting index of 0.
func NewPdfRange(rangeValues []float64) PdfRange {
	return NewPdfRangeAtIndex(rangeValues, 0)
}

// NewPdfRangeAtIndex creates a new PdfRange with an explicit starting index
// into the array. Some arrays specify multiple ranges (e.g., [0, 1, 0, 2, 2, 3]).
// The startingIndex selects which pair to use: Min = rangeValues[startingIndex*2],
// Max = rangeValues[startingIndex*2+1].
func NewPdfRangeAtIndex(rangeValues []float64, index int) PdfRange {
	arr := make([]float64, len(rangeValues))
	copy(arr, rangeValues)
	return PdfRange{rangeArray: arr, startingIndex: index}
}

// Min returns the minimum value of the range.
func (r PdfRange) Min() float64 {
	return r.rangeArray[r.startingIndex*2]
}

// Max returns the maximum value of the range.
func (r PdfRange) Max() float64 {
	return r.rangeArray[r.startingIndex*2+1]
}
