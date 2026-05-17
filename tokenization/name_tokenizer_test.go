package tokenization

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

/*
TestNameReadsName maps C# ReadsName.
Tests basic name tokenization from "/Type /XRef" — should return "Type".
*/
func TestNameReadsName(t *testing.T) {
	tokenizer := NewNameTokenizer()

	first, input := stringInput("/Type /XRef")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	nameToken := assertNameToken(t, token)
	if nameToken.Data() != "Type" {
		t.Errorf("expected %q, got %q", "Type", nameToken.Data())
	}
}

/*
TestNameReadsNameNoEndSpace maps C# ReadsNameNoEndSpace.
Tests "/Type/XRef" without space between names — should return "Type".
*/
func TestNameReadsNameNoEndSpace(t *testing.T) {
	tokenizer := NewNameTokenizer()

	first, input := stringInput("/Type/XRef")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	nameToken := assertNameToken(t, token)
	if nameToken.Data() != "Type" {
		t.Errorf("expected %q, got %q", "Type", nameToken.Data())
	}
}

/*
TestNameReadsNameNotAtForwardSlash maps C# ReadsName_NotAtForwardSlash_Throws.
Tests that a space as first byte returns false since it's not '/'.
*/
func TestNameReadsNameNotAtForwardSlash(t *testing.T) {
	tokenizer := NewNameTokenizer()

	first, input := stringInput(" /Type")
	token, ok := tokenizer.Tokenize(first, input)

	if ok {
		t.Error("expected false for space byte, got true")
	}
	if token != nil {
		t.Errorf("expected nil token, got %v", token)
	}
}

/*
TestNameReadsNameAtEndOfStream maps C# ReadsNameAtEndOfStream.
Tests "/XRef" at end of stream — should return "XRef".
*/
func TestNameReadsNameAtEndOfStream(t *testing.T) {
	tokenizer := NewNameTokenizer()

	first, input := stringInput("/XRef")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	nameToken := assertNameToken(t, token)
	if nameToken.Data() != "XRef" {
		t.Errorf("expected %q, got %q", "XRef", nameToken.Data())
	}
}

/*
TestNameFallsBackToUnescapedForEarlyPdfTypes maps C# FallsBackToUnescapedForEarlyPdfTypes.
Tests "/Priorto1.2#INvalidHexHash" — the # is followed by 'I' which is not a hex digit,
so it falls back to treating everything as literal characters.
*/
func TestNameFallsBackToUnescapedForEarlyPdfTypes(t *testing.T) {
	tokenizer := NewNameTokenizer()

	first, input := stringInput("/Priorto1.2#INvalidHexHash")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	nameToken := assertNameToken(t, token)
	if nameToken.Data() != "Priorto1.2#INvalidHexHash" {
		t.Errorf("expected %q, got %q", "Priorto1.2#INvalidHexHash", nameToken.Data())
	}
}

/*
TestNameReadsValidPdfNames maps C# ReadsValidPdfNames parameterized test.
Tests various valid PDF name objects with different character sets.
*/
func TestNameReadsValidPdfNames(t *testing.T) {
	tokenizer := NewNameTokenizer()

	tests := []struct {
		input    string
		expected string
	}{
		{"/Name1", "Name1"},
		{"/ASomewhatLongerName", "ASomewhatLongerName"},
		{"/A−Name_With;Various***Characters?", "A−Name_With;Various***Characters?"},
		{"/1.2", "1.2"},
		{"/$$", "$$"},
		{"/@pattern", "@pattern"},
		{"/.notdef", ".notdef"},
	}

	for _, tc := range tests {
		first, input := stringInput(tc.input)
		token, ok := tokenizer.Tokenize(first, input)

		if !ok {
			t.Errorf("expected true for %q, got false", tc.input)
			continue
		}
		nameToken := assertNameToken(t, token)
		if nameToken.Data() != tc.expected {
			t.Errorf("for %q: expected %q, got %q", tc.input, tc.expected, nameToken.Data())
		}
	}
}

/*
TestNameReadsHexNames maps C# ReadsHexNames parameterized test.
Tests hex-encoded characters within PDF names using #XX syntax.
*/
func TestNameReadsHexNames(t *testing.T) {
	tokenizer := NewNameTokenizer()

	tests := []struct {
		input    string
		expected string
	}{
		{"/Adobe#20Green", "Adobe Green"},
		{"/PANTONE#205757#20CV", "PANTONE 5757 CV"},
		{"/paired#28#29parentheses", "paired()parentheses"},
		{"/The_Key_of_F#23_Minor", "The_Key_of_F#_Minor"},
		{"/A#42", "AB"},
	}

	for _, tc := range tests {
		first, input := stringInput(tc.input)
		token, ok := tokenizer.Tokenize(first, input)

		if !ok {
			t.Errorf("expected true for %q, got false", tc.input)
			continue
		}
		nameToken := assertNameToken(t, token)
		if nameToken.Data() != tc.expected {
			t.Errorf("for %q: expected %q, got %q", tc.input, tc.expected, nameToken.Data())
		}
	}
}

/*
TestNameIgnoredInvalidHex maps C# IgnoredInvalidHex.
Tests "/Invalid#AZBadHex" — #A followed by 'Z' (not hex), so the # and A are
treated as literal characters. Result: "Invalid#AZBadHex".
*/
func TestNameIgnoredInvalidHex(t *testing.T) {
	tokenizer := NewNameTokenizer()

	first, input := stringInput("/Invalid#AZBadHex")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	nameToken := assertNameToken(t, token)
	if nameToken.Data() != "Invalid#AZBadHex" {
		t.Errorf("expected %q, got %q", "Invalid#AZBadHex", nameToken.Data())
	}
}

/*
TestNameIgnoreInvalidSingleHex maps C# IgnoreInvalidSingleHex.
Tests "/Invalid#Z" — # followed by single 'Z' which is not a valid hex pair.
Result: "Invalid#Z".
*/
func TestNameIgnoreInvalidSingleHex(t *testing.T) {
	tokenizer := NewNameTokenizer()

	first, input := stringInput("/Invalid#Z")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	nameToken := assertNameToken(t, token)
	if nameToken.Data() != "Invalid#Z" {
		t.Errorf("expected %q, got %q", "Invalid#Z", nameToken.Data())
	}
}

/*
TestNameEndsNameFollowingInvalidHex maps C# EndsNameFollowingInvalidHex.
Tests "/Hex#/Name" — # at end of "Hex#" followed by '/' which ends the name.
Result: "Hex#".
*/
func TestNameEndsNameFollowingInvalidHex(t *testing.T) {
	tokenizer := NewNameTokenizer()

	first, input := stringInput("/Hex#/Name")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	nameToken := assertNameToken(t, token)
	if nameToken.Data() != "Hex#" {
		t.Errorf("expected %q, got %q", "Hex#", nameToken.Data())
	}
}

func assertNameToken(t *testing.T, token tokens.Token) *tokens.NameToken {
	t.Helper()
	if token == nil {
		t.Fatal("expected non-nil token")
	}
	name, ok := token.(*tokens.NameToken)
	if !ok {
		t.Fatalf("expected *tokens.NameToken, got %T", token)
	}
	return name
}
