package outline

// ContainerBookmarkNode represents a pure container bookmark node: it has a title
// and child nodes but no destination or action. This is used to handle the common
// "grouping" bookmarks in PDFs.
type ContainerBookmarkNode struct {
	*BookmarkNode
}

// NewContainerBookmarkNode creates a new ContainerBookmarkNode.
func NewContainerBookmarkNode(title string, level int, children []*BookmarkNode) (*ContainerBookmarkNode, error) {
	node, err := NewBookmarkNode(title, level, children)
	if err != nil {
		return nil, err
	}

	return &ContainerBookmarkNode{
		BookmarkNode: node,
	}, nil
}
