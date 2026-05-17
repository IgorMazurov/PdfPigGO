//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/testutil"
)

// TestWindows1252Encoding verifies that Windows-1252 encoded characters are correctly
// decoded when extracting text from a PDF page. Matches C# EncodingsTests.Windows1252Encoding.
func TestWindows1252Encoding(t *testing.T) {
	path := testutil.GetDocumentPath("GHOSTSCRIPT-698363-0.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page := pageAny.(*content.Page)
	letters := page.Letters()

	var actual strings.Builder
	for _, l := range letters {
		actual.WriteString(l.Value)
	}

	expected := "ҘҹЧѥЧКጹѝঐܮ̂ҥ҇ҁӃ࿋\u0c0dҀғҊ˺෨ཌආр෨ཌ̂ҘҹЧѥЧКጹѝঐܮ̂ҥ҇ҁӃ࿋\u0c0dҀғҊ˺෨ཌආр෨ཌ̂ݰႺംࢥ༢࣭\u089aѽ̔ҫһҐ̔ݰႺംࢥ༢࣭\u089aѽ̔ҫһҐ̔"

	if actual.String() != expected {
		t.Errorf("text mismatch:\nexpected: %q\nactual:   %q", expected, actual.String())
	}
}

// TestIssue688 verifies that missing characters in a font with indexed DeviceRGB and
// JPXDecode are still correctly extracted. Matches C# EncodingsTests.Issue688.
func TestIssue688(t *testing.T) {
	path := testutil.GetDocumentPath("Indexed-DeviceRGB-JPXDecode-0-.0.0.255.0.-Font-F1_1_missing_char_255-1.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page := pageAny.(*content.Page)
	letters := page.Letters()

	if len(letters) == 0 {
		t.Fatal("expected non-empty letters")
	}

	rect := core.NewPdfRectangleFloat(207, 158, 229, 168.5)

	var missingChars []*content.Letter
	for _, l := range letters {
		bb := l.BoundingBox
		if rectContainsRect(rect, bb) {
			missingChars = append(missingChars, l)
		}
	}

	if len(missingChars) == 0 {
		t.Fatal("expected non-empty missing chars within rectangle")
	}

	if len(missingChars) != 2 {
		t.Errorf("expected exactly 2 missing chars, got %d", len(missingChars))
	}

	if len(missingChars) >= 1 && missingChars[0].Value != "先" {
		t.Errorf("missingChars[0].Value: expected %q, got %q", "先", missingChars[0].Value)
	}

	if len(missingChars) >= 2 && missingChars[1].Value != "祖" {
		t.Errorf("missingChars[1].Value: expected %q, got %q", "祖", missingChars[1].Value)
	}
}

// rectContainsRect reports whether outer fully contains inner by checking all four corners.
func rectContainsRect(outer, inner core.PdfRectangle) bool {
	return outer.Contains(inner.BottomLeft, true) &&
		outer.Contains(inner.TopLeft, true) &&
		outer.Contains(inner.TopRight, true) &&
		outer.Contains(inner.BottomRight, true)
}
