//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/fonts/systemfonts"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func getClipPathLetterFilterTest1() string {
	return testutil.GetDocumentPath("ClipPathLetterFilter-Test1.pdf", true)
}

func getClipPathLetterFilterTest2() string {
	return testutil.GetDocumentPath("ClipPathLetterFilter-Test2.pdf", true)
}

// TestCanFilterClippedLetters verifies that enabling ClipPaths reduces the letter count.
// Skipped on platforms where TimesNewRomanPSMT is not available (matching C# behavior).
func TestCanFilterClippedLetters(t *testing.T) {
	if systemfonts.Instance.GetTrueTypeFont("TimesNewRomanPSMT") == nil {
		t.Skip("Skipped because the font TimesNewRomanPSMT could not be found.")
	}
	if systemfonts.Instance.GetTrueTypeFont("TimesNewRomanPS-BoldMT") == nil {
		t.Skip("Skipped because the font TimesNewRomanPS-BoldMT could not be found.")
	}
	if systemfonts.Instance.GetTrueTypeFont("TimesNewRomanPS-ItalicMT") == nil {
		t.Skip("Skipped because the font TimesNewRomanPS-ItalicMT could not be found.")
	}

	path := getClipPathLetterFilterTest1()

	docFiltered, err := pdfpig.OpenFile(path, &content.ParsingOptions{ClipPaths: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(ClipPaths=true): %v", err)
	}
	defer docFiltered.Close()

	docAll, err := pdfpig.OpenFile(path, &content.ParsingOptions{ClipPaths: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(ClipPaths=false): %v", err)
	}
	defer docAll.Close()

	filteredPageAny, err := docFiltered.GetPage(5)
	if err != nil {
		t.Fatalf("GetPage(5) filtered: %v", err)
	}
	filteredPage, ok := filteredPageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(5): expected *content.Page, got %T", filteredPageAny)
	}

	allPageAny, err := docAll.GetPage(5)
	if err != nil {
		t.Fatalf("GetPage(5) all: %v", err)
	}
	allPage, ok := allPageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(5): expected *content.Page, got %T", allPageAny)
	}

	filteredLetterCount := len(filteredPage.Letters())
	allLetterCount := len(allPage.Letters())

	if filteredLetterCount >= allLetterCount {
		t.Errorf("expected filtered letter count (%d) to be less than non-filtered (%d)", filteredLetterCount, allLetterCount)
	}
}

// TestCanFilterClippedLettersCheckBleedInSpecificWord verifies that clip path filtering
// prevents letters from hidden columns being merged into words.
func TestCanFilterClippedLettersCheckBleedInSpecificWord(t *testing.T) {
	path := getClipPathLetterFilterTest2()

	docFiltered, err := pdfpig.OpenFile(path, &content.ParsingOptions{ClipPaths: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(ClipPaths=true): %v", err)
	}
	defer docFiltered.Close()

	docAll, err := pdfpig.OpenFile(path, &content.ParsingOptions{ClipPaths: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(ClipPaths=false): %v", err)
	}
	defer docAll.Close()

	allPageAny, err := docAll.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1) all: %v", err)
	}
	allPage, ok := allPageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", allPageAny)
	}

	filteredPageAny, err := docFiltered.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1) filtered: %v", err)
	}
	filteredPage, ok := filteredPageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", filteredPageAny)
	}

	allWords := allPage.GetWords()
	filteredWords := filteredPage.GetWords()

	wordToSearch := "ARISER"

	indexInAll := -1
	for i, w := range allWords {
		if w.Text == wordToSearch {
			indexInAll = i
			break
		}
	}
	if indexInAll == -1 || indexInAll+1 >= len(allWords) {
		t.Fatalf("word %q not found or no following word in non-filtered words", wordToSearch)
	}

	allWordAfterAriser := allWords[indexInAll+1].Text
	expectedAllBad := "MLIA0U01CP00O0I3N6G2"
	if allWordAfterAriser != expectedAllBad {
		t.Errorf("non-filtered word after ARISER = %q, want %q (letters from hidden columns merged)", allWordAfterAriser, expectedAllBad)
	}

	indexInFiltered := -1
	for i, w := range filteredWords {
		if w.Text == wordToSearch {
			indexInFiltered = i
			break
		}
	}
	if indexInFiltered == -1 || indexInFiltered+1 >= len(filteredWords) {
		t.Fatalf("word %q not found or no following word in filtered words", wordToSearch)
	}

	filteredWordAfterAriser := filteredWords[indexInFiltered+1].Text
	expectedFilteredGood := "ACOGUT"
	if filteredWordAfterAriser != expectedFilteredGood {
		t.Errorf("filtered word after ARISER = %q, want %q (clip path should prevent hidden column bleed)", filteredWordAfterAriser, expectedFilteredGood)
	}
}
