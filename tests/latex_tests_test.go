//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/export"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/text_extractor"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func getLaTexFilename() string {
	return testutil.GetDocumentPath("ICML03-081.pdf", true)
}

// TestLaTexCanReadContent verifies text content on pages 1 and 2.
func TestLaTexCanReadContent(t *testing.T) {
	doc, err := pdfpig.OpenFile(getLaTexFilename(), &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page1, ok := page1Any.(*content.Page)
	if !ok {
		t.Fatal("page 1 is not *content.Page")
	}

	expectedPage1 := "TacklingthePoorAssumptionsofNaiveBayesTextClassiﬁers"
	if !strings.Contains(page1.Text(), expectedPage1) {
		t.Errorf("Page 1 text does not contain %q", expectedPage1)
	}

	page2Any, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	page2, ok := page2Any.(*content.Page)
	if !ok {
		t.Fatal("page 2 is not *content.Page")
	}

	expectedPage2 := "is~θc={θc1,θc2,...,θcn},"
	if !strings.Contains(page2.Text(), expectedPage2) {
		t.Errorf("Page 2 text does not contain %q", expectedPage2)
	}
}

// TestLaTexLettersHaveHeight verifies the first letter has non-zero bounding box height.
func TestLaTexLettersHaveHeight(t *testing.T) {
	doc, err := pdfpig.OpenFile(getLaTexFilename(), nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
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
	if len(letters) == 0 {
		t.Fatal("no letters found on page 1")
	}

	if letters[0].BoundingBox.Height == 0 {
		t.Errorf("Letters[0].BoundingBox.Height = 0, want non-zero")
	}
}

// TestLaTexHasCorrectNumberOfPages verifies the document has 8 pages.
func TestLaTexHasCorrectNumberOfPages(t *testing.T) {
	doc, err := pdfpig.OpenFile(getLaTexFilename(), nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 8 {
		t.Errorf("NumberOfPages = %d, want 8", got)
	}
}

// TestLaTexLettersHaveCorrectPositionsXfinium compares letter positions against
// Xfinium reference data (height check disabled).
func TestLaTexLettersHaveCorrectPositionsXfinium(t *testing.T) {
	positions := getXfiniumPositionDataLatex()

	doc, err := pdfpig.OpenFile(getLaTexFilename(), nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
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
	for i := 0; i < len(letters); i++ {
		if i >= len(positions) {
			break
		}

		positions[i].AssertWithinTolerance(t, letters[i], false)
	}
}

// TestLaTexLettersHaveCorrectPositionsPdfBox compares letter positions against
// PdfBox reference data read from file (height check enabled).
func TestLaTexLettersHaveCorrectPositionsPdfBox(t *testing.T) {
	positions, err := getPdfBoxPositionDataLatex()
	if err != nil {
		t.Fatalf("getPdfBoxPositionData: %v", err)
	}

	doc, err := pdfpig.OpenFile(getLaTexFilename(), nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
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
	for i := 0; i < len(letters); i++ {
		if i >= len(positions) {
			break
		}

		positions[i].AssertWithinTolerance(t, letters[i], true)
	}
}

// TestLaTexPage1Words verifies the first N words on page 1 match expected text.
func TestLaTexPage1Words(t *testing.T) {
	expectedString := `Tackling the Poor Assumptions of Naive Bayes Text Classiﬁers
Jason D. M. Rennie jrennie@mit.edu
Lawrence Shih kai@mit.edu
Jaime Teevan teevan@mit.edu
David R. Karger karger@mit.edu
Artiﬁcial Intelligence Laboratory; Massachusetts Institute of Technology; Cambridge, MA 02139
Abstract amples. To balance the amount of training examples
used per estimate, we introduce a “complement class”
Naive Bayes is often used as a baseline in`

	expected := strings.Fields(expectedString)

	doc, err := pdfpig.OpenFile(getLaTexFilename(), nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
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

	words := page.GetWords()

	for i := 0; i < len(words); i++ {
		if i >= len(expected) {
			break
		}

		if words[i].Text != expected[i] {
			t.Errorf("Word at index %d: expected %q, got %q", i, expected[i], words[i].Text)
		}
	}
}

// TestLaTexCanGetMetadata verifies the document has XMP metadata with correct header.
func TestLaTexCanGetMetadata(t *testing.T) {
	doc, err := pdfpig.OpenFile(getLaTexFilename(), &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	metadata, ok, err := doc.TryGetXmpMetadata()
	if err != nil {
		t.Fatalf("TryGetXmpMetadata: %v", err)
	}

	if !ok {
		t.Fatal("document has no XMP metadata")
	}

	if metadata == nil {
		t.Fatal("XMP metadata is nil")
	}

	xmlBytes := metadata.GetXmlBytes()
	text := core.BytesAsLatin1String(xmlBytes)

	expectedPrefix := "<?xpacket begin='' id='W5M0MpCehiHzreSzNTczkc9d'"
	if !strings.HasPrefix(text, expectedPrefix) {
		t.Errorf("XMP metadata does not start with %q, got prefix %q", expectedPrefix, text[:min(len(text), len(expectedPrefix)+10)])
	}
}

// TestLaTexCanExportSvg verifies SVG export produces non-empty output.
func TestLaTexCanExportSvg(t *testing.T) {
	doc, err := pdfpig.OpenFile(getLaTexFilename(), &content.ParsingOptions{ClipPaths: true})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
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

	exporter := export.NewSvgTextExporter(export.DoNotCheck)
	svg := exporter.Get(page)

	if svg == "" {
		t.Error("SVG export returned empty string")
	}
}

// TestLaTexCanExtractContentOrderText verifies content order text extraction
// produces non-empty output for every page.
func TestLaTexCanExtractContentOrderText(t *testing.T) {
	doc, err := pdfpig.OpenFile(getLaTexFilename(), nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	pages, err := doc.GetPages()
	if err != nil {
		t.Fatalf("GetPages: %v", err)
	}

	for i, pg := range pages {
		page, ok := pg.(*content.Page)
		if !ok {
			t.Errorf("Page %d: expected *content.Page, got %T", i+1, pg)
			continue
		}
		text := text_extractor.GetText(page, false)
		if text == "" {
			t.Errorf("Page %d: ContentOrderTextExtractor returned empty text", i+1)
		}
	}
}

// getXfiniumPositionData returns the hardcoded Xfinium reference position data.
func getXfiniumPositionDataLatex() []*testutil.AssertablePositionData {
	data := `75.731	83.12866	11.218572	T	14.346	WDKAAR+CMBX12	9.956124
85.6153934	83.123866	7.847262	a	11.218572	WDKAAR+CMBX12	9.956124
93.462656	83.123866	7.173	c	11.218572	WDKAAR+CMBX12	9.956124
100.176584	83.123866	8.521524	k	11.218572	WDKAAR+CMBX12	9.956124
108.698108	83.123866	4.490298	l	11.218572	WDKAAR+CMBX12	9.956124`

	lines := strings.Split(data, "\n")
	var result []*testutil.AssertablePositionData
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pd, err := testutil.ParseAssertablePositionData(line)
		if err != nil {
			panic("invalid xfinium position data: " + err.Error())
		}
		result = append(result, pd)
	}

	return result
}

// getPdfBoxPositionData reads the PdfBox reference position data from file.
func getPdfBoxPositionDataLatex() ([]*testutil.AssertablePositionData, error) {
	path := filepath.Join(integrationDocRoot, "ICML03-081.Page1.Positions.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var result []*testutil.AssertablePositionData
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pd, err := testutil.ParseAssertablePositionData(line)
		if err != nil {
			return nil, err
		}
		result = append(result, pd)
	}

	return result, nil
}
