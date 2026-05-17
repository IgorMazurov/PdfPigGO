// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/content"
)

// Push saves the current graphics state on the graphics state stack.
type Push struct{}

const pushSymbol = "q"

// InstancePush is the singleton instance of Push.
var InstancePush = Push{}

// Operator returns the operator symbol for this operation.
func (o Push) Operator() string {
	return pushSymbol
}

// Run executes the operation on the given context, pushing the graphics state stack.
func (o Push) Run(ctx OperationContext) {
	ctx.PushState()
}

// Write writes the operator symbol to the given writer.
func (o Push) Write(w io.Writer) error {
	if _, err := io.WriteString(w, pushSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (Push) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o Push) String() string {
	return pushSymbol
}

var _ content.GraphicsStateOperation = Push{}
