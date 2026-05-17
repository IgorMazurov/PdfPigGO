//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func TestZapfDingbatsTrueTypeSimpleFont1(t *testing.T) {
	path := testutil.GetDocumentPath("capas", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(18)
	if err != nil {
		t.Fatalf("GetPage(18): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(18): expected *content.Page, got %T", pageAny)
	}

	letters := page.Letters()
	for _, l := range letters {
		if strings.ContainsRune(l.Value, ' ') {
			return
		}
	}

	t.Error("expected at least one letter containing space character from ZapfDingbats")
}

func TestZapfDingbatsType1Standard14Font1(t *testing.T) {
	path := testutil.GetDocumentPath("TIKA-469-0", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(2): expected *content.Page, got %T", pageAny)
	}

	letters := page.Letters()
	for _, l := range letters {
		if strings.ContainsRune(l.Value, '●') {
			return
		}
	}

	t.Error("expected at least one letter containing black circle character (\\u25cf)")
}

func TestZapfDingbatsType1Standard14Font2(t *testing.T) {
	path := testutil.GetDocumentPath("MOZILLA-LINK-5251-1", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
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
	values := make([]string, len(letters))
	for i, l := range letters {
		values[i] = l.Value
	}

	expectedChars := []rune{'✁', '✂', '✄', '☎', '✆', '✇'}
	for _, ch := range expectedChars {
		found := false
		for _, v := range values {
			if strings.ContainsRune(v, ch) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected letter containing character %q (\\u%04x), not found in any letter value", string(ch), ch)
		}
	}
}

func TestZapfDingbatsType1FontSimple1(t *testing.T) {
	path := testutil.GetDocumentPath("MOZILLA-2775-1", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(11)
	if err != nil {
		t.Fatalf("GetPage(11): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(11): expected *content.Page, got %T", pageAny)
	}

	letters := page.Letters()
	for _, l := range letters {
		if strings.ContainsRune(l.Value, '●') {
			return
		}
	}

	t.Error("expected at least one letter containing black circle character (\\u25cf)")
}

func TestZapfDingbatsType1FontSimple2(t *testing.T) {
	path := testutil.GetDocumentPath("PDFBOX-492-4.jar-8", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
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
	for _, l := range letters {
		if strings.ContainsRune(l.Value, '\u25a0') {
			return
		}
	}

	t.Error("expected at least one letter containing black square character (\\u25a0)")
}
