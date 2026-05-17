package type1_test

import (
	"os"
	"path/filepath"
	"testing"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func init() {
	testutil.IntegrationDocumentsRoot = "../../testdata/integration/Documents"
}

func TestCorrectBoundingBoxesFlexPoints(t *testing.T) {
	pointComparer := testutil.NewPointComparer(testutil.NewDoubleComparer(0.001))

	filePath := filepath.Join(testutil.IntegrationDocumentsRoot, "data.pdf")

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("could not read test file %q: %v", filePath, err)
	}

	doc, openErr := pdfpig.Open(data, &content.ParsingOptions{})
	if openErr != nil {
		t.Fatalf("Open(data): %v", openErr)
	}
	defer doc.Close()

	pageAny, pageErr := doc.GetPage(1)
	if pageErr != nil {
		t.Fatalf("GetPage(1): %v", pageErr)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	letters := page.Letters()

	if len(letters) < 2 {
		t.Fatalf("expected at least 2 letters, got %d", len(letters))
	}

	m := letters[0]
	if m.Value != "m" {
		t.Errorf("letters[0].Value = %q, want %q", m.Value, "m")
	}

	expectedMBottomLeft := core.NewPdfPoint(253.4458, 658.431)
	if !pointComparer.Equals(m.BoundingBox.BottomLeft, expectedMBottomLeft) {
		t.Errorf("letters[0].BoundingBox.BottomLeft = %v, want %v", m.BoundingBox.BottomLeft, expectedMBottomLeft)
	}

	expectedMTopRight := core.NewPdfPoint(261.22659, 662.83446)
	if !pointComparer.Equals(m.BoundingBox.TopRight, expectedMTopRight) {
		t.Errorf("letters[0].BoundingBox.TopRight = %v, want %v", m.BoundingBox.TopRight, expectedMTopRight)
	}

	p := letters[1]
	if p.Value != "p" {
		t.Errorf("letters[1].Value = %q, want %q", p.Value, "p")
	}

	expectedPBottomLeft := core.NewPdfPoint(261.70778, 656.49825)
	if !pointComparer.Equals(p.BoundingBox.BottomLeft, expectedPBottomLeft) {
		t.Errorf("letters[1].BoundingBox.BottomLeft = %v, want %v", p.BoundingBox.BottomLeft, expectedPBottomLeft)
	}

	expectedPTopRight := core.NewPdfPoint(266.6193, 662.83446)
	if !pointComparer.Equals(p.BoundingBox.TopRight, expectedPTopRight) {
		t.Errorf("letters[1].BoundingBox.TopRight = %v, want %v", p.BoundingBox.TopRight, expectedPTopRight)
	}
}
