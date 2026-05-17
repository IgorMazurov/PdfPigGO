package png

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"io"
)

// PngBuilder constructs PNG image files from pixel data.
type PngBuilder struct {
	rawData        []byte
	hasAlphaChannel bool
	width          int
	height         int
	bytesPerPixel  int
}

const deflate32KbWindow byte = 120
const checksumBits byte = 1

// Create creates a new PngBuilder for an image with the given dimensions and alpha channel setting.
func Create(width, height int, hasAlphaChannel bool) *PngBuilder {
	bpp := 3
	if hasAlphaChannel {
		bpp = 4
	}

	length := (height*width*bpp) + height

	return &PngBuilder{
		rawData:         make([]byte, length),
		hasAlphaChannel: hasAlphaChannel,
		width:           width,
		height:          height,
		bytesPerPixel:   bpp,
	}
}

// SetPixelRGB sets the RGB pixel value at column x and row y.
func (b *PngBuilder) SetPixelRGB(r, g, bval byte, x, y int) *PngBuilder {
	return b.SetPixel(NewRgbPixel(r, g, bval), x, y)
}

// SetPixel sets the pixel value at column x and row y.
func (b *PngBuilder) SetPixel(pixel Pixel, x, y int) *PngBuilder {
	start := (y*((b.width*b.bytesPerPixel)+1)) + 1 + (x*b.bytesPerPixel)

	b.rawData[start] = pixel.R
	start++
	b.rawData[start] = pixel.G
	start++
	b.rawData[start] = pixel.B
	start++

	if b.hasAlphaChannel {
		b.rawData[start] = pixel.A
	}

	return b
}

// SaveBytes returns the complete PNG file as a byte slice.
func (b *PngBuilder) SaveBytes() []byte {
	var buf bytes.Buffer
	b.Save(&buf)
	return buf.Bytes()
}

// Save writes the PNG file bytes to the provided writer.
func (b *PngBuilder) Save(w io.Writer) error {
	_, err := w.Write(expectedHeader)
	if err != nil {
		return fmt.Errorf("failed to write PNG header: %w", err)
	}

	stream, err := NewPngStreamWriteHelper(w)
	if err != nil {
		return fmt.Errorf("failed to create PNG stream helper: %w", err)
	}

	err = stream.WriteChunkLength(13)
	if err != nil {
		return fmt.Errorf("failed to write IHDR chunk length: %w", err)
	}

	err = stream.WriteChunkHeader(HeaderBytes)
	if err != nil {
		return fmt.Errorf("failed to write IHDR header: %w", err)
	}

	widthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(widthBuf, uint32(b.width))
	_, err = stream.Write(widthBuf)
	if err != nil {
		return fmt.Errorf("failed to write IHDR width: %w", err)
	}

	heightBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(heightBuf, uint32(b.height))
	_, err = stream.Write(heightBuf)
	if err != nil {
		return fmt.Errorf("failed to write IHDR height: %w", err)
	}

	colorType := ColorTypeColorUsed
	if b.hasAlphaChannel {
		colorType |= ColorTypeAlphaChannelUsed
	}

	_, err = stream.Write([]byte{8, byte(colorType), byte(DeflateWithSlidingWindow), byte(AdaptiveFiltering), byte(InterlaceMethodNone)})
	if err != nil {
		return fmt.Errorf("failed to write IHDR remaining fields: %w", err)
	}

	err = stream.WriteCrc()
	if err != nil {
		return fmt.Errorf("failed to write IHDR CRC: %w", err)
	}

	imageData, err := compress(b.rawData)
	if err != nil {
		return fmt.Errorf("failed to compress image data: %w", err)
	}

	err = stream.WriteChunkLength(int32(len(imageData)))
	if err != nil {
		return fmt.Errorf("failed to write IDAT chunk length: %w", err)
	}

	err = stream.WriteChunkHeader([]byte("IDAT"))
	if err != nil {
		return fmt.Errorf("failed to write IDAT header: %w", err)
	}

	_, err = stream.Write(imageData)
	if err != nil {
		return fmt.Errorf("failed to write IDAT data: %w", err)
	}

	err = stream.WriteCrc()
	if err != nil {
		return fmt.Errorf("failed to write IDAT CRC: %w", err)
	}

	err = stream.WriteChunkLength(0)
	if err != nil {
		return fmt.Errorf("failed to write IEND chunk length: %w", err)
	}

	err = stream.WriteChunkHeader([]byte("IEND"))
	if err != nil {
		return fmt.Errorf("failed to write IEND header: %w", err)
	}

	err = stream.WriteCrc()
	if err != nil {
		return fmt.Errorf("failed to write IEND CRC: %w", err)
	}

	return nil
}

func compress(data []byte) ([]byte, error) {
	headerLength := 2

	var compressBuf bytes.Buffer
	w, err := flate.NewWriter(&compressBuf, flate.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed to create flate writer: %w", err)
	}

	adler := newAdler32ChecksumWriter(w)
	_, _ = adler.Write(data)
	_ = adler.Close()
	_ = w.Close()

	compressed := compressBuf.Bytes()
	result := make([]byte, headerLength+len(compressed)+4)

	result[0] = deflate32KbWindow
	result[1] = checksumBits

	copy(result[headerLength:], compressed)

	checksum := adler.Checksum()
	offset := headerLength + len(compressed)
	binary.BigEndian.PutUint32(result[offset:offset+4], checksum)

	return result, nil
}

// Adler32ChecksumWriter wraps an io.Writer and computes the Adler-32 checksum of written data.
type Adler32ChecksumWriter struct {
	w          io.Writer
	checksum   uint32
	isClosed   bool
}

func newAdler32ChecksumWriter(w io.Writer) *Adler32ChecksumWriter {
	return &Adler32ChecksumWriter{w: w, checksum: 1}
}

// Write writes data to the underlying writer and updates the Adler-32 checksum.
func (a *Adler32ChecksumWriter) Write(p []byte) (int, error) {
	if a.isClosed {
		return 0, fmt.Errorf("adler32 checksum writer is closed")
	}
	n, err := a.w.Write(p)
	a.updateAdler32(p[:n])
	return n, err
}

// Close flushes the underlying writer and finalizes the checksum.
func (a *Adler32ChecksumWriter) Close() error {
	a.isClosed = true
	if closer, ok := a.w.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Checksum returns the computed Adler-32 checksum.
func (a *Adler32ChecksumWriter) Checksum() uint32 {
	return a.checksum
}

func (a *Adler32ChecksumWriter) updateAdler32(data []byte) {
	const modAdler = 65521
	s1 := uint32(a.checksum & 0xFFFF)
	s2 := uint32((a.checksum >> 16) & 0xFFFF)

	for _, b := range data {
		s1 = (s1 + uint32(b)) % modAdler
		s2 = (s2 + s1) % modAdler
	}

	a.checksum = (s2 << 16) | s1
}
