// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// InvokeNamedXObject paints the specified XObject. The operand name must appear as a key
// in the XObject subdictionary of the current resource dictionary. The associated value
// must be a stream whose Type entry, if present, is XObject. The effect depends on the
// value of the XObject's Subtype entry, which may be Image, Form or PS.
type InvokeNamedXObject struct {
	Name *tokens.NameToken
}

const invokeNamedXObjectSymbol = "Do"

// NewInvokeNamedXObject creates a new InvokeNamedXObject operation.
func NewInvokeNamedXObject(name *tokens.NameToken) (*InvokeNamedXObject, error) {
	if name == nil {
		return nil, fmt.Errorf("name cannot be nil")
	}
	return &InvokeNamedXObject{Name: name}, nil
}

// Operator returns the operator symbol for this operation.
func (o InvokeNamedXObject) Operator() string {
	return invokeNamedXObjectSymbol
}

// Run executes the operation on the given context, applying the XObject.
func (o *InvokeNamedXObject) Run(ctx OperationContext) {
	ctx.ApplyXObject(o.Name)
}

// Write writes the operator to the given writer.
func (o InvokeNamedXObject) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "/%s %s\n", o.Name.Data(), invokeNamedXObjectSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (InvokeNamedXObject) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o InvokeNamedXObject) String() string {
	return fmt.Sprintf("%s %s", o.Name, invokeNamedXObjectSymbol)
}
