//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

func getInvalidOperatorPath() string {
	return filepath.Join(specificTestDocRoot, "invalid-operator.pdf")
}

// TestInvalidOperatorThrowsExceptionIfNotUsingLenientParsing verifies that opening a PDF
// with an invalid operator throws an error when lenient parsing is disabled.
// This matches C# InvalidOperatorTests.InvalidOperatorThrowsExceptionIfNotUsingLenientParsing.
func TestInvalidOperatorThrowsExceptionIfNotUsingLenientParsing(t *testing.T) {
	path := getInvalidOperatorPath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	_, err = doc.GetPage(1)
	if err == nil {
		t.Fatal("expected error when getting page 1 with lenient parsing disabled, got nil")
	}
}

// TestInvalidOperatorDoesNotThrowExceptionIfUsingLenientParsing verifies that opening a PDF
// with an invalid operator succeeds and extracts text when lenient parsing is enabled.
// This matches C# InvalidOperatorTests.InvalidOperatorDoesNotThrowExceptionIfUsingLenientParsing.
func TestInvalidOperatorDoesNotThrowExceptionIfUsingLenientParsing(t *testing.T) {
	path := getInvalidOperatorPath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
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

	if !strings.Contains(text, "Text line 1") {
		t.Errorf("page text does not contain 'Text line 1'; got: %q", text)
	}

	if !strings.Contains(text, "Text line 2") {
		t.Errorf("page text does not contain 'Text line 2'; got: %q", text)
	}
}
