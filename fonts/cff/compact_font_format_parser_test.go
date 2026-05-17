package cff_test

import (
	"os"
	"testing"

	"github.com/uglytoad/pdfpig/go/fonts/cff"
	"github.com/uglytoad/pdfpig/go/fonts/cff/charstrings"
	cffcharset "github.com/uglytoad/pdfpig/go/fonts/cff_charset"
)

// parseCharStrings adapts charstrings.Parse to the CharStringsParseCallback signature.
func parseCharStrings(
	charStringBytes [][]byte,
	subroutinesSelector *cff.CompactFontFormatSubroutinesSelector,
	charset cffcharset.CompactFontFormatCharset,
) (cff.Type2CharStringsProvider, error) {
	return charstrings.Parse(charStringBytes, subroutinesSelector, charset)
}

// getFileBytes reads a test data file from the testdata directory.
func getFileBytes(t *testing.T, name string) []byte {
	t.Helper()
	path := "testdata/" + name
	buf, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read test file %q: %v", path, err)
	}
	return buf
}

func TestCanReadMinionPro(t *testing.T) {
	fileBytes := getFileBytes(t, "MinionPro.bin")

	collection, err := cff.Parse(cff.NewCompactFontFormatData(fileBytes), parseCharStrings)
	if err != nil {
		t.Fatalf("unexpected error parsing font: %v", err)
	}

	header := collection.Header()
	if header.MajorVersion != 1 {
		t.Errorf("expected MajorVersion == 1, got %d", header.MajorVersion)
	}

	fonts := collection.Fonts()
	if len(fonts) != 1 {
		t.Fatalf("expected exactly 1 font, got %d", len(fonts))
	}

	if _, ok := fonts["MinionPro-It"]; !ok {
		t.Error("expected font collection to contain key \"MinionPro-It\"")
	}
}

func TestCanInterpretPercentSymbol(t *testing.T) {
	fileBytes := getFileBytes(t, "MinionPro.bin")

	collection, err := cff.Parse(cff.NewCompactFontFormatData(fileBytes), parseCharStrings)
	if err != nil {
		t.Fatalf("unexpected error parsing font: %v", err)
	}

	// Calls a global subroutine
	box := collection.GetCharacterBoundingBox("percent")
	if box == nil {
		t.Error("expected non-nil bounding box for character \"percent\"")
	}
}

func TestCanInterpretNumberSignSymbol(t *testing.T) {
	fileBytes := getFileBytes(t, "MinionPro.bin")

	collection, err := cff.Parse(cff.NewCompactFontFormatData(fileBytes), parseCharStrings)
	if err != nil {
		t.Fatalf("unexpected error parsing font: %v", err)
	}

	// Calls a local subroutine
	box := collection.GetCharacterBoundingBox("numbersign")
	if box == nil {
		t.Error("expected non-nil bounding box for character \"numbersign\"")
	}
}

func TestCanInterpretPerThousandSymbol(t *testing.T) {
	fileBytes := getFileBytes(t, "MinionPro.bin")

	collection, err := cff.Parse(cff.NewCompactFontFormatData(fileBytes), parseCharStrings)
	if err != nil {
		t.Fatalf("unexpected error parsing font: %v", err)
	}

	// Calls a local subroutine which adds to the hints
	box := collection.GetCharacterBoundingBox("perthousand")
	if box == nil {
		t.Error("expected non-nil bounding box for character \"perthousand\"")
	}
}

func TestCanInterpretATildeSmallSymbol(t *testing.T) {
	fileBytes := getFileBytes(t, "MinionPro.bin")

	collection, err := cff.Parse(cff.NewCompactFontFormatData(fileBytes), parseCharStrings)
	if err != nil {
		t.Fatalf("unexpected error parsing font: %v", err)
	}

	// Calls a global subroutine which adds to the hints
	box := collection.GetCharacterBoundingBox("Atildesmall")
	if box == nil {
		t.Error("expected non-nil bounding box for character \"Atildesmall\"")
	}
}

func TestCanInterpretUniF687Symbol(t *testing.T) {
	fileBytes := getFileBytes(t, "MinionPro.bin")

	collection, err := cff.Parse(cff.NewCompactFontFormatData(fileBytes), parseCharStrings)
	if err != nil {
		t.Fatalf("unexpected error parsing font: %v", err)
	}

	// Calls hugely nested subroutines
	box := collection.GetCharacterBoundingBox("uniF687")
	if box == nil {
		t.Error("expected non-nil bounding box for character \"uniF687\"")
	}
}

func TestCanInterpretAllGlyphs(t *testing.T) {
	fileBytes := getFileBytes(t, "MinionPro.bin")

	collection, err := cff.Parse(cff.NewCompactFontFormatData(fileBytes), parseCharStrings)
	if err != nil {
		t.Fatalf("unexpected error parsing font: %v", err)
	}

	font := collection.Fonts()["MinionPro-It"]
	if font == nil {
		t.Fatal("expected font \"MinionPro-It\" to exist in collection")
	}

	charNames := font.CharacterNames()
	if charNames == nil {
		t.Fatal("expected MinionPro to expose character names (Type 2 CharStrings)")
	}

	for _, name := range charNames {
		_, gotPath := font.TryGetPath(name)
		if !gotPath {
			t.Errorf("expected path for glyph %q", name)
		}
	}
}
