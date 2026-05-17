package content

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Pages manages the page tree for a PDF document, providing access to all pages,
// page numbers, and their dictionaries. It caches page factories for creating
// typed page instances.
type Pages struct {
	pageFactoryCache   map[string]any
	defaultPageFactory any
	pdfScanner         tokenization.PdfTokenScanner
	pagesByNumber      map[int]*PageTreeNode
	pageTree           *PageTreeNode
}

// Count returns the total number of pages in the document.
func (p *Pages) Count() int {
	return len(p.pagesByNumber)
}

// PageTree returns the page tree root containing all pages, page numbers and their dictionaries.
func (p *Pages) PageTree() *PageTreeNode {
	return p.pageTree
}

// NewPages creates a new Pages instance with the given default page factory, scanner,
// page tree root, and number-to-node mapping.
func NewPages[T any](
	defaultFactory PageFactory[T],
	pdfScanner tokenization.PdfTokenScanner,
	pageTree *PageTreeNode,
	pagesByNumber map[int]*PageTreeNode,
) (*Pages, error) {
	if defaultFactory == nil {
		return nil, fmt.Errorf("default page factory cannot be nil")
	}

	if pdfScanner == nil {
		return nil, fmt.Errorf("pdf scanner cannot be nil")
	}

	pages := &Pages{
		pageFactoryCache:   make(map[string]any),
		defaultPageFactory: defaultFactory,
		pdfScanner:         pdfScanner,
		pagesByNumber:      pagesByNumber,
		pageTree:           pageTree,
	}

	pages.AddPageFactory(defaultFactory)

	return pages, nil
}

// GetPage returns the page at the given 1-based number using the default factory.
// The caller should type-assert the result to the expected page type.
func (p *Pages) GetPage(pageNumber int, namedDestinations any, parsingOptions *ParsingOptions) (any, error) {
	if p.defaultPageFactory == nil {
		return nil, fmt.Errorf("no default page factory configured")
	}

	// Try generic interface first
	type genericFactory interface {
		Create(number int, dictionary *tokens.DictionaryToken, pageTreeMembers *PageTreeMembers, namedDestinations any) (any, error)
	}

	if factory, ok := p.defaultPageFactory.(genericFactory); ok {
		return p.getPageWithFactory(factory.Create, pageNumber, namedDestinations, parsingOptions)
	}

	// Try parser.PageFactory-compatible signature
	type concreteFactory interface {
		Create(number int, dictionary *tokens.DictionaryToken, pageTreeMembers *PageTreeMembers, namedDestinations *destinations.NamedDestinations) (*Page, error)
	}

	if factory, ok := p.defaultPageFactory.(concreteFactory); ok {
		node := p.pagesByNumber[pageNumber]
		if node == nil {
			return nil, fmt.Errorf("page %d not found", pageNumber)
		}

		var result *Page
		var createErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					if pdfErr, ok := r.(*core.PdfDocumentFormatException); ok {
						createErr = pdfErr
					} else {
						panic(r)
					}
				}
			}()

			pageStack := collectAncestors(node)
			members := buildPageTreeMembers(pageStack, p.pdfScanner)
			nd := namedDestinations.(*destinations.NamedDestinations)
			res, err := factory.Create(pageNumber, node.NodeDictionary, members, nd)
			if err != nil {
				createErr = err
				return
			}
			result = res
		}()

		if createErr != nil {
			return nil, createErr
		}
		return result, nil
	}

	return nil, fmt.Errorf("default page factory is not compatible")
}

func (p *Pages) getPageWithFactory(
	createFunc func(number int, dictionary *tokens.DictionaryToken, pageTreeMembers *PageTreeMembers, namedDestinations any) (any, error),
	pageNumber int,
	namedDestinations any,
	parsingOptions *ParsingOptions,
) (any, error) {
	if pageNumber <= 0 || pageNumber > p.Count() {
		if parsingOptions != nil && parsingOptions.Logger != nil {
			parsingOptions.Logger.Error(fmt.Sprintf("Page %d requested but is out of range.", pageNumber))
		}
		return nil, fmt.Errorf("page number %d invalid, must be between 1 and %d", pageNumber, p.Count())
	}

	pageNode := p.GetPageNode(pageNumber)
	if pageNode == nil {
		return nil, fmt.Errorf("could not find page node for page number %d", pageNumber)
	}

	var result any
	var createErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				if pdfErr, ok := r.(*core.PdfDocumentFormatException); ok {
					createErr = pdfErr
				} else {
					panic(r)
				}
			}
		}()

		pageStack := collectAncestors(pageNode)
		pageTreeMembers := buildPageTreeMembers(pageStack, p.pdfScanner)

		res, err := createFunc(
			pageNumber,
			pageNode.NodeDictionary,
			pageTreeMembers,
			namedDestinations,
		)
		if err != nil {
			createErr = err
			return
		}
		result = res
	}()

	if createErr != nil {
		return nil, createErr
	}

	return result, nil
}

// collectAncestors walks from node up to root and returns them in order from root to leaf.
func collectAncestors(node *PageTreeNode) []*PageTreeNode {
	stack := make([]*PageTreeNode, 0)
	current := node
	for current != nil {
		stack = append(stack, current)
		current = current.Parent()
	}

	result := make([]*PageTreeNode, len(stack))
	for i := 0; i < len(stack); i++ {
		result[i] = stack[len(stack)-1-i]
	}

	return result
}

