package ccitffax

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestCcittFaxDecodeFilter_CanDecodeCCittFaxCompressedImageData(t *testing.T) {
	testDataDir := filepath.Join("testdata")

	encodedBytes, err := os.ReadFile(filepath.Join(testDataDir, "ccittfax-encoded.bin"))
	if err != nil {
		t.Fatalf("ReadFile(ccittfax-encoded.bin): %v", err)
	}

	expectedBytes, err := os.ReadFile(filepath.Join(testDataDir, "ccittfax-decoded.bin"))
	if err != nil {
		t.Fatalf("ReadFile(ccittfax-decoded.bin): %v", err)
	}

	decodeParmsDict, err := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.K:         tokens.MinusOne,
		tokens.Columns:   tokens.NewNumericTokenFromInt(1800),
		tokens.Rows:      tokens.NewNumericTokenFromInt(3113),
		tokens.BlackIs1:  tokens.True,
	})
	if err != nil {
		t.Fatalf("NewDictionary(decodeParms): %v", err)
	}

	streamDict, err := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.D: tokens.NewArrayToken([]tokens.Token{
			tokens.One,
			tokens.Zero,
		}),
		tokens.W:   tokens.NewNumericTokenFromInt(1800),
		tokens.H:   tokens.NewNumericTokenFromInt(3113),
		tokens.Bpc: tokens.One,
		tokens.F:   tokens.CcittfaxDecode,
		tokens.DecodeParms: decodeParmsDict,
	})
	if err != nil {
		t.Fatalf("NewDictionary(streamDict): %v", err)
	}

	filter := NewCcittFaxDecodeFilter()
	decodedBytes, err := filter.Decode(encodedBytes, streamDict, filters.TestFilterProviderInstance, 0)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if !bytes.Equal(expectedBytes, decodedBytes) {
		t.Errorf("Decoded bytes mismatch.\nexpected length: %d\n  actual length:   %d", len(expectedBytes), len(decodedBytes))
	}
}
