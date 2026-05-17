//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func getLittlePigInArmenianPath() string {
	return testutil.GetDocumentPath("little-pig-in-armenian.pdf", true)
}

// TestLittlePigInArmenianCanReadTextCorrectly verifies that Armenian text is extracted correctly.
// This matches C# LittlePigInArmenianTests.CanReadTextCorrectly.
func TestLittlePigInArmenianCanReadTextCorrectly(t *testing.T) {
	path := getLittlePigInArmenianPath()

	doc, err := pdfpig.OpenFile(path, nil)
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
		t.Fatal("page 1 is not *content.Page")
	}

	words := page.GetWords()

	textFromWords := strings.Join(func() []string {
		texts := make([]string, len(words))
		for i, w := range words {
			texts[i] = w.Text
		}
		return texts
	}(), " ")

	expected := "\u0583\u0578\u0584\u0580\u056B\u056F \u056D\u0578\u0566"
	if textFromWords != expected {
		t.Errorf("text = %q, want %q", textFromWords, expected)
	}
}
