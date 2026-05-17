// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
)

// SetTextMatrix sets the text matrix and text line matrix ("Tm" operator).
type SetTextMatrix struct {
	Value [6]float64
}

const setTextMatrixSymbol = "Tm"

// NewSetTextMatrix creates a new SetTextMatrix operation.
func NewSetTextMatrix(v0, v1, v2, v3, v4, v5 float64) *SetTextMatrix {
	return &SetTextMatrix{Value: [6]float64{v0, v1, v2, v3, v4, v5}}
}

// Operator returns the operator symbol for this operation.
func (o SetTextMatrix) Operator() string {
	return setTextMatrixSymbol
}

// Run executes the set text matrix operation on the given context.
// Sets both TextMatrix and TextLineMatrix to the provided transformation matrix.
func (o *SetTextMatrix) Run(ctx OperationContext) {
	newMatrix, _ := core.FromArray(o.Value[:])
	ctx.TextMatrices().TextMatrix = newMatrix
	ctx.TextMatrices().TextLineMatrix = newMatrix
}

// Write writes the operator to the given writer.
func (o SetTextMatrix) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %g %g %g %g %s\n", o.Value[0], o.Value[1], o.Value[2], o.Value[3], o.Value[4], o.Value[5], setTextMatrixSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetTextMatrix) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetTextMatrix) String() string {
	return fmt.Sprintf("%g %g %g %g %g %g %s", o.Value[0], o.Value[1], o.Value[2], o.Value[3], o.Value[4], o.Value[5], setTextMatrixSymbol)
}

var _ content.GraphicsStateOperation = (*SetTextMatrix)(nil)
