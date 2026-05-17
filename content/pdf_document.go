// Package content provides types for accessing PDF document structure and catalog data.
package content

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/acroforms"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/encryption"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AdvancedAccess provides access to rare or advanced features of the PDF specification.
// This interface breaks an import cycle with the document package.
type AdvancedAccess interface {
	Close() error
}

// BookmarksRetriever retrieves bookmarks from a catalog.
// This interface breaks an import cycle with the outline package.
type BookmarksRetriever interface {
	GetBookmarks(catalog *Catalog, allowContainerNode bool) (any, error)
}

// PdfDocument provides access to document level information for a PDF document
// as well as access to the pages contained in the document.
type PdfDocument struct {
	isDisposed           bool
	documentForm         func() (*acroforms.AcroForm, error)
	version              *HeaderVersion
	inputBytes           core.InputBytes
	encryptionDictionary *encryption.EncryptionDictionary
	pdfScanner           tokenization.PdfTokenScanner
	filterProvider       LookupFilterProvider
	bookmarksRetriever   BookmarksRetriever
	parsingOptions       *ParsingOptions
	pages                *Pages
	namedDestinations    any

	// Information is the metadata associated with this document.
	Information *DocumentInformation

	// Structure provides access to the underlying raw structure of the document.
	Structure *Structure

	// Advanced provides access to rare or advanced features of the PDF specification.
	Advanced AdvancedAccess
}

// NewPdfDocument creates a new PdfDocument instance from the given components.
func NewPdfDocument(
	inputBytes core.InputBytes,
	version *HeaderVersion,
	catalog *Catalog,
	information *DocumentInformation,
	encryptionDictionary *encryption.EncryptionDictionary,
	pdfScanner tokenization.PdfTokenScanner,
	filterProvider LookupFilterProvider,
	acroFormFactory *acroforms.AcroFormFactory,
	bookmarksRetriever BookmarksRetriever,
	parsingOptions *ParsingOptions,
	advanced AdvancedAccess,
) (*PdfDocument, error) {
	if version == nil {
		return nil, fmt.Errorf("version cannot be null")
	}
	if pdfScanner == nil {
		return nil, fmt.Errorf("pdf scanner cannot be null")
	}
	if filterProvider == nil {
		return nil, fmt.Errorf("filter provider cannot be null")
	}
	if bookmarksRetriever == nil {
		return nil, fmt.Errorf("bookmarks retriever cannot be null")
	}
	if information == nil {
		return nil, fmt.Errorf("information cannot be null")
	}

	structure, err := NewStructure(catalog, pdfScanner)
	if err != nil {
		return nil, fmt.Errorf("failed to create structure: %w", err)
	}

	docFormFactory := func() (*acroforms.AcroForm, error) {
		if acroFormFactory == nil {
			return nil, nil
		}
		return acroFormFactory.GetAcroForm(&acroformCatalog{catalog: catalog})
	}

	doc := &PdfDocument{
		inputBytes:           inputBytes,
		version:              version,
		encryptionDictionary: encryptionDictionary,
		pdfScanner:           pdfScanner,
		filterProvider:       filterProvider,
		bookmarksRetriever:   bookmarksRetriever,
		parsingOptions:       parsingOptions,
		pages:                catalog.Pages(),
		namedDestinations:    catalog.NamedDestinations(),
		Information:          information,
		Structure:            structure,
		Advanced:             advanced,
		documentForm:         docFormFactory,
	}

	return doc, nil
}

// Version returns the version number of the PDF specification which this file conforms to.
func (d *PdfDocument) Version() float64 {
	if d.version == nil {
		return 0
	}
	return d.version.Version
}

// NumberOfPages returns the number of pages in this document.
func (d *PdfDocument) NumberOfPages() int {
	if d.pages == nil {
		return 0
	}
	return d.pages.Count()
}

// IsEncrypted indicates whether the document content is encrypted.
func (d *PdfDocument) IsEncrypted() bool {
	return d.encryptionDictionary != nil
}

// GetPage returns the page with the specified page number (1 indexed).
func (d *PdfDocument) GetPage(pageNumber int) (any, error) {
	if d.isDisposed {
		return nil, fmt.Errorf("cannot access page after the document is disposed")
	}

	if d.parsingOptions != nil && d.parsingOptions.Logger != nil {
		d.parsingOptions.Logger.Debug(fmt.Sprintf("Accessing page %d.", pageNumber))
	}

	page, err := d.pages.GetPage(pageNumber, d.namedDestinations, d.parsingOptions)
	if err != nil {
		if d.IsEncrypted() {
			return nil, encryption.NewPdfDocumentEncryptedExceptionWithDictAndInner(
				"document was encrypted which may have caused error when retrieving page",
				d.encryptionDictionary,
				err,
			)
		}
		return nil, err
	}

	return page, nil
}

