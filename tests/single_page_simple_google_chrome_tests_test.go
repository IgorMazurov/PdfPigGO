//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"math"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

const singlePageSimpleGoogleChromeDocName = "Single Page Simple - from google drive.pdf"

func getSinglePageSimpleGoogleChromePath() string {
	return filepath.Join(integrationDocRoot, singlePageSimpleGoogleChromeDocName)
}

var ignoredHiddenCharacters = map[string]bool{
	"\u200B": true,
}

// TestSinglePageSimpleGoogleChromeHasCorrectNumberOfPages verifies the document has 1 page.
func TestSinglePageSimpleGoogleChromeHasCorrectNumberOfPages(t *testing.T) {
	path := getSinglePageSimpleGoogleChromePath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 1 {
		t.Errorf("NumberOfPages = %d, want 1", got)
	}
}

// TestSinglePageSimpleGoogleChromeCanAccessPage verifies page 1 can be accessed.
func TestSinglePageSimpleGoogleChromeCanAccessPage(t *testing.T) {
	path := getSinglePageSimpleGoogleChromePath()

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

	if page == nil {
		t.Error("page is nil")
	}

	if got := page.Number(); got != 1 {
		t.Errorf("page.Number = %d, want 1", got)
	}
}

// TestSinglePageSimpleGoogleChromeAccessPageLowerThanOneThrows verifies GetPage(0) returns error.
func TestSinglePageSimpleGoogleChromeAccessPageLowerThanOneThrows(t *testing.T) {
	path := getSinglePageSimpleGoogleChromePath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	_, err = doc.GetPage(0)
	if err == nil {
		t.Error("GetPage(0) expected error, got nil")
	}
}

// TestSinglePageSimpleGoogleChromePageHasCorrectDimensions verifies page dimensions are 612x792.
func TestSinglePageSimpleGoogleChromePageHasCorrectDimensions(t *testing.T) {
	path := getSinglePageSimpleGoogleChromePath()

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

	if got := page.Width(); got != 612 {
		t.Errorf("page.Width = %g, want 612", got)
	}

	if got := page.Height(); got != 792 {
		t.Errorf("page.Height = %g, want 792", got)
	}
}

// TestSinglePageSimpleGoogleChromePageHasCorrectTextIgnoringHiddenCharacters verifies extracted text.
func TestSinglePageSimpleGoogleChromePageHasCorrectTextIgnoringHiddenCharacters(t *testing.T) {
	path := getSinglePageSimpleGoogleChromePath()

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
		if !ignoredHiddenCharacters[l.Value] {
			text.WriteString(l.Value)
		}
	}

	expected := "This is the document title  There is some lede text here  And then another line of text. "
	if got := text.String(); got != expected {
		t.Errorf("Text = %q, want %q", got, expected)
	}
}

