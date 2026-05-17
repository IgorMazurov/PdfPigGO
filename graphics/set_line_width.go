// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// SetLineWidth sets the line width in the graphics state.
type SetLineWidth struct {
	Width float64
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setLineWidthSymbol = "w"

// NewSetLineWidth creates a new SetLineWidth operation.
func NewSetLineWidth(width float64) *SetLineWidth {
	return &SetLineWidth{Width: width}
}

// Operator returns the operator symbol for this operation.
func (o SetLineWidth) Operator() string {
	return setLineWidthSymbol
}

// Run executes the operation on the given context, setting the line width.
func (o *SetLineWidth) Run(ctx OperationContext) {
	ctx.SetLineWidth(o.Width)
}

// Write writes the operator to the given writer.
func (o SetLineWidth) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %s\n", o.Width, setLineWidthSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetLineWidth) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetLineWidth) String() string {
	return fmt.Sprintf("%g %s", o.Width, setLineWidthSymbol)
}
