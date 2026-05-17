//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/outline"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func pigReproductionPowerpointPath() string {
	return testutil.GetDocumentPath("Pig Reproduction Powerpoint.pdf", true)
}

func TestPigReproductionPowerpointCanReadContent(t *testing.T) {
	path := pigReproductionPowerpointPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page := pageAny.(*content.Page)

	expectedText := "Pigs per sow per year: 18 to 27"
	if !strings.Contains(page.Text(), expectedText) {
		t.Errorf("page text does not contain %q", expectedText)
	}
}

func TestPigReproductionPowerpointHasCorrectNumberOfPages(t *testing.T) {
	path := pigReproductionPowerpointPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	expectedPages := 35
	if got := doc.NumberOfPages(); got != expectedPages {
		t.Errorf("expected %d pages, got %d", expectedPages, got)
	}
}

func TestPigReproductionPowerpointCanReadAllPages(t *testing.T) {
	path := pigReproductionPowerpointPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for i := 1; i <= doc.NumberOfPages(); i++ {
		_, err := doc.GetPage(i)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", i, err)
		}
	}
}

func TestPigReproductionPowerpointCanGetBookmarks(t *testing.T) {
	path := pigReproductionPowerpointPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	bookmarksAny, found, err := doc.TryGetBookmarks(false)
	if err != nil {
		t.Fatalf("TryGetBookmarks: %v", err)
	}
	if !found {
		t.Fatal("expected bookmarks to be found")
	}

	bookmarks, ok := bookmarksAny.(*outline.Bookmarks)
	if !ok {
		t.Fatalf("expected *outline.Bookmarks; got %T", bookmarksAny)
	}

	expectedCount := 35
	if got := len(bookmarks.Roots()); got != expectedCount {
		t.Errorf("expected %d bookmark roots, got %d", expectedCount, got)
	}

	allNodes := bookmarks.GetNodes()
	if got := len(allNodes); got != expectedCount {
		t.Errorf("expected %d total bookmark nodes, got %d", expectedCount, got)
	}
}
