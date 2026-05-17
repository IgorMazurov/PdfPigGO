package content

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// pageCounter tracks the number of pages discovered during page tree traversal.
type pageCounter struct {
	pageCount int
}

func (c *pageCounter) increment() {
	c.pageCount++
}

// CreatePages builds a Pages instance from the root /Pages node in the document's page tree.
// It traverses the page tree, assigns page numbers, and returns a fully populated Pages object.
func CreatePages[T any](
	pagesReference core.IndirectReference,
	pagesDictionary *tokens.DictionaryToken,
	scanner tokenization.PdfTokenScanner,
	pageFactory PageFactory[T],
	log logging.Log,
	isLenientParsing bool,
) (*Pages, error) {
	if pagesDictionary == nil {
		return nil, fmt.Errorf("pages dictionary cannot be nil")
	}

	if scanner == nil {
		return nil, fmt.Errorf("scanner cannot be nil")
	}

	pageNumber := &pageCounter{}

	parentRef, err := core.NewIndirectReference(1, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create parent reference: %w", err)
	}

	pageTree, err := processPagesNode(
		pagesReference,
		pagesDictionary,
		parentRef,
		true,
		scanner,
		isLenientParsing,
		pageNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to process page tree: %w", err)
	}

	if !pageTree.IsRoot() {
		return nil, fmt.Errorf("page tree must be the root page tree node")
	}

	pagesByNumber := make(map[int]*PageTreeNode)
	populatePageByNumberDictionary(pageTree, pagesByNumber)

	dictionaryPageCount := getDictionaryPageCount(pagesDictionary)
	if dictionaryPageCount != len(pagesByNumber) {
		log.Warn(fmt.Sprintf("Dictionary Page Count %d different to discovered pages %d. Using %d.",
			dictionaryPageCount, len(pagesByNumber), len(pagesByNumber)))
	}

	pages, err := NewPages(pageFactory, scanner, pageTree, pagesByNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to create pages: %w", err)
	}

	return pages, nil
}

func getDictionaryPageCount(dict *tokens.DictionaryToken) int {
	if countRaw, found := dict.TryGet(tokens.Count); found {
		if numericTok, ok := countRaw.(*tokens.NumericToken); ok {
			return numericTok.IntVal()
		}
	}
	return 0
}

