// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// SetCharacterSpacing sets the character spacing in the text state.
type SetCharacterSpacing struct {
	Spacing float64
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setCharacterSpacingSymbol = "Tc"

// NewSetCharacterSpacing creates a new SetCharacterSpacing operation.
func NewSetCharacterSpacing(spacing float64) *SetCharacterSpacing {
	return &SetCharacterSpacing{Spacing: spacing}
}

// Operator returns the operator symbol for this operation.
func (o SetCharacterSpacing) Operator() string {
	return setCharacterSpacingSymbol
}

// Run executes the operation on the given context, setting the character spacing.
func (o *SetCharacterSpacing) Run(ctx OperationContext) {
	ctx.SetCharacterSpacing(o.Spacing)
}

// Write writes the operator to the given writer.
func (o SetCharacterSpacing) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %s\n", o.Spacing, setCharacterSpacingSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetCharacterSpacing) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetCharacterSpacing) String() string {
	return fmt.Sprintf("%g %s", o.Spacing, setCharacterSpacingSymbol)
}
