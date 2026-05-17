// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
)

// SetLineJoin sets the line join style in the graphics state.
type SetLineJoin struct {
	Join graphiccore.LineJoinStyle
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setLineJoinSymbol = "j"

// NewSetLineJoin creates a new SetLineJoin operation with validation.
func NewSetLineJoin(join graphiccore.LineJoinStyle) (*SetLineJoin, error) {
	if join < 0 || join > 2 {
		return nil, fmt.Errorf("invalid argument passed for line join style: should be 0, 1 or 2; instead got: %d", join)
	}
	return &SetLineJoin{Join: join}, nil
}

// Operator returns the operator symbol for this operation.
func (o SetLineJoin) Operator() string {
	return setLineJoinSymbol
}

// Run executes the operation on the given context, setting the line join style.
func (o *SetLineJoin) Run(ctx OperationContext) {
	ctx.SetLineJoin(o.Join)
}

// Write writes the operator to the given writer.
func (o SetLineJoin) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%d %s\n", o.Join, setLineJoinSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetLineJoin) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetLineJoin) String() string {
	return fmt.Sprintf("%d %s", o.Join, setLineJoinSymbol)
}
