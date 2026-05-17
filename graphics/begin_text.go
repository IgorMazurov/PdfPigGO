// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
)

// BeginText begins a text object, initializing the text matrix and the text line matrix to the identity matrix.
// Text objects cannot be nested.
type BeginText struct{}

const beginTextSymbol = "BT"

// InstanceBeginText is the singleton instance of BeginText.
var InstanceBeginText = BeginText{}

// Operator returns the operator symbol for this operation.
func (o BeginText) Operator() string {
	return beginTextSymbol
}

// Run executes the operation on the given context, resetting text matrices to identity.
func (o BeginText) Run(ctx OperationContext) {
	ctx.TextMatrices().TextMatrix = core.Identity
	ctx.TextMatrices().TextLineMatrix = core.Identity
}

// Write writes the operator symbol to the given writer.
func (o BeginText) Write(w io.Writer) error {
	if _, err := io.WriteString(w, beginTextSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (BeginText) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o BeginText) String() string {
	return beginTextSymbol
}

var _ content.GraphicsStateOperation = BeginText{}
