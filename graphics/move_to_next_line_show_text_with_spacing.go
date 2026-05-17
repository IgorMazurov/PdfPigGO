// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
)

// MoveToNextLineShowTextWithSpacing moves to the next line with spacing and shows a text string ("\"" operator).
type MoveToNextLineShowTextWithSpacing struct {
	WordSpacing      float64
	CharacterSpacing float64
	Text             string
	Bytes            []byte
}

const moveToNextLineShowTextWithSpacingSymbol = "\""

// NewMoveToNextLineShowTextWithSpacing creates a new operation from a text string.
func NewMoveToNextLineShowTextWithSpacing(wordSpacing, characterSpacing float64, text string) *MoveToNextLineShowTextWithSpacing {
	return &MoveToNextLineShowTextWithSpacing{
		WordSpacing:      wordSpacing,
		CharacterSpacing: characterSpacing,
		Text:             text,
	}
}

// NewMoveToNextLineShowTextWithSpacingBytes creates a new operation from hex bytes.
func NewMoveToNextLineShowTextWithSpacingBytes(wordSpacing, characterSpacing float64, hexBytes []byte) *MoveToNextLineShowTextWithSpacing {
	result := make([]byte, len(hexBytes))
	copy(result, hexBytes)
	return &MoveToNextLineShowTextWithSpacing{
		WordSpacing:      wordSpacing,
		CharacterSpacing: characterSpacing,
		Bytes:            result,
	}
}

// Operator returns the operator symbol for this operation.
func (o MoveToNextLineShowTextWithSpacing) Operator() string {
	return moveToNextLineShowTextWithSpacingSymbol
}

// Run executes the move-to-next-line-with-spacing-and-show-text operation on the given context.
// Sets word spacing, character spacing, moves to next line, then shows text.
func (o *MoveToNextLineShowTextWithSpacing) Run(ctx OperationContext) {
	ctx.SetWordSpacing(o.WordSpacing)
	ctx.SetCharacterSpacing(o.CharacterSpacing)
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
func (o MoveToNextLineShowTextWithSpacing) Write(w io.Writer) error {
	if o.Text != "" {
		_, err := fmt.Fprintf(w, "%g %g (%s) %s\n", o.WordSpacing, o.CharacterSpacing, o.Text, moveToNextLineShowTextWithSpacingSymbol)
		return err
	}

	hexStr := fmt.Sprintf("<%X>", o.Bytes)
	_, err := fmt.Fprintf(w, "%g %g %s %s\n", o.WordSpacing, o.CharacterSpacing, hexStr, moveToNextLineShowTextWithSpacingSymbol)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (MoveToNextLineShowTextWithSpacing) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o MoveToNextLineShowTextWithSpacing) String() string {
	if o.Text != "" {
		return fmt.Sprintf("%g %g %s %s", o.WordSpacing, o.CharacterSpacing, o.Text, moveToNextLineShowTextWithSpacingSymbol)
	}
	return fmt.Sprintf("%g %g <%X> %s", o.WordSpacing, o.CharacterSpacing, o.Bytes, moveToNextLineShowTextWithSpacingSymbol)
}

var _ content.GraphicsStateOperation = (*MoveToNextLineShowTextWithSpacing)(nil)
