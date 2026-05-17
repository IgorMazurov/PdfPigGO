// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// AppendRectangle appends a rectangle to the current path as a complete subpath.
type AppendRectangle struct {
	LowerLeftX float64
	LowerLeftY float64
	Width      float64
	Height     float64
}

const appendRectangleSymbol = "re"

// NewAppendRectangle creates a new AppendRectangle operation.
func NewAppendRectangle(x, y, width, height float64) *AppendRectangle {
	return &AppendRectangle{LowerLeftX: x, LowerLeftY: y, Width: width, Height: height}
}

// Operator returns the operator symbol for this operation.
func (o AppendRectangle) Operator() string {
	return appendRectangleSymbol
}

// Run executes the operation on the given context, adding a rectangle subpath.
func (o *AppendRectangle) Run(ctx OperationContext) {
	ctx.Rectangle(o.LowerLeftX, o.LowerLeftY, o.Width, o.Height)
}

// Write writes the operator to the given writer.
func (o AppendRectangle) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %g %g %s\n", o.LowerLeftX, o.LowerLeftY, o.Width, o.Height, appendRectangleSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (AppendRectangle) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o AppendRectangle) String() string {
	return fmt.Sprintf("%g %g %g %g %s", o.LowerLeftX, o.LowerLeftY, o.Width, o.Height, appendRectangleSymbol)
}
