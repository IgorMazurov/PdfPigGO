package filters

import (
	"bytes"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestFlateFilter_EncodeAndDecodePreservesInput(t *testing.T) {
	filter := NewFlateFilter()
	input := []byte{67, 69, 69, 10, 4, 20, 6, 19, 120, 64, 64, 64, 32}

	inputStream := bytes.NewReader(input)
	result, err := filter.Encode(inputStream)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}

	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})
	decoded, err := filter.Decode(result, dict, Instance, 0)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	if !bytes.Equal(input, decoded) {
		t.Errorf("Encode/Decode roundtrip mismatch.\nexpected: %v\n  actual: %v", input, decoded)
	}
}

func TestFlateFilter_CanDecodeCorruptedInputIssue1235(t *testing.T) {
	hexStr := "789C958D5D0AC2400C844FB077980B74BB7FD9D982F820B43E8B7B03C542C187EAFDC1F84B7D1164200999E49BD9044C6653D10E1E443DA1AF6636ED76EF315E7572968E1ECDAB7FB7506C4C59C0AEB3912EE270366AAAF4E36D364BF7911450DC274A5112B1AC9751D77A58680B51A4D8AE433D62953C037396E0F290FBE098B267A43051725AA34E77E44EF50B1B52B42C90E4ADF83FB94FDD0000000000"

	hex := tokens.NewHexToken([]rune(hexStr))
	filter := NewFlateFilter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	result, err := filter.Decode(hex.Bytes(), dict, Instance, 0)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	text := core.BytesAsLatin1String(result)

	if !strings.HasPrefix(text, "q") {
		t.Errorf("Expected decoded text to start with 'q', got: %q", text)
	}
}
