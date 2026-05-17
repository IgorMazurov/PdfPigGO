//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

var documentsToIgnore = map[string]bool{
	"issue_671.pdf":            true,
	"GHOSTSCRIPT-698363-0.pdf": true,
	"ErcotFacts.pdf":           true,
	"cmap-parsing-exception.pdf": true,
}

func getIntegrationDocuments(t *testing.T) []string {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(integrationDocRoot, "*.pdf"))
	if err != nil {
		t.Fatalf("Glob: %v", err)
	}

	result := make([]string, 0, len(files))
	for _, f := range files {
		name := filepath.Base(f)
		if !documentsToIgnore[name] {
			result = append(result, name)
		}
	}

	return result
}

func resolveIntegrationDocumentPath(t *testing.T, documentName string) string {
	t.Helper()

	path := filepath.Join(integrationDocRoot, documentName)
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("document not found: %s", path)
	}

	return path
}

// TestCheckGlyphLooseBoundingBoxes verifies that for every letter on every page,
// if the bounding box height is positive then the glyph loose rectangle height
// is also positive. This matches C# IntegrationDocumentTests.CheckGlyphLooseBoundingBoxes.
func TestCheckGlyphLooseBoundingBoxes(t *testing.T) {
	documents := getIntegrationDocuments(t)

	for _, documentName := range documents {
		t.Run(documentName, func(t *testing.T) {
			fullPath := resolveIntegrationDocumentPath(t, documentName)

			opts := &content.ParsingOptions{
				UseLenientParsing: true,
			}

			doc, err := pdfpig.OpenFile(fullPath, opts)
			if err != nil {
				t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
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

				letters := page.Letters()
				for _, letter := range letters {
					bbHeight := letter.BoundingBox.Height
					if bbHeight > 0 {
						glyphLooseHeight := letter.GlyphRectangleLoose.Height
						if glyphLooseHeight <= 0 {
							font := letter.GetFont()
							if font != nil {
								_ = font.GetAscent()
							}
						}

						if glyphLooseHeight <= 0 {
							t.Errorf("page %d: GlyphRectangleLoose.Height (%g) is not positive while BoundingBox.Height (%g) is positive", i+1, glyphLooseHeight, bbHeight)
						}
					}
				}
			}
		})
	}
}

// TestCanReadAllPages verifies that every page in each document can be read and
// annotations can be retrieved without error. This matches C# IntegrationDocumentTests.CanReadAllPages.
func TestCanReadAllPages(t *testing.T) {
	documents := getIntegrationDocuments(t)

	for _, documentName := range documents {
		t.Run(documentName, func(t *testing.T) {
			fullPath := resolveIntegrationDocumentPath(t, documentName)

			opts := &content.ParsingOptions{
				UseLenientParsing: false,
			}

			doc, err := pdfpig.OpenFile(fullPath, opts)
			if err != nil {
				t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
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

				annotations := page.GetAnnotations()
				if annotations == nil {
					t.Errorf("page %d: GetAnnotations returned nil", i+1)
				}
			}
		})
	}
}

// TestCanUseStreamForFirstPage verifies that documents can be opened from a byte slice
// (simulating a memory stream) and all pages are readable. This matches C# IntegrationDocumentTests.CanUseStreamForFirstPage.
func TestCanUseStreamForFirstPage(t *testing.T) {
	documents := getIntegrationDocuments(t)

	for _, documentName := range documents {
		t.Run(documentName, func(t *testing.T) {
			fullPath := resolveIntegrationDocumentPath(t, documentName)

			bytes, err := os.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("ReadFile(%q): %v", fullPath, err)
			}

			opts := &content.ParsingOptions{
				UseLenientParsing: false,
			}

			doc, err := pdfpig.Open(bytes, opts)
			if err != nil {
				t.Fatalf("Open(%q): %v", fullPath, err)
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

				annotations := page.GetAnnotations()
				if annotations == nil {
					t.Errorf("page %d: GetAnnotations returned nil", i+1)
				}
			}
		})
	}
}

// TestCanTokenizeAllAccessibleObjects verifies that the document structure catalog
// is accessible after opening each document. This matches C# IntegrationDocumentTests.CanTokenizeAllAccessibleObjects.
func TestCanTokenizeAllAccessibleObjects(t *testing.T) {
	documents := getIntegrationDocuments(t)

	for _, documentName := range documents {
		t.Run(documentName, func(t *testing.T) {
			fullPath := resolveIntegrationDocumentPath(t, documentName)

			opts := &content.ParsingOptions{
				UseLenientParsing: false,
			}

			doc, err := pdfpig.OpenFile(fullPath, opts)
			if err != nil {
				t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
			}
			defer doc.Close()

			if doc.Structure == nil {
				t.Fatal("document structure is nil")
			}

			catalog := doc.Structure.Catalog()
			if catalog == nil {
				t.Fatal("catalog is nil")
			}
		})
	}
}

// TestCanAccessImagesOnEveryPage verifies that images on every page can be accessed
// and have positive width and height. This matches C# IntegrationDocumentTests.CanAccessImagesOnEveryPage.
func TestCanAccessImagesOnEveryPage(t *testing.T) {
	documents := getIntegrationDocuments(t)

	for _, documentName := range documents {
		t.Run(documentName, func(t *testing.T) {
			fullPath := resolveIntegrationDocumentPath(t, documentName)

			opts := &content.ParsingOptions{
				UseLenientParsing: false,
			}

			doc, err := pdfpig.OpenFile(fullPath, opts)
			if err != nil {
				t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
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

				images := page.GetImages()
				if images == nil {
					t.Errorf("page %d: GetImages returned nil", i+1)
					continue
				}

				for _, img := range images {
					if img.WidthInSamples() <= 0 {
						t.Errorf("page %d: image had width of zero (%d)", i+1, img.WidthInSamples())
					}
					if img.HeightInSamples() <= 0 {
						t.Errorf("page %d: image had height of zero (%d)", i+1, img.HeightInSamples())
					}
				}
			}
		})
	}
}
