//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

func singlePageFormContentIText1Filename() string {
	return filepath.Join(integrationDocRoot, "Single Page Form Content - from itext 1_1.pdf")
}

// TestSinglePageFormContentIText1HasCorrectNumberOfPages verifies the document has 1 page.
func TestSinglePageFormContentIText1HasCorrectNumberOfPages(t *testing.T) {
	path := singlePageFormContentIText1Filename()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 1 {
		t.Errorf("NumberOfPages = %d, want 1", got)
	}
}

// TestSinglePageFormContentIText1HasCorrectPageSize verifies page 1 is A4.
func TestSinglePageFormContentIText1HasCorrectPageSize(t *testing.T) {
	path := singlePageFormContentIText1Filename()

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

	if page.Size() != content.PageSizeA4 {
		t.Errorf("page.Size() = %v, want PageSizeA4 (%v)", page.Size(), content.PageSizeA4)
	}
}

// TestSinglePageFormContentIText1ExtractsText verifies the extracted text matches expected content.
func TestSinglePageFormContentIText1ExtractsText(t *testing.T) {
	path := singlePageFormContentIText1Filename()

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

	expected := "Ich bin Marios test.$Anhang=/home/im/tmp/AGB.pdf$$Auftragsnummer=2004001452-001$$Mandant=1$$MailTo=mario@ops.co.at$$User=mario$Ende"

	if got := page.Text(); got != expected {
		t.Errorf("page.Text() = %q, want %q", got, expected)
	}
}

type letterPosition struct {
	index  int
	letter string
	x      float64
	y      float64
}

// TestSinglePageFormContentIText1TextHasCorrectPositions verifies specific letters are at expected coordinates.
func TestSinglePageFormContentIText1TextHasCorrectPositions(t *testing.T) {
	positionData := `   0|I|56.88|774.08
                                              1|c|60.72|774.08
                                              2|h|66.00|774.08
                                             20|$|56.88|744.80
                                             26|g|94.80|744.80
                                             49|$|56.88|730.16`

	expectedData := parsePositionData(positionData)

	path := singlePageFormContentIText1Filename()

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

	letters := page.Letters()

	for _, exp := range expectedData {
		if exp.index >= len(letters) {
			t.Errorf("letter index %d out of range, total letters = %d", exp.index, len(letters))
			continue
		}

		letter := letters[exp.index]

		if letter.Value != exp.letter {
			t.Errorf("letters[%d].Value = %q, want %q", exp.index, letter.Value, exp.letter)
		}

		loc := letter.Location()
		roundedX := math.Round(loc.X*100) / 100
		roundedY := math.Round(loc.Y*100) / 100

		if roundedX != exp.x {
			t.Errorf("letters[%d].Location.X = %.2f (raw: %g), want %.2f", exp.index, roundedX, loc.X, exp.x)
		}
		if roundedY != exp.y {
			t.Errorf("letters[%d].Location.Y = %.2f (raw: %g), want %.2f", exp.index, roundedY, loc.Y, exp.y)
		}
	}
}

func parsePositionData(positionData string) []letterPosition {
	lines := strings.Split(positionData, "\n")
	var result []letterPosition

	for _, line := range lines {
		stripped := strings.TrimSpace(line)
		if stripped == "" {
			continue
		}

		parts := strings.Split(stripped, "|")
		if len(parts) < 4 {
			continue
		}

		index, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}

		xVal, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
		if err != nil {
			continue
		}

		yVal, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
		if err != nil {
			continue
		}

		result = append(result, letterPosition{
			index:  index,
			letter: parts[1],
			x:      xVal,
			y:      yVal,
		})
	}

	return result
}
