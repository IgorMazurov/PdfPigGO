// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// SetWordSpacing sets the word spacing in unscaled text space units.
type SetWordSpacing struct {
	Spacing float64
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setWordSpacingSymbol = "Tw"

// NewSetWordSpacing creates a new SetWordSpacing operation.
func NewSetWordSpacing(spacing float64) *SetWordSpacing {
	return &SetWordSpacing{Spacing: spacing}
}

// Operator returns the operator symbol for this operation.
func (o SetWordSpacing) Operator() string {
	return setWordSpacingSymbol
}

// Run executes the operation on the given context, setting the word spacing.
func (o *SetWordSpacing) Run(ctx OperationContext) {
	ctx.SetWordSpacing(o.Spacing)
}

// Write writes the operator to the given writer.
func (o SetWordSpacing) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %s\n", o.Spacing, setWordSpacingSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetWordSpacing) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetWordSpacing) String() string {
	return fmt.Sprintf("%g %s", o.Spacing, setWordSpacingSymbol)
}
