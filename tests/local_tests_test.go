//go:build integration

// Package pdfpig provides integration test stubs for local files not checked into source control.
package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

// TestLocalFiles is a placeholder for testing PDF files that are not in source control.
// To use, set the LOCAL_PDF_DIR environment variable to a directory containing PDFs.
func TestLocalFiles(t *testing.T) {
	dir := os.Getenv("LOCAL_PDF_DIR")
	if dir == "" {
		t.Skip("LOCAL_PDF_DIR not set")
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.pdf"))
	if err != nil {
		t.Fatalf("Glob: %v", err)
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			opts := &content.ParsingOptions{
				UseLenientParsing: false,
			}

			doc, err := pdfpig.OpenFile(file, opts)
			if err != nil {
				t.Fatalf("pdfpig.OpenFile(%q): %v", file, err)
			}
			defer doc.Close()

			for i := 0; i < doc.NumberOfPages(); i++ {
				pageAny, err := doc.GetPage(i + 1)
				if err != nil {
					t.Errorf("GetPage(%d): %v", i+1, err)
					continue
				}

				page, ok := pageAny.(*content.Page)
				if !ok {
					t.Errorf("GetPage(%d): expected *content.Page, got %T", i+1, pageAny)
					continue
				}

				_ = page.Letters
			}
		})
	}
}
