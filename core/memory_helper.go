package core

import (
	"bytes"
)

// AsReadOnlyMemoryStream creates a read-only stream from a byte slice.
// The returned *bytes.Reader shares the underlying array with no additional allocation,
// matching the zero-copy attempt in the C# TryGetArray path.
func AsReadOnlyMemoryStream(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}

// AsBytesFromBuffer extracts the byte slice from a *bytes.Buffer.
// The returned slice shares the underlying buffer with no additional allocation,
// matching the zero-copy attempt in the C# TryGetBuffer path.
func AsBytesFromBuffer(buf *bytes.Buffer) []byte {
	return buf.Bytes()
}
