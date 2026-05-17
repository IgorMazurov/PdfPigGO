//go:build integration

package geometry_test

import (
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/parser"
)

const integrationDocRoot = "testdata/integration/Documents"

func TestContainsRectangleEvenOdd(t *testing.T) {
	docPath := filepath.Join(integrationDocRoot, "SPARC - v9 Architecture Manual.pdf")

	opts := &content.ParsingOptions{
		ClipPaths: true,
	}

	doc, err := parser.OpenFile(docPath, opts)
	if err != nil {
		t.Fatalf("OpenFile(%q): %v", docPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(45)
	if err != nil {
		t.Fatalf("GetPage(45): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(45): expected *content.Page, got %T", pageAny)
	}

	paths := page.Paths()
	if len(paths) != 28 {
		t.Errorf("expected 28 paths on page 45, got %d", len(paths))
	}
}
