// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"
)

// BeginInlineImage begins an inline image object.
type BeginInlineImage struct{}

// Symbol is the operator symbol for this operation in a PDF content stream.
const beginInlineImageSymbol = "BI"

// InstanceBeginInlineImage is the singleton instance of BeginInlineImage.
var InstanceBeginInlineImage = BeginInlineImage{}

// Operator returns the operator symbol for this operation.
func (o BeginInlineImage) Operator() string {
	return beginInlineImageSymbol
}

// Run executes the operation on the given context, starting inline image collection.
func (o BeginInlineImage) Run(ctx OperationContext) {
	ctx.BeginInlineImage()
}

// Write writes the operator symbol to the given writer.
func (o BeginInlineImage) Write(w io.Writer) error {
	if _, err := io.WriteString(w, beginInlineImageSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (BeginInlineImage) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o BeginInlineImage) String() string {
	return beginInlineImageSymbol
}
