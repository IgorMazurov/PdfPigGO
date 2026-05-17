//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

const singlePageSimpleOpenOfficeDocName = "Single Page Simple - from open office.pdf"

func getSinglePageSimpleOpenOfficePath() string {
	return filepath.Join(integrationDocRoot, singlePageSimpleOpenOfficeDocName)
}

// TestSinglePageSimpleOpenOfficeHasCorrectNumberOfPages verifies the document has 1 page.
func TestSinglePageSimpleOpenOfficeHasCorrectNumberOfPages(t *testing.T) {
	path := getSinglePageSimpleOpenOfficePath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 1 {
		t.Errorf("NumberOfPages = %d, want 1", got)
	}
}

// TestSinglePageSimpleOpenOfficeHasCorrectPageSize verifies page size is Letter.
func TestSinglePageSimpleOpenOfficeHasCorrectPageSize(t *testing.T) {
	path := getSinglePageSimpleOpenOfficePath()

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

// TestSinglePageSimpleOpenOfficeHasCorrectLetterBoundingBoxes verifies letter bounding box positions.
func TestSinglePageSimpleOpenOfficeHasCorrectLetterBoundingBoxes(t *testing.T) {
	path := getSinglePageSimpleOpenOfficePath()

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
	comparer := testutil.NewDoubleComparer(3)

	if len(letters) == 0 {
		t.Fatal("expected at least one letter")
	}

	// Letter[0] = "I"
	if got := letters[0].Value; got != "I" {
		t.Errorf("Letters[0].Value = %q, want %q", got, "I")
	}

	if !comparer.Equals(letters[0].BoundingBox.BottomLeft.X, 90.1) {
		t.Errorf("Letters[0].BoundingBox.BottomLeft.X = %g, want ~90.1", letters[0].BoundingBox.BottomLeft.X)
	}

	if !comparer.Equals(letters[0].BoundingBox.BottomLeft.Y, 709.2) {
		t.Errorf("Letters[0].BoundingBox.BottomLeft.Y = %g, want ~709.2", letters[0].BoundingBox.BottomLeft.Y)
	}

	if !comparer.Equals(letters[0].BoundingBox.TopRight.X, 94.0) {
		t.Errorf("Letters[0].BoundingBox.TopRight.X = %g, want ~94.0", letters[0].BoundingBox.TopRight.X)
	}

	if !comparer.Equals(letters[0].BoundingBox.TopRight.Y, 719.89) {
		t.Errorf("Letters[0].BoundingBox.TopRight.Y = %g, want ~719.89", letters[0].BoundingBox.TopRight.Y)
	}

	// Letter[5] = "a"
	if len(letters) < 6 {
		t.Fatalf("expected at least 6 letters, got %d", len(letters))
	}

	if got := letters[5].Value; got != "a" {
		t.Errorf("Letters[5].Value = %q, want %q", got, "a")
	}

	if !comparer.Equals(letters[5].BoundingBox.BottomLeft.X, 114.5) {
		t.Errorf("Letters[5].BoundingBox.BottomLeft.X = %g, want ~114.5", letters[5].BoundingBox.BottomLeft.X)
	}

	if !comparer.Equals(letters[5].BoundingBox.BottomLeft.Y, 709.2) {
		t.Errorf("Letters[5].BoundingBox.BottomLeft.Y = %g, want ~709.2", letters[5].BoundingBox.BottomLeft.Y)
	}

	if !comparer.Equals(letters[5].BoundingBox.TopRight.X, 119.82) {
		t.Errorf("Letters[5].BoundingBox.TopRight.X = %g, want ~119.82", letters[5].BoundingBox.TopRight.X)
	}

	if !comparer.Equals(letters[5].BoundingBox.TopRight.Y, 714.89) {
		t.Errorf("Letters[5].BoundingBox.TopRight.Y = %g, want ~714.89", letters[5].BoundingBox.TopRight.Y)
	}

	// Letter[16] = "f"
	if len(letters) < 17 {
		t.Fatalf("expected at least 17 letters, got %d", len(letters))
	}

	if got := letters[16].Value; got != "f" {
		t.Errorf("Letters[16].Value = %q, want %q", got, "f")
	}

	if !comparer.Equals(letters[16].BoundingBox.BottomLeft.X, 169.9) {
		t.Errorf("Letters[16].BoundingBox.BottomLeft.X = %g, want ~169.9", letters[16].BoundingBox.BottomLeft.X)
	}

	if !comparer.Equals(letters[16].BoundingBox.BottomLeft.Y, 709.2) {
		t.Errorf("Letters[16].BoundingBox.BottomLeft.Y = %g, want ~709.2", letters[16].BoundingBox.BottomLeft.Y)
	}

	if !comparer.Equals(letters[16].BoundingBox.TopRight.X, 176.89) {
		t.Errorf("Letters[16].BoundingBox.TopRight.X = %g, want ~176.89", letters[16].BoundingBox.TopRight.X)
	}

	if !comparer.Equals(letters[16].BoundingBox.TopRight.Y, 719.89) {
		t.Errorf("Letters[16].BoundingBox.TopRight.Y = %g, want ~719.89", letters[16].BoundingBox.TopRight.Y)
	}
}

// TestSinglePageSimpleOpenOfficeGetsCorrectPageTextIgnoringHiddenCharacters verifies extracted text.
func TestSinglePageSimpleOpenOfficeGetsCorrectPageTextIgnoringHiddenCharacters(t *testing.T) {
	path := getSinglePageSimpleOpenOfficePath()

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

	var sb strings.Builder
	for _, letter := range page.Letters() {
		sb.WriteString(letter.Value)
	}
	text := sb.String()

	if got := text; got != "I am a simple pdf." {
		t.Errorf("text = %q, want %q", got, "I am a simple pdf.")
	}
}

// TestSinglePageSimpleOpenOfficeTryGetBookmarksFalse verifies document has no bookmarks.
func TestSinglePageSimpleOpenOfficeTryGetBookmarksFalse(t *testing.T) {
	path := getSinglePageSimpleOpenOfficePath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	_, found, err := doc.TryGetBookmarks(false)
	if err != nil {
		t.Fatalf("TryGetBookmarks: %v", err)
	}

	if found {
		t.Error("expected TryGetBookmarks to return false, got true")
	}
}

// TestSinglePageSimpleOpenOfficeStartXRefNotNearEnd verifies parsing with appended empty trailer bytes.
func TestSinglePageSimpleOpenOfficeStartXRefNotNearEnd(t *testing.T) {
	path := getSinglePageSimpleOpenOfficePath()

	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", path, err)
	}

	emptyTrailer := make([]byte, 2026)
	emptyTrailer[0] = 10

	bytes = append(bytes, emptyTrailer...)

	doc, err := pdfpig.Open(bytes, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("Open(bytes): %v", err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 1 {
		t.Errorf("NumberOfPages = %d, want 1", got)
	}
}
