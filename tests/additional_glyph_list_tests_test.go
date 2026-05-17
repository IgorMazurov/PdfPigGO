//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"slices"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func TestType1FontSimple1(t *testing.T) {
	path := testutil.GetDocumentPath("2108.11480", true)

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(2): expected *content.Page, got %T", pageAny)
	}

	letters := page.Letters()
	values := make([]string, 0, len(letters))
	for _, l := range letters {
		values = append(values, l.Value)
	}

	if !slices.Contains(values, "\u22c3") {
		t.Errorf("expected letter value U+22C3 (N-ARY COPRODUCT) not found on page 2")
	}
}

func TestType1FontSimple2(t *testing.T) {
	path := testutil.GetDocumentPath("ICML03-081", true)

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(2): expected *content.Page, got %T", pageAny)
	}

	letters := page.Letters()
	values := make([]string, 0, len(letters))
	for _, l := range letters {
		values = append(values, l.Value)
	}

	expectedChars := []string{"\u2211", "\u220f", "[", "]"}
	for _, ch := range expectedChars {
		if !slices.Contains(values, ch) {
			t.Errorf("expected letter value %q not found on page 2", ch)
		}
	}
}

func TestType1FontSimple3(t *testing.T) {
	path := testutil.GetDocumentPath("Math119FakingData", true)

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(4)
	if err != nil {
		t.Fatalf("GetPage(4): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(4): expected *content.Page, got %T", pageAny)
	}

	letters := page.Letters()
	values := make([]string, 0, len(letters))
	for _, l := range letters {
		values = append(values, l.Value)
	}

	expectedChars := []string{"(", ")", "\u2211"}
	for _, ch := range expectedChars {
		if !slices.Contains(values, ch) {
			t.Errorf("expected letter value %q not found on page 4", ch)
		}
	}
}
