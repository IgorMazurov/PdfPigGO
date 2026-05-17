// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// FillPathEvenOddRuleAndStroke fills the path using the even-odd rule and then strokes it.
type FillPathEvenOddRuleAndStroke struct{}

const fillPathEvenOddRuleAndStrokeSymbol = "B*"

// InstanceFillPathEvenOddRuleAndStroke is the singleton instance of FillPathEvenOddRuleAndStroke.
var InstanceFillPathEvenOddRuleAndStroke = FillPathEvenOddRuleAndStroke{}

// Operator returns the operator symbol for this operation.
func (o FillPathEvenOddRuleAndStroke) Operator() string {
	return fillPathEvenOddRuleAndStrokeSymbol
}

// Run executes the operation on the given context, filling with even-odd rule and stroking without closing subpath.
func (o FillPathEvenOddRuleAndStroke) Run(ctx OperationContext) {
	ctx.FillStrokePath(core.FillingRuleEvenOdd, false)
}

// Write writes the operator symbol to the given writer.
func (o FillPathEvenOddRuleAndStroke) Write(w io.Writer) error {
	if _, err := io.WriteString(w, fillPathEvenOddRuleAndStrokeSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (FillPathEvenOddRuleAndStroke) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o FillPathEvenOddRuleAndStroke) String() string {
	return fillPathEvenOddRuleAndStrokeSymbol
}
