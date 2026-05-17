package core

import (
	"errors"
	"fmt"
	"strings"
)

// LineDashPattern controls the pattern of dashes and gaps used to stroke paths.
// It is specified by a dash array and a dash phase.
type LineDashPattern struct {
	// Phase is the distance into the dash pattern at which to start the dash.
	Phase int

	// Array specifies the lengths of alternating dashes and gaps.
	Array []float64
}

// Solid is the default solid line with no dashes.
var Solid = LineDashPattern{Phase: 0, Array: []float64{}}

// NewLineDashPattern creates a new LineDashPattern with the given phase and dash array.
// Returns an error if array is nil.
func NewLineDashPattern(phase int, array []float64) (LineDashPattern, error) {
	if array == nil {
		return LineDashPattern{}, errors.New("array cannot be nil")
	}

	return LineDashPattern{Phase: phase, Array: array}, nil
}

// String returns the string representation of the line dash pattern.
func (l LineDashPattern) String() string {
	parts := make([]string, len(l.Array))

	for i, v := range l.Array {
		parts[i] = fmt.Sprintf("%.2f", v)
	}

	return fmt.Sprintf("[%s] %d.", strings.Join(parts, " "), l.Phase)
}

var _ fmt.Stringer = LineDashPattern{}
