package filters

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// Jbig2DecodeFilter represents the JBIG2 filter for monochrome image data.
// This filter is not implemented and will not be used during parsing.
type Jbig2DecodeFilter struct{}

// NewJbig2DecodeFilter creates a new Jbig2DecodeFilter instance.
func NewJbig2DecodeFilter() *Jbig2DecodeFilter {
	return &Jbig2DecodeFilter{}
}

// IsSupported reports whether this filter is supported. Always false since JBIG2 decoding is not implemented.
func (f *Jbig2DecodeFilter) IsSupported() bool {
	return false
}

var errJbig2NotSupported = errors.New("JBIG2 filter for monochrome image data is not currently supported; try accessing the raw compressed data directly")

// Decode returns an error because JBIG2 decoding is not implemented.
func (f *Jbig2DecodeFilter) Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider FilterProvider, filterIndex int) ([]byte, error) {
	return nil, errJbig2NotSupported
}

var _ Filter = (*Jbig2DecodeFilter)(nil)
