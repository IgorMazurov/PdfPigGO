// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
)

// EndText ends a text object, discarding the text matrix.
type EndText struct{}

const endTextSymbol = "ET"

// InstanceEndText is the singleton instance of EndText.
var InstanceEndText = EndText{}

// Operator returns the operator symbol for this operation.
func (o EndText) Operator() string {
	return endTextSymbol
}

// Run executes the operation on the given context, resetting text matrices to identity.
func (o EndText) Run(ctx OperationContext) {
	ctx.TextMatrices().TextMatrix = core.Identity
	ctx.TextMatrices().TextLineMatrix = core.Identity
}

// Write writes the operator symbol to the given writer.
func (o EndText) Write(w io.Writer) error {
	if _, err := io.WriteString(w, endTextSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (EndText) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o EndText) String() string {
	return endTextSymbol
}

var _ content.GraphicsStateOperation = EndText{}
