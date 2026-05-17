//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func TestType1FontSimpleIssue807(t *testing.T) {
	path := testutil.GetDocumentPath("Diacritics_export.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
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

	words := page.GetWords()

	if len(words) != 3 {
		t.Fatalf("expected 3 words, got %d", len(words))
	}

	expectedTexts := []string{"Espinosa", "Spínola", "Moraña,"}
	for i, expected := range expectedTexts {
		if words[i].Text != expected {
			t.Errorf("words[%d].Text = %q, want %q", i, words[i].Text, expected)
		}
	}
}
