//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/testutil"
)

// TestReplacesObjectsFunc verifies that ReplaceIndirectObjectWithFunc can replace
// the page contents stream with an empty one, resulting in zero extracted letters.
// This matches C# AdvancedPdfDocumentAccessTests.ReplacesObjectsFunc.
func TestReplacesObjectsFunc(t *testing.T) {
	path := testutil.GetDocumentPath("Single Page Simple - from inkscape.pdf", true)

	opts := &content.ParsingOptions{}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	// Access the AdvancedPdfDocumentAccess through type assertion.
	adv, ok := doc.Advanced.(*document.AdvancedPdfDocumentAccess)
	if !ok || adv == nil {
		t.Skip("AdvancedPdfDocumentAccess not wired into PdfDocument; skipping test")
	}

	pageNode := doc.Structure.Catalog().Pages().GetPageNode(1)
	contentsToken, ok := pageNode.NodeDictionary.TryGet(tokens.Contents)
	if !ok {
		t.Fatal("page dictionary has no /Contents entry")
	}

	irToken, ok := contentsToken.(*tokens.IndirectReferenceToken)
	if !ok {
		t.Fatalf("/Contents is not an IndirectReferenceToken, got %T", contentsToken)
	}

	err = adv.ReplaceIndirectObjectWithFunc(irToken.Data(), func(existing tokens.Token) tokens.Token {
		dictData := map[*tokens.NameToken]tokens.Token{
			tokens.Length: tokens.Zero,
		}
		dict, dictErr := tokens.NewDictionary(dictData)
		if dictErr != nil {
			t.Errorf("NewDictionary failed: %v", dictErr)
			return existing
		}
		stream, streamErr := tokens.NewStreamToken(dict, []byte{})
		if streamErr != nil {
			t.Errorf("NewStreamToken failed: %v", streamErr)
			return existing
		}
		return stream
	})
	if err != nil {
		t.Fatalf("ReplaceIndirectObjectWithFunc: %v", err)
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage returned %T, expected *content.Page", pageAny)
	}

	letters := page.Letters()
	if len(letters) != 0 {
		t.Errorf("expected 0 letters after replacing contents with empty stream, got %d", len(letters))
	}
}

// TestReplacesObjects verifies that ReplaceIndirectObject can replace the page
// contents stream with an empty one, resulting in zero extracted letters.
// This matches C# AdvancedPdfDocumentAccessTests.ReplacesObjects.
func TestReplacesObjects(t *testing.T) {
	path := testutil.GetDocumentPath("Single Page Simple - from inkscape.pdf", true)

	opts := &content.ParsingOptions{}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	// Access the AdvancedPdfDocumentAccess through type assertion.
	adv, ok := doc.Advanced.(*document.AdvancedPdfDocumentAccess)
	if !ok || adv == nil {
		t.Skip("AdvancedPdfDocumentAccess not wired into PdfDocument; skipping test")
	}

	dictData := map[*tokens.NameToken]tokens.Token{
		tokens.Length: tokens.Zero,
	}
	dict, dictErr := tokens.NewDictionary(dictData)
	if dictErr != nil {
		t.Fatalf("NewDictionary failed: %v", dictErr)
	}
	replacement, streamErr := tokens.NewStreamToken(dict, []byte{})
	if streamErr != nil {
		t.Fatalf("NewStreamToken failed: %v", streamErr)
	}

	pageNode := doc.Structure.Catalog().Pages().GetPageNode(1)
	contentsToken, ok := pageNode.NodeDictionary.TryGet(tokens.Contents)
	if !ok {
		t.Fatal("page dictionary has no /Contents entry")
	}

	irToken, ok := contentsToken.(*tokens.IndirectReferenceToken)
	if !ok {
		t.Fatalf("/Contents is not an IndirectReferenceToken, got %T", contentsToken)
	}

	err = adv.ReplaceIndirectObject(irToken.Data(), replacement)
	if err != nil {
		t.Fatalf("ReplaceIndirectObject: %v", err)
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage returned %T, expected *content.Page", pageAny)
	}

	letters := page.Letters()
	if len(letters) != 0 {
		t.Errorf("expected 0 letters after replacing contents with empty stream, got %d", len(letters))
	}
}
