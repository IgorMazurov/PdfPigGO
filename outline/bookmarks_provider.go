package outline

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/util"
)

// BookmarksProvider extracts bookmarks (outlines) from a PDF document.
type BookmarksProvider struct {
	log        logging.Log
	pdfScanner tokenization.PdfTokenScanner
}

// NewBookmarksProvider creates a new BookmarksProvider.
func NewBookmarksProvider(log logging.Log, pdfScanner tokenization.PdfTokenScanner) *BookmarksProvider {
	return &BookmarksProvider{
		log:        log,
		pdfScanner: pdfScanner,
	}
}

// GetBookmarks extracts the bookmarks (outlines) from the given catalog, if any exist.
func (p *BookmarksProvider) GetBookmarks(catalog *content.Catalog, allowContainerNode bool) (*Bookmarks, error) {
	if catalog == nil {
		return nil, fmt.Errorf("catalog cannot be nil")
	}

	outlinesDict := p.resolveDictionaryToken(catalog.CatalogDictionary, tokens.Outlines)
	if outlinesDict == nil {
		return nil, nil
	}

	typeName, hasType := outlinesDict.TryGet(tokens.Type)
	if hasType {
		if nameToken, ok := typeName.(*tokens.NameToken); ok && nameToken.Data() != tokens.Outlines.Data() {
			p.log.Error(fmt.Sprintf("Outlines (bookmarks) dictionary did not have correct type specified: %s.", nameToken.Data()))
		}
	}

	next := p.resolveDictionaryToken(outlinesDict, tokens.First)
	if next == nil {
		return nil, nil
	}

	roots := make([]BookmarkNodeTyped, 0)
	seen := make(map[core.IndirectReference]bool)
	namedDests := catalog.NamedDestinations()

	for next != nil {
		err := p.readBookmarksRecursively(next, 0, false, seen, namedDests, &roots, allowContainerNode)
		if err != nil {
			return nil, err
		}

		nextRefToken, hasNext := next.TryGet(tokens.Next)
		if !hasNext {
			break
		}

		nextRef, ok := nextRefToken.(*tokens.IndirectReferenceToken)
		if !ok {
			break
		}

		if seen[nextRef.Data()] {
			break
		}
		seen[nextRef.Data()] = true

		next = p.getDictionaryByRef(nextRef.Data())
	}

	return NewBookmarksTyped(roots), nil
}

// BookmarkNodeTyped is an interface satisfied by all concrete bookmark node types.
type BookmarkNodeTyped interface {
	bookmarkNodeTyped()
	GetBookmarkNode() *BookmarkNode
}

func (*DocumentBookmarkNode) bookmarkNodeTyped()   {}
func (*ExternalBookmarkNode) bookmarkNodeTyped()   {}
func (*UriBookmarkNode) bookmarkNodeTyped()        {}
func (*ContainerBookmarkNode) bookmarkNodeTyped()  {}
func (*EmbeddedBookmarkNode) bookmarkNodeTyped()   {}

func (n *DocumentBookmarkNode) GetBookmarkNode() *BookmarkNode    { return n.BookmarkNode }
func (n *ExternalBookmarkNode) GetBookmarkNode() *BookmarkNode    { return n.DocumentBookmarkNode.BookmarkNode }
func (n *UriBookmarkNode) GetBookmarkNode() *BookmarkNode         { return n.BookmarkNode }
func (n *ContainerBookmarkNode) GetBookmarkNode() *BookmarkNode   { return n.BookmarkNode }
func (n *EmbeddedBookmarkNode) GetBookmarkNode() *BookmarkNode    { return n.DocumentBookmarkNode.BookmarkNode }

