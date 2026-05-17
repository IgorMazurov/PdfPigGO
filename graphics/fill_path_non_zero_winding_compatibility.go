// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// FillPathNonZeroWindingCompatibility fills the path using the nonzero winding number rule.
// Equivalent to FillPathNonZeroWinding, included only for compatibility. Although PDF consumer
// applications must be able to accept this operator (F), PDF producer applications should use
// FillPathNonZeroWinding (f) instead.
type FillPathNonZeroWindingCompatibility struct{}

const fillPathNonZeroWindingCompatSymbol = "F"

// InstanceFillPathNonZeroWindingCompatibility is the singleton instance of FillPathNonZeroWindingCompatibility.
var InstanceFillPathNonZeroWindingCompatibility = FillPathNonZeroWindingCompatibility{}

// Operator returns the operator symbol for this operation.
func (o FillPathNonZeroWindingCompatibility) Operator() string {
	return fillPathNonZeroWindingCompatSymbol
}

// Run executes the operation on the given context, filling the path with nonzero winding rule without closing subpath.
func (o FillPathNonZeroWindingCompatibility) Run(ctx OperationContext) {
	ctx.FillPath(core.FillingRuleNonZeroWinding, false)
}

// Write writes the operator symbol to the given writer.
// Although PDF reader applications shall be able to accept this operator,
// PDF writer applications should use f instead.
func (o FillPathNonZeroWindingCompatibility) Write(w io.Writer) error {
	if _, err := io.WriteString(w, fillPathNonZeroWindingSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (FillPathNonZeroWindingCompatibility) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o FillPathNonZeroWindingCompatibility) String() string {
	return fillPathNonZeroWindingCompatSymbol
}
