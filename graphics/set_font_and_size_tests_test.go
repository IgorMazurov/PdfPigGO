// Package graphics provides tests for SetFontAndSize (Tf) operation.
package graphics

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

var font1Name = tokens.Create("Font1")

func TestSetFontAndSizeHasCorrectSymbol(t *testing.T) {
	op := NewSetFontAndSize(font1Name, 12)
	if got := op.Operator(); got != "Tf" {
		t.Errorf("expected symbol 'Tf', got '%s'", got)
	}
}

func TestSetFontAndSizeSetsValues(t *testing.T) {
	op := NewSetFontAndSize(font1Name, 12.75)

	if op.Font.Data() != "Font1" {
		t.Errorf("expected font name 'Font1', got '%s'", op.Font.Data())
	}
	if op.Size != 12.75 {
		t.Errorf("expected size 12.75, got %f", op.Size)
	}
}

func TestSetFontAndSizeHasCorrectOperator(t *testing.T) {
	op := NewSetFontAndSize(font1Name, 12)
	if got := op.Operator(); got != "Tf" {
		t.Errorf("expected operator 'Tf', got '%s'", got)
	}
}

func TestSetFontAndSizeStringRepresentationIsCorrect(t *testing.T) {
	op := NewSetFontAndSize(font1Name, 12.76)
	expected := "/Font1 12.76 Tf"
	if got := op.String(); got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}

func TestSetFontAndSizeNameNullPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when font is nil")
		}
	}()
	_ = NewSetFontAndSize(nil, 6)
	t.Fatal("expected panic when font is nil")
}

func TestSetFontAndSizeRunSetsFontAndFontSize(t *testing.T) {
	ctx := NewTestOperationContext()

	op := NewSetFontAndSize(font1Name, 69.42)
	op.Run(ctx)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if state.FontState.FontSize != 69.42 {
		t.Errorf("expected font size 69.42, got %f", state.FontState.FontSize)
	}
	if state.FontState.FontName != font1Name {
		t.Errorf("expected font name Font1, got %v", state.FontState.FontName)
	}
}
