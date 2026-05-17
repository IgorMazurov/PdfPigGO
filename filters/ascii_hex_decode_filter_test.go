package filters

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestAsciiHexDecodeFilter_DecodesEncodedTextProperly(t *testing.T) {
	filter := NewAsciiHexDecodeFilter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	input := []byte("7368652073656C6C73207365617368656C6C73206F6E20746865207365612073686F7265")

	result, err := filter.Decode(input, dict, Instance, 1)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	expected := "she sells seashells on the sea shore"
	got := string(result)
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestAsciiHexDecodeFilter_DecodesEncodedTextWithBracesProperly(t *testing.T) {
	filter := NewAsciiHexDecodeFilter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	input := []byte("<7368652073656C6C73207365617368656C6C73206F6E20746865207365612073686F7265>")

	result, err := filter.Decode(input, dict, Instance, 1)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	expected := "she sells seashells on the sea shore"
	got := string(result)
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestAsciiHexDecodeFilter_DecodesEncodedTextWithWhitespaceProperly(t *testing.T) {
	filter := NewAsciiHexDecodeFilter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	input := []byte("6F6E6365207      5706F6E206120     74696D6520696E\n    20612067616C6178792046617220466172204177    6179")

	result, err := filter.Decode(input, dict, Instance, 1)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	expected := "once upon a time in a galaxy Far Far Away"
	got := string(result)
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestAsciiHexDecodeFilter_DecodesEncodedTextLowercaseProperly(t *testing.T) {
	filter := NewAsciiHexDecodeFilter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	input := []byte("6f6e63652075706f6e20612074696d6520696e20612067616c61787920466172204661722041776179")

	result, err := filter.Decode(input, dict, Instance, 1)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	expected := "once upon a time in a galaxy Far Far Away"
	got := string(result)
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestAsciiHexDecodeFilter_DecodeWithInvalidCharactersReturnsError(t *testing.T) {
	filter := NewAsciiHexDecodeFilter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	tests := []struct {
		name  string
		input string
	}{
		{name: "ZA", input: "ZA"},
		{name: "AM", input: "AM"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := []byte(tc.input)
			_, err := filter.Decode(input, dict, Instance, 1)
			if err == nil {
				t.Fatal("expected error for invalid characters, got nil")
			}
		})
	}
}

func TestAsciiHexDecodeFilter_SubstitutesZeroForLastByte(t *testing.T) {
	filter := NewAsciiHexDecodeFilter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	input := []byte("AE5>")

	result, err := filter.Decode(input, dict, Instance, 1)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	got := string(result)
	if got != "\xAE\x50" {
		t.Errorf("expected \\xAE\\x50, got %q", got)
	}
}

func TestAsciiHexDecodeFilter_DecodesEncodedTextStoppingAtLastBrace(t *testing.T) {
	filter := NewAsciiHexDecodeFilter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	input := []byte("6f6e63652075706f6e20612074696d6520696e20612067616c61787920466172204661722041776179> There is stuff following the EOD.")

	result, err := filter.Decode(input, dict, Instance, 1)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	expected := "once upon a time in a galaxy Far Far Away"
	got := string(result)
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
