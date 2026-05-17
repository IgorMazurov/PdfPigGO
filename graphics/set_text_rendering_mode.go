// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// SetTextRenderingMode sets the text rendering mode.
type SetTextRenderingMode struct {
	Mode core.TextRenderingMode
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setTextRenderingModeSymbol = "Tr"

// NewSetTextRenderingMode creates a new SetTextRenderingMode operation from an int value.
func NewSetTextRenderingMode(mode int) *SetTextRenderingMode {
	return &SetTextRenderingMode{Mode: core.TextRenderingMode(mode)}
}

// Operator returns the operator symbol for this operation.
func (o SetTextRenderingMode) Operator() string {
	return setTextRenderingModeSymbol
}

// Run executes the operation on the given context, setting the text rendering mode.
func (o *SetTextRenderingMode) Run(ctx OperationContext) {
	ctx.SetTextRenderingMode(o.Mode)
}

// Write writes the operator to the given writer.
func (o SetTextRenderingMode) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%d %s\n", o.Mode, setTextRenderingModeSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetTextRenderingMode) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetTextRenderingMode) String() string {
	return fmt.Sprintf("%d %s", o.Mode, setTextRenderingModeSymbol)
}
