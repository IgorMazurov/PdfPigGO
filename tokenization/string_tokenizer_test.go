package tokenization

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func assertStringToken(t *testing.T, token tokens.Token) *tokens.StringToken {
	t.Helper()
	if token == nil {
		t.Fatal("expected non-nil token")
	}
	st, ok := token.(*tokens.StringToken)
	if !ok {
		t.Fatalf("expected *tokens.StringToken, got %T", token)
	}
	return st
}

/*
TestStrNullInputReturnsFalse maps C# NullInput_ReturnsFalse.
Tests that passing a non-parenthesis byte with nil input returns false.
*/
func TestStrNullInputReturnsFalse(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	_, ok := tokenizer.Tokenize('A', nil)

	if ok {
		t.Error("expected false for nil input, got true")
	}
}

/*
TestStrDoesNotStartWithOpenBracketReturnsFalse maps C# DoesNotStartWithOpenBracket_ReturnsFalse.
Tests that bytes other than '(' cause the tokenizer to return false immediately.
*/
func TestStrDoesNotStartWithOpenBracketReturnsFalse(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	firstBytes := []byte{'<', '\\', 'A', '[', '{', '|', '>', ' ', 'y', '^'}

	for _, b := range firstBytes {
		input := core.NewMemoryInputBytes([]byte{b})
		input.MoveNext()

		token, ok := tokenizer.Tokenize(b, input)

		if ok {
			t.Errorf("expected false for byte %q, got true", string(b))
		}
		if token != nil {
			t.Errorf("expected nil token for byte %q, got %v", string(b), token)
		}
	}
}

