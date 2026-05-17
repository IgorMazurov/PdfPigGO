// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
)

// ShowText shows a text string ("Tj" operator).
type ShowText struct {
	Text  string
	Bytes []byte
}

const showTextSymbol = "Tj"

// NewShowText creates a new ShowText from a string.
func NewShowText(text string) *ShowText {
	return &ShowText{Text: text}
}

// NewShowTextBytes creates a new ShowText from hex bytes.
func NewShowTextBytes(hexBytes []byte) *ShowText {
	result := make([]byte, len(hexBytes))
	copy(result, hexBytes)
	return &ShowText{Bytes: result}
}

// Operator returns the operator symbol for this operation.
func (o ShowText) Operator() string {
	return showTextSymbol
}

// Run executes the show-text operation on the given context.
func (o *ShowText) Run(ctx OperationContext) {
	var bytes core.InputBytes
	if o.Text != "" {
		bytes = core.NewMemoryInputBytes(core.StringAsLatin1Bytes(o.Text))
	} else {
		bytes = core.NewMemoryInputBytes(o.Bytes)
	}
	ctx.ShowText(bytes)
}

// Write writes the operator to the given writer.
func (o ShowText) Write(w io.Writer) error {
	if len(o.Bytes) > 0 {
		hexStr := fmt.Sprintf("<%X>", o.Bytes)
		if _, err := fmt.Fprintf(w, "%s %s\n", hexStr, showTextSymbol); err != nil {
			return err
		}
		return nil
	}

	escaped := escapeText(o.Text)
	_, err := fmt.Fprintf(w, "(%s) %s\n", escaped, showTextSymbol)
	return err
}

// escapeText escapes special characters '(', ')' and '\\' per PDF spec 1.7 page 408.
func escapeText(text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, "(", "\\(")
	text = strings.ReplaceAll(text, ")", "\\)")
	return text
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (ShowText) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o ShowText) String() string {
	if o.Text != "" {
		return fmt.Sprintf("%s %s", o.Text, showTextSymbol)
	}
	return fmt.Sprintf("<%X> %s", o.Bytes, showTextSymbol)
}

var _ content.GraphicsStateOperation = (*ShowText)(nil)
