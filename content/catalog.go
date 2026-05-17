package content

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Catalog is the root of the document's object hierarchy. It contains
// references to objects defining the contents, outline, named destinations
// and more.
type Catalog struct {
	// CatalogDictionary is the catalog dictionary containing assorted information.
	CatalogDictionary *tokens.DictionaryToken

	pages             *Pages
	namedDestinations *destinations.NamedDestinations
}

// NewCatalog creates a new Catalog instance.
func NewCatalog(catalogDictionary *tokens.DictionaryToken, pages *Pages, namedDestinations *destinations.NamedDestinations) (*Catalog, error) {
	if catalogDictionary == nil {
		return nil, errors.New("catalog dictionary cannot be nil")
	}

	if pages == nil {
		return nil, errors.New("pages cannot be nil")
	}

	return &Catalog{
		CatalogDictionary: catalogDictionary,
		pages:             pages,
		namedDestinations: namedDestinations,
	}, nil
}

// Pages returns the page tree for this catalog.
func (c *Catalog) Pages() *Pages {
	return c.pages
}

// NamedDestinations returns the named destinations in the document.
func (c *Catalog) NamedDestinations() *destinations.NamedDestinations {
	return c.namedDestinations
}
