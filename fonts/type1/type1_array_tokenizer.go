package type1

import (
	"strconv"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Type1ArrayTokenizer tokenizes PostScript-style arrays delimited by curly braces ({...}).
// It reads the content between { and }, splits on whitespace, and produces an ArrayToken
// containing NumericTokens, NameTokens, StringTokens, or OperatorTokens.
type Type1ArrayTokenizer struct{}

var _ tokenization.Tokenizer = (*Type1ArrayTokenizer)(nil)

// NewType1ArrayTokenizer creates a new Type1ArrayTokenizer.
func NewType1ArrayTokenizer() *Type1ArrayTokenizer {
	return &Type1ArrayTokenizer{}
}

// ReadsNextByte returns false because this tokenizer consumes the closing brace itself.
func (t *Type1ArrayTokenizer) ReadsNextByte() bool {
	return false
}

// Tokenize attempts to tokenize a PostScript array starting with '{'.
// It reads bytes until '}' is found, splits the content by whitespace, and converts
// each part into the appropriate token type. Returns an ArrayToken and true on success,
// or nil and false if currentByte is not '{'.
func (t *Type1ArrayTokenizer) Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool) {
	if currentByte != '{' {
		return nil, false
	}

	var builder strings.Builder

	for input.MoveNext() {
		b := input.CurrentByte()
		if b == '}' {
			break
		}
		builder.WriteByte(b)
	}

	parts := strings.Fields(builder.String())

	tokensList := make([]tokens.Token, 0, len(parts))

	for _, part := range parts {
		firstByte := part[0]

		if (firstByte >= '0' && firstByte <= '9') || firstByte == '-' {
			if value, err := strconv.ParseFloat(part, 64); err == nil {
				tokensList = append(tokensList, tokens.NewNumericToken(value))
			} else {
				tokensList = append(tokensList, tokens.CreateOperatorToken(part))
			}
			continue
		}

		if firstByte == '/' {
			tokensList = append(tokensList, tokens.Create(part[1:]))
			continue
		}

		if len(part) >= 2 && part[0] == '(' && part[len(part)-1] == ')' {
			tokensList = append(tokensList, tokens.NewStringToken(part))
			continue
		}

		tokensList = append(tokensList, tokens.CreateOperatorToken(part))
	}

	return tokens.NewArrayToken(tokensList), true
}
