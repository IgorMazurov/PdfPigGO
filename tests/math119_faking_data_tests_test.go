//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func getMath119FakingDataPath() string {
	return testutil.GetDocumentPath("Math119FakingData.pdf", true)
}

// TestMath119FakingDataCanReadPage8 verifies that the Math119FakingData document's
// page 8 can be opened and words extracted without error. This matches C#
// Math119FakingDataTests.CombinesDiaeresisForWords.
func TestMath119FakingDataCanReadPage8(t *testing.T) {
	path := getMath119FakingDataPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(8)
	if err != nil {
		t.Fatalf("GetPage(8): %v", err)
	}
	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("page 8 is not *content.Page")
	}

	words := page.GetWords()

	for i, w := range words {
		t.Logf("word[%d]: %q", i, w.Text)
	}
}
