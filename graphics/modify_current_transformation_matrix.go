// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// ModifyCurrentTransformationMatrix modifies the current transformation matrix by
// concatenating the specified 6-value transformation matrix ("cm" operator).
type ModifyCurrentTransformationMatrix struct {
	Value [6]float64
}

const modifyCurrentTransformationMatrixSymbol = "cm"

// NewModifyCurrentTransformationMatrix creates a new ModifyCurrentTransformationMatrix operation.
func NewModifyCurrentTransformationMatrix(v0, v1, v2, v3, v4, v5 float64) *ModifyCurrentTransformationMatrix {
	return &ModifyCurrentTransformationMatrix{Value: [6]float64{v0, v1, v2, v3, v4, v5}}
}

// Operator returns the operator symbol for this operation.
func (o ModifyCurrentTransformationMatrix) Operator() string {
	return modifyCurrentTransformationMatrixSymbol
}

// Run executes the operation on the given context, modifying the current transformation matrix.
func (o *ModifyCurrentTransformationMatrix) Run(ctx OperationContext) {
	matrix, _ := core.FromArray(o.Value[:])
	ctx.ModifyCurrentTransformationMatrix(matrix)
}

// Write writes the operator to the given writer.
func (o ModifyCurrentTransformationMatrix) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %g %g %g %g %s\n", o.Value[0], o.Value[1], o.Value[2], o.Value[3], o.Value[4], o.Value[5], modifyCurrentTransformationMatrixSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (ModifyCurrentTransformationMatrix) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o ModifyCurrentTransformationMatrix) String() string {
	return fmt.Sprintf("%g %g %g %g %g %g %s", o.Value[0], o.Value[1], o.Value[2], o.Value[3], o.Value[4], o.Value[5], modifyCurrentTransformationMatrixSymbol)
}
