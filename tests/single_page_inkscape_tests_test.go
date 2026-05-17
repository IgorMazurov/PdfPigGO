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

const singlePageInkscapeDocName = "Single Page Simple - from inkscape.pdf"

func getSinglePageInkscapePath() string {
	return filepath.Join(integrationDocRoot, singlePageInkscapeDocName)
}

// TestSinglePageInkscapeLettersHaveCorrectPositionsPdfBox verifies letter positions match PDFBox reference data.
func TestSinglePageInkscapeLettersHaveCorrectPositionsPdfBox(t *testing.T) {
	path := getSinglePageInkscapePath()

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
	positions := getSinglePageInkscapePdfBoxData(t)

	if len(letters) != len(positions) {
		t.Fatalf("Letter count mismatch: got %d, want %d", len(letters), len(positions))
	}

	for i := range letters {
		positions[i].AssertWithinTolerance(t, letters[i], false)
	}
}

func getSinglePageInkscapePdfBoxData(t *testing.T) []*testutil.AssertablePositionData {
	t.Helper()

	const data = "100.57143\t687.4286\t31.616001\tW\t32.0\tKTICVV+DejaVuSans\t47.776\n" +
		"130.74742\t687.4286\t13.152\tr\t32.0\tKTICVV+DejaVuSans\t36.704002\n" +
		"143.89941\t687.4286\t8.864\ti\t32.0\tKTICVV+DejaVuSans\t49.792004\n" +
		"152.76341\t687.4286\t12.544001\tt\t32.0\tKTICVV+DejaVuSans\t46.016003\n" +
		"165.30742\t687.4286\t19.68\te\t32.0\tKTICVV+DejaVuSans\t37.632\n" +
		"184.98743\t687.4286\t10.144\t \t32.0\tKTICVV+DejaVuSans\t0.0\n" +
		"195.13142\t687.4286\t16.640001\ts\t32.0\tKTICVV+DejaVuSans\t37.632\n" +
		"211.77142\t687.4286\t19.552\to\t32.0\tKTICVV+DejaVuSans\t37.632\n" +
		"231.32343\t687.4286\t31.168001\tm\t32.0\tKTICVV+DejaVuSans\t36.704002\n" +
		"262.49142\t687.4286\t19.68\te\t32.0\tKTICVV+DejaVuSans\t37.632\n" +
		"282.17142\t687.4286\t12.544001\tt\t32.0\tKTICVV+DejaVuSans\t46.016003\n" +
		"294.71542\t687.4286\t20.256\th\t32.0\tKTICVV+DejaVuSans\t49.792004\n" +
		"315.06744\t687.4286\t8.864\ti\t32.0\tKTICVV+DejaVuSans\t49.792004\n" +
		"323.93146\t687.4286\t20.256\tn\t32.0\tKTICVV+DejaVuSans\t36.704002\n" +
		"344.18747\t687.4286\t20.288\tg\t32.0\tKTICVV+DejaVuSans\t50.336002\n" +
		"364.47546\t687.4286\t10.144\t \t32.0\tKTICVV+DejaVuSans\t0.0\n" +
		"374.61948\t687.4286\t8.864\ti\t32.0\tKTICVV+DejaVuSans\t49.792004\n" +
		"383.4835\t687.4286\t20.256\tn\t32.0\tKTICVV+DejaVuSans\t36.704002\n" +
		"100.57143\t647.4286\t9.408\tI\t32.0\tKTICVV+DejaVuSans\t47.776\n" +
		"109.97942\t647.4286\t20.256\tn\t32.0\tKTICVV+DejaVuSans\t36.704002\n" +
		"130.23543\t647.4286\t18.528002\tk\t32.0\tKTICVV+DejaVuSans\t49.792004\n" +
		"148.76343\t647.4286\t16.640001\ts\t32.0\tKTICVV+DejaVuSans\t37.632\n" +
		"165.40343\t647.4286\t17.568\tc\t32.0\tKTICVV+DejaVuSans\t37.632\n" +
		"183.06743\t647.4286\t19.584002\ta\t32.0\tKTICVV+DejaVuSans\t37.632\n" +
		"202.65143\t647.4286\t20.288\tp\t32.0\tKTICVV+DejaVuSans\t50.336002\n" +
		"222.93942\t647.4286\t19.68\te\t32.0\tKTICVV+DejaVuSans\t37.632\n"

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
