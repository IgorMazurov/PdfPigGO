package reading_order_detector_test

import (
	"sort"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	dla "github.com/uglytoad/pdfpig/go/document_layout_analysis"
	rod "github.com/uglytoad/pdfpig/go/document_layout_analysis/reading_order_detector"
)

func createFakeTextBlock(boundingBox core.PdfRectangle) *dla.TextBlock {
	letter := content.NewLetterWithDetails(
		"a",
		boundingBox,
		boundingBox,
		boundingBox.BottomLeft,
		boundingBox.BottomRight,
		10, 1, nil, core.NeitherFillNorStrokeButClip, nil, nil, 0, 0)

	word, _ := content.NewWord([]*content.Letter{letter})
	line, _ := dla.NewTextLine([]*content.Word{word}, " ")
	block, _ := dla.NewTextBlock([]*dla.TextLine{line}, "\n")

	return block
}

func TestReadingOrderOrdersItemsOnTheSameRowContents(t *testing.T) {
	leftTextBlock := createFakeTextBlock(core.NewPdfRectangle(
		core.NewPdfPoint(0, 0), core.NewPdfPoint(10, 10)))
	rightTextBlock := createFakeTextBlock(core.NewPdfRectangle(
		core.NewPdfPoint(100, 0), core.NewPdfPoint(110, 10)))

	textBlocks := []*dla.TextBlock{rightTextBlock, leftTextBlock}

	detector := rod.NewUnsupervisedReadingOrderDetector(5, rod.RowWise, true)
	orderedBlocks := detector.Get(textBlocks)

	ordered := make([]*dla.TextBlock, len(orderedBlocks))
	copy(ordered, orderedBlocks)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].ReadingOrder < ordered[j].ReadingOrder
	})

	if got := ordered[0].BoundingBox.Left(); got != 0 {
		t.Errorf("expected first block Left == 0, got %g", got)
	}
	if got := ordered[1].BoundingBox.Left(); got != 100 {
		t.Errorf("expected second block Left == 100, got %g", got)
	}

	_ = orderedBlocks
}

func TestDocumentTest(t *testing.T) {
	title := createFakeTextBlock(core.NewPdfRectangle(
		core.NewPdfPoint(42.6, 709.06), core.NewPdfPoint(42.6, 709.06)))
	line1Left := createFakeTextBlock(core.NewPdfRectangle(
		core.NewPdfPoint(42.6, 668.86), core.NewPdfPoint(42.6, 668.86)))
	line1Right := createFakeTextBlock(core.NewPdfRectangle(
		core.NewPdfPoint(302.21, 668.86), core.NewPdfPoint(302.21, 668.86)))
	line2Left := createFakeTextBlock(core.NewPdfRectangle(
		core.NewPdfPoint(42.6, 608.26), core.NewPdfPoint(42.6, 608.26)))
	line2TallerRight := createFakeTextBlock(core.NewPdfRectangle(
		core.NewPdfPoint(302.21, 581.35), core.NewPdfPoint(302.21, 581.35)))
	line3 := createFakeTextBlock(core.NewPdfRectangle(
		core.NewPdfPoint(42.6, 515.83), core.NewPdfPoint(42.6, 515.83)))
	line4Left := createFakeTextBlock(core.NewPdfRectangle(
		core.NewPdfPoint(42.6, 490.27), core.NewPdfPoint(42.6, 490.27)))
	line4Right := createFakeTextBlock(core.NewPdfRectangle(
		core.NewPdfPoint(302.21, 491.59), core.NewPdfPoint(302.21, 491.59)))

	textBlocks := []*dla.TextBlock{
		title, line4Left, line2TallerRight, line4Right,
		line1Right, line1Left, line3, line2Left,
	}

	detector := rod.NewUnsupervisedReadingOrderDetector(5, rod.RowWise, true)
	orderedBlocks := detector.Get(textBlocks)

	ordered := make([]*dla.TextBlock, len(orderedBlocks))
	copy(ordered, orderedBlocks)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].ReadingOrder < ordered[j].ReadingOrder
	})

	expected := []*dla.TextBlock{
		title, line1Left, line1Right, line2Left,
		line2TallerRight, line3, line4Left, line4Right,
	}

	for i, exp := range expected {
		if !ordered[i].BoundingBox.Equals(exp.BoundingBox) {
			t.Errorf("index %d: expected BoundingBox %v, got %v",
				i, exp.BoundingBox, ordered[i].BoundingBox)
		}
	}

	_ = orderedBlocks
}
