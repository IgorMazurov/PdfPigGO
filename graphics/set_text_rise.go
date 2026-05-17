// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// SetTextRise sets the vertical displacement of the baseline from the default position.
type SetTextRise struct {
	Rise float64
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setTextRiseSymbol = "Ts"

// NewSetTextRise creates a new SetTextRise operation.
func NewSetTextRise(rise float64) *SetTextRise {
	return &SetTextRise{Rise: rise}
}

// Operator returns the operator symbol for this operation.
func (o SetTextRise) Operator() string {
	return setTextRiseSymbol
}

// Run executes the operation on the given context, setting the text rise.
func (o *SetTextRise) Run(ctx OperationContext) {
	ctx.SetTextRise(o.Rise)
}

// Write writes the operator to the given writer.
func (o SetTextRise) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %s\n", o.Rise, setTextRiseSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetTextRise) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetTextRise) String() string {
	return fmt.Sprintf("%g %s", o.Rise, setTextRiseSymbol)
}
