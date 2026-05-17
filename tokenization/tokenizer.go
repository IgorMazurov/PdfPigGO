package tokenization

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Tokenizer reads tokens from input data.
type Tokenizer interface {
	// ReadsNextByte reports whether this tokenizer type reads the byte
	// following the token itself to detect if the token has ended.
	ReadsNextByte() bool

	// Tokenize tries to read a token of the corresponding type from the input.
	// The currentByte parameter is the byte that detected this is the correct
	// tokenizer to use. Returns the token and true on success, or nil and false
	// if tokenization failed.
	Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool)
}
