package tokens

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
)

func TestDictionaryTokenNullConstructorThrows(t *testing.T) {
	_, err := NewDictionary(nil)

	if err == nil {
		t.Fatal("expected error when constructing DictionaryToken with nil data")
	}
}

func TestDictionaryTokenEmptyValid(t *testing.T) {
	dict, err := NewDictionary(make(map[*NameToken]Token))
	if err != nil {
		t.Fatalf("unexpected error creating empty dictionary: %v", err)
	}

	if len(dict.Data()) != 0 {
		t.Errorf("expected empty Data(), got length %d", len(dict.Data()))
	}
}

func TestDictionaryTokenTryGetByNameEmptyDictionary(t *testing.T) {
	dict, err := NewDictionary(make(map[*NameToken]Token))
	if err != nil {
		t.Fatalf("unexpected error creating empty dictionary: %v", err)
	}

	token, found := dict.TryGet(ActualText)

	if found {
		t.Error("expected TryGet to return false for empty dictionary")
	}

	if token != nil {
		t.Error("expected nil token when key not found")
	}
}

// Matches C# TryGetByName_NullName_Throws: passing null name throws ArgumentNullException.
func TestDictionaryTokenTryGetByNameNullNameThrows(t *testing.T) {
	dict, err := NewDictionary(make(map[*NameToken]Token))
	if err != nil {
		t.Fatalf("unexpected error creating empty dictionary: %v", err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic (ArgumentNullException) when calling TryGet with nil name")
		}
	}()

	dict.TryGet(nil)
}

func TestDictionaryTokenTryGetByNameNonEmptyNotContainingKey(t *testing.T) {
	registryName := Create("Registry")
	dictData := map[*NameToken]Token{
		registryName: NewStringToken("None"),
	}

	dict, err := NewDictionary(dictData)
	if err != nil {
		t.Fatalf("unexpected error creating dictionary: %v", err)
	}

	token, found := dict.TryGet(ActualText)

	if found {
		t.Error("expected TryGet to return false for non-existent key")
	}

	if token != nil {
		t.Error("expected nil token when key not present")
	}
}

func TestDictionaryTokenTryGetByNameContainingKey(t *testing.T) {
	fishName := Create("Fish")
	dictData := map[*NameToken]Token{
		fishName:   NewNumericToken(420),
		Registry:   NewStringToken("None"),
	}

	dict, err := NewDictionary(dictData)
	if err != nil {
		t.Fatalf("unexpected error creating dictionary: %v", err)
	}

	token, found := dict.TryGet(Registry)

	if !found {
		t.Fatal("expected TryGet to return true for existing key Registry")
	}

	strToken, ok := token.(*StringToken)
	if !ok {
		t.Fatalf("expected StringToken, got %T", token)
	}

	if strToken.Data() != "None" {
		t.Errorf("expected token data %q, got %q", "None", strToken.Data())
	}
}

// Matches C# GetWithObjectNotOfTypeOrReferenceThrows: calling Get<T> when the stored
// value is not of type T (and not an indirect reference) throws PdfDocumentFormatException.
func TestDictionaryTokenGetTypedWrongTypeThrows(t *testing.T) {
	dictData := map[*NameToken]Token{
		Count: NewStringToken("twelve"),
	}

	dict, err := NewDictionary(dictData)
	if err != nil {
		t.Fatalf("unexpected error creating dictionary: %v", err)
	}

	_, err = GetTyped[*NumericToken](dict, Count)

	if err == nil {
		t.Fatal("expected PdfDocumentFormatException when getting wrong type")
	}

	if _, ok := any(err).(*core.PdfDocumentFormatException); !ok {
		t.Errorf("expected PdfDocumentFormatException, got %T: %v", err, err)
	}
}

// Matches C# WithCorrectlyAddsKey: uses Get<T> (not TryGet) to verify the added key.
func TestDictionaryTokenWithCorrectlyAddsKey(t *testing.T) {
	dictData := map[*NameToken]Token{
		Count: NewStringToken("12"),
	}

	dict, err := NewDictionary(dictData)
	if err != nil {
		t.Fatalf("unexpected error creating dictionary: %v", err)
	}

	newDict := dict.With(ActualText, NewStringToken("The text"))

	if !newDict.ContainsKey(ActualText) {
		t.Fatal("expected new dictionary to contain ActualText key")
	}

	strToken, err := GetTyped[*StringToken](newDict, ActualText)
	if err != nil {
		t.Fatalf("unexpected error getting ActualText: %v", err)
	}

	if strToken.Data() != "The text" {
		t.Errorf("expected token data %q, got %q", "The text", strToken.Data())
	}
}
