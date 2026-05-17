//go:build integration

package geometry_test

import (
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/parser"
)

func BenchmarkClipTwoRectangles(b *testing.B) {
	docPath := filepath.Join(integrationDocRoot, "path_ext_oddeven.pdf")

	opts := &content.ParsingOptions{
		ClipPaths:     false,
		MaxStackDepth: 256,
	}

	doc, err := parser.OpenFile(docPath, opts)
	if err != nil {
		b.Fatalf("OpenFile(%q): %v", docPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		b.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		b.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	paths := page.Paths()
	if len(paths) < 2 {
		b.Fatalf("expected at least 2 paths, got %d", len(paths))
	}

	path0 := paths[0].PdfPath
	path1 := paths[1].PdfPath

	b.Logf("path0: subpaths=%d, isClipping=%v, fillingRule=%v", path0.Len(), path0.IsClipping(), path0.FillingRule())
	b.Logf("path1: subpaths=%d, isClipping=%v, fillingRule=%v", path1.Len(), path1.IsClipping(), path1.FillingRule())

	for i := 0; i < b.N; i++ {
		_ = path0.Clip(path1, logging.NoopLog)
	}
}

func BenchmarkClipRaw(b *testing.B) {
	// Open PDF once to get real paths from the document
	docPath := filepath.Join(integrationDocRoot, "path_ext_oddeven.pdf")

	opts := &content.ParsingOptions{
		ClipPaths:     false,
		MaxStackDepth: 256,
	}

	doc, err := parser.OpenFile(docPath, opts)
	if err != nil {
		b.Fatalf("OpenFile(%q): %v", docPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		b.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		b.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	paths := page.Paths()
	if len(paths) < 2 {
		b.Fatalf("expected at least 2 paths, got %d", len(paths))
	}

	path0 := paths[0].PdfPath
	path1 := paths[1].PdfPath

	// Print path details for debugging
	for i, sp := range path0.Subpaths() {
		b.Logf("path0 subpath[%d]: commands=%d, closed=%v", i, len(sp.Commands()), sp.IsClosed())
	}
	for i, sp := range path1.Subpaths() {
		b.Logf("path1 subpath[%d]: commands=%d, closed=%v", i, len(sp.Commands()), sp.IsClosed())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = geometry.Clip(path0, path1, logging.NoopLog)
	}
}
