package reading_order_detector

import (
	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
)

// SpatialReasoningRules defines the spatial reasoning constraints used to
// determine the reading order of text blocks on a page.
type SpatialReasoningRules int

const (
	// Basic applies basic spatial reasoning: in western culture the reading
	// order is from left to right and from top to bottom.
	Basic SpatialReasoningRules = iota

	// RowWise treats text-blocks as read in rows from left-to-right,
	// top-to-bottom. The diagonal direction 'left-bottom to top-right' cannot
	// be present among the Basic relations allowed.
	RowWise

	// ColumnWise treats text-blocks as read in columns, from top-to-bottom
	// and from left-to-right. The diagonal direction 'right-top to bottom-left'
	// cannot be present among the Basic relations allowed.
	ColumnWise
)

// UnsupervisedReadingOrderDetector retrieves blocks' reading order using
// spatial reasoning (Allen's interval relations) and optionally the rendering
// order (TextSequence). See section 4.1 of 'Unsupervised document structure
// analysis of digital scientific articles' by S. Klampfl et al., and
// 'Document Understanding for a Broad Class of Documents' by L. Todoran et al.
type UnsupervisedReadingOrderDetector struct {
	// UseRenderingOrder indicates whether to also use the rendering order
	// as indicated by the TextSequence.
	UseRenderingOrder bool

	// SpatialReasoningRule is the rule encoding the spatial reasoning constraints.
	SpatialReasoningRule SpatialReasoningRules

	// T is the tolerance parameter: if two coordinates are closer than T they
	// are considered equal. This flexibility is necessary because text blocks in
	// the same column might not be exactly aligned due to noise in PDF extraction.
	T float64

	getBeforeInMethod func(a, b *document_layout_analysis.TextBlock, t float64) bool
}

var _ ReadingOrderDetector = (*UnsupervisedReadingOrderDetector)(nil)

// UnsupervisedInstance is the default singleton instance of UnsupervisedReadingOrderDetector
// with T=5, ColumnWise spatial reasoning, and rendering order enabled.
var UnsupervisedInstance = NewUnsupervisedReadingOrderDetector(5, ColumnWise, true)

// NewUnsupervisedReadingOrderDetector creates a new detector with the given
// tolerance, spatial reasoning rule, and rendering order flag.
func NewUnsupervisedReadingOrderDetector(T float64, spatialReasoningRule SpatialReasoningRules, useRenderingOrder bool) *UnsupervisedReadingOrderDetector {
	d := &UnsupervisedReadingOrderDetector{
		T:                   T,
		SpatialReasoningRule: spatialReasoningRule,
		UseRenderingOrder:   useRenderingOrder,
	}

	switch d.SpatialReasoningRule {
	case ColumnWise:
		if d.UseRenderingOrder {
			d.getBeforeInMethod = func(a, b *document_layout_analysis.TextBlock, t float64) bool {
				return getBeforeInReadingVertical(a, b, t) || getBeforeInRendering(a, b)
			}
		} else {
			d.getBeforeInMethod = getBeforeInReadingVertical
		}

	case RowWise:
		if d.UseRenderingOrder {
			d.getBeforeInMethod = func(a, b *document_layout_analysis.TextBlock, t float64) bool {
				return getBeforeInReadingHorizontal(a, b, t) || getBeforeInRendering(a, b)
			}
		} else {
			d.getBeforeInMethod = getBeforeInReadingHorizontal
		}

	case Basic:
		fallthrough
	default:
		if d.UseRenderingOrder {
			d.getBeforeInMethod = func(a, b *document_layout_analysis.TextBlock, t float64) bool {
				return getBeforeInReading(a, b, t) || getBeforeInRendering(a, b)
			}
		} else {
			d.getBeforeInMethod = getBeforeInReading
		}
	}

	return d
}

