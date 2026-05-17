// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"
)

// EndPath ends the current path construction without filling or stroking.
type EndPath struct{}

const endPathSymbol = "n"

// InstanceEndPath is the singleton instance of EndPath.
var InstanceEndPath = EndPath{}

// Operator returns the operator symbol for this operation.
func (o EndPath) Operator() string {
	return endPathSymbol
}

// Run executes the operation on the given context, ending the current path.
func (o EndPath) Run(ctx OperationContext) {
	ctx.EndPath()
}

// Write writes the operator symbol to the given writer.
func (o EndPath) Write(w io.Writer) error {
	if _, err := io.WriteString(w, endPathSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (EndPath) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o EndPath) String() string {
	return endPathSymbol
}
