//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

type boldItalicData struct {
	value    string
	isBold   bool
	isItalic bool
}

var dataBoldItalic = []boldItalicData{
	{"L", false, false},
	{"o", false, false},
	{"r", false, false},
	{"e", false, false},
	{"m", false, false},
	{" ", false, false},
	{"i", true, false},
	{"p", true, false},
	{"s", true, false},
	{"u", true, false},
	{"m", true, false},
	{" ", false, false},
	{"d", false, false},
	{"o", false, false},
	{"l", false, false},
	{"o", false, false},
	{"r", false, false},
	{" ", false, false},
	{"s", false, true},
	{"i", false, true},
	{"t", false, true},
	{" ", false, false},
	{"a", true, true},
	{"m", true, true},
	{"e", true, true},
	{"t", true, true},
	{",", false, false},
	{" ", false, false},
	{"c", true, true},
	{"o", true, true},
	{"n", true, true},
	{"s", true, true},
	{"e", true, true},
	{"c", true, true},
	{"t", true, true},
	{"e", true, true},
	{"t", true, true},
	{"u", true, true},
	{"r", true, true},
	{" ", false, false},
	{"a", false, false},
	{"d", false, false},
	{"i", false, false},
	{"p", false, false},
	{"i", false, false},
	{"s", false, false},
	{"c", false, false},
	{"i", false, false},
	{"n", false, false},
	{"g", false, false},
	{" ", false, false},
	{"e", false, false},
	{"l", false, false},
	{"i", false, false},
	{"t", false, false},
	{".", false, false},
	{" ", false, false},
}

func resolveBoldItalicPath(t *testing.T) string {
	t.Helper()

	candidates := []string{
		filepath.Join(integrationDocRoot, "bold-italic.pdf"),
		filepath.FromSlash("testdata/integration/Documents/bold-italic.pdf"),
		filepath.FromSlash("../PdfPig-Source/src/UglyToad.PdfPig.Tests/Integration/Documents/bold-italic.pdf"),
	}

	for _, candidate := range candidates {
		if abs, err := filepath.Abs(candidate); err == nil {
			candidate = abs
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	t.Fatalf("bold-italic.pdf not found in any known location")
	return ""
}

func TestGetsCorrectBoldItalic(t *testing.T) {
	filename := resolveBoldItalicPath(t)

	doc, err := pdfpig.OpenFile(filename, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", filename, err)
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
	if len(letters) != len(dataBoldItalic) {
		t.Fatalf("expected %d letters, got %d", len(dataBoldItalic), len(letters))
	}

	for i, expected := range dataBoldItalic {
		letter := letters[i]

		if letter.Value != expected.value {
			t.Errorf("letter[%d].Value: expected %q, got %q", i, expected.value, letter.Value)
		}

		isBold := false
		isItalic := false
		if letter.FontDetails != nil {
			isBold = letter.FontDetails.IsBold
			isItalic = letter.FontDetails.IsItalic
		}

		if isBold != expected.isBold {
			t.Errorf("letter[%d].FontDetails.IsBold: expected %v, got %v (value=%q)", i, expected.isBold, isBold, letter.Value)
		}

		if isItalic != expected.isItalic {
			t.Errorf("letter[%d].FontDetails.IsItalic: expected %v, got %v (value=%q)", i, expected.isItalic, isItalic, letter.Value)
		}
	}
}
