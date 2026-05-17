package png_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	"github.com/uglytoad/pdfpig/go/images/png"
	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/testutil"
)

var (
	rgbBlack  = []byte{0, 0, 0}
	rgbWhite  = []byte{255, 255, 255}
	rgbPalette = [][]byte{rgbBlack, rgbWhite}

	cmykBlack = []byte{0, 0, 0, 255}
	cmykWhite = []byte{0, 0, 0, 0}

	grayscaleBlack byte = 0
	grayscaleWhite byte = 255
)

func flattenBytes(slices [][]byte) []byte {
	total := 0
	for _, s := range slices {
		total += len(s)
	}
	result := make([]byte, total)
	pos := 0
	for _, s := range slices {
		pos += copy(result[pos:], s)
	}
	return result
}

func loadImage(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", name, err)
	}
	return data
}

func TestCanGeneratePngFromDeviceRgbImageData(t *testing.T) {
	testutil.SetFilesFolder("testdata")

	pixels := [][]byte{
		rgbWhite, rgbBlack, rgbWhite,
		rgbBlack, rgbWhite, rgbBlack,
		rgbWhite, rgbBlack, rgbWhite,
	}

	image := testutil.NewTestPdfImage()
	image.SetColorSpaceDetails(&colors.DeviceRgbColorSpaceDetails)
	image.SetDecodedBytes(flattenBytes(pixels))
	image.SetWidthInSamples(3)
	image.SetHeightInSamples(3)

	bytes, ok := png.TryGenerate(image, nil)
	if !ok {
		t.Fatal("TryGenerate returned false for DeviceRGB image")
	}

	equal, err := testutil.ImagesAreEqual(loadImage(t, "3x3.png"), bytes)
	if err != nil {
		t.Fatalf("ImagesAreEqual error: %v", err)
	}
	if !equal {
		t.Error("Generated PNG does not match expected 3x3.png")
	}
}

func TestCanGeneratePngFromDeviceCMYKImageData(t *testing.T) {
	testutil.SetFilesFolder("testdata")

	pixels := [][]byte{
		cmykWhite, cmykBlack, cmykWhite,
		cmykBlack, cmykWhite, cmykBlack,
		cmykWhite, cmykBlack, cmykWhite,
	}

	image := testutil.NewTestPdfImage()
	image.SetColorSpaceDetails(&colors.DeviceCmykColorSpaceDetails)
	image.SetDecodedBytes(flattenBytes(pixels))
	image.SetWidthInSamples(3)
	image.SetHeightInSamples(3)

	bytes, ok := png.TryGenerate(image, nil)
	if !ok {
		t.Fatal("TryGenerate returned false for DeviceCMYK image")
	}

	equal, err := testutil.ImagesAreEqual(loadImage(t, "3x3.png"), bytes)
	if err != nil {
		t.Fatalf("ImagesAreEqual error: %v", err)
	}
	if !equal {
		t.Error("Generated PNG does not match expected 3x3.png")
	}
}

func TestCanGeneratePngFromDeviceGrayscaleImageData(t *testing.T) {
	testutil.SetFilesFolder("testdata")

	pixels := []byte{
		grayscaleWhite, grayscaleBlack, grayscaleWhite,
		grayscaleBlack, grayscaleWhite, grayscaleBlack,
		grayscaleWhite, grayscaleBlack, grayscaleWhite,
	}

	image := testutil.NewTestPdfImage()
	image.SetColorSpaceDetails(&colors.DeviceGrayColorSpaceDetails)
	image.SetDecodedBytes(pixels)
	image.SetWidthInSamples(3)
	image.SetHeightInSamples(3)

	bytes, ok := png.TryGenerate(image, nil)
	if !ok {
		t.Fatal("TryGenerate returned false for DeviceGray image")
	}

	equal, err := testutil.ImagesAreEqual(loadImage(t, "3x3.png"), bytes)
	if err != nil {
		t.Fatalf("ImagesAreEqual error: %v", err)
	}
	if !equal {
		t.Error("Generated PNG does not match expected 3x3.png")
	}
}

func TestCanGeneratePngFromIndexedImageData8bpc(t *testing.T) {
	testutil.SetFilesFolder("testdata")

	indices := []byte{
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
	}

	indexedCS := colors.NewIndexedColorSpaceDetails(colors.DeviceRgbColorSpaceDetails, 1, flattenBytes(rgbPalette))
	image := testutil.NewTestPdfImage()
	image.SetColorSpaceDetails(&indexedCS)
	image.SetDecodedBytes(indices)
	image.SetWidthInSamples(3)
	image.SetHeightInSamples(3)

	bytes, ok := png.TryGenerate(image, nil)
	if !ok {
		t.Fatal("TryGenerate returned false for indexed image (8 bpc)")
	}

	equal, err := testutil.ImagesAreEqual(loadImage(t, "3x3.png"), bytes)
	if err != nil {
		t.Fatalf("ImagesAreEqual error: %v", err)
	}
	if !equal {
		t.Error("Generated PNG does not match expected 3x3.png")
	}
}