// AddPageFactory registers an additional page factory for creating typed pages.
func (d *PdfDocument) AddPageFactory(factory any) {
	if d.pages != nil && factory != nil {
		d.pages.AddPageFactory(factory)
	}
}

// GetTypedPage returns the page with the specified page number (1 indexed) using a factory that produces type T.
// The factory must have been registered via AddPageFactory and implement TypedPageFactory.
// Note: Go does not support generic methods on structs, so this is a standalone function.
func GetTypedPage[T any](d *PdfDocument, pageNumber int) (*T, error) {
	if d.isDisposed {
		return nil, fmt.Errorf("cannot access page after the document is disposed")
	}
	return GetPage[T](d.pages, pageNumber, d.namedDestinations, d.parsingOptions)
}

// GetPages returns all pages in this document in order.
func (d *PdfDocument) GetPages() ([]any, error) {
	result := make([]any, 0, d.NumberOfPages())
	for i := 1; i <= d.NumberOfPages(); i++ {
		page, err := d.GetPage(i)
		if err != nil {
			return result, err
		}
		result = append(result, page)
	}
	return result, nil
}

// TryGetXmpMetadata gets the document level metadata if present.
// The metadata is XML in the XMP (Extensible Metadata Platform) format.
func (d *PdfDocument) TryGetXmpMetadata() (*XmpMetadata, bool, error) {
	if d.isDisposed {
		return nil, false, fmt.Errorf("cannot access the document metadata after the document is disposed")
	}

	catalogDict := d.Structure.Catalog().CatalogDictionary
	if catalogDict == nil {
		return nil, false, nil
	}

	xmpStreamTokenRaw, ok := catalogDict.TryGet(tokens.Metadata)
	if !ok {
		return nil, false, nil
	}

	var streamToken *tokens.StreamToken
	switch t := xmpStreamTokenRaw.(type) {
	case *tokens.StreamToken:
		streamToken = t
	case *tokens.IndirectReferenceToken:
		obj, err := d.Structure.GetObject(t.Data())
		if err != nil || obj == nil {
			return nil, false, nil
		}
		streamToken, ok = obj.Data().(*tokens.StreamToken)
		if !ok {
			return nil, false, nil
		}
	default:
		return nil, false, nil
	}

	metadata, err := NewXmpMetadata(streamToken, d.filterProvider, d.pdfScanner)
	if err != nil {
		return nil, false, err
	}

	return metadata, true, nil
}

// TryGetBookmarks gets the bookmarks if this document contains some.
func (d *PdfDocument) TryGetBookmarks(allowContainerNode bool) (any, bool, error) {
	if d.isDisposed {
		return nil, false, fmt.Errorf("cannot access the bookmarks after the document is disposed")
	}

	catalog := d.Structure.Catalog()
	if catalog == nil {
		return nil, false, nil
	}

	bookmarks, err := d.bookmarksRetriever.GetBookmarks(catalog, allowContainerNode)
	if err != nil {
		return nil, false, err
	}

	if bookmarks == nil {
		return nil, false, nil
	}

	return bookmarks, true, nil
}

// TryGetForm gets the form if this document contains one.
func (d *PdfDocument) TryGetForm() (*acroforms.AcroForm, bool, error) {
	if d.isDisposed {
		return nil, false, fmt.Errorf("cannot access the form after the document is disposed")
	}

	if d.documentForm == nil {
		return nil, false, nil
	}

	form, err := d.documentForm()
	if err != nil {
		return nil, false, err
	}

	if form == nil {
		return nil, false, nil
	}

	return form, true, nil
}

// Close releases resources held by the document and closes any unmanaged resources.
func (d *PdfDocument) Close() error {
	d.isDisposed = true

	var lastErr error

	if d.Advanced != nil {
		if err := d.Advanced.Close(); err != nil {
			lastErr = err
		}
	}

	if d.pdfScanner != nil {
		if err := d.pdfScanner.Close(); err != nil {
			lastErr = err
		}
	}

	if closer, ok := any(d.inputBytes).(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			lastErr = err
		}
	}

	if lastErr != nil && d.parsingOptions != nil && d.parsingOptions.Logger != nil {
		d.parsingOptions.Logger.Error(fmt.Sprintf("Failed disposing the PdfDocument due to an error: %v", lastErr))
	}

	return lastErr
}

// acroformCatalog adapts *Catalog to acroforms.Catalog interface.
type acroformCatalog struct {
	catalog *Catalog
}

func (a *acroformCatalog) CatalogDictionary() *tokens.DictionaryToken {
	if a.catalog == nil {
		return nil
	}
	return a.catalog.CatalogDictionary
}

func (a *acroformCatalog) GetPageNumberByReference(ref core.IndirectReference) *int {
	if a.catalog == nil || a.catalog.Pages() == nil {
		return nil
	}
	pageRef := a.catalog.Pages().GetPageByReference(ref)
	if pageRef == nil {
		return nil
	}
	return pageRef.PageNumber()
}