// TestSinglePageSimpleGoogleChromeGetsCorrectPageSize verifies page size is Letter.
func TestSinglePageSimpleGoogleChromeGetsCorrectPageSize(t *testing.T) {
	path := getSinglePageSimpleGoogleChromePath()

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

// TestSinglePageSimpleGoogleChromeLettersHavePdfBoxPositions compares letter positions against PdfBox data.
func TestSinglePageSimpleGoogleChromeLettersHavePdfBoxPositions(t *testing.T) {
	path := getSinglePageSimpleGoogleChromePath()

	pdfBoxData := getSinglePageSimpleGoogleChromePdfBoxData(t)
	index := 0

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

	for _, letter := range letters {
		if ignoredHiddenCharacters[letter.Value] || strings.TrimSpace(letter.Value) == "" {
			continue
		}

		if index >= len(pdfBoxData) {
			break
		}

		datum := pdfBoxData[index]

		for ignoredHiddenCharacters[datum.Text] && index < len(pdfBoxData)-1 {
			index++
			datum = pdfBoxData[index]
		}

		if letter.Value != datum.Text {
			t.Errorf("Letter[%d] Text = %q, want %q", index, letter.Value, datum.Text)
		}

		if math.Abs(letter.Location().X-datum.X) > 0.01 {
			t.Errorf("Letter[%d] X = %g, want %g (precision 2)", index, letter.Location().X, datum.X)
		}

		transformed := page.Height() - letter.Location().Y
		if math.Abs(transformed-datum.Y) > 0.01 {
			t.Errorf("Letter[%d] Y = %g, want %g (precision 2)", index, transformed, datum.Y)
		}

		if math.Abs(letter.Width-datum.Width) > 0.01 {
			t.Errorf("Letter[%d] Width = %g, want %g (precision 2)", index, letter.Width, datum.Width)
		}

		if letter.FontName() != datum.FontName {
			t.Errorf("Letter[%d] FontName = %q, want %q", index, letter.FontName(), datum.FontName)
		}

		index++
	}
}

// TestSinglePageSimpleGoogleChromeLettersHaveOtherProviderPositions compares letter positions against other provider data.
func TestSinglePageSimpleGoogleChromeLettersHaveOtherProviderPositions(t *testing.T) {
	path := getSinglePageSimpleGoogleChromePath()

	pdfBoxData := getSinglePageSimpleGoogleChromeOtherData(t)
	index := 0

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

	for _, letter := range letters {
		if ignoredHiddenCharacters[letter.Value] || strings.TrimSpace(letter.Value) == "" {
			continue
		}

		if index >= len(pdfBoxData) {
			break
		}

		datum := pdfBoxData[index]

		for (ignoredHiddenCharacters[datum.Text] || datum.Text == " ") && index < len(pdfBoxData)-1 {
			index++
			datum = pdfBoxData[index]
		}

		if letter.Value != datum.Text {
			t.Errorf("Letter[%d] Text = %q, want %q", index, letter.Value, datum.Text)
		}

		if math.Abs(letter.Location().X-datum.X) > 0.01 {
			t.Errorf("Letter[%d] X = %g, want %g (precision 2)", index, letter.Location().X, datum.X)
		}

		transformed := page.Height() - letter.Location().Y
		if math.Abs(transformed-datum.Y) > 0.01 {
			t.Errorf("Letter[%d] Y = %g, want %g (precision 2)", index, transformed, datum.Y)
		}

		if math.Abs(datum.Width-letter.Width) >= 0.03 {
			t.Errorf("Letter[%d] Width = %g, want within 0.03 of %g", index, letter.Width, datum.Width)
		}

		index++
	}
}

// TestSinglePageSimpleGoogleChromeHandleCorruptedFileOffsets verifies text extraction from corrupted file.
func TestSinglePageSimpleGoogleChromeHandleCorruptedFileOffsets(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "Single Page Broken Offsets - from google drive.pdf")

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

	text := page.Text()

	if text == "" {
		t.Error("expected non-empty text from corrupted file")
	}
}

