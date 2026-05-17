//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/word_extractor"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	"github.com/uglytoad/pdfpig/go/testutil"
)

const colorSpaceOutputFolder = "ColorSpaceTests"

func initColorSpaceOutputFolder(t *testing.T) {
	t.Helper()
	if err := os.MkdirAll(colorSpaceOutputFolder, 0755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", colorSpaceOutputFolder, err)
	}
}

func assertEqual[T comparable](t *testing.T, expected, got T, msg string) {
	t.Helper()
	if got != expected {
		t.Errorf("%s: expected %v, got %v", msg, expected, got)
	}
}

func assertTrue(t *testing.T, cond bool, msg string) {
	t.Helper()
	if !cond {
		t.Errorf("%s: expected true", msg)
	}
}

func assertNotNull[T comparable](t *testing.T, val T, msg string) {
	t.Helper()
	var zero T
	if val == zero {
		t.Fatalf("%s: expected non-nil", msg)
	}
}

func assertAlmostEqual(t *testing.T, expected, got float64, precision int, msg string) {
	t.Helper()
	diff := math.Abs(expected - got)
	tolerance := math.Pow(10, -float64(precision))
	if diff > tolerance {
		t.Errorf("%s: expected %v, got %v (diff=%v)", msg, expected, got, diff)
	}
}

func convertToByte(componentValue float64) byte {
	r := math.Round(componentValue * 255)
	if r < 0 {
		r = 0
	}
	if r > 255 {
		r = 255
	}
	return byte(r)
}

