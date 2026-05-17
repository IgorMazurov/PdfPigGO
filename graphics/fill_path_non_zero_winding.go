// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// FillPathNonZeroWinding fills the path using the nonzero winding number rule to determine the region to fill.
type FillPathNonZeroWinding struct{}

const fillPathNonZeroWindingSymbol = "f"

// InstanceFillPathNonZeroWinding is the singleton instance of FillPathNonZeroWinding.
var InstanceFillPathNonZeroWinding = FillPathNonZeroWinding{}

// Operator returns the operator symbol for this operation.
func (o FillPathNonZeroWinding) Operator() string {
	return fillPathNonZeroWindingSymbol
}

// Run executes the operation on the given context, filling the path with nonzero winding rule.
func (o FillPathNonZeroWinding) Run(ctx OperationContext) {
	ctx.FillPath(core.FillingRuleNonZeroWinding, false)
}

// Write writes the operator symbol to the given writer.
func (o FillPathNonZeroWinding) Write(w io.Writer) error {
	if _, err := io.WriteString(w, fillPathNonZeroWindingSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (FillPathNonZeroWinding) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o FillPathNonZeroWinding) String() string {
	return fillPathNonZeroWindingSymbol
}
