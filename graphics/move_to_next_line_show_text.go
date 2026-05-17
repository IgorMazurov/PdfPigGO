// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
)

// MoveToNextLineShowText moves to the next line and shows a text string ("'" operator).
type MoveToNextLineShowText struct {
	Text  string
	Bytes []byte
}

const moveToNextLineShowTextSymbol = "'"

// NewMoveToNextLineShowText creates a new MoveToNextLineShowText from a string.
func NewMoveToNextLineShowText(text string) *MoveToNextLineShowText {
	return &MoveToNextLineShowText{Text: text}
}

// NewMoveToNextLineShowTextBytes creates a new MoveToNextLineShowText from hex bytes.
func NewMoveToNextLineShowTextBytes(hexBytes []byte) *MoveToNextLineShowText {
	result := make([]byte, len(hexBytes))
	copy(result, hexBytes)
	return &MoveToNextLineShowText{Bytes: result}
}

// Operator returns the operator symbol for this operation.
func (o MoveToNextLineShowText) Operator() string {
	return moveToNextLineShowTextSymbol
}

// Run executes the move-to-next-line-and-show-text operation on the given context.
// First moves to the next line, then shows the text.
func (o *MoveToNextLineShowText) Run(ctx OperationContext) {
	ctx.MoveToNextLineWithOffset()

	var bytes core.InputBytes
	if o.Text != "" {
		bytes = core.NewMemoryInputBytes(core.StringAsLatin1Bytes(o.Text))
	} else {
		bytes = core.NewMemoryInputBytes(o.Bytes)
	}
	ctx.ShowText(bytes)
}

// Write writes the operator to the given writer.
func (o MoveToNextLineShowText) Write(w io.Writer) error {
	if o.Text != "" {
		_, err := fmt.Fprintf(w, "(%s) %s\n", o.Text, moveToNextLineShowTextSymbol)
		return err
	}

	hexStr := fmt.Sprintf("<%X>", o.Bytes)
	_, err := fmt.Fprintf(w, "%s %s\n", hexStr, moveToNextLineShowTextSymbol)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (MoveToNextLineShowText) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o MoveToNextLineShowText) String() string {
	if o.Text != "" {
		return fmt.Sprintf("%s %s", o.Text, moveToNextLineShowTextSymbol)
	}
	return fmt.Sprintf("<%X> %s", o.Bytes, moveToNextLineShowTextSymbol)
}

var _ content.GraphicsStateOperation = (*MoveToNextLineShowText)(nil)
