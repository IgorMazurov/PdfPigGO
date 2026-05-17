// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"
)

// SetNonStrokeColorDeviceCmyk sets the non-stroking color space to DeviceCMYK and set the color.
type SetNonStrokeColorDeviceCmyk struct {
	C float64
	M float64
	Y float64
	K float64
}

const setNonStrokeColorDeviceCmykSymbol = "k"

// NewSetNonStrokeColorDeviceCmyk creates a new SetNonStrokeColorDeviceCmyk operation.
func NewSetNonStrokeColorDeviceCmyk(c, m, y, k float64) *SetNonStrokeColorDeviceCmyk {
	return &SetNonStrokeColorDeviceCmyk{C: c, M: m, Y: y, K: k}
}

// Operator returns the operator symbol for this operation.
func (o SetNonStrokeColorDeviceCmyk) Operator() string {
	return setNonStrokeColorDeviceCmykSymbol
}

// Run executes the operation on the given context, setting the non-stroking color to CMYK.
func (o *SetNonStrokeColorDeviceCmyk) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetNonStrokingColorCmyk(o.C, o.M, o.Y, o.K)
}

// Write writes the operator to the given writer.
func (o SetNonStrokeColorDeviceCmyk) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %g %g", o.C, o.M, o.Y, o.K); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, " "+setNonStrokeColorDeviceCmykSymbol); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetNonStrokeColorDeviceCmyk) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetNonStrokeColorDeviceCmyk) String() string {
	return strings.Join([]string{
		fmt.Sprintf("%g", o.C),
		fmt.Sprintf("%g", o.M),
		fmt.Sprintf("%g", o.Y),
		fmt.Sprintf("%g", o.K),
		setNonStrokeColorDeviceCmykSymbol,
	}, " ")
}
