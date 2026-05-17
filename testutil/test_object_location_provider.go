package testutil

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

var _ tokenization.ObjectLocationProvider = (*TestObjectLocationProvider)(nil)

// TestObjectLocationProvider is a simple test implementation of ObjectLocationProvider
// that stores offsets in a map and does not cache object tokens.
type TestObjectLocationProvider struct {
	Offsets map[core.IndirectReference]core.XrefLocation
}

// NewTestObjectLocationProvider creates a new TestObjectLocationProvider with an initialized offsets map.
func NewTestObjectLocationProvider() *TestObjectLocationProvider {
	return &TestObjectLocationProvider{
		Offsets: make(map[core.IndirectReference]core.XrefLocation),
	}
}

// TryGetOffset attempts to retrieve the cross-reference location for the given reference.
func (p *TestObjectLocationProvider) TryGetOffset(reference core.IndirectReference) (core.XrefLocation, bool) {
	offset, ok := p.Offsets[reference]
	return offset, ok
}

// UpdateOffset stores or updates the cross-reference location for a reference.
func (p *TestObjectLocationProvider) UpdateOffset(reference core.IndirectReference, offset core.XrefLocation) {
	p.Offsets[reference] = offset
}

// TryGetCached always returns nil and false; this provider does not cache object tokens.
func (p *TestObjectLocationProvider) TryGetCached(reference core.IndirectReference) (*tokens.ObjectToken, bool) {
	return nil, false
}

// Cache is a no-op; this provider does not cache object tokens.
func (p *TestObjectLocationProvider) Cache(objectToken *tokens.ObjectToken, force bool) {
}
