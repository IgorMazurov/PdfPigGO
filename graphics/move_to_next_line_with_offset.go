// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
)

// MoveToNextLineWithOffset moves to the start of the next line offset by (Tx, Ty).
// PDF operator: Td
type MoveToNextLineWithOffset struct {
	Tx float64
	Ty float64
}

const moveToNextLineWithOffsetSymbol = "Td"

// NewMoveToNextLineWithOffset creates a new MoveToNextLineWithOffset operation.
func NewMoveToNextLineWithOffset(tx, ty float64) *MoveToNextLineWithOffset {
	return &MoveToNextLineWithOffset{Tx: tx, Ty: ty}
}

// Operator returns the operator symbol for this operation.
func (o MoveToNextLineWithOffset) Operator() string {
	return moveToNextLineWithOffsetSymbol
}

// Run executes the move-to-next-line-with-offset operation on the given context.
// Tm = Tlm = [1 0 0 1 Tx Ty] * Tlm
func (o *MoveToNextLineWithOffset) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil {
		return
	}

	currentTextLineMatrix := ctx.TextMatrices().TextLineMatrix
	matrix := core.FromValues(1, 0, 0, 1, o.Tx, o.Ty)
	transformed := matrix.Multiply(currentTextLineMatrix)

	ctx.TextMatrices().TextLineMatrix = transformed
	ctx.TextMatrices().TextMatrix = transformed
}

// Write writes the operator to the given writer.
func (o MoveToNextLineWithOffset) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %s\n", o.Tx, o.Ty, moveToNextLineWithOffsetSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (MoveToNextLineWithOffset) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o MoveToNextLineWithOffset) String() string {
	return fmt.Sprintf("%g %g %s", o.Tx, o.Ty, moveToNextLineWithOffsetSymbol)
}

var _ content.GraphicsStateOperation = (*MoveToNextLineWithOffset)(nil)
