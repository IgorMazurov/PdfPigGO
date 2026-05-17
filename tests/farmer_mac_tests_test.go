//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func getFarmerMacFilename() string {
	return testutil.GetDocumentPath("FarmerMac.pdf", true)
}

// TestFarmerMacHasCorrectNumberOfPages verifies the FarmerMac PDF has 5 pages.
func TestFarmerMacHasCorrectNumberOfPages(t *testing.T) {
	doc, err := pdfpig.OpenFile(getFarmerMacFilename(), &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 5 {
		t.Errorf("NumberOfPages = %d, want 5", got)
	}
}

// TestFarmerMacGetPagesWorks verifies GetPages returns the correct page count.
func TestFarmerMacGetPagesWorks(t *testing.T) {
	doc, err := pdfpig.OpenFile(getFarmerMacFilename(), &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	pages, err := doc.GetPages()
	if err != nil {
		t.Fatalf("GetPages: %v", err)
	}

	if got := len(pages); got != 5 {
		t.Errorf("len(GetPages()) = %d, want 5", got)
	}
}

// TestFarmerMacHasCorrectContentAfterReadingPreviousPage verifies that reading
// page 1 then page 2 yields the expected text on page 2.
func TestFarmerMacHasCorrectContentAfterReadingPreviousPage(t *testing.T) {
	doc, err := pdfpig.OpenFile(getFarmerMacFilename(), &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page2Any, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}

	page1, ok := page1Any.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", page1Any)
	}

	page2, ok := page2Any.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(2): expected *content.Page, got %T", page2Any)
	}

	expected := "financial results for the fiscal quarter ended June 30, 2017 and (2) a conference call to discuss those results and Farmer Mac"

	if !strings.Contains(page2.Text(), expected) {
		t.Errorf("page 2 text does not contain expected substring.\n"+
			"want substring: %q\n"+
			"got text:\n%s", expected, page2.Text())
	}

	_ = page1
}
