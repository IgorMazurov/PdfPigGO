// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// SetGraphicsStateParametersFromDictionary sets the specified parameters in the graphics state
// using the ExtGState subdictionary with the given name ("gs" operator).
type SetGraphicsStateParametersFromDictionary struct {
	Name *tokens.NameToken
}

// symbolSetGraphicsState is the operator symbol for this operation in a PDF content stream.
const symbolSetGraphicsState = "gs"

// NewSetGraphicsStateParametersFromDictionary creates a new SetGraphicsStateParametersFromDictionary operation.
func NewSetGraphicsStateParametersFromDictionary(name *tokens.NameToken) *SetGraphicsStateParametersFromDictionary {
	return &SetGraphicsStateParametersFromDictionary{Name: name}
}

// Operator returns the operator symbol for this operation.
func (o SetGraphicsStateParametersFromDictionary) Operator() string {
	return symbolSetGraphicsState
}

// Run executes the operation on the given context, setting named graphics state parameters.
func (o *SetGraphicsStateParametersFromDictionary) Run(ctx OperationContext) {
	ctx.SetNamedGraphicsState(o.Name)
}

// Write writes the operator to the given writer.
func (o SetGraphicsStateParametersFromDictionary) Write(w io.Writer) error {
	if _, err := io.WriteString(w, "/"+o.Name.Data()); err != nil {
		return err
	}
	if _, err := io.WriteString(w, " "); err != nil {
		return err
	}
	if _, err := io.WriteString(w, symbolSetGraphicsState); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetGraphicsStateParametersFromDictionary) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetGraphicsStateParametersFromDictionary) String() string {
	return o.Name.Data() + " " + symbolSetGraphicsState
}
