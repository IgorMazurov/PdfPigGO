//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

func getIssue881Path() string {
	return filepath.Join(integrationDocRoot, "issue_881.pdf")
}

// TestTrueTypeIssue881 verifies that words extracted from a document using TrueType
// fonts are correctly recognized and their text content matches expected values.
// This matches C# UglyToad.PdfPig.Tests.Integration.TrueTypeTests.Issue881.
func TestTrueTypeIssue881(t *testing.T) {
	path := getIssue881Path()

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

	words := page.GetWords()

	expectedTexts := []string{"IDNR:", "4174", "/", "06.08.2018"}

	if len(words) != len(expectedTexts) {
		t.Fatalf("expected %d words, got %d", len(expectedTexts), len(words))
	}

	for i, expected := range expectedTexts {
		if words[i].Text != expected {
			t.Errorf("words[%d].Text: expected %q, got %q", i, expected, words[i].Text)
		}
	}
}
