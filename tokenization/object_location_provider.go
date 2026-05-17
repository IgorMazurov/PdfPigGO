package tokenization

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// ObjectLocationProvider provides methods for tracking object locations
// and caching parsed object tokens during PDF scanning.
type ObjectLocationProvider interface {
	// TryGetOffset attempts to retrieve the cross-reference location for the given reference.
	// Returns the location and true if found, or zero XrefLocation and false otherwise.
	TryGetOffset(reference core.IndirectReference) (core.XrefLocation, bool)

	// UpdateOffset stores or updates the cross-reference location for a reference.
	UpdateOffset(reference core.IndirectReference, offset core.XrefLocation)

	// TryGetCached attempts to retrieve a cached object token for the given reference.
	// Returns the token and true if found, or nil and false otherwise.
	TryGetCached(reference core.IndirectReference) (*tokens.ObjectToken, bool)

	// Cache stores an object token in the cache.
	// When force is true, the token is cached regardless of existing entries.
	Cache(objectToken *tokens.ObjectToken, force bool)
}
