package filters

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestRunLengthFilter_CanDecodeRunLengthEncodedData(t *testing.T) {
	filter := NewRunLengthFilter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	data := []byte{
		// Write the following 6 bytes literally
		5, 0, 1, 2, 69, 12, 9,
		// Repeat 52 (257 - 254) 3 times
		254, 52,
		// Write the following 3 bytes literally
		2, 60, 61, 16,
		// Repeat 12 (257 - 250) 7 times
		250, 12,
		// Write the following 2 bytes literally
		1, 10, 19,
	}

	result, err := filter.Decode(data, dict, Instance, 1)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	expected := []byte{
		0, 1, 2, 69, 12, 9,
		52, 52, 52,
		60, 61, 16,
		12, 12, 12, 12, 12, 12, 12,
		10, 19,
	}

	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("byte[%d]: expected %d, got %d", i, expected[i], result[i])
		}
	}
}

func TestRunLengthFilter_StopsAtEndOfDataByte(t *testing.T) {
	filter := NewRunLengthFilter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	data := []byte{
		// Repeat 7 (257 - 254) 3 times
		254, 7,
		// Write the following 2 bytes literally
		1, 128, 50,
		// End of Data Byte
		128,
		// Ignore these
		90, 6, 7,
	}

	result, err := filter.Decode(data, dict, Instance, 0)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	expected := []byte{
		7, 7, 7,
		128, 50,
	}

	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("byte[%d]: expected %d, got %d", i, expected[i], result[i])
		}
	}
}
