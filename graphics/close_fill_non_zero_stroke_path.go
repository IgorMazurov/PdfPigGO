// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// CloseFillNonZeroStrokePath closes the current subpath, fills using the nonzero winding number rule,
// and then strokes the path.
type CloseFillNonZeroStrokePath struct{}

const closeFillNonZeroStrokePathSymbol = "b"

// InstanceCloseFillNonZeroStrokePath is the singleton instance of CloseFillNonZeroStrokePath.
var InstanceCloseFillNonZeroStrokePath = CloseFillNonZeroStrokePath{}

// Operator returns the operator symbol for this operation.
func (o CloseFillNonZeroStrokePath) Operator() string {
	return closeFillNonZeroStrokePathSymbol
}

// Run executes the operation on the given context, closing the subpath, filling with nonzero winding rule, and stroking.
func (o CloseFillNonZeroStrokePath) Run(ctx OperationContext) {
	ctx.FillStrokePath(core.FillingRuleNonZeroWinding, true)
}

// Write writes the operator symbol to the given writer.
func (o CloseFillNonZeroStrokePath) Write(w io.Writer) error {
	if _, err := io.WriteString(w, closeFillNonZeroStrokePathSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (CloseFillNonZeroStrokePath) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o CloseFillNonZeroStrokePath) String() string {
	return closeFillNonZeroStrokePathSymbol
}
