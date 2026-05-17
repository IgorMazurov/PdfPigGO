package images_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/uglytoad/pdfpig/go/images"
)

func TestCanGetJpegInformation(t *testing.T) {
	data, err := loadJpg("testdata/jpg/218995467-ccb746b0-dc28-4616-bcb1-4ad685f81876.jpg")
	if err != nil {
		t.Fatalf("loadJpg error: %v", err)
	}

	ms := bytes.NewReader(data)

	jpegInfo, err := images.GetInformation(ms)
	if err != nil {
		t.Fatalf("GetInformation error: %v", err)
	}

	if jpegInfo.BitsPerComponent != 8 {
		t.Errorf("BitsPerComponent = %d, want 8", jpegInfo.BitsPerComponent)
	}

	if jpegInfo.Height != 2290 {
		t.Errorf("Height = %d, want 2290", jpegInfo.Height)
	}

	if jpegInfo.Width != 1648 {
		t.Errorf("Width = %d, want 1648", jpegInfo.Width)
	}
}

func loadJpg(name string) ([]byte, error) {
	return os.ReadFile(name)
}
