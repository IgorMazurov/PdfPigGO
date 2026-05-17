package outline

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/outline/destinations"
)

// DocumentBookmarkNode represents a bookmark node that links to a destination
// within the same PDF document.
type DocumentBookmarkNode struct {
	*BookmarkNode

	// Destination is the explicit destination this bookmark points to.
	Destination destinations.ExplicitDestination
}

// PageNumber returns the 1-based page number where the bookmark is located.
func (n *DocumentBookmarkNode) PageNumber() int {
	return n.Destination.PageNumber
}

// NewDocumentBookmarkNode creates a new DocumentBookmarkNode.
func NewDocumentBookmarkNode(title string, level int, dest destinations.ExplicitDestination, children []*BookmarkNode) (*DocumentBookmarkNode, error) {
	node, err := NewBookmarkNode(title, level, children)
	if err != nil {
		return nil, err
	}

	return &DocumentBookmarkNode{
		BookmarkNode: node,
		Destination:  dest,
	}, nil
}

// String returns a string representation of the document bookmark.
func (n *DocumentBookmarkNode) String() string {
	return fmt.Sprintf("page #%d, %d, %s", n.PageNumber(), n.Level, n.Title)
}
