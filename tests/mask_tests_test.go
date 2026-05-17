//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"bytes"
	"image/color"
	"image/png"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func TestPigProductionHandbook(t *testing.T) {
	path := testutil.GetDocumentPath("Pig Production Handbook.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true, SkipMissingFonts: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page := pageAny.(*content.Page)

	images := page.GetImages()

	checkMaskedImageAlpha(t, images, 1, "page1 image[1]")
	checkMaskedImageAlpha(t, images, 2, "page1 image[2]")
}

func TestMozillaLink3264_0(t *testing.T) {
	path := testutil.GetDocumentPath("MOZILLA-LINK-3264-0.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true, SkipMissingFonts: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page1 := page1Any.(*content.Page)

	images1 := page1.GetImages()
	checkMaskedImageAlpha(t, images1, 1, "page1 image[1]")

	page2Any, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	page2 := page2Any.(*content.Page)

	images2 := page2.GetImages()
	checkMaskedImageHasMask(t, images2, 1, "page2 image[1]")
}

func checkMaskedImageAlpha(t *testing.T, images []content.PdfImage, idx int, label string) {
	t.Helper()

	if len(images) <= idx {
		t.Fatalf("%s: expected at least %d images, got %d", label, idx+1, len(images))
	}

	img := images[idx]

	maskImg := img.MaskImage()
	if maskImg == nil {
		t.Errorf("%s: MaskImage is nil", label)
		return
	}

	pngBytes, ok := img.TryGetPng()
	if !ok {
		t.Errorf("%s: TryGetPng returned false", label)
		return
	}

	imgData, err := png.Decode(bytes.NewReader(pngBytes))
	if err != nil {
		t.Errorf("%s: png.Decode: %v", label, err)
		return
	}

	rgba := color.NRGBAModel.Convert(imgData.At(0, 0)).(color.NRGBA)
	alpha := rgba.A
	if alpha != 0 {
		t.Errorf("%s: expected pixel (0,0) alpha == 0, got %d", label, alpha)
	}
}

func checkMaskedImageHasMask(t *testing.T, images []content.PdfImage, idx int, label string) {
	t.Helper()

	if len(images) <= idx {
		t.Fatalf("%s: expected at least %d images, got %d", label, idx+1, len(images))
	}

	img := images[idx]

	maskImg := img.MaskImage()
	if maskImg == nil {
		t.Errorf("%s: MaskImage is nil", label)
		return
	}

	pngBytes, ok := img.TryGetPng()
	if !ok {
		t.Errorf("%s: TryGetPng returned false", label)
		return
	}
	_ = pngBytes
}
