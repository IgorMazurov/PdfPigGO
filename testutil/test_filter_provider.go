package testutil

import (
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// TestFilterProvider is a no-op implementation of LookupFilterProvider for tests.
// All methods return empty filter lists.
var _ filters.LookupFilterProvider = (*TestFilterProvider)(nil)

type TestFilterProvider struct{}

// Instance is the singleton shared test filter provider.
var Instance = &TestFilterProvider{}

func (p *TestFilterProvider) GetFilters(_ *tokens.DictionaryToken) ([]filters.Filter, error) {
	return nil, nil
}

func (p *TestFilterProvider) GetNamedFilters(_ []*tokens.NameToken) ([]filters.Filter, error) {
	return nil, nil
}

func (p *TestFilterProvider) GetAllFilters() []filters.Filter {
	return nil
}

func (p *TestFilterProvider) GetFiltersWithScanner(_ *tokens.DictionaryToken, _ tokenization.PdfTokenScanner) ([]filters.Filter, error) {
	return nil, nil
}
