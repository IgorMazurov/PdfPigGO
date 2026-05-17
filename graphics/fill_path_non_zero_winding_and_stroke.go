// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// FillPathNonZeroWindingAndStroke fills the path using the nonzero winding number rule and then strokes it.
type FillPathNonZeroWindingAndStroke struct{}

const fillPathNonZeroWindingAndStrokeSymbol = "B"

// InstanceFillPathNonZeroWindingAndStroke is the singleton instance of FillPathNonZeroWindingAndStroke.
var InstanceFillPathNonZeroWindingAndStroke = FillPathNonZeroWindingAndStroke{}

// Operator returns the operator symbol for this operation.
func (o FillPathNonZeroWindingAndStroke) Operator() string {
	return fillPathNonZeroWindingAndStrokeSymbol
}

// Run executes the operation on the given context, filling with nonzero winding rule and stroking without closing subpath.
func (o FillPathNonZeroWindingAndStroke) Run(ctx OperationContext) {
	ctx.FillStrokePath(core.FillingRuleNonZeroWinding, false)
}

// Write writes the operator symbol to the given writer.
func (o FillPathNonZeroWindingAndStroke) Write(w io.Writer) error {
	if _, err := io.WriteString(w, fillPathNonZeroWindingAndStrokeSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (FillPathNonZeroWindingAndStroke) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o FillPathNonZeroWindingAndStroke) String() string {
	return fillPathNonZeroWindingAndStrokeSymbol
}
