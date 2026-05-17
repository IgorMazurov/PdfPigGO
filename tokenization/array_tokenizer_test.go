package tokenization

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func stringInput(s string) (byte, core.InputBytes) {
	input := core.NewMemoryInputBytes([]byte(s))
	input.MoveNext()
	return input.CurrentByte(), input
}

func assertArrayToken(t *testing.T, token tokens.Token) *tokens.ArrayToken {
	t.Helper()
	if token == nil {
		t.Fatal("expected non-nil token")
	}
	arr, ok := token.(*tokens.ArrayToken)
	if !ok {
		t.Fatalf("expected *tokens.ArrayToken, got %T", token)
	}
	return arr
}

func assertDataNumeric(t *testing.T, index int, array *tokens.ArrayToken) *tokens.NumericToken {
	t.Helper()
	if len(array.Data()) <= index {
		t.Fatalf("array has only %d elements, need index %d", len(array.Data()), index)
	}
	token := array.Data()[index]
	num, ok := token.(*tokens.NumericToken)
	if !ok {
		t.Fatalf("expected *tokens.NumericToken at index %d, got %T", index, token)
	}
	return num
}

func assertDataName(t *testing.T, index int, array *tokens.ArrayToken) *tokens.NameToken {
	t.Helper()
	if len(array.Data()) <= index {
		t.Fatalf("array has only %d elements, need index %d", len(array.Data()), index)
	}
	token := array.Data()[index]
	name, ok := token.(*tokens.NameToken)
	if !ok {
		t.Fatalf("expected *tokens.NameToken at index %d, got %T", index, token)
	}
	return name
}

func assertDataString(t *testing.T, index int, array *tokens.ArrayToken) *tokens.StringToken {
	t.Helper()
	if len(array.Data()) <= index {
		t.Fatalf("array has only %d elements, need index %d", len(array.Data()), index)
	}
	token := array.Data()[index]
	str, ok := token.(*tokens.StringToken)
	if !ok {
		t.Fatalf("expected *tokens.StringToken at index %d, got %T", index, token)
	}
	return str
}

func assertDataBoolean(t *testing.T, index int, array *tokens.ArrayToken) *tokens.BooleanToken {
	t.Helper()
	if len(array.Data()) <= index {
		t.Fatalf("array has only %d elements, need index %d", len(array.Data()), index)
	}
	token := array.Data()[index]
	b, ok := token.(*tokens.BooleanToken)
	if !ok {
		t.Fatalf("expected *tokens.BooleanToken at index %d, got %T", index, token)
	}
	return b
}

func newTestStackGuard() *core.StackDepthGuard {
	guard, _ := core.NewStackDepthGuard(256)
	return guard
}

func TestInvalidFirstCharacterReturnsFalse(t *testing.T) {
	testCases := []string{"]", "<", " [", "a", "\x00"}

	tokenizer := NewArrayTokenizer(true, newTestStackGuard(), false)

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

func TestSingleElementArray(t *testing.T) {
	tokenizer := NewArrayTokenizer(true, newTestStackGuard(), false)

	numerics := []struct {
		input    string
		expected float64
	}{
		{"[12]", 12},
		{"[ 12 ]", 12},
		{"[\n2948344 ]", 2948344},
	}

	for _, tc := range numerics {
		first, input := stringInput(tc.input)
		token, ok := tokenizer.Tokenize(first, input)

		if !ok {
			t.Fatalf("expected true for %q, got false", tc.input)
		}

		arr := assertArrayToken(t, token)

		if len(arr.Data()) != 1 {
			t.Fatalf("for %q: expected 1 element, got %d", tc.input, len(arr.Data()))
		}

		num := assertDataNumeric(t, 0, arr)
		if num.Data() != tc.expected {
			t.Errorf("for %q: expected %.3f, got %.3f", tc.input, tc.expected, num.Data())
		}
	}

	strTests := []struct {
		input    string
		expected string
	}{
		{ "[(Bertrand) \t]", "Bertrand" },
		{ "[ <AE>\r\n]", "®" },
	}

	for _, tc := range strTests {
		first, input := stringInput(tc.input)
		token, ok := tokenizer.Tokenize(first, input)

		if !ok {
			t.Fatalf("expected true for %q, got false", tc.input)
		}

		arr := assertArrayToken(t, token)

		if len(arr.Data()) != 1 {
			t.Fatalf("for %q: expected 1 element, got %d", tc.input, len(arr.Data()))
		}

		dataToken := arr.Data()[0]
		var dataStr string
		switch dt := dataToken.(type) {
		case *tokens.StringToken:
			dataStr = dt.Data()
		case *tokens.HexToken:
			dataStr = dt.Data()
		default:
			t.Fatalf("for %q: expected StringToken or HexToken, got %T", tc.input, dataToken)
		}

		if dataStr != tc.expected {
			t.Errorf("for %q: expected %q, got %q", tc.input, tc.expected, dataStr)
		}
	}
}

func TestEmptyArray(t *testing.T) {
	testCases := []string{"[]", "[ ]", "[\r\n\r\n\t]"}

	tokenizer := NewArrayTokenizer(true, newTestStackGuard(), false)

	for _, s := range testCases {
		first, input := stringInput(s)
		token, ok := tokenizer.Tokenize(first, input)

		if !ok {
			t.Fatalf("expected true for %q, got false", s)
		}

		arr := assertArrayToken(t, token)

		if len(arr.Data()) != 0 {
			t.Errorf("for %q: expected empty array, got %d elements", s, len(arr.Data()))
		}
	}
}

func TestNestedArray(t *testing.T) {
	s := "[ 12 +10.453 /Fonts [ /F1 /F3 ] (Moreover) ]"

	tokenizer := NewArrayTokenizer(true, newTestStackGuard(), false)
	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)
	if !ok {
		t.Fatal("expected true, got false")
	}

	arr := assertArrayToken(t, token)

	if num := assertDataNumeric(t, 0, arr); num.Data() != 12 {
		t.Errorf("index 0: expected 12, got %.3f", num.Data())
	}

	if num := assertDataNumeric(t, 1, arr); num.Data() != 10.453 {
		t.Errorf("index 1: expected 10.453, got %.6f", num.Data())
	}

	if name := assertDataName(t, 2, arr); name.Data() != "Fonts" {
		t.Errorf("index 2: expected Fonts, got %s", name.Data())
	}

	inner := assertArrayToken(t, arr.Data()[3])

	if name := assertDataName(t, 0, inner); name.Data() != "F1" {
		t.Errorf("inner index 0: expected F1, got %s", name.Data())
	}

	if name := assertDataName(t, 1, inner); name.Data() != "F3" {
		t.Errorf("inner index 1: expected F3, got %s", name.Data())
	}

	if str := assertDataString(t, 4, arr); str.Data() != "Moreover" {
		t.Errorf("index 4: expected Moreover, got %s", str.Data())
	}
}

