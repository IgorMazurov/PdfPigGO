// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"
)

// SetStrokeColorDeviceRgb sets the stroking color space to DeviceRGB and set the color.
type SetStrokeColorDeviceRgb struct {
	R float64
	G float64
	B float64
}

const setStrokeColorDeviceRgbSymbol = "RG"

// NewSetStrokeColorDeviceRgb creates a new SetStrokeColorDeviceRgb operation.
func NewSetStrokeColorDeviceRgb(r, g, b float64) *SetStrokeColorDeviceRgb {
	return &SetStrokeColorDeviceRgb{R: r, G: g, B: b}
}

// Operator returns the operator symbol for this operation.
func (o SetStrokeColorDeviceRgb) Operator() string {
	return setStrokeColorDeviceRgbSymbol
}

// Run executes the operation on the given context, setting the stroking color to RGB.
func (o *SetStrokeColorDeviceRgb) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetStrokingColorRgb(o.R, o.G, o.B)
}

// Write writes the operator to the given writer.
func (o SetStrokeColorDeviceRgb) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %g", o.R, o.G, o.B); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, " "+setStrokeColorDeviceRgbSymbol); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetStrokeColorDeviceRgb) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetStrokeColorDeviceRgb) String() string {
	return strings.Join([]string{
		fmt.Sprintf("%g", o.R),
		fmt.Sprintf("%g", o.G),
		fmt.Sprintf("%g", o.B),
		setStrokeColorDeviceRgbSymbol,
	}, " ")
}
