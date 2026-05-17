//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func getSinglePageHyperlinksOpenOfficeFilename() string {
	return testutil.GetDocumentPath("Single Page Hyperlinks - from open office.pdf", true)
}

func TestSinglePageHyperlinksOpenOfficeGetsCorrectText(t *testing.T) {
	filePath := getSinglePageHyperlinksOpenOfficeFilename()

	doc, err := pdfpig.OpenFile(filePath, &pdfpig.LenientParsingOff)
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
		t.Fatal("expected *content.Page from GetPage")
	}

	expectedText := "https://duckduckgo.com/ a link aboveGitHub"
	if got := page.Text(); got != expectedText {
		t.Errorf("page.Text() = %q, want %q", got, expectedText)
	}
}

func TestSinglePageHyperlinksOpenOfficeGetsHyperlinks(t *testing.T) {
	filePath := getSinglePageHyperlinksOpenOfficeFilename()

	doc, err := pdfpig.OpenFile(filePath, &pdfpig.LenientParsingOff)
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
		t.Fatal("expected *content.Page from GetPage")
	}

	links := page.GetHyperlinks()

	if got := len(links); got != 2 {
		t.Fatalf("len(links) = %d, want 2", got)
	}

	ddg := links[0]

	if got := ddg.Text; got != "https://duckduckgo.com/" {
		t.Errorf("ddg.Text = %q, want %q", got, "https://duckduckgo.com/")
	}

	if got := ddg.Uri; got != "https://duckduckgo.com/" {
		t.Errorf("ddg.Uri = %q, want %q", got, "https://duckduckgo.com/")
	}

	expectedLetterCount := len("https://duckduckgo.com/ ")
	if got := len(ddg.Letters); got != expectedLetterCount {
		t.Errorf("len(ddg.Letters) = %d, want %d", got, expectedLetterCount)
	}

	if ddg.Annotation == nil {
		t.Error("ddg.Annotation is nil, want non-nil")
	}

	github := links[1]

	if got := github.Text; got != "GitHub" {
		t.Errorf("github.Text = %q, want %q", got, "GitHub")
	}

	if got := github.Uri; got != "https://github.com/" {
		t.Errorf("github.Uri = %q, want %q", got, "https://github.com/")
	}

	if got := len(github.Letters); got != 6 {
		t.Errorf("len(github.Letters) = %d, want 6", got)
	}

	if github.Annotation == nil {
		t.Error("github.Annotation is nil, want non-nil")
	}
}
