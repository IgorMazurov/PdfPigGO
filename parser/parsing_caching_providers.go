package parser

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/content"
)

// ParsingCachingProviders holds document-scoped caching providers.
type ParsingCachingProviders struct {
	ResourceContainer content.ResourceStore
}

// NewParsingCachingProviders creates a new ParsingCachingProviders with the given resource store.
func NewParsingCachingProviders(resourceContainer content.ResourceStore) (*ParsingCachingProviders, error) {
	if resourceContainer == nil {
		return nil, errors.New("resourceContainer cannot be nil")
	}

	return &ParsingCachingProviders{
		ResourceContainer: resourceContainer,
	}, nil
}
