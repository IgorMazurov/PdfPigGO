package outline

// Bookmarks represents the bookmarks (outlines) in a PDF document.
type Bookmarks struct {
	roots []any // stores concrete types: *DocumentBookmarkNode, *UriBookmarkNode, etc.
}

// Roots returns the root bookmark nodes preserving their concrete types.
func (b *Bookmarks) Roots() []any {
	return b.roots
}

// NewBookmarks creates a new Bookmarks instance from any-slice roots.
func NewBookmarks(roots []any) *Bookmarks {
	if roots == nil {
		roots = make([]any, 0)
	}

	return &Bookmarks{
		roots: roots,
	}
}

// NewBookmarksTyped creates a new Bookmarks instance from typed bookmark nodes,
// converting them to any for storage while preserving concrete type information.
func NewBookmarksTyped(roots []BookmarkNodeTyped) *Bookmarks {
	if roots == nil {
		roots = make([]BookmarkNodeTyped, 0)
	}

	asAny := make([]any, len(roots))
	for i, r := range roots {
		asAny[i] = r
	}

	return &Bookmarks{
		roots: asAny,
	}
}

// GetNodes returns all bookmark nodes in the tree via depth-first traversal.
func (b *Bookmarks) GetNodes() []*BookmarkNode {
	var result []*BookmarkNode

	for _, root := range b.roots {
		result = append(result, getNodesFromAny(root)...)
	}

	return result
}

// getNodesFromAny recursively collects all nodes under the given any value.
func getNodesFromAny(node any) []*BookmarkNode {
	var result []*BookmarkNode

	switch n := node.(type) {
	case *DocumentBookmarkNode:
		result = append(result, n.BookmarkNode)
		for _, child := range n.Children() {
			result = append(result, collectChildren(child)...)
		}
	case *ExternalBookmarkNode:
		result = append(result, n.DocumentBookmarkNode.BookmarkNode)
		for _, child := range n.DocumentBookmarkNode.Children() {
			result = append(result, collectChildren(child)...)
		}
	case *EmbeddedBookmarkNode:
		result = append(result, n.DocumentBookmarkNode.BookmarkNode)
		for _, child := range n.DocumentBookmarkNode.Children() {
			result = append(result, collectChildren(child)...)
		}
	case *UriBookmarkNode:
		result = append(result, n.BookmarkNode)
		for _, child := range n.Children() {
			result = append(result, collectChildren(child)...)
		}
	case *ContainerBookmarkNode:
		result = append(result, n.BookmarkNode)
		for _, child := range n.Children() {
			result = append(result, collectChildren(child)...)
		}
	case *BookmarkNode:
		result = append(result, n)
		for _, child := range n.Children() {
			result = append(result, collectChildren(child)...)
		}
	}

	return result
}

// collectChildren collects a BookmarkNode and all its descendants recursively.
func collectChildren(node *BookmarkNode) []*BookmarkNode {
	var result []*BookmarkNode
	result = append(result, node)
	for _, child := range node.Children() {
		result = append(result, collectChildren(child)...)
	}
	return result
}