func processPagesNode(
	referenceInput core.IndirectReference,
	nodeDictionaryInput *tokens.DictionaryToken,
	parentReferenceInput core.IndirectReference,
	isRoot bool,
	pdfTokenScanner tokenization.PdfTokenScanner,
	isLenientParsing bool,
	pageNumber *pageCounter,
) (*PageTreeNode, error) {
	isPage := checkIfIsPage(nodeDictionaryInput, parentReferenceInput, isRoot, pdfTokenScanner, isLenientParsing)

	if isPage {
		pageNumber.increment()
		pn := pageNumber.pageCount
		node, err := NewPageTreeNode(nodeDictionaryInput, referenceInput, true, &pn)
		if err != nil {
			return nil, fmt.Errorf("failed to create page tree node: %w", err)
		}
		result, err := node.WithChildren([]*PageTreeNode{})
		if err != nil {
			return nil, fmt.Errorf("failed to set children on page node: %w", err)
		}
		return result, nil
	}

	const infiniteLoopWorkingWindow = 1000
	visitedTokens := make(map[int64]map[int]bool)
	visitedTokensQueue := make([]indirectRefID, 0, infiniteLoopWorkingWindow)

	type processItem struct {
		thisPage       *PageTreeNode
		reference      core.IndirectReference
		nodeDictionary *tokens.DictionaryToken
		parentReference core.IndirectReference
		nodeChildren   *[]*PageTreeNode
	}

	firstPage, err := NewPageTreeNode(nodeDictionaryInput, referenceInput, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create first page node: %w", err)
	}

var setChildren []func()
	firstPageChildren := make([]*PageTreeNode, 0)

	toProcess := []processItem{
		{
			thisPage:       firstPage,
			reference:      referenceInput,
			nodeDictionary: nodeDictionaryInput,
			parentReference: parentReferenceInput,
			nodeChildren:   &firstPageChildren,
		},
	}

	for len(toProcess) > 0 {
		current := &toProcess[0]
		toProcess = toProcess[1:]

		currentRefObjNum := current.reference.ObjectNumber()
		currentRefGen := current.reference.Generation()

		if gens, exists := visitedTokens[currentRefObjNum]; exists {
			if gens[currentRefGen] {
				continue
			} else {
				gens[currentRefGen] = true
			}
		} else {
			visitedTokens[currentRefObjNum] = map[int]bool{currentRefGen: true}
			visitedTokensQueue = append(visitedTokensQueue, indirectRefID{objNum: currentRefObjNum, gen: currentRefGen})

			if len(visitedTokensQueue) >= infiniteLoopWorkingWindow {
				toRemove := visitedTokensQueue[0]
				visitedTokensQueue = visitedTokensQueue[1:]
				gens := visitedTokens[toRemove.objNum]
				delete(gens, toRemove.gen)
				if len(gens) == 0 {
					delete(visitedTokens, toRemove.objNum)
				}
			}
		}

		kidsArray, found := current.nodeDictionary.TryGet(tokens.Kids)
		var kids *tokens.ArrayToken

		if !found {
			if !isLenientParsing {
				return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("Pages node in the document pages tree did not define a kids array: %s.", current.nodeDictionary))
			}
			kids = tokens.NewArrayToken(nil)
		} else {
			kids = tryGetToken[*tokens.ArrayToken](kidsArray, pdfTokenScanner)
			if kids == nil {
				if !isLenientParsing {
					return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("Kids entry is not an array: %s.", current.nodeDictionary))
				}
				kids = tokens.NewArrayToken(nil)
			}
		}

		for _, kid := range kids.Data() {
			kidRef, ok := kid.(*tokens.IndirectReferenceToken)
			if !ok {
				return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("Kids array contained invalid entry (must be indirect reference): %v.", kid))
			}

			kidDictionaryToken := tryGetDictionary(kidRef, pdfTokenScanner)

			if kidDictionaryToken == nil && !isLenientParsing {
				return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("Could not find dictionary associated with reference in pages kids array: %v.", kidRef))
			}

			if kidDictionaryToken == nil {
				kidDictionaryToken, _ = tokens.WithMap(make(map[string]tokens.Token))
			}

			isChildPage := checkIfIsPage(kidDictionaryToken, current.reference, false, pdfTokenScanner, isLenientParsing)

			if isChildPage {
				kidPageNode, err := NewPageTreeNode(kidDictionaryToken, kidRef.Data(), true, nil)
				if err != nil {
					return nil, fmt.Errorf("failed to create child page node: %w", err)
				}
				result, err := kidPageNode.WithChildren([]*PageTreeNode{})
				if err != nil {
					return nil, fmt.Errorf("failed to set children on child page node: %w", err)
				}
				*current.nodeChildren = append(*current.nodeChildren, result)
			} else {
				kidChildNode, err := NewPageTreeNode(kidDictionaryToken, kidRef.Data(), false, nil)
				if err != nil {
					return nil, fmt.Errorf("failed to create child node: %w", err)
				}
				kidChildren := make([]*PageTreeNode, 0)

				toProcess = append(toProcess, processItem{
					thisPage:       kidChildNode,
					reference:      kidRef.Data(),
					nodeDictionary: kidDictionaryToken,
					parentReference: current.reference,
					nodeChildren:   &kidChildren,
				})

				setChildren = append(setChildren, func() {
					for _, child := range kidChildren {
						child.SetParent(kidChildNode)
					}
					kidChildNode.children = kidChildren
				})

				*current.nodeChildren = append(*current.nodeChildren, kidChildNode)
			}
		}
	}

	for _, action := range setChildren {
		action()
	}

	// Set firstPage's children from the collected slice
	for _, child := range firstPageChildren {
		child.SetParent(firstPage)
	}
	firstPage.children = firstPageChildren

	childrenInOrder := toRecursiveOrderList(firstPage)
	for _, child := range childrenInOrder {
		if child.IsPage && child.pageNumber == nil {
			pageNumber.increment()
			pn := pageNumber.pageCount
			child.pageNumber = &pn
		}
	}

	return firstPage, nil
}

type indirectRefID struct {
	objNum int64
	gen    int
}

