package parser

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/cmap"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestCanParseWithArray(t *testing.T) {
	inputBytes := core.NewMemoryInputBytes([]byte("<0003> <0004> [<0020> <0041>]"))
	guard, err := core.NewStackDepthGuard(256)
	if err != nil {
		t.Fatalf("unexpected error creating StackDepthGuard: %v", err)
	}
	scanner := tokenization.NewCoreTokenScanner(inputBytes, false, guard, tokenization.ScannerScopeNone, nil, false, false)

	parser := NewBaseFontRangeParser()
	builder := cmap.NewCharacterMapBuilder()

	err = parser.Parse(tokens.NewNumericTokenFromInt(1), scanner, builder)
	if err != nil {
		t.Fatalf("unexpected error from Parse: %v", err)
	}

	bfcm := builder.BaseFontCharacterMap()
	if len(bfcm) != 2 {
		t.Fatalf("expected 2 base font characters, got %d", len(bfcm))
	}

	if bfcm[3] != " " {
		t.Errorf("expected base font char at code 3 to be \" \", got %q", bfcm[3])
	}
	if bfcm[4] != "A" {
		t.Errorf("expected base font char at code 4 to be \"A\", got %q", bfcm[4])
	}
}

func TestCanParseWithHex(t *testing.T) {
	inputBytes := core.NewMemoryInputBytes([]byte("<8141> <8147> <8141>"))
	guard, err := core.NewStackDepthGuard(256)
	if err != nil {
		t.Fatalf("unexpected error creating StackDepthGuard: %v", err)
	}
	scanner := tokenization.NewCoreTokenScanner(inputBytes, false, guard, tokenization.ScannerScopeNone, nil, false, false)

	parser := NewBaseFontRangeParser()
	builder := cmap.NewCharacterMapBuilder()

	err = parser.Parse(tokens.NewNumericTokenFromInt(1), scanner, builder)
	if err != nil {
		t.Fatalf("unexpected error from Parse: %v", err)
	}

	bfcm := builder.BaseFontCharacterMap()
	if len(bfcm) != 7 {
		t.Fatalf("expected 7 base font characters, got %d", len(bfcm))
	}

	if bfcm[33089] != "\u8141" {
		t.Errorf("expected base font char at code 33089 to be U+8141, got %q", bfcm[33089])
	}
	if bfcm[33090] != "\u8142" {
		t.Errorf("expected base font char at code 33090 to be U+8142, got %q", bfcm[33090])
	}
}

func TestCanParseTwoRowsWithDifferentFormat(t *testing.T) {
	inputStr := "<0019> <001B> <3C>\n<0001> <0003> [/happy /feet /penguin]"
	inputBytes := core.NewMemoryInputBytes([]byte(inputStr))
	guard, err := core.NewStackDepthGuard(256)
	if err != nil {
		t.Fatalf("unexpected error creating StackDepthGuard: %v", err)
	}
	scanner := tokenization.NewCoreTokenScanner(inputBytes, false, guard, tokenization.ScannerScopeNone, nil, false, false)

	parser := NewBaseFontRangeParser()
	builder := cmap.NewCharacterMapBuilder()

	err = parser.Parse(tokens.NewNumericTokenFromInt(2), scanner, builder)
	if err != nil {
		t.Fatalf("unexpected error from Parse: %v", err)
	}

	bfcm := builder.BaseFontCharacterMap()
	if len(bfcm) != 6 {
		t.Fatalf("expected 6 base font characters, got %d", len(bfcm))
	}

	if bfcm[1] != "happy" {
		t.Errorf("expected base font char at code 1 to be \"happy\", got %q", bfcm[1])
	}
	if bfcm[2] != "feet" {
		t.Errorf("expected base font char at code 2 to be \"feet\", got %q", bfcm[2])
	}
	if bfcm[3] != "penguin" {
		t.Errorf("expected base font char at code 3 to be \"penguin\", got %q", bfcm[3])
	}
	if bfcm[25] != "<" {
		t.Errorf("expected base font char at code 25 to be \"<\", got %q", bfcm[25])
	}
}
