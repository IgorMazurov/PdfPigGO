package word_extractor_test

import (
	"strings"
	"testing"

	dla "github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/content"
	pdfpig "github.com/uglytoad/pdfpig/go"
)

func TestReadWordsFromDataPdfPage3(t *testing.T) {
	docPath := dla.GetDocumentPath("data.pdf", true)

	opts := &content.ParsingOptions{
		UseLenientParsing: true,
	}

	doc, err := pdfpig.OpenFile(docPath, opts)
	if err != nil {
		t.Fatalf("OpenFile(%q): %v", docPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(3)
	if err != nil {
		t.Fatalf("GetPage(3): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}

	words := page.GetWords()

	texts := make([]string, len(words))
	for i, w := range words {
		texts[i] = w.Text
	}
	text := strings.Join(texts, " ")

	expected := "len supp dose 4.2 VC 0.5 11.5 VC 0.5 7.3 VC 0.5 5.8 VC 0.5 6.4 VC 0.5 10.0 VC 0.5 11.2 VC 0.5 11.2 VC 0.5 5.2 VC 0.5 7.0 VC 0.5" +
		" 16.5 VC 1.0 16.5 VC 1.0 15.2 VC 1.0 17.3 VC 1.0 22.5 VC 1.0 3"

	if text != expected {
		t.Errorf("expected %q, got %q", expected, text)
	}
}

func TestReadWordsFromOldGutnishPage1(t *testing.T) {
	docPath := dla.GetDocumentPath("Old Gutnish Internet Explorer.pdf", true)

	opts := &content.ParsingOptions{
		UseLenientParsing: true,
	}

	doc, err := pdfpig.OpenFile(docPath, opts)
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

	words := page.GetWords()

	texts := make([]string, len(words))
	for i, w := range words {
		texts[i] = w.Text
	}
	text := strings.Join(texts, " ")

	expectedPrefix := "Old Gutnish - Wikipedia Page 1 of 3 Old Gutnish Old Gutnish was the dialect of Old Norse that was spoken on the Baltic island of Gotland." +
		" It shows sufficient differences from the Old West Norse and Old East Norse dialects that it is considered to be a separate branch." +
		" Gutnish is still spoken in some parts of Gotland and on the adjoining island of Fårö."

	if !strings.HasPrefix(text, expectedPrefix) {
		t.Errorf("text does not start with expected prefix.\nexpected prefix: %q\ngot: %q", expectedPrefix, text)
	}
}

func TestReadWordsFromTikka1552Page8(t *testing.T) {
	docPath := dla.GetDocumentPath("TIKA-1552-0.pdf", true)

	opts := &content.ParsingOptions{
		UseLenientParsing: true,
	}

	doc, err := pdfpig.OpenFile(docPath, opts)
	if err != nil {
		t.Fatalf("OpenFile(%q): %v", docPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(8)
	if err != nil {
		t.Fatalf("GetPage(8): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}

	words := page.GetWords()

	texts := make([]string, len(words))
	for i, w := range words {
		texts[i] = w.Text
	}
	text := strings.Join(texts, " ")

	expectedPrefix := "2 THE BUDGET MESSAGE OF THE PRESIDENT Administration’s SelectUSA initiative to help draw businesses and investment from around the world to our shores." +
		" If we want to make the best products, we also have to invest in the best ideas. That is why the Budget maintains a world-class commitment to science and research," +
		" targeting resources to those areas most likely to contribute directly to the creation of transformational technologies that can create the businesses and jobs of the future."

	if !strings.HasPrefix(text, expectedPrefix) {
		t.Errorf("text does not start with expected prefix.\nexpected prefix: %q\ngot: %q", expectedPrefix, text)
	}
}
