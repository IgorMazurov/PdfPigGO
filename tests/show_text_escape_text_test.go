//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/writer"
)

func getShowTextEscapeFilename() string {
	return filepath.Join(integrationDocRoot, "ShowTextOpWithUnbalancedRoundBrackets.pdf")
}

func TestPdfCopyShowTextOpUsesEscapedText(t *testing.T) {
	filePath := getShowTextEscapeFilename()

	sourceDocument, err := pdfpig.OpenFile(filePath, nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer sourceDocument.Close()

	pdfBuilder := writer.NewPdfDocumentBuilder()

	pageNumber := 1
	sourcePageAny, err := sourceDocument.GetPage(pageNumber)
	if err != nil {
		t.Fatalf("GetPage(%d): %v", pageNumber, err)
	}
	sourcePage, ok := sourcePageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page from GetPage")
	}

	pageBuilder, err := pdfBuilder.AddPageWithSize(sourcePage.Width(), sourcePage.Height())
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	if _, err = pageBuilder.CopyFrom(sourcePage); err != nil {
		t.Fatalf("CopyFrom: %v", err)
	}

	pdfBytes, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	document, err := pdfpig.Open(pdfBytes, nil)
	if err != nil {
		t.Fatalf("Open from bytes: %v", err)
	}
	defer document.Close()

	pageAny, err := document.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page from GetPage")
	}

	words := page.GetWords()

	found := false
	for _, w := range words {
		if strings.Contains(w.Text, "wander") {
			found = true
			break
		}
	}

	if !found {
		var allText string
		for i, w := range words {
			if i > 0 {
				allText += " "
			}
			allText += w.Text
		}
		t.Errorf("expected word containing 'wander' in copied PDF; got words: %s", allText)
	}
}
