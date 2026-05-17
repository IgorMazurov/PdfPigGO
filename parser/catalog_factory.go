package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CreateCatalog builds a Catalog from the root catalog dictionary.
// It validates the /Type entry, resolves the /Pages reference, creates the page tree,
// and reads named destinations. The type parameter T determines the page factory type.
func CreateCatalog[T any](
	rootReference core.IndirectReference,
	dictionary *tokens.DictionaryToken,
	scanner tokenization.PdfTokenScanner,
	pageFactory content.PageFactory[T],
	log logging.Log,
	isLenientParsing bool,
) (*content.Catalog, error) {
	if dictionary == nil {
		return nil, fmt.Errorf("dictionary cannot be nil")
	}

	typeTok, hasType := dictionary.TryGet(tokens.Type)
	if hasType && !isLenientParsing {
		if nameTok, ok := typeTok.(*tokens.NameToken); ok {
			if !nameTok.Equals(tokens.Catalog) {
				return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("The type of the catalog dictionary was not Catalog: %s.", dictionary))
			}
		}
	}

	pagesValue, hasPages := dictionary.TryGet(tokens.Pages)
	if !hasPages {
		return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("No pages entry was found in the catalog dictionary: %s.", dictionary))
	}

	pagesDictionary := resolveToDictionary(pagesValue, scanner)
	pagesReference := rootReference

	if refTok, ok := pagesValue.(*tokens.IndirectReferenceToken); ok {
		pagesReference = refTok.Data()
		if pagesDictionary == nil {
			pagesDictionary = resolveToDictionaryByRef(refTok.Data(), scanner)
		}
	} else if dictTok, ok := pagesValue.(*tokens.DictionaryToken); ok {
		pagesDictionary = dictTok
	}

	if pagesDictionary == nil {
		if isLenientParsing {
			var err error
			pagesDictionary, err = tokens.WithMap(make(map[string]tokens.Token))
			if err != nil {
				return nil, fmt.Errorf("failed to create empty pages dictionary: %w", err)
			}
		} else {
			return nil, core.NewPdfDocumentFormatException("Pages entry is null")
		}
	}

	pages, err := content.CreatePages[T](pagesReference, pagesDictionary, scanner, pageFactory, log, isLenientParsing)
	if err != nil {
		return nil, fmt.Errorf("failed to create pages: %w", err)
	}

	namedDestinations := destinations.Read(dictionary, scanner, pages, log)

	catalog, err := content.NewCatalog(dictionary, pages, namedDestinations)
	if err != nil {
		return nil, fmt.Errorf("failed to create catalog: %w", err)
	}

	return catalog, nil
}

// resolveToDictionary resolves a token to a *tokens.DictionaryToken, following indirect references.
func resolveToDictionary(token tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.DictionaryToken {
	if dict, ok := token.(*tokens.DictionaryToken); ok {
		return dict
	}

	if ref, ok := token.(*tokens.IndirectReferenceToken); ok {
		obj := scanner.Get(ref.Data())
		if obj == nil {
			return nil
		}
		return resolveToDictionary(obj.Data(), scanner)
	}

	return nil
}

// resolveToDictionaryByRef resolves an IndirectReference to a *tokens.DictionaryToken via the scanner.
func resolveToDictionaryByRef(ref core.IndirectReference, scanner tokenization.PdfTokenScanner) *tokens.DictionaryToken {
	obj := scanner.Get(ref)
	if obj == nil {
		return nil
	}

	data := obj.Data()
	if _, isNull := any(data).(*tokens.NullToken); isNull {
		return nil
	}

	if dict, ok := any(data).(*tokens.DictionaryToken); ok {
		return dict
	}

	if nestedRef, ok := any(data).(*tokens.IndirectReferenceToken); ok {
		return resolveToDictionaryByRef(nestedRef.Data(), scanner)
	}

	return nil
}
