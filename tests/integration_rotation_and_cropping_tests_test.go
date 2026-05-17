//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"math"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

func sparcV9Path() string {
	return filepath.Join(integrationDocRoot, "SPARC - v9 Architecture Manual.pdf")
}

func publicationOfAwardPath() string {
	return filepath.Join(integrationDocRoot, "Publication_of_award_of_Bids_for_Transport_Sector__August_2016.pdf")
}

func smallCropboxPath() string {
	return filepath.Join(integrationDocRoot, "SmallCropbox.pdf")
}

// TestCroppedPageHasCorrectTextCoordinates verifies that a cropped page has the
// expected dimensions and text coordinates. Matches C# CroppedPageHasCorrectTextCoordinates.
func TestCroppedPageHasCorrectTextCoordinates(t *testing.T) {
	path := sparcV9Path()

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

	// Due to cropping
	if page.Width() != 612 {
		t.Errorf("page.Width() = %.0f, want 612", page.Width())
	}
	if page.Height() != 792 {
		t.Errorf("page.Height() = %.0f, want 792", page.Height())
	}

	letters := page.Letters()
	if len(letters) == 0 {
		t.Fatal("page has no letters")
	}

	minX := letters[0].BoundingBox.Left()
	maxX := letters[0].BoundingBox.Right()
	for _, l := range letters[1:] {
		if l.BoundingBox.Left() < minX {
			minX = l.BoundingBox.Left()
		}
		if l.BoundingBox.Right() > maxX {
			maxX = l.BoundingBox.Right()
		}
	}

	// If cropping is not applied correctly, these values will be off
	// C# uses Assert.Equal(74, minX, 0) which rounds to integer precision (precision=0).
	if math.Round(minX) != 74 {
		t.Errorf("minX = %.4f, want 74", minX)
	}
	if math.Round(maxX) != 540 {
		t.Errorf("maxX = %.4f, want 540", maxX)
	}

	// page.Content is not exposed as a public getter in Go; the fact that Letters() returns non-nil
	// implicitly confirms content was loaded. The C# Assert.NotNull(page.Content) is satisfied by
	// the len(letters) > 0 check above.
}

// TestWrongPathCount verifies path count and page dimensions with ClipPaths enabled.
// Matches C# WrongPathCount.
func TestWrongPathCount(t *testing.T) {
	path := publicationOfAwardPath()

	opts := &content.ParsingOptions{
		ClipPaths: true,
	}

	doc, err := pdfpig.OpenFile(path, opts)
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

	if page.Height() != 612 {
		t.Errorf("page.Height() = %.0f, want 612", page.Height())
	}

	paths := page.Paths()
	if len(paths) != 224 {
		t.Errorf("len(page.Paths()) = %d, want 224", len(paths))
	}
}

// TestIssue665 verifies rotation, dimensions, crop box and media box for SmallCropbox.pdf.
// Matches C# Issue665.
func TestIssue665(t *testing.T) {
	path := smallCropboxPath()

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

	// Clockwise rotation
	if page.Rotation().Value != 270 {
		t.Errorf("page.Rotation().Value = %d, want 270", page.Rotation().Value)
	}

	if int(page.Height()) != 680 {
		t.Errorf("int(page.Height()) = %d, want 680", int(page.Height()))
	}
	if int(page.Width()) != 433 {
		t.Errorf("int(page.Width()) = %d, want 433", int(page.Width()))
	}

	if page.Size() != content.PageSizeCustom {
		t.Errorf("page.Size() = %v, want PageSize.Custom (%v)", page.Size(), content.PageSizeCustom)
	}

	if len(page.Letters()) != 2429 {
		t.Errorf("len(page.Letters()) = %d, want 2429", len(page.Letters()))
	}

	cropBox := page.CropBox().Bounds
	if cropBox.Rotation() != 0 {
		t.Errorf("cropBox.Rotation() = %.0f, want 0", cropBox.Rotation())
	}
	if int(cropBox.Height) != 680 {
		t.Errorf("int(cropBox.Height) = %d, want 680", int(cropBox.Height))
	}
	if int(cropBox.Width) != 433 {
		t.Errorf("int(cropBox.Width) = %d, want 433", int(cropBox.Width))
	}
	if int(cropBox.Bottom()) != 0 {
		t.Errorf("int(cropBox.Bottom()) = %d, want 0", int(cropBox.Bottom()))
	}
	if int(cropBox.Left()) != 0 {
		t.Errorf("int(cropBox.Left()) = %d, want 0", int(cropBox.Left()))
	}
	if int(cropBox.Right()) != 433 {
		t.Errorf("int(cropBox.Right()) = %d, want 433", int(cropBox.Right()))
	}
	if int(cropBox.Top()) != 680 {
		t.Errorf("int(cropBox.Top()) = %d, want 680", int(cropBox.Top()))
	}

	mediaBox := page.MediaBox().Bounds
	if mediaBox.Rotation() != 0 {
		t.Errorf("mediaBox.Rotation() = %.0f, want 0", mediaBox.Rotation())
	}
	if int(mediaBox.Height) != 680 {
		t.Errorf("int(mediaBox.Height) = %d, want 680", int(mediaBox.Height))
	}
	if int(mediaBox.Width) != 433 {
		t.Errorf("int(mediaBox.Width) = %d, want 433", int(mediaBox.Width))
	}
	if int(mediaBox.Bottom()) != 0 {
		t.Errorf("int(mediaBox.Bottom()) = %d, want 0", int(mediaBox.Bottom()))
	}
	if int(mediaBox.Left()) != 0 {
		t.Errorf("int(mediaBox.Left()) = %d, want 0", int(mediaBox.Left()))
	}
	if int(mediaBox.Right()) != 433 {
		t.Errorf("int(mediaBox.Right()) = %d, want 433", int(mediaBox.Right()))
	}
	if int(mediaBox.Top()) != 680 {
		t.Errorf("int(mediaBox.Top()) = %d, want 680", int(mediaBox.Top()))
	}
}
