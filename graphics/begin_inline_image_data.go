// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// BeginInlineImageData begins the image data for an inline image object.
type BeginInlineImageData struct {
	Dictionary map[*tokens.NameToken]tokens.Token
}

// beginInlineImageDataSymbol is the operator symbol for this operation in a PDF content stream.
const beginInlineImageDataSymbol = "ID"

// NewBeginInlineImageData creates a new BeginInlineImageData operation.
func NewBeginInlineImageData(dictionary map[*tokens.NameToken]tokens.Token) *BeginInlineImageData {
	if dictionary == nil {
		return nil
	}
	return &BeginInlineImageData{Dictionary: dictionary}
}

// Operator returns the operator symbol for this operation.
func (o BeginInlineImageData) Operator() string {
	return beginInlineImageDataSymbol
}

// Run executes the operation on the given context, setting inline image properties.
func (o *BeginInlineImageData) Run(ctx OperationContext) {
	ctx.SetInlineImageProperties(o.Dictionary)
}

// Write writes the operator to the given writer.
func (o BeginInlineImageData) Write(w io.Writer) error {
	for name, value := range o.Dictionary {
		if _, err := fmt.Fprintf(w, "%s ", name); err != nil {
			return err
		}
		if str, ok := value.(fmt.Stringer); ok {
			if _, err := fmt.Fprint(w, str.String()); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprint(w, value); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(w, beginInlineImageDataSymbol)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (BeginInlineImageData) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o BeginInlineImageData) String() string {
	return beginInlineImageDataSymbol
}
