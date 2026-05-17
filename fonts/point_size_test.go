package fonts_test

import (
	"testing"

	dla "github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/content"
	pdfpig "github.com/uglytoad/pdfpig/go"
)

func TestRotatedText(t *testing.T) {
	docPath := dla.GetDocumentPath("complex rotated", true)

	doc, err := pdfpig.OpenFile(docPath, nil)
	if err != nil {
		t.Fatalf("OpenFile(%q): %v", docPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}

	for _, letter := range page.Letters() {
		if letter.PointSize != 12 {
			t.Errorf("expected PointSize == 12, got %g for letter %q", letter.PointSize, letter.Value)
		}
	}
}
