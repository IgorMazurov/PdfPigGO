package tokenization

import (
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// PlainTokenizer tokenizes plain text tokens such as booleans, null, and operators.
type PlainTokenizer struct{}

var _ Tokenizer = (*PlainTokenizer)(nil)

// NewPlainTokenizer creates a new PlainTokenizer.
func NewPlainTokenizer() *PlainTokenizer {
	return &PlainTokenizer{}
}

// ReadsNextByte returns true because this tokenizer reads ahead past the token.
func (t *PlainTokenizer) ReadsNextByte() bool {
	return true
}

// Tokenize attempts to tokenize a plain text token (boolean, null, or operator).
// Returns the token and true on success, or nil and false if currentByte is whitespace.
func (t *PlainTokenizer) Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool) {
	if core.IsWhitespace(currentByte) {
		return nil, false
	}

	var builder strings.Builder
	builder.Grow(16)
	builder.WriteByte(currentByte)

	for input.MoveNext() {
		b := input.CurrentByte()

		if core.IsWhitespace(b) {
			break
		}

		// Use if statements instead of switch to ensure break exits the for loop.
		// In Go, 'break' inside a switch only breaks the switch, not the enclosing loop.
		if b == '<' || b == '[' || b == '/' || b == ']' || b == '>' || b == '(' || b == ')' {
			break
		}

		builder.WriteByte(b)
	}

	text := builder.String()

	switch text {
	case "true":
		return tokens.True, true
	case "false":
		return tokens.False, true
	case "null":
		return tokens.Null, true
	default:
		return tokens.CreateOperatorToken(text), true
	}
}
