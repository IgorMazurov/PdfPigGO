// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// SetMiterLimit sets the miter limit in the graphics state.
type SetMiterLimit struct {
	Limit float64
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setMiterLimitSymbol = "M"

// NewSetMiterLimit creates a new SetMiterLimit operation.
func NewSetMiterLimit(limit float64) *SetMiterLimit {
	return &SetMiterLimit{Limit: limit}
}

// Operator returns the operator symbol for this operation.
func (o SetMiterLimit) Operator() string {
	return setMiterLimitSymbol
}

// Run executes the operation on the given context, setting the miter limit.
func (o *SetMiterLimit) Run(ctx OperationContext) {
	ctx.SetMiterLimit(o.Limit)
}

// Write writes the operator to the given writer.
func (o SetMiterLimit) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %s\n", o.Limit, setMiterLimitSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetMiterLimit) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetMiterLimit) String() string {
	return fmt.Sprintf("%g %s", o.Limit, setMiterLimitSymbol)
}
