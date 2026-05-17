// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// DesignateMarkedContentPoint designates a single marked-content point in the content stream.
type DesignateMarkedContentPoint struct {
	Name *tokens.NameToken
}

const designateMarkedContentPointSymbol = "MP"

// NewDesignateMarkedContentPoint creates a new DesignateMarkedContentPoint operation.
func NewDesignateMarkedContentPoint(name *tokens.NameToken) (*DesignateMarkedContentPoint, error) {
	if name == nil {
		return nil, fmt.Errorf("name cannot be nil")
	}
	return &DesignateMarkedContentPoint{Name: name}, nil
}

// Operator returns the operator symbol for this operation.
func (o DesignateMarkedContentPoint) Operator() string {
	return designateMarkedContentPointSymbol
}

// Run executes the operation on the given context. This is a no-op as per PDF spec.
func (o *DesignateMarkedContentPoint) Run(ctx OperationContext) {
}

// Write writes the operator to the given writer.
func (o DesignateMarkedContentPoint) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s %s\n", o.Name, designateMarkedContentPointSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (DesignateMarkedContentPoint) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o DesignateMarkedContentPoint) String() string {
	return fmt.Sprintf("%s %s", o.Name, designateMarkedContentPointSymbol)
}
