// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// SetFlatnessTolerance sets the flatness tolerance in the graphics state.
type SetFlatnessTolerance struct {
	Tolerance float64
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setFlatnessToleranceSymbol = "i"

// NewSetFlatnessTolerance creates a new SetFlatnessTolerance operation.
func NewSetFlatnessTolerance(tolerance float64) *SetFlatnessTolerance {
	return &SetFlatnessTolerance{Tolerance: tolerance}
}

// Operator returns the operator symbol for this operation.
func (o SetFlatnessTolerance) Operator() string {
	return setFlatnessToleranceSymbol
}

// Run executes the operation on the given context, setting the flatness tolerance.
func (o *SetFlatnessTolerance) Run(ctx OperationContext) {
	ctx.SetFlatnessTolerance(o.Tolerance)
}

// Write writes the operator to the given writer.
func (o SetFlatnessTolerance) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %s\n", o.Tolerance, setFlatnessToleranceSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetFlatnessTolerance) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetFlatnessTolerance) String() string {
	return fmt.Sprintf("%g %s", o.Tolerance, setFlatnessToleranceSymbol)
}
