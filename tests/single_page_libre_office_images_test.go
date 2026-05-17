//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"math"
	"sort"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func getSinglePageLibreOfficeImagesPath() string {
	return testutil.GetDocumentPath("Single Page Images - from libre office.pdf", true)
}

func TestSinglePageLibreOfficeImagesHas3Images(t *testing.T) {
	filePath := getSinglePageLibreOfficeImagesPath()

	doc, err := pdfpig.OpenFile(filePath, &content.ParsingOptions{UseLenientParsing: false})
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
		t.Fatal("expected *content.Page from GetPage")
	}

	images := page.GetImages()

	if got := len(images); got != 3 {
		t.Errorf("len(images) = %d, want 3", got)
	}
}

func TestSinglePageLibreOfficeImagesHaveCorrectDimensionsAndLocations(t *testing.T) {
	const tolerance = 0.1

	filePath := getSinglePageLibreOfficeImagesPath()

	doc, err := pdfpig.OpenFile(filePath, &content.ParsingOptions{UseLenientParsing: false})
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
		t.Fatal("expected *content.Page from GetPage")
	}

	images := page.GetImages()

	sort.Slice(images, func(i, j int) bool {
		return images[i].BoundingBox().Width < images[j].BoundingBox().Width
	})

	pdfPigSquare := images[0]
	bb0 := pdfPigSquare.BoundingBox()

	assertNear(t, bb0.Width, 148.3, tolerance, "pdfPigSquare BoundingBox.Width")
	assertNear(t, bb0.Height, 148.3, tolerance, "pdfPigSquare BoundingBox.Height")
	assertNear(t, bb0.Left(), 60.1, tolerance, "pdfPigSquare BoundingBox.Left")
	assertNear(t, bb0.Top(), 765.8, tolerance, "pdfPigSquare BoundingBox.Top")

	pdfPigSquished := images[1]
	bb1 := pdfPigSquished.BoundingBox()

	assertNear(t, bb1.Width, 206.8, tolerance, "pdfPigSquished BoundingBox.Width")
	assertNear(t, bb1.Height, 83.2, tolerance, "pdfPigSquished BoundingBox.Height")
	assertNear(t, bb1.Left(), 309.8, tolerance, "pdfPigSquished BoundingBox.Left")
	assertNear(t, bb1.Top(), 552.1, tolerance, "pdfPigSquished BoundingBox.Top")

	birthdayPigs := images[2]
	bb2 := birthdayPigs.BoundingBox()

	assertNear(t, bb2.Width, 391.0, tolerance, "birthdayPigs BoundingBox.Width")
	assertNear(t, bb2.Height, 267.1, tolerance, "birthdayPigs BoundingBox.Height")
	assertNear(t, bb2.Left(), 102.2, tolerance, "birthdayPigs BoundingBox.Left")
	assertNear(t, bb2.Top(), 426.3, tolerance, "birthdayPigs BoundingBox.Top")
}

func TestSinglePageLibreOfficeImagesHasCorrectText(t *testing.T) {
	filePath := getSinglePageLibreOfficeImagesPath()

	doc, err := pdfpig.OpenFile(filePath, &content.ParsingOptions{UseLenientParsing: false})
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
		t.Fatal("expected *content.Page from GetPage")
	}

	expectedText := "Oink oink"
	if got := page.Text(); got != expectedText {
		t.Errorf("page.Text() = %q, want %q", got, expectedText)
	}
}

func TestSinglePageLibreOfficeImagesCanAccessImageBytes(t *testing.T) {
	filePath := getSinglePageLibreOfficeImagesPath()

	doc, err := pdfpig.OpenFile(filePath, &content.ParsingOptions{UseLenientParsing: false})
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
		t.Fatal("expected *content.Page from GetPage")
	}

	for i, image := range page.GetImages() {
		bytes, ok := image.TryGetBytes()
		if ok {
			if len(bytes) == 0 {
				t.Errorf("image[%d]: TryGetBytes succeeded but returned empty slice", i)
			}
		} else {
			if len(image.RawMemory()) == 0 {
				t.Errorf("image[%d]: TryGetBytes failed and RawMemory is also empty", i)
			}
		}
	}
}

func assertNear(t *testing.T, got, want, tolerance float64, msg string) {
	t.Helper()
	if math.Abs(got-want) > tolerance {
		t.Errorf("%s: got %v, want %v (diff=%v, tolerance=%v)", msg, got, want, math.Abs(got-want), tolerance)
	}
}
