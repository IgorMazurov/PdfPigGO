package tokenization

import (
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CommentTokenizer tokenizes PDF comments starting with '%'.
type CommentTokenizer struct{}

var _ Tokenizer = (*CommentTokenizer)(nil)

// NewCommentTokenizer creates a new CommentTokenizer.
func NewCommentTokenizer() *CommentTokenizer {
	return &CommentTokenizer{}
}

// ReadsNextByte returns true because this tokenizer reads ahead past the comment.
func (t *CommentTokenizer) ReadsNextByte() bool {
	return true
}

// Tokenize attempts to tokenize a comment starting with '%'.
// Returns the CommentToken and true on success, or nil and false if currentByte is not '%'.
func (t *CommentTokenizer) Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool) {
	if currentByte != '%' {
		return nil, false
	}

	var builder strings.Builder

	for input.MoveNext() && !core.IsEndOfLineByte(input.CurrentByte()) {
		builder.WriteByte(input.CurrentByte())
	}

	return tokens.NewCommentToken(builder.String()), true
}
