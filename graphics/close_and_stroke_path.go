// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"
)

// CloseAndStrokePath closes the current subpath by appending a straight line segment from
// the current point to the starting point of the subpath, then strokes the path.
type CloseAndStrokePath struct{}

const closeAndStrokePathSymbol = "s"

// InstanceCloseAndStrokePath is the singleton instance of CloseAndStrokePath.
var InstanceCloseAndStrokePath = CloseAndStrokePath{}

// Operator returns the operator symbol for this operation.
func (o CloseAndStrokePath) Operator() string {
	return closeAndStrokePathSymbol
}

// Run executes the operation on the given context, closing the subpath and stroking it.
func (o CloseAndStrokePath) Run(ctx OperationContext) {
	ctx.StrokePath(true)
}

// Write writes the operator symbol to the given writer.
func (o CloseAndStrokePath) Write(w io.Writer) error {
	if _, err := io.WriteString(w, closeAndStrokePathSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (CloseAndStrokePath) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o CloseAndStrokePath) String() string {
	return closeAndStrokePathSymbol
}
