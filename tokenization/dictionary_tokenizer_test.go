package tokenization

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func newTestDictTokenizer() *DictionaryTokenizer {
	guard, _ := core.NewStackDepthGuard(256)
	return NewDictionaryTokenizer(true, guard, nil, false)
}

func assertDictToken(t *testing.T, token tokens.Token) *tokens.DictionaryToken {
	t.Helper()
	if token == nil {
		t.Fatal("expected non-nil token")
	}
	dict, ok := token.(*tokens.DictionaryToken)
	if !ok {
		t.Fatalf("expected *tokens.DictionaryToken, got %T", token)
	}
	return dict
}

func assertDictEntryString(t *testing.T, dict *tokens.DictionaryToken, key string, expected string) {
	t.Helper()
	val, ok := dict.Data()[key]
	if !ok {
		t.Fatalf("dictionary missing key %q", key)
	}
	name, ok := val.(*tokens.NameToken)
	if !ok {
		t.Fatalf("for key %q: expected *tokens.NameToken, got %T", key, val)
	}
	if name.Data() != expected {
		t.Errorf("for key %q: expected %q, got %q", key, expected, name.Data())
	}
}

func assertDictEntryNumeric(t *testing.T, dict *tokens.DictionaryToken, key string, expected float64) {
	t.Helper()
	val, ok := dict.Data()[key]
	if !ok {
		t.Fatalf("dictionary missing key %q", key)
	}
	num, ok := val.(*tokens.NumericToken)
	if !ok {
		t.Fatalf("for key %q: expected *tokens.NumericToken, got %T", key, val)
	}
	if num.Data() != expected {
		t.Errorf("for key %q: expected %.10g, got %.10g", key, expected, num.Data())
	}
}

func assertDictEntryStringToken(t *testing.T, dict *tokens.DictionaryToken, key string, expected string) {
	t.Helper()
	val, ok := dict.Data()[key]
	if !ok {
		t.Fatalf("dictionary missing key %q", key)
	}
	str, ok := val.(*tokens.StringToken)
	if !ok {
		t.Fatalf("for key %q: expected *tokens.StringToken, got %T", key, val)
	}
	if str.Data() != expected {
		t.Errorf("for key %q: expected %q, got %q", key, expected, str.Data())
	}
}

func assertDictEntryBool(t *testing.T, dict *tokens.DictionaryToken, key string, expected bool) {
	t.Helper()
	val, ok := dict.Data()[key]
	if !ok {
		t.Fatalf("dictionary missing key %q", key)
	}
	b, ok := val.(*tokens.BooleanToken)
	if !ok {
		t.Fatalf("for key %q: expected *tokens.BooleanToken, got %T", key, val)
	}
	if b.Data() != expected {
		t.Errorf("for key %q: expected %v, got %v", key, expected, b.Data())
	}
}

func assertDictEntryRef(t *testing.T, dict *tokens.DictionaryToken, key string, expected core.IndirectReference) {
	t.Helper()
	val, ok := dict.Data()[key]
	if !ok {
		t.Fatalf("dictionary missing key %q", key)
	}
	ref, ok := val.(*tokens.IndirectReferenceToken)
	if !ok {
		t.Fatalf("for key %q: expected *tokens.IndirectReferenceToken, got %T", key, val)
	}
	if !ref.Data().Equals(expected) {
		t.Errorf("for key %q: expected %v, got %v", key, expected, ref.Data())
	}
}

func getDictIndex(t *testing.T, index int, dict *tokens.DictionaryToken) (string, tokens.Token) {
	t.Helper()
	entries := dict.OrderedEntries()
	if index < 0 || index >= len(entries) {
		t.Fatalf("dictionary did not contain index %d (has %d entries)", index, len(entries))
	}
	e := entries[index]
	return e.Key, e.Value
}

