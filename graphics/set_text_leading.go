// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// SetTextLeading sets the text leading in unscaled text space units.
type SetTextLeading struct {
	Leading float64
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setTextLeadingSymbol = "TL"

// NewSetTextLeading creates a new SetTextLeading operation.
func NewSetTextLeading(leading float64) *SetTextLeading {
	return &SetTextLeading{Leading: leading}
}

// Operator returns the operator symbol for this operation.
func (o SetTextLeading) Operator() string {
	return setTextLeadingSymbol
}

// Run executes the operation on the given context, setting the text leading.
func (o *SetTextLeading) Run(ctx OperationContext) {
	ctx.SetTextLeading(o.Leading)
}

// Write writes the operator to the given writer.
func (o SetTextLeading) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %s\n", o.Leading, setTextLeadingSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetTextLeading) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetTextLeading) String() string {
	return fmt.Sprintf("%g %s", o.Leading, setTextLeadingSymbol)
}
