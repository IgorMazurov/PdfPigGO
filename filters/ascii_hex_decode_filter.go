package filters

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

var asciiHexReverse = [128]int8{
	0:   -1, 1:  -1, 2:  -1, 3:  -1, 4:  -1, 5:  -1, 6:  -1, 7:  -1, 8:  -1, 9:  -1,
	10:  -1, 11: -1, 12: -1, 13: -1, 14: -1, 15: -1, 16: -1, 17: -1, 18: -1, 19: -1,
	20:  -1, 21: -1, 22: -1, 23: -1, 24: -1, 25: -1, 26: -1, 27: -1, 28: -1, 29: -1,
	30:  -1, 31: -1, 32: -1, 33: -1, 34: -1, 35: -1, 36: -1, 37: -1, 38: -1, 39: -1,
	40:  -1, 41: -1, 42: -1, 43: -1, 44: -1, 45: -1, 46: -1, 47: -1, 48: 0,  49: 1,
	50:  2,  51: 3,  52: 4,  53: 5,  54: 6,  55: 7,  56: 8,  57: 9,  58: -1, 59: -1,
	60:  -1, 61: -1, 62: -1, 63: -1, 64: -1, 65: 10, 66: 11, 67: 12, 68: 13, 69: 14,
	70:  15, 71: -1, 72: -1, 73: -1, 74: -1, 75: -1, 76: -1, 77: -1, 78: -1, 79: -1,
	80:  -1, 81: -1, 82: -1, 83: -1, 84: -1, 85: -1, 86: -1, 87: -1, 88: -1, 89: -1,
	90:  -1, 91: -1, 92: -1, 93: -1, 94: -1, 95: -1, 96: -1, 97: 10, 98: 11, 99: 12,
	100: 13, 101: 14, 102: 15,
}

// AsciiHexDecodeFilter decodes ASCII hexadecimal encoded PDF stream data.
type AsciiHexDecodeFilter struct{}

// NewAsciiHexDecodeFilter creates a new AsciiHexDecodeFilter instance.
func NewAsciiHexDecodeFilter() *AsciiHexDecodeFilter {
	return &AsciiHexDecodeFilter{}
}

// IsSupported reports whether this filter is supported. Always true.
func (f *AsciiHexDecodeFilter) IsSupported() bool {
	return true
}

// Decode decodes ASCII hex-encoded input bytes.
func (f *AsciiHexDecodeFilter) Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider FilterProvider, filterIndex int) ([]byte, error) {
	pair := [2]byte{}
	index := 0

	writer := core.NewArrayPoolBufferWriterWithSize(len(input))
	defer writer.Dispose()

	for i := 0; i < len(input); i++ {
		if input[i] == '>' {
			break
		}

		if isAsciiHexWhitespace(input[i]) || input[i] == '<' {
			continue
		}

		pair[index] = input[i]
		index++

		if index == 2 {
			if !writeAsciiHexByte(&pair, writer) {
				writer.Dispose()
				return nil, fmt.Errorf("invalid ASCII hex encoding at position %d", i)
			}

			index = 0
		}
	}

	if index > 0 {
		if index == 1 {
			pair[1] = '0'
		}

		if !writeAsciiHexByte(&pair, writer) {
			writer.Dispose()
			return nil, fmt.Errorf("invalid ASCII hex encoding for trailing character")
		}
	}

	result := make([]byte, len(writer.WrittenMemory()))
	copy(result, writer.WrittenMemory())
	return result, nil
}

func writeAsciiHexByte(hexBytes *[2]byte, writer *core.ArrayPoolBufferWriter) bool {
	first := asciiHexReverse[hexBytes[0]]
	second := asciiHexReverse[hexBytes[1]]

	if first == -1 || second == -1 {
		return false
	}

	value := byte(first*16 + second)
	writer.WriteSingle(value)
	return true
}

func isAsciiHexWhitespace(c byte) bool {
	return c == 0 || c == '\t' || c == '\n' || c == '\f' || c == '\r' || c == ' '
}

var _ Filter = (*AsciiHexDecodeFilter)(nil)
