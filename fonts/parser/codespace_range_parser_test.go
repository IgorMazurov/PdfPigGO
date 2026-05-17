package parser

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/cmap"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestCodespaceRangeParserError(t *testing.T) {
	p := NewCodespaceRangeParser()
	byteArrayInput := core.NewMemoryInputBytes(core.StringAsLatin1Bytes("1 begincodespacerange\nendcodespacerange"))
	sdGuard, err := core.NewStackDepthGuard(256)
	if err != nil {
		t.Fatalf("unexpected error creating StackDepthGuard: %v", err)
	}
	scanner := tokenization.NewCoreTokenScanner(byteArrayInput, false, sdGuard, tokenization.ScannerScopeNone, nil, false, false)

	if !scanner.Advance() {
		t.Fatal("expected first Advance to succeed")
	}

	num, ok := scanner.Current().(*tokens.NumericToken)
	if !ok {
		t.Fatal("expected first token to be NumericToken")
	}

	if !scanner.Advance() {
		t.Fatal("expected second Advance to succeed")
	}

	op, ok := scanner.Current().(*tokens.OperatorToken)
	if !ok {
		t.Fatal("expected second token to be OperatorToken")
	}
	if op.Data() != "begincodespacerange" {
		t.Errorf("expected operator data \"begincodespacerange\", got %q", op.Data())
	}

	builder := cmap.NewCharacterMapBuilder()

	err = p.Parse(num, scanner, builder)
	if err != nil {
		t.Fatalf("unexpected error from Parse: %v", err)
	}

	if len(builder.CodespaceRanges()) > 0 {
		t.Errorf("expected empty CodespaceRanges, got %d ranges", len(builder.CodespaceRanges()))
	}
}