// readBookmarksRecursively reads bookmark nodes recursively from the outline tree.
func (p *BookmarksProvider) readBookmarksRecursively(
	nodeDict *tokens.DictionaryToken,
	level int,
	readSiblings bool,
	seen map[core.IndirectReference]bool,
	namedDests *destinations.NamedDestinations,
	list *[]BookmarkNodeTyped,
	allowContainerNode bool,
) error {
	title, ok := p.tryGetStringDirect(nodeDict, tokens.Title)
	if !ok {
		return core.NewPdfDocumentFormatException(fmt.Sprintf("Invalid title for outline (bookmark) node: %s.", nodeDict))
	}

	childrenList := make([]BookmarkNodeTyped, 0)
	firstChild := p.resolveDictionaryToken(nodeDict, tokens.First)
	if firstChild != nil {
		if err := p.readBookmarksRecursively(firstChild, level+1, true, seen, namedDests, &childrenList, allowContainerNode); err != nil {
			return err
		}
	}

	children := typedToBookmarkNodes(childrenList)

	destProvider := destinations.DestinationProvider{}
	if dest, ok := destProvider.TryGetDestination(nodeDict, tokens.Dest, namedDests, p.pdfScanner, p.log, false); ok {
		docNode, err := NewDocumentBookmarkNode(title, level, dest, children)
		if err == nil && docNode != nil {
			*list = append(*list, docNode)
		}
	} else if bm, err := p.tryCreateBookmarkFromAction(nodeDict, title, level, children, namedDests, allowContainerNode); err != nil {
		return err
	} else if bm != nil {
		*list = append(*list, bm)
	} else if allowContainerNode {
		p.log.Warn(fmt.Sprintf("No /Dest(ination) or /A(ction) entry found for bookmark node: %s.", nodeDict))
		containerNode, err := NewContainerBookmarkNode(title, level, children)
		if err == nil && containerNode != nil {
			*list = append(*list, containerNode)
		}
	} else {
		p.log.Error(fmt.Sprintf("No /Dest(ination) or /A(ction) entry found for bookmark node: %s.", nodeDict))
		return nil
	}

	if !readSiblings {
		return nil
	}

	current := nodeDict
	for {
		nextRefToken, hasNext := current.TryGet(tokens.Next)
		if !hasNext {
			break
		}

		nextRef, ok := nextRefToken.(*tokens.IndirectReferenceToken)
		if !ok {
			break
		}

		if seen[nextRef.Data()] {
			break
		}
		seen[nextRef.Data()] = true

		current = p.getDictionaryByRef(nextRef.Data())
		if current == nil {
			break
		}

		if err := p.readBookmarksRecursively(current, level, false, seen, namedDests, list, allowContainerNode); err != nil {
			return err
		}
	}

	return nil
}

// tryCreateBookmarkFromAction attempts to create a bookmark node from an /A (action) entry.
// This mirrors the logic in actions.ActionProvider but returns typed bookmark nodes directly.
func (p *BookmarksProvider) tryCreateBookmarkFromAction(
	nodeDict *tokens.DictionaryToken,
	title string,
	level int,
	children []*BookmarkNode,
	namedDests *destinations.NamedDestinations,
	allowContainerNode bool,
) (BookmarkNodeTyped, error) {
	actionDict := p.resolveDictionaryToken(nodeDict, tokens.A)
	if actionDict == nil {
		return nil, nil
	}

	actionTypeToken, hasActionType := actionDict.TryGet(tokens.S)
	if !hasActionType {
		return nil, nil
	}

	actionName, ok := actionTypeToken.(*tokens.NameToken)
	if !ok {
		return nil, nil
	}

	destProvider := destinations.DestinationProvider{}

	switch actionName.Data() {
	case tokens.GoTo.Data():
		dest, ok := destProvider.TryGetDestination(actionDict, tokens.D, namedDests, p.pdfScanner, p.log, false)
		if !ok {
			return nil, nil
		}
		docNode, err := NewDocumentBookmarkNode(title, level, dest, children)
		if err != nil || docNode == nil {
			return nil, err
		}
		return docNode, nil

	case tokens.GoToR.Data():
		filename, filenameOk := util.TryGetOptionalStringDirect(actionDict, tokens.F, p.pdfScanner)
		if !filenameOk {
			return nil, nil
		}
		dest, ok := destProvider.TryGetDestination(actionDict, tokens.D, namedDests, p.pdfScanner, p.log, true)
		if !ok {
			return nil, nil
		}
		extNode, err := NewExternalBookmarkNode(title, level, dest, children, filename)
		if err != nil || extNode == nil {
			return nil, err
		}
		return extNode, nil

	case tokens.GoToE.Data():
		dest, ok := destProvider.TryGetDestination(actionDict, tokens.D, namedDests, p.pdfScanner, p.log, true)
		if !ok {
			return nil, nil
		}
		fileSpec, fileSpecOk := util.TryGetOptionalStringDirect(actionDict, tokens.F, p.pdfScanner)
		if !fileSpecOk {
			fileSpec = ""
		}
		if fileSpec == "" {
			return nil, nil
		}
		embNode, err := NewEmbeddedBookmarkNode(title, level, dest, children, fileSpec)
		if err != nil || embNode == nil {
			return nil, err
		}
		return embNode, nil

	case tokens.Uri.Data():
		uri, uriOk := util.TryGetOptionalStringDirect(actionDict, tokens.Uri, p.pdfScanner)
		if !uriOk {
			return nil, nil
		}
		uriNode, err := NewUriBookmarkNode(title, level, uri, children)
		if err != nil || uriNode == nil {
			return nil, err
		}
		return uriNode, nil

	default:
		if allowContainerNode {
			containerNode, err := NewContainerBookmarkNode(title, level, children)
			if err == nil && containerNode != nil {
				return containerNode, nil
			}
		}
		return nil, nil
	}
}

