//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

func getPlosOneFilename() string {
	return filepath.Join(integrationDocRoot, "journal.pone.0196757.pdf")
}

// TestPlosOneCanReadPageOneContent verifies that page 1 of the PLoS ONE article
// can be opened and its extracted text has more than 50 characters.
// This matches C# UglyToad.PdfPig.Tests.Integration.PlosOneTests.CanReadPageOneContent.
func TestPlosOneCanReadPageOneContent(t *testing.T) {
	path := getPlosOneFilename()

	opts := &content.ParsingOptions{
		UseLenientParsing: false,
	}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
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

	text := page.Text()
	if len(text) <= 50 {
		t.Errorf("expected page text length > 50, got %d", len(text))
	}
}
