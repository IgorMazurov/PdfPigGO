package outline

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/outline/destinations"
)

// EmbeddedBookmarkNode represents a bookmark node that corresponds to a location
// in an embedded file within the PDF document.
type EmbeddedBookmarkNode struct {
	*DocumentBookmarkNode

	// FileSpecification is the file specification for the embedded file.
	FileSpecification string
}

// NewEmbeddedBookmarkNode creates a new EmbeddedBookmarkNode.
func NewEmbeddedBookmarkNode(title string, level int, dest destinations.ExplicitDestination, children []*BookmarkNode, fileSpecification string) (*EmbeddedBookmarkNode, error) {
	if fileSpecification == "" {
		return nil, fmt.Errorf("fileSpecification cannot be empty")
	}

	docNode, err := NewDocumentBookmarkNode(title, level, dest, children)
	if err != nil {
		return nil, err
	}

	return &EmbeddedBookmarkNode{
		DocumentBookmarkNode: docNode,
		FileSpecification:    fileSpecification,
	}, nil
}

// String returns a string representation of the embedded bookmark.
func (n *EmbeddedBookmarkNode) String() string {
	return fmt.Sprintf("embedded file '%s', %d, %s", n.FileSpecification, n.Level, n.Title)
}