func checkIfIsPage(
	nodeDictionary *tokens.DictionaryToken,
	parentReference core.IndirectReference,
	isRoot bool,
	pdfTokenScanner tokenization.PdfTokenScanner,
	isLenientParsing bool,
) bool {
	isPage := false

	typeTokenRaw, foundType := nodeDictionary.TryGet(tokens.Type)

	if !foundType {
		if !isLenientParsing {
			panic(core.NewPdfDocumentFormatException(fmt.Sprintf("Node in the document pages tree did not define a type: %s.", nodeDictionary)))
		}

		kidsFound, _ := nodeDictionary.TryGet(tokens.Kids)
		if kidsFound == nil {
			isPage = true
		}
	} else {
		nameTok := tryGetToken[*tokens.NameToken](typeTokenRaw, pdfTokenScanner)
		if nameTok == nil {
			if !isLenientParsing {
				panic(core.NewPdfDocumentFormatException(fmt.Sprintf("Type value is not a name token: %s.", nodeDictionary)))
			}
			return false
		}

		isPage = nameTok.Data() == tokens.Page.Data()

		if !isPage && nameTok.Data() != tokens.Pages.Data() && !isLenientParsing {
			panic(core.NewPdfDocumentFormatException(fmt.Sprintf("Node in the document pages tree defined invalid type: %s.", nodeDictionary)))
		}
	}

	if !isLenientParsing && !isRoot {
		parentToken, foundParent := nodeDictionary.TryGet(tokens.Parent)

		if !foundParent {
			panic(core.NewPdfDocumentFormatException(fmt.Sprintf("Could not find parent indirect reference token on pages tree node: %s.", nodeDictionary)))
		}

		parentRefTok, ok := parentToken.(*tokens.IndirectReferenceToken)
		if !ok {
			panic(core.NewPdfDocumentFormatException(fmt.Sprintf("Parent is not an indirect reference token: %s.", nodeDictionary)))
		}

		if parentRefTok.Data() != parentReference {
			panic(core.NewPdfDocumentFormatException(fmt.Sprintf("Pages tree node parent reference %v did not match actual parent %v.", parentRefTok.Data(), parentReference)))
		}
	}

	return isPage
}

func populatePageByNumberDictionary(node *PageTreeNode, result map[int]*PageTreeNode) {
	if node.IsPage {
		if pn := node.PageNumber(); pn == nil {
			panic(fmt.Errorf("node was page but did not have page number: %s", node))
		} else {
			result[*pn] = node
		}
		return
	}

	for _, child := range node.Children() {
		populatePageByNumberDictionary(child, result)
	}
}

// toRecursiveOrderList performs a depth-first traversal of the page tree and returns
// all nodes in document order (children before siblings at each level).
func toRecursiveOrderList(node *PageTreeNode) []*PageTreeNode {
	result := make([]*PageTreeNode, 0)
	var traverse func(*PageTreeNode)
	traverse = func(n *PageTreeNode) {
		for _, child := range n.Children() {
			traverse(child)
		}
		result = append(result, n)
	}
	traverse(node)
	return result
}

// tryGetDictionary attempts to resolve an IndirectReferenceToken to a DictionaryToken
// using the scanner, following the DirectObjectFinder pattern. Returns nil on failure.
func tryGetDictionary(ref *tokens.IndirectReferenceToken, scanner tokenization.PdfTokenScanner) *tokens.DictionaryToken {
	return tryGetToken[*tokens.DictionaryToken](ref, scanner)
}

// tryGetToken attempts to resolve a token to the target type T, following indirect references via the scanner.
func tryGetToken[T tokens.Token](token tokens.Token, scanner tokenization.PdfTokenScanner) T {
	var zero T

	if result, ok := any(token).(T); ok {
		return result
	}

	ref, ok := token.(*tokens.IndirectReferenceToken)
	if !ok {
		return zero
	}

	return resolveIndirect[T](ref.Data(), scanner)
}

// resolveIndirect follows an IndirectReference through the scanner, handling nested references.
func resolveIndirect[T tokens.Token](reference core.IndirectReference, scanner tokenization.PdfTokenScanner) T {
	var zero T

	obj := scanner.Get(reference)
	if obj == nil {
		return zero
	}

	data := obj.Data()

	if _, isNull := any(data).(*tokens.NullToken); isNull {
		return zero
	}

	if result, ok := any(data).(T); ok {
		return result
	}

	if nestedRef, ok := any(data).(*tokens.IndirectReferenceToken); ok {
		return resolveIndirect[T](nestedRef.Data(), scanner)
	}

	return zero
}
