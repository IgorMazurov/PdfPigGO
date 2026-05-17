// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"
)

// SetNonStrokeColorDeviceGray sets the non-stroking color to a gray level.
type SetNonStrokeColorDeviceGray struct {
	Gray float64
}

const setNonStrokeColorDeviceGraySymbol = "g"

// NewSetNonStrokeColorDeviceGray creates a new SetNonStrokeColorDeviceGray operation.
func NewSetNonStrokeColorDeviceGray(gray float64) *SetNonStrokeColorDeviceGray {
	return &SetNonStrokeColorDeviceGray{Gray: gray}
}

// Operator returns the operator symbol for this operation.
func (o SetNonStrokeColorDeviceGray) Operator() string {
	return setNonStrokeColorDeviceGraySymbol
}

// Run executes the operation on the given context, setting the non-stroking color to gray.
func (o *SetNonStrokeColorDeviceGray) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetNonStrokingColorGray(o.Gray)
}

// Write writes the operator to the given writer.
func (o SetNonStrokeColorDeviceGray) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g", o.Gray); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, " "+setNonStrokeColorDeviceGraySymbol); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetNonStrokeColorDeviceGray) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetNonStrokeColorDeviceGray) String() string {
	return strings.Join([]string{
		fmt.Sprintf("%g", o.Gray),
		setNonStrokeColorDeviceGraySymbol,
	}, " ")
}
