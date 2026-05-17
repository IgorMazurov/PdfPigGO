// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// CloseFillEvenOddStrokePath closes the current subpath, fills using the even-odd rule,
// and then strokes the path.
type CloseFillEvenOddStrokePath struct{}

const closeFillEvenOddStrokePathSymbol = "b*"

// InstanceCloseFillEvenOddStrokePath is the singleton instance of CloseFillEvenOddStrokePath.
var InstanceCloseFillEvenOddStrokePath = CloseFillEvenOddStrokePath{}

// Operator returns the operator symbol for this operation.
func (o CloseFillEvenOddStrokePath) Operator() string {
	return closeFillEvenOddStrokePathSymbol
}

// Run executes the operation on the given context, closing the subpath, filling with even-odd rule, and stroking.
func (o CloseFillEvenOddStrokePath) Run(ctx OperationContext) {
	ctx.FillStrokePath(core.FillingRuleEvenOdd, true)
}

// Write writes the operator symbol to the given writer.
func (o CloseFillEvenOddStrokePath) Write(w io.Writer) error {
	if _, err := io.WriteString(w, closeFillEvenOddStrokePathSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (CloseFillEvenOddStrokePath) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o CloseFillEvenOddStrokePath) String() string {
	return closeFillEvenOddStrokePathSymbol
}
