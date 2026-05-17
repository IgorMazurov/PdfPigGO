package content

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// PageTreeNode represents a node in the PDF document's page tree.
// Nodes may either be of type 'Page' — a single page, or 'Pages' — a container for multiple child Page
// or Pages nodes.
type PageTreeNode struct {
	// NodeDictionary is the dictionary for this node in the page tree.
	NodeDictionary *tokens.DictionaryToken

	// Reference is the indirect reference for this node in the page tree.
	Reference core.IndirectReference

	// IsPage indicates whether this node is a page or not. If false it must be a /Pages container.
	IsPage bool

	// pageNumber is the number of this page if IsPage is true.
	pageNumber *int

	// children are the child nodes if IsPage is false.
	children []*PageTreeNode

	// parent is the parent node, unless this is the root node.
	parent *PageTreeNode
}

// IsRoot reports whether this node is the root node.
func (n *PageTreeNode) IsRoot() bool {
	return n.parent == nil
}

// Children returns the child nodes of this node if IsPage is false.
func (n *PageTreeNode) Children() []*PageTreeNode {
	return n.children
}

// Parent returns the parent node of this node, unless it is the root node.
func (n *PageTreeNode) Parent() *PageTreeNode {
	return n.parent
}

// SetParent sets the parent node. This method is intended for internal use only.
func (n *PageTreeNode) SetParent(parent *PageTreeNode) {
	n.parent = parent
}

// NewPageTreeNode creates a new PageTreeNode.
func NewPageTreeNode(nodeDictionary *tokens.DictionaryToken, reference core.IndirectReference, isPage bool, pageNumber *int) (*PageTreeNode, error) {
	if nodeDictionary == nil {
		return nil, fmt.Errorf("nodeDictionary cannot be nil")
	}

	if !isPage && pageNumber != nil {
		return nil, fmt.Errorf("cannot define page number for a pages node")
	}

	return &PageTreeNode{
		NodeDictionary: nodeDictionary,
		Reference:      reference,
		IsPage:         isPage,
		pageNumber:     pageNumber,
	}, nil
}

// WithChildren sets the child nodes of this Pages container node.
func (n *PageTreeNode) WithChildren(children []*PageTreeNode) (*PageTreeNode, error) {
	if children == nil {
		return nil, fmt.Errorf("children cannot be nil")
	}

	if n.IsPage && len(children) > 0 {
		return nil, fmt.Errorf("cannot define children on a page node")
	}

	n.children = children

	for _, child := range children {
		child.parent = n
	}

	return n, nil
}

// String returns the string representation of the page tree node.
func (n *PageTreeNode) String() string {
	if n.IsPage {
		if pn := n.PageNumber(); pn != nil {
			return fmt.Sprintf("Page #%d: %s.", *pn, n.NodeDictionary)
		}
		return fmt.Sprintf("Page: %s.", n.NodeDictionary)
	}

	childCount := 0
	if n.children != nil {
		childCount = len(n.children)
	}

	return fmt.Sprintf("Pages (%d children): %s", childCount, n.NodeDictionary)
}

// PageNumber returns the page number, satisfying destinations.PageRef.
func (n *PageTreeNode) PageNumber() *int {
	return n.pageNumber
}

var _ fmt.Stringer = (*PageTreeNode)(nil)
