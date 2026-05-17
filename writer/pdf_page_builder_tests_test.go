package writer_test

import (
	"os"
	"path/filepath"
	"testing"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	standard14fonts "github.com/uglytoad/pdfpig/go/fonts/standard14_fonts"
	"github.com/uglytoad/pdfpig/go/writer"
)

func TestCanAddPng(t *testing.T) {
	var pdfBytes []byte

	pdfBuilder := writer.NewPdfDocumentBuilder()

	page1, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize(page1): %v", err)
	}
	dataPNG := loadPng(t, "1-16bitRGBA-Issue550.png")
	if _, addErr := page1.AddPng(dataPNG, core.NewPdfRectangleFloat(0, 0, 595, 842)); addErr != nil {
		t.Fatalf("page1.AddPng: %v", addErr)
	}

	page2, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize(page2): %v", err)
	}
	dataPNG = loadPng(t, "2-16bitRGB.png")
	if _, addErr := page2.AddPng(dataPNG, core.NewPdfRectangleFloat(0, 0, 595, 842)); addErr != nil {
		t.Fatalf("page2.AddPng: %v", addErr)
	}

	page3, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize(page3): %v", err)
	}
	dataPNG = loadPng(t, "3-16bitGray.png")
	if _, addErr := page3.AddPng(dataPNG, core.NewPdfRectangleFloat(0, 0, 595, 842)); addErr != nil {
		t.Fatalf("page3.AddPng: %v", addErr)
	}

	page4, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize(page4): %v", err)
	}
	dataPNG = loadPng(t, "4-16bitRGBA.png")
	if _, addErr := page4.AddPng(dataPNG, core.NewPdfRectangleFloat(0, 0, 595, 842)); addErr != nil {
		t.Fatalf("page4.AddPng: %v", addErr)
	}

	pdfBytes, err = pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	if writeErr := os.WriteFile("PdfPageBuilderTests_CanAddPng.pdf", pdfBytes, 0o644); writeErr != nil {
		t.Logf("WriteAllBytes: %v", writeErr)
	}

	doc, openErr := pdfpig.Open(pdfBytes, nil)
	if openErr != nil {
		t.Fatalf("Open(pdfBytes): %v", openErr)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 4 {
		t.Errorf("NumberOfPages = %d, want 4", got)
	}

	assertImageOnPage(t, doc, 1, 1170, 2532, 8)
	assertImageOnPage(t, doc, 2, 900, 900, 8)
	assertImageOnPage(t, doc, 3, 900, 900, 8)
	assertImageOnPage(t, doc, 4, 900, 900, 8)
}

