package testutil

import "github.com/uglytoad/pdfpig/go/core"

// PointComparer provides equality comparison for PdfPoint values
// using an underlying DoubleComparer for coordinate-wise comparison.
type PointComparer struct {
	doubleComparer *DoubleComparer
}

// NewPointComparer creates a new PointComparer with the given double comparer.
func NewPointComparer(doubleComparer *DoubleComparer) *PointComparer {
	return &PointComparer{doubleComparer: doubleComparer}
}

// Equals reports whether two PdfPoints are equal by comparing both X and Y
// coordinates using the configured DoubleComparer.
func (c *PointComparer) Equals(a, b core.PdfPoint) bool {
	return c.doubleComparer.Equals(a.X, b.X) && c.doubleComparer.Equals(a.Y, b.Y)
}
