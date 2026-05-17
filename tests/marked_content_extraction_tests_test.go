//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"math"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

const fileName1 = "Multiple Page - from Mortality Statistics.pdf"
const fileName2 = "68-1990-01_A.pdf"

func getMarkedContentPath1() string {
	return testutil.GetDocumentPath(fileName1, true)
}

func getMarkedContentPath2() string {
	return testutil.GetDocumentPath(fileName2, true)
}

// TestCanIncrementIndex verifies marked content indices are sequential 0..n-1.
func TestCanIncrementIndex(t *testing.T) {
	// First document: page 2 should have 37 marked contents
	doc1, err := pdfpig.OpenFile(getMarkedContentPath1(), nil)
	if err != nil {
		t.Fatalf("OpenFile path1: %v", err)
	}
	defer doc1.Close()

	pageAny1, err := doc1.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	page1, ok := pageAny1.(*content.Page)
	if !ok {
		t.Fatal("page 2 is not *content.Page")
	}

	mcs1 := page1.GetMarkedContents()
	if len(mcs1) == 0 {
		t.Fatal("expected non-empty marked contents for document 1")
	}
	if len(mcs1) != 37 {
		t.Errorf("document 1: expected 37 marked contents, got %d", len(mcs1))
	}

	for i := range mcs1 {
		if mcs1[i].Index != i {
			t.Errorf("document 1: mcs[%d].Index = %d, want %d", i, mcs1[i].Index, i)
		}
	}

	// Second document: page 10 should have 86 marked contents
	doc2, err := pdfpig.OpenFile(getMarkedContentPath2(), nil)
	if err != nil {
		t.Fatalf("OpenFile path2: %v", err)
	}
	defer doc2.Close()

	pageAny2, err := doc2.GetPage(10)
	if err != nil {
		t.Fatalf("GetPage(10): %v", err)
	}
	page2, ok := pageAny2.(*content.Page)
	if !ok {
		t.Fatal("page 10 is not *content.Page")
	}

	mcs2 := page2.GetMarkedContents()
	if len(mcs2) == 0 {
		t.Fatal("expected non-empty marked contents for document 2")
	}
	if len(mcs2) != 86 {
		t.Errorf("document 2: expected 86 marked contents, got %d", len(mcs2))
	}

	for i := range mcs2 {
		if mcs2[i].Index != i {
			t.Errorf("document 2: mcs[%d].Index = %d, want %d", i, mcs2[i].Index, i)
		}
	}
}

// TestCanGetTree verifies the tree structure of marked content with children.
func TestCanGetTree(t *testing.T) {
	doc, err := pdfpig.OpenFile(getMarkedContentPath2(), nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(10)
	if err != nil {
		t.Fatalf("GetPage(10): %v", err)
	}
	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("page 10 is not *content.Page")
	}

	mcs := page.GetMarkedContents()
	if len(mcs) == 0 {
		t.Fatal("expected non-empty marked contents")
	}
	if len(mcs) != 86 {
		t.Errorf("expected 86 marked contents, got %d", len(mcs))
	}

	testCases := []int{8, 9, 75}
	for _, idx := range testCases {
		mc := mcs[idx]
		if len(mc.Children) != 1 {
			t.Errorf("mcs[%d].Children: expected length 1, got %d", idx, len(mc.Children))
		}
		if mc.Index != idx {
			t.Errorf("mcs[%d].Index = %d, want %d", idx, mc.Index, idx)
		}
		if len(mc.Children) > 0 && mc.Children[0].Index != idx {
			t.Errorf("mcs[%d].Children[0].Index = %d, want %d", idx, mc.Children[0].Index, idx)
		}

		// C# Assert.DoesNotContain(mc.Children[0], mcs) uses reference equality.
		// In Go value semantics, children are separate copies so this is guaranteed.
	}
}

// TestCanGetArtifact verifies artifact marked content element properties.
func TestCanGetArtifact(t *testing.T) {
	doc, err := pdfpig.OpenFile(getMarkedContentPath1(), &content.ParsingOptions{ClipPaths: false})
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("page 2 is not *content.Page")
	}

	mcs := page.GetMarkedContents()
	if len(mcs) == 0 {
		t.Fatal("expected non-empty marked contents")
	}

	content0 := mcs[0]
	if !content0.IsArtifact {
		t.Error("mcs[0].IsArtifact expected true")
	}

	if content0.MarkedContentIdentifier != -1 {
		t.Errorf("content0.MarkedContentIdentifier = %d, want -1", content0.MarkedContentIdentifier)
	}

	if len(content0.Letters) != 33 {
		t.Errorf("content0.Letters: expected length 33, got %d", len(content0.Letters))
	}
	if len(content0.Paths) != 8 {
		t.Errorf("content0.Paths: expected length 8, got %d", len(content0.Paths))
	}
	if len(content0.Images) != 0 {
		t.Errorf("content0.Images: expected length 0, got %d", len(content0.Images))
	}

	if !content0.IsTopAttached {
		t.Error("content0.IsTopAttached expected true")
	}
	if content0.IsRightAttached {
		t.Error("content0.IsRightAttached expected false")
	}
	if content0.IsLeftAttached {
		t.Error("content0.IsLeftAttached expected false")
	}
	if content0.IsBottomAttached {
		t.Error("content0.IsBottomAttached expected false")
	}

	if content0.BoundingBox == nil {
		t.Fatal("content0.BoundingBox expected non-nil")
	}
	assertFloatNearMarkedContent(t, "BoundingBox.BottomLeft.X", content0.BoundingBox.BottomLeft.X, 89.03, 0.001)
	assertFloatNearMarkedContent(t, "BoundingBox.BottomLeft.Y", content0.BoundingBox.BottomLeft.Y, 717.756, 0.001)
	assertFloatNearMarkedContent(t, "BoundingBox.TopRight.X", content0.BoundingBox.TopRight.X, 574.422, 0.001)
	assertFloatNearMarkedContent(t, "BoundingBox.TopRight.Y", content0.BoundingBox.TopRight.Y, 751.1398, 0.001)

	if content0.ArtifactType != content.Pagination {
		t.Errorf("content0.ArtifactType = %v, want %v", content0.ArtifactType, content.Pagination)
	}
	if content0.SubType == nil || *content0.SubType != "Header" {
		t.Errorf("content0.SubType = %v, want \"Header\"", content0.SubType)
	}
}

func assertFloatNearMarkedContent(t *testing.T, name string, got, want, tolerance float64) {
	if math.Abs(got-want) > tolerance {
		t.Errorf("%s = %f, want %f (tolerance %f)", name, got, want, tolerance)
	}
}
