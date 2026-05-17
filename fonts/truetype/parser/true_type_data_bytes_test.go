package truetypeparser_test

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
)

func TestReadUnsignedInt(t *testing.T) {
	data := truetypeparser.NewTrueTypeDataBytes([]byte{
		220,
		43,
		250,
		6,
	})

	result := data.ReadUnsignedInt()

	expected := uint32(3693869574)
	if result != expected {
		t.Errorf("ReadUnsignedInt() = %d, want %d", result, expected)
	}
}
