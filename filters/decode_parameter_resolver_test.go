package filters

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestGetFilterParameters_NilDictionaryReturnsError(t *testing.T) {
	_, err := GetFilterParameters(nil, 0)
	if err == nil {
		t.Fatal("expected error for nil stream dictionary, got nil")
	}
}

func TestGetFilterParameters_NegativeIndexReturnsError(t *testing.T) {
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	_, err := GetFilterParameters(dict, -1)
	if err == nil {
		t.Fatal("expected error for negative index, got nil")
	}
}

func TestGetFilterParameters_EmptyDictionaryReturnsEmpty(t *testing.T) {
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	result, err := GetFilterParameters(dict, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data()) != 0 {
		t.Errorf("expected empty dictionary, got %d entries", len(result.Data()))
	}
}

func TestGetFilterParameters_SingleFilterReturnsParameterDictionary(t *testing.T) {
	filterParams, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.K:         tokens.MinusOne,
		tokens.Columns:   tokens.NewNumericTokenFromInt(1800),
		tokens.Rows:      tokens.NewNumericTokenFromInt(3113),
		tokens.BlackIs1:  tokens.True,
	})

	streamDict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.Filter:     tokens.CcittfaxDecode,
		tokens.DecodeParms: filterParams,
	})

	result, err := GetFilterParameters(streamDict, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Equals(filterParams) {
		t.Error("expected result to equal filter parameters dictionary")
	}
}

func TestGetFilterParameters_SingleFilterInArrayReturnsParameterDictionary(t *testing.T) {
	filterParams, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.K:         tokens.MinusOne,
		tokens.Columns:   tokens.NewNumericTokenFromInt(1800),
		tokens.Rows:      tokens.NewNumericTokenFromInt(3113),
		tokens.BlackIs1:  tokens.True,
	})

	streamDict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.F:          tokens.NewArrayToken([]tokens.Token{tokens.CcittfaxDecode}),
		tokens.DecodeParms: tokens.NewArrayToken([]tokens.Token{filterParams}),
	})

	result, err := GetFilterParameters(streamDict, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Equals(filterParams) {
		t.Error("expected result to equal filter parameters dictionary")
	}
}

func TestGetFilterParameters_MultipleFiltersNullParameterReturnsEmpty(t *testing.T) {
	filter2Params, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.K:         tokens.MinusOne,
		tokens.Columns:   tokens.NewNumericTokenFromInt(1800),
		tokens.Rows:      tokens.NewNumericTokenFromInt(3113),
		tokens.BlackIs1:  tokens.True,
	})

	streamDict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.F:          tokens.NewArrayToken([]tokens.Token{tokens.FlateDecode, tokens.CcittfaxDecode}),
		tokens.DecodeParms: tokens.NewArrayToken([]tokens.Token{tokens.Null, filter2Params}),
	})

	result, err := GetFilterParameters(streamDict, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data()) != 0 {
		t.Errorf("expected empty dictionary for null parameter, got %d entries", len(result.Data()))
	}
}

func TestGetFilterParameters_MultipleFiltersReturnsCorrectParameterDictionary(t *testing.T) {
	filter2Params, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.K:         tokens.MinusOne,
		tokens.Columns:   tokens.NewNumericTokenFromInt(1800),
		tokens.Rows:      tokens.NewNumericTokenFromInt(3113),
		tokens.BlackIs1:  tokens.True,
	})

	streamDict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.F:          tokens.NewArrayToken([]tokens.Token{tokens.FlateDecode, tokens.CcittfaxDecode}),
		tokens.DecodeParms: tokens.NewArrayToken([]tokens.Token{tokens.Null, filter2Params}),
	})

	result, err := GetFilterParameters(streamDict, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Equals(filter2Params) {
		t.Error("expected result to equal second filter parameters dictionary")
	}
}
