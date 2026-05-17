package filters

import (
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Filter represents a PDF stream filter for decoding compressed data.
type Filter interface {
	// IsSupported reports whether this filter is supported by the current build.
	IsSupported() bool

	// Decode decodes the given input bytes using this filter.
	Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider FilterProvider, filterIndex int) ([]byte, error)
}

// FilterProvider provides filter implementations for decoding PDF data.
// Corresponds to C# IFilterProvider.
type FilterProvider interface {
	// GetFilters returns the filters from the /Filter or /F entry in the stream dictionary.
	GetFilters(dictionary *tokens.DictionaryToken) ([]Filter, error)

	// GetNamedFilters returns the filters corresponding to the given filter names.
	GetNamedFilters(names []*tokens.NameToken) ([]Filter, error)

	// GetAllFilters returns all registered filters in the provider.
	GetAllFilters() []Filter
}
