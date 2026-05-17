//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

func TestCanReadDocumentWithInvalidObjectLoop(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "invalid-xref-loop.pdf")

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for i := 1; i <= doc.NumberOfPages(); i++ {
		pageAny, err := doc.GetPage(i)
		if err != nil {
			t.Errorf("GetPage(%d): %v", i, err)
			continue
		}

		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Errorf("GetPage(%d): expected *content.Page, got %T", i, pageAny)
			continue
		}

		if page.Letters() == nil {
			t.Errorf("GetPage(%d): Content is nil", i)
		}
	}
}