/*
TestStrCanHandleEscapedParentheses maps C# CanHandleEscapedParentheses.
Tests that escaped parentheses inside a string are handled correctly.
*/
func TestStrCanHandleEscapedParentheses(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := "(this string \\)contains escaped \\( parentheses)"
	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := `this string )contains escaped ( parentheses`
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrCanReadValidStrings maps C# CanReadValidStrings parameterized test.
Tests various valid PDF literal strings including newlines, balanced parens, and special chars.
*/
func TestStrCanReadValidStrings(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	tests := []struct {
		input    string
		expected string
	}{
		{"(This is a string)", "This is a string"},
		{"(Strings may contain newlines\r\nand such.)", "Strings may contain newlines\r\nand such."},
		{"(Strings may contain balanced parentheses () and special characters (*!*&}^% and so on).)",
			"Strings may contain balanced parentheses () and special characters (*!*&}^% and so on)."},
		{"()", ""},
	}

	for _, tc := range tests {
		first, input := stringInput(tc.input)

		token, ok := tokenizer.Tokenize(first, input)

		if !ok {
			t.Errorf("expected true for %q, got false", tc.input)
			continue
		}
		got := assertStringToken(t, token).Data()
		if got != tc.expected {
			t.Errorf("for input %q: expected %q, got %q", tc.input, tc.expected, got)
		}
	}
}

/*
TestStrCanHandleNestedParentheses maps C# CanHandleNestedParentheses.
Tests nested parentheses at two levels deep.
*/
func TestStrCanHandleNestedParentheses(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := "(this string (contains nested (two levels)) parentheses)"
	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := "this string (contains nested (two levels)) parentheses"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrCanHandleAngleBrackets maps C# CanHandleAngleBrackets.
Tests that angle brackets inside a literal string are preserved.
*/
func TestStrCanHandleAngleBrackets(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := "(this string <contains>)"
	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := "this string <contains>"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrSkipsEscapedEndLines maps C# SkipsEscapedEndLines.
Tests that a backslash followed by newline is treated as line continuation (skipped).
*/
func TestStrSkipsEscapedEndLines(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := `(These \
two strings \
are the same.)`

	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := "These two strings are the same."
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrTreatsEndLinesAsNewline maps C# TreatsEndLinesAsNewline.
Tests that a bare newline inside a string is preserved as \n.
*/
func TestStrTreatsEndLinesAsNewline(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := "(So does this one.\n)"

	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := "So does this one.\n"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrConvertsFullOctal maps C# ConvertsFullOctal.
Tests 3-digit octal escape sequences: \245 -> ¥, \307 -> Ç.
*/
func TestStrConvertsFullOctal(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := `(This string contains \245two octal characters\307.)`

	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := "This string contains ¥two octal charactersÇ."
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrConvertsFullOctalFollowedByNormalNumber maps C# ConvertsFullOctalFollowedByNormalNumber.
Tests \245 followed by digit '1' — only 3 octal digits consumed, '1' is literal.
*/
func TestStrConvertsFullOctalFollowedByNormalNumber(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := `(This string contains \2451 octal character.)`

	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := "This string contains ¥1 octal character."
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrConvertsPartialOctal maps C# ConvertsPartialOctal.
Tests \53 (two octal digits) -> '+' character.
*/
func TestStrConvertsPartialOctal(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := `(This string has a plus: \53 as octal)`

	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := "This string has a plus: + as octal"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrConvertsTwoPartialOctalsInARow maps C# ConvertsTwoPartialOctalsInARow.
Tests consecutive 2-digit octal escapes: \53 -> '+', \326 -> 'Ö'.
*/
func TestStrConvertsTwoPartialOctalsInARow(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := `(This string has two \53\326ctals)`

	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := "This string has two +Öctals"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrHandlesEscapedBackslash maps C# HandlesEscapedBackslash.
Tests that \\ produces a single backslash in output.
*/
func TestStrHandlesEscapedBackslash(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := `(listen\\learn)`

	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := `listen\learn`
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrWritesEscapedCharactersToOutput maps C# WritesEscapedCharactersToOutput.
Tests standard escape sequences: \n, \r, \t, \b, \f are converted to actual control chars.
*/
func TestStrWritesEscapedCharactersToOutput(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	tests := []struct {
		input    string
		expected string
	}{
		{`(new line \n)`, "new line \n"},
		{`(carriage return \r)`, "carriage return \r"},
		{`(tab \t)`, "tab \t"},
		{`(bell \b)`, "bell \b"},
		{`(uhmmm \f)`, "uhmmm \f"},
	}

	for _, tc := range tests {
		first, input := stringInput(tc.input)

		token, ok := tokenizer.Tokenize(first, input)

		if !ok {
			t.Errorf("expected true for %q, got false", tc.input)
			continue
		}
		got := assertStringToken(t, token).Data()
		if got != tc.expected {
			t.Errorf("for input %q: expected %q, got %q", tc.input, tc.expected, got)
		}
	}
}

/*
TestStrEscapedNonEscapeCharacterWritesPlainCharacter maps C# EscapedNonEscapeCharacterWritesPlainCharacter.
Tests that \e (not a recognized escape) produces just 'e' in output.
*/
func TestStrEscapedNonEscapeCharacterWritesPlainCharacter(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := `(this does not need escaping \e)`

	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := `this does not need escaping e`
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrReachesEndOfInputAssumesEndOfString maps C# ReachesEndOfInputAssumesEndOfString.
Tests that a string without closing ')' consumes to end of input and still succeeds.
*/
func TestStrReachesEndOfInputAssumesEndOfString(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := `(this does not end with bracket`

	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := `this does not end with bracket`
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrHandlesEscapedEscapeCharacters maps C# HandlesEscapedEscapeCharacters.
Tests complex escaping with escaped parens and nested escape sequences.
*/
func TestStrHandlesEscapedEscapeCharacters(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	s := `(   \(sleep 1; printf ""QUIT\\r\\n""\) | )`

	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := `   (sleep 1; printf ""QUIT\r\n"") | `
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrHandlesUtf16Strings maps C# HandlesUtf16Strings.
Tests UTF-16BE string with BOM 0xFE 0xFF decoding "Mic".
*/
func TestStrHandlesUtf16Strings(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	input := core.NewMemoryInputBytes([]byte{
		0xFE, 0xFF, 0x00, 0x4D, 0x00, 0x69, 0x00,
		0x63, 0x29,
	})

	token, ok := tokenizer.Tokenize(0x28, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := `Mic`
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

/*
TestStrHandlesUtf16BigEndianStrings maps C# HandlesUtf16BigEndianStrings.
Tests UTF-16LE string with BOM 0xFF 0xFE decoding "Mic".
*/
func TestStrHandlesUtf16BigEndianStrings(t *testing.T) {
	tokenizer := NewStringTokenizer(true)

	input := core.NewMemoryInputBytes([]byte{
		0xFF, 0xFE, 0x4D, 0x00, 0x69, 0x00, 0x63,
		0x00, 0x29,
	})

	token, ok := tokenizer.Tokenize(0x28, input)

	if !ok {
		t.Fatal("expected true, got false")
	}
	got := assertStringToken(t, token).Data()
	expected := `Mic`
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
