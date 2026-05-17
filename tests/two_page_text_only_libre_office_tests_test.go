//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func getTwoPageTextOnlyLibreOfficePath() string {
	return testutil.GetDocumentPath("Two Page Text Only - from libre office.pdf", true)
}

func TestTwoPageTextOnlyLibreOfficeHasCorrectNumberOfPages(t *testing.T) {
	filePath := getTwoPageTextOnlyLibreOfficePath()

	doc, err := pdfpig.OpenFile(filePath, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 2 {
		t.Errorf("NumberOfPages = %d, want 2", got)
	}
}

func TestTwoPageTextOnlyLibreOfficeHasCorrectPageSize(t *testing.T) {
	filePath := getTwoPageTextOnlyLibreOfficePath()

	doc, err := pdfpig.OpenFile(filePath, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	for pageNum := 1; pageNum <= 2; pageNum++ {
		pageAny, err := doc.GetPage(pageNum)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", pageNum, err)
		}
		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("page %d: expected *content.Page from GetPage", pageNum)
		}

		if got := page.Size(); got != content.PageSizeA4 {
			t.Errorf("page %d Size = %v, want PageSizeA4 (%v)", pageNum, got, content.PageSizeA4)
		}
	}
}

func TestTwoPageTextOnlyLibreOfficePagesStartWithCorrectText(t *testing.T) {
	filePath := getTwoPageTextOnlyLibreOfficePath()

	doc, err := pdfpig.OpenFile(filePath, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	tests := []struct {
		pageNum    int
		startsWith string
	}{
		{1, "Apache License"},
		{2, "2. Grant of Copyright"},
	}

	for _, tc := range tests {
		t.Run(string(rune(tc.pageNum)), func(t *testing.T) {
			pageAny, err := doc.GetPage(tc.pageNum)
			if err != nil {
				t.Fatalf("GetPage(%d): %v", tc.pageNum, err)
			}
			page, ok := pageAny.(*content.Page)
			if !ok {
				t.Fatal("expected *content.Page from GetPage")
			}

			text := page.Text()
			if !strings.HasPrefix(text, tc.startsWith) {
				t.Errorf("page %d Text does not start with %q; got prefix %q", tc.pageNum, tc.startsWith, text[:min(len(text), len(tc.startsWith)+20)])
			}
		})
	}
}
