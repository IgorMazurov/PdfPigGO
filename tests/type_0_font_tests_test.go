//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

const type0FontDocName = "Type0 Font.pdf"

func getType0FontPath() string {
	return filepath.Join(integrationDocRoot, type0FontDocName)
}

// TestType0FontHasCorrectNumberOfPages verifies the document has 1 page.
func TestType0FontHasCorrectNumberOfPages(t *testing.T) {
	path := getType0FontPath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 1 {
		t.Errorf("NumberOfPages = %d, want 1", got)
	}
}

// TestType0FontHasCorrectPageSize verifies page size is Letter.
func TestType0FontHasCorrectPageSize(t *testing.T) {
	path := getType0FontPath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
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
		t.Fatal("page is not *content.Page")
	}

	if got := page.Size(); got != content.PageSizeLetter {
		t.Errorf("Size = %v, want %v", got, content.PageSizeLetter)
	}
}

// TestType0FontGetsCorrectPageTextIgnoringHiddenCharacters verifies the extracted
// page text contains "Powder River Examiner".
func TestType0FontGetsCorrectPageTextIgnoringHiddenCharacters(t *testing.T) {
	path := getType0FontPath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
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
		t.Fatal("page is not *content.Page")
	}

	var text strings.Builder
	for _, letter := range page.Letters() {
		text.WriteString(letter.Value)
	}

	pageText := text.String()
	expectedSubstr := "Powder River Examiner"
	if !strings.Contains(pageText, expectedSubstr) {
		t.Errorf("page text does not contain %q; got length=%d", expectedSubstr, len(pageText))
	}
}

// TestType0FontHasLetterWidthsAndHeights verifies that letters have nonzero
// bounding box and glyph rectangle dimensions.
func TestType0FontHasLetterWidthsAndHeights(t *testing.T) {
	path := getType0FontPath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
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
		t.Fatal("page is not *content.Page")
	}

	letters := page.Letters()

	hasNonzeroBBWidth := false
	hasNonzeroBBHeight := false
	hasNonzeroGRLooseWidth := false
	hasNonzeroGRLooseHeight := false

	for _, letter := range letters {
		if letter.BoundingBox.Width != 0 {
			hasNonzeroBBWidth = true
		}
		if letter.BoundingBox.Height != 0 {
			hasNonzeroBBHeight = true
		}
		if letter.GlyphRectangleLoose.Width != 0 {
			hasNonzeroGRLooseWidth = true
		}
		if letter.GlyphRectangleLoose.Height != 0 {
			hasNonzeroGRLooseHeight = true
		}
		if hasNonzeroBBWidth && hasNonzeroBBHeight &&
			hasNonzeroGRLooseWidth && hasNonzeroGRLooseHeight {
			break
		}
	}

	if !hasNonzeroBBWidth {
		t.Error("no letter has nonzero BoundingBox.Width")
	}
	if !hasNonzeroBBHeight {
		t.Error("no letter has nonzero BoundingBox.Height")
	}
	if !hasNonzeroGRLooseWidth {
		t.Error("no letter has nonzero GlyphRectangleLoose.Width")
	}
	if !hasNonzeroGRLooseHeight {
		t.Error("no letter has nonzero GlyphRectangleLoose.Height")
	}
}
