package filters

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// DctDecodeFilter represents the DCT (Discrete Cosine Transform) filter used for JPEG-encoded data.
// This filter is not implemented; raw JPEG data can be extracted and passed to a JPEG library.
type DctDecodeFilter struct{}

// NewDctDecodeFilter creates a new DctDecodeFilter instance.
func NewDctDecodeFilter() *DctDecodeFilter {
	return &DctDecodeFilter{}
}

// IsSupported reports whether this filter is supported. Always false since JPEG decoding is not implemented.
func (f *DctDecodeFilter) IsSupported() bool {
	return false
}

var errDctNotSupported = errors.New("DCT (Discrete Cosine Transform) filter indicates data is encoded in JPEG format; this filter is not currently supported but the raw data can be supplied to JPEG supporting libraries")

// Decode returns an error because DCT decoding is not implemented.
func (f *DctDecodeFilter) Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider FilterProvider, filterIndex int) ([]byte, error) {
	return nil, errDctNotSupported
}

var _ Filter = (*DctDecodeFilter)(nil)
