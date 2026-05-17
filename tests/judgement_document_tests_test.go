//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

func resolveJudgementDocumentPath() string {
	return filepath.Join(integrationDocRoot, "Judgement Document.pdf")
}

// TestJudgementDocumentHasCorrectNumberOfPages verifies the document has 13 pages.
func TestJudgementDocumentHasCorrectNumberOfPages(t *testing.T) {
	path := resolveJudgementDocumentPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 13 {
		t.Errorf("NumberOfPages = %d, want 13", got)
	}
}

// TestJudgementDocumentHasCorrectPageContents verifies text content on pages 1 and 2.
func TestJudgementDocumentHasCorrectPageContents(t *testing.T) {
	path := resolveJudgementDocumentPath()

	doc, err := pdfpig.OpenFile(path, nil)
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
		t.Fatal("page is not *content.Page")
	}

	expectedPage1 := "Royal Courts of Justice, Rolls Building Fetter Lane, London, EC4A 1NL"
	if !strings.Contains(page.Text(), expectedPage1) {
		t.Errorf("Page 1 text does not contain %q", expectedPage1)
	}

	pageAny, err = doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	page, ok = pageAny.(*content.Page)
	if !ok {
		t.Fatal("page is not *content.Page")
	}

	expectedPage2 := "The reference to BAR is to another trade organisation of which CMUK was"
	if !strings.Contains(page.Text(), expectedPage2) {
		t.Errorf("Page 2 text does not contain %q", expectedPage2)
	}
}

// TestJudgementDocumentHasCorrectPageSize verifies all 13 pages are A4.
func TestJudgementDocumentHasCorrectPageSize(t *testing.T) {
	path := resolveJudgementDocumentPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for i := 1; i <= 13; i++ {
		pageAny, err := doc.GetPage(i)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", i, err)
		}
		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("Page %d is not *content.Page", i)
		}

		if size := page.Size(); size != content.PageSizeA4 {
			t.Errorf("Page %d Size = %v, want PageSizeA4", i, size)
		}
	}
}
