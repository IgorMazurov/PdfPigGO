// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"
)

// SetNonStrokeColorDeviceRgb sets the non-stroking color space to DeviceRGB and set the color.
type SetNonStrokeColorDeviceRgb struct {
	R float64
	G float64
	B float64
}

const setNonStrokeColorDeviceRgbSymbol = "rg"

// NewSetNonStrokeColorDeviceRgb creates a new SetNonStrokeColorDeviceRgb operation.
func NewSetNonStrokeColorDeviceRgb(r, g, b float64) *SetNonStrokeColorDeviceRgb {
	return &SetNonStrokeColorDeviceRgb{R: r, G: g, B: b}
}

// Operator returns the operator symbol for this operation.
func (o SetNonStrokeColorDeviceRgb) Operator() string {
	return setNonStrokeColorDeviceRgbSymbol
}

// Run executes the operation on the given context, setting the non-stroking color to RGB.
func (o *SetNonStrokeColorDeviceRgb) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetNonStrokingColorRgb(o.R, o.G, o.B)
}

// Write writes the operator to the given writer.
func (o SetNonStrokeColorDeviceRgb) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %g", o.R, o.G, o.B); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, " "+setNonStrokeColorDeviceRgbSymbol); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetNonStrokeColorDeviceRgb) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetNonStrokeColorDeviceRgb) String() string {
	return strings.Join([]string{
		fmt.Sprintf("%g", o.R),
		fmt.Sprintf("%g", o.G),
		fmt.Sprintf("%g", o.B),
		setNonStrokeColorDeviceRgbSymbol,
	}, " ")
}
