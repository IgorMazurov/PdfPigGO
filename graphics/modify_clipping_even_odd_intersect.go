// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// ModifyClippingByEvenOddIntersect modifies the current clipping path by intersecting it with
// the current path, using the even-odd rule to determine which regions lie inside.
type ModifyClippingByEvenOddIntersect struct{}

// Symbol is the operator symbol for this operation in a PDF content stream.
const modifyClippingByEvenOddIntersectSymbol = "W*"

// Instance is the singleton instance of ModifyClippingByEvenOddIntersect.
var InstanceModifyClippingByEvenOddIntersect = ModifyClippingByEvenOddIntersect{}

// Operator returns the operator symbol for this operation.
func (o ModifyClippingByEvenOddIntersect) Operator() string {
	return modifyClippingByEvenOddIntersectSymbol
}

// Run executes the operation on the given context.
func (o ModifyClippingByEvenOddIntersect) Run(ctx OperationContext) {
	ctx.ModifyClippingIntersect(core.FillingRuleEvenOdd)
}

// Write writes the operator symbol to the given writer.
func (o ModifyClippingByEvenOddIntersect) Write(w io.Writer) error {
	if _, err := io.WriteString(w, modifyClippingByEvenOddIntersectSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (ModifyClippingByEvenOddIntersect) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o ModifyClippingByEvenOddIntersect) String() string {
	return modifyClippingByEvenOddIntersectSymbol
}
