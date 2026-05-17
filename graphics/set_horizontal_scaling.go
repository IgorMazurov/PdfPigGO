// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// SetHorizontalScaling sets the horizontal text scaling percentage.
type SetHorizontalScaling struct {
	Scale float64
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setHorizontalScalingSymbol = "Tz"

// NewSetHorizontalScaling creates a new SetHorizontalScaling operation.
func NewSetHorizontalScaling(scale float64) *SetHorizontalScaling {
	return &SetHorizontalScaling{Scale: scale}
}

// Operator returns the operator symbol for this operation.
func (o SetHorizontalScaling) Operator() string {
	return setHorizontalScalingSymbol
}

// Run executes the operation on the given context, setting the horizontal scaling percentage.
func (o *SetHorizontalScaling) Run(ctx OperationContext) {
	ctx.SetHorizontalScaling(o.Scale)
}

// Write writes the operator to the given writer.
func (o SetHorizontalScaling) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %s\n", o.Scale, setHorizontalScalingSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetHorizontalScaling) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetHorizontalScaling) String() string {
	return fmt.Sprintf("%g %s", o.Scale, setHorizontalScalingSymbol)
}
