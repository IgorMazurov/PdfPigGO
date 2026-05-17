package cmap

// CharacterIdentifierSystemInfoBuilder builds a CID system info by tracking
// which fields have been explicitly set.
type CharacterIdentifierSystemInfoBuilder struct {
	registry      *string
	ordering      *string
	supplement    int
	hasRegistry   bool
	hasOrdering   bool
	hasSupplement bool
}

// NewCharacterIdentifierSystemInfoBuilder creates a new builder instance.
func NewCharacterIdentifierSystemInfoBuilder() *CharacterIdentifierSystemInfoBuilder {
	return &CharacterIdentifierSystemInfoBuilder{}
}

// Registry gets the registry value. Returns "" if not set.
func (b *CharacterIdentifierSystemInfoBuilder) Registry() string {
	if b.registry == nil {
		return ""
	}
	return *b.registry
}

// SetRegistry sets the registry and marks it as present.
func (b *CharacterIdentifierSystemInfoBuilder) SetRegistry(value string) {
	b.registry = &value
	b.hasRegistry = true
}

// HasRegistry reports whether Registry has been set.
func (b *CharacterIdentifierSystemInfoBuilder) HasRegistry() bool {
	return b.hasRegistry
}

// Ordering gets the ordering value. Returns "" if not set.
func (b *CharacterIdentifierSystemInfoBuilder) Ordering() string {
	if b.ordering == nil {
		return ""
	}
	return *b.ordering
}

// SetOrdering sets the ordering and marks it as present.
func (b *CharacterIdentifierSystemInfoBuilder) SetOrdering(value string) {
	b.ordering = &value
	b.hasOrdering = true
}

// HasOrdering reports whether Ordering has been set.
func (b *CharacterIdentifierSystemInfoBuilder) HasOrdering() bool {
	return b.hasOrdering
}

// Supplement gets the supplement value. Returns 0 if not set.
func (b *CharacterIdentifierSystemInfoBuilder) Supplement() int {
	return b.supplement
}

// SetSupplement sets the supplement and marks it as present.
func (b *CharacterIdentifierSystemInfoBuilder) SetSupplement(value int) {
	b.supplement = value
	b.hasSupplement = true
}

// HasSupplement reports whether Supplement has been set.
func (b *CharacterIdentifierSystemInfoBuilder) HasSupplement() bool {
	return b.hasSupplement
}
