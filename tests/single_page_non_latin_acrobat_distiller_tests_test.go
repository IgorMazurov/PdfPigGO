//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

const singlePageNonLatinAcrobatDistillerDocName = "Single Page Non Latin - from acrobat distiller.pdf"

func getSinglePageNonLatinPath() string {
	return filepath.Join(integrationDocRoot, singlePageNonLatinAcrobatDistillerDocName)
}

// TestSinglePageNonLatinHasCorrectNumberOfPages verifies the document has 1 page.
func TestSinglePageNonLatinHasCorrectNumberOfPages(t *testing.T) {
	path := getSinglePageNonLatinPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 1 {
		t.Errorf("NumberOfPages = %d, want 1", got)
	}
}

// TestSinglePageNonLatinHasCorrectPageSize verifies page size is Letter.
func TestSinglePageNonLatinHasCorrectPageSize(t *testing.T) {
	path := getSinglePageNonLatinPath()

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
		t.Fatal("page is not *content.Page")
	}

	if size := page.Size(); size != content.PageSizeLetter {
		t.Errorf("Size = %v, want PageSizeLetter", size)
	}
}

// TestSinglePageNonLatinGetsCorrectPageTextIgnoringHiddenCharacters verifies extracted text.
func TestSinglePageNonLatinGetsCorrectPageTextIgnoringHiddenCharacters(t *testing.T) {
	path := getSinglePageNonLatinPath()

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
		t.Fatal("page is not *content.Page")
	}

	letters := page.Letters()
	var text strings.Builder
	for _, l := range letters {
		text.WriteString(l.Value)
	}

	expected := "Hello ﺪﻤﺤﻣ World. "
	if got := text.String(); got != expected {
		t.Errorf("Text = %q, want %q", got, expected)
	}
}

// TestSinglePageNonLatinLetterPositionsAreCorrectPdfBox compares letter positions against PdfBox data.
func TestSinglePageNonLatinLetterPositionsAreCorrectPdfBox(t *testing.T) {
	path := getSinglePageNonLatinPath()

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
		t.Fatal("page is not *content.Page")
	}

	letters := page.Letters()
	positions := getPdfBoxPositionDataNonLatin(t)

	comparer := testutil.NewDoubleComparer(3)

	index := 0
	for _, letter := range letters {
		if index >= len(positions) {
			break
		}

		myLetter := letter.Value
		theirLetter := positions[index].Text

		if myLetter == " " && theirLetter != " " {
			continue
		}

		if myLetter != theirLetter {
			t.Errorf("Letter[%d] Text = %q, want %q", index, myLetter, theirLetter)
		}

		if !comparer.Equals(positions[index].X, letter.Location().X) {
			t.Errorf("Letter[%d] X = %g, want %g (tolerance 3)", index, letter.Location().X, positions[index].X)
		}

		if !comparer.Equals(positions[index].Width, letter.Width) {
			t.Errorf("Letter[%d] Width = %g, want %g (tolerance 3)", index, letter.Width, positions[index].Width)
		}

		index++
	}
}

// TestSinglePageNonLatinLetterPositionsAreCorrectXfinium compares letter positions against Xfinium data.
func TestSinglePageNonLatinLetterPositionsAreCorrectXfinium(t *testing.T) {
	path := getSinglePageNonLatinPath()

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
		t.Fatal("page is not *content.Page")
	}

	letters := page.Letters()
	positions := getXfiniumPositionDataNonLatin(t)

	comparer := testutil.NewDoubleComparer(1)

	index := 0
	for _, letter := range letters {
		if index >= len(positions) {
			break
		}

		myLetter := letter.Value
		theirLetter := positions[index].Text

		if myLetter == " " && theirLetter != " " {
			continue
		}

		if myLetter != theirLetter {
			t.Errorf("Letter[%d] Text = %q, want %q", index, myLetter, theirLetter)
		}

		if !comparer.Equals(positions[index].X, letter.Location().X) {
			t.Errorf("Letter[%d] X = %g, want %g (tolerance 1)", index, letter.Location().X, positions[index].X)
		}

		if !comparer.Equals(positions[index].Width, letter.Width) {
			t.Errorf("Letter[%d] Width = %g, want %g (tolerance 1)", index, letter.Width, positions[index].Width)
		}

		index++
	}
}