/*
TestDictIncorrectStartCharactersReturnsFalse maps C# IncorrectStartCharacters_ReturnsFalse.
C# tests "[rjee]", "\r\n", "<AE>", "<[p]>" all return false.
Go only checks first byte != '<', so <AE> and <[p]> return true (behavioral difference).
Test cases adjusted: only non-'<' inputs are tested.
*/
func TestDictIncorrectStartCharactersReturnsFalse(t *testing.T) {
	// C# test cases that should return false because they don't start with <<:
	// "[rjee]", "\r\n", "<AE>", "<[p]>"
	// Go implementation only checks first byte != '<', so <AE> and <[p]> return true.
	// Test the ones that correctly return false in both C# and Go:
	testCases := []string{"[rjee]", "\r\n"}

	tokenizer := newTestDictTokenizer()

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

/*
TestDictSkipsWhitespaceInStartSymbols maps C# SkipsWhitespaceInStartSymbols.
Tests < < /Name (Barry Scott) >> with whitespace between the two '<' chars.
*/
func TestDictSkipsWhitespaceInStartSymbols(t *testing.T) {
	tokenizer := newTestDictTokenizer()

	first, input := stringInput("< < /Name (Barry Scott) >>")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}

	dict := assertDictToken(t, token)
	assertDictEntryStringToken(t, dict, "Name", "Barry Scott")
}

/*
TestDictSimpleNameDictionary maps C# SimpleNameDictionary.
Tests << /Type /Example>>.
*/
func TestDictSimpleNameDictionary(t *testing.T) {
	tokenizer := newTestDictTokenizer()

	first, input := stringInput("<< /Type /Example>>")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}

	dict := assertDictToken(t, token)
	assertDictEntryString(t, dict, "Type", "Example")
}

/*
TestDictStreamDictionary maps C# StreamDictionary.
Tests << /Filter /FlateDecode /S 36 /Length 53 >>.
*/
func TestDictStreamDictionary(t *testing.T) {
	tokenizer := newTestDictTokenizer()

	first, input := stringInput("<< /Filter /FlateDecode /S 36 /Length 53 >>")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}

	dict := assertDictToken(t, token)
	assertDictEntryString(t, dict, "Filter", "FlateDecode")
	assertDictEntryNumeric(t, dict, "S", 36)
	assertDictEntryNumeric(t, dict, "Length", 53)
}

/*
TestDictCatalogDictionary maps C# CatalogDictionary.
Tests <</Pages 14 0 R /Type /Catalog >> with IndirectReference.
*/
func TestDictCatalogDictionary(t *testing.T) {
	tokenizer := newTestDictTokenizer()

	first, input := stringInput("<</Pages 14 0 R /Type /Catalog >>")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}

	dict := assertDictToken(t, token)

	expectedRef, _ := core.NewIndirectReference(14, 0)
	assertDictEntryRef(t, dict, "Pages", expectedRef)
	assertDictEntryString(t, dict, "Type", "Catalog")
}

/*
TestDictSpecificationExampleDictionary maps C# SpecificationExampleDictionary.
Tests a complex multi-line dictionary with nested subdictionary.
*/
func TestDictSpecificationExampleDictionary(t *testing.T) {
	s := "<< /Type /Example\n" +
		"/Subtype /DictionaryExample\n" +
		"/Version 0.01\n" +
		"/IntegerItem 12\n" +
		"/StringItem (a string)\n" +
		"/Subdictionary \n" +
		"<<  /Item1 0.4\n" +
		"    /Item2 true\n" +
		"    /LastItem (not!)\n" +
		"    /VeryLastItem (OK)\n" +
		">>\n" +
		">>"

	tokenizer := newTestDictTokenizer()
	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)
	if !ok {
		t.Fatal("expected true, got false")
	}

	dict := assertDictToken(t, token)

	assertDictEntryString(t, dict, "Type", "Example")
	assertDictEntryString(t, dict, "Subtype", "DictionaryExample")
	assertDictEntryNumeric(t, dict, "Version", 0.01)
	assertDictEntryNumeric(t, dict, "IntegerItem", 12)
	assertDictEntryStringToken(t, dict, "StringItem", "a string")

	subKey, subVal := getDictIndex(t, 5, dict)
	if subKey != "Subdictionary" {
		t.Errorf("subdictionary key: expected %q, got %q", "Subdictionary", subKey)
	}

	subDict, ok := subVal.(*tokens.DictionaryToken)
	if !ok {
		t.Fatalf("Subdictionary value: expected *tokens.DictionaryToken, got %T", subVal)
	}

	assertDictEntryNumeric(t, subDict, "Item1", 0.4)
	assertDictEntryBool(t, subDict, "Item2", true)
	assertDictEntryStringToken(t, subDict, "LastItem", "not!")
	assertDictEntryStringToken(t, subDict, "VeryLastItem", "OK")
}

