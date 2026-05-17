package document_layout_analysis

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"sync"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/geometry"
)

// numbersPattern matches digits or Roman numerals (case-insensitive, matching C# RegexOptions.IgnoreCase).
var numbersPattern = regexp.MustCompile(`(\d+)|(\b([MDCLXVImdclxvi]+)\b)`)

// BlockSegmenter abstracts page segmentation to avoid an import cycle between
// document_layout_analysis and its page_segmenter subpackage.
type BlockSegmenter interface {
	GetBlocks(words []*content.Word) []*TextBlock
}

// DecorationTextBlockClassifier retrieves blocks labelled as decoration (e.g. headers, footers)
// for each page in the document using content and geometric similarity measures.
type DecorationTextBlockClassifier struct{}

var Instance = &DecorationTextBlockClassifier{}

func replaceNumbers(text string) string {
	return numbersPattern.ReplaceAllStringFunc(text, func(match string) string {
		return "@"
	})
}

// GetFromPages returns decoration blocks for each page given a list of pages, word extractor,
// and block segmenter. Uses the default MinimumEditDistanceNormalised distance function.
// similarityThreshold defaults to 0.25, n defaults to 5, maxDegreeOfParallelism defaults to unlimited (-1).
func (c *DecorationTextBlockClassifier) GetFromPages(
	pages []*content.Page,
	wordExtractor content.WordExtractor,
	segmenter BlockSegmenter,
	similarityThreshold float64,
	n int,
	maxDegreeOfParallelism int,
) [][]*TextBlock {
	return c.GetFromPagesWithDistance(pages, wordExtractor, segmenter, MinimumEditDistanceNormalised, similarityThreshold, n, maxDegreeOfParallelism)
}

// GetFromPagesWithDistance returns decoration blocks for each page given a list of pages,
// word extractor, block segmenter, and custom distance function.
func (c *DecorationTextBlockClassifier) GetFromPagesWithDistance(
	pages []*content.Page,
	wordExtractor content.WordExtractor,
	segmenter BlockSegmenter,
	minimumEditDistanceNormalised func(string, string) float64,
	similarityThreshold float64,
	n int,
	maxDegreeOfParallelism int,
) [][]*TextBlock {
	if len(pages) < 2 {
		return nil
	}

	type pageBlocks struct {
		index int
		blocks []*TextBlock
	}

	ch := make(chan pageBlocks, len(pages))
	var wg sync.WaitGroup
	semaphoreSize := maxDegreeOfParallelism
	if semaphoreSize <= 0 {
		semaphoreSize = math.MaxInt32
	}
	semaphore := make(chan struct{}, semaphoreSize)

	for i, page := range pages {
		wg.Add(1)
		go func(idx int, p *content.Page) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			words := p.GetWordsWithExtractor(wordExtractor)
			blocks := segmenter.GetBlocks(words)
			ch <- pageBlocks{index: idx, blocks: blocks}
		}(i, page)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	allPages := make(map[int][]*TextBlock, len(pages))
	for pb := range ch {
		allPages[pb.index] = pb.blocks
	}

	orderedPages := make([][]*TextBlock, len(pages))
	for i := range orderedPages {
		orderedPages[i] = allPages[i]
	}

	return c.GetFromBlocksWithDistance(orderedPages, minimumEditDistanceNormalised, similarityThreshold, n, maxDegreeOfParallelism)
}

// GetFromBlocks returns decoration blocks for each page given pre-extracted text blocks.
// Uses the default MinimumEditDistanceNormalised distance function.
func (c *DecorationTextBlockClassifier) GetFromBlocks(
	pagesTextBlocks [][]*TextBlock,
	similarityThreshold float64,
	n int,
	maxDegreeOfParallelism int,
) [][]*TextBlock {
	return c.GetFromBlocksWithDistance(pagesTextBlocks, MinimumEditDistanceNormalised, similarityThreshold, n, maxDegreeOfParallelism)
}

