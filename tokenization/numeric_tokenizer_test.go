package tokenization

import (
	"math"
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

/*
TestNumFirstByteInvalidReturnsFalse maps C# FirstByteInvalid_ReturnsFalse.
Tests that letters and non-numeric characters as first byte return false.
*/
func TestNumFirstByteInvalidReturnsFalse(t *testing.T) {
	tokenizer := NewNumericTokenizer()

	testCases := []string{"a", "b", "A", "|", "z", "e", "E", "\n"}

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
TestNumParsesValidNumbers maps C# ParsesValidNumbers parameterized test.
Tests comprehensive set of valid PDF numeric literals including integers,
negatives, decimals, and scientific notation.
*/
func TestNumParsesValidNumbers(t *testing.T) {
	tokenizer := NewNumericTokenizer()

	tests := []struct {
		input    string
		expected float64
		tolerance *float64
	}{
		{"0", 0, nil},
		{"0003", 3, nil},
		{"1", 1, nil},
		{"2", 2, nil},
		{"3", 3, nil},
		{"4", 4, nil},
		{"5", 5, nil},
		{"6", 6, nil},
		{"7", 7, nil},
		{"8", 8, nil},
		{"9", 9, nil},
		{"10", 10, nil},
		{"11", 11, nil},
		{"29", 29, nil},
		{"-0", 0, nil},
		{"-0123", -123, nil},
		{"-6.9000", -6.9, nil},
		{"57473.3458382", 57473.3458382, nil},
		{"123", 123, nil},
		{"43445", 43445, nil},
		{"+17", 17, nil},
		{"-98", -98, nil},
		{"34.5", 34.5, nil},
		{"-3.62", -3.62, nil},
		{"+123.6", 123.6, nil},
		{"4.", 4, nil},
		{"-.002", -0.002, nil},
		{"0.0", 0, nil},
		{"1.57e3", 1570, nil},
		{"1.57e-3", 0.00157, floatPtr(0.0000001)},
		{"1.24e1", 12.4, nil},
		{"1.457E2", 145.7, nil},
	}

	for _, tc := range tests {
		first, input := stringInput(tc.input)
		token, ok := tokenizer.Tokenize(first, input)

		if !ok {
			t.Errorf("expected true for %q, got false", tc.input)
			continue
		}
		numToken := assertNumericToken(t, token)
		got := numToken.Data()
		if tc.tolerance != nil {
			if math.Abs(got-tc.expected) > *tc.tolerance {
				t.Errorf("for %q: expected %.10g (tolerance %.0e), got %.10g", tc.input, tc.expected, *tc.tolerance, got)
			}
		} else {
			if got != tc.expected {
				t.Errorf("for %q: expected %.10g, got %.10g", tc.input, tc.expected, got)
			}
		}
	}
}

/*
TestNumOnlyParsesNumberPart maps C# OnlyParsesNumberPart.
Tests "135.6654/Type" — should parse only the numeric part and leave
the cursor at '/'.
*/
func TestNumOnlyParsesNumberPart(t *testing.T) {
	tokenizer := NewNumericTokenizer()

	first, input := stringInput("135.6654/Type")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	numToken := assertNumericToken(t, token)
	if numToken.Data() != 135.6654 {
		t.Errorf("expected %.6g, got %.6g", 135.6654, numToken.Data())
	}
	if input.CurrentByte() != '/' {
		t.Errorf("expected current byte '/', got %c", input.CurrentByte())
	}
}

/*
TestNumHandlesDash maps C# HandlesDash.
Tests standalone "-" which is parsed as 0.
*/
func TestNumHandlesDash(t *testing.T) {
	tokenizer := NewNumericTokenizer()

	first, input := stringInput("-")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	numToken := assertNumericToken(t, token)
	if numToken.Data() != 0 {
		t.Errorf("expected 0, got %.10g", numToken.Data())
	}
}

/*
TestNumHandleDoubleDashedNumber maps C# HandleDoubleDashedNumber.
Tests "--10.25" — a double-dashed number seen in the wild, parsed as -10.25.
*/
func TestNumHandleDoubleDashedNumber(t *testing.T) {
	tokenizer := NewNumericTokenizer()

	first, input := stringInput("--10.25")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	numToken := assertNumericToken(t, token)
	if numToken.Data() != -10.25 {
		t.Errorf("expected %.4g, got %.10g", -10.25, numToken.Data())
	}
}

/*
TestNumHandlesDot maps C# HandlesDot.
Tests standalone "." which is parsed as 0.
*/
func TestNumHandlesDot(t *testing.T) {
	tokenizer := NewNumericTokenizer()

	first, input := stringInput(".")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	numToken := assertNumericToken(t, token)
	if numToken.Data() != 0 {
		t.Errorf("expected 0, got %.10g", numToken.Data())
	}
}

func assertNumericToken(t *testing.T, token tokens.Token) *tokens.NumericToken {
	t.Helper()
	if token == nil {
		t.Fatal("expected non-nil token")
	}
	num, ok := token.(*tokens.NumericToken)
	if !ok {
		t.Fatalf("expected *tokens.NumericToken, got %T", token)
	}
	return num
}

func floatPtr(f float64) *float64 {
	return &f
}
