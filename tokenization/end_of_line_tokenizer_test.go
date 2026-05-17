package tokenization

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

/*
TestEOLCurrentByteIsNotEndOfLineFalse maps C# CurrentByteIsNotEndOfLineFalse.
Passes a null byte as currentByte with input "\r something \n". The tokenizer
should reject the null byte since it is neither '\r' nor '\n'.
*/
func TestEOLCurrentByteIsNotEndOfLineFalse(t *testing.T) {
	tokenizer := NewEndOfLineTokenizer()

	input := core.NewMemoryInputBytes([]byte("\r something \n"))
	input.MoveNext()

	token, ok := tokenizer.Tokenize('\x00', input)

	if ok {
		t.Error("expected false for null byte, got true")
	}
	if token != nil {
		t.Errorf("expected nil token, got %v", token)
	}
}

/*
TestEOLCurrentByteIsCarriageReturnTrue maps C# CurrentByteIsCarriageReturnTrue.
Passes '\r' as currentByte with input "\r". The tokenizer should return the
singleton EOLToken and true.
*/
func TestEOLCurrentByteIsCarriageReturnTrue(t *testing.T) {
	tokenizer := NewEndOfLineTokenizer()

	input := core.NewMemoryInputBytes([]byte("\r"))
	input.MoveNext()

	token, ok := tokenizer.Tokenize('\r', input)

	if !ok {
		t.Error("expected true for carriage return, got false")
	}
	if token != tokens.EOLToken {
		t.Errorf("expected EOLToken singleton, got %v", token)
	}
}

/*
TestEOLCurrentByteIsNewlineTrue maps C# CurrentByteIsEndOfLineTrue.
Passes '\n' as currentByte with input "\n". The tokenizer should return the
singleton EOLToken and true.
*/
func TestEOLCurrentByteIsNewlineTrue(t *testing.T) {
	tokenizer := NewEndOfLineTokenizer()

	input := core.NewMemoryInputBytes([]byte("\n"))
	input.MoveNext()

	token, ok := tokenizer.Tokenize('\n', input)

	if !ok {
		t.Error("expected true for newline, got false")
	}
	if token != tokens.EOLToken {
		t.Errorf("expected EOLToken singleton, got %v", token)
	}
}