func getSinglePageSimpleGoogleChromePdfBoxData(t *testing.T) []*testutil.AssertablePositionData {
	t.Helper()

	const fromPdfBox = "72\t105\t9.771912\tT\t21\tArialMT\n" +
		"81.77106\t105\t8.897049\th\t21\tArialMT\n" +
		"90.66733\t105\t3.554138\ti\t21\tArialMT\n" +
		"94.22115\t105\t7.998741\ts\t21\tArialMT\n" +
		"102.2192\t105\t0\t\u200B\t21\tGautami\n" +
		"106.6634\t105\t0\t\u200B\t21\tGautami\n" +
		"106.6634\t105\t3.554131\ti\t21\tArialMT\n" +
		"110.2173\t105\t7.998749\ts\t21\tArialMT\n" +
		"118.2153\t105\t0\t\u200B\t21\tGautami\n" +
		"122.6595\t105\t0\t\u200B\t21\tGautami\n" +
		"122.6595\t105\t4.444618\tt\t21\tArialMT\n" +
		"127.1038\t105\t8.897049\th\t21\tArialMT\n" +
		"136\t105\t8.897049\te\t21\tArialMT\n" +
		"144.8963\t105\t0\t\u200B\t21\tGautami\n" +
		"149.3405\t105\t0\t\u200B\t21\tGautami\n" +
		"149.3405\t105\t8.897049\td\t21\tArialMT\n" +
		"158.2368\t105\t8.897049\to\t21\tArialMT\n" +
		"167.1331\t105\t7.998749\tc\t21\tArialMT\n" +
		"175.1311\t105\t8.897049\tu\t21\tArialMT\n" +
		"184.0274\t105\t13.32605\tm\t21\tArialMT\n" +
		"197.3523\t105\t8.897049\te\t21\tArialMT\n" +
		"206.2485\t105\t8.897049\tn\t21\tArialMT\n" +
		"215.1448\t105\t4.444611\tt\t21\tArialMT\n" +
		"219.5891\t105\t0\t\u200B\t21\tGautami\n" +
		"224.0333\t105\t0\t\u200B\t21\tGautami\n" +
		"224.0333\t105\t4.444611\tt\t21\tArialMT\n" +
		"228.4775\t105\t3.554138\ti\t21\tArialMT\n" +
		"232.0313\t105\t4.444611\tt\t21\tArialMT\n" +
		"236.4756\t105\t3.554123\tl\t21\tArialMT\n" +
		"240.0294\t105\t8.897049\te\t21\tArialMT\n" +
		"72\t143.25\t6.716187\tT\t14\tArialMT\n" +
		"78.71446\t143.25\t6.114899\th\t14\tArialMT\n" +
		"84.8278\t143.25\t6.114891\te\t14\tArialMT\n" +
		"90.94113\t143.25\t3.661423\tr\t14\tArialMT\n" +
		"94.60161\t143.25\t6.114899\te\t14\tArialMT\n" +
		"100.7149\t143.25\t0\t\u200B\t14\tGautami\n" +
		"103.7689\t143.25\t0\t\u200B\t14\tGautami\n" +
		"103.7689\t143.25\t2.442749\ti\t14\tArialMT\n" +
		"106.211\t143.25\t5.497505\ts\t14\tArialMT\n" +
		"111.7071\t143.25\t0\t\u200B\t14\tGautami\n" +
		"114.7611\t143.25\t0\t\u200B\t14\tGautami\n" +
		"114.7611\t143.25\t5.497505\ts\t14\tArialMT\n" +
		"120.2572\t143.25\t6.114899\to\t14\tArialMT\n" +
		"126.3705\t143.25\t9.158928\tm\t14\tArialMT\n" +
		"135.5271\t143.25\t6.114899\te\t14\tArialMT\n" +
		"141.6404\t143.25\t0\t\u200B\t14\tGautami\n" +
		"144.6944\t143.25\t0\t\u200B\t14\tGautami\n" +
		"144.6944\t143.25\t2.442749\tl\t14\tArialMT\n" +
		"147.1365\t143.25\t6.114899\te\t14\tArialMT\n" +
		"153.2499\t143.25\t6.114899\td\t14\tArialMT\n" +
		"159.3632\t143.25\t6.114899\te\t14\tArialMT\n" +
		"165.4765\t143.25\t0\t\u200B\t14\tGautami\n" +
		"168.5305\t143.25\t0\t\u200B\t14\tGautami\n" +
		"168.5305\t143.25\t3.054749\tt\t14\tArialMT\n" +
		"171.5845\t143.25\t6.114899\te\t14\tArialMT\n" +
		"177.6978\t143.25\t5.497498\tx\t14\tArialMT\n" +
		"183.1939\t143.25\t3.054764\tt\t14\tArialMT\n" +
		"186.2479\t143.25\t0\t\u200B\t14\tGautami\n" +
		"189.3019\t143.25\t0\t\u200B\t14\tGautami\n" +
		"189.3019\t143.25\t6.114899\th\t14\tArialMT\n" +
		"195.4152\t143.25\t6.114899\te\t14\tArialMT\n" +
		"201.5285\t143.25\t3.661423\tr\t14\tArialMT\n" +
		"205.189\t143.25\t6.114899\te\t14\tArialMT\n" +
		"72\t173.25\t7.33358\tA\t14\tArialMT\n" +
		"79.3317\t173.25\t6.114891\tn\t14\tArialMT\n" +
		"85.44504\t173.25\t6.114891\td\t14\tArialMT\n" +
		"91.55836\t173.25\t0\t\u200B\t14\tGautami\n" +
		"94.61235\t173.25\t0\t\u200B\t14\tGautami\n" +
		"94.61235\t173.25\t3.054756\tt\t14\tArialMT\n" +
		"97.66633\t173.25\t6.114899\th\t14\tArialMT\n" +
		"103.7797\t173.25\t6.114899\te\t14\tArialMT\n" +
		"109.893\t173.25\t6.114899\tn\t14\tArialMT\n" +
		"116.0063\t173.25\t0\t\u200B\t14\tGautami\n" +
		"119.0603\t173.25\t0\t\u200B\t14\tGautami\n" +
		"119.0603\t173.25\t6.114899\ta\t14\tArialMT\n" +
		"125.1736\t173.25\t6.114899\tn\t14\tArialMT\n" +
		"131.287\t173.25\t6.114899\to\t14\tArialMT\n" +
		"137.4003\t173.25\t3.054749\tt\t14\tArialMT\n" +
		"140.4543\t173.25\t6.114899\th\t14\tArialMT\n" +
		"146.5676\t173.25\t6.114899\te\t14\tArialMT\n" +
		"152.6809\t173.25\t3.661423\tr\t14\tArialMT\n" +
		"156.3414\t173.25\t0\t\u200B\t14\tGautami\n" +
		"159.3954\t173.25\t0\t\u200B\t14\tGautami\n" +
		"159.3954\t173.25\t2.442749\tl\t14\tArialMT\n" +
		"161.8375\t173.25\t2.442734\ti\t14\tArialMT\n" +
		"164.2796\t173.25\t6.114899\tn\t14\tArialMT\n" +
		"170.393\t173.25\t6.114899\te\t14\tArialMT\n" +
		"176.5063\t173.25\t0\t\u200B\t14\tGautami\n" +
		"179.5603\t173.25\t0\t\u200B\t14\tGautami\n" +
		"179.5603\t173.25\t6.114899\to\t14\tArialMT\n" +
		"185.6736\t173.25\t3.054764\tf\t14\tArialMT\n" +
		"188.7276\t173.25\t0\t\u200B\t14\tGautami\n" +
		"191.7816\t173.25\t0\t\u200B\t14\tGautami\n" +
		"191.7816\t173.25\t3.054764\tt\t14\tArialMT\n" +
		"194.8355\t173.25\t6.114899\te\t14\tArialMT\n" +
		"200.9489\t173.25\t5.497482\tx\t14\tArialMT\n" +
		"206.445\t173.25\t3.054764\tt\t14\tArialMT\n" +
		"209.499\t173.25\t3.054764\t.\t14\tArialMT"

	return parseSinglePageSimpleGoogleChromePositionData(t, fromPdfBox)
}

