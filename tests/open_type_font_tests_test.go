//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	dla "github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/word_extractor"
	"github.com/uglytoad/pdfpig/go/testutil"
)

// TestIssue672 verifies that a PDF with an embedded OpenType CFF font parsed
// by TrueTypeFontParser can still be opened and text extracted without error.
func TestIssue672(t *testing.T) {
	path := testutil.GetDocumentPath("Why.does.this.not.work", true)

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
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

	blocks := page_segmenter.Instance.GetBlocks(words)

	var lines []*dla.TextLine
	for _, block := range blocks {
		lines = append(lines, block.TextLines...)
	}

	if len(lines) != 3 {
		t.Fatalf("lines count = %d, want 3", len(lines))
	}

	expected := []string{
		"THIS TEST SEEMS TO BREAK THE PARSER....",
		"This is just some test text.",
		"SO DOES THIS",
	}

	for i, exp := range expected {
		if lines[i].Text != exp {
			t.Errorf("lines[%d].Text = %q, want %q", i, lines[i].Text, exp)
		}
	}
}

// TestIssue672ok verifies that a PDF with multiple embedded OpenType fonts
// (Semplicita Pro, Verdana) can be opened and text extracted correctly.
func TestIssue672ok(t *testing.T) {
	path := testutil.GetDocumentPath("Test.Doc", true)

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
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

	blocks := page_segmenter.Instance.GetBlocks(words)

	var lines []*dla.TextLine
	for _, block := range blocks {
		lines = append(lines, block.TextLines...)
	}

	if len(lines) != 4 {
		t.Fatalf("lines count = %d, want 4", len(lines))
	}

	expected := []string{
		"This is just a bunch of boring text...",
		"THIS IS SOME SEMPLICITA PRO FONT",
		"Hopefully font that are not embedded on the server.",
		"And a bit of Verdana for good measure.",
	}

	for i, exp := range expected {
		if lines[i].Text != exp {
			t.Errorf("lines[%d].Text = %q, want %q", i, lines[i].Text, exp)
		}
	}
}

// TestSo74165171 verifies that a PDF with an embedded OpenType CFF font can be
// opened and words extracted without error. Regression test for
// https://stackoverflow.com/questions/74165171/embedded-opentype-cff-font-in-a-pdf-shows-strange-behaviour-in-some-viewers
func TestSo74165171(t *testing.T) {
	path := testutil.GetDocumentPath("test-2_so_74165171", true)

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
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

	if len(words) != 2 {
		t.Errorf("words count = %d, want 2", len(words))
	}
}
