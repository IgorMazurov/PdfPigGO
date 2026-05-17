package filters

import (
	"bytes"
	"fmt"

	"github.com/uglytoad/pdfpig/go/filters/lzw"
	"github.com/uglytoad/pdfpig/go/tokens"
)

const (
	lzwDefaultColors            = 1
	lzwDefaultBitsPerComponent  = 8
	lzwDefaultColumns           = 1

	lzwClearTable      = 256
	lzwEodMarker       = 257

	lzwNineBitBoundary  = 511
	lzwTenBitBoundary   = 1023
	lzwElevenBitBoundary = 2047
)

// LzwFilter decodes LZW (Lempel-Ziv-Welch) compressed PDF stream data.
// The filter is a variable-length, adaptive compression method adopted as one
// of the standard compression methods in the TIFF standard.
type LzwFilter struct{}

// NewLzwFilter creates a new LzwFilter instance.
func NewLzwFilter() *LzwFilter {
	return &LzwFilter{}
}

// IsSupported reports whether this filter is supported by the current build.
func (f *LzwFilter) IsSupported() bool {
	return true
}

// Decode decodes LZW-compressed input bytes using adaptive dictionary decompression.
// If a predictor is specified in the decode parameters, it is applied during
// decompression to further reconstruct image data. Returns the original input
// on any error.
func (f *LzwFilter) Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider FilterProvider, filterIndex int) ([]byte, error) {
	params, err := GetFilterParameters(streamDictionary, filterIndex)
	if err != nil {
		return input, nil
	}

	predictor := getIntOrDefault(params, tokens.Predictor, -1)
	earlyChange := getIntOrDefault(params, tokens.EarlyChange, 1)
	colors := min(getIntOrDefault(params, tokens.Colors, lzwDefaultColors), 32)
	bitsPerComponent := getIntOrDefault(params, tokens.BitsPerComponent, lzwDefaultBitsPerComponent)
	columns := getIntOrDefault(params, tokens.Columns, lzwDefaultColumns)

	result, err := decodeLzw(input, earlyChange == 1, predictor, colors, bitsPerComponent, columns)
	if err != nil {
		return input, fmt.Errorf("lzw decode: %w", err)
	}

	return result, nil
}

func decodeLzw(input []byte, isEarlyChange bool, predictor, colors, bitsPerComponent, columns int) ([]byte, error) {
	output := &bytes.Buffer{}
	output.Grow(len(input) * 3 / 2)

	outStream := wrapPredictor(output, predictor, colors, bitsPerComponent, columns)

	table := getDefaultTable()

	codeBits := 9
	data := lzw.NewBitStream(input)

	codeOffset := 0
	if !isEarlyChange {
		codeOffset = 1
	}

	previous := -1

	for {
		nextVal, err := data.Get(codeBits)
		if err != nil {
			break
		}
		next := int(nextVal)

		if next == lzwEodMarker {
			break
		}

		if next == lzwClearTable {
			table = getDefaultTable()
			previous = -1
			codeBits = 9
			continue
		}

		b, ok := table[next]
		if ok {
			outStream.Write(b)

			if previous >= 0 {
				lastSequence := table[previous]

				newSequence := make([]byte, len(lastSequence)+1)
				copy(newSequence, lastSequence)
				newSequence[len(lastSequence)] = b[0]

				table[len(table)] = newSequence
			}
		} else {
			lastSequence := table[previous]

			newSequence := make([]byte, len(lastSequence)+1)
			copy(newSequence, lastSequence)
			newSequence[len(lastSequence)] = lastSequence[0]

			outStream.Write(newSequence)
			table[len(table)] = newSequence
		}

		previous = next

		if len(table) >= lzwElevenBitBoundary+codeOffset {
			codeBits = 12
		} else if len(table) >= lzwTenBitBoundary+codeOffset {
			codeBits = 11
		} else if len(table) >= lzwNineBitBoundary+codeOffset {
			codeBits = 10
		} else {
			codeBits = 9
		}
	}

	if flusher, ok := outStream.(interface{ Flush() error }); ok {
		if err := flusher.Flush(); err != nil {
			return nil, fmt.Errorf("lzw decode: flush failed: %w", err)
		}
	}

	return output.Bytes(), nil
}

func getDefaultTable() map[int][]byte {
	table := make(map[int][]byte, 258)

	for i := 0; i < 256; i++ {
		table[i] = []byte{byte(i)}
	}

	table[lzwClearTable] = nil
	table[lzwEodMarker] = nil

	return table
}

var _ Filter = (*LzwFilter)(nil)
