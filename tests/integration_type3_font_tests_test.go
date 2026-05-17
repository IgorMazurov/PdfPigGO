//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

func getType3FontZeroHeightPath() string {
	return filepath.Join(integrationDocRoot, "type3-font-zero-height.pdf")
}

// TestType3FontHasLetterWidthsAndHeights verifies that letters extracted from a document
// using Type3 fonts have non-zero bounding box width and height.
// This matches C# UglyToad.PdfPig.Tests.Integration.Type3FontTests.HasLetterWidthsAndHeights.
func TestType3FontHasLetterWidthsAndHeights(t *testing.T) {
	path := getType3FontZeroHeightPath()

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

	letters := page.Letters()

	hasNonZeroWidth := false
	hasNonZeroHeight := false

	for _, letter := range letters {
		if letter.BoundingBox.Width != 0 {
			hasNonZeroWidth = true
		}
		if letter.BoundingBox.Height != 0 {
			hasNonZeroHeight = true
		}
	}

	if !hasNonZeroWidth {
		t.Error("expected at least one letter with non-zero BoundingBox.Width")
	}

	if !hasNonZeroHeight {
		t.Error("expected at least one letter with non-zero BoundingBox.Height")
	}
}
