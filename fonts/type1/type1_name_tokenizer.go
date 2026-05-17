package type1

import (
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Type1NameTokenizer tokenizes PostScript name objects starting with '/'.
// It reads bytes after the leading slash until a delimiter character is found,
// producing a NameToken.
type Type1NameTokenizer struct{}

var _ tokenization.Tokenizer = (*Type1NameTokenizer)(nil)

// NewType1NameTokenizer creates a new Type1NameTokenizer.
func NewType1NameTokenizer() *Type1NameTokenizer {
	return &Type1NameTokenizer{}
}

// ReadsNextByte returns true because this tokenizer reads ahead past the name
// to detect its end boundary.
func (t *Type1NameTokenizer) ReadsNextByte() bool {
	return true
}

// Tokenize attempts to tokenize a PostScript name starting with '/'.
// It reads bytes until a whitespace or delimiter character ({, <, /, [, () is found,
// producing an interned NameToken. Returns the token and true on success,
// or nil and false if currentByte is not '/'.
func (t *Type1NameTokenizer) Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool) {
	if currentByte != '/' {
		return nil, false
	}

	var builder strings.Builder

	for input.MoveNext() {
		b := input.CurrentByte()
		if core.IsWhitespace(b) || b == '{' || b == '<' || b == '/' || b == '[' || b == '(' {
			break
		}
		builder.WriteByte(b)
	}

	return tokens.Create(builder.String()), true
}
