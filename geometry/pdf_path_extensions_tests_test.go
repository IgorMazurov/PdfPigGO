//go:build integration

package geometry_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/parser"
)

func TestPathExtensionsContainsRectangleEvenOdd(t *testing.T) {
	docPath := filepath.Join(integrationDocRoot, "path_ext_oddeven.pdf")

	opts := &content.ParsingOptions{
		ClipPaths:     true,
		MaxStackDepth: 256,
	}

	doc, err := parser.OpenFile(docPath, opts)
	if err != nil {
		t.Fatalf("OpenFile(%q): %v", docPath, err)
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

	words := page.GetWords()

	for _, path := range page.Paths() {
		if path.FillingRule() == core.FillingRuleNonZeroWinding {
			t.Errorf("expected filling rule != NonZeroWinding, got NonZeroWinding")
		}

		var remaining []*content.Word
		for _, w := range words {
			if geometry.PathContainsRect(path.PdfPath, w.BoundingBox(), true) {
				parts := strings.Split(w.Text, "_")
				lastPart := parts[len(parts)-1]
				if lastPart != "in" {
					t.Errorf("word %q inside path expected to end with 'in', got '%s'", w.Text, lastPart)
				}
			} else {
				remaining = append(remaining, w)
			}
		}
		words = remaining
	}

	for _, w := range words {
		parts := strings.Split(w.Text, "_")
		lastPart := parts[len(parts)-1]
		if lastPart == "in" {
			t.Errorf("remaining word %q should not end with 'in'", w.Text)
		}
	}
}
