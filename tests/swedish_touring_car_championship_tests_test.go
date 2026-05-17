//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/images/png"
)

func swedishTouringCarChampionshipFilename() string {
	return filepath.Join(integrationDocRoot, "2006_Swedish_Touring_Car_Championship.pdf")
}

// TestSwedishTouringCarChampionshipHasCorrectNumberOfPages verifies the document has 4 pages.
func TestSwedishTouringCarChampionshipHasCorrectNumberOfPages(t *testing.T) {
	path := swedishTouringCarChampionshipFilename()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 4 {
		t.Errorf("NumberOfPages = %d, want 4", got)
	}
}

// TestSwedishTouringCarChampionshipHasCorrectVersion verifies the PDF version is 1.4.
func TestSwedishTouringCarChampionshipHasCorrectVersion(t *testing.T) {
	path := swedishTouringCarChampionshipFilename()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.Version(); got != 1.4 {
		t.Errorf("Version = %.1f, want 1.4", got)
	}
}

// TestSwedishTouringCarChampionshipGetsFirstPageContent verifies page 1 text and A4 size.
func TestSwedishTouringCarChampionshipGetsFirstPageContent(t *testing.T) {
	path := swedishTouringCarChampionshipFilename()

	doc, err := pdfpig.OpenFile(path, nil)
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

	expectedSubstring := "A privateers championship named Caran Cup was created for drivers using cars constructed in 2003 or earlier"
	if !strings.Contains(page.Text(), expectedSubstring) {
		t.Errorf("page.Text() does not contain %q", expectedSubstring)
	}

	if got := page.Size(); got != content.PageSizeA4 {
		t.Errorf("page.Size() = %v, want PageSizeA4 (%v)", got, content.PageSizeA4)
	}
}

// TestSwedishTouringCarChampionshipGetsSwedishCharacters verifies Swedish characters on pages 2 and 3.
func TestSwedishTouringCarChampionshipGetsSwedishCharacters(t *testing.T) {
	path := swedishTouringCarChampionshipFilename()

	doc, err := pdfpig.OpenFile(path, nil)
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

	if !strings.Contains(page.Text(), "Vålerbanen") {
		t.Error("page 2 Text does not contain \"Vålerbanen\"")
	}

	pageAny, err = doc.GetPage(3)
	if err != nil {
		t.Fatalf("GetPage(3): %v", err)
	}

	page, ok = pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(3): expected *content.Page, got %T", pageAny)
	}

	if !strings.Contains(page.Text(), "Söderberg") {
		t.Error("page 3 Text does not contain \"Söderberg\"")
	}
}

// TestSwedishTouringCarChampionshipGetsHyperlinks verifies hyperlinks on page 1.
func TestSwedishTouringCarChampionshipGetsHyperlinks(t *testing.T) {
	path := swedishTouringCarChampionshipFilename()

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

	links := page.GetHyperlinks()

	if got := len(links); got != 4 {
		t.Fatalf("len(links) = %d, want 4", got)
	}

	pageLink := links[0]

	if got := pageLink.Text; got != "Swedish Touring Car Championship" {
		t.Errorf("links[0].Text = %q, want %q", got, "Swedish Touring Car Championship")
	}

	if got := pageLink.Uri; got != "https://en.wikipedia.org/wiki/Swedish_Touring_Car_Championship" {
		t.Errorf("links[0].Uri = %q, want %q", got, "https://en.wikipedia.org/wiki/Swedish_Touring_Car_Championship")
	}

	year2005 := links[1]

	if got := year2005.Text; got != "2005" {
		t.Errorf("links[1].Text = %q, want %q", got, "2005")
	}

	if got := year2005.Uri; got != "https://en.wikipedia.org/wiki/2005_Swedish_Touring_Car_Championship" {
		t.Errorf("links[1].Uri = %q, want %q", got, "https://en.wikipedia.org/wiki/2005_Swedish_Touring_Car_Championship")
	}

	year2007 := links[2]

	if got := year2007.Text; got != "2007" {
		t.Errorf("links[2].Text = %q, want %q", got, "2007")
	}

	if got := year2007.Uri; got != "https://en.wikipedia.org/wiki/2007_Swedish_Touring_Car_Championship" {
		t.Errorf("links[2].Uri = %q, want %q", got, "https://en.wikipedia.org/wiki/2007_Swedish_Touring_Car_Championship")
	}

	fullLink := links[3]

	expectedFullText := "The 2006 Swedish Touring Car Championship season was the 11th Swedish Touring Car Championship (STCC) season. In total nine racing weekends at six different circuits were held; each"
	if got := fullLink.Text; got != expectedFullText {
		t.Errorf("links[3].Text = %q, want %q", got, expectedFullText)
	}

	if got := fullLink.Uri; got != "https://en.wikipedia.org/wiki/Swedish_Touring_Car_Championship" {
		t.Errorf("links[3].Uri = %q, want %q", got, "https://en.wikipedia.org/wiki/Swedish_Touring_Car_Championship")
	}
}

// TestSwedishTouringCarChampionshipGetsImagesAsPng verifies all images can be converted to PNG.
func TestSwedishTouringCarChampionshipGetsImagesAsPng(t *testing.T) {
	path := swedishTouringCarChampionshipFilename()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pagesAny, err := doc.GetPages()
	if err != nil {
		t.Fatalf("GetPages: %v", err)
	}

	for pi, pageAny := range pagesAny {
		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("page index %d: expected *content.Page, got %T", pi, pageAny)
		}

		images := page.GetImages()

		for ii, image := range images {
			if _, ok := image.TryGetBytes(); !ok {
				continue
			}

			pngBytes, pngOk := image.TryGetPng()
			if !pngOk {
				t.Errorf("page %d, image %d: TryGetPng returned false", pi+1, ii)
				continue
			}

			pngActual, err := png.Open(bytes.NewReader(pngBytes), nil)
			if err != nil {
				t.Errorf("page %d, image %d: Png.Open error: %v", pi+1, ii, err)
				continue
			}

			if pngActual == nil {
				t.Errorf("page %d, image %d: Png.Open returned nil", pi+1, ii)
			}
		}
	}
}