// typedToBookmarkNodes converts a slice of BookmarkNodeTyped to []*BookmarkNode.
func typedToBookmarkNodes(items []BookmarkNodeTyped) []*BookmarkNode {
	result := make([]*BookmarkNode, 0, len(items))
	for _, item := range items {
		result = append(result, item.GetBookmarkNode())
	}
	return result
}

// resolveDictionaryToken tries to get a DictionaryToken from the dictionary by name,
// resolving indirect references through the scanner using direct object lookup.
func (p *BookmarksProvider) resolveDictionaryToken(dict *tokens.DictionaryToken, name *tokens.NameToken) *tokens.DictionaryToken {
	token, ok := dict.TryGet(name)
	if !ok {
		return nil
	}

	// If it's an indirect reference, resolve it via direct lookup (matching C# DirectObjectFinder).
	if indRef, isIndirect := token.(*tokens.IndirectReferenceToken); isIndirect {
		obj := p.pdfScanner.Get(indRef.Data())
		if obj == nil {
			return nil
		}
		token = obj.Data()
	}

	result, ok := token.(*tokens.DictionaryToken)
	if !ok {
		return nil
	}
	return result
}

// tryGetStringDirect tries to get a string value directly from the dictionary without
// resolving indirect references. Handles both StringToken and HexToken.
func (p *BookmarksProvider) tryGetStringDirect(dict *tokens.DictionaryToken, name *tokens.NameToken) (string, bool) {
	token, ok := dict.TryGet(name)
	if !ok {
		return "", false
	}

	switch t := token.(type) {
	case *tokens.StringToken:
		return t.Data(), true
	case *tokens.HexToken:
		return t.Data(), true
	}

	return "", false
}

// getDictionaryByRef resolves an indirect reference to a DictionaryToken.
func (p *BookmarksProvider) getDictionaryByRef(ref core.IndirectReference) *tokens.DictionaryToken {
	obj := p.pdfScanner.Get(ref)
	if obj == nil {
		return nil
	}

	token := obj.Data()
	result, ok := token.(*tokens.DictionaryToken)
	if !ok {
		return nil
	}
	return result
}

// unwrapIndirect unwraps an IndirectReferenceToken by reading through the scanner.
func unwrapIndirect(token tokens.Token, scanner tokenization.PdfTokenScanner) tokens.Token {
	if token == nil || scanner == nil {
		return token
	}

	_, isIndirect := token.(*tokens.IndirectReferenceToken)
	if !isIndirect {
		return token
	}

	if scanner.Advance() {
		return scanner.Current()
	}
	return token
}
