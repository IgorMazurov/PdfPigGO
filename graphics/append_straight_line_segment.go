// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// AppendStraightLineSegment appends a straight line segment from the current point to (x, y).
type AppendStraightLineSegment struct {
	X float64
	Y float64
}

const appendStraightLineSegmentSymbol = "l"

// NewAppendStraightLineSegment creates a new AppendStraightLineSegment operation.
func NewAppendStraightLineSegment(x, y float64) *AppendStraightLineSegment {
	return &AppendStraightLineSegment{X: x, Y: y}
}

// Operator returns the operator symbol for this operation.
func (o AppendStraightLineSegment) Operator() string {
	return appendStraightLineSegmentSymbol
}

// Run executes the operation on the given context, drawing a line to (x, y).
func (o *AppendStraightLineSegment) Run(ctx OperationContext) {
	ctx.LineTo(o.X, o.Y)
}

// Write writes the operator to the given writer.
func (o AppendStraightLineSegment) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %s\n", o.X, o.Y, appendStraightLineSegmentSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (AppendStraightLineSegment) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o AppendStraightLineSegment) String() string {
	return fmt.Sprintf("%g %g %s", o.X, o.Y, appendStraightLineSegmentSymbol)
}
