// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// ModifyClippingByNonZeroWindingIntersect modifies the current clipping path by intersecting it with
// the current path, using the nonzero winding number rule to determine which regions lie inside.
type ModifyClippingByNonZeroWindingIntersect struct{}

// Symbol is the operator symbol for this operation in a PDF content stream.
const modifyClippingByNonZeroWindingIntersectSymbol = "W"

// InstanceModifyClippingByNonZeroWindingIntersect is the singleton instance of ModifyClippingByNonZeroWindingIntersect.
var InstanceModifyClippingByNonZeroWindingIntersect = ModifyClippingByNonZeroWindingIntersect{}

// Operator returns the operator symbol for this operation.
func (o ModifyClippingByNonZeroWindingIntersect) Operator() string {
	return modifyClippingByNonZeroWindingIntersectSymbol
}

// Run executes the operation on the given context.
func (o ModifyClippingByNonZeroWindingIntersect) Run(ctx OperationContext) {
	ctx.ModifyClippingIntersect(core.FillingRuleNonZeroWinding)
}

// Write writes the operator symbol to the given writer.
func (o ModifyClippingByNonZeroWindingIntersect) Write(w io.Writer) error {
	if _, err := io.WriteString(w, modifyClippingByNonZeroWindingIntersectSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (ModifyClippingByNonZeroWindingIntersect) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o ModifyClippingByNonZeroWindingIntersect) String() string {
	return modifyClippingByNonZeroWindingIntersectSymbol
}
