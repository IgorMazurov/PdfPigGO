package tokenization

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestCommentInvalidFirstCharacterReturnsFalse(t *testing.T) {
	testCases := []string{"(%not a comment)", "\\%not a comment)", "\u2030"}

	tokenizer := NewCommentTokenizer()

	for _, s := range testCases {
		first, input := stringInput(s)
		token, ok := tokenizer.Tokenize(first, input)

		if ok {
			t.Errorf("expected false for first byte %q (%c), got true", s, s[0])
		}
		if token != nil {
			t.Errorf("expected nil token for first byte %q, got %v", s, token)
		}
	}
}

func TestCommentParsesComment(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"%Resource-CMAP\n%AnotherComment", "Resource-CMAP"},
		{"%%PDF 1.5", "%PDF 1.5"},
		{"% comment {/%) blah blah blah\n            123", " comment {/%) blah blah blah"},
		{"%comment\rNot comment", "comment"},
		{"%comment\r\nNot comment", "comment"},
		{"%comment\nNot comment", "comment"},
	}

	tokenizer := NewCommentTokenizer()

	for _, tc := range testCases {
		first, input := stringInput(tc.input)
		token, ok := tokenizer.Tokenize(first, input)

		if !ok {
			t.Fatalf("expected true for %q, got false", tc.input)
		}

		comment, ok := token.(*tokens.CommentToken)
		if !ok {
			t.Fatalf("for %q: expected *tokens.CommentToken, got %T", tc.input, token)
		}

		if comment.Data() != tc.expected {
			t.Errorf("for %q: expected %q, got %q", tc.input, tc.expected, comment.Data())
		}
	}
}
