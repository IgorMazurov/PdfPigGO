// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// BeginMarkedContent begins a marked-content sequence terminated by a balancing EndMarkedContent operator.
type BeginMarkedContent struct {
	Name *tokens.NameToken
}

const beginMarkedContentSymbol = "BMC"

// NewBeginMarkedContent creates a new BeginMarkedContent operation.
func NewBeginMarkedContent(name *tokens.NameToken) (*BeginMarkedContent, error) {
	if name == nil {
		return nil, fmt.Errorf("name cannot be nil")
	}
	return &BeginMarkedContent{Name: name}, nil
}

// Operator returns the operator symbol for this operation.
func (o BeginMarkedContent) Operator() string {
	return beginMarkedContentSymbol
}

// Run executes the operation on the given context, starting a marked content section.
func (o *BeginMarkedContent) Run(ctx OperationContext) {
	ctx.BeginMarkedContent(o.Name, nil, nil)
}

// Write writes the operator to the given writer.
func (o BeginMarkedContent) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s %s\n", o.Name, beginMarkedContentSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (BeginMarkedContent) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o BeginMarkedContent) String() string {
	return fmt.Sprintf("%s %s", o.Name, beginMarkedContentSymbol)
}
