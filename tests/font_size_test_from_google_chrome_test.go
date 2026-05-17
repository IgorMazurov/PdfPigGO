//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

const fontSizeTestDocName = "Font Size Test - from google chrome print pdf.pdf"

func resolveFontSizeTestDocumentPath(t *testing.T) string {
	t.Helper()

	path := filepath.Join(integrationDocRoot, fontSizeTestDocName)
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("document not found: %s", path)
	}

	return path
}

func TestGetsCorrectNumberOfPages(t *testing.T) {
	fullPath := resolveFontSizeTestDocumentPath(t)

	doc, err := pdfpig.OpenFile(fullPath, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
	}
	defer doc.Close()

	pageCount := doc.NumberOfPages()
	if pageCount != 1 {
		t.Errorf("expected 1 page, got %d", pageCount)
	}
}

func TestGetsCorrectPageWidthAndHeight(t *testing.T) {
	fullPath := resolveFontSizeTestDocumentPath(t)

	doc, err := pdfpig.OpenFile(fullPath, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
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

	if page.Width() != 595 {
		t.Errorf("expected width 595, got %g", page.Width())
	}

	if page.Height() != 842 {
		t.Errorf("expected height 842, got %g", page.Height())
	}
}

func TestGetsCorrectPageTextIgnoringHiddenCharacters(t *testing.T) {
	fullPath := resolveFontSizeTestDocumentPath(t)

	doc, err := pdfpig.OpenFile(fullPath, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
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

	letters := page.Letters()
	text := strings.Clone("")
	for _, letter := range letters {
		text += letter.Value
	}

	expected := "Hello, this is 16ptHello, this is 16px"
	if text != expected {
		t.Errorf("expected text %q, got %q", expected, text)
	}
}

func TestGetsCorrectPageSize(t *testing.T) {
	fullPath := resolveFontSizeTestDocumentPath(t)

	doc, err := pdfpig.OpenFile(fullPath, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
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

	if page.Size() != content.PageSizeA4 {
		t.Errorf("expected PageSizeA4, got %d", page.Size())
	}
}
