package truetypeparser_test

import (
	"bytes"
	"testing"

	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	truetypeparser "github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
)

// TestOs2WritesSameTableAsRead verifies that serializing the parsed OS/2 table
// produces byte-for-byte identical output to the original table data in the font file.
func TestOs2WritesSameTableAsRead(t *testing.T) {
	fonts := []string{
		"Andada-Regular",
		"Roboto-Regular",
		"PMingLiU",
	}

	for _, fontFile := range fonts {
		t.Run(fontFile, func(t *testing.T) {
			fontBytes := readFileBytes(t, fontFile)

			input := truetypeparser.NewTrueTypeDataBytes(fontBytes)
			parsed := truetypeparser.Parse(input)

			if parsed == nil {
				t.Fatal("parsed font is nil")
			}

			os2 := parsed.TableRegister().Os2Table
			os2Header, ok := parsed.TableHeaders()[truetype.Os2]
			if !ok {
				t.Fatalf("OS/2 table header not found in %q", fontFile)
			}

			os2InputBytes := fontBytes[os2Header.Offset:os2Header.Offset+os2Header.Length]

			var buf bytes.Buffer
			if err := os2.Write(&buf); err != nil {
				t.Fatalf("failed to write OS/2 table: %v", err)
			}

			result := buf.Bytes()

			if len(result) != len(os2InputBytes) {
				t.Errorf("OS/2 table length mismatch: got %d, want %d", len(result), len(os2InputBytes))
			}

			for i := 0; i < len(os2InputBytes); i++ {
				if i >= len(result) {
					break
				}
				if result[i] != os2InputBytes[i] {
					t.Errorf("byte mismatch at index %d: got 0x%02x, want 0x%02x", i, result[i], os2InputBytes[i])
					break
				}
			}
		})
	}
}
