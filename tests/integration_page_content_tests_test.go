//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func TestDetectPageContents(t *testing.T) {
	path := testutil.GetDocumentPath("Various Content Types", true)

	doc, err := pdfpig.OpenFile(path, &pdfpig.LenientParsingOff)
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

	foundStroke := false
	for _, l := range letters {
		if l.RenderingMode == core.StrokeText {
			foundStroke = true
			break
		}
	}
	if !foundStroke {
		t.Error("expected at least one letter with RenderingMode Stroke (for REGULAR TEXT)")
	}

	foundNeither := false
	for _, l := range letters {
		if l.RenderingMode == core.NeitherFillNorStroke {
			foundNeither = true
			break
		}
	}
	if !foundNeither {
		t.Error("expected at least one letter with RenderingMode NeitherFillNorStroke (for INVISIBLE TEXT)")
	}

	images := page.GetImages()
	if len(images) == 0 {
		t.Error("expected non-empty images from page content")
	}

	paths := page.Paths()
	if len(paths) == 0 {
		t.Error("expected non-empty paths from page content")
	}
}