func TestCanExtractPngFromPdfWithIndexedImageData8bpc(t *testing.T) {
	pdfPath := filepath.Join("testdata", "indexed-png-with-mask.pdf")

	doc, err := pdfpig.OpenFile(pdfPath, nil)
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
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	imgs := page.GetImages()
	if len(imgs) == 0 {
		t.Fatal("No images found on page")
	}

	img := imgs[0]
	bytes, ok := img.TryGetPng()
	if !ok {
		t.Fatal("TryGetPng returned false")
	}

	outputPath := filepath.Join("testdata", "indexed-png-with-mask-extracted.png")
	if err := os.WriteFile(outputPath, bytes, 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func TestCanGeneratePngFromIndexedImageData1bpc(t *testing.T) {
	testutil.SetFilesFolder("testdata")

	lines := make([]byte, 3)
	lines[0] |= (1 << 7)
	lines[0] |= (1 << 5)
	lines[1] |= (1 << 6)
	lines[2] |= (1 << 7)
	lines[2] |= (1 << 5)

	colorTable := flattenBytes(rgbPalette)
	indexedCS := colors.NewIndexedColorSpaceDetails(colors.DeviceRgbColorSpaceDetails, 1, colorTable)
	image := testutil.NewTestPdfImage()
	image.SetColorSpaceDetails(&indexedCS)
	image.SetDecodedBytes(lines)
	image.SetWidthInSamples(3)
	image.SetHeightInSamples(3)
	image.SetBitsPerComponent(1)

	bytes, ok := png.TryGenerate(image, nil)
	if !ok {
		t.Fatal("TryGenerate returned false for indexed image (1 bpc)")
	}

	equal, err := testutil.ImagesAreEqual(loadImage(t, "3x3.png"), bytes)
	if err != nil {
		t.Fatalf("ImagesAreEqual error: %v", err)
	}
	if !equal {
		t.Error("Generated PNG does not match expected 3x3.png")
	}
}

func TestCanGeneratePngFromCcittFaxDecodedImageData(t *testing.T) {
	testutil.SetFilesFolder("testdata")

	decodedBytes, err := testutil.LoadFileBytes("ccittfax-decoded.bin", false)
	if err != nil {
		t.Fatalf("LoadFileBytes(ccittfax-decoded.bin): %v", err)
	}

	stencilCS := colors.StencilIndexedColorSpace(colors.DeviceGrayColorSpaceDetails)
	image := testutil.NewTestPdfImage()
	image.SetColorSpaceDetails(&stencilCS)
	image.SetDecodedBytes(decodedBytes)
	image.SetWidthInSamples(1800)
	image.SetHeightInSamples(3113)
	image.SetBitsPerComponent(1)
	image.SetDecode([]float64{1.0, 0})

	bytes, ok := png.TryGenerate(image, nil)
	if !ok {
		t.Fatal("TryGenerate returned false for CCITT fax decoded image")
	}

	equal, err := testutil.ImagesAreEqual(loadImage(t, "ccittfax.png"), bytes)
	if err != nil {
		t.Fatalf("ImagesAreEqual error: %v", err)
	}
	if !equal {
		t.Error("Generated PNG does not match expected ccittfax.png")
	}
}

func TestCanGeneratePngFromICCBasedImageData(t *testing.T) {
	testutil.SetFilesFolder("testdata")

	decodedBytes, err := testutil.LoadFileBytes("iccbased-decoded.bin", false)
	if err != nil {
		t.Fatalf("LoadFileBytes(iccbased-decoded.bin): %v", err)
	}

	iccCS, err := colors.NewICCBasedColorSpaceDetails(3, colors.DeviceRgbColorSpaceDetails, []float64{0, 1, 0, 1, 0, 1})
	if err != nil {
		t.Fatalf("NewICCBasedColorSpaceDetails: %v", err)
	}

	image := testutil.NewTestPdfImage()
	image.SetColorSpaceDetails(&iccCS)
	image.SetDecodedBytes(decodedBytes)
	image.SetWidthInSamples(1)
	image.SetHeightInSamples(1)
	image.SetBitsPerComponent(8)

	bytes, ok := png.TryGenerate(image, nil)
	if !ok {
		t.Fatal("TryGenerate returned false for ICCBased image")
	}

	equal, err := testutil.ImagesAreEqual(loadImage(t, "iccbased.png"), bytes)
	if err != nil {
		t.Fatalf("ImagesAreEqual error: %v", err)
	}
	if !equal {
		t.Error("Generated PNG does not match expected iccbased.png")
	}
}

func TestAlternateColorSpaceDetailsIsCurrentlyUsedInPdfPigWhenGeneratingPngsFromICCBasedImageData(t *testing.T) {
	testutil.SetFilesFolder("testdata")

	decodedBytes, err := testutil.LoadFileBytes("iccbased-decoded.bin", false)
	if err != nil {
		t.Fatalf("LoadFileBytes(iccbased-decoded.bin): %v", err)
	}

	iccCS, err := colors.NewICCBasedColorSpaceDetails(3, colors.DeviceRgbColorSpaceDetails, []float64{0, 1, 0, 1, 0, 1})
	if err != nil {
		t.Fatalf("NewICCBasedColorSpaceDetails: %v", err)
	}

	iccImage := testutil.NewTestPdfImage()
	iccImage.SetColorSpaceDetails(&iccCS)
	iccImage.SetDecodedBytes(decodedBytes)
	iccImage.SetWidthInSamples(1)
	iccImage.SetHeightInSamples(1)
	iccImage.SetBitsPerComponent(8)

	deviceRGBImage := testutil.NewTestPdfImage()
	deviceRGBImage.SetColorSpaceDetails(&colors.DeviceRgbColorSpaceDetails)
	deviceRGBImage.SetDecodedBytes(decodedBytes)
	deviceRGBImage.SetWidthInSamples(1)
	deviceRGBImage.SetHeightInSamples(1)
	deviceRGBImage.SetBitsPerComponent(8)

	iccPngBytes, ok := png.TryGenerate(iccImage, nil)
	if !ok {
		t.Fatal("TryGenerate returned false for ICCBased image")
	}

	deviceRgbBytes, ok := png.TryGenerate(deviceRGBImage, nil)
	if !ok {
		t.Fatal("TryGenerate returned false for DeviceRGB image")
	}

	if !equalBytes(iccPngBytes, deviceRgbBytes) {
		t.Error("ICCBased PNG bytes do not match DeviceRGB PNG bytes")
	}
}

func TestCanGeneratePngFromCalRGBImageData(t *testing.T) {
	testutil.SetFilesFolder("testdata")

	decodedBytes, err := testutil.LoadFileBytes("calrgb-decoded.bin", false)
	if err != nil {
		t.Fatalf("LoadFileBytes(calrgb-decoded.bin): %v", err)
	}

	calRGBCS, err := colors.NewCalRGBColorSpaceDetails(
		[]float64{0.95043, 1, 1.09},
		nil,
		[]float64{2.2, 2.2, 2.2},
		[]float64{
			0.41239, 0.21264, 0.01933,
			0.35758, 0.71517, 0.11919,
			0.18045, 0.07218, 0.9504,
		},
	)
	if err != nil {
		t.Fatalf("NewCalRGBColorSpaceDetails: %v", err)
	}

	image := testutil.NewTestPdfImage()
	image.SetColorSpaceDetails(&calRGBCS)
	image.SetDecodedBytes(decodedBytes)
	image.SetWidthInSamples(153)
	image.SetHeightInSamples(83)
	image.SetBitsPerComponent(8)

	bytes, ok := png.TryGenerate(image, nil)
	if !ok {
		t.Fatal("TryGenerate returned false for CalRGB image")
	}

	equal, err := testutil.ImagesAreEqual(loadImage(t, "calrgb.png"), bytes)
	if err != nil {
		t.Fatalf("ImagesAreEqual error: %v", err)
	}
	if !equal {
		t.Error("Generated PNG does not match expected calrgb.png")
	}
}

func TestCanGeneratePngFromCalGrayImageData(t *testing.T) {
	testutil.SetFilesFolder("testdata")

	decodedBytes, err := testutil.LoadFileBytes("calgray-decoded.bin", true)
	if err != nil {
		t.Fatalf("LoadFileBytes(calgray-decoded.bin): %v", err)
	}

	calGrayCS, err := colors.NewCalGrayColorSpaceDetails(
		[]float64{0.9505000114, 1, 1.0889999866},
		nil,
		2.2000000477,
	)
	if err != nil {
		t.Fatalf("NewCalGrayColorSpaceDetails: %v", err)
	}

	image := testutil.NewTestPdfImage()
	image.SetColorSpaceDetails(&calGrayCS)
	image.SetDecodedBytes(decodedBytes)
	image.SetWidthInSamples(2480)
	image.SetHeightInSamples(1748)
	image.SetBitsPerComponent(8)

	bytes, ok := png.TryGenerate(image, nil)
	if !ok {
		t.Fatal("TryGenerate returned false for CalGray image")
	}

	equal, err := testutil.ImagesAreEqual(loadImage(t, "calgray.png"), bytes)
	if err != nil {
		t.Fatalf("ImagesAreEqual error: %v", err)
	}
	if !equal {
		t.Error("Generated PNG does not match expected calgray.png")
	}
}

func equalBytes(a, b []byte) bool {
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
