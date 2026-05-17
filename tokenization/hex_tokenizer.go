package tokenization

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// HexTokenizer tokenizes PDF hexadecimal strings enclosed in angle brackets.
type HexTokenizer struct{}

var _ Tokenizer = (*HexTokenizer)(nil)

// NewHexTokenizer creates a new HexTokenizer.
func NewHexTokenizer() *HexTokenizer {
	return &HexTokenizer{}
}

// ReadsNextByte returns false because this tokenizer does not read ahead.
func (t *HexTokenizer) ReadsNextByte() bool {
	return false
}

// Tokenize attempts to tokenize a hexadecimal string starting with '<'.
// Returns the HexToken and true on success, or nil and false if currentByte is not '<'
// or an invalid hex character is encountered.
func (t *HexTokenizer) Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool) {
	if currentByte != '<' {
		return nil, false
	}

	var chars []rune

	for input.MoveNext() {
		current := input.CurrentByte()

		if core.IsWhitespace(current) {
			continue
		}

		if current == '>' {
			break
		}

		if !isValidHexChar(current) {
			return nil, false
		}

		chars = append(chars, rune(current))
	}

	return tokens.NewHexToken(chars), true
}

func isValidHexChar(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
}
