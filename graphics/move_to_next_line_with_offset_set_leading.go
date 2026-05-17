// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/content"
)

// MoveToNextLineWithOffsetSetLeading moves to the start of the next line offset by (Tx, Ty)
// and sets the leading parameter in the text state.
// PDF operator: TD
type MoveToNextLineWithOffsetSetLeading struct {
	Tx float64
	Ty float64
}

const moveToNextLineWithOffsetSetLeadingSymbol = "TD"

// NewMoveToNextLineWithOffsetSetLeading creates a new MoveToNextLineWithOffsetSetLeading operation.
func NewMoveToNextLineWithOffsetSetLeading(tx, ty float64) *MoveToNextLineWithOffsetSetLeading {
	return &MoveToNextLineWithOffsetSetLeading{Tx: tx, Ty: ty}
}

// Operator returns the operator symbol for this operation.
func (o MoveToNextLineWithOffsetSetLeading) Operator() string {
	return moveToNextLineWithOffsetSetLeadingSymbol
}

// Run executes the move-to-next-line-with-offset-and-set-leading operation on the given context.
// First sets text leading to -Ty, then moves to next line with offset (Tx, Ty).
func (o *MoveToNextLineWithOffsetSetLeading) Run(ctx OperationContext) {
	ctx.SetTextLeading(-o.Ty)

	moveOp := NewMoveToNextLineWithOffset(o.Tx, o.Ty)
	moveOp.Run(ctx)
}

// Write writes the operator to the given writer.
func (o MoveToNextLineWithOffsetSetLeading) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %s\n", o.Tx, o.Ty, moveToNextLineWithOffsetSetLeadingSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (MoveToNextLineWithOffsetSetLeading) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o MoveToNextLineWithOffsetSetLeading) String() string {
	return fmt.Sprintf("%g %g %s", o.Tx, o.Ty, moveToNextLineWithOffsetSetLeadingSymbol)
}

var _ content.GraphicsStateOperation = (*MoveToNextLineWithOffsetSetLeading)(nil)
