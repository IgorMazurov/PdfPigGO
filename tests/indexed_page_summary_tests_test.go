//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"strconv"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func indexedPageSummaryFilename() string {
	return testutil.GetDocumentPath("FICTIF_TABLE_INDEX.pdf", true)
}

// TestIndexedPageSummaryHasCorrectNumberOfPages verifies the FICTIF_TABLE_INDEX PDF has 14 pages.
func TestIndexedPageSummaryHasCorrectNumberOfPages(t *testing.T) {
	doc, err := pdfpig.OpenFile(indexedPageSummaryFilename(), &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 14 {
		t.Errorf("NumberOfPages = %d, want 14", got)
	}
}

// TestIndexedPageSummaryGetPagesWorks verifies GetPages returns the correct page count.
func TestIndexedPageSummaryGetPagesWorks(t *testing.T) {
	doc, err := pdfpig.OpenFile(indexedPageSummaryFilename(), &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	pages, err := doc.GetPages()
	if err != nil {
		t.Fatalf("GetPages: %v", err)
	}

	if got := len(pages); got != 14 {
		t.Errorf("len(GetPages()) = %d, want 14", got)
	}
}

type indexedNameTest struct {
	name       string
	pageNumber int
}

var indexedNameTestData = []indexedNameTest{
	{"M. HERNANDEZ DANIEL", 1},
	{"M. HERNANDEZ DANIEL", 2},
	{"Mme ALIBERT CHLOE AA", 3},
	{"Mme ALIBERT CHLOE AA", 4},
	{"M. SIMPSON BART AAA", 5},
	{"M. SIMPSON BART AAA", 6},
	{"M. BOND JAMES A", 7},
	{"M. BOND JAMES A", 8},
	{"M. DE BALZAC HONORE", 9},
	{"M. DE BALZAC HONORE", 10},
	{"M. STALLONE SILVESTER", 11},
	{"M. STALLONE SILVESTER", 12},
	{"M. SCOTT MICHAEL", 13},
	{"M. SCOTT MICHAEL", 14},
}

// TestIndexedPageSummaryCheckSpecificNamesPresence verifies that specific names
// appear on their expected indexed page numbers in FICTIF_TABLE_INDEX.pdf.
func TestIndexedPageSummaryCheckSpecificNamesPresence(t *testing.T) {
	for _, tc := range indexedNameTestData {
		t.Run(tc.name+"-p"+strconv.Itoa(tc.pageNumber), func(t *testing.T) {
			doc, err := pdfpig.OpenFile(indexedPageSummaryFilename(), &content.ParsingOptions{})
			if err != nil {
				t.Fatalf("OpenFile: %v", err)
			}
			defer doc.Close()

			pageAny, err := doc.GetPage(tc.pageNumber)
			if err != nil {
				t.Fatalf("GetPage(%d): %v", tc.pageNumber, err)
			}

			page, ok := pageAny.(*content.Page)
			if !ok {
				t.Fatalf("GetPage(%d): expected *content.Page, got %T", tc.pageNumber, pageAny)
			}

			if !strings.Contains(page.Text(), tc.name) {
				t.Errorf("page %d text does not contain %q", tc.pageNumber, tc.name)
			}
		})
	}
}
