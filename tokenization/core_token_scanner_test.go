package tokenization

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func newTestScanner(input core.InputBytes, isStream bool) *CoreTokenScanner {
	guard, _ := core.NewStackDepthGuard(256)
	return NewCoreTokenScanner(input, true, guard, ScannerScopeNone, nil, false, isStream)
}

func collectTokens(scanner *CoreTokenScanner) []tokens.Token {
	var tokensList []tokens.Token
	for scanner.Advance() {
		tokensList = append(tokensList, scanner.Current())
	}
	return tokensList
}

func assertNumeric(t *testing.T, token tokens.Token, expected float64) {
	t.Helper()
	num, ok := token.(*tokens.NumericToken)
	if !ok {
		t.Fatalf("expected *tokens.NumericToken, got %T", token)
	}
	if num.Data() != expected {
		t.Errorf("expected %.6f, got %.6f", expected, num.Data())
	}
}

func assertBoolean(t *testing.T, token tokens.Token, expected bool) {
	t.Helper()
	b, ok := token.(*tokens.BooleanToken)
	if !ok {
		t.Fatalf("expected *tokens.BooleanToken, got %T", token)
	}
	if b.Data() != expected {
		t.Errorf("expected %v, got %v", expected, b.Data())
	}
}

func assertString(t *testing.T, token tokens.Token, expected string) {
	t.Helper()
	st, ok := token.(*tokens.StringToken)
	if !ok {
		t.Fatalf("expected *tokens.StringToken, got %T", token)
	}
	if st.Data() != expected {
		t.Errorf("expected %q, got %q", expected, st.Data())
	}
}

func assertName(t *testing.T, token tokens.Token, expected string) {
	t.Helper()
	nm, ok := token.(*tokens.NameToken)
	if !ok {
		t.Fatalf("expected *tokens.NameToken, got %T", token)
	}
	if nm.Data() != expected {
		t.Errorf("expected %q, got %q", expected, nm.Data())
	}
}

func assertOperator(t *testing.T, token tokens.Token, expected string) {
	t.Helper()
	op, ok := token.(*tokens.OperatorToken)
	if !ok {
		t.Fatalf("expected *tokens.OperatorToken, got %T", token)
	}
	if op.Data() != expected {
		t.Errorf("expected %q, got %q", expected, op.Data())
	}
}

func assertArray(t *testing.T, token tokens.Token) *tokens.ArrayToken {
	t.Helper()
	arr, ok := token.(*tokens.ArrayToken)
	if !ok {
		t.Fatalf("expected *tokens.ArrayToken, got %T", token)
	}
	return arr
}

func TestScansSpecificationArrayExampleContents(t *testing.T) {
	s := "549 3.14 false (Ralph) /SomeName"
	input := core.NewMemoryInputBytes([]byte(s))
	scanner := newTestScanner(input, false)
	tokensList := collectTokens(scanner)

	if len(tokensList) < 5 {
		t.Fatalf("expected at least 5 tokens, got %d", len(tokensList))
	}

	assertNumeric(t, tokensList[0], 549)
	assertNumeric(t, tokensList[1], 3.14)
	assertBoolean(t, tokensList[2], false)
	assertString(t, tokensList[3], "Ralph")
	assertName(t, tokensList[4], "SomeName")
}

func TestScansSpecificationSimpleDictionaryExampleContents(t *testing.T) {
	s := "/Type /Example\n/Subtype /DictionaryExample\n/Version 0.01\n/IntegerItem 12\n/StringItem(a string)"
	input := core.NewMemoryInputBytes([]byte(s))
	scanner := newTestScanner(input, false)
	tokensList := collectTokens(scanner)

	if len(tokensList) < 10 {
		t.Fatalf("expected at least 10 tokens, got %d", len(tokensList))
	}

	assertName(t, tokensList[0], "Type")
	assertName(t, tokensList[1], "Example")
	assertName(t, tokensList[2], "Subtype")
	assertName(t, tokensList[3], "DictionaryExample")
	assertName(t, tokensList[4], "Version")
	assertNumeric(t, tokensList[5], 0.01)
	assertName(t, tokensList[6], "IntegerItem")
	assertNumeric(t, tokensList[7], 12)
	assertName(t, tokensList[8], "StringItem")
	assertString(t, tokensList[9], "a string")
}

