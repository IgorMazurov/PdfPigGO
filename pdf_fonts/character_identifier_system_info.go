package pdffonts

import "fmt"

// CidFontSystemInfo provides access to the character collection definition
// for a CID font (registry, ordering). Defined here to be shared between
// fonts/cidfonts and pdf_fonts/parser/handlers packages without import cycles.
type CidFontSystemInfo interface {
	RegistryName() string
	OrderingName() string
	String() string
}

// CharacterIdentifierSystemInfo specifies the character collection associated
// with a CID font (CIDFont).
type CharacterIdentifierSystemInfo struct {
	Registry   string
	Ordering   string
	Supplement int
}

var _ CidFontSystemInfo = CharacterIdentifierSystemInfo{}

// NewCharacterIdentifierSystemInfo creates a new CharacterIdentifierSystemInfo.
func NewCharacterIdentifierSystemInfo(registry, ordering string, supplement int) CharacterIdentifierSystemInfo {
	return CharacterIdentifierSystemInfo{
		Registry:   registry,
		Ordering:   ordering,
		Supplement: supplement,
	}
}

// RegistryName returns the registry name.
func (i CharacterIdentifierSystemInfo) RegistryName() string {
	return i.Registry
}

// OrderingName returns the ordering name.
func (i CharacterIdentifierSystemInfo) OrderingName() string {
	return i.Ordering
}

// String returns the registry-ordering-supplement representation.
func (i CharacterIdentifierSystemInfo) String() string {
	return fmt.Sprintf("%s-%s-%d", i.Registry, i.Ordering, i.Supplement)
}