// TestIndexedDeviceNColorSpaceImages verifies that images with Indexed/DeviceN color spaces
// can be decoded to PNG. Matches C# ColorSpaceTests.IndexedDeviceNColorSpaceImages.
func TestIndexedDeviceNColorSpaceImages(t *testing.T) {
	initColorSpaceOutputFolder(t)

	path := testutil.GetDocumentPath("MOZILLA-3136-0.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	assertNoError(t, err, "GetPage(1)")
	page1 := page1Any.(*content.Page)

	images1 := page1.GetImages()
	if len(images1) < 14 {
		t.Fatalf("expected at least 14 images on page 1, got %d", len(images1))
	}

	image12 := images1[12]
	cs12 := image12.ColorSpaceDetails()
	assertNotNull(t, cs12, "image12 ColorSpaceDetails")
	assertEqual(t, colors.Indexed, (*cs12).Type(), "image12 Type == Indexed")

	indexedCs12, ok := (*cs12).(colors.IndexedColorSpaceAccessor)
	assertTrue(t, ok, "image12 should be IndexedColorSpaceAccessor")
	assertEqual(t, colors.DeviceN, indexedCs12.BaseColorSpace().Type(), "image12 BaseColorSpace.Type == DeviceN")

	bytes1_12, ok := image12.TryGetPng()
	assertTrue(t, ok, "image12 TryGetPng should succeed (Cyan square)")
	if ok {
		outPath := filepath.Join(colorSpaceOutputFolder, "MOZILLA-3136-0_1_12.png")
		writeFile(t, outPath, bytes1_12)
	}

	image13 := images1[13]
	cs13 := image13.ColorSpaceDetails()
	assertNotNull(t, cs13, "image13 ColorSpaceDetails")
	assertEqual(t, colors.Indexed, (*cs13).Type(), "image13 Type == Indexed")

	indexedCs13, ok := (*cs13).(colors.IndexedColorSpaceAccessor)
	assertTrue(t, ok, "image13 should be IndexedColorSpaceAccessor")
	assertEqual(t, colors.DeviceN, indexedCs13.BaseColorSpace().Type(), "image13 BaseColorSpace.Type == DeviceN")

	bytes1_13, ok := image13.TryGetPng()
	assertTrue(t, ok, "image13 TryGetPng should succeed (Cyan square)")
	if ok {
		outPath := filepath.Join(colorSpaceOutputFolder, "MOZILLA-3136-0_1_13.png")
		writeFile(t, outPath, bytes1_13)
	}
}

// TestBitsPerComponents16 verifies that 16-bit per component images can be decoded.
// Matches C# ColorSpaceTests.BitsPerComponents16.
func TestBitsPerComponents16(t *testing.T) {
	initColorSpaceOutputFolder(t)

	path := testutil.GetDocumentPath("MOZILLA-3136-0.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(3)
	assertNoError(t, err, "GetPage(3)")
	page1 := page1Any.(*content.Page)

	images1 := page1.GetImages()
	if len(images1) < 10 {
		t.Fatalf("expected at least 10 images on page 3, got %d", len(images1))
	}

	image9 := images1[9]
	assertEqual(t, 16, image9.BitsPerComponent(), "image9 BitsPerComponent == 16")

	bytes_3_9, ok := image9.TryGetPng()
	assertTrue(t, ok, "image9 TryGetPng should succeed")
	if ok {
		outPath := filepath.Join(colorSpaceOutputFolder, "MOZILLA-3136-0_3_9_16bits.png")
		writeFile(t, outPath, bytes_3_9)
	}
}

// TestDeviceNColorSpaceImages verifies that DeviceN color space images with ICC alternate
// can be decoded. Matches C# ColorSpaceTests.DeviceNColorSpaceImages.
func TestDeviceNColorSpaceImages(t *testing.T) {
	initColorSpaceOutputFolder(t)

	path := testutil.GetDocumentPath("DeviceN_CS_test.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	page3Any, err := doc.GetPage(3)
	assertNoError(t, err, "GetPage(3)")
	page3 := page3Any.(*content.Page)

	images3 := page3.GetImages()
	assertDeviceNImage(t, images3[0], "DeviceN_CS_test_3_0")
	assertDeviceNImage(t, images3[2], "DeviceN_CS_test_3_2")

	page6Any, err := doc.GetPage(6)
	assertNoError(t, err, "GetPage(6)")
	page6 := page6Any.(*content.Page)

	images6 := page6.GetImages()
	assertDeviceNImage(t, images6[0], "DeviceN_CS_test_6_0")
	assertDeviceNImage(t, images6[1], "DeviceN_CS_test_6_1")
	assertDeviceNImage(t, images6[2], "DeviceN_CS_test_6_2")
}

func assertDeviceNImage(t *testing.T, img content.PdfImage, name string) {
	t.Helper()
	cs := img.ColorSpaceDetails()
	assertNotNull(t, cs, name+" ColorSpaceDetails")

	deviceNCs, ok := (*cs).(colors.DeviceNColorSpaceAccessor)
	assertTrue(t, ok, name+" should be DeviceNColorSpaceAccessor")

	altCS := deviceNCs.AlternateColorSpace()
	if altCS == nil {
		t.Fatalf("%s: AlternateColorSpace is nil", name)
	}
	assertEqual(t, colors.ICCBased, altCS.Type(), name+" AlternateColorSpace.Type == ICCBased")

	bytes, ok := img.TryGetPng()
	assertTrue(t, ok, name+" TryGetPng should succeed")
	if ok {
		outPath := filepath.Join(colorSpaceOutputFolder, name+".png")
		writeFile(t, outPath, bytes)
	}
}

// TestSeparationColorSpaceImages verifies that Separation color space images with
// DeviceCMYK alternate can be decoded. Matches C# ColorSpaceTests.SeparationColorSpaceImages.
func TestSeparationColorSpaceImages(t *testing.T) {
	initColorSpaceOutputFolder(t)

	path := testutil.GetDocumentPath("MOZILLA-7375-0.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	assertNoError(t, err, "GetPage(1)")
	page1 := page1Any.(*content.Page)

	images := page1.GetImages()
	cs0 := images[0].ColorSpaceDetails()
	assertNotNull(t, cs0, "image 0 ColorSpaceDetails")

	separationCs, ok := (*cs0).(colors.SeparationColorSpaceAccessor)
	assertTrue(t, ok, "image 0 should be SeparationColorSpaceAccessor")

	altCS := separationCs.AlternateColorSpace()
	if altCS == nil {
		t.Fatal("Separation AlternateColorSpace is nil")
	}
	assertEqual(t, colors.DeviceCMYK, altCS.Type(), "AlternateColorSpace.Type == DeviceCMYK")

	for i, img := range images {
		png, ok := img.TryGetPng()
		if ok {
			outPath := filepath.Join(colorSpaceOutputFolder, "MOZILLA-7375-0_1_"+string(rune('0'+i))+".png")
			writeFile(t, outPath, png)
		}
	}
}

// TestSeparationColorSpace verifies Separation and Indexed/Separation color spaces.
// Matches C# ColorSpaceTests.SeparationColorSpace.
func TestSeparationColorSpace(t *testing.T) {
	initColorSpaceOutputFolder(t)

	path := testutil.GetDocumentPath("MOZILLA-3136-0.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(4)
	assertNoError(t, err, "GetPage(4)")
	page := pageAny.(*content.Page)

	images := page.GetImages()
	if len(images) < 10 {
		t.Fatalf("expected at least 10 images on page 4, got %d", len(images))
	}

	image4 := images[4]
	cs4 := image4.ColorSpaceDetails()
	assertNotNull(t, cs4, "image4 ColorSpaceDetails")

	separation, ok := (*cs4).(colors.SeparationColorSpaceAccessor)
	assertTrue(t, ok, "image4 should be SeparationColorSpaceAccessor")
	_ = separation

	png4, ok := image4.TryGetPng()
	assertTrue(t, ok, "image4 TryGetPng should succeed")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-3136-0_4_separation.png"), png4)
	}

	image9 := images[9]
	cs9 := image9.ColorSpaceDetails()
	assertNotNull(t, cs9, "image9 ColorSpaceDetails")

	indexedCs, ok := (*cs9).(colors.IndexedColorSpaceAccessor)
	assertTrue(t, ok, "image9 should be IndexedColorSpaceAccessor")
	assertEqual(t, colors.Separation, indexedCs.BaseColorSpace().Type(), "image9 BaseColorSpace.Type == Separation")

	png9, ok := image9.TryGetPng()
	assertTrue(t, ok, "image9 TryGetPng should succeed (Green dolphin)")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-3136-0_9_separation.png"), png9)
	}
}

// TestIndexedCalRgbColorSpaceImages verifies Indexed/CalRGB color space images.
// Matches C# ColorSpaceTests.IndexedCalRgbColorSpaceImages.
func TestIndexedCalRgbColorSpaceImages(t *testing.T) {
	initColorSpaceOutputFolder(t)

	path := testutil.GetDocumentPath("MOZILLA-10084-0.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	assertNoError(t, err, "GetPage(1)")
	page1 := page1Any.(*content.Page)

	images1 := page1.GetImages()

	image0 := images1[0]
	cs0 := image0.ColorSpaceDetails()
	assertNotNull(t, cs0, "image0 ColorSpaceDetails")
	assertEqual(t, colors.Indexed, (*cs0).Type(), "image0 Type == Indexed")

	indexedCs0, ok := (*cs0).(colors.IndexedColorSpaceAccessor)
	assertTrue(t, ok, "image0 should be IndexedColorSpaceAccessor")
	assertEqual(t, colors.CalRGB, indexedCs0.BaseColorSpace().Type(), "image0 BaseColorSpace.Type == CalRGB")

	bytes0, ok := image0.TryGetPng()
	assertTrue(t, ok, "image0 TryGetPng should succeed")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-10084-0_1_0.png"), bytes0)
	}

	image1 := images1[1]
	cs1 := image1.ColorSpaceDetails()
	assertNotNull(t, cs1, "image1 ColorSpaceDetails")
	assertEqual(t, colors.Indexed, (*cs1).Type(), "image1 Type == Indexed")

	indexedCs1, ok := (*cs1).(colors.IndexedColorSpaceAccessor)
	assertTrue(t, ok, "image1 should be IndexedColorSpaceAccessor")
	assertEqual(t, colors.CalRGB, indexedCs1.BaseColorSpace().Type(), "image1 BaseColorSpace.Type == CalRGB")

	bytes1, ok := image1.TryGetPng()
	assertTrue(t, ok, "image1 TryGetPng should succeed")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-10084-0_1_1.png"), bytes1)
	}
}

// TestStencilIndexedIccColorSpaceImages verifies stencil and Indexed/ICCBased images.
// Matches C# ColorSpaceTests.StencilIndexedIccColorSpaceImages.
func TestStencilIndexedIccColorSpaceImages(t *testing.T) {
	initColorSpaceOutputFolder(t)

	path := testutil.GetDocumentPath("MOZILLA-10225-0.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	page2Any, err := doc.GetPage(2)
	assertNoError(t, err, "GetPage(2)")
	page2 := page2Any.(*content.Page)

	images2 := page2.GetImages()

	image2_0 := images2[0]
	cs2_0 := image2_0.ColorSpaceDetails()
	assertNotNull(t, cs2_0, "image 2/0 ColorSpaceDetails")
	assertEqual(t, colors.Indexed, (*cs2_0).Type(), "image 2/0 Type == Indexed (ICCBased)")
	bytes2_0, ok := image2_0.TryGetPng()
	assertTrue(t, ok, "image 2/0 TryGetPng should succeed")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-10225-0_1_0.png"), bytes2_0)
	}

	image2_1 := images2[1]
	cs2_1 := image2_1.ColorSpaceDetails()
	assertNotNull(t, cs2_1, "image 2/1 ColorSpaceDetails")
	assertEqual(t, colors.Indexed, (*cs2_1).Type(), "image 2/1 Type == Indexed (stencil)")
	bytes2_1, ok := image2_1.TryGetPng()
	assertTrue(t, ok, "image 2/1 TryGetPng should succeed")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-10225-0_1_1.png"), bytes2_1)
	}

	page23Any, err := doc.GetPage(23)
	assertNoError(t, err, "GetPage(23)")
	page23 := page23Any.(*content.Page)

	images23 := page23.GetImages()
	image23_0 := images23[0]
	cs23_0 := image23_0.ColorSpaceDetails()
	assertNotNull(t, cs23_0, "image 23/0 ColorSpaceDetails")
	assertEqual(t, colors.Indexed, (*cs23_0).Type(), "image 23/0 Type == Indexed")
	bytes23_0, ok := image23_0.TryGetPng()
	assertTrue(t, ok, "image 23/0 TryGetPng should succeed")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-10225-0_23_0.png"), bytes23_0)
	}

	page332Any, err := doc.GetPage(332)
	assertNoError(t, err, "GetPage(332)")
	page332 := page332Any.(*content.Page)

	images332 := page332.GetImages()
	image332_0 := images332[0]
	cs332_0 := image332_0.ColorSpaceDetails()
	assertNotNull(t, cs332_0, "image 332/0 ColorSpaceDetails")
	assertEqual(t, colors.ICCBased, (*cs332_0).Type(), "image 332/0 Type == ICCBased")
	bytes332_0, ok := image332_0.TryGetPng()
	assertTrue(t, ok, "image 332/0 TryGetPng should succeed")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-10225-0_332_0.png"), bytes332_0)
	}

	page338Any, err := doc.GetPage(338)
	assertNoError(t, err, "GetPage(338)")
	page338 := page338Any.(*content.Page)

	images338 := page338.GetImages()
	image338_1 := images338[1]
	cs338_1 := image338_1.ColorSpaceDetails()
	assertNotNull(t, cs338_1, "image 338/1 ColorSpaceDetails")
	assertEqual(t, colors.Indexed, (*cs338_1).Type(), "image 338/1 Type == Indexed")
	bytes338_1, ok := image338_1.TryGetPng()
	assertTrue(t, ok, "image 338/1 TryGetPng should succeed")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-10225-0_338_1.png"), bytes338_1)
	}

	page339Any, err := doc.GetPage(339)
	assertNoError(t, err, "GetPage(339)")
	page339 := page339Any.(*content.Page)

	images339 := page339.GetImages()
	image339_0 := images339[0]
	cs339_0 := image339_0.ColorSpaceDetails()
	assertNotNull(t, cs339_0, "image 339/0 ColorSpaceDetails")
	assertEqual(t, colors.Indexed, (*cs339_0).Type(), "image 339/0 Type == Indexed")
	bytes339_0, ok := image339_0.TryGetPng()
	assertTrue(t, ok, "image 339/0 TryGetPng should succeed")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-10225-0_339_0.png"), bytes339_0)
	}

	image339_1 := images339[1]
	cs339_1 := image339_1.ColorSpaceDetails()
	assertNotNull(t, cs339_1, "image 339/1 ColorSpaceDetails")
	assertEqual(t, colors.Indexed, (*cs339_1).Type(), "image 339/1 Type == Indexed")
	bytes339_1, ok := image339_1.TryGetPng()
	assertTrue(t, ok, "image 339/1 TryGetPng should succeed")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-10225-0_339_1.png"), bytes339_1)
	}

	page341Any, err := doc.GetPage(341)
	assertNoError(t, err, "GetPage(341)")
	page341 := page341Any.(*content.Page)

	images341 := page341.GetImages()
	image341_0 := images341[0]
	cs341_0 := image341_0.ColorSpaceDetails()
	assertNotNull(t, cs341_0, "image 341/0 ColorSpaceDetails")
	assertEqual(t, colors.Indexed, (*cs341_0).Type(), "image 341/0 Type == Indexed")
	bytes341_0, ok := image341_0.TryGetPng()
	assertTrue(t, ok, "image 341/0 TryGetPng should succeed")
	if ok {
		writeFile(t, filepath.Join(colorSpaceOutputFolder, "MOZILLA-10225-0_341_0.png"), bytes341_0)
	}
}

