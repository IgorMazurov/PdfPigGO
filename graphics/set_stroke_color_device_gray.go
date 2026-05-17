// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"
)

// SetStrokeColorDeviceGray sets the stroking color to a gray level.
type SetStrokeColorDeviceGray struct {
	Gray float64
}

const setStrokeColorDeviceGraySymbol = "G"

// NewSetStrokeColorDeviceGray creates a new SetStrokeColorDeviceGray operation.
func NewSetStrokeColorDeviceGray(gray float64) *SetStrokeColorDeviceGray {
	return &SetStrokeColorDeviceGray{Gray: gray}
}

// Operator returns the operator symbol for this operation.
func (o SetStrokeColorDeviceGray) Operator() string {
	return setStrokeColorDeviceGraySymbol
}

// Run executes the operation on the given context, setting the stroking color to gray.
func (o *SetStrokeColorDeviceGray) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetStrokingColorGray(o.Gray)
}

// Write writes the operator to the given writer.
func (o SetStrokeColorDeviceGray) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g", o.Gray); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, " "+setStrokeColorDeviceGraySymbol); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetStrokeColorDeviceGray) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetStrokeColorDeviceGray) String() string {
	return strings.Join([]string{
		fmt.Sprintf("%g", o.Gray),
		setStrokeColorDeviceGraySymbol,
	}, " ")
}
