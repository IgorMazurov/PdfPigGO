package png

import (
	"bytes"
	"errors"
	"io"
	"os"
)

// Png represents a parsed PNG image. Call Open to open from stream, bytes, or file path.
type Png struct {
	header              *ImageHeader
	data                *RawPngData
	hasTransparencyChunk bool
}

// Header returns the header data from the PNG image.
func (p *Png) Header() *ImageHeader {
	return p.header
}

// Width returns the width of the image in pixels.
func (p *Png) Width() int {
	return p.header.Width()
}

// Height returns the height of the image in pixels.
func (p *Png) Height() int {
	return p.header.Height()
}

// HasAlphaChannel reports whether the image has an alpha (transparency) layer.
func (p *Png) HasAlphaChannel() bool {
	return (p.header.ColorType()&ColorTypeAlphaChannelUsed) != 0 || p.hasTransparencyChunk
}

// GetPixel returns the pixel at the given column and row coordinates.
// Pixel values are generated on demand from the underlying data to prevent holding
// many items in memory at once, so consumers should cache values if they're going
// to be looped over many times.
func (p *Png) GetPixel(x, y int) Pixel {
	return p.data.GetPixel(x, y)
}

// NewPng creates a new Png from the given header, raw data, and transparency flag.
func NewPng(header *ImageHeader, data *RawPngData, hasTransparencyChunk bool) (*Png, error) {
	if header == nil {
		return nil, errors.New("header must not be nil")
	}

	if data == nil {
		return nil, errors.New("data must not be nil")
	}

	return &Png{
		header:               header,
		data:                 data,
		hasTransparencyChunk: hasTransparencyChunk,
	}, nil
}

// Open reads a PNG image from an io.Reader stream. chunkVisitor is optional and may be nil.
func Open(stream io.Reader, chunkVisitor ChunkVisitor) (*Png, error) {
	return OpenWithSettings(stream, &PngOpenerSettings{
		ChunkVisitor: chunkVisitor,
	})
}

// OpenBytes reads a PNG image from a byte slice. chunkVisitor is optional and may be nil.
func OpenBytes(data []byte, chunkVisitor ChunkVisitor) (*Png, error) {
	reader := bytes.NewReader(data)
	return OpenWithSettings(reader, &PngOpenerSettings{
		ChunkVisitor: chunkVisitor,
	})
}

// OpenFile reads a PNG image from the given file path. The file will be locked during reading.
// chunkVisitor is optional and may be nil.
func OpenFile(filePath string, chunkVisitor ChunkVisitor) (*Png, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return OpenWithSettings(file, &PngOpenerSettings{
		ChunkVisitor: chunkVisitor,
	})
}
