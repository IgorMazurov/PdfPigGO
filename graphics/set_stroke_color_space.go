// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// SetStrokeColorSpace sets the current color space for stroking operations.
type SetStrokeColorSpace struct {
	Name *tokens.NameToken
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setStrokeColorSpaceSymbol = "CS"

// NewSetStrokeColorSpace creates a new SetStrokeColorSpace operation.
func NewSetStrokeColorSpace(name *tokens.NameToken) *SetStrokeColorSpace {
	return &SetStrokeColorSpace{Name: name}
}

// Operator returns the operator symbol for this operation.
func (o SetStrokeColorSpace) Operator() string {
	return setStrokeColorSpaceSymbol
}

// Run executes the operation on the given context, setting the stroking color space.
func (o *SetStrokeColorSpace) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetStrokingColorspace(o.Name, nil)
}

// Write writes the operator to the given writer.
func (o SetStrokeColorSpace) Write(w io.Writer) error {
	if _, err := io.WriteString(w, "/"+o.Name.Data()); err != nil {
		return err
	}
	if _, err := io.WriteString(w, " "); err != nil {
		return err
	}
	if _, err := io.WriteString(w, setStrokeColorSpaceSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetStrokeColorSpace) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetStrokeColorSpace) String() string {
	return o.Name.Data() + " " + setStrokeColorSpaceSymbol
}
