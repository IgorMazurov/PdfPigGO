// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
)

// SetLineCap sets the line cap style in the graphics state.
type SetLineCap struct {
	Cap graphiccore.LineCapStyle
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setLineCapSymbol = "J"

// NewSetLineCap creates a new SetLineCap operation with validation.
func NewSetLineCap(cap graphiccore.LineCapStyle) (*SetLineCap, error) {
	if cap < 0 || cap > 2 {
		return nil, fmt.Errorf("invalid argument passed for line cap style: should be 0, 1 or 2; instead got: %d", cap)
	}
	return &SetLineCap{Cap: cap}, nil
}

// Operator returns the operator symbol for this operation.
func (o SetLineCap) Operator() string {
	return setLineCapSymbol
}

// Run executes the operation on the given context, setting the line cap style.
func (o *SetLineCap) Run(ctx OperationContext) {
	ctx.SetLineCap(o.Cap)
}

// Write writes the operator to the given writer.
func (o SetLineCap) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%d %s\n", o.Cap, setLineCapSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetLineCap) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetLineCap) String() string {
	return fmt.Sprintf("%d %s", o.Cap, setLineCapSymbol)
}
