//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/annotations"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
)

func getCatGeneticsPath() string {
	return filepath.Join(integrationDocRoot, "cat-genetics.pdf")
}

// TestCanReadContent verifies that the cat-genetics document's first page contains
// the text "catus". This matches C# CatGeneticsTests.CanReadContent.
func TestCanReadContent(t *testing.T) {
	doc, err := pdfpig.OpenFile(getCatGeneticsPath(), nil)
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

	if !strings.Contains(page.Text(), "catus") {
		t.Errorf("page text does not contain 'catus'; got: %q", page.Text())
	}
}

// TestCanGetAnnotations verifies that the cat-genetics document has annotations,
// including Highlight-type annotations with QuadPoints. This matches C# CatGeneticsTests.CanGetAnnotations.
func TestCanGetAnnotations(t *testing.T) {
	doc, err := pdfpig.OpenFile(getCatGeneticsPath(), &content.ParsingOptions{UseLenientParsing: false})
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

	annots := page.GetAnnotations()
	if len(annots) == 0 {
		t.Fatal("expected non-empty annotations")
	}

	var highlights []content.LinkAnnotationIface
	for _, a := range annots {
		if at, ok := a.Type().(annotations.AnnotationType); ok && at == annotations.Highlight {
			highlights = append(highlights, a)
		}
	}

	for i, h := range highlights {
		qp := h.QuadPoints()
		if qp == nil {
			t.Errorf("highlight[%d]: expected non-empty QuadPoints, got nil", i)
		}
	}
}

// TestCanSupportPageInformationNotFoundInLenientMode verifies that a document with
// a null Pages entry can be processed in lenient mode but fails in strict mode.
// This matches C# CatGeneticsTests.CanSupportPageInformationNotFoundInLenientMode.
func TestCanSupportPageInformationNotFoundInLenientMode(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "pages-indirect-to-null.pdf")

	// Lenient parsing on -> can process
	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile (lenient): %v", err)
	}

	if doc.NumberOfPages() != 1 {
		t.Errorf("expected 1 page, got %d", doc.NumberOfPages())
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	if pageAny == nil {
		t.Fatal("GetPage(1) returned nil")
	}
	doc.Close()

	// Lenient parsing off -> error with "Pages entry is null"
	_, err = pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err == nil {
		t.Fatal("expected error when opening with lenient parsing off")
	}

	var formatErr *core.PdfDocumentFormatException
	if !errors.As(err, &formatErr) {
		t.Fatalf("expected *core.PdfDocumentFormatException, got %T: %v", err, err)
	}
	if formatErr.Message != "Pages entry is null" {
		t.Errorf("expected message 'Pages entry is null', got %q", formatErr.Message)
	}
}

// TestCanSupportPageKidsObjectNotBeingAPage verifies that a document with a non-page
// object in the pages kids array can be processed in lenient mode but fails in strict mode.
// This matches C# CatGeneticsTests.CanSupportPageKidsObjectNotBeingAPage.
func TestCanSupportPageKidsObjectNotBeingAPage(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "pages-kids-not-page.pdf")

	// Lenient parsing on -> can process
	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile (lenient): %v", err)
	}

	if doc.NumberOfPages() != 1 {
		t.Errorf("expected 1 page, got %d", doc.NumberOfPages())
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	if pageAny == nil {
		t.Fatal("GetPage(1) returned nil")
	}
	doc.Close()

	// Lenient parsing off -> error about pages kids array
	_, err = pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err == nil {
		t.Fatal("expected error when opening with lenient parsing off")
	}

	var formatErr *core.PdfDocumentFormatException
	if !errors.As(err, &formatErr) {
		t.Fatalf("expected *core.PdfDocumentFormatException, got %T: %v", err, err)
	}
	expectedMsg := "Could not find dictionary associated with reference in pages kids array: 3 0."
	if formatErr.Message != expectedMsg {
		t.Errorf("expected message %q, got %q", expectedMsg, formatErr.Message)
	}
}
