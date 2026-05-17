// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"
)

// SetStrokeColor sets the stroking color based on the current color space.
type SetStrokeColor struct {
	Operands []float64
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setStrokeColorSymbol = "SC"

// NewSetStrokeColor creates a new SetStrokeColor operation.
func NewSetStrokeColor(operands []float64) *SetStrokeColor {
	return &SetStrokeColor{Operands: operands}
}

// Operator returns the operator symbol for this operation.
func (o SetStrokeColor) Operator() string {
	return setStrokeColorSymbol
}

// Run executes the operation on the given context, setting the stroking color.
func (o *SetStrokeColor) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetStrokingColor(o.Operands, nil)
}

// Write writes the operator to the given writer.
func (o SetStrokeColor) Write(w io.Writer) error {
	for _, operand := range o.Operands {
		if _, err := fmt.Fprintf(w, "%g ", operand); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(w, setStrokeColorSymbol); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetStrokeColor) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetStrokeColor) String() string {
	parts := make([]string, len(o.Operands))
	for i, operand := range o.Operands {
		parts[i] = fmt.Sprintf("%g", operand)
	}
	return strings.Join(parts, " ") + " " + setStrokeColorSymbol
}
