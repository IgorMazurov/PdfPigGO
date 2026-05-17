// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// PaintShading paints the shape and color shading described by a shading dictionary,
// subject to the current clipping path. The current color in the graphics state is
// neither used nor altered. The effect is different from that of painting a path
// using a shading pattern as the current color.
type PaintShading struct {
	Name *tokens.NameToken
}

const paintShadingSymbol = "sh"

// NewPaintShading creates a new PaintShading operation.
func NewPaintShading(name *tokens.NameToken) (*PaintShading, error) {
	if name == nil {
		return nil, fmt.Errorf("name cannot be nil")
	}
	return &PaintShading{Name: name}, nil
}

// Operator returns the operator symbol for this operation.
func (o PaintShading) Operator() string {
	return paintShadingSymbol
}

// Run executes the operation on the given context, painting the shading fill.
func (o *PaintShading) Run(ctx OperationContext) {
	ctx.PaintShading(o.Name)
}

// Write writes the operator to the given writer.
func (o PaintShading) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "/%s %s\n", o.Name.Data(), paintShadingSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (PaintShading) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o PaintShading) String() string {
	return fmt.Sprintf("%s %s", o.Name, paintShadingSymbol)
}
