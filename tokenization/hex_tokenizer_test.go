package tokenization

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

/*
TestHexCannotTokenizeInvalidBytes maps C# CannotTokenizeInvalidBytes.
Tests that inputs not starting with '<' return false and nil token.
*/
func TestHexCannotTokenizeInvalidBytes(t *testing.T) {
	tokenizer := NewHexTokenizer()

	testCases := []string{
		">not hex",
		"\\<not hex",
		"not hex",
		"AE1094 still not hex",
	}

	for _, s := range testCases {
		first, input := stringInput(s)
		token, ok := tokenizer.Tokenize(first, input)

		if ok {
			t.Errorf("expected false for %q, got true", s)
		}
		if token != nil {
			t.Errorf("expected nil token for %q, got %v", s, token)
		}
	}
}

/*
TestHexTokenizesHexStringsCorrectly maps C# TokenizesHexStringsCorrectly.
Tests basic hex string decoding: "<00>" produces empty string (null byte skipped),
"<A1>" produces the character at byte 0xA1 cast to char (PDFDocEncoding).
*/
func TestHexTokenizesHexStringsCorrectly(t *testing.T) {
	tokenizer := NewHexTokenizer()

	tests := []struct {
		input    string
		expected string
	}{
		{"<00>", ""},
		{"<A1>", "¡"},
	}

	for _, tc := range tests {
		first, input := stringInput(tc.input)
		token, ok := tokenizer.Tokenize(first, input)

		if !ok {
			t.Errorf("expected true for %q, got false", tc.input)
			continue
		}

		hexToken := assertHexToken(t, token)
		if hexToken.Data() != tc.expected {
			t.Errorf("for %q: expected data %q, got %q", tc.input, tc.expected, hexToken.Data())
		}
	}
}

/*
TestHexHandlesUtf16Strings maps C# HandlesUtf16Strings.
Tests UTF-16BE BOM detection and decoding in hex strings.
*/
func TestHexHandlesUtf16Strings(t *testing.T) {
	tokenizer := NewHexTokenizer()

	tests := []struct {
		input    string
		expected string
	}{
		{"<FEFF004C0069006200720065004F0066006600690063006500200036002E0031>", "LibreOffice 6.1"},
		{"<FEFF30533093306B3061306F4E16754C>", "こんにちは世界"},
	}

	for _, tc := range tests {
		first, input := stringInput(tc.input)
		token, ok := tokenizer.Tokenize(first, input)

		if !ok {
			t.Errorf("expected true for %q, got false", tc.input)
			continue
		}

		hexToken := assertHexToken(t, token)
		if hexToken.Data() != tc.expected {
			t.Errorf("for %q: expected data %q, got %q", tc.input, tc.expected, hexToken.Data())
		}
	}
}

func assertHexToken(t *testing.T, token tokens.Token) *tokens.HexToken {
	t.Helper()
	if token == nil {
		t.Fatal("expected non-nil token")
	}
	hex, ok := token.(*tokens.HexToken)
	if !ok {
		t.Fatalf("expected *tokens.HexToken, got %T", token)
	}
	return hex
}