func getPdfBoxPositionDataNonLatin(t *testing.T) []*testutil.AssertablePositionData {
	t.Helper()

	const data = "90\t90.65997\t14.42556\tH\t19\tFFJICI+TimesNewRomanPSMT\n" +
		"104.4395\t90.65997\t8.871117\te\t19\tFFJICI+TimesNewRomanPSMT\n" +
		"113.3247\t90.65997\t5.554443\tl\t19\tFFJICI+TimesNewRomanPSMT\n" +
		"118.8931\t90.65997\t5.554443\tl\t19\tFFJICI+TimesNewRomanPSMT\n" +
		"124.4615\t90.65997\t9.989998\to\t19\tFFJICI+TimesNewRomanPSMT\n" +
		"139.4505\t90.65997\t6.733261\tﺪ\t19\tFFJIAH+TimesNewRomanPSMT\n" +
		"146.1778\t90.65997\t7.872116\tﻤ\t19\tFFJIAH+TimesNewRomanPSMT\n" +
		"154.0439\t90.65997\t10.5894\tﺤ\t19\tFFJIAH+TimesNewRomanPSMT\n" +
		"164.6273\t90.65997\t7.872116\tﻣ\t19\tFFJIAH+TimesNewRomanPSMT\n" +
		"177.4964\t90.65997\t18.86111\tW\t19\tFFJICI+TimesNewRomanPSMT\n" +
		"196.3575\t90.65997\t9.990005\to\t19\tFFJICI+TimesNewRomanPSMT\n" +
		"206.4275\t90.65997\t6.653336\tr\t19\tFFJICI+TimesNewRomanPSMT\n" +
		"213.0808\t90.65997\t5.554443\tl\t19\tFFJICI+TimesNewRomanPSMT\n" +
		"218.6352\t90.65997\t9.990005\td\t19\tFFJICI+TimesNewRomanPSMT\n" +
		"228.6252\t90.65997\t4.994995\t.\t19\tFFJICI+TimesNewRomanPSMT"

	lines := strings.Split(data, "\n")
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

func getXfiniumPositionDataNonLatin(t *testing.T) []*testutil.AssertablePositionData {
	t.Helper()

	const data = "90\t90.66\t14.439546\tH\t19\tFFJICI+TimesNewRomanPSMT\t17.802\n" +
		"104.4395\t90.66\t8.885106\te\t19\tFFJICI+TimesNewRomanPSMT\t17.80218\n" +
		"113.3247\t90.66\t5.568426\tl\t19\tFFJICI+TimesNewRomanPSMT\t17.80218\n" +
		"118.8931\t90.66\t5.568426\tl\t19\tFFJICI+TimesNewRomanPSMT\t17.80218\n" +
		"124.4615\t90.66\t10.003986\to\t19\tFFJICI+TimesNewRomanPSMT\t17.80218\n" +
		"139.4505\t90.66\t6.727266\tﺪ\t19\tFFJIAH+TimesNewRomanPSMT\t17.80218\n" +
		"146.1778\t90.66\t7.866126\tﻤ\t19\tFFJIAH+TimesNewRomanPSMT\t17.80218\n" +
		"154.0439\t90.66\t10.583406\tﺤ\t19\tFFJIAH+TimesNewRomanPSMT\t17.80218\n" +
		"164.6273\t90.66\t7.866126\tﻣ\t19\tFFJIAH+TimesNewRomanPSMT\t17.80218\n" +
		"177.4964\t90.66\t18.86112\tW\t19\tFFJICI+TimesNewRomanPSMT\t17.80218"

	lines := strings.Split(data, "\n")
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
