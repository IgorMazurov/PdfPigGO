package testutil

import "math"

// DoubleComparer provides approximate equality comparison for float64 values.
type DoubleComparer struct {
	precision float64
}

// NewDoubleComparer creates a new DoubleComparer with the given precision.
func NewDoubleComparer(precision float64) *DoubleComparer {
	return &DoubleComparer{precision: precision}
}

// Equals returns true if x and y are approximately equal within the configured precision.
func (c *DoubleComparer) Equals(x, y float64) bool {
	return math.Abs(x-y) < c.precision
}
