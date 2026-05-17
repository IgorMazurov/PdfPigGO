// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"
)

// EndInlineImage ends an inline image object.
type EndInlineImage struct {
	// ImageData is the raw byte data of the inline image.
	ImageData []byte
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const endInlineImageSymbol = "EI"

// NewEndInlineImage creates a new EndInlineImage operation with the given image data.
func NewEndInlineImage(imageData []byte) *EndInlineImage {
	return &EndInlineImage{ImageData: imageData}
}

// Operator returns the operator symbol for this operation.
func (o EndInlineImage) Operator() string {
	return endInlineImageSymbol
}

// Run executes the operation on the given context, ending inline image collection.
func (o EndInlineImage) Run(ctx OperationContext) {
	ctx.EndInlineImage(o.ImageData)
}

// Write writes the operator symbol to the given writer.
func (o EndInlineImage) Write(w io.Writer) error {
	if _, err := io.WriteString(w, endInlineImageSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (EndInlineImage) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o EndInlineImage) String() string {
	return endInlineImageSymbol
}
