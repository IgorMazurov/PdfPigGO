// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"

	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// SetColorRenderingIntent sets the color rendering intent in the graphics state.
type SetColorRenderingIntent struct {
	RenderingIntent *tokens.NameToken
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setColorRenderingIntentSymbol = "ri"

// NewSetColorRenderingIntent creates a new SetColorRenderingIntent operation.
func NewSetColorRenderingIntent(renderingIntent *tokens.NameToken) *SetColorRenderingIntent {
	return &SetColorRenderingIntent{RenderingIntent: renderingIntent}
}

// Operator returns the operator symbol for this operation.
func (o SetColorRenderingIntent) Operator() string {
	return setColorRenderingIntentSymbol
}

// Run executes the operation on the given context, setting the color rendering intent.
func (o *SetColorRenderingIntent) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil {
		return
	}
	state.RenderingIntent = graphiccore.ParseRenderingIntent(o.RenderingIntent.Data())
}

// Write writes the operator to the given writer.
func (o SetColorRenderingIntent) Write(w io.Writer) error {
	if _, err := io.WriteString(w, "/"+o.RenderingIntent.Data()); err != nil {
		return err
	}
	if _, err := io.WriteString(w, " "); err != nil {
		return err
	}
	if _, err := io.WriteString(w, setColorRenderingIntentSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetColorRenderingIntent) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetColorRenderingIntent) String() string {
	return o.RenderingIntent.Data() + " " + setColorRenderingIntentSymbol
}
