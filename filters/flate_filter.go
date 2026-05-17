package filters

import (
	"bytes"
	compress "compress/flate"
	"fmt"
	"io"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// flateReaderPool caches decompressors to avoid re-allocating the internal
// dictionary decoder (~32KB+) on every decompression call.
var flateReaderPool = sync.Pool{
	New: func() any {
		r := compress.NewReader(nil)
		return r
	},
}

const (
	defaultColors             = 1
	defaultBitsPerComponent   = 8
	defaultColumns            = 1
	deflate32KbWindow         = byte(120)
	checksumBits              = byte(1)
)

// FlateFilter decodes Flate (zlib/deflate) compressed PDF stream data.
// The filter is based on the public-domain zlib/deflate compression method,
// a variable-length Lempel-Ziv adaptive compression cascaded with adaptive
// Huffman coding as defined in RFC 1950 and RFC 1951.
type FlateFilter struct{}

// NewFlateFilter creates a new FlateFilter instance.
func NewFlateFilter() *FlateFilter {
	return &FlateFilter{}
}

// IsSupported reports whether this filter is supported by the current build.
func (f *FlateFilter) IsSupported() bool {
	return true
}

// Decode decodes Flate-compressed input bytes using deflate decompression.
// If a predictor is specified in the decode parameters, it is applied during
// decompression to further reconstruct image data. Returns the original input
// on any error.
func (f *FlateFilter) Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider FilterProvider, filterIndex int) ([]byte, error) {
	params, err := GetFilterParameters(streamDictionary, filterIndex)
	if err != nil {
		return input, nil
	}

	predictor := getIntOrDefault(params, tokens.Predictor, -1)

	colors := min(getIntOrDefault(params, tokens.Colors, defaultColors), 32)
	bitsPerComponent := getIntOrDefault(params, tokens.BitsPerComponent, defaultBitsPerComponent)
	columns := getIntOrDefault(params, tokens.Columns, defaultColumns)

	result, err := decompress(input, predictor, colors, bitsPerComponent, columns)
	if err != nil {
		return input, nil
	}

	return result, nil
}

// Encode compresses raw data using Flate (zlib/deflate) encoding with an
// Adler-32 checksum appended as required by the zlib format specification.
func (f *FlateFilter) Encode(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("flate encode: read input: %w", err)
	}

	compressedBuf := &bytes.Buffer{}
	deflater, err := compress.NewWriter(compressedBuf, compress.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("flate encode: create compressor: %w", err)
	}
	adlerStream := NewAdler32ChecksumStream(deflater)

	if _, err := adlerStream.Write(data); err != nil {
		return nil, fmt.Errorf("flate encode: write compressed data: %w", err)
	}

	if err := adlerStream.Close(); err != nil {
		return nil, fmt.Errorf("flate encode: close compressor: %w", err)
	}

	compressed := compressedBuf.Bytes()
	checksum := adlerStream.Checksum()

	headerLen := 2
	checksumLen := 4
	result := make([]byte, headerLen+len(compressed)+checksumLen)

	result[0] = deflate32KbWindow
	result[1] = checksumBits

	copy(result[headerLen:], compressed)

	offset := headerLen + len(compressed)
	result[offset] = byte(checksum >> 24)
	result[offset+1] = byte(checksum >> 16)
	result[offset+2] = byte(checksum >> 8)
	result[offset+3] = byte(checksum)

	return result, nil
}

func decompress(input []byte, predictor, colors, bitsPerComponent, columns int) ([]byte, error) {
	reader := core.AsReadOnlyMemoryStream(input)

	_, err := reader.ReadByte()
	if err != nil {
		return nil, fonts.NewCorruptCompressedDataException("invalid Flate compressed stream: missing header byte 1")
	}
	_, err = reader.ReadByte()
	if err != nil {
		return nil, fonts.NewCorruptCompressedDataException("invalid Flate compressed stream: missing header byte 2")
	}

	deflateReader := flateReaderPool.Get().(io.ReadCloser)
	defer flateReaderPool.Put(deflateReader)

	if resetter, ok := deflateReader.(interface{ Reset(io.Reader, []byte) error }); ok {
		if err := resetter.Reset(reader, nil); err != nil {
			return nil, fonts.NewCorruptCompressedDataException("failed to reset flate reader: " + err.Error())
		}
	}

	// Pre-allocate output buffer. Content streams (predictor=-1) are typically
	// text with low compression ratio (~1:1), while image data can be 2-3x.
	capacity := len(input) * 3 / 2
	if predictor == -1 {
		capacity = len(input)
	}

	output := bytes.NewBuffer(make([]byte, 0, capacity))
	outStream := wrapPredictor(output, predictor, colors, bitsPerComponent, columns)

	if _, err := io.Copy(outStream, deflateReader); err != nil {
		return nil, fonts.NewCorruptCompressedDataExceptionWithInner("invalid Flate compressed stream encountered", err)
	}

	if flusher, ok := outStream.(interface{ Flush() error }); ok {
		if err := flusher.Flush(); err != nil {
			return nil, fonts.NewCorruptCompressedDataExceptionWithInner("invalid Flate compressed stream: flush failed", err)
		}
	}

	return output.Bytes(), nil
}

func getIntOrDefault(dict *tokens.DictionaryToken, name *tokens.NameToken, defaultValue int) int {
	if dict == nil {
		return defaultValue
	}
	if token, ok := dict.TryGet(name); ok {
		if num, ok := token.(*tokens.NumericToken); ok {
			return int(num.Data())
		}
	}
	return defaultValue
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func wrapPredictor(out io.Writer, predictor, colors, bitsPerComponent, columns int) io.Writer {
	return WrapPredictor(out, predictor, colors, bitsPerComponent, columns)
}

var _ Filter = (*FlateFilter)(nil)