func getSinglePageSimpleGoogleChromeOtherData(t *testing.T) []*testutil.AssertablePositionData {
	t.Helper()

	const fromOther = "72\t105\t9.758476\tT\t0\tArialMT\n" +
		"81.77106\t105\t8.894608\th\t0\tArialMT\n" +
		"90.66733\t105\t3.551445\ti\t0\tArialMT\n" +
		"94.22115\t105\t7.998749\ts\t0\tArialMT\n" +
		"102.2192\t105\t4.431305\t \t0\tArialMT\n" +
		"102.2192\t105\t0\t\u200B\t0\tArialMT\n" +
		"106.6634\t105\t3.551445\ti\t0\tArialMT\n" +
		"106.6634\t105\t0\t\u200B\t0\tArialMT\n" +
		"110.2173\t105\t7.998749\ts\t0\tArialMT\n" +
		"118.2153\t105\t0\t\u200B\t0\tArialMT\n" +
		"118.2153\t105\t4.431305\t \t0\tArialMT\n" +
		"122.6595\t105\t4.431305\tt\t0\tArialMT\n" +
		"122.6595\t105\t0\t\u200B\t0\tArialMT\n" +
		"127.1038\t105\t8.894608\th\t0\tArialMT\n" +
		"136\t105\t8.894608\te\t0\tArialMT\n" +
		"144.8963\t105\t4.431305\t \t0\tArialMT\n" +
		"144.8963\t105\t0\t\u200B\t0\tArialMT\n" +
		"149.3405\t105\t8.894608\td\t0\tArialMT\n" +
		"149.3405\t105\t0\t\u200B\t0\tArialMT\n" +
		"158.2368\t105\t8.894608\to\t0\tArialMT\n" +
		"167.1331\t105\t7.998749\tc\t0\tArialMT\n" +
		"175.1311\t105\t8.894608\tu\t0\tArialMT\n" +
		"184.0274\t105\t13.32591\tm\t0\tArialMT\n" +
		"197.3523\t105\t8.894608\te\t0\tArialMT\n" +
		"206.2485\t105\t8.894608\tn\t0\tArialMT\n" +
		"215.1448\t105\t4.431305\tt\t0\tArialMT\n" +
		"219.5891\t105\t4.431305\t \t0\tArialMT\n" +
		"219.5891\t105\t0\t\u200B\t0\tArialMT\n" +
		"224.0333\t105\t4.431305\tt\t0\tArialMT\n" +
		"224.0333\t105\t0\t\u200B\t0\tArialMT\n" +
		"228.4775\t105\t3.551453\ti\t0\tArialMT\n" +
		"232.0313\t105\t4.431305\tt\t0\tArialMT\n" +
		"236.4756\t105\t3.551453\tl\t0\tArialMT\n" +
		"240.0294\t105\t8.894608\te\t0\tArialMT\n" +
		"248.918\t105\t4.431305\t \t0\tArialMT\n" +
		"72\t128.25\t3.045616\t \t0\tArialMT\n" +
		"72\t143.25\t6.706947\tT\t0\tArialMT\n" +
		"78.71446\t143.25\t6.11322\th\t0\tArialMT\n" +
		"84.8278\t143.25\t6.11322\te\t0\tArialMT\n" +
		"90.94113\t143.25\t3.661331\tr\t0\tArialMT\n" +
		"94.60161\t143.25\t6.11322\te\t0\tArialMT\n" +
		"100.7149\t143.25\t3.045616\t \t0\tArialMT\n" +
		"100.7149\t143.25\t0\t\u200B\t0\tArialMT\n" +
		"103.7689\t143.25\t2.440887\ti\t0\tArialMT\n" +
		"103.7689\t143.25\t0\t\u200B\t0\tArialMT\n" +
		"106.211\t143.25\t5.497498\ts\t0\tArialMT\n" +
		"111.7071\t143.25\t3.045616\t \t0\tArialMT\n" +
		"111.7071\t143.25\t0\t\u200B\t0\tArialMT\n" +
		"114.7611\t143.25\t0\t\u200B\t0\tArialMT\n" +
		"114.7611\t143.25\t5.497498\ts\t0\tArialMT\n" +
		"120.2572\t143.25\t6.11322\to\t0\tArialMT\n" +
		"126.3705\t143.25\t9.158836\tm\t0\tArialMT\n" +
		"135.5271\t143.25\t6.11322\te\t0\tArialMT\n" +
		"141.6404\t143.25\t0\t\u200B\t0\tArialMT\n" +
		"141.6404\t143.25\t3.045609\t \t0\tArialMT\n" +
		"144.6944\t143.25\t2.440887\tl\t0\tArialMT\n" +
		"144.6944\t143.25\t0\t\u200B\t0\tArialMT\n" +
		"147.1365\t143.25\t6.11322\te\t0\tArialMT\n" +
		"153.2499\t143.25\t6.11322\td\t0\tArialMT\n" +
		"159.3632\t143.25\t6.11322\te\t0\tArialMT\n" +
		"165.4765\t143.25\t0\t\u200B\t0\tArialMT\n" +
		"165.4765\t143.25\t3.045609\t \t0\tArialMT\n" +
		"168.5305\t143.25\t3.045609\tt\t0\tArialMT\n" +
		"168.5305\t143.25\t0\t\u200B\t0\tArialMT\n" +
		"171.5845\t143.25\t6.11322\te\t0\tArialMT\n" +
		"177.6978\t143.25\t5.497498\tx\t0\tArialMT\n" +
		"183.1939\t143.25\t3.045609\tt\t0\tArialMT\n" +
		"186.2479\t143.25\t0\t\u200B\t0\tArialMT\n" +
		"186.2479\t143.25\t3.045609\t \t0\tArialMT\n" +
		"189.3019\t143.25\t6.11322\th\t0\tArialMT\n" +
		"189.3019\t143.25\t0\t\u200B\t0\tArialMT\n" +
		"195.4152\t143.25\t6.11322\te\t0\tArialMT\n" +
		"201.5285\t143.25\t3.661331\tr\t0\tArialMT\n" +
		"205.189\t143.25\t6.11322\te\t0\tArialMT\n" +
		"211.3008\t143.25\t3.045609\t \t0\tArialMT\n" +
		"72\t158.25\t3.045616\t \t0\tArialMT\n" +
		"72\t173.25\t7.32267\tA\t0\tArialMT\n" +
		"79.3317\t173.25\t6.11322\tn\t0\tArialMT\n" +
		"85.44504\t173.25\t6.11322\td\t0\tArialMT\n" +
		"91.55836\t173.25\t3.045616\t \t0\tArialMT\n" +
		"91.55836\t173.25\t0\t\u200B\t0\tArialMT\n" +
		"94.61235\t173.25\t0\t\u200B\t0\tArialMT\n" +
		"94.61235\t173.25\t3.045616\tt\t0\tArialMT\n" +
		"97.66633\t173.25\t6.11322\th\t0\tArialMT\n" +
		"103.7797\t173.25\t6.11322\te\t0\tArialMT\n" +
		"109.893\t173.25\t6.11322\tn\t0\tArialMT\n" +
		"116.0063\t173.25\t0\t\u200B\t0\tArialMT\n" +
		"116.0063\t173.25\t3.045616\t \t0\tArialMT\n" +
		"119.0603\t173.25\t6.11322\ta\t0\tArialMT\n" +
		"119.0603\t173.25\t0\t\u200B\t0\tArialMT\n" +
		"125.1736\t173.25\t6.11322\tn\t0\tArialMT\n" +
		"131.287\t173.25\t6.11322\to\t0\tArialMT\n" +
		"137.4003\t173.25\t3.045609\tt\t0\tArialMT\n" +
		"140.4543\t173.25\t6.11322\th\t0\tArialMT\n" +
		"146.5676\t173.25\t6.11322\te\t0\tArialMT\n" +
		"152.6809\t173.25\t3.661331\tr\t0\tArialMT\n" +
		"156.3414\t173.25\t3.045609\t \t0\tArialMT\n" +
		"156.3414\t173.25\t0\t\u200B\t0\tArialMT\n" +
		"159.3954\t173.25\t2.440887\tl\t0\tArialMT\n" +
		"159.3954\t173.25\t0\t\u200B\t0\tArialMT\n" +
		"161.8375\t173.25\t2.440887\ti\t0\tArialMT\n" +
		"164.2796\t173.25\t6.11322\tn\t0\tArialMT\n" +
		"170.393\t173.25\t6.11322\te\t0\tArialMT\n" +
		"176.5063\t173.25\t3.045609\t \t0\tArialMT\n" +
		"176.5063\t173.25\t0\t\u200B\t0\tArialMT\n" +
		"179.5603\t173.25\t6.11322\to\t0\tArialMT\n" +
		"179.5603\t173.25\t0\t\u200B\t0\tArialMT\n" +
		"185.6736\t173.25\t3.045609\tf\t0\tArialMT\n" +
		"188.7276\t173.25\t0\t\u200B\t0\tArialMT\n" +
		"188.7276\t173.25\t3.045609\t \t0\tArialMT\n" +
		"191.7816\t173.25\t3.045609\tt\t0\tArialMT\n" +
		"191.7816\t173.25\t0\t\u200B\t0\tArialMT\n" +
		"194.8355\t173.25\t6.11322\te\t0\tArialMT\n" +
		"200.9489\t173.25\t5.497498\tx\t0\tArialMT\n" +
		"206.445\t173.25\t3.045609\tt\t0\tArialMT\n" +
		"209.499\t173.25\t3.045609\t.\t0\tArialMT\n" +
		"212.543\t173.25\t3.045609\t \t0\tArialMT"

	return parseSinglePageSimpleGoogleChromePositionData(t, fromOther)
}

func parseSinglePageSimpleGoogleChromePositionData(t *testing.T, data string) []*testutil.AssertablePositionData {
	t.Helper()

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
