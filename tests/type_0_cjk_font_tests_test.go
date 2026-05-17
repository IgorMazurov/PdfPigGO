//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

const type0CjkFontDocName = "Type0_CJK_Font.pdf"

func getType0CjkFontPath() string {
	return filepath.Join(integrationDocRoot, type0CjkFontDocName)
}

// TestType0CjkFontHasCorrectNumberOfPages verifies the document has 95 pages.
func TestType0CjkFontHasCorrectNumberOfPages(t *testing.T) {
	path := getType0CjkFontPath()

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 95 {
		t.Errorf("NumberOfPages = %d, want 95", got)
	}
}

// TestType0CjkFontHasCorrectChineseCharacters verifies the extracted page text
// contains expected Chinese character sequences.
func TestType0CjkFontHasCorrectChineseCharacters(t *testing.T) {
	path := getType0CjkFontPath()

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

	expectedSubstrings := []string{
		"中航动力控制股份有限公司",
		"年半年度报告",
		"中航动力控制股份有限公司董事会",
		"2010年8月17日",
	}

	for _, substr := range expectedSubstrings {
		if !strings.Contains(pageText, substr) {
			t.Errorf("page text does not contain %q; got length=%d", substr, len(pageText))
		}
	}
}
