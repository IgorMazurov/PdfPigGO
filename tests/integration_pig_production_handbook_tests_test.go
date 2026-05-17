//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func pigProductionHandbookPath() string {
	return testutil.GetDocumentPath("Pig Production Handbook.pdf", true)
}

func TestPigProductionHandbookCanReadContent(t *testing.T) {
	path := pigProductionHandbookPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page := pageAny.(*content.Page)

	if !strings.Contains(page.Text(), "For the small holders at village level") {
		t.Errorf("page text does not contain expected substring")
	}
}

func TestPigProductionHandbookLettersHaveCorrectColors(t *testing.T) {
	path := pigProductionHandbookPath()

	doc, err := pdfpig.OpenFile(path, nil)
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
	if len(letters) < 77 {
		t.Fatalf("expected at least 77 letters, got %d", len(letters))
	}

	pigAssertRGBNear(t, "letter[0] pinkish", letters[0].Color.ToRGBValues(), 1, 0.914, 0.765)
	pigAssertRGBNear(t, "letter[37] white", letters[37].Color.ToRGBValues(), 1, 1, 1)
	pigAssertRGBNear(t, "letter[76] blackish", letters[76].Color.ToRGBValues(), 0.137, 0.122, 0.125)
}

func TestPigProductionHandbookPage1HasCorrectWords(t *testing.T) {
	expected := []string{
		"European",
		"Comission",
		"Farmer's",
		"Hand",
		"Book",
		"on",
		"Pig",
		"Production",
		"(For",
		"the",
		"small",
		"holders",
		"at",
		"village",
		"level)",
		"GCP/NEP/065/EC",
		"Food",
		"and",
		"Agriculture",
		"Organization",
		"of",
		"the",
		"United",
		"Nations",
	}

	path := pigProductionHandbookPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page := pageAny.(*content.Page)

	words := page.GetWords()

	wordTexts := make([]string, 0, len(words))
	for _, w := range words {
		wordTexts = append(wordTexts, w.Text)
	}

	pigAssertStringSliceEqual(t, "page 1 word texts", expected, wordTexts)
}

func TestPigProductionHandbookPage4HasCorrectWords(t *testing.T) {
	const wordsPage4 = `Disclaimer
The designations employed end the presentation of the material in this information
product do not imply the expression of any opinion whatsoever on the part of the
Food and Agriculture Organization of the United Nations (FAO) concerning the
legal or development status of any country, territory, city or area of its authorities,
or concerning the delimitation of its frontiers or boundaries. The mention of
specific companies or products of manufacturers, whether or not these have been
patented, does not imply that these have been endorsed or recommended by FAO
in preference to others of similar nature that are not mentioned.
The views expressed in this publication are those of the author(s) and do not
necessarily reflects the views of FAO.
All rights reserved. Reproduction and dissemination of materials in this information
product for educational or other non-commercial purposes are authorized without
any prior written permission from the copyright holders provided the source is
fully acknowledged. Reproduction in this information product for resale or other
commercial purposes is prohibited without written permission of the copyright
holders. Applications for such permission should be addressed to: Chief, Electronic
Publishing Policy and Support Branch Communication Division, FAO, Viale delle
Terme di Caracalla, 00153 Rome, Italy or by e-mail to: copyright@fao.org
FAO 2009
design&print: wps, eMail: printnepal@gmail.com`

	expected := strings.Fields(wordsPage4)

	path := pigProductionHandbookPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(4)
	if err != nil {
		t.Fatalf("GetPage(4): %v", err)
	}
	page := pageAny.(*content.Page)

	words := page.GetWords()

	wordTexts := make([]string, 0, len(words))
	for _, w := range words {
		wordTexts = append(wordTexts, w.Text)
	}

	pigAssertStringSliceEqual(t, "page 4 word texts", expected, wordTexts)
}

func TestPigProductionHandbookCanReadPage9(t *testing.T) {
	path := pigProductionHandbookPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(9)
	if err != nil {
		t.Fatalf("GetPage(9): %v", err)
	}
	page := pageAny.(*content.Page)

	expected := "BreedsNative breeds of pig can be found throughout the country. They are a small body size compared to other exotic and crosses pig types. There name varies from region to region, for example"

	if !strings.Contains(page.Text(), expected) {
		t.Errorf("page 9 text does not contain expected substring")
	}
}

func TestPigProductionHandbookHasCorrectNumberOfPages(t *testing.T) {
	path := pigProductionHandbookPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 86 {
		t.Errorf("expected 86 pages, got %d", got)
	}
}

func TestPigProductionHandbookLettersHaveCorrectPosition(t *testing.T) {
	path := pigProductionHandbookPath()

	doc, err := pdfpig.OpenFile(path, nil)
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
	positions := getPigProductionHandbookPositions(t)

	if len(letters) != len(positions) {
		t.Fatalf("expected %d letters, got %d", len(positions), len(letters))
	}

	for i := 0; i < len(letters); i++ {
		positions[i].AssertWithinTolerance(t, letters[i], false)
	}
}

func getPigProductionHandbookPositions(t *testing.T) []*testutil.AssertablePositionData {
	t.Helper()

	posPath := filepath.Join(testutil.IntegrationDocumentsRoot, "Pig Production Handbook.Page1.Positions.txt")

	data, err := os.ReadFile(posPath)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", posPath, err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	var result []*testutil.AssertablePositionData
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pd, err := testutil.ParseAssertablePositionData(line)
		if err != nil {
			t.Fatalf("ParseAssertablePositionData(%q): %v", line, err)
		}
		result = append(result, pd)
	}

	return result
}

func pigAssertRGBNear(t *testing.T, name string, rgb colors.RGBValues, rWant, gWant, bWant float64) {
	t.Helper()

	const tol = 0.005

	pigAssertFloatNear(t, name+" R", rgb.R, rWant, tol)
	pigAssertFloatNear(t, name+" G", rgb.G, gWant, tol)
	pigAssertFloatNear(t, name+" B", rgb.B, bWant, tol)
}

func pigAssertFloatNear(t *testing.T, name string, got, want, tolerance float64) {
	t.Helper()

	if pigAbs(got-want) > tolerance {
		t.Errorf("%s: expected %g, got %g (diff=%g)", name, want, got, pigAbs(got-want))
	}
}

func pigAssertStringSliceEqual(t *testing.T, name string, expected, actual []string) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("%s: expected %d items, got %d", name, len(expected), len(actual))
		return
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("%s[%d]: expected %q, got %q", name, i, expected[i], actual[i])
		}
	}
}

func pigAbs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
