// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/content"
)

// MoveToNextLine moves to the start of the next line (T* operator).
// Performs: 0 -Tl Td
type MoveToNextLine struct{}

const moveToNextLineSymbol = "T*"

// InstanceMoveToNextLine is the singleton instance of MoveToNextLine.
var InstanceMoveToNextLine = MoveToNextLine{}

// Operator returns the operator symbol for this operation.
func (o MoveToNextLine) Operator() string {
	return moveToNextLineSymbol
}

// Run executes the operation on the given context, moving to the next line with offset.
func (o MoveToNextLine) Run(ctx OperationContext) {
	ctx.MoveToNextLineWithOffset()
}

// Write writes the operator symbol to the given writer.
func (o MoveToNextLine) Write(w io.Writer) error {
	if _, err := io.WriteString(w, moveToNextLineSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (MoveToNextLine) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o MoveToNextLine) String() string {
	return moveToNextLineSymbol
}

var _ content.GraphicsStateOperation = MoveToNextLine{}
