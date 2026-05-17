package outline

import "fmt"

// BookmarkNode represents a node in the PDF document's bookmarks (also known as outlines).
type BookmarkNode struct {
	// Title is the text displayed for this node.
	Title string

	// Level is the node's level in the hierarchy.
	Level int

	children []*BookmarkNode
}

// Children returns the bookmark's sub-bookmarks.
func (n *BookmarkNode) Children() []*BookmarkNode {
	return n.children
}

// IsLeaf reports whether this node is a leaf node (has no children).
func (n *BookmarkNode) IsLeaf() bool {
	return len(n.children) == 0
}

// NewBookmarkNode creates a new BookmarkNode.
func NewBookmarkNode(title string, level int, children []*BookmarkNode) (*BookmarkNode, error) {
	if children == nil {
		return nil, fmt.Errorf("children cannot be nil")
	}

	return &BookmarkNode{
		Title:    title,
		Level:    level,
		children: children,
	}, nil
}
