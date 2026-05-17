//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

// TestCanReadDocumentWithMissingWhitespaceAfterXRef verifies that a document with
// no whitespace after the xref table can be opened and its pages counted.
// This matches C# CrossReferenceParserTests.CanReadDocumentWithMissingWhitespaceAfterXRef.
func TestCanReadDocumentWithMissingWhitespaceAfterXRef(t *testing.T) {
	path := testutil.GetSpecificTestDocumentPath("xref-with-no-whitespace.pdf", true)

	opts := &content.ParsingOptions{
		UseLenientParsing: false,
	}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 3 {
		t.Errorf("NumberOfPages = %d, want 3", got)
	}
}

// TestCanReadDocumentWithCircularXRef verifies that a document with circular xref
// references can be opened without entering an infinite loop.
// This matches C# CrossReferenceParserTests.CanReadDocumentWithCircularXRef.
func TestCanReadDocumentWithCircularXRef(t *testing.T) {
	path := testutil.GetSpecificTestDocumentPath("B17-2000-transportation-fuels.pdf", true)

	opts := &content.ParsingOptions{
		UseLenientParsing: false,
	}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 1 {
		t.Errorf("NumberOfPages = %d, want 1", got)
	}
}