func TestScansIndirectObjectExampleContents(t *testing.T) {
	s := "12 0 obj\n(Brillig)\nendobj"
	input := core.NewMemoryInputBytes([]byte(s))
	scanner := newTestScanner(input, false)
	tokensList := collectTokens(scanner)

	if len(tokensList) < 5 {
		t.Fatalf("expected at least 5 tokens, got %d", len(tokensList))
	}

	assertNumeric(t, tokensList[0], 12)
	assertNumeric(t, tokensList[1], 0)
	assertOperator(t, tokensList[2], "obj")
	assertString(t, tokensList[3], "Brillig")
	assertOperator(t, tokensList[4], "endobj")
}

func TestScansArrayInSequence(t *testing.T) {
	s := "/Bounds [12 15 19 1455.3]/Font /F1 /Name (Bob)[16]"
	input := core.NewMemoryInputBytes([]byte(s))
	scanner := newTestScanner(input, false)
	tokensList := collectTokens(scanner)

	if len(tokensList) < 7 {
		t.Fatalf("expected at least 7 tokens, got %d", len(tokensList))
	}

	assertName(t, tokensList[0], "Bounds")
	assertArray(t, tokensList[1])
	assertName(t, tokensList[2], "Font")
	assertName(t, tokensList[3], "F1")
	assertName(t, tokensList[4], "Name")
	assertString(t, tokensList[5], "Bob")
	assertArray(t, tokensList[6])
}

func TestCorrectlyScansArrayWithEscapedStrings(t *testing.T) {
	s := "<0078>Tj\n/TT0 1 Tf\n0.463 0 Td\n( )Tj\n-0.002 Tc 0.007 Tw 11.04 -0 0 11.04 180 695.52 Tm\n[(R)2.6(eg)-11.3(i)2.7(s)-2(t)4.2(r)-5.9(at)-6.6(i)2.6(on S)2(e)10.5(r)-6(v)8.9(i)2.6(c)-2(e S)1.9(o)10.6(f)-17.5(t)4.3(w)13.4(ar)-6(e \\()-6(R)2.6(S)2(S)1.9(\\))]TJ\n0 Tc 0 Tw 16.12 0 Td"
	input := core.NewMemoryInputBytes([]byte(s))
	scanner := newTestScanner(input, false)
	tokensList := collectTokens(scanner)

	if len(tokensList) != 30 {
		t.Fatalf("expected 30 tokens, got %d", len(tokensList))
	}

	assertOperator(t, tokensList[29], "Td")
	assertNumeric(t, tokensList[28], 0)
	assertNumeric(t, tokensList[27], 16.12)
	assertOperator(t, tokensList[26], "Tw")

	array := assertArray(t, tokensList[21])
	arrData := array.Data()
	if len(arrData) < 2 {
		t.Fatalf("array has only %d elements", len(arrData))
	}
	assertString(t, arrData[len(arrData)-1], ")")
	assertNumeric(t, arrData[len(arrData)-2], 1.9)
}

func TestScansStringWithoutWhitespacePreceding(t *testing.T) {
	s := "T*() Tj\n-91"
	input := core.NewMemoryInputBytes([]byte(s))
	scanner := newTestScanner(input, false)
	tokensList := collectTokens(scanner)

	if len(tokensList) != 4 {
		t.Fatalf("expected 4 tokens, got %d", len(tokensList))
	}

	assertOperator(t, tokensList[0], "T*")
	assertString(t, tokensList[1], "")
	assertOperator(t, tokensList[2], "Tj")
	assertNumeric(t, tokensList[3], -91)
}

