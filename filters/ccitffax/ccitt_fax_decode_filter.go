package ccitffax

import (
	"errors"
	"io"

	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/iostream"
	"github.com/uglytoad/pdfpig/go/tokens"
)

var errUnsupported = errors.New("unsupported operation")

// readOnlyBytes is a minimal io.ReadWriteSeeker wrapper around a byte slice.
type readOnlyBytes struct {
	data []byte
	pos  int64
}

func newReadOnlyBytes(data []byte) *readOnlyBytes {
	return &readOnlyBytes{data: data, pos: 0}
}

func (r *readOnlyBytes) Read(p []byte) (int, error) {
	if r.pos >= int64(len(r.data)) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += int64(n)
	return n, nil
}

func (r *readOnlyBytes) Write(p []byte) (int, error) {
	return 0, errUnsupported
}

func (r *readOnlyBytes) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.pos = offset
	case io.SeekCurrent:
		r.pos += offset
	case io.SeekEnd:
		r.pos = int64(len(r.data)) + offset
	default:
		return 0, io.ErrUnexpectedEOF
	}
	if r.pos < 0 {
		r.pos = 0
	}
	return r.pos, nil
}

// CcittFaxDecodeFilter decodes image data encoded using CCITT Group 3 or Group 4 fax compression.
type CcittFaxDecodeFilter struct{}

// NewCcittFaxDecodeFilter creates a new CcittFaxDecodeFilter instance.
func NewCcittFaxDecodeFilter() *CcittFaxDecodeFilter {
	return &CcittFaxDecodeFilter{}
}

// IsSupported reports whether this filter is supported. Always true.
func (f *CcittFaxDecodeFilter) IsSupported() bool {
	return true
}

// Decode decodes CCITT fax-compressed input bytes.
func (f *CcittFaxDecodeFilter) Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider filters.FilterProvider, filterIndex int) ([]byte, error) {
	if streamDictionary == nil || len(input) == 0 {
		return nil, nil
	}

	decodeParms := getFilterParameters(streamDictionary, filterIndex)

	cols := getIntOrDefault(decodeParms, tokens.Columns, 1728)
	rows := getIntOrDefault(decodeParms, tokens.Rows, 0)
	height := getIntOrDefault(streamDictionary, tokens.Height, 0)
	if height == 0 {
		hToken, ok := streamDictionary.TryGet(tokens.H)
		if ok {
			if hNumeric, ok := hToken.(*tokens.NumericToken); ok {
				height = hNumeric.IntVal()
			}
		}
	}

	if rows > 0 && height > 0 {
		rows = height
	} else if height > rows {
		rows = height
	}

	k := getIntOrDefault(decodeParms, tokens.K, 0)
	encodedByteAlign := getBoolOrDefault(decodeParms, tokens.EncodedByteAlign, false)
	compressionType := determineCompressionType(input, k)

	stream := NewCcittFaxDecoderStream(
		iostream.NewStreamWrapper(newReadOnlyBytes(input)),
		cols,
		compressionType,
		encodedByteAlign,
	)

	arraySize := (cols + 7) / 8 * rows
	decompressed := make([]byte, arraySize)
	readFromDecoderStream(stream, decompressed)

	blackIsOne := getBoolOrDefault(decodeParms, tokens.BlackIs1, false)
	if !blackIsOne {
		invertBitmap(decompressed)
	}

	return decompressed, nil
}

func determineCompressionType(input []byte, k int) CcittFaxCompressionType {
	if k == 0 {
		compressionType := Group3_1D

		if len(input) >= 2 && (input[0] != 0 || (input[1]>>4 != 1 && input[1] != 1)) {
			compressionType = ModifiedHuffman
			b := ((int(input[0]) << 8) | int(input[1]&0xff)) >> 4
			for i := 12; i < 160 && i/8 < len(input); i++ {
				b = (b << 1) | int((input[i/8]>>(7-i%8))&0x01)
				if b&0xFFF == 1 {
					return Group3_1D
				}
			}
		}

		return compressionType
	}

	if k > 0 {
		return Group3_2D
	}

	return Group4_2D
}

func readFromDecoderStream(d *CcittFaxDecoderStream, result []byte) {
	pos := 0
	for pos < len(result) {
		n, _ := d.Read(result, pos, len(result)-pos)
		if n <= 0 {
			break
		}
		pos += n
	}
}

func invertBitmap(data []byte) {
	for i := range data {
		data[i] = ^data[i]
	}
}

// getFilterParameters resolves the decode parameters dictionary for a given filter index.
// This mirrors the C# DecodeParameterResolver.GetFilterParameters logic.
func getFilterParameters(streamDictionary *tokens.DictionaryToken, filterIndex int) *tokens.DictionaryToken {
	if streamDictionary == nil || filterIndex < 0 {
		return emptyDict()
	}

	filter := getObjectOrDefault(streamDictionary, tokens.Filter, tokens.F)
	parameters := getObjectOrDefault(streamDictionary, tokens.DecodeParms, tokens.Dp)

	switch filter.(type) {
	case *tokens.NameToken:
		if dict, ok := parameters.(*tokens.DictionaryToken); ok {
			return dict
		}
	case *tokens.ArrayToken:
		if arr, ok := parameters.(*tokens.ArrayToken); ok {
			if filterIndex < len(arr.Data()) && filterIndex >= 0 {
				if dict, ok := arr.Data()[filterIndex].(*tokens.DictionaryToken); ok {
					return dict
				}
			}
		}
	}

	return emptyDict()
}

func getObjectOrDefault(dict *tokens.DictionaryToken, names ...*tokens.NameToken) tokens.Token {
	for _, name := range names {
		if token, ok := dict.TryGet(name); ok {
			return token
		}
	}
	return nil
}

func getIntOrDefault(dict *tokens.DictionaryToken, name *tokens.NameToken, defaultValue int) int {
	token, ok := dict.TryGet(name)
	if !ok {
		return defaultValue
	}
	if numeric, ok := token.(*tokens.NumericToken); ok {
		return numeric.IntVal()
	}
	return defaultValue
}

func getBoolOrDefault(dict *tokens.DictionaryToken, name *tokens.NameToken, defaultValue bool) bool {
	token, ok := dict.TryGet(name)
	if !ok {
		return defaultValue
	}
	if boolean, ok := token.(*tokens.BooleanToken); ok {
		return boolean.Data()
	}
	return defaultValue
}

var emptyDictInstance = func() *tokens.DictionaryToken {
	d, _ := tokens.NewDictionary(make(map[*tokens.NameToken]tokens.Token))
	return d
}()

func emptyDict() *tokens.DictionaryToken {
	return emptyDictInstance
}

var _ filters.Filter = (*CcittFaxDecodeFilter)(nil)
