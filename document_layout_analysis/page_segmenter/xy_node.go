package page_segmenter

import (
	"log"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/reading_order_detector"
	"github.com/uglytoad/pdfpig/go/geometry"
)

// XYNode is an internal node in the RecursiveXYCut algorithm tree.
type XYNode struct {
	BoundingBox core.PdfRectangle
	Children    []*XYNode
}

// NewXYNode creates a new internal XYNode from children.
func NewXYNode(children ...*XYNode) *XYNode {
	return NewXYNodes(children)
}

// NewXYNodes creates a new internal XYNode from a slice of children.
func NewXYNodes(children []*XYNode) *XYNode {
	if len(children) == 0 {
		return &XYNode{
			Children: nil,
		}
	}

	minLeft := children[0].BoundingBox.Left()
	minBottom := children[0].BoundingBox.Bottom()
	maxRight := children[0].BoundingBox.Right()
	maxTop := children[0].BoundingBox.Top()

	for _, child := range children[1:] {
		left := child.BoundingBox.Left()
		if left < minLeft {
			minLeft = left
		}
		bottom := child.BoundingBox.Bottom()
		if bottom < minBottom {
			minBottom = bottom
		}
		right := child.BoundingBox.Right()
		if right > maxRight {
			maxRight = right
		}
		top := child.BoundingBox.Top()
		if top > maxTop {
			maxTop = top
		}
	}

	return &XYNode{
		BoundingBox: core.NewPdfRectangleFloat(minLeft, minBottom, maxRight, maxTop),
		Children:    children,
	}
}

// CountWords recursively counts the words included in this node.
func (n *XYNode) CountWords() int {
	if n.Children == nil {
		return 0
	}

	count := 0
	for _, child := range n.Children {
		switch c := any(child).(type) {
		case *XYLeaf:
			count += len(c.Words)
		case *XYNode:
			count += c.CountWords()
		}
	}
	return count
}

// GetLeaves recursively collects the leaf nodes of this node in reading order.
func (n *XYNode) GetLeaves() []*XYLeaf {
	if n.Children == nil || len(n.Children) == 0 {
		return nil
	}

	var leaves []*XYLeaf
	getLeavesRecursive(n.Children, &leaves, 0)
	return leaves
}

func getLeavesRecursive(children []*XYNode, leaves *[]*XYLeaf, level int) {
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		if leaf, ok := any(child).(*XYLeaf); ok {
			*leaves = append(*leaves, leaf)
		}
	}

	level++
	isVerticalCut := (level-1)%2 == 0

	var notLeaves []*XYNode
	for _, child := range children {
		if _, ok := any(child).(*XYLeaf); !ok {
			notLeaves = append(notLeaves, child)
		}
	}

	if isVerticalCut {
		sortByLeftAsc(notLeaves)
	} else {
		sortByTopDesc(notLeaves)
	}

	for _, node := range notLeaves {
		getLeavesRecursive(node.Children, leaves, level)
	}
}

func sortByLeftAsc(nodes []*XYNode) {
	for i := 1; i < len(nodes); i++ {
		for j := 0; j < len(nodes)-i; j++ {
			if nodes[j].BoundingBox.Left() > nodes[j+1].BoundingBox.Left() {
				nodes[j], nodes[j+1] = nodes[j+1], nodes[j]
			}
		}
	}
}

func sortByTopDesc(nodes []*XYNode) {
	for i := 1; i < len(nodes); i++ {
		for j := 0; j < len(nodes)-i; j++ {
			if nodes[j].BoundingBox.Top() < nodes[j+1].BoundingBox.Top() {
				nodes[j], nodes[j+1] = nodes[j+1], nodes[j]
			}
		}
	}
}

// String returns a debug string for the node.
func (n *XYNode) String() string {
	return "Node"
}

// XYLeaf is a leaf node in the RecursiveXYCut algorithm, representing a block of words.
type XYLeaf struct {
	BoundingBox core.PdfRectangle
	Words       []*content.Word
}

// NewXYLeaf creates a new XYLeaf from words.
func NewXYLeaf(words ...*content.Word) (*XYLeaf, error) {
	return NewXYLeafs(words)
}

// NewXYLeafs creates a new XYLeaf from a slice of words.
func NewXYLeafs(words []*content.Word) (*XYLeaf, error) {
	if len(words) == 0 {
		return nil, core.NewPdfDocumentFormatException("XYLeaf: words cannot be empty")
	}

	var minLeft, minBottom, maxRight, maxTop float64
	first := true

	for _, w := range words {
		bb := geometry.RectangleNormalise(w.BoundingBox())
		if first {
			minLeft = bb.Left()
			minBottom = bb.Bottom()
			maxRight = bb.Right()
			maxTop = bb.Top()
			first = false
		} else {
			left := bb.Left()
			if left < minLeft {
				minLeft = left
			}
			bottom := bb.Bottom()
			if bottom < minBottom {
				minBottom = bottom
			}
			right := bb.Right()
			if right > maxRight {
				maxRight = right
			}
			top := bb.Top()
			if top > maxTop {
				maxTop = top
			}
		}
	}

	return &XYLeaf{
		BoundingBox: core.NewPdfRectangleFloat(minLeft, minBottom, maxRight, maxTop),
		Words:       words,
	}, nil
}

// CountWords returns the number of words in this leaf.
func (l *XYLeaf) CountWords() int {
	return len(l.Words)
}

// GetLines groups the words by bottom coordinate and creates TextLine instances ordered by reading order.
func (l *XYLeaf) GetLines(wordSeparator string) []*document_layout_analysis.TextLine {
	groups := make(map[float64][]*content.Word)

	for _, w := range l.Words {
		bottom := w.BoundingBox().Bottom()
		groups[bottom] = append(groups[bottom], w)
	}

	lines := make([]*document_layout_analysis.TextLine, 0, len(groups))

	for _, groupWords := range groups {
		ordered, err := reading_order_detector.OrderWordsByReadingOrder(groupWords)
		if err != nil {
			log.Printf("warning: OrderWordsByReadingOrder failed in XYLeaf.GetLines: %v", err)
			continue
		}

		line, err := document_layout_analysis.NewTextLine(ordered, wordSeparator)
		if err != nil {
			log.Printf("warning: NewTextLine failed in XYLeaf.GetLines: %v", err)
			continue
		}

		lines = append(lines, line)
	}

	orderedLines, err := reading_order_detector.OrderLinesByReadingOrder(lines)
	if err != nil {
		log.Printf("warning: OrderLinesByReadingOrder failed in XYLeaf.GetLines: %v", err)
		return lines
	}

	return orderedLines
}

// String returns a debug string for the leaf.
func (l *XYLeaf) String() string {
	return "Leaf"
}
