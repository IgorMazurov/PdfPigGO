//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

func popBugzilla37292Path() string {
	return filepath.Join(integrationDocRoot, "pop-bugzilla37292.pdf")
}

// TestPopBugzilla37292CanReadPages verifies that every page can be read and has
// a non-nil Letters collection. This matches C# PopBugzilla37292Tests.CanReadPages.
func TestPopBugzilla37292CanReadPages(t *testing.T) {
	path := popBugzilla37292Path()

	opts := &content.ParsingOptions{
		UseLenientParsing: false,
	}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for i := 0; i < doc.NumberOfPages(); i++ {
		pageAny, err := doc.GetPage(i + 1)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", i+1, err)
		}

		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page, got %T", i+1, pageAny)
		}

		if page.Letters() == nil {
			t.Errorf("page %d: Letters is nil", i+1)
		}
	}
}
