// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// SetNonStrokeColorSpace sets the current color space for non-stroking operations.
type SetNonStrokeColorSpace struct {
	Name *tokens.NameToken
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setNonStrokeColorSpaceSymbol = "cs"

// NewSetNonStrokeColorSpace creates a new SetNonStrokeColorSpace operation.
func NewSetNonStrokeColorSpace(name *tokens.NameToken) *SetNonStrokeColorSpace {
	return &SetNonStrokeColorSpace{Name: name}
}

// Operator returns the operator symbol for this operation.
func (o SetNonStrokeColorSpace) Operator() string {
	return setNonStrokeColorSpaceSymbol
}

// Run executes the operation on the given context, setting the non-stroking color space.
func (o *SetNonStrokeColorSpace) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetNonStrokingColorspace(o.Name, nil)
}

// Write writes the operator to the given writer.
func (o SetNonStrokeColorSpace) Write(w io.Writer) error {
	if _, err := io.WriteString(w, "/"+o.Name.Data()); err != nil {
		return err
	}
	if _, err := io.WriteString(w, " "); err != nil {
		return err
	}
	if _, err := io.WriteString(w, setNonStrokeColorSpaceSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetNonStrokeColorSpace) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetNonStrokeColorSpace) String() string {
	return o.Name.Data() + " " + setNonStrokeColorSpaceSymbol
}
