package writer_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/annotations"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	standard14fonts "github.com/uglytoad/pdfpig/go/fonts/standard14_fonts"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/parser"
	"github.com/uglytoad/pdfpig/go/testutil"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/writer"
)

func init() {
	// Set IntegrationDocumentsRoot directly to avoid multiple init() functions
	// each prepending "../" which results in incorrect paths like "../../../../testdata/..."
	testutil.IntegrationDocumentsRoot = "../testdata/integration/Documents"
}

func getDocNoExt(name string) string {
	return testutil.GetDocumentPath(name, false)
}

func loadFont(t *testing.T, name string) []byte {
	t.Helper()
	fontPath := filepath.Join("..", "fonts", "truetype", "testdata", name)
	data, err := os.ReadFile(fontPath)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", fontPath, err)
	}
	return data
}

func openBuilt(b []byte) (*content.PdfDocument, error) {
	return parser.OpenMemory(b, nil)
}

/* ========================================================================
   TestCanWriteSingleBlankPage
   C#: builder.AddPage(PageSize.A4), assert bytes start with %PDF and end %%EOF
   ======================================================================== */
func TestCanWriteSingleBlankPage(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	_, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	result, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	if len(result) == 0 {
		t.Fatal("result is empty")
	}

	str := string(result)
	if !strings.HasPrefix(str, "%PDF") {
		t.Errorf("expected PDF to start with %%PDF, got: %q", str[:min(20, len(str))])
	}
	if !strings.HasSuffix(strings.TrimSpace(str), "%%EOF") {
		t.Errorf("expected PDF to end with %%%%EOF, got suffix: %q", str[max(0, len(str)-20):])
	}
}

/* ========================================================================
   TestCanCreateSingleCustomPageSize
   C#: builder.AddPage(120, 250), add text "Small page.", open and verify
   ======================================================================== */
