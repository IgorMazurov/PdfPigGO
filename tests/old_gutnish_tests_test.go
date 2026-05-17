//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func getOldGutnishFilename() string {
	return testutil.GetDocumentPath("Old Gutnish Internet Explorer.pdf", true)
}

// TestOldGutnishHasCorrectNumberOfPages verifies the Old Gutnish PDF has 3 pages.
func TestOldGutnishHasCorrectNumberOfPages(t *testing.T) {
	doc, err := pdfpig.OpenFile(getOldGutnishFilename(), &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 3 {
		t.Errorf("NumberOfPages = %d, want 3", got)
	}
}

// TestOldGutnishHasCorrectContentAfterReadingPreviousPage verifies that reading
// page 1 then page 2 yields the expected Old Gutnish text on each page.
func TestOldGutnishHasCorrectContentAfterReadingPreviousPage(t *testing.T) {
	doc, err := pdfpig.OpenFile(getOldGutnishFilename(), &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page1, ok := page1Any.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", page1Any)
	}

	expectedPage1 := "Þissi þieluar hafþi ann sun sum hit hafþi. En hafþa cuna hit huita stierna"
	if !strings.Contains(page1.Text(), expectedPage1) {
		t.Errorf("page 1 text does not contain expected substring.\n"+
			"want substring: %q\n"+
			"got text:\n%s", expectedPage1, page1.Text())
	}

	page2Any, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}

	page2, ok := page2Any.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(2): expected *content.Page, got %T", page2Any)
	}

	expectedPage2 := "Greipur sem annar hét; og Gunnfjón sá þriðji"
	if !strings.Contains(page2.Text(), expectedPage2) {
		t.Errorf("page 2 text does not contain expected substring.\n"+
			"want substring: %q\n"+
			"got text:\n%s", expectedPage2, page2.Text())
	}
}

// TestOldGutnishGetsImageOnPageOne verifies that page 1 contains exactly one image.
func TestOldGutnishGetsImageOnPageOne(t *testing.T) {
	doc, err := pdfpig.OpenFile(getOldGutnishFilename(), &content.ParsingOptions{})
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

	images := page.GetImages()
	if len(images) != 1 {
		t.Errorf("page 1 GetImages count = %d, want 1", len(images))
	}
}