func TestScansStringWithWeirdDoubleSymbolNumerics(t *testing.T) {
	content := "\n\t\t0.00 --21.72 TD\n\t\t/F1 8.00 Tf"
	input := core.NewMemoryInputBytes([]byte(content))
	scanner := newTestScanner(input, false)
	tokensList := collectTokens(scanner)

	if len(tokensList) != 6 {
		t.Fatalf("expected 6 tokens, got %d", len(tokensList))
	}

	assertNumeric(t, tokensList[0], 0)
	assertNumeric(t, tokensList[1], -21.72)
	assertOperator(t, tokensList[2], "TD")
	assertName(t, tokensList[3], "F1")
	assertNumeric(t, tokensList[4], 8)
	assertOperator(t, tokensList[5], "Tf")
}

func TestSkipsCommentsInStreams(t *testing.T) {
	content := "% 641 0 obj\n<<\n/Type /Encoding\n/Differences [16/quotedblleft/quotedblright 21/endash 27/ff/fi/fl/ffi 39/quoteright/parenleft/parenright 43/plus/comma/hyphen/period/slash/zero/one/two/three/four/five/six/seven/eight/nine/colon 64/at/A/B/C/D/E/F/G/H/I/J/K/L/M/N/O/P/Q/R/S/T/U/V/W/X/Y/Z/bracketleft 93/bracketright 97/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/braceleft 125/braceright 225/aacute 232/egrave/eacute 252/udieresis]\n>>\n% 315 0 obj\n<<\n/Type /Font\n/Subtype /Type1\n/BaseFont /IXNPPI+CMEX10\n/FontDescriptor 661 0 R\n/FirstChar 80\n/LastChar 88\n/Widths 644 0 R\n/ToUnicode 699 0 R\n>>\n% 306 0 obj\n<<\n/Type /Font\n/Subtype /Type1\n/BaseFont /MSNKTF+CMMI10\n/FontDescriptor 663 0 R\n/FirstChar 58\n/LastChar 119\n/Widths 651 0 R\n/ToUnicode 700 0 R\n>>"

	input := core.NewMemoryInputBytes([]byte(content))
	scanner := newTestScanner(input, true)
	tokensList := collectTokens(scanner)

	if len(tokensList) != 3 {
		t.Fatalf("stream scanner: expected 3 tokens, got %d", len(tokensList))
	}

	for _, tok := range tokensList {
		if _, ok := tok.(*tokens.DictionaryToken); !ok {
			t.Errorf("expected *tokens.DictionaryToken, got %T", tok)
		}
	}

	input2 := core.NewMemoryInputBytes([]byte(content))
	scanner2 := newTestScanner(input2, false)
	tokensList2 := collectTokens(scanner2)

	if len(tokensList2) != 6 {
		t.Fatalf("non-stream scanner: expected 6 tokens, got %d", len(tokensList2))
	}

	commentCount := 0
	dictCount := 0
	for _, tok := range tokensList2 {
		switch tok.(type) {
		case *tokens.CommentToken:
			commentCount++
		case *tokens.DictionaryToken:
			dictCount++
		}
	}

	if commentCount != 3 {
		t.Errorf("expected 3 CommentTokens, got %d", commentCount)
	}
	if dictCount != 3 {
		t.Errorf("expected 3 DictionaryTokens, got %d", dictCount)
	}
}

func TestDocument006324Test(t *testing.T) {
	content := "q\n\t1 0 0 1 248.6304 572.546 cm\n\t0 0 m\n\t\t0.021 -0.007 l\n\t\t3 -0.003 -0.01 0 0 0 c\n\tf\nQ\nq\n\t1 0 0 1 2489394 57249855 cm\n\t0 0 m\n\t\t-0.046 -0.001 -0.609 0.029 -0.286 -0.014 c\n\t\t-02.61 -0.067 -0.286 -0 .61 -0 0 c\n\tf\nQ\nq\n\t1 0 0 1 24862464 572. .836 cm\n\t0 0 m\n\t\t0.936 -0.029 l\n\t\t0.038 -0.021 0.55 -0.014 0 0 c\n\tf\nQ"

	input := core.NewMemoryInputBytes([]byte(content))
	scanner := newTestScanner(input, true)
	collectTokens(scanner)
}
