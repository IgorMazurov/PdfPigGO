package core

import (
	"io"
)

// StreamWrite writes a byte slice to an io.Writer using a pooled temporary buffer.
// This mirrors the C# StreamExtensions.Write(Stream, ReadOnlySpan<byte>) polyfill
// that bridges ReadOnlySpan<byte> to Stream.Write(byte[], int, int) via ArrayPool.
// In Go, io.Writer already accepts []byte directly, so this delegates without copying.
func StreamWrite(w io.Writer, data []byte) (int, error) {
	return w.Write(data)
}

// StreamRead reads from an io.Reader into a byte slice using a pooled temporary buffer.
// This mirrors the C# StreamExtensions.Read(Stream, Span<byte>) polyfill
// that bridges Span<byte> to Stream.Read(byte[], int, int) via ArrayPool.
// In Go, io.Reader already writes into []byte directly, so this delegates without copying.
func StreamRead(r io.Reader, data []byte) (int, error) {
	return r.Read(data)
}