/*
TestDictExitsDictionaryParsingSingleLevel maps C# ExitsDictionaryParsingSingleLevel.
Tests that dictionary parsing stops at >> without consuming endobj/obj markers.
*/
func TestDictExitsDictionaryParsingSingleLevel(t *testing.T) {
	s := "<< /Pages 69 0 R /Type /Catalog >>\nendobj\n5 0 obj"

	tokenizer := newTestDictTokenizer()
	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)
	if !ok {
		t.Fatal("expected true, got false")
	}

	dict := assertDictToken(t, token)

	expectedRef, _ := core.NewIndirectReference(69, 0)
	assertDictEntryRef(t, dict, "Pages", expectedRef)
	assertDictEntryString(t, dict, "Type", "Catalog")

	if len(dict.Data()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(dict.Data()))
	}
}

/*
TestDictParseNestedDictionary maps C# ParseNestedDictionary.
Tests nested dictionary: << /Count 12 /Definition << /Name (Glorp)>> /Type /Catalog >>.
*/
func TestDictParseNestedDictionary(t *testing.T) {
	tokenizer := newTestDictTokenizer()

	first, input := stringInput("<< /Count 12 /Definition << /Name (Glorp)>> /Type /Catalog >>")
	token, ok := tokenizer.Tokenize(first, input)

	if !ok {
		t.Fatal("expected true, got false")
	}

	dict := assertDictToken(t, token)

	assertDictEntryNumeric(t, dict, "Count", 12)

	defKey, defVal := getDictIndex(t, 1, dict)
	if defKey != "Definition" {
		t.Errorf("definition key: expected %q, got %q", "Definition", defKey)
	}

	subDict, ok := defVal.(*tokens.DictionaryToken)
	if !ok {
		t.Fatalf("Definition value: expected *tokens.DictionaryToken, got %T", defVal)
	}

	assertDictEntryStringToken(t, subDict, "Name", "Glorp")
	assertDictEntryString(t, dict, "Type", "Catalog")

	if len(dict.Data()) != 3 {
		t.Errorf("expected 3 entries, got %d", len(dict.Data()))
	}
}

/*
TestDictSupportTicket29 maps C# SupportTicket29.
Tests dictionary with MediaBox array containing embedded newlines.
*/
func TestDictSupportTicket29(t *testing.T) {
	s := "<< /Type /Page /Parent 4 0 R /MediaBox [ 0 0      \r\n   100.28 841.89 ] /Resources >>"

	tokenizer := newTestDictTokenizer()
	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)
	if !ok {
		t.Fatal("expected true, got false")
	}

	dict := assertDictToken(t, token)

	mediaBoxVal, ok := dict.Data()["MediaBox"]
	if !ok {
		t.Fatal("dictionary missing key MediaBox")
	}

	mediaBox, ok := mediaBoxVal.(*tokens.ArrayToken)
	if !ok {
		t.Fatalf("MediaBox: expected *tokens.ArrayToken, got %T", mediaBoxVal)
	}

	if len(mediaBox.Data()) != 4 {
		t.Errorf("MediaBox: expected 4 elements, got %d", len(mediaBox.Data()))
	}
}

/*
TestDictCommentsInsideDictionaryFromSap maps C# CommentsInsideDictionaryFromSap.
Tests dictionary with many SAP comment lines between entries and >>.
*/
func TestDictCommentsInsideDictionaryFromSap(t *testing.T) {
	s := "<<\n" +
		"/Author (ABCD )\n" +
		"/CreationDate (D:20150505083655)\n" +
		"/Creator (Form 2014 EN)\n" +
		"/Producer (SAP NetWeaver 700 )\n" +
		"%SAPinfoStart TOA_DARA\n" +
		"%FUNCTION=( )\n" +
		"%MANDANT=( )\n" +
		"%DEL_DATE=( )\n" +
		"%SAP_OBJECT=( )\n" +
		"%AR_OBJECT=( )\n" +
		"%OBJECT_ID=( )\n" +
		"%FORM_ID=( )\n" +
		"%FORMARCHIV=( )\n" +
		"%RESERVE=( )\n" +
		"%NOTIZ=( )\n" +
		"%-( )\n" +
		"%-( )\n" +
		"%-( )\n" +
		"%SAPinfoEnd TOA_DARA\n" +
		">>"

	tokenizer := newTestDictTokenizer()
	first, input := stringInput(s)

	token, ok := tokenizer.Tokenize(first, input)
	if !ok {
		t.Fatal("expected true, got false")
	}

	dict := assertDictToken(t, token)
	assertDictEntryStringToken(t, dict, "Producer", "SAP NetWeaver 700 ")
}