// GetFromBlocksWithDistance returns decoration blocks for each page given pre-extracted text
// blocks and a custom distance function. This is the core algorithm implementation.
func (c *DecorationTextBlockClassifier) GetFromBlocksWithDistance(
	pagesTextBlocks [][]*TextBlock,
	minimumEditDistanceNormalised func(string, string) float64,
	similarityThreshold float64,
	n int,
	maxDegreeOfParallelism int,
) [][]*TextBlock {
	if len(pagesTextBlocks) < 2 {
		return nil
	}

	pageDecorations := make([][]*TextBlock, len(pagesTextBlocks))
	for i := range pageDecorations {
		pageDecorations[i] = make([]*TextBlock, 0)
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphoreSize := maxDegreeOfParallelism
	if semaphoreSize <= 0 {
		semaphoreSize = math.MaxInt32
	}
	semaphore := make(chan struct{}, semaphoreSize)

	for p := 0; p < len(pagesTextBlocks); p++ {
		wg.Add(1)
		go func(pageIdx int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			pMinus1 := getPreviousPageNumber(pageIdx, len(pagesTextBlocks))
			pPlus1 := getNextPageNumber(pageIdx, len(pagesTextBlocks))

			previousPage := make([]*TextBlock, len(pagesTextBlocks[pMinus1]))
			copy(previousPage, pagesTextBlocks[pMinus1])

			currentPage := make([]*TextBlock, len(pagesTextBlocks[pageIdx]))
			copy(currentPage, pagesTextBlocks[pageIdx])

			nextPage := make([]*TextBlock, len(pagesTextBlocks[pPlus1]))
			copy(nextPage, pagesTextBlocks[pPlus1])

			nCurrent := n
			if len(currentPage) < nCurrent {
				nCurrent = len(currentPage)
			}

			decorations := classifyPageDecorations(previousPage, currentPage, nextPage, minimumEditDistanceNormalised, similarityThreshold, nCurrent)

			mu.Lock()
			pageDecorations[pageIdx] = decorations
			mu.Unlock()
		}(p)
	}

	wg.Wait()
	return pageDecorations
}

// classifyPageDecorations identifies decoration blocks on a single page by comparing against
// previous and next pages using four different sort orders.
func classifyPageDecorations(
	previousPage, currentPage, nextPage []*TextBlock,
	minimumEditDistanceNormalised func(string, string) float64,
	similarityThreshold float64,
	nCurrent int,
) []*TextBlock {
	set := NewOrderedSet[*TextBlock]()

	// Sort from top to bottom (based on minimum y coordinate, i.e., Bottom)
	sortByBottom := func(blocks []*TextBlock) {
		sort.SliceStable(blocks, func(i, j int) bool {
			bi := blocks[i].BoundingBox.Bottom()
			bj := blocks[j].BoundingBox.Bottom()
			if bi != bj {
				return bi > bj
			}
			return blocks[i].BoundingBox.Left() < blocks[j].BoundingBox.Left()
		})
	}

	// Sort from bottom to top (maximum y coordinate, i.e., Top)
	sortByTop := func(blocks []*TextBlock) {
		sort.SliceStable(blocks, func(i, j int) bool {
			ti := blocks[i].BoundingBox.Top()
			tj := blocks[j].BoundingBox.Top()
			if ti != tj {
				return ti < tj
			}
			return blocks[i].BoundingBox.Left() < blocks[j].BoundingBox.Left()
		})
	}

	// Sort from left to right (minimum x coordinate)
	sortByLeft := func(blocks []*TextBlock) {
		sort.SliceStable(blocks, func(i, j int) bool {
			li := blocks[i].BoundingBox.Left()
			lj := blocks[j].BoundingBox.Left()
			if li != lj {
				return li < lj
			}
			return blocks[i].BoundingBox.Top() < blocks[j].BoundingBox.Top()
		})
	}

	// Sort from right to left (maximum x coordinate)
	sortByRight := func(blocks []*TextBlock) {
		sort.SliceStable(blocks, func(i, j int) bool {
			ri := blocks[i].BoundingBox.Right()
			rj := blocks[j].BoundingBox.Right()
			if ri != rj {
				return ri > rj
			}
			return blocks[i].BoundingBox.Top() < blocks[j].BoundingBox.Top()
		})
	}

	sorters := []func([]*TextBlock){sortByBottom, sortByTop, sortByLeft, sortByRight}

	for _, sorter := range sorters {
		prevCopy := make([]*TextBlock, len(previousPage))
		copy(prevCopy, previousPage)

		currCopy := make([]*TextBlock, len(currentPage))
		copy(currCopy, currentPage)

		nextCopy := make([]*TextBlock, len(nextPage))
		copy(nextCopy, nextPage)

		sorter(prevCopy)
		sorter(currCopy)
		sorter(nextCopy)

		for i := 0; i < nCurrent; i++ {
			current := currCopy[i]
			score := scoreBlock(current, prevCopy, nextCopy, minimumEditDistanceNormalised, similarityThreshold, nCurrent)
			if score >= similarityThreshold {
				set.TryAdd(current)
			}
		}
	}

	return set.GetList()
}

// contentSimilarity calculates similarity from the normalized edit distance between two
// content strings where digits are replaced with "@" chars. A value of 1 means identical.
func contentSimilarity(b1, b2 *TextBlock, minimumEditDistanceNormalised func(string, string) float64) float64 {
	text1 := replaceNumbers(b1.Text)
	text2 := replaceNumbers(b2.Text)
	return 1.0 - minimumEditDistanceNormalised(text1, text2)
}

// geomSimilarity returns the area of intersection between two bounding boxes divided by
// the larger of the two bounding box areas. Returns 0 if no intersection exists.
func geomSimilarity(b1, b2 *TextBlock) float64 {
	intersect := geometry.RectangleIntersect(b1.BoundingBox, b2.BoundingBox)
	if intersect == nil {
		return 0
	}
	maxArea := math.Max(b1.BoundingBox.Area(), b2.BoundingBox.Area())
	if maxArea == 0 {
		return 0
	}
	return intersect.Area() / maxArea
}

// similarity returns the product of content and geometric similarity for two text blocks.
func similarity(b1, b2 *TextBlock, minimumEditDistanceNormalised func(string, string) float64) float64 {
	return contentSimilarity(b1, b2, minimumEditDistanceNormalised) * geomSimilarity(b1, b2)
}

// scoreI computes the average similarity between current block and matching blocks on
// previous and next pages.
func scoreI(current, previous, next *TextBlock, minimumEditDistanceNormalised func(string, string) float64) float64 {
	return 0.5 * (similarity(current, next, minimumEditDistanceNormalised)+similarity(current, previous, minimumEditDistanceNormalised))
}

// scoreBlock scores a text block against corresponding blocks from previous and next pages.
// Returns early once the threshold is reached for efficiency.
func scoreBlock(
	current *TextBlock,
	previous, next []*TextBlock,
	minimumEditDistanceNormalised func(string, string) float64,
	threshold float64,
	n int,
) float64 {
	effectiveN := n
	if len(previous) < effectiveN {
		effectiveN = len(previous)
	}
	if len(next) < effectiveN {
		effectiveN = len(next)
	}

	score := 0.0
	for i := 0; i < effectiveN; i++ {
		s := scoreI(current, previous[i], next[i], minimumEditDistanceNormalised)
		if s > score {
			score = s
		}
		if score >= threshold {
			return score
		}
	}
	return score
}

// getPreviousPageNumber returns the index of the previous page to compare against.
// For documents with more than three pages, it skips one additional page to account
// for two-sided layouts (comparing even with even, odd with odd).
func getPreviousPageNumber(currentPage, pagesCount int) int {
	pMinus1 := currentPage - 1
	if pMinus1 < 0 {
		pMinus1 = pagesCount - 1
	}
	if pagesCount > 3 {
		pMinus1--
		if pMinus1 < 0 {
			pMinus1 = pagesCount - 1
		}
	}
	return pMinus1
}

// getNextPageNumber returns the index of the next page to compare against.
// For documents with more than three pages, it skips one additional page to account
// for two-sided layouts (comparing even with even, odd with odd).
func getNextPageNumber(currentPage, pagesCount int) int {
	pPlus1 := currentPage + 1
	if pPlus1 >= pagesCount {
		pPlus1 = 0
	}
	if pagesCount > 3 {
		pPlus1++
		if pPlus1 >= pagesCount {
			pPlus1 = 0
		}
	}
	return pPlus1
}

// ErrLessThanTwoPages is returned when the decoration classifier cannot operate on fewer than 2 pages.
var ErrLessThanTwoPages = fmt.Errorf("the algorithm cannot be used with a document of less than 2 pages")