func TestManyNestedArrays(t *testing.T) {
	s := "[ /Bounds [ [19 -69.] [7 64.625]] (More) [[[15]]]]"

	tokenizer := NewArrayTokenizer(true, newTestStackGuard(), false)
	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)
	if !ok {
		t.Fatal("expected true, got false")
	}

	arr := assertArrayToken(t, token)

	if name := assertDataName(t, 0, arr); name.Data() != "Bounds" {
		t.Errorf("index 0: expected Bounds, got %s", name.Data())
	}

	firstInner := assertArrayToken(t, arr.Data()[1])
	firstFirstInner := assertArrayToken(t, firstInner.Data()[0])

	if num := assertDataNumeric(t, 0, firstFirstInner); num.Data() != 19 {
		t.Errorf("deep index 0: expected 19, got %.3f", num.Data())
	}

	if num := assertDataNumeric(t, 1, firstFirstInner); num.Data() != -69 {
		t.Errorf("deep index 1: expected -69, got %.3f", num.Data())
	}

	secondFirstInner := assertArrayToken(t, firstInner.Data()[1])

	if num := assertDataNumeric(t, 0, secondFirstInner); num.Data() != 7 {
		t.Errorf("deep2 index 0: expected 7, got %.3f", num.Data())
	}

	if num := assertDataNumeric(t, 1, secondFirstInner); num.Data() != 64.625 {
		t.Errorf("deep2 index 1: expected 64.625, got %.6f", num.Data())
	}

	if str := assertDataString(t, 2, arr); str.Data() != "More" {
		t.Errorf("index 2: expected More, got %s", str.Data())
	}

	secondInner := assertArrayToken(t, arr.Data()[3])
	firstSecondInner := assertArrayToken(t, secondInner.Data()[0])
	firstFirstSecondInner := assertArrayToken(t, firstSecondInner.Data()[0])

	if num := assertDataNumeric(t, 0, firstFirstSecondInner); num.Data() != 15 {
		t.Errorf("deepest index 0: expected 15, got %.3f", num.Data())
	}
}

func TestSpecificationExampleArray(t *testing.T) {
	s := "[549 3.14 false (Ralph) /SomeName]"

	tokenizer := NewArrayTokenizer(true, newTestStackGuard(), false)
	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)
	if !ok {
		t.Fatal("expected true, got false")
	}

	arr := assertArrayToken(t, token)

	if num := assertDataNumeric(t, 0, arr); num.Data() != 549 {
		t.Errorf("index 0: expected 549, got %.3f", num.Data())
	}

	if num := assertDataNumeric(t, 1, arr); num.Data() != 3.14 {
		t.Errorf("index 1: expected 3.14, got %.6f", num.Data())
	}

	if b := assertDataBoolean(t, 2, arr); b.Data() != false {
		t.Errorf("index 2: expected false, got %v", b.Data())
	}

	if str := assertDataString(t, 3, arr); str.Data() != "Ralph" {
		t.Errorf("index 3: expected Ralph, got %s", str.Data())
	}

	if name := assertDataName(t, 4, arr); name.Data() != "SomeName" {
		t.Errorf("index 4: expected SomeName, got %s", name.Data())
	}
}