func TestCanAddPngTestPattern1(t *testing.T) {
	const subfolderName = "TestPattern1"

	pdfBuilder := writer.NewPdfDocumentBuilder()
	courierFont, err := pdfBuilder.AddStandard14Font(standard14fonts.Courier)
	if err != nil {
		t.Fatalf("AddStandard14Font(Courier): %v", err)
	}

	addPageWithImageCalls := []struct {
		fileName   string
		imageHeight float64
	}{
		{"tp1-001-8bitRGB-withExif~Thumbnail~ColorProfile.png", 150},
		{"tp1-002-8bitRGBA-withExif~Thumbnail~ColorProfile.png", 200},
		{"tp1-003-8bitRGB-Interlaced-withExif~ColorProfile.png", 150},
		{"tp1-004-8bitRGBA-Interlaced-withExif~ColorProfile.png", 200},
		{"tp1-101-16bitRGB-withExif~Thumbnail~ColorProfile.png", 150},
		{"tp1-102-16bitRGBA-withExif~Thumbnail~ColorProfile.png", 200},
		{"tp1-201-32bitRGB-withExif~Thumbnail~ColorProfile.png", 150},
		{"tp1-301-16bitFloatRGB-withExif~Thumbnail~ColorProfile.png", 150},
		{"tp1-401-32bitFloatRGB-withExif~Thumbnail~ColorProfile.png", 150},
	}

	for _, c := range addPageWithImageCalls {
		addPageWithImage(t, pdfBuilder, subfolderName, c.fileName, c.imageHeight, courierFont)
	}

	pdfBytes, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if closeErr := pdfBuilder.Close(); closeErr != nil {
		t.Logf("Close: %v", closeErr)
	}

	if writeErr := os.WriteFile("PdfPageBuilderTests_CanAddPngTestPattern1.pdf", pdfBytes, 0o644); writeErr != nil {
		t.Logf("WriteAllBytes: %v", writeErr)
	}

	doc, openErr := pdfpig.Open(pdfBytes, nil)
	if openErr != nil {
		t.Fatalf("Open(pdfBytes): %v", openErr)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 9 {
		t.Errorf("NumberOfPages = %d, want 9", got)
	}

	assertImageOnPage(t, doc, 1, 200, 150, 8)
	assertImageOnPage(t, doc, 2, 200, 200, 8)
	assertImageOnPage(t, doc, 3, 200, 150, 8)
	assertImageOnPage(t, doc, 4, 200, 200, 8)
	assertImageOnPage(t, doc, 5, 200, 150, 8)
	assertImageOnPage(t, doc, 6, 200, 200, 8)
	assertImageOnPage(t, doc, 7, 200, 150, 8)
	assertImageOnPage(t, doc, 8, 200, 150, 8)
	assertImageOnPage(t, doc, 9, 200, 150, 8)
}

func addPageWithImage(t *testing.T, pdfBuilder *writer.PdfDocumentBuilder, subfolderName, imageFileName string, imageHeight float64, font *writer.AddedFont) {
	t.Helper()

	imageBottom := 842.0 - 600.0
	imagePlacement := core.NewPdfRectangleFloat(0, imageBottom, 595, 842)
	borderPlacement := imagePlacement.BottomLeft
	labelPlacement := core.NewPdfPoint(50, imagePlacement.BottomLeft.Y-50)

	page, err := pdfBuilder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	dataPNG := loadPngSubfolder(t, imageFileName, subfolderName)
	page.DrawRectangle(borderPlacement, imagePlacement.Width, imagePlacement.Height, 3, true)
	if _, addErr := page.AddPng(dataPNG, imagePlacement); addErr != nil {
		t.Fatalf("AddPng(%q): %v", imageFileName, addErr)
	}
	if _, addTextErr := page.AddText(imageFileName, 12, labelPlacement, font); addTextErr != nil {
		t.Logf("AddText(%q): %v", imageFileName, addTextErr)
	}
}

func assertImageOnPage(t *testing.T, doc *content.PdfDocument, pageNum int, wantWidth, wantHeight, wantBpc int) {
	t.Helper()

	pageAny, err := doc.GetPage(pageNum)
	if err != nil {
		t.Errorf("GetPage(%d): %v", pageNum, err)
		return
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Errorf("GetPage(%d): expected *content.Page, got %T", pageNum, pageAny)
		return
	}

	images := page.GetImages()
	if len(images) == 0 {
		t.Errorf("page %d: GetImages returned empty slice", pageNum)
		return
	}

	img := images[0]
	if img.WidthInSamples() != wantWidth {
		t.Errorf("page %d: WidthInSamples = %d, want %d", pageNum, img.WidthInSamples(), wantWidth)
	}
	if img.HeightInSamples() != wantHeight {
		t.Errorf("page %d: HeightInSamples = %d, want %d", pageNum, img.HeightInSamples(), wantHeight)
	}
	if img.BitsPerComponent() != wantBpc {
		t.Errorf("page %d: BitsPerComponent = %d, want %d", pageNum, img.BitsPerComponent(), wantBpc)
	}
}

func loadPng(t *testing.T, name string) []byte {
	t.Helper()
	return loadPngSubfolder(t, name, "")
}

func loadPngSubfolder(t *testing.T, name, subfolderName string) []byte {
	t.Helper()

	pngFilesFolder := filepath.Join("..", "testdata", "Png")
	if subfolderName != "" {
		pngFilesFolder = filepath.Join(pngFilesFolder, subfolderName)
	}

	pngFilePath := filepath.Join(pngFilesFolder, name)
	data, err := os.ReadFile(pngFilePath)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", pngFilePath, err)
	}
	return data
}
