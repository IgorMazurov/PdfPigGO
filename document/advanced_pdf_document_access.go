// Package document provides types for accessing PDF document-level features.
package document

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/tokenization"
)

// AdvancedPdfDocumentAccess provides access to rare or advanced features from the PDF specification.
type AdvancedPdfDocumentAccess struct {
	scanner    tokenization.PdfTokenScanner
	filterProv content.LookupFilterProvider
	catalog    *content.Catalog
	disposed   bool
}

// NewAdvancedPdfDocumentAccess creates a new AdvancedPdfDocumentAccess.
func NewAdvancedPdfDocumentAccess(
	scanner tokenization.PdfTokenScanner,
	filterProv content.LookupFilterProvider,
	catalog *content.Catalog,
) (*AdvancedPdfDocumentAccess, error) {
	if scanner == nil {
		return nil, fmt.Errorf("scanner cannot be nil")
	}

	if filterProv == nil {
		return nil, fmt.Errorf("filter provider cannot be nil")
	}

	if catalog == nil {
		return nil, fmt.Errorf("catalog cannot be nil")
	}

	return &AdvancedPdfDocumentAccess{
		scanner:    scanner,
		filterProv: filterProv,
		catalog:    catalog,
	}, nil
}

// TryGetEmbeddedFiles gets any embedded files contained in this PDF document.
// Since PDF 1.3 any external file referenced by the document may have its contents
// embedded within the referring PDF file, allowing its contents to be stored or
// transmitted along with the PDF file.
func (a *AdvancedPdfDocumentAccess) TryGetEmbeddedFiles() ([]*content.EmbeddedFile, bool, error) {
	if err := a.guardDisposed(); err != nil {
		return nil, false, err
	}

	catalogDict := a.catalog.CatalogDictionary
	if catalogDict == nil {
		return nil, false, nil
	}

	namesToken, ok := catalogDict.TryGet(tokens.Names)
	if !ok {
		return nil, false, nil
	}

	namesDict, ok := parts.TryGet[*tokens.DictionaryToken](namesToken, a.scanner)
	if !ok {
		return nil, false, nil
	}

	embeddedFilesToken, ok := namesDict.TryGet(tokens.EmbeddedFiles)
	if !ok {
		return nil, false, nil
	}

	embeddedFileNamesDict, ok := parts.TryGet[*tokens.DictionaryToken](embeddedFilesToken, a.scanner)
	if !ok {
		return nil, false, nil
	}

	embeddedFileNames := parts.FlattenNameTreeToDictionary[tokens.Token](embeddedFileNamesDict, a.scanner, func(t tokens.Token) tokens.Token {
		return t
	})

	if len(embeddedFileNames) == 0 {
		return nil, false, nil
	}

	result := make([]*content.EmbeddedFile, 0, len(embeddedFileNames))

	for name, value := range embeddedFileNames {
		fileDescDict, resolved := parts.TryGet[*tokens.DictionaryToken](value, a.scanner)
		if !resolved {
			continue
		}

		efToken, ok := fileDescDict.TryGet(tokens.Ef)
		if !ok {
			continue
		}

		efDict, ok := parts.TryGet[*tokens.DictionaryToken](efToken, a.scanner)
		if !ok {
			continue
		}

		fileStreamToken, ok := efDict.TryGet(tokens.F)
		if !ok {
			continue
		}

		stream, ok := parts.TryGet[*tokens.StreamToken](fileStreamToken, a.scanner)
		if !ok {
			continue
		}

		fileSpecification := ""
		if fToken, ok := fileDescDict.TryGet(tokens.F); ok {
			if dataToken, ok := fToken.(interface{ Data() string }); ok {
				fileSpecification = dataToken.Data()
			}
		}

		fileBytes := a.filterProv.DecodeStream(stream, a.scanner)

		embeddedFile, err := content.NewEmbeddedFile(name, fileSpecification, fileBytes, stream)
		if err != nil {
			continue
		}

		result = append(result, embeddedFile)
	}

	if len(result) == 0 {
		return nil, false, nil
	}

	return result, true, nil
}

// ReplaceIndirectObjectFunc is a function that takes an existing token and returns a replacement.
type ReplaceIndirectObjectFunc func(tokens.Token) tokens.Token

// ReplaceIndirectObjectWithFunc replaces the token in an internal cache that will be returned
// instead of scanning the source PDF data for future requests. The replacer function receives
// the existing token's data and returns the new token to use as replacement.
func (a *AdvancedPdfDocumentAccess) ReplaceIndirectObjectWithFunc(reference core.IndirectReference, replacer ReplaceIndirectObjectFunc) error {
	if err := a.guardDisposed(); err != nil {
		return err
	}

	obj := a.scanner.Get(reference)
	if obj == nil {
		return fmt.Errorf("object with reference %v not found", reference)
	}

	replacement := replacer(obj.Data())
	a.scanner.ReplaceToken(reference, replacement)
	return nil
}

// ReplaceIndirectObject replaces the token in an internal cache that will be returned
// instead of scanning the source PDF data for future requests.
func (a *AdvancedPdfDocumentAccess) ReplaceIndirectObject(reference core.IndirectReference, replacement tokens.Token) error {
	if err := a.guardDisposed(); err != nil {
		return err
	}

	a.scanner.ReplaceToken(reference, replacement)
	return nil
}

func (a *AdvancedPdfDocumentAccess) guardDisposed() error {
	if a.disposed {
		return fmt.Errorf("%w: %s", ErrDisposed, "AdvancedPdfDocumentAccess")
	}
	return nil
}

// Close releases resources held by this access object.
func (a *AdvancedPdfDocumentAccess) Close() error {
	a.disposed = true
	if a.scanner != nil {
		return a.scanner.Close()
	}
	return nil
}

// ErrDisposed is the error returned when operating on a disposed AdvancedPdfDocumentAccess.
var ErrDisposed = fmt.Errorf("advanced PDF document access has been disposed")
