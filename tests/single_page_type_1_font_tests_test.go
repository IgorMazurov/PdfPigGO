//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

const singlePageType1FontDocName = "Single Page Type 1 Font.pdf"

func getSinglePageType1FontPath() string {
	return filepath.Join(integrationDocRoot, singlePageType1FontDocName)
}

// TestSinglePageType1FontHasCorrectNumberOfPages verifies the document has 1 page.
func TestSinglePageType1FontHasCorrectNumberOfPages(t *testing.T) {
	path := getSinglePageType1FontPath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 1 {
		t.Errorf("NumberOfPages = %d, want 1", got)
	}
}

// TestSinglePageType1FontHasCorrectPageSize verifies page size is Letter.
func TestSinglePageType1FontHasCorrectPageSize(t *testing.T) {
	path := getSinglePageType1FontPath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
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

	if got := page.Size(); got != content.PageSizeLetter {
		t.Errorf("Size = %v, want %v", got, content.PageSizeLetter)
	}
}

// TestSinglePageType1FontHasCorrectPageText verifies the extracted page text matches expected value.
func TestSinglePageType1FontHasCorrectPageText(t *testing.T) {
	path := getSinglePageType1FontPath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
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

	got := strings.TrimSpace(page.Text())
	want := "PDF test that contains the word yadda"
	if got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
}
