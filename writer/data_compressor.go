package writer

import (
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/writer/interfaces"
)

// CompressBytes compresses the given byte slice using Flate (zlib/deflate).
func CompressBytes(data []byte) []byte {
	return interfaces.CompressBytes(data)
}

// DataCompressorCompressBytes is the legacy name for CompressBytes, kept for backward compatibility.
func DataCompressorCompressBytes(data []byte) []byte {
	return CompressBytes(data)
}

// CompressToStream creates a StreamToken containing the Flate-compressed data
// with appropriate dictionary entries for Length, Length1 (original length),
// and Filter.
func CompressToStream(data []byte) *tokens.StreamToken {
	return interfaces.CompressToStream(data)
}
