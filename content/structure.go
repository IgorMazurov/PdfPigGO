// Package content provides types for accessing PDF document structure and catalog data.
package content

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Structure provides access to explore and retrieve the underlying PDF objects from the document.
type Structure struct {
	catalog    *Catalog
	tokenScanner tokenization.PdfTokenScanner
}

// NewStructure creates a new Structure instance.
func NewStructure(catalog *Catalog, scanner tokenization.PdfTokenScanner) (*Structure, error) {
	if catalog == nil {
		return nil, core.NewPdfDocumentFormatException("catalog cannot be null")
	}
	if scanner == nil {
		return nil, core.NewPdfDocumentFormatException("scanner cannot be null")
	}
	return &Structure{
		catalog:    catalog,
		tokenScanner: scanner,
	}, nil
}

// Catalog returns the root of the document's hierarchy providing access to the page tree.
func (s *Structure) Catalog() *Catalog {
	return s.catalog
}

// GetObject retrieves the tokenized object with the specified reference number.
func (s *Structure) GetObject(reference core.IndirectReference) (*tokens.ObjectToken, error) {
	obj := s.tokenScanner.Get(reference)
	if obj == nil {
		return nil, core.NewPdfDocumentFormatException("could not find the object with reference: " + reference.String())
	}
	return obj, nil
}

// Scanner returns the underlying token scanner for advanced operations.
func (s *Structure) Scanner() tokenization.PdfTokenScanner {
	return s.tokenScanner
}
