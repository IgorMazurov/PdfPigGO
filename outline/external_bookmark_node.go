package outline

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/outline/destinations"
)

// ExternalBookmarkNode represents a bookmark node that links to a destination
// in an external PDF file.
type ExternalBookmarkNode struct {
	*DocumentBookmarkNode

	// Filename is the name of the external PDF file containing this bookmark.
	Filename string
}

// NewExternalBookmarkNode creates a new ExternalBookmarkNode.
func NewExternalBookmarkNode(title string, level int, dest destinations.ExplicitDestination, children []*BookmarkNode, filename string) (*ExternalBookmarkNode, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename cannot be empty")
	}

	docNode, err := NewDocumentBookmarkNode(title, level, dest, children)
	if err != nil {
		return nil, err
	}

	return &ExternalBookmarkNode{
		DocumentBookmarkNode: docNode,
		Filename:             filename,
	}, nil
}

// String returns a string representation of the external bookmark.
func (n *ExternalBookmarkNode) String() string {
	return fmt.Sprintf("file '%s', %d, %s", n.Filename, n.Level, n.Title)
}
