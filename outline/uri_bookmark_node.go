package outline

import "fmt"

// UriBookmarkNode represents a bookmark node that links to a URI on the Internet.
type UriBookmarkNode struct {
	*BookmarkNode

	// Uri is the uniform resource identifier to resolve.
	Uri string
}

// NewUriBookmarkNode creates a new UriBookmarkNode.
func NewUriBookmarkNode(title string, level int, uri string, children []*BookmarkNode) (*UriBookmarkNode, error) {
	if uri == "" {
		return nil, fmt.Errorf("uri cannot be empty")
	}

	node, err := NewBookmarkNode(title, level, children)
	if err != nil {
		return nil, err
	}

	return &UriBookmarkNode{
		BookmarkNode: node,
		Uri:          uri,
	}, nil
}

// String returns a string representation of the URI bookmark.
func (n *UriBookmarkNode) String() string {
	return fmt.Sprintf("URI '%s', %d, %s", n.Uri, n.Level, n.Title)
}
