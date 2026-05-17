package parts_test

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/testutil"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestTryGetCanFollowMultipleReferenceLinks(t *testing.T) {
	scanner := testutil.NewTestPdfTokenScanner()

	ref1, _ := core.NewIndirectReference(7, 0)
	ref2, _ := core.NewIndirectReference(9, 0)

	scanner.Objects[ref1] = tokens.NewObjectToken(core.File(10), ref1, tokens.NewIndirectReferenceToken(ref2))
	scanner.Objects[ref2] = tokens.NewObjectToken(core.File(12), ref2, tokens.NewNumericToken(69))

	result, ok := parts.TryGet[*tokens.NumericToken](tokens.NewIndirectReferenceToken(ref1), scanner)
	if !ok {
		t.Fatal("expected TryGet to succeed following multiple reference links")
	}

	if result == nil {
		t.Fatal("expected non-nil NumericToken result")
	}

	if got := int(result.Data()); got != 69 {
		t.Errorf("expected result.Data() == 69, got %d", got)
	}
}

func TestGetByRefCanFollowMultipleReferenceLinks(t *testing.T) {
	scanner := testutil.NewTestPdfTokenScanner()

	ref1, _ := core.NewIndirectReference(7, 0)
	ref2, _ := core.NewIndirectReference(9, 0)

	scanner.Objects[ref1] = tokens.NewObjectToken(core.File(10), ref1, tokens.NewIndirectReferenceToken(ref2))
	scanner.Objects[ref2] = tokens.NewObjectToken(core.File(12), ref2, tokens.NewNumericToken(69))

	result, err := parts.GetByRef[*tokens.NumericToken](ref1, scanner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil NumericToken result")
	}

	if got := int(result.Data()); got != 69 {
		t.Errorf("expected result.Data() == 69, got %d", got)
	}
}

func TestGetByTokenCanFollowMultipleReferenceLinks(t *testing.T) {
	scanner := testutil.NewTestPdfTokenScanner()

	ref1, _ := core.NewIndirectReference(7, 0)
	ref2, _ := core.NewIndirectReference(9, 0)

	scanner.Objects[ref1] = tokens.NewObjectToken(core.File(10), ref1, tokens.NewIndirectReferenceToken(ref2))
	scanner.Objects[ref2] = tokens.NewObjectToken(core.File(12), ref2, tokens.NewNumericToken(69))

	result, err := parts.GetByToken[*tokens.NumericToken](tokens.NewIndirectReferenceToken(ref1), scanner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil NumericToken result")
	}

	if got := int(result.Data()); got != 69 {
		t.Errorf("expected result.Data() == 69, got %d", got)
	}
}

func TestGetByRefReturnsSingleItemFromArray(t *testing.T) {
	scanner := testutil.NewTestPdfTokenScanner()

	ref, _ := core.NewIndirectReference(10, 0)

	expected := "Goopy"
	scanner.Objects[ref] = tokens.NewObjectToken(core.File(10), ref, tokens.NewArrayToken([]tokens.Token{
		tokens.NewStringToken(expected),
	}))

	result, err := parts.GetByRef[*tokens.StringToken](ref, scanner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil StringToken result")
	}

	if got := result.Data(); got != expected {
		t.Errorf("expected result.Data() == %q, got %q", expected, got)
	}
}

func TestGetByRefFollowsSingleIndirectReferenceFromArray(t *testing.T) {
	scanner := testutil.NewTestPdfTokenScanner()

	ref, _ := core.NewIndirectReference(10, 0)
	ref2, _ := core.NewIndirectReference(69, 0)

	expected := "Goopy"
	scanner.Objects[ref] = tokens.NewObjectToken(core.File(10), ref, tokens.NewArrayToken([]tokens.Token{
		tokens.NewIndirectReferenceToken(ref2),
	}))

	scanner.Objects[ref2] = tokens.NewObjectToken(core.File(69), ref2, tokens.NewStringToken(expected))

	result, err := parts.GetByRef[*tokens.StringToken](ref, scanner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil StringToken result")
	}

	if got := result.Data(); got != expected {
		t.Errorf("expected result.Data() == %q, got %q", expected, got)
	}
}

func TestGetByRefThrowsOnInvalidArray(t *testing.T) {
	scanner := testutil.NewTestPdfTokenScanner()

	ref, _ := core.NewIndirectReference(10, 0)

	scanner.Objects[ref] = tokens.NewObjectToken(core.File(10), ref, tokens.NewArrayToken([]tokens.Token{
		tokens.NewNumericToken(5),
		tokens.NewNumericToken(6),
		tokens.NewNumericToken(0),
	}))

	_, err := parts.GetByRef[*tokens.StringToken](ref, scanner)

	if err == nil {
		t.Fatal("expected error for invalid array type mismatch")
	}

	_, ok := err.(*core.PdfDocumentFormatException)
	if !ok {
		t.Errorf("expected PdfDocumentFormatException, got %T: %v", err, err)
	}
}
