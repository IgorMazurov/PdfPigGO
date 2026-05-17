// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// FillPathEvenOddRule fills the path using the even-odd rule to determine the region to fill.
type FillPathEvenOddRule struct{}

const fillPathEvenOddRuleSymbol = "f*"

// InstanceFillPathEvenOddRule is the singleton instance of FillPathEvenOddRule.
var InstanceFillPathEvenOddRule = FillPathEvenOddRule{}

// Operator returns the operator symbol for this operation.
func (o FillPathEvenOddRule) Operator() string {
	return fillPathEvenOddRuleSymbol
}

// Run executes the operation on the given context, filling the path with even-odd rule.
func (o FillPathEvenOddRule) Run(ctx OperationContext) {
	ctx.FillPath(core.FillingRuleEvenOdd, false)
}

// Write writes the operator symbol to the given writer.
func (o FillPathEvenOddRule) Write(w io.Writer) error {
	if _, err := io.WriteString(w, fillPathEvenOddRuleSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (FillPathEvenOddRule) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o FillPathEvenOddRule) String() string {
	return fillPathEvenOddRuleSymbol
}
