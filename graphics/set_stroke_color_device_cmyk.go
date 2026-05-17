// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"
)

// SetStrokeColorDeviceCmyk sets the stroking color space to DeviceCMYK and set the color.
type SetStrokeColorDeviceCmyk struct {
	C float64
	M float64
	Y float64
	K float64
}

const setStrokeColorDeviceCmykSymbol = "K"

// NewSetStrokeColorDeviceCmyk creates a new SetStrokeColorDeviceCmyk operation.
func NewSetStrokeColorDeviceCmyk(c, m, y, k float64) *SetStrokeColorDeviceCmyk {
	return &SetStrokeColorDeviceCmyk{C: c, M: m, Y: y, K: k}
}

// Operator returns the operator symbol for this operation.
func (o SetStrokeColorDeviceCmyk) Operator() string {
	return setStrokeColorDeviceCmykSymbol
}

// Run executes the operation on the given context, setting the stroking color to CMYK.
func (o *SetStrokeColorDeviceCmyk) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetStrokingColorCmyk(o.C, o.M, o.Y, o.K)
}

// Write writes the operator to the given writer.
func (o SetStrokeColorDeviceCmyk) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %g %g", o.C, o.M, o.Y, o.K); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, " "+setStrokeColorDeviceCmykSymbol); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetStrokeColorDeviceCmyk) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetStrokeColorDeviceCmyk) String() string {
	return strings.Join([]string{
		fmt.Sprintf("%g", o.C),
		fmt.Sprintf("%g", o.M),
		fmt.Sprintf("%g", o.Y),
		fmt.Sprintf("%g", o.K),
		setStrokeColorDeviceCmykSymbol,
	}, " ")
}
