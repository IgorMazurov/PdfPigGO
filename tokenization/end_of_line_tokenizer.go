package tokenization

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// EndOfLineTokenizer tokenizes end-of-line markers (\r or \n).
type EndOfLineTokenizer struct{}

var _ Tokenizer = (*EndOfLineTokenizer)(nil)

// NewEndOfLineTokenizer creates a new EndOfLineTokenizer.
func NewEndOfLineTokenizer() *EndOfLineTokenizer {
	return &EndOfLineTokenizer{}
}

// ReadsNextByte returns false because this tokenizer does not read ahead.
func (t *EndOfLineTokenizer) ReadsNextByte() bool {
	return false
}

// Tokenize attempts to tokenize an end-of-line marker.
// Returns the EOL singleton token and true if currentByte is \r or \n,
// or nil and false otherwise.
func (t *EndOfLineTokenizer) Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool) {
	if currentByte != '\r' && currentByte != '\n' {
		return nil, false
	}

	return tokens.EOLToken, true
}
