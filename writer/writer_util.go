package writer

import (
	"bytes"
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// seekableBuffer wraps bytes.Buffer to also implement io.Seeker.
type seekableBuffer struct{ bytes.Buffer }

func (s *seekableBuffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		return 0, nil
	case io.SeekCurrent:
		return int64(s.Len()), nil
	case io.SeekEnd:
		return int64(s.Len()), nil
	default:
		return 0, fmt.Errorf("unsupported whence: %d", whence)
	}
}

// tokenScannerForCopy defines the minimal scanner operations needed for CopyToken.
type tokenScannerForCopy interface {
	Get(reference core.IndirectReference) *tokens.ObjectToken
}

// GetOrCreateDict retrieves or creates a mutable string-keyed dictionary from the given map.
// If an existing entry is found, it follows indirect reference chains (up to 100 hops) using
// sourceScanner to resolve the underlying DictionaryToken. If no entry exists for key, a new
// empty dictionary is created and inserted.
func GetOrCreateDict[T comparable](
	dict map[T]tokens.Token,
	key T,
	sourceScanner tokenScannerForCopy,
) (map[string]tokens.Token, error) {
	item, ok := dict[key]
	if !ok {
		created := make(map[string]tokens.Token)
		dictToken, err := tokens.WithMap(created)
		if err != nil {
			return nil, fmt.Errorf("cannot create dictionary token: %w", err)
		}
		dict[key] = dictToken
		return created, nil
	}

	itemChain := item
	for chainCount := 0; chainCount < 100; chainCount++ {
		if ir, ok := itemChain.(*tokens.IndirectReferenceToken); ok && sourceScanner != nil {
			objToken := sourceScanner.Get(ir.Data())
			if objToken == nil {
				break
			}
			itemChain = objToken.Data()
		} else {
			break
		}
	}

	dt, ok := itemChain.(*tokens.DictionaryToken)
	if !ok {
		return nil, fmt.Errorf("while trying to copy token called %v which should have been a dictionary token we found a token of type %T", key, itemChain)
	}

	mutable := dt.Data()
	dict[key] = dt
	return mutable, nil
}

// CopyToken copies a token from the source document into the writer's stream, resolving all
// indirect references. Previously copied references are tracked in referencesFromDocument to
// avoid duplication. callstack tracks in-progress copies to handle cyclic references.
func CopyToken(
	writer PdfStreamWriter,
	tokenToCopy tokens.Token,
	tokenScanner tokenScannerForCopy,
	referencesFromDocument map[core.IndirectReference]*tokens.IndirectReferenceToken,
	callstack map[core.IndirectReference]*tokens.IndirectReferenceToken,
) tokens.Token {
	if callstack == nil {
		callstack = make(map[core.IndirectReference]*tokens.IndirectReferenceToken)
	}

	switch t := tokenToCopy.(type) {
	case *tokens.DictionaryToken:
		newContent := make(map[string]tokens.Token, len(t.Data()))
		for k, v := range t.Data() {
			name := tokens.Create(k)
			copied := CopyToken(writer, v, tokenScanner, referencesFromDocument, callstack)
			newContent[name.Data()] = copied
		}
		dictToken, err := tokens.WithMap(newContent)
		if err != nil {
			return tokenToCopy
		}
		return dictToken

	case *tokens.ArrayToken:
		newArray := make([]tokens.Token, 0, t.Length())
		for _, token := range t.Data() {
			newArray = append(newArray, CopyToken(writer, token, tokenScanner, referencesFromDocument, callstack))
		}
		return tokens.NewArrayToken(newArray)

	case *tokens.IndirectReferenceToken:
		ref := t.Data()
		if newRef, ok := referencesFromDocument[ref]; ok {
			return newRef
		}

		if existing, ok := callstack[ref]; ok && existing == nil {
			newRef := writer.ReserveObjectNumber()
			callstack[ref] = newRef
			referencesFromDocument[ref] = newRef
			return newRef
		}

		callstack[ref] = nil

		objToken := tokenScanner.Get(ref)
		if objToken == nil {
			return nil
		}

		tokenObject := objToken.Data()
		result := CopyToken(writer, tokenObject, tokenScanner, referencesFromDocument, callstack)

		if callstack[ref] != nil {
			return writer.WriteTokenAt(result, callstack[ref])
		}

		newRef := writer.WriteToken(result)
		referencesFromDocument[ref] = newRef
		return newRef

	case *tokens.StreamToken:
		properties := CopyToken(writer, t.StreamDictionary, tokenScanner, referencesFromDocument, callstack)
		if dictToken, ok := properties.(*tokens.DictionaryToken); ok {
			bytes := t.Data()
			streamToken, err := tokens.NewStreamToken(dictToken, bytes)
			if err != nil {
				return tokenToCopy
			}
			return streamToken
		}
		return tokenToCopy

	case *tokens.ObjectToken:
		return tokenToCopy

	default:
		return tokenToCopy
	}
}

// WalkTreeResult holds a node dictionary, its page number, and its ancestor dictionaries from the page tree walk.
type WalkTreeResult struct {
	Node       *tokens.DictionaryToken
	PageNumber *int
	Parents    []*tokens.DictionaryToken
}

// WalkTree recursively walks a PageTreeNode, yielding each leaf page along with its
// chain of parent /Pages dictionaries. For a page node it returns immediately; for a
// pages container it descends into all children and accumulates the ancestor path.
func WalkTree(node *content.PageTreeNode) []WalkTreeResult {
	var results []WalkTreeResult
	walkInner(node, nil, &results)
	return results
}

func walkInner(node *content.PageTreeNode, parents []*tokens.DictionaryToken, results *[]WalkTreeResult) {
	if node == nil {
		return
	}

	if node.IsPage {
		*results = append(*results, WalkTreeResult{
			Node:       node.NodeDictionary,
			PageNumber: node.PageNumber(),
			Parents:    parents,
		})
		return
	}

	newParents := make([]*tokens.DictionaryToken, len(parents)+1)
	copy(newParents, parents)
	newParents[len(parents)] = node.NodeDictionary

	for _, child := range node.Children() {
		walkInner(child, newParents, results)
	}
}
