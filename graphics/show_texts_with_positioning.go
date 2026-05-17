// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// ShowTextsWithPositioning shows text with individual glyph positioning ("TJ" operator).
// Each element of the array can be a string (shows the string) or a number
// (adjusts the text position by that amount in thousandths of a unit).
type ShowTextsWithPositioning struct {
	Array []tokens.Token
}

const showTextsWithPositioningSymbol = "TJ"

// NewShowTextsWithPositioning creates a new operation from an array of tokens.
func NewShowTextsWithPositioning(array []tokens.Token) (*ShowTextsWithPositioning, error) {
	if array == nil {
		return nil, fmt.Errorf("array must not be nil")
	}

	for i, token := range array {
		switch token.(type) {
		case *tokens.StringToken, *tokens.NumericToken, *tokens.HexToken:
			// valid token types
		default:
			return nil, fmt.Errorf("found invalid token at index %d for showing texts with position: %v", i, token)
		}
	}

	return &ShowTextsWithPositioning{Array: array}, nil
}

// Operator returns the operator symbol for this operation.
func (o ShowTextsWithPositioning) Operator() string {
	return showTextsWithPositioningSymbol
}

// Run executes the show-texts-with-positioning operation on the given context.
func (o *ShowTextsWithPositioning) Run(ctx OperationContext) {
	ctx.ShowPositionedText(o.Array)
}

// writeTJToken writes a single token from a TJ array with proper escaping for content streams.
// StringTokens are escaped per PDF spec (backslash, open/close parens), matching C# TokenWriter behavior.
func writeTJToken(token tokens.Token, w io.Writer) error {
	switch t := token.(type) {
	case *tokens.StringToken:
		if _, err := fmt.Fprint(w, "("); err != nil {
			return err
		}
		for _, r := range t.Data() {
			switch r {
			case '\\':
				if _, err := fmt.Fprint(w, "\\\\"); err != nil {
					return err
				}
			case '(':
				if _, err := fmt.Fprint(w, "\\("); err != nil {
					return err
				}
			case ')':
				if _, err := fmt.Fprint(w, "\\)"); err != nil {
					return err
				}
			default:
				if _, err := fmt.Fprint(w, string(r)); err != nil {
					return err
				}
			}
		}
		if _, err := fmt.Fprint(w, ") "); err != nil {
			return err
		}
	case *tokens.NumericToken:
		if _, err := fmt.Fprintf(w, "%g ", t.Data()); err != nil {
			return err
		}
	case *tokens.HexToken:
		if _, err := fmt.Fprintf(w, "<%X> ", t.GetHexString()); err != nil {
			return err
		}
	default:
		if _, err := fmt.Fprint(w, token); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, " "); err != nil {
			return err
		}
	}
	return nil
}

// Write writes the operator to the given writer.
func (o ShowTextsWithPositioning) Write(w io.Writer) error {
	if _, err := fmt.Fprint(w, "["); err != nil {
		return err
	}

	for _, token := range o.Array {
		if err := writeTJToken(token, w); err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w, "] %s\n", showTextsWithPositioningSymbol)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (ShowTextsWithPositioning) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o ShowTextsWithPositioning) String() string {
	parts := make([]string, len(o.Array))
	for i, token := range o.Array {
		parts[i] = fmt.Sprint(token)
	}
	return fmt.Sprintf("[%s] %s", joinStrings(parts, " "), showTextsWithPositioningSymbol)
}

func joinStrings(ss []string, sep string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

var _ content.GraphicsStateOperation = (*ShowTextsWithPositioning)(nil)
