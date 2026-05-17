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

const libreOfficeFontSizeDocName = "Font Size Test - from libre office.pdf"

func resolveLibreOfficeFontSizeDocPath(t *testing.T) string {
	t.Helper()

	path := filepath.Join(integrationDocRoot, libreOfficeFontSizeDocName)
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("document not found: %s", path)
	}

	return path
}

func TestLibreOfficeFontSizeGetsCorrectNumberOfPages(t *testing.T) {
	fullPath := resolveLibreOfficeFontSizeDocPath(t)

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

func TestLibreOfficeFontSizeGetsCorrectPageSize(t *testing.T) {
	fullPath := resolveLibreOfficeFontSizeDocPath(t)

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

func TestLibreOfficeFontSizeGetsCorrectPageTextIgnoringHiddenCharacters(t *testing.T) {
	fullPath := resolveLibreOfficeFontSizeDocPath(t)

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

	expected := "36pt font14 pt font6pt font"
	if text != expected {
		t.Errorf("expected text %q, got %q", expected, text)
	}
}
