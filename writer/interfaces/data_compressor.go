package interfaces

import (
	"bytes"

	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CompressBytes compresses the given byte slice using Flate (zlib/deflate).
func CompressBytes(data []byte) []byte {
	flater := filters.NewFlateFilter()
	compressed, _ := flater.Encode(bytes.NewReader(data))
	return compressed
}

// CompressToStream creates a StreamToken containing the Flate-compressed data
// with appropriate dictionary entries for Length, Length1 (original length),
// and Filter.
func CompressToStream(data []byte) *tokens.StreamToken {
	compressed := CompressBytes(data)

	dictEntries := map[*tokens.NameToken]tokens.Token{
		tokens.Length:     tokens.NewNumericTokenFromInt(len(compressed)),
		tokens.Length1:    tokens.NewNumericTokenFromInt(len(data)),
		tokens.Filter:     tokens.NewArrayToken([]tokens.Token{tokens.FlateDecode}),
	}

	dict, _ := tokens.NewDictionary(dictEntries)
	stream, _ := tokens.NewStreamToken(dict, compressed)

	return stream
}
