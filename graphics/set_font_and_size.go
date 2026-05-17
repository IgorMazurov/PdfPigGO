// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// SetFontAndSize sets the font and the font size ("Tf" operator).
type SetFontAndSize struct {
	Font *tokens.NameToken
	Size float64
}

const setFontAndSizeSymbol = "Tf"

// NewSetFontAndSize creates a new SetFontAndSize operation.
func NewSetFontAndSize(font *tokens.NameToken, size float64) *SetFontAndSize {
	if font == nil {
		panic("font cannot be nil")
	}
	return &SetFontAndSize{Font: font, Size: size}
}

// Operator returns the operator symbol for this operation.
func (o SetFontAndSize) Operator() string {
	return setFontAndSizeSymbol
}

// Run executes the operation on the given context, setting the font and size.
func (o *SetFontAndSize) Run(ctx OperationContext) {
	ctx.SetFontAndSize(o.Font, o.Size)
}

// Write writes the operator to the given writer.
func (o SetFontAndSize) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "/%s %g %s\n", o.Font.Data(), o.Size, setFontAndSizeSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetFontAndSize) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetFontAndSize) String() string {
	return fmt.Sprintf("/%s %g %s", o.Font.Data(), o.Size, setFontAndSizeSymbol)
}

var _ content.GraphicsStateOperation = (*SetFontAndSize)(nil)
