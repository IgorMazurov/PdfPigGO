// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"
)

// BeginCompatibilitySection begins a compatibility section. Unrecognized operators
// (along with their operands) are ignored without error.
type BeginCompatibilitySection struct{}

// Symbol is the operator symbol for this operation in a PDF content stream.
const beginCompatibilitySectionSymbol = "BX"

// InstanceBeginCompatibilitySection is the singleton instance of BeginCompatibilitySection.
var InstanceBeginCompatibilitySection = BeginCompatibilitySection{}

// Operator returns the operator symbol for this operation.
func (o BeginCompatibilitySection) Operator() string {
	return beginCompatibilitySectionSymbol
}

// Run executes the operation on the given context. This is a no-op as per PDF spec.
func (o BeginCompatibilitySection) Run(ctx OperationContext) {
}

// Write writes the operator symbol to the given writer.
func (o BeginCompatibilitySection) Write(w io.Writer) error {
	if _, err := io.WriteString(w, beginCompatibilitySectionSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (BeginCompatibilitySection) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o BeginCompatibilitySection) String() string {
	return beginCompatibilitySectionSymbol
}
