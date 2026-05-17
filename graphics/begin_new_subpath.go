// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// BeginNewSubpath begins a new subpath by moving to (x, y), omitting any connecting line segment.
type BeginNewSubpath struct {
	X float64
	Y float64
}

const beginNewSubpathSymbol = "m"

// NewBeginNewSubpath creates a new BeginNewSubpath operation.
func NewBeginNewSubpath(x, y float64) *BeginNewSubpath {
	return &BeginNewSubpath{X: x, Y: y}
}

// Operator returns the operator symbol for this operation.
func (o BeginNewSubpath) Operator() string {
	return beginNewSubpathSymbol
}

// Run executes the operation on the given context, moving to (x, y).
func (o *BeginNewSubpath) Run(ctx OperationContext) {
	ctx.MoveTo(o.X, o.Y)
}

// Write writes the operator to the given writer.
func (o BeginNewSubpath) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %s\n", o.X, o.Y, beginNewSubpathSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (BeginNewSubpath) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o BeginNewSubpath) String() string {
	return fmt.Sprintf("%g %g %s", o.X, o.Y, beginNewSubpathSymbol)
}
