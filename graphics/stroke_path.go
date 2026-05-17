// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"
)

// StrokePath strokes the current path without closing it.
type StrokePath struct{}

const strokePathSymbol = "S"

// InstanceStrokePath is the singleton instance of StrokePath.
var InstanceStrokePath = StrokePath{}

// Operator returns the operator symbol for this operation.
func (o StrokePath) Operator() string {
	return strokePathSymbol
}

// Run executes the operation on the given context, stroking the path without closing.
func (o StrokePath) Run(ctx OperationContext) {
	ctx.StrokePath(false)
}

// Write writes the operator symbol to the given writer.
func (o StrokePath) Write(w io.Writer) error {
	if _, err := io.WriteString(w, strokePathSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (StrokePath) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o StrokePath) String() string {
	return strokePathSymbol
}
