package filters

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// JpxDecodeFilter represents the JPX (JPEG 2000) filter for image data.
// This filter is not implemented; raw JPEG 2000 compressed data can be accessed directly.
type JpxDecodeFilter struct{}

// NewJpxDecodeFilter creates a new JpxDecodeFilter instance.
func NewJpxDecodeFilter() *JpxDecodeFilter {
	return &JpxDecodeFilter{}
}

// IsSupported reports whether this filter is supported. Always false since JPEG 2000 decoding is not implemented.
func (f *JpxDecodeFilter) IsSupported() bool {
	return false
}

var errJpxNotSupported = errors.New("JPX filter (JPEG 2000) for image data is not currently supported; try accessing the raw compressed data directly")

// Decode returns an error because JPX decoding is not implemented.
func (f *JpxDecodeFilter) Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider FilterProvider, filterIndex int) ([]byte, error) {
	return nil, errJpxNotSupported
}

var _ Filter = (*JpxDecodeFilter)(nil)
