package writer_test

import (
	"bytes"
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/writer"
)

func TestEscapeSpecialCharacter(t *testing.T) {
	w := writer.NewTokenWriter()
	var buf bytes.Buffer

	if err := w.WriteToken(tokens.NewStringToken("\\"), &buf); err != nil {
		t.Fatalf("WriteToken(backslash): %v", err)
	}
	if err := w.WriteToken(tokens.NewStringToken("(Hello)"), &buf); err != nil {
		t.Fatalf("WriteToken(parentheses): %v", err)
	}

	expected := "(\\\\) (\\(Hello\\)) "
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}
