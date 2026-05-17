package core

import "sync"

// bufferPool caches slices to reduce allocations.
var bufferPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 0, 256)
		return &buf
	},
}

// BufferWriter is the Go equivalent of .NET's IBufferWriter[T].
type BufferWriter interface {
	Advance(count int)
	GetMemory(sizeHint int) []byte
	GetSpan(sizeHint int) []byte
}

// ArrayPoolBufferWriter provides a pooled, growable byte buffer writer.
// It rents buffers from a sync.Pool and returns them on Dispose.
type ArrayPoolBufferWriter struct {
	buffer   []byte
	position int
}

const defaultBufferSize = 256

// NewArrayPoolBufferWriter creates a new pooled buffer writer with the default initial size.
func NewArrayPoolBufferWriter() *ArrayPoolBufferWriter {
	buf := bufferPool.Get().(*[]byte)
	*buf = (*buf)[:0]
	if cap(*buf) < defaultBufferSize {
		*buf = make([]byte, 0, defaultBufferSize)
	}
	return &ArrayPoolBufferWriter{buffer: *buf}
}

// NewArrayPoolBufferWriterWithSize creates a new pooled buffer writer with the given initial size.
func NewArrayPoolBufferWriterWithSize(size int) *ArrayPoolBufferWriter {
	buf := make([]byte, 0, size)
	return &ArrayPoolBufferWriter{buffer: buf}
}

// Advance moves the write position forward by count bytes.
func (w *ArrayPoolBufferWriter) Advance(count int) {
	w.position += count
}

// WriteSingle writes a single byte value to the buffer.
func (w *ArrayPoolBufferWriter) WriteSingle(value byte) {
	span := w.GetSpan(1)
	span[0] = value
	w.position++
}

// WriteBytes writes a slice of bytes to the buffer.
func (w *ArrayPoolBufferWriter) WriteBytes(values []byte) {
	span := w.GetSpan(len(values))
	copy(span, values)
	w.position += len(values)
}

// GetMemory returns a writable memory region starting at the current position.
// In Go this is represented as a slice.
func (w *ArrayPoolBufferWriter) GetMemory(sizeHint int) []byte {
	w.ensureCapacity(sizeHint)
	needed := w.position + sizeHint
	if len(w.buffer) < needed {
		w.buffer = w.buffer[:needed]
	}
	return w.buffer[w.position:]
}

// GetSpan returns a writable span starting at the current position.
func (w *ArrayPoolBufferWriter) GetSpan(sizeHint int) []byte {
	w.ensureCapacity(sizeHint)
	needed := w.position + sizeHint
	if len(w.buffer) < needed {
		w.buffer = w.buffer[:needed]
	}
	return w.buffer[w.position:]
}

// WrittenCount returns the number of bytes written to the buffer.
func (w *ArrayPoolBufferWriter) WrittenCount() int {
	return w.position
}

// WrittenMemory returns the committed data as a read-only slice.
func (w *ArrayPoolBufferWriter) WrittenMemory() []byte {
	return w.buffer[:w.position]
}

// WrittenSpan returns the committed data as a read-only slice.
func (w *ArrayPoolBufferWriter) WrittenSpan() []byte {
	return w.buffer[:w.position]
}

func (w *ArrayPoolBufferWriter) ensureCapacity(sizeHint int) {
	if sizeHint == 0 {
		sizeHint = 1
	}

	remaining := cap(w.buffer) - w.position
	if sizeHint > remaining {
		newSize := w.position + sizeHint
		if newSize < 512 {
			newSize = 512
		}
		newBuf := make([]byte, w.position, newSize)
		copy(newBuf, w.buffer[:w.position])
		w.returnBuffer()
		w.buffer = newBuf
	}
}

func (w *ArrayPoolBufferWriter) returnBuffer() {
	if len(w.buffer) != 0 && cap(w.buffer) >= defaultBufferSize {
		w.buffer = w.buffer[:0]
		buf := w.buffer
		bufferPool.Put(&buf)
	}
	w.buffer = nil
}

// Reset resets the internal state so the instance can be reused before disposal.
// If clearArray is true, written bytes are zeroed.
func (w *ArrayPoolBufferWriter) Reset(clearArray bool) {
	if clearArray && w.position > 0 {
		for i := 0; i < w.position; i++ {
			w.buffer[i] = 0
		}
	}
	w.position = 0
}

// Dispose returns the rented buffer to the pool for reuse.
func (w *ArrayPoolBufferWriter) Dispose() {
	w.returnBuffer()
}
