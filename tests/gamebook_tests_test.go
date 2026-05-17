//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

// TestIssue682 verifies that every page in Gamebook.pdf can be opened and is not nil.
// This matches C# Integration/GamebookTests.Issue682.
func TestIssue682(t *testing.T) {
	path := testutil.GetDocumentPath("Gamebook", true)

	opts := &content.ParsingOptions{
		UseLenientParsing: false,
	}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for i := 0; i < doc.NumberOfPages(); i++ {
		pageAny, err := doc.GetPage(i + 1)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", i+1, err)
		}

		if pageAny == nil {
			t.Errorf("page %d is nil", i+1)
		}
	}
}