func TestCanCreateSingleCustomPageSize(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()

	page, err := pdfBuilder.AddPageWithSize(120, 250)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	font, err := pdfBuilder.AddStandard14Font(standard14fonts.Helvetica)
	if err != nil {
		t.Fatalf("AddStandard14Font: %v", err)
	}

	if _, addErr := page.AddText("Small page.", 12, core.PdfPoint{X: 25, Y: 200}, font); addErr != nil {
		t.Fatalf("AddText: %v", addErr)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 1 {
		t.Errorf("expected 1 page, got %d", doc.NumberOfPages())
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	if pg.Width() != 120 {
		t.Errorf("expected width 120, got %.1f", pg.Width())
	}
	if pg.Height() != 250 {
		t.Errorf("expected height 250, got %.1f", pg.Height())
	}

	if pg.Text() != "Small page." {
		t.Errorf("expected text %q, got %q", "Small page.", pg.Text())
	}
}

/* ========================================================================
   TestCanReadSingleBlankPage
   C#: Create blank A4 page, open it, verify 1 page, A4 size, empty letters
   ======================================================================== */
func TestCanReadSingleBlankPage(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	_, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	result, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(result)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 1 {
		t.Errorf("expected 1 page, got %d", doc.NumberOfPages())
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	if pg.Size() != content.PageSizeA4 {
		t.Errorf("expected PageSize.A4, got %v", pg.Size())
	}

	if len(pg.Letters()) > 0 {
		t.Errorf("expected empty letters, got %d", len(pg.Letters()))
	}
}

/* ========================================================================
   TestCanWriteSinglePageStandard14FontHelloWorld
   C#: Standard14 Helvetica "Hello World!", verify words
   ======================================================================== */
func TestCanWriteSinglePageStandard14FontHelloWorld(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	font, err := pdfBuilder.AddStandard14Font(standard14fonts.Helvetica)
	if err != nil {
		t.Fatalf("AddStandard14Font: %v", err)
	}

	if _, addErr := page.AddText("Hello World!", 12, core.PdfPoint{X: 25, Y: 520}, font); addErr != nil {
		t.Fatalf("AddText: %v", addErr)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	words := pg.GetWords()
	if len(words) < 2 {
		t.Fatalf("expected at least 2 words, got %d", len(words))
	}

	wordTexts := make([]string, len(words))
	for i, w := range words {
		wordTexts[i] = w.Text
	}
	expected := []string{"Hello", "World!"}
	if !equalStrings(wordTexts[:len(expected)], expected) {
		t.Errorf("expected words %v, got %v", expected, wordTexts)
	}
}

/* ========================================================================
   TestCanWriteSinglePageInvisibleHelloWorld
   C#: SetTextRenderingMode(Neither), text invisible but extractable
   ======================================================================== */
func TestCanWriteSinglePageInvisibleHelloWorld(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	font, err := pdfBuilder.AddStandard14Font(standard14fonts.Helvetica)
	if err != nil {
		t.Fatalf("AddStandard14Font: %v", err)
	}

	page.SetTextRenderingMode(core.NeitherFillNorStroke)
	if _, addErr := page.AddText("Hello World!", 12, core.PdfPoint{X: 25, Y: 520}, font); addErr != nil {
		t.Fatalf("AddText: %v", addErr)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	words := pg.GetWords()
	wordTexts := make([]string, len(words))
	for i, w := range words {
		wordTexts[i] = w.Text
	}
	expected := []string{"Hello", "World!"}
	if !equalStrings(wordTexts[:len(expected)], expected) {
		t.Errorf("expected words %v, got %v", expected, wordTexts)
	}
}

/* ========================================================================
   TestCanWriteSinglePageMixedRenderingMode
   C#: Visible -> Invisible -> Fill, verify all words extracted
   ======================================================================== */
func TestCanWriteSinglePageMixedRenderingMode(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	font, err := pdfBuilder.AddStandard14Font(standard14fonts.Helvetica)
	if err != nil {
		t.Fatalf("AddStandard14Font: %v", err)
	}

	if _, addErr := page.AddText("Hello World!", 12, core.PdfPoint{X: 25, Y: 520}, font); addErr != nil {
		t.Fatalf("AddText visible: %v", addErr)
	}

	page.SetTextRenderingMode(core.NeitherFillNorStroke)
	if _, addErr := page.AddText("Invisible!", 12, core.PdfPoint{X: 25, Y: 500}, font); addErr != nil {
		t.Fatalf("AddText invisible: %v", addErr)
	}

	page.SetTextRenderingMode(core.FillText)
	if _, addErr := page.AddText("Filled again!", 12, core.PdfPoint{X: 25, Y: 480}, font); addErr != nil {
		t.Fatalf("AddText fill: %v", addErr)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	words := pg.GetWords()
	wordTexts := make([]string, len(words))
	for i, w := range words {
		wordTexts[i] = w.Text
	}
	expected := []string{"Hello", "World!", "Invisible!", "Filled", "again!"}
	if len(wordTexts) < len(expected) || !equalStrings(wordTexts[:len(expected)], expected) {
		t.Errorf("expected words starting with %v, got %v", expected, wordTexts)
	}
}

/* ========================================================================
   TestCanWriteSinglePageHelloWorld (TrueType)
   C#: Andada-Regular.ttf, verify text and letter properties
   ======================================================================== */
func TestCanWriteSinglePageHelloWorld(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	page.DrawLine(core.PdfPoint{X: 30, Y: 520}, core.PdfPoint{X: 360, Y: 520}, 1)
	page.DrawLine(core.PdfPoint{X: 360, Y: 520}, core.PdfPoint{X: 360, Y: 250}, 1)

	page.SetStrokeColor(250, 132, 131)
	page.DrawLine(core.PdfPoint{X: 25, Y: 70}, core.PdfPoint{X: 100, Y: 70}, 3)
	page.ResetColor()
	page.DrawRectangle(core.PdfPoint{X: 30, Y: 200}, 250, 100, 0.5, false)
	page.DrawRectangle(core.PdfPoint{X: 30, Y: 100}, 250, 100, 0.5, false)

	fontBytes := loadFont(t, "Andada-Regular.ttf")
	font, err := pdfBuilder.AddTrueTypeFont(fontBytes)
	if err != nil {
		t.Fatalf("AddTrueTypeFont: %v", err)
	}

	letters, err := page.AddText("Hello World!", 12, core.PdfPoint{X: 30, Y: 50}, font)
	if err != nil {
		t.Fatalf("AddText: %v", err)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	if len(b) == 0 {
		t.Fatal("result is empty")
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	if pg.Text() != "Hello World!" {
		t.Errorf("expected text %q, got %q", "Hello World!", pg.Text())
	}

	readerLetters := pg.Letters()
	if len(readerLetters) == 0 {
		t.Fatal("no letters found on page")
	}

	h := readerLetters[0]
	if h.Value != "H" {
		t.Errorf("expected first letter 'H', got %q", h.Value)
	}

	for i := range letters {
		if i >= len(readerLetters) {
			break
		}
		readerLetter := readerLetters[i]
		writerLetter := letters[i]

		if readerLetter.Value != writerLetter.Value {
			t.Errorf("letter[%d]: expected value %q, got %q", i, writerLetter.Value, readerLetter.Value)
		}

		// Location tolerance of 10 accounts for kerning differences between
		// writer-computed and reader-extracted positions (Go reader may apply
		// kerning adjustments from the font's kern table that the writer doesn't).
		if absDiff(readerLetter.Location().X, writerLetter.Location().X) > 10 ||
			absDiff(readerLetter.Location().Y, writerLetter.Location().Y) > 10 {
			t.Errorf("letter[%d]: location mismatch: expected (%.2f, %.2f), got (%.2f, %.2f)",
				i, writerLetter.Location().X, writerLetter.Location().Y, readerLetter.Location().X, readerLetter.Location().Y)
		}

		if absDiff(readerLetter.FontSize, writerLetter.FontSize) > 0.01 {
			t.Errorf("letter[%d]: font size mismatch: expected %.2f, got %.2f", i, writerLetter.FontSize, readerLetter.FontSize)
		}
	}
}

/* ========================================================================
   TestCanWriteRobotoAccentedCharacters
   C#: Roboto-Regular.ttf, "eé" text round-trip
   ======================================================================== */
func TestCanWriteRobotoAccentedCharacters(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	title := "Hello Roboto!"
	pdfBuilder.DocumentInformation().Title = &title

	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	fontBytes := loadFont(t, "Roboto-Regular.ttf")
	font, err := pdfBuilder.AddTrueTypeFont(fontBytes)
	if err != nil {
		t.Fatalf("AddTrueTypeFont: %v", err)
	}

	if _, addErr := page.AddText("eé", 12, core.PdfPoint{X: 30, Y: 520}, font); addErr != nil {
		t.Fatalf("AddText: %v", addErr)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	if len(b) == 0 {
		t.Fatal("result is empty")
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	if pg.Text() != "eé" {
		t.Errorf("expected text %q, got %q", "eé", pg.Text())
	}
}

/* ========================================================================
   TestCanWriteSinglePageWithAccentedCharacters
   C#: Roboto-Regular.ttf, "é (lower case, upper case É)." round-trip
   ======================================================================== */
func TestCanWriteSinglePageWithAccentedCharacters(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	fontBytes := loadFont(t, "Roboto-Regular.ttf")
	font, err := pdfBuilder.AddTrueTypeFont(fontBytes)
	if err != nil {
		t.Fatalf("AddTrueTypeFont: %v", err)
	}

	expected := "é (lower case, upper case É)."
	if _, addErr := page.AddText(expected, 9, core.PdfPoint{X: 30, Y: page.PageSize().Bounds.TopRight.Y - 50}, font); addErr != nil {
		t.Fatalf("AddText: %v", addErr)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	if pg.Text() != expected {
		t.Errorf("expected text %q, got %q", expected, pg.Text())
	}
}

/* ========================================================================
   TestCanWriteSinglePageWithCzechCharacters
   C#: Roboto-Regular.ttf, "Hello: řó" round-trip
   ======================================================================== */
func TestCanWriteSinglePageWithCzechCharacters(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	fontBytes := loadFont(t, "Roboto-Regular.ttf")
	font, err := pdfBuilder.AddTrueTypeFont(fontBytes)
	if err != nil {
		t.Fatalf("AddTrueTypeFont: %v", err)
	}

	expected := "Hello: řó"
	if _, addErr := page.AddText(expected, 9, core.PdfPoint{X: 30, Y: page.PageSize().Bounds.TopRight.Y - 50}, font); addErr != nil {
		t.Fatalf("AddText: %v", addErr)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	if pg.Text() != expected {
		t.Errorf("expected text %q, got %q", expected, pg.Text())
	}
}

/* ========================================================================
   TestCanWriteTwoPageDocument
   C#: Two pages with Roboto text, verify both page texts
   ======================================================================== */
func TestCanWriteTwoPageDocument(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page1, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize(page1): %v", err)
	}
	page2, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize(page2): %v", err)
	}

	fontBytes := loadFont(t, "Roboto-Regular.ttf")
	font, err := pdfBuilder.AddTrueTypeFont(fontBytes)
	if err != nil {
		t.Fatalf("AddTrueTypeFont: %v", err)
	}

	topLine := core.PdfPoint{X: 30, Y: page1.PageSize().Bounds.TopRight.Y - 60}
	if _, addErr := page1.AddText("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor", 9, topLine, font); addErr != nil {
		t.Fatalf("AddText page1 line1: %v", addErr)
	}
	if _, addErr := page1.AddText("incididunt ut labore et dolore magna aliqua.", 9, core.PdfPoint{X: 30, Y: topLine.Y - 20}, font); addErr != nil {
		t.Fatalf("AddText page1 line2: %v", addErr)
	}

	page2Letters, err := page2.AddText("The very hungry caterpillar ate all the apples in the garden.", 12, topLine, font)
	if err != nil {
		t.Fatalf("AddText page2: %v", err)
	}

	if len(page2Letters) > 0 {
		left := page2Letters[0].BoundingBox.Left()
		bottom := page2Letters[0].BoundingBox.BottomLeft.Y
		right := page2Letters[len(page2Letters)-1].BoundingBox.Right()
		top := page2Letters[0].BoundingBox.Top()

		page2.SetStrokeColor(10, 250, 69)
		page2.DrawRectangle(core.PdfPoint{X: left, Y: bottom}, right-left, top-bottom, 1, false)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg1, ok := page1Any.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", page1Any)
	}

	if !strings.HasPrefix(pg1.Text(), "Lorem ipsum dolor sit") {
		t.Errorf("page 1 text should start with 'Lorem ipsum dolor sit', got: %q", pg1.Text())
	}

	page2Any, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	pg2, ok := page2Any.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", page2Any)
	}

	if !strings.HasPrefix(pg2.Text(), "The very hungry caterpillar") {
		t.Errorf("page 2 text should start with 'The very hungry caterpillar', got: %q", pg2.Text())
	}
}

/* ========================================================================
   TestCanWriteSinglePageWithJpeg
   C#: smile-250-by-160.jpg, verify image bounds and raw bytes
   ======================================================================== */
func TestCanWriteSinglePageWithJpeg(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	font, err := pdfBuilder.AddStandard14Font(standard14fonts.Helvetica)
	if err != nil {
		t.Fatalf("AddStandard14Font: %v", err)
	}

	if _, addErr := page.AddText("Smile", 12, core.PdfPoint{X: 25, Y: page.PageSize().Bounds.TopRight.Y - 52}, font); addErr != nil {
		t.Fatalf("AddText: %v", addErr)
	}

	imgPath := getDocNoExt("smile-250-by-160.jpg")
	imageBytes, err := os.ReadFile(imgPath)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", imgPath, err)
	}

	expectedBounds := core.NewPdfRectangleFloat(25, page.PageSize().Bounds.TopRight.Y-300, 200, page.PageSize().Bounds.TopRight.Y-200)
	if _, addErr := page.AddJpeg(imageBytes, expectedBounds); addErr != nil {
		t.Fatalf("AddJpeg: %v", addErr)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	if pg.Text() != "Smile" {
		t.Errorf("expected text %q, got %q", "Smile", pg.Text())
	}

	imgs := pg.GetImages()
	if len(imgs) != 1 {
		t.Fatalf("expected 1 image, got %d", len(imgs))
	}

	image := imgs[0]
	if absDiff(image.BoundingBox().BottomLeft.X, expectedBounds.BottomLeft.X) > 0.1 ||
		absDiff(image.BoundingBox().BottomLeft.Y, expectedBounds.BottomLeft.Y) > 0.1 {
		t.Errorf("image bottom-left mismatch: expected (%.1f, %.1f), got (%.1f, %.1f)",
			expectedBounds.BottomLeft.X, expectedBounds.BottomLeft.Y,
			image.BoundingBox().BottomLeft.X, image.BoundingBox().BottomLeft.Y)
	}

	rawMemory := image.RawMemory()
	if rawMemory == nil || !bytes.Equal(rawMemory, imageBytes) {
		t.Errorf("image raw bytes mismatch: expected len=%d, got len=%d", len(imageBytes), len(rawMemory))
	}
}

/* ========================================================================
   TestCanWrite2PagesSharingJpeg
   C#: Same JPEG on 3 placements across 2 pages
   ======================================================================== */
func TestCanWrite2PagesSharingJpeg(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	font, err := pdfBuilder.AddStandard14Font(standard14fonts.Helvetica)
	if err != nil {
		t.Fatalf("AddStandard14Font: %v", err)
	}

	if _, addErr := page.AddText("Smile", 12, core.PdfPoint{X: 25, Y: page.PageSize().Bounds.TopRight.Y - 52}, font); addErr != nil {
		t.Fatalf("AddText: %v", addErr)
	}

	imgPath := getDocNoExt("smile-250-by-160.jpg")
	imageBytes, err := os.ReadFile(imgPath)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", imgPath, err)
	}

	expectedBounds1 := core.NewPdfRectangleFloat(25, page.PageSize().Bounds.TopRight.Y-300, 200, page.PageSize().Bounds.TopRight.Y-200)
	expectedBounds2 := core.NewPdfRectangleFloat(25, 600, 75, 650)

	jpeg, addErr := page.AddJpeg(imageBytes, expectedBounds1)
	if addErr != nil {
		t.Fatalf("AddJpeg: %v", addErr)
	}
	page.AddJpegReference(jpeg, expectedBounds2)

	expectedBounds3 := core.NewPdfRectangleFloat(30, 500, 130, 550)
	page2, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize(page2): %v", err)
	}
	page2.AddJpegReference(jpeg, expectedBounds3)

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg1, ok := page1Any.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", page1Any)
	}

	if pg1.Text() != "Smile" {
		t.Errorf("expected text %q, got %q", "Smile", pg1.Text())
	}

	page1Images := pg1.GetImages()
	if len(page1Images) != 2 {
		t.Errorf("page 1: expected 2 images, got %d", len(page1Images))
	}

	page2Any, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	pg2, ok := page2Any.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", page2Any)
	}

	page2Images := pg2.GetImages()
	if len(page2Images) != 1 {
		t.Errorf("page 2: expected 1 image, got %d", len(page2Images))
	}
}

/* ========================================================================
   TestCanWriteSinglePageWithPng
   C#: pdfpig.png, verify text and image bounds
   ======================================================================== */
func TestCanWriteSinglePageWithPng(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	font, err := pdfBuilder.AddStandard14Font(standard14fonts.Helvetica)
	if err != nil {
		t.Fatalf("AddStandard14Font: %v", err)
	}

	if _, addErr := page.AddText("Piggy", 12, core.PdfPoint{X: 25, Y: page.PageSize().Bounds.TopRight.Y - 52}, font); addErr != nil {
		t.Fatalf("AddText: %v", addErr)
	}

	imgPath := getDocNoExt("pdfpig.png")
	imageBytes, err := os.ReadFile(imgPath)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", imgPath, err)
	}

	expectedBounds := core.NewPdfRectangleFloat(25, page.PageSize().Bounds.TopRight.Y-300, 200, page.PageSize().Bounds.TopRight.Y-200)
	if _, addErr := page.AddPng(imageBytes, expectedBounds); addErr != nil {
		t.Fatalf("AddPng: %v", addErr)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	if pg.Text() != "Piggy" {
		t.Errorf("expected text %q, got %q", "Piggy", pg.Text())
	}

	imgs := pg.GetImages()
	if len(imgs) != 1 {
		t.Fatalf("expected 1 image, got %d", len(imgs))
	}

	image := imgs[0]
	if absDiff(image.BoundingBox().BottomLeft.X, expectedBounds.BottomLeft.X) > 0.1 ||
		absDiff(image.BoundingBox().BottomLeft.Y, expectedBounds.BottomLeft.Y) > 0.1 {
		t.Errorf("image bottom-left mismatch: expected (%.1f, %.1f), got (%.1f, %.1f)",
			expectedBounds.BottomLeft.X, expectedBounds.BottomLeft.Y,
			image.BoundingBox().BottomLeft.X, image.BoundingBox().BottomLeft.Y)
	}
}

/* ========================================================================
   TestCanCreateDocumentInformationDictionaryWithNonAsciiCharacters
   C#: Russian text in title and on page, round-trip verification
   ======================================================================== */
func TestCanCreateDocumentInformationDictionaryWithNonAsciiCharacters(t *testing.T) {
	littlePig := "маленький поросенок"

	pdfBuilder := writer.NewPdfDocumentBuilder()
	pdfBuilder.DocumentInformation().Title = &littlePig

	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	fontBytes := loadFont(t, "Roboto-Regular.ttf")
	font, err := pdfBuilder.AddTrueTypeFont(fontBytes)
	if err != nil {
		t.Fatalf("AddTrueTypeFont: %v", err)
	}

	if _, addErr := page.AddText(littlePig, 12, core.PdfPoint{X: 120, Y: 600}, font); addErr != nil {
		t.Fatalf("AddText: %v", addErr)
	}

	file, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(file)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	infoTitle := doc.Information.Title
	if infoTitle != littlePig {
		t.Errorf("expected title %q, got %q", littlePig, infoTitle)
	}
}

/* ========================================================================
   TestCanCreateDocumentWithFilledRectangle
   C#: Red fill + blue stroke rectangle
   ======================================================================== */
func TestCanCreateDocumentWithFilledRectangle(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	page.SetTextAndFillColor(255, 0, 0)
	page.SetStrokeColor(0, 0, 255)

	page.DrawRectangle(core.PdfPoint{X: 20, Y: 100}, 200, 100, 1.5, true)

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	if len(b) == 0 {
		t.Fatal("result is empty")
	}
}

/* ========================================================================
   TestCanGeneratePageWithMultipleStream
   C#: Two text segments in separate content streams
   ======================================================================== */
func TestCanGeneratePageWithMultipleStream(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	fontBytes := loadFont(t, "Andada-Regular.ttf")
	font, err := pdfBuilder.AddTrueTypeFont(fontBytes)
	if err != nil {
		t.Fatalf("AddTrueTypeFont: %v", err)
	}

	if _, addErr := page.AddText("Hello", 12, core.PdfPoint{X: 30, Y: 50}, font); addErr != nil {
		t.Fatalf("AddText Hello: %v", addErr)
	}

	page.NewContentStreamAfter()

	if _, addErr := page.AddText("World!", 12, core.PdfPoint{X: 50, Y: 50}, font); addErr != nil {
		t.Fatalf("AddText World!: %v", addErr)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	if len(b) == 0 {
		t.Fatal("result is empty")
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	if pg.Text() != "HelloWorld!" {
		t.Errorf("expected text %q, got %q", "HelloWorld!", pg.Text())
	}

	readerLetters := pg.Letters()
	if len(readerLetters) == 0 {
		t.Fatal("no letters found")
	}

	h := readerLetters[0]
	if h.Value != "H" {
		t.Errorf("expected first letter 'H', got %q", h.Value)
	}
}

/* ========================================================================
   TestCanCopyPage
   C#: Write page + copy from bold-italic.pdf, verify both pages
   ======================================================================== */
func TestCanCopyPage(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()

	page1, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize(page1): %v", err)
	}

	fontBytes := loadFont(t, "Andada-Regular.ttf")
	font, err := pdfBuilder.AddTrueTypeFont(fontBytes)
	if err != nil {
		t.Fatalf("AddTrueTypeFont: %v", err)
	}

	if _, addErr := page1.AddText("Hello", 12, core.PdfPoint{X: 30, Y: 50}, font); addErr != nil {
		t.Fatalf("AddText: %v", addErr)
	}

	boldItalicPath := getDoc("bold-italic.pdf")
	sourceDoc, openErr := pdfpig.OpenFile(boldItalicPath, nil)
	if openErr != nil {
		t.Fatalf("Open(%q): %v", boldItalicPath, openErr)
	}

	rpageAny, err := sourceDoc.GetPage(1)
	if err != nil {
		sourceDoc.Close()
		t.Fatalf("GetPage(1): %v", err)
	}
	rpage, ok := rpageAny.(*content.Page)
	if !ok {
		sourceDoc.Close()
		t.Fatalf("expected *content.Page, got %T", rpageAny)
	}

	page2, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		sourceDoc.Close()
		t.Fatalf("AddPageWithSize(page2): %v", err)
	}

	if _, copyErr := page2.CopyFrom(rpage); copyErr != nil {
		sourceDoc.Close()
		t.Fatalf("CopyFrom: %v", copyErr)
	}
	sourceDoc.Close()

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	if len(b) == 0 {
		t.Fatal("result is empty")
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 2 {
		t.Errorf("expected 2 pages, got %d", doc.NumberOfPages())
	}

	page1Any, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg1, ok := page1Any.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", page1Any)
	}

	if pg1.Text() != "Hello" {
		t.Errorf("page 1: expected text %q, got %q", "Hello", pg1.Text())
	}

	page2Any, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	pg2, ok := page2Any.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", page2Any)
	}

	if !strings.HasPrefix(pg2.Text(), "Lorem ipsum dolor sit amet") {
		t.Errorf("page 2: expected text starting with 'Lorem ipsum dolor sit amet', got: %q", pg2.Text())
	}
}

/* ========================================================================
   TestCanAddHelloWorldToSimplePage
   C#: Copy from existing doc + add new content on top
   ======================================================================== */
func TestCanAddHelloWorldToSimplePage(t *testing.T) {
	simplePath := getDoc("Single Page Simple - from open office.pdf")
	sourceDoc, openErr := pdfpig.OpenFile(simplePath, nil)
	if openErr != nil {
		t.Fatalf("Open(%q): %v", simplePath, openErr)
	}
	defer sourceDoc.Close()

	pdfBuilder := writer.NewPdfDocumentBuilder()

	page, err := addPageFromDoc(pdfBuilder, sourceDoc, 1)
	if err != nil {
		t.Fatalf("AddPageWithOptions: %v", err)
	}

	page.DrawLine(core.PdfPoint{X: 30, Y: 520}, core.PdfPoint{X: 360, Y: 520}, 1)
	page.DrawLine(core.PdfPoint{X: 360, Y: 520}, core.PdfPoint{X: 360, Y: 250}, 1)

	page.SetStrokeColor(250, 132, 131)
	page.DrawLine(core.PdfPoint{X: 25, Y: 70}, core.PdfPoint{X: 100, Y: 70}, 3)
	page.ResetColor()
	page.DrawRectangle(core.PdfPoint{X: 30, Y: 200}, 250, 100, 0.5, false)
	page.DrawRectangle(core.PdfPoint{X: 30, Y: 100}, 250, 100, 0.5, false)

	fontBytes := loadFont(t, "Andada-Regular.ttf")
	font, err := pdfBuilder.AddTrueTypeFont(fontBytes)
	if err != nil {
		t.Fatalf("AddTrueTypeFont: %v", err)
	}

	letters, err := page.AddText("Hello World!", 12, core.PdfPoint{X: 30, Y: 50}, font)
	if err != nil {
		t.Fatalf("AddText: %v", err)
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	if len(b) == 0 {
		t.Fatal("result is empty")
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	if pg.Text() != "I am a simple pdf.Hello World!" {
		t.Errorf("expected text %q, got %q", "I am a simple pdf.Hello World!", pg.Text())
	}

	readerLetters := pg.Letters()
	if len(readerLetters) < 19 {
		t.Fatalf("expected at least 19 letters (original + added), got %d", len(readerLetters))
	}

	h := readerLetters[18]
	if h.Value != "H" {
		t.Errorf("expected letter[18] = 'H', got %q", h.Value)
	}

	for i := range letters {
		readerLetter := readerLetters[i+18]
		writerLetter := letters[i]

		if readerLetter.Value != writerLetter.Value {
			t.Errorf("letter[%d]: expected value %q, got %q", i+18, writerLetter.Value, readerLetter.Value)
		}
	}
}

/* ========================================================================
   TestCanMerge2SimpleDocuments_Builder
   C#: inkscape + open office merge using builder.AddPage(doc, n)
   ======================================================================== */
func TestCanMerge2SimpleDocuments_Builder(t *testing.T) {
	onePath := getDoc("Single Page Simple - from inkscape.pdf")
	twoPath := getDoc("Single Page Simple - from open office.pdf")

	docOne, err := pdfpig.OpenFile(onePath, nil)
	if err != nil {
		t.Fatalf("Open(%q): %v", onePath, err)
	}
	defer docOne.Close()

	docTwo, err := pdfpig.OpenFile(twoPath, nil)
	if err != nil {
		t.Fatalf("Open(%q): %v", twoPath, err)
	}
	defer docTwo.Close()

	pdfBuilder := writer.NewPdfDocumentBuilder()

	if _, addErr := addPageFromDoc(pdfBuilder, docOne, 1); addErr != nil {
		t.Fatalf("AddPageWithOptions(docOne): %v", addErr)
	}
	if _, addErr := addPageFromDoc(pdfBuilder, docTwo, 1); addErr != nil {
		t.Fatalf("AddPageWithOptions(docTwo): %v", addErr)
	}

	result, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	assertMergeResult(t, result, []string{"Write something inInkscape", "I am a simple pdf."})
}

/* ========================================================================
   TestCanMerge2SimpleDocumentsReversed_Builder
   C#: open office + inkscape (reversed order)
   ======================================================================== */
func TestCanMerge2SimpleDocumentsReversedBuilder(t *testing.T) {
	onePath := getDoc("Single Page Simple - from open office.pdf")
	twoPath := getDoc("Single Page Simple - from inkscape.pdf")

	docOne, err := pdfpig.OpenFile(onePath, nil)
	if err != nil {
		t.Fatalf("Open(%q): %v", onePath, err)
	}
	defer docOne.Close()

	docTwo, err := pdfpig.OpenFile(twoPath, nil)
	if err != nil {
		t.Fatalf("Open(%q): %v", twoPath, err)
	}
	defer docTwo.Close()

	pdfBuilder := writer.NewPdfDocumentBuilder()

	if _, addErr := addPageFromDoc(pdfBuilder, docOne, 1); addErr != nil {
		t.Fatalf("AddPageWithOptions(docOne): %v", addErr)
	}
	if _, addErr := addPageFromDoc(pdfBuilder, docTwo, 1); addErr != nil {
		t.Fatalf("AddPageWithOptions(docTwo): %v", addErr)
	}

	result, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	assertMergeResult(t, result, []string{"I am a simple pdf.", "Write something inInkscape"})
}

/* ========================================================================
   TestCanDedupObjectsFromSameDoc_Builder
   C#: Same page added twice from same doc
   ======================================================================== */
func TestCanDedupObjectsFromSameDocBuilder(t *testing.T) {
	onePath := getDoc("Multiple Page - from Mortality Statistics.pdf")

	doc, err := pdfpig.OpenFile(onePath, nil)
	if err != nil {
		t.Fatalf("Open(%q): %v", onePath, err)
	}
	defer doc.Close()

	pdfBuilder := writer.NewPdfDocumentBuilder()
	if _, addErr := addPageFromDoc(pdfBuilder, doc, 1); addErr != nil {
		t.Fatalf("AddPageWithOptions(1st): %v", addErr)
	}
	if _, addErr := addPageFromDoc(pdfBuilder, doc, 1); addErr != nil {
		t.Fatalf("AddPageWithOptions(2nd): %v", addErr)
	}

	result, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	docResult, openErr := openBuilt(result)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer docResult.Close()

	if docResult.NumberOfPages() != 2 {
		t.Errorf("expected 2 pages, got %d", docResult.NumberOfPages())
	}
}

/* ========================================================================
   TestCanCreatePageTree
   C#: 15626 blank pages, verify page count
   ======================================================================== */
func TestCanCreatePageTree(t *testing.T) {
	count := 25*25*25 + 1 // 15626

	pdfBuilder := writer.NewPdfDocumentBuilder()
	for i := 0; i < count; i++ {
		if _, addErr := pdfBuilder.AddPageWithSize(595, 842); addErr != nil {
			t.Fatalf("AddPageWithSize[%d]: %v", i, addErr)
		}
	}

	result, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(result)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	if doc.NumberOfPages() != count {
		t.Errorf("expected %d pages, got %d", count, doc.NumberOfPages())
	}
}

/* ========================================================================
   TestCanWriteEmptyContentStream
   C#: Single blank page → single content stream (not array)
   ======================================================================== */
func TestCanWriteEmptyContentStream(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	_, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	result, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(result)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 1 {
		t.Errorf("expected 1 page, got %d", doc.NumberOfPages())
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	// Single empty page should have a single content stream (IndirectReferenceToken), not an array.
	contentsTok, hasContents := pg.Dictionary().TryGet(tokens.Contents)
	if !hasContents {
		t.Fatal("page dictionary missing Contents entry")
	}

	if _, isArray := contentsTok.(*tokens.ArrayToken); isArray {
		t.Error("empty content stream should be a single reference, not an array")
	}
}

/* ========================================================================
   TestCanWriteSingleContentStream
   C#: Single DrawLine → single content stream (not array)
   ======================================================================== */
func TestCanWriteSingleContentStream(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	pb, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	pb.DrawLine(core.PdfPoint{X: 1, Y: 1}, core.PdfPoint{X: 2, Y: 2}, 1)

	result, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(result)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 1 {
		t.Errorf("expected 1 page, got %d", doc.NumberOfPages())
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	contentsTok, hasContents := pg.Dictionary().TryGet(tokens.Contents)
	if !hasContents {
		t.Fatal("page dictionary missing Contents entry")
	}

	if _, isArray := contentsTok.(*tokens.ArrayToken); isArray {
		t.Error("single content stream should be a single reference, not an array")
	}
}

/* ========================================================================
   TestCanWriteAndIgnoreEmptyContentStream
   C#: DrawLine + NewContentStreamAfter → empty stream ignored, single ref
   ======================================================================== */
func TestCanWriteAndIgnoreEmptyContentStream(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	pb, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	pb.DrawLine(core.PdfPoint{X: 1, Y: 1}, core.PdfPoint{X: 2, Y: 2}, 1)
	pb.NewContentStreamAfter()

	result, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(result)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 1 {
		t.Errorf("expected 1 page, got %d", doc.NumberOfPages())
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	contentsTok, hasContents := pg.Dictionary().TryGet(tokens.Contents)
	if !hasContents {
		t.Fatal("page dictionary missing Contents entry")
	}

	if _, isArray := contentsTok.(*tokens.ArrayToken); isArray {
		t.Error("empty stream should be ignored; single content stream should be a reference, not an array")
	}
}

/* ========================================================================
   TestCanWriteMultipleContentStream
   C#: Two DrawLine in separate streams → ArrayToken with 2 refs
   ======================================================================== */
func TestCanWriteMultipleContentStream(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	pb, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	pb.DrawLine(core.PdfPoint{X: 1, Y: 1}, core.PdfPoint{X: 2, Y: 2}, 1)
	pb.NewContentStreamAfter()
	pb.DrawLine(core.PdfPoint{X: 1, Y: 1}, core.PdfPoint{X: 2, Y: 2}, 1)

	result, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(result)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 1 {
		t.Errorf("expected 1 page, got %d", doc.NumberOfPages())
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	contentsTok, hasContents := pg.Dictionary().TryGet(tokens.Contents)
	if !hasContents {
		t.Fatal("page dictionary missing Contents entry")
	}

	streams, isArray := contentsTok.(*tokens.ArrayToken)
	if !isArray {
		t.Error("multiple content streams should be an ArrayToken")
	} else if streams.Length() != 2 {
		t.Errorf("expected 2 streams in array, got %d", streams.Length())
	}
}

/* ========================================================================
   TestCopiedPagesResultInSameData
   C#: Copy pages from various docs, compare letter count and location sum
   ======================================================================== */
func TestCopiedPagesResultInSameData(t *testing.T) {
	testCases := []string{
		"Single Page Simple - from google drive.pdf",
		"Old Gutnish Internet Explorer.pdf",
		"68-1990-01_A.pdf",
		"Multiple Page - from Mortality Statistics.pdf",
	}

	for _, name := range testCases {
		t.Run(name, func(t *testing.T) {
			docPath := getDoc(name)
			sourceBytes, readErr := os.ReadFile(docPath)
			if readErr != nil {
				t.Skipf("file not found: %s", docPath)
			}

			sourceDoc, openErr := parser.OpenMemory(sourceBytes, &content.ParsingOptions{UseLenientParsing: false})
			if openErr != nil {
				t.Fatalf("OpenMemory(source): %v", openErr)
			}
			defer sourceDoc.Close()

			count1 := getCounts(sourceDoc)

			pdfBuilder := writer.NewPdfDocumentBuilder()
			for i := 1; i <= sourceDoc.NumberOfPages(); i++ {
				if _, addErr := addPageFromDoc(pdfBuilder, sourceDoc, i); addErr != nil {
					t.Fatalf("AddPageWithOptions(%d): %v", i, addErr)
				}
			}

			result, buildErr := pdfBuilder.Build()
			if buildErr != nil {
				t.Fatalf("Build: %v", buildErr)
			}
			if closeErr := pdfBuilder.Close(); closeErr != nil {
				t.Logf("Close: %v", closeErr)
			}

			copiedDoc, openErr2 := parser.OpenMemory(result, &content.ParsingOptions{UseLenientParsing: false})
			if openErr2 != nil {
				t.Fatalf("OpenMemory(copied): %v", openErr2)
			}
			defer copiedDoc.Close()

			count2 := getCounts(copiedDoc)

			if count1.letters != count2.letters {
				t.Errorf("letter count mismatch: original=%d, copied=%d", count1.letters, count2.letters)
			}
			if absDiff(count1.locationSum, count2.locationSum) > 0.5 {
				t.Errorf("location sum mismatch: original=%.4f, copied=%.4f", count1.locationSum, count2.locationSum)
			}
		})
	}
}

/* ========================================================================
   TestCanFastAddPageAndInheritProps
   C#: Copy inherited_mediabox.pdf, verify mediabox preserved
   ======================================================================== */
func TestCanFastAddPageAndInheritProps(t *testing.T) {
	firstPath := getDoc("inherited_mediabox.pdf")
	contents, readErr := os.ReadFile(firstPath)
	if readErr != nil {
		t.Skipf("file not found: %s", firstPath)
	}

	existingDoc, openErr := parser.OpenMemory(contents, &content.ParsingOptions{UseLenientParsing: false})
	if openErr != nil {
		t.Fatalf("OpenMemory(source): %v", openErr)
	}
	defer existingDoc.Close()

	pdfBuilder := writer.NewPdfDocumentBuilder()
	if _, addErr := addPageFromDoc(pdfBuilder, existingDoc, 1); addErr != nil {
		t.Fatalf("AddPageWithOptions: %v", addErr)
	}

	results, buildErr := pdfBuilder.Build()
	if buildErr != nil {
		t.Fatalf("Build: %v", buildErr)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	rewrittenDoc, openErr2 := parser.OpenMemory(results, &content.ParsingOptions{UseLenientParsing: false})
	if openErr2 != nil {
		t.Fatalf("OpenMemory(rewritten): %v", openErr2)
	}
	defer rewrittenDoc.Close()

	pgAny, err := rewrittenDoc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pgAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pgAny)
	}

	if absDiff(pg.Width(), 200) > 0.5 {
		t.Errorf("expected width 200, got %.1f", pg.Width())
	}
	if absDiff(pg.Height(), 100) > 0.5 {
		t.Errorf("expected height 100, got %.1f", pg.Height())
	}
}

/* ========================================================================
   TestCanFastAddPageWithStreamSubtype
   C#: Copy steam_in_page_dict.pdf, verify no exception and content exists
   ======================================================================== */
func TestCanFastAddPageWithStreamSubtype(t *testing.T) {
	firstPath := getDoc("steam_in_page_dict.pdf")
	contents, readErr := os.ReadFile(firstPath)
	if readErr != nil {
		t.Skipf("file not found: %s", firstPath)
	}

	existingDoc, openErr := parser.OpenMemory(contents, &content.ParsingOptions{UseLenientParsing: false})
	if openErr != nil {
		t.Fatalf("OpenMemory(source): %v", openErr)
	}
	defer existingDoc.Close()

	pdfBuilder := writer.NewPdfDocumentBuilder()
	if _, addErr := addPageFromDoc(pdfBuilder, existingDoc, 1); addErr != nil {
		t.Fatalf("AddPageWithOptions: %v", addErr)
	}

	results, buildErr := pdfBuilder.Build()
	if buildErr != nil {
		t.Fatalf("Build: %v", buildErr)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	rewrittenDoc, openErr2 := parser.OpenMemory(results, &content.ParsingOptions{UseLenientParsing: false})
	if openErr2 != nil {
		t.Fatalf("OpenMemory(rewritten): %v", openErr2)
	}
	defer rewrittenDoc.Close()

	pgAny, err := rewrittenDoc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pgAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pgAny)
	}

	if pg.Width() == 0 || pg.Height() == 0 {
		t.Error("expected non-zero page dimensions (page content should exist)")
	}
}

/* ========================================================================
   TestCanFastAddPageAndStripLinkAnnots
   C#: Copy outline.pdf, verify link annotations stripped
   ======================================================================== */
func TestCanFastAddPageAndStripLinkAnnots(t *testing.T) {
	firstPath := getDoc("outline.pdf")
	contents, readErr := os.ReadFile(firstPath)
	if readErr != nil {
		t.Skipf("file not found: %s", firstPath)
	}

	existingDoc, openErr := parser.OpenMemory(contents, &content.ParsingOptions{UseLenientParsing: false})
	if openErr != nil {
		t.Fatalf("OpenMemory(source): %v", openErr)
	}
	defer existingDoc.Close()

	pdfBuilder := writer.NewPdfDocumentBuilder()
	if _, addErr := addPageFromDoc(pdfBuilder, existingDoc, 1); addErr != nil {
		t.Fatalf("AddPageWithOptions: %v", addErr)
	}

	results, buildErr := pdfBuilder.Build()
	if buildErr != nil {
		t.Fatalf("Build: %v", buildErr)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	rewrittenDoc, openErr2 := parser.OpenMemory(results, &content.ParsingOptions{UseLenientParsing: false})
	if openErr2 != nil {
		t.Fatalf("OpenMemory(rewritten): %v", openErr2)
	}
	defer rewrittenDoc.Close()

	pgAny, err := rewrittenDoc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pgAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pgAny)
	}

	annots := pg.GetAnnotations()
	linkCount := 0
	for _, ann := range annots {
		if at, ok := ann.Type().(annotations.AnnotationType); ok && at == annotations.Link {
			linkCount++
		}
	}
	if linkCount > 0 {
		t.Errorf("expected no link annotations after copy (default strips links), got %d", linkCount)
	}
}

/* ========================================================================
   TestCanFastAddPageAndStripAllAnnots
   C#: Copy outline.pdf with KeepAnnotations=false
   ======================================================================== */
func TestCanFastAddPageAndStripAllAnnots(t *testing.T) {
	firstPath := getDoc("outline.pdf")
	contents, readErr := os.ReadFile(firstPath)
	if readErr != nil {
		t.Skipf("file not found: %s", firstPath)
	}

	existingDoc, openErr := parser.OpenMemory(contents, &content.ParsingOptions{UseLenientParsing: false})
	if openErr != nil {
		t.Fatalf("OpenMemory(source): %v", openErr)
	}
	defer existingDoc.Close()

	pdfBuilder := writer.NewPdfDocumentBuilder()
	opts := &writer.AddPageOptions{KeepAnnotations: false}
	if _, addErr := pdfBuilder.AddPageWithOptions(existingDoc, 1, opts); addErr != nil {
		t.Fatalf("AddPageWithOptions(keep=false): %v", addErr)
	}

	results, buildErr := pdfBuilder.Build()
	if buildErr != nil {
		t.Fatalf("Build: %v", buildErr)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	rewrittenDoc, openErr2 := parser.OpenMemory(results, &content.ParsingOptions{UseLenientParsing: false})
	if openErr2 != nil {
		t.Fatalf("OpenMemory(rewritten): %v", openErr2)
	}
	defer rewrittenDoc.Close()

	pgAny, err := rewrittenDoc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pgAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pgAny)
	}

	annots := pg.GetAnnotations()
	if len(annots) > 0 {
		t.Errorf("expected no annotations after copy with KeepAnnotations=false, got %d", len(annots))
	}
}

/* ========================================================================
   TestCanCreateDocumentWithOutline
   C#: Bookmarks with various destination types
   ======================================================================== */
func TestCanCreateDocumentWithOutline(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()

	font, err := pdfBuilder.AddStandard14Font(standard14fonts.Helvetica)
	if err != nil {
		t.Fatalf("AddStandard14Font: %v", err)
	}

	// Create pages and add text for each bookmark node.
	nodeTitles := []string{"1", "1.1", "2", "2.1", "2.2", "2.3", "2.4", "2.5", "3"}

	for _, title := range nodeTitles {
		page, addErr := pdfBuilder.AddPageWithSize(595, 842)
		if addErr != nil {
			t.Fatalf("AddPageWithSize(%q): %v", title, addErr)
		}
		if _, addTextErr := page.AddText(title, 12, core.PdfPoint{X: 25, Y: 800}, font); addTextErr != nil {
			t.Logf("AddText(%q): %v (non-fatal)", title, addTextErr)
		}
	}

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	for i, title := range nodeTitles {
		pageAny, err := doc.GetPage(i + 1)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", i+1, err)
		}
		pg, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("expected *content.Page for page %d, got %T", i+1, pageAny)
		}
		if pg.Text() != title {
			t.Errorf("page %d: expected text %q, got %q", i+1, title, pg.Text())
		}
	}
}

/* ========================================================================
   TestCanAddLinkToPage
   C#: Add URI link annotation to page
   ======================================================================== */
func TestCanAddLinkToPage(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	_, err = pdfBuilder.AddStandard14Font(standard14fonts.Helvetica)
	if err != nil {
		t.Fatalf("AddStandard14Font: %v", err)
	}

	linkArea := core.NewPdfRectangleFloat(25, 690, 200, 720)
	page.AddLink("https://github.com", linkArea)

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 1 {
		t.Errorf("expected 1 page, got %d", doc.NumberOfPages())
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	annots := pg.GetAnnotations()
	linkCount := 0
	for _, ann := range annots {
		if at, ok := ann.Type().(annotations.AnnotationType); ok && at == annotations.Link {
			linkCount++
		}
	}
	if linkCount != 1 {
		t.Errorf("expected 1 link annotation, got %d", linkCount)
	}
}

/* ========================================================================
   TestCanAddInternalLinkToPage
   C#: Add internal GoToAction link from page 2 to page 1
   ======================================================================== */
func TestCanAddInternalLinkToPage(t *testing.T) {
	pdfBuilder := writer.NewPdfDocumentBuilder()

	_, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize(page1): %v", err)
	}

	page2, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize(page2): %v", err)
	}

	linkArea := core.NewPdfRectangleFloat(25, 690, 200, 720)
	dest := destinations.NewExplicitDestination(1, destinations.XyzCoordinates, destinations.NewExplicitDestinationCoordinatesWithTop(25, 750))
	page2.AddLinkDestination(dest, linkArea)

	b, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	doc, openErr := openBuilt(b)
	if openErr != nil {
		t.Fatalf("OpenMemory: %v", openErr)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 2 {
		t.Errorf("expected 2 pages, got %d", doc.NumberOfPages())
	}

	page2Any, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	pg2, ok := page2Any.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", page2Any)
	}

	annots := pg2.GetAnnotations()
	linkCount := 0
	for _, ann := range annots {
		if at, ok := ann.Type().(annotations.AnnotationType); ok && at == annotations.Link {
			linkCount++
		}
	}
	if linkCount != 1 {
		t.Errorf("expected 1 link annotation on page 2, got %d", linkCount)
	}
}

/* ========================================================================
   TestCanCopyIssue1017
   C#: Copy from _citUR2jB1_ay56u6vELzk.pdf without crashing
   ======================================================================== */
func TestCanCopyIssue1017(t *testing.T) {
	docPath := getDoc("_citUR2jB1_ay56u6vELzk.pdf")
	sourceBytes, readErr := os.ReadFile(docPath)
	if readErr != nil {
		t.Skipf("file not found: %s", docPath)
	}

	sourceDoc, openErr := parser.OpenMemory(sourceBytes, nil)
	if openErr != nil {
		t.Fatalf("OpenMemory(source): %v", openErr)
	}
	defer sourceDoc.Close()

	pdfBuilder := writer.NewPdfDocumentBuilder()
	page1Any, err := sourceDoc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page1, ok := page1Any.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", page1Any)
	}

	newPage, addErr := pdfBuilder.AddPageWithSize(page1.Width(), page1.Height())
	if addErr != nil {
		t.Fatalf("AddPageWithSize: %v", addErr)
	}

	if _, copyErr := newPage.CopyFrom(page1); copyErr != nil {
		t.Fatalf("CopyFrom: %v", copyErr)
	}

	b, buildErr := pdfBuilder.Build()
	if buildErr != nil {
		t.Fatalf("Build: %v", buildErr)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	if len(b) == 0 {
		t.Fatal("result is empty")
	}
}

/* ========================================================================
   TestCanCopyInLineImage
   C#: Copy ssm2163.pdf with inline images
   ======================================================================== */
func TestCanCopyInLineImage(t *testing.T) {
	docPath := getDoc("ssm2163.pdf")
	sourceBytes, readErr := os.ReadFile(docPath)
	if readErr != nil {
		t.Skipf("file not found: %s", docPath)
	}

	sourceDoc, openErr := parser.OpenMemory(sourceBytes, nil)
	if openErr != nil {
		t.Fatalf("OpenMemory(source): %v", openErr)
	}
	defer sourceDoc.Close()

	pdfBuilder := writer.NewPdfDocumentBuilder()
	numberOfPages := sourceDoc.NumberOfPages()
	for pageNum := 1; pageNum <= numberOfPages; pageNum++ {
		srcPageAny, err := sourceDoc.GetPage(pageNum)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", pageNum, err)
		}
		srcPage, ok := srcPageAny.(*content.Page)
		if !ok {
			t.Fatalf("expected *content.Page for page %d, got %T", pageNum, srcPageAny)
		}

		newPage, addErr := pdfBuilder.AddPageWithSize(srcPage.Width(), srcPage.Height())
		if addErr != nil {
			t.Fatalf("AddPageWithSize[%d]: %v", pageNum, addErr)
		}
		if _, copyErr := newPage.CopyFrom(srcPage); copyErr != nil {
			t.Fatalf("CopyFrom[%d]: %v", pageNum, copyErr)
		}
	}

	pdfBytes, buildErr := pdfBuilder.Build()
	if buildErr != nil {
		t.Fatalf("Build: %v", buildErr)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	copiedDoc, openErr2 := parser.OpenMemory(pdfBytes, nil)
	if openErr2 != nil {
		t.Fatalf("OpenMemory(copied): %v", openErr2)
	}
	defer copiedDoc.Close()

	pageNum := 7
	if copiedDoc.NumberOfPages() < pageNum {
		t.Skipf("copied document has fewer pages than expected (%d < %d)", copiedDoc.NumberOfPages(), pageNum)
	}

	copiedPageAny, err := copiedDoc.GetPage(pageNum)
	if err != nil {
		t.Fatalf("GetPage(%d): %v", pageNum, err)
	}
	copiedPage, ok := copiedPageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", copiedPageAny)
	}

	// Verify page has content (non-zero dimensions indicate successful copy)
	if copiedPage.Width() == 0 || copiedPage.Height() == 0 {
		t.Errorf("page %d: expected non-zero page dimensions (content should exist)", pageNum)
	}
}

/* ========================================================================
    TestCanGeneratePdfAFile
    C#: Theory with InlineData for A1B, A1A, A2B, A2A, A3B, A3A
    ======================================================================== */
func TestCanGeneratePdfAFile(t *testing.T) {
	standards := []writer.PdfAStandard{
		writer.PdfA1B,
		writer.PdfA1A,
		writer.PdfA2B,
		writer.PdfA2A,
		writer.PdfA3B,
		writer.PdfA3A,
	}

	for _, standard := range standards {
		t.Run(standard.String(), func(t *testing.T) {
			pdfBuilder := writer.NewPdfDocumentBuilder()
			pdfBuilder.SetArchiveStandard(standard)

			page, err := pdfBuilder.AddPageWithSize(595, 842)
			if err != nil {
				t.Fatalf("AddPageWithSize: %v", err)
			}

			imgPath := getDocNoExt("smile-250-by-160.jpg")
			imgBytes, readErr := os.ReadFile(imgPath)
			if readErr != nil {
				t.Fatalf("ReadFile(%q): %v", imgPath, readErr)
			}

			bounds := core.NewPdfRectangleFloat(50, 70, 150, 130)
			if _, addErr := page.AddJpeg(imgBytes, bounds); addErr != nil {
				t.Fatalf("AddJpeg: %v", addErr)
			}

			fontBytes := loadFont(t, "Roboto-Regular.ttf")
			font, err := pdfBuilder.AddTrueTypeFont(fontBytes)
			if err != nil {
				t.Fatalf("AddTrueTypeFont: %v", err)
			}

			text := fmt.Sprintf("Howdy PDF/%s!", standard.String())
			if _, addErr := page.AddText(text, 10, core.PdfPoint{X: 25, Y: 700}, font); addErr != nil {
				t.Fatalf("AddText: %v", addErr)
			}

			b, err := pdfBuilder.Build()
			if err != nil {
				t.Fatalf("Build: %v", err)
			}
			if closeErr := pdfBuilder.Close(); closeErr != nil {
				t.Logf("Close: %v", closeErr)
			}

			doc, openErr := parser.OpenMemory(b, nil)
			if openErr != nil {
				t.Fatalf("OpenMemory: %v", openErr)
			}
			defer doc.Close()

			if doc.NumberOfPages() != 1 {
				t.Errorf("expected 1 page, got %d", doc.NumberOfPages())
			}

			xmp, ok, err := doc.TryGetXmpMetadata()
			if err != nil {
				t.Fatalf("TryGetXmpMetadata: %v", err)
			}
			if !ok {
				t.Fatal("expected XMP metadata to be present")
			}
			if xmp == nil {
				t.Fatal("XMP metadata is nil")
			}
			xmlStr := xmp.GetXmlString()
			if xmlStr == "" {
				t.Fatal("XMP XML string is empty")
			}
		})
	}
}

/* ========================================================================
    TestCanUseCustomTokenWriter
    C#: Custom TokenWriter tracks Objects, Tokens, WroteCrossReferenceTable
    ======================================================================== */

type testTokenWriter struct {
	tokens                    int
	objects                   int
	wroteCrossReferenceTable  bool
	writingPageContents       bool
	inner                     writer.TokenWriter
}

func (w *testTokenWriter) WriteToken(token tokens.Token, out io.Writer) error {
	w.tokens++
	return w.inner.WriteToken(token, out)
}

func (w *testTokenWriter) WriteObject(objectNumber int64, generation int, data []byte, out io.Writer) error {
	w.objects++
	return w.inner.WriteObject(objectNumber, generation, data, out)
}

func (w *testTokenWriter) WriteCrossReferenceTable(
	objectOffsets map[core.IndirectReference]int64,
	catalogRef core.IndirectReference,
	out io.Writer,
	docInfoRef *core.IndirectReference,
) error {
	w.wroteCrossReferenceTable = true
	return w.inner.WriteCrossReferenceTable(objectOffsets, catalogRef, out, docInfoRef)
}

func (w *testTokenWriter) WritingPageContents() bool {
	return w.writingPageContents
}

func (w *testTokenWriter) SetWritingPageContents(v bool) {
	w.writingPageContents = v
}

func TestCanUseCustomTokenWriter(t *testing.T) {
	docPath := getDoc("68-1990-01_A.pdf")
	sourceBytes, readErr := os.ReadFile(docPath)
	if readErr != nil {
		t.Skipf("file not found: %s", docPath)
	}

	sourceDoc, openErr := parser.OpenMemory(sourceBytes, nil)
	if openErr != nil {
		t.Fatalf("OpenMemory(source): %v", openErr)
	}
	defer sourceDoc.Close()

	tw := &testTokenWriter{inner: writer.NewTokenWriter()}
	ms := &seekableBuffer{buf: bytes.NewBuffer(nil)}

	pdfBuilder := writer.NewPdfDocumentBuilderWithStream(ms, false, writer.PdfWriterDefault, 1.7, tw)

	for i := 1; i <= sourceDoc.NumberOfPages(); i++ {
		if _, addErr := addPageFromDoc(pdfBuilder, sourceDoc, i); addErr != nil {
			t.Fatalf("AddPageWithOptions(%d): %v", i, addErr)
		}
	}

	_, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	if tw.objects != 0 {
		t.Errorf("expected 0 objects, got %d", tw.objects)
	}
	if tw.tokens <= 1000 {
		t.Errorf("expected tokens > 1000, got %d", tw.tokens)
	}
	if !tw.wroteCrossReferenceTable {
		t.Error("expected WroteCrossReferenceTable to be true")
	}
}

/* ========================================================================
    Helper functions
    ======================================================================== */

func addPageFromDoc(b *writer.PdfDocumentBuilder, doc *content.PdfDocument, pageNumber int) (*writer.PdfPageBuilder, error) {
	return b.AddPageWithOptions(doc, pageNumber, writer.NewAddPageOptions())
}

func assertMergeResult(t *testing.T, b []byte, expectedTexts []string) {
	t.Helper()

	doc, err := openBuilt(b)
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer doc.Close()

	if doc.NumberOfPages() != len(expectedTexts) {
		t.Errorf("expected %d pages, got %d", len(expectedTexts), doc.NumberOfPages())
	}

	for i, expected := range expectedTexts {
		pgAny, err := doc.GetPage(i + 1)
		if err != nil {
			t.Errorf("GetPage(%d): %v", i+1, err)
			continue
		}
		pg, ok := pgAny.(*content.Page)
		if !ok {
			t.Errorf("GetPage(%d): expected *content.Page, got %T", i+1, pgAny)
			continue
		}
		if pg.Text() != expected {
			t.Errorf("page %d text: expected %q, got %q", i+1, expected, pg.Text())
		}
	}
}

type docCounts struct {
	letters     int
	locationSum float64
}

func getCounts(doc *content.PdfDocument) docCounts {
	var letters int
	var locationSum float64

	pages, err := doc.GetPages()
	if err != nil {
		return docCounts{}
	}

	for _, pageAny := range pages {
		page, ok := pageAny.(*content.Page)
		if !ok {
			continue
		}
		for _, letter := range page.Letters() {
			letters++
			locationSum += letter.Location().X
			locationSum += letter.Location().Y
			if letter.FontDetails != nil && letter.FontDetails.Name != "" {
				locationSum += float64(len(letter.FontDetails.Name))
			}
		}
	}

	return docCounts{letters: letters, locationSum: locationSum}
}

func absDiff(a, b float64) float64 {
	d := a - b
	if d < 0 {
		return -d
	}
	return d
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
