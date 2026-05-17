package png

import (
	"encoding/binary"
	"errors"
	"io"
)

// PngStreamWriteHelper wraps an io.Writer for writing PNG chunk structures,
// accumulating CRC-32 data incrementally on each write.
type PngStreamWriteHelper struct {
	inner io.Writer
	crc   *Crc32State
}

var errInnerNil = errors.New("inner stream must not be nil")

// NewPngStreamWriteHelper creates a new helper wrapping the given writer.
func NewPngStreamWriteHelper(inner io.Writer) (*PngStreamWriteHelper, error) {
	if inner == nil {
		return nil, errInnerNil
	}
	return &PngStreamWriteHelper{
		inner: inner,
		crc:   NewCrc32State(),
	}, nil
}

// CanRead returns whether the underlying stream supports reading.
func (h *PngStreamWriteHelper) CanRead() bool {
	_, ok := h.inner.(io.Reader)
	return ok
}

// CanSeek returns whether the underlying stream supports seeking.
func (h *PngStreamWriteHelper) CanSeek() bool {
	_, ok := h.inner.(io.Seeker)
	return ok
}

// CanWrite returns whether the underlying stream supports writing.
func (h *PngStreamWriteHelper) CanWrite() bool {
	_, ok := h.inner.(io.Writer)
	return ok
}

// Seek delegates to the inner stream if it supports seeking.
func (h *PngStreamWriteHelper) Seek(offset int64, whence int) (int64, error) {
	if seeker, ok := h.inner.(io.Seeker); ok {
		return seeker.Seek(offset, whence)
	}
	return 0, errors.New("inner stream does not support seeking")
}

// Read delegates to the inner stream if it supports reading.
func (h *PngStreamWriteHelper) Read(p []byte) (int, error) {
	if reader, ok := h.inner.(io.Reader); ok {
		return reader.Read(p)
	}
	return 0, errors.New("inner stream does not support reading")
}

// Write writes data to the inner writer and accumulates it into the running CRC-32.
func (h *PngStreamWriteHelper) Write(p []byte) (int, error) {
	h.crc.Append(p)
	return h.inner.Write(p)
}

// Flush delegates to the inner stream if it supports flushing.
func (h *PngStreamWriteHelper) Flush() error {
	if flusher, ok := h.inner.(interface{ Flush() error }); ok {
		return flusher.Flush()
	}
	return nil
}

// SetLength sets the length of the underlying stream if it supports truncation.
func (h *PngStreamWriteHelper) SetLength(value int64) error {
	if truncater, ok := h.inner.(interface{ Truncate(size int64) error }); ok {
		return truncater.Truncate(value)
	}
	return errors.New("inner stream does not support setting length")
}

// WriteChunkLength writes a 4-byte big-endian signed integer representing the chunk length.
func (h *PngStreamWriteHelper) WriteChunkLength(length int32) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(length))
	_, err := h.inner.Write(buf)
	return err
}

// WriteChunkHeader resets the CRC accumulator and writes the chunk type header bytes.
func (h *PngStreamWriteHelper) WriteChunkHeader(header []byte) error {
	h.crc.Reset()
	_, err := h.Write(header)
	return err
}

// WriteCrc computes the current CRC-32 hash from accumulated data and writes it as a 4-byte big-endian uint32.
func (h *PngStreamWriteHelper) WriteCrc() error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, h.crc.CurrentHash())
	_, err := h.inner.Write(buf)
	return err
}
