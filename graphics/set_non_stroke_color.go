// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"
)

// SetNonStrokeColor sets the non-stroking color based on the current color space.
type SetNonStrokeColor struct {
	Operands []float64
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setNonStrokeColorSymbol = "sc"

// NewSetNonStrokeColor creates a new SetNonStrokeColor operation.
func NewSetNonStrokeColor(operands []float64) *SetNonStrokeColor {
	return &SetNonStrokeColor{Operands: operands}
}

// Operator returns the operator symbol for this operation.
func (o SetNonStrokeColor) Operator() string {
	return setNonStrokeColorSymbol
}

// Run executes the operation on the given context, setting the non-stroking color.
func (o *SetNonStrokeColor) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetNonStrokingColor(o.Operands, nil)
}

// Write writes the operator to the given writer.
func (o SetNonStrokeColor) Write(w io.Writer) error {
	for _, operand := range o.Operands {
		if _, err := fmt.Fprintf(w, "%g ", operand); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(w, setNonStrokeColorSymbol); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetNonStrokeColor) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetNonStrokeColor) String() string {
	parts := make([]string, len(o.Operands))
	for i, operand := range o.Operands {
		parts[i] = fmt.Sprintf("%g", operand)
	}
	return strings.Join(parts, " ") + " " + setNonStrokeColorSymbol
}