// buildPageTreeMembers extracts inherited properties from ancestor nodes in the page tree.
func buildPageTreeMembers(nodes []*PageTreeNode, scanner tokenization.PdfTokenScanner) *PageTreeMembers {
	members := &PageTreeMembers{
		ParentResources: make([]tokens.Token, 0),
	}

	for _, node := range nodes {
		if resourcesRaw, found := node.NodeDictionary.TryGet(tokens.Resources); found {
			resDict := resolveToDictionary(resourcesRaw, scanner)
			if resDict != nil {
				members.ParentResources = append(members.ParentResources, resDict)
			} else {
				members.ParentResources = append(members.ParentResources, resourcesRaw)
			}
		}

		if mediaBoxRaw, found := node.NodeDictionary.TryGet(tokens.MediaBox); found {
			mediaArray := resolveToArray(mediaBoxRaw, scanner)
			if mediaArray != nil {
				data := mediaArray.Data()
				if len(data) == 4 {
					rect := arrayToRectangle(data, scanner)
					if rect != nil {
						members.MediaBoxValue = NewMediaBox(*rect)
					}
				}
			}
		}

		if rotateRaw, found := node.NodeDictionary.TryGet(tokens.Rotate); found {
			if numericTok, ok := rotateRaw.(*tokens.NumericToken); ok {
				rot := numericTok.IntVal()
				members.Rotation = &rot
			}
		}
	}

	return members
}

// TypedPageFactory is a page factory that identifies its output type by name.
// Custom page factories should implement this interface to be discoverable via GetPage[T]().
type TypedPageFactory interface {
	// OutputTypeName returns the Go type name of the page type produced by Create(),
	// e.g., "*pdfpig.SimplePage" or "content.PageInformation".
	OutputTypeName() string
	// Create builds a page from the given PDF dictionary and context.
	Create(number int, dictionary *tokens.DictionaryToken, pageTreeMembers *PageTreeMembers, namedDestinations any) (any, error)
}

// AddPageFactory registers an additional page factory.
func (p *Pages) AddPageFactory(factory any) {
	if factory == nil {
		return
	}

	name := fmt.Sprintf("%T", factory)
	if _, exists := p.pageFactoryCache[name]; exists {
		panic(fmt.Errorf("could not add page factory for type '%s' as it was already added", name))
	}

	p.pageFactoryCache[name] = factory

	if tf, ok := factory.(TypedPageFactory); ok {
		p.pageFactoryCache[tf.OutputTypeName()] = factory
	}
}

// GetPage returns the page at the given 1-based number using a factory that produces type T.
// The factory must have been registered via AddPageFactory and implement TypedPageFactory.
// Note: Go does not support generic methods on structs, so this is a standalone function.
func GetPage[T any](p *Pages, pageNumber int, namedDestinations any, parsingOptions *ParsingOptions) (*T, error) {
	var zero T
	typeName := fmt.Sprintf("%T", &zero)

	factoryAny, found := p.pageFactoryCache[typeName]
	if !found {
		return nil, fmt.Errorf("could not find page factory of type %q for output type %s", typeName, typeName)
	}

	factory, ok := factoryAny.(TypedPageFactory)
	if !ok {
		return nil, fmt.Errorf("page factory registered under %q does not implement TypedPageFactory", typeName)
	}

	if pageNumber <= 0 || pageNumber > p.Count() {
		if parsingOptions != nil && parsingOptions.Logger != nil {
			parsingOptions.Logger.Error(fmt.Sprintf("Page %d requested but is out of range.", pageNumber))
		}
		return nil, fmt.Errorf("page number %d invalid, must be between 1 and %d", pageNumber, p.Count())
	}

	pageNode := p.GetPageNode(pageNumber)
	if pageNode == nil {
		return nil, fmt.Errorf("could not find page node for page number %d", pageNumber)
	}

	pageStack := collectAncestors(pageNode)
	pageTreeMembers := buildPageTreeMembers(pageStack, p.pdfScanner)

	resultAny, err := factory.Create(pageNumber, pageNode.NodeDictionary, pageTreeMembers, namedDestinations)
	if err != nil {
		return nil, fmt.Errorf("failed to create page %d: %w", pageNumber, err)
	}

	result, ok := resultAny.(*T)
	if !ok {
		return nil, fmt.Errorf("factory created type %T but expected *%T", resultAny, zero)
	}

	return result, nil
}

// GetPageNode returns the PageTreeNode for the given 1-based page number.
func (p *Pages) GetPageNode(pageNumber int) *PageTreeNode {
	node, ok := p.pagesByNumber[pageNumber]
	if !ok {
		panic(fmt.Errorf("could not find page node by number for: %d", pageNumber))
	}

	return node
}

// GetPageByReference finds and returns the page tree node whose reference matches the given one.
// Returns nil if no matching page is found.
func (p *Pages) GetPageByReference(reference core.IndirectReference) destinations.PageRef {
	for _, node := range p.pagesByNumber {
		if node.Reference == reference {
			return node
		}
	}

	return nil
}
