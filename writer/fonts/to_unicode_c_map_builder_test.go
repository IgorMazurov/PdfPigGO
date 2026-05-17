package fonts_test

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/parser"
	writerfonts "github.com/uglytoad/pdfpig/go/writer/fonts"
)

func TestWritesValidCMap(t *testing.T) {
	mappings := map[rune]byte{
		'1':  1,
		'=':  2,
		'H':  7,
		'a':  12,
		'2':  25,
	}

	builder := &writerfonts.ToUnicodeCMapBuilder{}
	cmapStream, err := builder.ConvertToCMapStream(mappings)
	if err != nil {
		t.Fatalf("ConvertToCMapStream: %v", err)
	}

	str := core.BytesAsLatin1String(cmapStream)
	if str == "" {
		t.Fatal("expected non-empty CMap string")
	}

	result, err := parser.NewCMapParser().Parse(core.NewMemoryInputBytes(cmapStream))
	if err != nil {
		t.Fatalf("CMapParser.Parse: %v", err)
	}

	codespaceRanges := result.CodespaceRanges()
	if len(codespaceRanges) != 1 {
		t.Fatalf("expected 1 codespace range, got %d", len(codespaceRanges))
	}

	range0 := codespaceRanges[0]

	if range0.CodeLength != 1 {
		t.Errorf("expected CodeLength == 1, got %d", range0.CodeLength)
	}

	if range0.StartInt != 0 {
		t.Errorf("expected StartInt == 0, got %d", range0.StartInt)
	}

	if range0.EndInt != 255 {
		t.Errorf("expected EndInt == 255 (byte.MaxValue), got %d", range0.EndInt)
	}

	bfcMap := result.BaseFontCharacterMap()
	if len(bfcMap) != len(mappings) {
		t.Fatalf("expected BaseFontCharacterMap count == %d, got %d", len(mappings), len(bfcMap))
	}

	for codeInt, unicodeStr := range bfcMap {
		charCode := byte(codeInt)

		var foundRune rune
		for r, c := range mappings {
			if c == charCode {
				foundRune = r
				break
			}
		}

		if len(unicodeStr) == 0 || rune(unicodeStr[0]) != foundRune {
			t.Errorf("mismatch: code %d -> unicode %q, expected rune U+%04X from mappings", codeInt, unicodeStr, foundRune)
		}
	}
}
