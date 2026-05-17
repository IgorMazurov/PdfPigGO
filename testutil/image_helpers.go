package testutil

import (
	"bytes"
	"compress/flate"
	"os"
	"path/filepath"

	"github.com/uglytoad/pdfpig/go/images/png"
)

// FilesFolder is the base directory for image test fixture files.
var FilesFolder string

// SetFilesFolder sets the base directory for image test fixture files.
func SetFilesFolder(folder string) {
	FilesFolder = folder
}

// LoadFileBytes reads a file from the configured FilesFolder and returns its contents as bytes.
// If isCompressed is true, the file content is decompressed using raw DEFLATE before returning.
func LoadFileBytes(filename string, isCompressed bool) ([]byte, error) {
	filePath := filepath.Join(FilesFolder, filename)

	if isCompressed {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		reader := flate.NewReader(bytes.NewReader(data))
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		reader.Close()
		if err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}

	return os.ReadFile(filePath)
}

// ImagesAreEqual compares two PNG images pixel by pixel.
// Returns true if both images have identical dimensions, alpha channel presence, and every pixel value matches.
func ImagesAreEqual(first, second []byte) (bool, error) {
	png1, err := png.OpenBytes(first, nil)
	if err != nil {
		return false, err
	}

	png2, err := png.OpenBytes(second, nil)
	if err != nil {
		return false, err
	}

	if png1.Width() != png2.Width() || png1.Height() != png2.Height() || png1.HasAlphaChannel() != png2.HasAlphaChannel() {
		return false, nil
	}

	for y := 0; y < png1.Height(); y++ {
		for x := 0; x < png1.Width(); x++ {
			if !png1.GetPixel(x, y).Equals(png2.GetPixel(x, y)) {
				return false, nil
			}
		}
	}

	return true, nil
}
