package filters

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// LookupFilterProvider extends FilterProvider with scanner-aware filter resolution.
// Corresponds to C# ILookupFilterProvider : IFilterProvider.
type LookupFilterProvider interface {
	FilterProvider

	// GetFiltersWithScanner resolves filters from the dictionary using the scanner
	// for indirect reference dereferencing.
	GetFiltersWithScanner(dictionary *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) ([]Filter, error)
}

// FilterProviderWithLookup wraps an inner FilterProvider and adds scanner-aware
// filter resolution that can follow indirect references. Corresponds to C#
// FilterProviderWithLookup : ILookupFilterProvider.
type FilterProviderWithLookup struct {
	inner FilterProvider
}

// NewFilterProviderWithLookup creates a new FilterProviderWithLookup wrapping the given provider.
func NewFilterProviderWithLookup(inner FilterProvider) *FilterProviderWithLookup {
	if inner == nil {
		return nil
	}
	return &FilterProviderWithLookup{inner: inner}
}

// GetFilters returns the filters from the /Filter or /F entry in the stream dictionary.
func (p *FilterProviderWithLookup) GetFilters(dictionary *tokens.DictionaryToken) ([]Filter, error) {
	return p.inner.GetFilters(dictionary)
}

// GetNamedFilters returns the filters corresponding to the given filter names.
func (p *FilterProviderWithLookup) GetNamedFilters(names []*tokens.NameToken) ([]Filter, error) {
	return p.inner.GetNamedFilters(names)
}

// GetAllFilters returns all registered filters in the provider.
func (p *FilterProviderWithLookup) GetAllFilters() []Filter {
	return p.inner.GetAllFilters()
}

// GetFiltersWithScanner resolves filters from the dictionary using the scanner for
// indirect reference dereferencing. The /Filter or /F entry may be a name token,
// an array of name tokens, or an indirect reference to either.
func (p *FilterProviderWithLookup) GetFiltersWithScanner(dictionary *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) ([]Filter, error) {
	if dictionary == nil {
		return nil, fmt.Errorf("dictionary cannot be nil")
	}

	token := getFilterOrFallback(dictionary)
	if token == nil {
		return []Filter{}, nil
	}

	switch t := token.(type) {
	case *tokens.ArrayToken:
		data := t.Data()
		result := make([]*tokens.NameToken, len(data))
		for i, filterToken := range data {
			nameToken, ok := filterToken.(*tokens.NameToken)
			if !ok {
				return nil, core.NewPdfDocumentFormatException(
					fmt.Sprintf("the filter for the stream was not a valid object. Expected name or array, instead got: %v", filterToken))
			}
			result[i] = nameToken
		}
		return p.GetNamedFilters(result)

	case *tokens.NameToken:
		return p.GetNamedFilters([]*tokens.NameToken{t})

	case *tokens.IndirectReferenceToken:
		resolved, ok := resolveIndirectToken(t.Data(), scanner)
		if !ok {
			return nil, core.NewPdfDocumentFormatException(
				fmt.Sprintf("the filter for the stream was not a valid object. Expected name or array, instead got: %v", token))
		}

		switch rt := resolved.(type) {
		case *tokens.NameToken:
			return p.GetNamedFilters([]*tokens.NameToken{rt})
		case *tokens.ArrayToken:
			data := rt.Data()
			result := make([]*tokens.NameToken, 0, len(data))
			for _, x := range data {
				nameToken, ok := x.(*tokens.NameToken)
				if !ok {
					return nil, core.NewPdfDocumentFormatException(
						fmt.Sprintf("the filter for the stream was not a valid object. Expected name or array, instead got: %v", x))
				}
				result = append(result, nameToken)
			}
			return p.GetNamedFilters(result)
		default:
			return nil, core.NewPdfDocumentFormatException(
				fmt.Sprintf("the filter for the stream was not a valid object. Expected name or array, instead got: %v", token))
		}

	default:
		return nil, core.NewPdfDocumentFormatException(
			fmt.Sprintf("the filter for the stream was not a valid object. Expected name or array, instead got: %v", token))
	}
}

// resolveIndirectToken resolves an indirect reference through the scanner,
// following nested indirect references recursively. Returns (nil, false) on failure.
func resolveIndirectToken(reference core.IndirectReference, scanner tokenization.PdfTokenScanner) (tokens.Token, bool) {
	obj := scanner.Get(reference)
	if obj == nil || obj.Data() == nil {
		return nil, false
	}

	data := obj.Data()

	if data == nil {
		return nil, false
	}

	switch t := data.(type) {
	case *tokens.NameToken:
		return t, true
	case *tokens.ArrayToken:
		return t, true
	case *tokens.IndirectReferenceToken:
		return resolveIndirectToken(t.Data(), scanner)
	default:
		return data, true
	}
}
