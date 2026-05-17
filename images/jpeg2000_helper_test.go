package images_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/images"
	pdfpig "github.com/uglytoad/pdfpig/go"
)

var jp2TestRoot = "testdata/jp2"

func getJp2Files(t *testing.T) []string {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(jp2TestRoot, "*.jp2"))
	if err != nil {
		t.Fatalf("Glob: %v", err)
	}

	names := make([]string, 0, len(files))
	for _, f := range files {
		names = append(names, filepath.Base(f))
	}
	return names
}

func TestGetBitsPerComponent_ReturnsCorrectBitsPerComponent_WhenValidInput(t *testing.T) {
	jp2Files := getJp2Files(t)
	if len(jp2Files) == 0 {
		t.Fatal("no .jp2 test files found")
	}

	for _, name := range jp2Files {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(jp2TestRoot, name))
			if err != nil {
				t.Fatalf("ReadFile(%q): %v", name, err)
			}

			bpc, err := images.GetBitsPerComponent(data)
			if err != nil {
				t.Fatalf("GetBitsPerComponent: %v", err)
			}

			if bpc != 8 {
				t.Errorf("GetBitsPerComponent = %d, want 8", bpc)
			}
		})
	}
}

func TestGetBitsPerComponent_ReturnsError_WhenInputIsTooShort(t *testing.T) {
	_, err := images.GetBitsPerComponent(make([]byte, 11))
	if err == nil {
		t.Error("expected error for input shorter than 12 bytes")
	}
}

func TestGetBitsPerComponent_ReturnsError_WhenSignatureBoxIsInvalid(t *testing.T) {
	_, err := images.GetBitsPerComponent(make([]byte, 12))
	if err == nil {
		t.Error("expected error for invalid JP2 signature")
	}
}

func TestGetBitsPerComponentJ2K(t *testing.T) {
	path := filepath.Join("..", "testdata", "integration", "SpecificTestDocuments", "GHOSTSCRIPT-688999-2.pdf")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", path, err)
	}

	doc, err := pdfpig.Open(data, nil)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	imgs := page.GetImages()
	if len(imgs) != 1 {
		t.Fatalf("expected 1 image on page 1, got %d", len(imgs))
	}

	bpc, err := images.GetBitsPerComponent(imgs[0].RawBytes())
	if err != nil {
		t.Fatalf("GetBitsPerComponent: %v", err)
	}

	if bpc != 8 {
		t.Errorf("GetBitsPerComponent = %d, want 8", bpc)
	}
}