// TestSeparationLabColorSpace verifies Separation/Lab color space rendering.
// Matches C# ColorSpaceTests.SeparationLabColorSpace.
func TestSeparationLabColorSpace(t *testing.T) {
	initColorSpaceOutputFolder(t)

	path := testutil.GetDocumentPath("TIKA-1552-0.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	assertNoError(t, err, "GetPage(1)")
	page1 := page1Any.(*content.Page)

	paths := page1.Paths()
	if len(paths) == 0 {
		t.Fatal("expected at least one path on page 1")
	}

	background := paths[0]
	assertTrue(t, background.IsFilled(), "background path should be filled")

	fillColor := background.FillColor()
	rgbValues := fillColor.ToRGBValues()

	rByte := convertToByte(rgbValues.R)
	gByte := convertToByte(rgbValues.G)
	bByte := convertToByte(rgbValues.B)

	assertEqual(t, byte(10), rByte, "R should be ~10 (Pantone 289 C)")
	assertEqual(t, byte(34), gByte, "G should be 34")
	assertEqual(t, byte(64), bByte, "B should be 64")
}

// TestCanGetAllPagesImages verifies that all images on every page can be extracted.
// Matches C# ColorSpaceTests.CanGetAllPagesImages.
func TestCanGetAllPagesImages(t *testing.T) {
	initColorSpaceOutputFolder(t)

	path := testutil.GetDocumentPath("Pig Production Handbook.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for p := 0; p < doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p + 1)
		assertNoError(t, err, "GetPage("+string(rune('0'+p+1))+")")
		page := pageAny.(*content.Page)

		images := page.GetImages()
		for i, img := range images {
			png, ok := img.TryGetPng()
			if ok {
				outPath := filepath.Join(colorSpaceOutputFolder, "Pig_Production_Handbook_"+string(rune('0'+p+1))+"_"+string(rune('0'+i))+".png")
				writeFile(t, outPath, png)
			}
		}
	}
}

// TestSeparationIccColorSpacesWithForm verifies Separation/ICC color spaces with forms.
// Matches C# ColorSpaceTests.SeparationIccColorSpacesWithForm.
func TestSeparationIccColorSpacesWithForm(t *testing.T) {
	initColorSpaceOutputFolder(t)

	path := testutil.GetDocumentPath("68-1990-01_A.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	assertNoError(t, err, "GetPage(1)")
	page1 := page1Any.(*content.Page)

	paths1 := page1.Paths()
	var filledPaths1 []content.PdfPath
	for _, p := range paths1 {
		if p.IsFilled() {
			filledPaths1 = append(filledPaths1, p)
		}
	}
	if len(filledPaths1) == 0 {
		t.Fatal("expected at least one filled path on page 1")
	}

	reflexRed := filledPaths1[0].FillColor().ToRGBValues()
	assertAlmostEqual(t, 0.930496, reflexRed.R, 6, "Reflex Red R component")
	assertAlmostEqual(t, 0.111542, reflexRed.G, 6, "Reflex Red G component")
	assertAlmostEqual(t, 0.142197, reflexRed.B, 6, "Reflex Red B component")

	page2Any, err := doc.GetPage(2)
	assertNoError(t, err, "GetPage(2)")
	page2 := page2Any.(*content.Page)

	words := page2.GetWordsWithExtractor(word_extractor.DefaultInstance)
	var urlWord *content.Word
	for _, w := range words {
		if strings.Contains(w.Text, "www.extron.com") {
			urlWord = w
			break
		}
	}
	assertNotNull(t, urlWord, "should find word containing www.extron.com")

	firstLetter := urlWord.Letters[0]
	assertEqual(t, "w", firstLetter.Value, "first letter should be 'w'")
	letterRGB := firstLetter.Color.ToRGBValues()
	assertAlmostEqual(t, 0, letterRGB.R, 0, "first letter R == 0 (Blue)")
	assertAlmostEqual(t, 0, letterRGB.G, 0, "first letter G == 0 (Blue)")
	assertAlmostEqual(t, 1, letterRGB.B, 0, "first letter B == 1 (Blue)")

	paths2 := page2.Paths()
	var filledPaths2 []content.PdfPath
	for _, p := range paths2 {
		if p.IsFilled() {
			filledPaths2 = append(filledPaths2, p)
		}
	}

	var filledRects []content.PdfPath
	for _, p := range filledPaths2 {
		if p.Len() == 1 && p.Subpaths()[0].IsDrawnAsRectangle() {
			filledRects = append(filledRects, p)
		}
	}

	sortFilledRectsByPosition(filledRects)

	lightRed := colors.RGBValues{R: 0.985, G: 0.942, B: 0.921}
	lightRed2 := colors.RGBValues{R: 1, G: 0.95, B: 0.95}
	lightOrange := colors.RGBValues{R: 0.993, G: 0.964, B: 0.929}

	for i, rect := range filledRects {
		colorVal := rect.FillColor()
		assertEqual(t, colors.DeviceRGB, colorVal.ColorSpace(), "filled rect "+string(rune('0'+i))+" ColorSpace == DeviceRGB")

		rgbVals := colorVal.ToRGBValues()

		if i%2 == 0 {
			if i == 2 {
				assertAlmostEqual(t, lightRed2.R, rgbVals.R, 3, "rect "+string(rune('0'+i))+" R (lightRed2)")
				assertAlmostEqual(t, lightRed2.G, rgbVals.G, 3, "rect "+string(rune('0'+i))+" G (lightRed2)")
				assertAlmostEqual(t, lightRed2.B, rgbVals.B, 3, "rect "+string(rune('0'+i))+" B (lightRed2)")
			} else {
				assertAlmostEqual(t, lightRed.R, rgbVals.R, 3, "rect "+string(rune('0'+i))+" R (lightRed)")
				assertAlmostEqual(t, lightRed.G, rgbVals.G, 3, "rect "+string(rune('0'+i))+" G (lightRed)")
				assertAlmostEqual(t, lightRed.B, rgbVals.B, 3, "rect "+string(rune('0'+i))+" B (lightRed)")
			}
		} else {
			assertAlmostEqual(t, lightOrange.R, rgbVals.R, 3, "rect "+string(rune('0'+i))+" R (lightOrange)")
			assertAlmostEqual(t, lightOrange.G, rgbVals.G, 3, "rect "+string(rune('0'+i))+" G (lightOrange)")
			assertAlmostEqual(t, lightOrange.B, rgbVals.B, 3, "rect "+string(rune('0'+i))+" B (lightOrange)")
		}
	}
}

func sortFilledRectsByPosition(rects []content.PdfPath) {
	for i := 0; i < len(rects); i++ {
		for j := i + 1; j < len(rects); j++ {
			bbI := rects[i].GetBoundingRectangle()
			bbJ := rects[j].GetBoundingRectangle()
			if bbI == nil || bbJ == nil {
				continue
			}

			if bbI.Left() > bbJ.Left() {
				rects[i], rects[j] = rects[j], rects[i]
			} else if bbI.Left() == bbJ.Left() {
				if bbI.Top() < bbJ.Top() {
					rects[i], rects[j] = rects[j], rects[i]
				}
			}
		}
	}
}

// TestIssue724 verifies that a specific PDF with Type 4 function can be opened without error.
// Matches C# ColorSpaceTests.Issue724.
func TestIssue724(t *testing.T) {
	path := testutil.GetDocumentPath("11194059_2017-11_de_s.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	assertNoError(t, err, "GetPage(1)")
	page1 := page1Any.(*content.Page)
	assertNotNull(t, page1, "page 1 should not be nil")

	page2Any, err := doc.GetPage(2)
	assertNoError(t, err, "GetPage(2)")
	page2 := page2Any.(*content.Page)
	assertNotNull(t, page2, "page 2 should not be nil")
}

func assertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Errorf("WriteFile(%q): %v", path, err)
	}
}
