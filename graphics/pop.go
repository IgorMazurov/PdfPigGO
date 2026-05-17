// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/content"
)

// Pop restores the graphics state by removing the most recently saved state from the stack
// and making it the current state.
type Pop struct{}

const popSymbol = "Q"

// InstancePop is the singleton instance of Pop.
var InstancePop = Pop{}

// Operator returns the operator symbol for this operation.
func (o Pop) Operator() string {
	return popSymbol
}

// Run executes the operation on the given context, popping the graphics state stack.
func (o Pop) Run(ctx OperationContext) {
	ctx.PopState()
}

// Write writes the operator symbol to the given writer.
func (o Pop) Write(w io.Writer) error {
	if _, err := io.WriteString(w, popSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (Pop) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o Pop) String() string {
	return popSymbol
}

var _ content.GraphicsStateOperation = Pop{}
