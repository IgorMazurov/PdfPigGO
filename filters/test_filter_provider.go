package filters

import (
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// TestFilterProvider is a no-op filter provider that returns empty lists for all
// lookup methods. Useful in tests where filter decoding should be skipped.
type TestFilterProvider struct{}

// TestFilterProviderInstance is the singleton shared instance of TestFilterProvider.
var TestFilterProviderInstance = &TestFilterProvider{}

// GetFilters returns an empty slice, matching the C# TestFilterProvider behavior.
func (p *TestFilterProvider) GetFilters(_ *tokens.DictionaryToken) ([]Filter, error) {
	return []Filter{}, nil
}

// GetNamedFilters returns an empty slice regardless of the names provided.
func (p *TestFilterProvider) GetNamedFilters(_ []*tokens.NameToken) ([]Filter, error) {
	return []Filter{}, nil
}

// GetAllFilters returns an empty slice since no filters are registered.
func (p *TestFilterProvider) GetAllFilters() []Filter {
	return []Filter{}
}

// GetFiltersWithScanner returns an empty slice regardless of the dictionary or scanner.
func (p *TestFilterProvider) GetFiltersWithScanner(_ *tokens.DictionaryToken, _ tokenization.PdfTokenScanner) ([]Filter, error) {
	return []Filter{}, nil
}