// Get returns the text blocks in reading order and sets each block's ReadingOrder.
func (d *UnsupervisedReadingOrderDetector) Get(textBlocks []*document_layout_analysis.TextBlock) []*document_layout_analysis.TextBlock {
	result := make([]*document_layout_analysis.TextBlock, 0, len(textBlocks))

	graph := buildGraph(textBlocks, d.T, d.getBeforeInMethod)

	readingOrder := 0

	for len(graph) > 0 {
		maxCount := -1
		var currentKey int

		for k, v := range graph {
			if len(v) > maxCount {
				maxCount = len(v)
				currentKey = k
			}
		}

		delete(graph, currentKey)

		for g := range graph {
			graph[g] = removeIndex(graph[g], currentKey)
		}

		block := textBlocks[currentKey]
		block.SetReadingOrder(readingOrder) //nolint:errcheck
		readingOrder++
		result = append(result, block)
	}

	return result
}

func buildGraph(textBlocks []*document_layout_analysis.TextBlock, T float64, compare func(*document_layout_analysis.TextBlock, *document_layout_analysis.TextBlock, float64) bool) map[int][]int {
	graph := make(map[int][]int)

	for i := 0; i < len(textBlocks); i++ {
		graph[i] = []int{}
	}

	for i := 0; i < len(textBlocks); i++ {
		a := textBlocks[i]
		for j := 0; j < len(textBlocks); j++ {
			if i == j {
				continue
			}
			b := textBlocks[j]
			if compare(a, b, T) {
				graph[i] = append(graph[i], j)
			}
		}
	}

	return graph
}

func getBeforeInRendering(a, b *document_layout_analysis.TextBlock) bool {
	return avgTextSequence(a) < avgTextSequence(b)
}

// getBeforeInReading applies basic spatial reasoning: in western culture the
// reading order is from left to right and from top to bottom.
func getBeforeInReading(a, b *document_layout_analysis.TextBlock, T float64) bool {
	xRelation := GetRelationX(a.BoundingBox, b.BoundingBox, T)
	yRelation := GetRelationY(a.BoundingBox, b.BoundingBox, T)

	return xRelation == Precedes ||
		yRelation == Precedes ||
		xRelation == Meets ||
		yRelation == Meets ||
		xRelation == Overlaps ||
		yRelation == Overlaps
}

// getBeforeInReadingVertical applies column-wise reasoning: text-blocks are
// read in columns, from top-to-bottom and from left-to-right.
func getBeforeInReadingVertical(a, b *document_layout_analysis.TextBlock, T float64) bool {
	xRelation := GetRelationX(a.BoundingBox, b.BoundingBox, T)
	yRelation := GetRelationY(a.BoundingBox, b.BoundingBox, T)

	return xRelation == Precedes ||
		xRelation == Meets ||
		(xRelation == Overlaps && (yRelation == Precedes ||
			yRelation == Meets ||
			yRelation == Overlaps)) ||
		((yRelation == Precedes || yRelation == Meets || yRelation == Overlaps) &&
			(xRelation == Precedes ||
				xRelation == Meets ||
				xRelation == Overlaps ||
				xRelation == Starts ||
				xRelation == FinishesI ||
				xRelation == Equals ||
				xRelation == During ||
				xRelation == DuringI ||
				xRelation == Finishes ||
				xRelation == StartsI ||
				xRelation == OverlapsI))
}

// getBeforeInReadingHorizontal applies row-wise reasoning: text-blocks are
// read in rows from left-to-right, top-to-bottom.
func getBeforeInReadingHorizontal(a, b *document_layout_analysis.TextBlock, T float64) bool {
	xRelation := GetRelationX(a.BoundingBox, b.BoundingBox, T)
	yRelation := GetRelationY(a.BoundingBox, b.BoundingBox, T)

	return yRelation == Precedes ||
		yRelation == Meets ||
		(yRelation == Overlaps && (xRelation == Precedes ||
			xRelation == Meets ||
			xRelation == Overlaps)) ||
		((xRelation == Precedes || xRelation == Meets || xRelation == Overlaps) &&
			(yRelation == Precedes ||
				yRelation == Meets ||
				yRelation == Overlaps ||
				yRelation == Starts ||
				yRelation == FinishesI ||
				yRelation == Equals ||
				yRelation == During ||
				yRelation == DuringI ||
				yRelation == Finishes ||
				yRelation == StartsI ||
				yRelation == OverlapsI))
}

func removeIndex(slice []int, idx int) []int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == idx {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
