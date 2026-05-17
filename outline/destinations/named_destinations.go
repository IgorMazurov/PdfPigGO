package destinations

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// PageRef represents a reference to a single page within the document.
type PageRef interface {
	// PageNumber returns the 1-based page number, or nil if no valid number exists.
	PageNumber() *int
}

// Pages is the minimal interface for page collections needed by destination resolution.
// Implemented by content.Pages without creating an import cycle.
type Pages interface {
	Count() int
	// GetPageByReference returns the page for the given indirect reference, or nil if not found.
	GetPageByReference(ref core.IndirectReference) PageRef
}

// NamedDestinations holds named destination mappings within a PDF document.
type NamedDestinations struct {
	namedDestinations map[string]ExplicitDestination
	pages             Pages
}

// NewNamedDestinations creates a new NamedDestinations instance with the given
// destination map and page collection.
func NewNamedDestinations(
	namedDestinations map[string]ExplicitDestination,
	pages Pages,
) *NamedDestinations {
	return &NamedDestinations{
		namedDestinations: namedDestinations,
		pages:             pages,
	}
}

// TryGet looks up an explicit destination by its string name key.
func (n *NamedDestinations) TryGet(name string) (ExplicitDestination, bool) {
	dest, ok := n.namedDestinations[name]
	return dest, ok
}

// TryGetName looks up an explicit destination by its string name key.
// Alias for TryGet, used by DestinationProvider.
func (n *NamedDestinations) TryGetName(name string) (ExplicitDestination, bool) {
	return n.TryGet(name)
}

// TryGetExplicitDestination attempts to parse an explicit destination from an array token.
func (n *NamedDestinations) TryGetExplicitDestination(
	array *tokens.ArrayToken,
	log logging.Log,
	isRemoteDestination bool,
) (ExplicitDestination, bool) {
	return TryGetExplicitDestination(array, n.pages, log, isRemoteDestination)
}
