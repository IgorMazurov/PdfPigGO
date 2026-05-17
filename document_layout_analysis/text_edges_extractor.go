package document_layout_analysis

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
)

// EdgeType defines the type of text edge.
type EdgeType int

const (
	// LeftEdge is a text edge where words have their BoundingBox's left coordinate aligned on the same vertical line.
	LeftEdge EdgeType = iota
	// MidEdge is a text edge where words have their BoundingBox's mid coordinate aligned on the same vertical line.
	MidEdge
	// RightEdge is a text edge where words have their BoundingBox's right coordinate aligned on the same vertical line.
	RightEdge
)

type edgeFunc struct {
	typ EdgeType
	fn  func(core.PdfRectangle) float64
}

var edgesFuncs = []edgeFunc{
	{LeftEdge, func(r core.PdfRectangle) float64 { return math.Round(r.Left()) }},
	{MidEdge, func(r core.PdfRectangle) float64 { return math.Round(r.Left() + r.Width/2) }},
	{RightEdge, func(r core.PdfRectangle) float64 { return math.Round(r.Right()) }},
}

// GetEdges returns the text edges. Text edges are where words have either their
// BoundingBox's left, right or mid coordinates aligned on the same vertical line.
// Useful to detect text columns, tables, justified text, lists, etc.
func GetEdges(pageWords []*content.Word, minimumElements int, maxDegreeOfParallelism int) (map[EdgeType][]core.PdfLine, error) {
	if minimumElements < 0 {
		return nil, fmt.Errorf("TextEdgesExtractor.GetEdges(): The minimum number of elements should be positive")
	}

	cleanWords := make([]*content.Word, 0, len(pageWords))
	for _, w := range pageWords {
		if strings.TrimSpace(w.Text) != "" {
			cleanWords = append(cleanWords, w)
		}
	}

	result := make(map[EdgeType][]core.PdfLine)

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxDegreeOfParallelism)

	for _, ef := range edgesFuncs {
		wg.Add(1)
		go func(e edgeFunc) {
			defer wg.Done()
			if maxDegreeOfParallelism > 0 {
				sem <- struct{}{}
				defer func() { <-sem }()
			}
			result[e.typ] = getVerticalEdges(cleanWords, e.fn, minimumElements)
		}(ef)
	}

	wg.Wait()

	return result, nil
}

func getVerticalEdges(pageWords []*content.Word, funcEdge func(core.PdfRectangle) float64, minimumElements int) []core.PdfLine {
	edges := make(map[float64][]*content.Word)
	for _, w := range pageWords {
		key := funcEdge(w.BoundingBox())
		edges[key] = append(edges[key], w)
	}

	filteredEdges := make(map[float64][]*content.Word)
	for k, v := range edges {
		if len(v) >= minimumElements {
			filteredEdges[k] = make([]*content.Word, len(v))
			copy(filteredEdges[k], v)
		}
	}

	cleanEdges := make(map[float64][][]*content.Word)

	for edgeKey, edgeWords := range filteredEdges {
		sortedEdges := make([]*content.Word, len(edgeWords))
		copy(sortedEdges, edgeWords)
		sortByBottom(sortedEdges)
		cleanEdges[edgeKey] = nil

		cuttings := findCuttings(pageWords, edgeWords, edgeKey)
		sortByBottom(cuttings)

		if len(cuttings) > 0 {
			for _, cut := range cuttings {
				group1 := make([]*content.Word, 0)
				for _, w := range sortedEdges {
					if w.BoundingBox().Top() < cut.BoundingBox().Bottom() {
						group1 = append(group1, w)
					}
				}
				if len(group1) >= minimumElements {
					cleanEdges[edgeKey] = append(cleanEdges[edgeKey], group1)
				}
				sortedEdges = removeFromSlice(sortedEdges, group1)
			}
			if len(sortedEdges) >= minimumElements {
				cleanEdges[edgeKey] = append(cleanEdges[edgeKey], sortedEdges)
			}
		} else {
			cleanEdges[edgeKey] = append(cleanEdges[edgeKey], sortedEdges)
		}
	}

	var lines []core.PdfLine
	for edgeX, groups := range cleanEdges {
		for _, group := range groups {
			minBottom := group[0].BoundingBox().Bottom()
			maxTop := group[0].BoundingBox().Top()
			for _, w := range group {
				bb := w.BoundingBox()
				if bb.Bottom() < minBottom {
					minBottom = bb.Bottom()
				}
				if bb.Top() > maxTop {
					maxTop = bb.Top()
				}
			}
			lines = append(lines, core.NewPdfLineFromCoords(edgeX, minBottom, edgeX, maxTop))
		}
	}

	return lines
}

func findCuttings(allWords, edgeWords []*content.Word, edgeKey float64) []*content.Word {
	edgeSet := make(map[*content.Word]bool)
	for _, w := range edgeWords {
		edgeSet[w] = true
	}

	var minBottom, maxTop float64
	first := true
	for _, w := range edgeWords {
		bb := w.BoundingBox()
		if first {
			minBottom = bb.Bottom()
			maxTop = bb.Top()
			first = false
		} else {
			if bb.Bottom() < minBottom {
				minBottom = bb.Bottom()
			}
			if bb.Top() > maxTop {
				maxTop = bb.Top()
			}
		}
	}

	var result []*content.Word
	for _, w := range allWords {
		if edgeSet[w] {
			continue
		}
		bb := w.BoundingBox()
		if bb.Left() < edgeKey && bb.Right() > edgeKey {
			if bb.Bottom() > minBottom && bb.Top() < maxTop {
				result = append(result, w)
			}
		}
	}

	return result
}

func sortByBottom(words []*content.Word) {
	for i := 0; i < len(words)-1; i++ {
		for j := i + 1; j < len(words); j++ {
			if words[j].BoundingBox().Bottom() < words[i].BoundingBox().Bottom() {
				words[i], words[j] = words[j], words[i]
			}
		}
	}
}

func removeFromSlice(slice, remove []*content.Word) []*content.Word {
	removeSet := make(map[*content.Word]bool)
	for _, w := range remove {
		removeSet[w] = true
	}
	result := make([]*content.Word, 0)
	for _, w := range slice {
		if !removeSet[w] {
			result = append(result, w)
		}
	}
	return result
}
