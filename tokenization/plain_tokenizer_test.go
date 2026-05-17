package tokenization

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

/*
TestPlainTextNullReturnsNullToken maps C# TextNullReturnsNullToken.
Tests that "null " is recognized as the PDF null literal and returns the
singleton Null token.
*/
func TestPlainTextNullReturnsNullToken(t *testing.T) {
	tokenizer := NewPlainTokenizer()

	first, input := stringInput("null ")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	if token != tokens.Null {
		t.Errorf("expected Null singleton, got %v", token)
	}
}

/*
TestPlainTryTokenizeWhitespaceFalse maps C# TryTokenizeWhitespaceFalse.
Tests that whitespace as first byte returns false and nil token.
*/
func TestPlainTryTokenizeWhitespaceFalse(t *testing.T) {
	tokenizer := NewPlainTokenizer()

	first, input := stringInput("    something")
	token, ok := tokenizer.Tokenize(first, input)

	if ok {
		t.Error("expected false for whitespace byte, got true")
	}
	if token != nil {
		t.Errorf("expected nil token, got %v", token)
	}
}
