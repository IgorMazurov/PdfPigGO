package tokenization

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// objectLocationProvider implements ObjectLocationProvider.
type objectLocationProvider struct {
	cache              map[core.IndirectReference]*tokens.ObjectToken
	bytes              core.InputBytes
	bruteForcedOffsets map[core.IndirectReference]core.XrefLocation
	offsets            map[core.IndirectReference]core.XrefLocation
	bruteForceFunc     func(core.InputBytes) (map[core.IndirectReference]core.XrefLocation, error)
}

// NewObjectLocationProviderImpl creates a new ObjectLocationProvider implementation.
func NewObjectLocationProviderImpl(
	xrefOffsets map[core.IndirectReference]core.XrefLocation,
	bruteForceOffsets map[core.IndirectReference]core.XrefLocation,
	inputBytes core.InputBytes,
) ObjectLocationProvider {
	offsets := make(map[core.IndirectReference]core.XrefLocation)
	for k, v := range xrefOffsets {
		offsets[k] = v
	}

	return &objectLocationProvider{
		cache:              make(map[core.IndirectReference]*tokens.ObjectToken),
		bytes:              inputBytes,
		bruteForcedOffsets: bruteForceOffsets,
		offsets:            offsets,
	}
}

// WithBruteForce sets the lazy brute-force location resolver on the provider.
func WithBruteForce(p ObjectLocationProvider, fn func(core.InputBytes) (map[core.IndirectReference]core.XrefLocation, error)) ObjectLocationProvider {
	if impl, ok := p.(*objectLocationProvider); ok {
		impl.bruteForceFunc = fn
	}
	return p
}

func (p *objectLocationProvider) TryGetOffset(reference core.IndirectReference) (core.XrefLocation, bool) {
	if p.bruteForcedOffsets != nil {
		if offset, ok := p.bruteForcedOffsets[reference]; ok {
			return offset, true
		}
	}

	if offset, ok := p.offsets[reference]; ok {
		return offset, true
	}

	if p.bruteForcedOffsets == nil && p.bruteForceFunc != nil {
		locations, err := p.bruteForceFunc(p.bytes)
		if err == nil {
			p.bruteForcedOffsets = locations
		} else {
			p.bruteForcedOffsets = make(map[core.IndirectReference]core.XrefLocation)
		}

		if offset, ok := p.bruteForcedOffsets[reference]; ok {
			return offset, true
		}
	}

	var zero core.XrefLocation
	return zero, false
}

func (p *objectLocationProvider) UpdateOffset(reference core.IndirectReference, offset core.XrefLocation) {
	p.offsets[reference] = offset
}

func (p *objectLocationProvider) TryGetCached(reference core.IndirectReference) (*tokens.ObjectToken, bool) {
	tok, ok := p.cache[reference]
	return tok, ok
}

func (p *objectLocationProvider) Cache(objectToken *tokens.ObjectToken, force bool) {
	if objectToken == nil {
		return
	}

	if !force {
		expected, exists := p.offsets[objectToken.Number()]
		if exists && (objectToken.Position().Type != expected.Type || objectToken.Position().Value1 != expected.Value1) {
			return
		}
	}

	p.cache[objectToken.Number()] = objectToken
}
