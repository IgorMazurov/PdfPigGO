package filters

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// DefaultFilterProvider is the default implementation of FilterProvider that
// registers all standard PDF stream filters. It is a singleton accessed via Instance.
type DefaultFilterProvider struct {
	BaseFilterProvider
}

// Instance is the single shared instance of DefaultFilterProvider.
var Instance = newDefaultFilterProvider()

func newDefaultFilterProvider() *DefaultFilterProvider {
	return &DefaultFilterProvider{
		BaseFilterProvider: *NewBaseFilterProvider(buildFilterDictionary()),
	}
}

func buildFilterDictionary() map[string]Filter {
	ascii85 := NewAscii85Filter()
	asciiHex := NewAsciiHexDecodeFilter()
	ccitt := &CcittFaxDecodeFilter{}
	dct := &DctDecodeFilter{}
	flate := NewFlateFilter()
	jbig2 := NewJbig2DecodeFilter()
	jpx := NewJpxDecodeFilter()
	runLength := NewRunLengthFilter()
	lzw := NewLzwFilter()

	return map[string]Filter{
		tokens.Ascii85Decode.Data():              ascii85,
		tokens.Ascii85DecodeAbbreviation.Data():  ascii85,
		tokens.AsciiHexDecode.Data():             asciiHex,
		tokens.AsciiHexDecodeAbbreviation.Data(): asciiHex,
		tokens.CcittfaxDecode.Data():             ccitt,
		tokens.CcittfaxDecodeAbbreviation.Data(): ccitt,
		tokens.DctDecode.Data():                  dct,
		tokens.DctDecodeAbbreviation.Data():      dct,
		tokens.FlateDecode.Data():                flate,
		tokens.FlateDecodeAbbreviation.Data():    flate,
		tokens.Jbig2Decode.Data():                jbig2,
		tokens.JpxDecode.Data():                  jpx,
		tokens.RunLengthDecode.Data():            runLength,
		tokens.RunLengthDecodeAbbreviation.Data(): runLength,
		tokens.LzwDecode.Data():                  lzw,
		tokens.LzwDecodeAbbreviation.Data():      lzw,
	}
}

// RunLengthFilter decodes data encoded with the Run Length filter.
// The encoded data is a sequence of runs, where each run consists of a length
// byte followed by 1 to 128 bytes of data. A length byte in the range 0-127
// indicates that the following length+1 bytes should be copied literally.
// A length byte in the range 129-255 indicates that the single following byte
// should be repeated 257-length times (between 2 and 128 repetitions).
// A length byte of 128 signals end of data.
type RunLengthFilter struct{}

// NewRunLengthFilter creates a new RunLengthFilter instance.
func NewRunLengthFilter() *RunLengthFilter {
	return &RunLengthFilter{}
}

// IsSupported reports whether this filter is supported by the current build.
func (f *RunLengthFilter) IsSupported() bool {
	return true
}

// Decode decodes Run Length-encoded input bytes.
func (f *RunLengthFilter) Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider FilterProvider, filterIndex int) ([]byte, error) {
	// Guard against decompression bombs: limit output to min(input*100+1MB, 50MB).
	maxDecodedSize := len(input)*100 + 1024*1024
	if maxDecodedSize > 50*1024*1024 {
		maxDecodedSize = 50 * 1024 * 1024
	}

	// Use simple slice-based buffer to avoid O(n²) copy overhead of ArrayPoolBufferWriter.
	result := make([]byte, 0, len(input))

	i := 0
	for i < len(input) {
		runLength := input[i]

		if runLength == 128 {
			break
		}

		if runLength <= 127 {
			count := int(runLength) + 1
			i++
			for count > 0 && i < len(input) {
				if len(result) >= maxDecodedSize {
					goto exceeded
				}
				result = append(result, input[i])
				i++
				count--
			}
		} else {
			if i+1 >= len(input) {
				break
			}
			repeatCount := 257 - int(runLength)
			if len(result)+repeatCount > maxDecodedSize {
				goto exceeded
			}
			for j := 0; j < repeatCount; j++ {
				result = append(result, input[i+1])
			}
			i += 2
		}
	}

	return result, nil

exceeded:
	panic(core.NewPdfDocumentFormatException(
		fmt.Sprintf("Decoded stream size exceeds the estimated maximum size. Current decoded stream length: %d, 1 filters applied out of 1",
			len(result))))
}


// CcittFaxDecodeFilter implements the CCITT Group 3/4 fax compression filter.
type CcittFaxDecodeFilter struct{}

// IsSupported reports whether this filter is supported by the current build.
func (f *CcittFaxDecodeFilter) IsSupported() bool {
	return true
}

// Decode decodes CCITT fax-compressed input bytes.
func (f *CcittFaxDecodeFilter) Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider FilterProvider, filterIndex int) ([]byte, error) {
	return ccittDecode(input, streamDictionary, filterIndex)
}
