//go:build integration

package pdfpig_test


import (
	"github.com/uglytoad/pdfpig/go/testutil"
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/word_extractor"
)

func TestCanReadDocumentWithIndirectReferences(t *testing.T) {
	path := testutil.GetSpecificTestDocumentPath("93101_1", true)

	opts := &content.ParsingOptions{
		UseLenientParsing: false,
	}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	words := word_extractor.DefaultInstance.GetWords(page.Letters())
	if len(words) == 0 {
		t.Fatal("expected non-empty words on page 1")
	}

	if words[0].Text != "Railway" {
		t.Errorf("words[0].Text = %q, want %q", words[0].Text, "Railway")
	}

	for p := 2; p <= doc.NumberOfPages(); p++ {
		pageAny, err = doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}

		page, ok = pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page, got %T", p, pageAny)
		}

		letters := page.Letters()
		if len(letters) == 0 {
			t.Errorf("page %d: expected non-empty letters", p)
		}
	}
}
