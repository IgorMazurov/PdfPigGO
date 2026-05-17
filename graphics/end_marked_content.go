// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"
)

// EndMarkedContent ends a marked-content sequence.
type EndMarkedContent struct{}

// Symbol is the operator symbol for this operation in a PDF content stream.
const endMarkedContentSymbol = "EMC"

// InstanceEndMarkedContent is the singleton instance of EndMarkedContent.
var InstanceEndMarkedContent = EndMarkedContent{}

// Operator returns the operator symbol for this operation.
func (o EndMarkedContent) Operator() string {
	return endMarkedContentSymbol
}

// Run executes the operation on the given context, ending the current marked content sequence.
func (o EndMarkedContent) Run(ctx OperationContext) {
	ctx.EndMarkedContent()
}

// Write writes the operator symbol to the given writer.
func (o EndMarkedContent) Write(w io.Writer) error {
	if _, err := io.WriteString(w, endMarkedContentSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (EndMarkedContent) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o EndMarkedContent) String() string {
	return endMarkedContentSymbol
}
