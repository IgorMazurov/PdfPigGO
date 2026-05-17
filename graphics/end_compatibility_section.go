// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"
)

// EndCompatibilitySection ends a compatibility section. Unrecognized operators
// (along with their operands) are ignored without error.
type EndCompatibilitySection struct{}

// Symbol is the operator symbol for this operation in a PDF content stream.
const endCompatibilitySectionSymbol = "EX"

// InstanceEndCompatibilitySection is the singleton instance of EndCompatibilitySection.
var InstanceEndCompatibilitySection = EndCompatibilitySection{}

// Operator returns the operator symbol for this operation.
func (o EndCompatibilitySection) Operator() string {
	return endCompatibilitySectionSymbol
}

// Run executes the operation on the given context. This is a no-op as per PDF spec.
func (o EndCompatibilitySection) Run(ctx OperationContext) {
}

// Write writes the operator symbol to the given writer.
func (o EndCompatibilitySection) Write(w io.Writer) error {
	if _, err := io.WriteString(w, endCompatibilitySectionSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (EndCompatibilitySection) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o EndCompatibilitySection) String() string {
	return endCompatibilitySectionSymbol
}
