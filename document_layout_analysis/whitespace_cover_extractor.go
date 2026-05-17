package document_layout_analysis

import (
	"sort"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
)

// WhitespaceCoverExtractor provides a top-down algorithm that finds a cover of the
// background whitespace of a document in terms of maximal empty rectangles.
// See Section 3.2 of "High precision text extraction from PDF documents" by
// Oyvind Raddum Berg and Section 2 of "Two geometric algorithms for layout analysis"
// by Thomas M. Breuel.

const defaultMaxRectangleCount = 40
const defaultWhitespaceFuzziness = 0.15

// GetWhitespacesFromWords finds whitespace rectangles from words, auto-computing
// minimum dimensions from the letter bounding box mode.
func GetWhitespacesFromWords(words []*content.Word, images []content.PdfImage, maxRectangleCount, maxBoundQueueSize int) []core.PdfRectangle {
	if maxRectangleCount == 0 {
		maxRectangleCount = defaultMaxRectangleCount
	}

	var widths, heights []float64
	for _, w := range words {
		for _, l := range w.Letters {
			widths = append(widths, l.BoundingBox.Width)
			heights = append(heights, l.BoundingBox.Height)
		}
	}

	minWidth := ModeFloat64(widths) * 1.25
	minHeight := ModeFloat64(heights) * 1.25

	return GetWhitespacesFromWordsAndDims(words, images, minWidth, minHeight, maxRectangleCount, defaultWhitespaceFuzziness, maxBoundQueueSize)
}

// GetWhitespacesFromWordsAndDims finds whitespace rectangles from words and images
// with explicit minimum rectangle dimensions.
func GetWhitespacesFromWordsAndDims(words []*content.Word, images []content.PdfImage, minWidth, minHeight float64, maxRectangleCount int, whitespaceFuzziness float64, maxBoundQueueSize int) []core.PdfRectangle {
	if maxRectangleCount == 0 {
		maxRectangleCount = defaultMaxRectangleCount
	}
	if whitespaceFuzziness == 0 {
		whitespaceFuzziness = defaultWhitespaceFuzziness
	}

	var bboxes []core.PdfRectangle
	for _, w := range words {
		bb := w.BoundingBox()
		if bb.Width > 0 && bb.Height > 0 {
			bboxes = append(bboxes, bb)
		}
	}

	if len(images) > 0 {
		for _, img := range images {
			bb := img.BoundingBox()
			if bb.Width > 0 && bb.Height > 0 {
				bboxes = append(bboxes, bb)
			}
		}
	}

	return GetWhitespacesFromRects(bboxes, minWidth, minHeight, maxRectangleCount, defaultWhitespaceFuzziness, maxBoundQueueSize)
}

// GetWhitespacesFromRects finds whitespace rectangles from a list of obstacle bounding boxes.
func GetWhitespacesFromRects(boundingboxes []core.PdfRectangle, minWidth, minHeight float64, maxRectangleCount int, whitespaceFuzziness float64, maxBoundQueueSize int) []core.PdfRectangle {
	if maxRectangleCount == 0 {
		maxRectangleCount = defaultMaxRectangleCount
	}
	if whitespaceFuzziness == 0 {
		whitespaceFuzziness = defaultWhitespaceFuzziness
	}

	if len(boundingboxes) == 0 {
		return nil
	}

	obstacles := deduplicateRects(boundingboxes)
	pageBound := getBound(obstacles)

	return getMaximalRectangles(pageBound, obstacles, minWidth, minHeight, maxRectangleCount, whitespaceFuzziness, maxBoundQueueSize)
}

func deduplicateRects(rects []core.PdfRectangle) []core.PdfRectangle {
	set := make(map[rectKey]struct{})
	var result []core.PdfRectangle
	for _, r := range rects {
		k := rectToKey(r)
		if _, exists := set[k]; !exists {
			set[k] = struct{}{}
			result = append(result, r)
		}
	}
	return result
}

type rectKey struct {
	left, bottom, right, top float64
}

func rectToKey(r core.PdfRectangle) rectKey {
	return rectKey{r.Left(), r.Bottom(), r.Right(), r.Top()}
}

// rectIntersection computes the intersection of two axis-aligned rectangles.
// Returns (intersection, true) if they overlap, or (zero rectangle, false) otherwise.
func rectIntersection(r1, r2 core.PdfRectangle) (core.PdfRectangle, bool) {
	left := max(r1.Left(), r2.Left())
	right := min(r1.Right(), r2.Right())
	bottom := max(r1.Bottom(), r2.Bottom())
	top := min(r1.Top(), r2.Top())

	if left >= right || bottom >= top {
		return core.PdfRectangle{}, false
	}

	return core.NewPdfRectangleFloat(left, bottom, right, top), true
}

// getMaximalRectangles runs the main maximal empty rectangle algorithm using a priority queue.
func getMaximalRectangles(bound core.PdfRectangle, obstacles []core.PdfRectangle, minWidth, minHeight float64, maxRectangleCount int, whitespaceFuzziness float64, maxBoundQueueSize int) []core.PdfRectangle {
	queue := newSortedQueue(maxBoundQueueSize)
	queue.enqueue(newQueueEntry(bound, obstacles, whitespaceFuzziness))

	var selected []core.PdfRectangle
	var holdList []*queueEntry

	for queue.any() {
		current := queue.dequeue()

		if current.isEmptyEnough(obstacles) {
			if isInsideAny(selected, current.bound) {
				continue
			}

			if !isAdjacentToPageBounds(bound, current.bound) &&
				!anyAdjacent(selected, current.bound) {
				holdList = append(holdList, current)
				continue
			}

			selected = append(selected, current.bound)

			if len(selected) >= maxRectangleCount {
				return selected
			}

			obstacles = append(obstacles, current.bound)

			for _, hold := range holdList {
				queue.enqueue(hold)
			}
			holdList = nil

			for _, entry := range queue.all() {
				if overlapsHard(current.bound, entry.bound) {
					entry.addWhitespace(current.bound)
				}
			}

			continue
		}

		pivot := current.getPivot(obstacles)
		b := current.bound

		rRight := core.NewPdfRectangleFloat(pivot.Right(), b.Bottom(), b.Right(), b.Top())
		if b.Right() > pivot.Right() && rRight.Height > minHeight && rRight.Width > minWidth {
			rightObstacles := filterOverlapping(current.obstacles, rRight)
			queue.enqueue(newQueueEntry(rRight, rightObstacles, whitespaceFuzziness))
		}

		rLeft := core.NewPdfRectangleFloat(b.Left(), b.Bottom(), pivot.Left(), b.Top())
		if b.Left() < pivot.Left() && rLeft.Height > minHeight && rLeft.Width > minWidth {
			leftObstacles := filterOverlapping(current.obstacles, rLeft)
			queue.enqueue(newQueueEntry(rLeft, leftObstacles, whitespaceFuzziness))
		}

		rAbove := core.NewPdfRectangleFloat(b.Left(), b.Bottom(), b.Right(), pivot.Bottom())
		if b.Bottom() < pivot.Bottom() && rAbove.Height > minHeight && rAbove.Width > minWidth {
			aboveObstacles := filterOverlapping(current.obstacles, rAbove)
			queue.enqueue(newQueueEntry(rAbove, aboveObstacles, whitespaceFuzziness))
		}

		rBelow := core.NewPdfRectangleFloat(b.Left(), pivot.Top(), b.Right(), b.Top())
		if b.Top() > pivot.Top() && rBelow.Height > minHeight && rBelow.Width > minWidth {
			belowObstacles := filterOverlapping(current.obstacles, rBelow)
			queue.enqueue(newQueueEntry(rBelow, belowObstacles, whitespaceFuzziness))
		}
	}

	return selected
}

func isInsideAny(selected []core.PdfRectangle, rect core.PdfRectangle) bool {
	for _, s := range selected {
		if inside(s, rect) {
			return true
		}
	}
	return false
}

func isAdjacent(rectangle1, rectangle2 core.PdfRectangle) bool {
	if rectangle1.Left() > rectangle2.Right() ||
		rectangle2.Left() > rectangle1.Right() ||
		rectangle1.Top() < rectangle2.Bottom() ||
		rectangle2.Top() < rectangle1.Bottom() {
		return false
	}

	return AlmostEquals(rectangle1.Left(), rectangle2.Right()) ||
		AlmostEquals(rectangle1.Right(), rectangle2.Left()) ||
		AlmostEquals(rectangle1.Bottom(), rectangle2.Top()) ||
		AlmostEquals(rectangle1.Top(), rectangle2.Bottom())
}

func isAdjacentToPageBounds(pageBound, rectangle core.PdfRectangle) bool {
	return AlmostEquals(rectangle.Bottom(), pageBound.Bottom()) ||
		AlmostEquals(rectangle.Top(), pageBound.Top()) ||
		AlmostEquals(rectangle.Left(), pageBound.Left()) ||
		AlmostEquals(rectangle.Right(), pageBound.Right())
}

func anyAdjacent(selected []core.PdfRectangle, rect core.PdfRectangle) bool {
	for _, s := range selected {
		if isAdjacent(s, rect) {
			return true
		}
	}
	return false
}

func overlapsHard(rectangle1, rectangle2 core.PdfRectangle) bool {
	return rectangle1.Left() < rectangle2.Right() &&
		rectangle2.Left() < rectangle1.Right() &&
		rectangle1.Top() > rectangle2.Bottom() &&
		rectangle2.Top() > rectangle1.Bottom()
}

func inside(rectangle1, rectangle2 core.PdfRectangle) bool {
	return rectangle2.Right() <= rectangle1.Right() &&
		rectangle2.Left() >= rectangle1.Left() &&
		rectangle2.Top() <= rectangle1.Top() &&
		rectangle2.Bottom() >= rectangle1.Bottom()
}

func getBound(obstacles []core.PdfRectangle) core.PdfRectangle {
	if len(obstacles) == 0 {
		return core.PdfRectangle{}
	}

	minLeft := obstacles[0].Left()
	minBottom := obstacles[0].Bottom()
	maxRight := obstacles[0].Right()
	maxTop := obstacles[0].Top()

	for _, o := range obstacles[1:] {
		if o.Left() < minLeft {
			minLeft = o.Left()
		}
		if o.Bottom() < minBottom {
			minBottom = o.Bottom()
		}
		if o.Right() > maxRight {
			maxRight = o.Right()
		}
		if o.Top() > maxTop {
			maxTop = o.Top()
		}
	}

	return core.NewPdfRectangleFloat(minLeft, minBottom, maxRight, maxTop)
}

func filterOverlapping(obstacles []core.PdfRectangle, rect core.PdfRectangle) []core.PdfRectangle {
	var result []core.PdfRectangle
	for _, o := range obstacles {
		if overlapsHard(rect, o) {
			result = append(result, o)
		}
	}
	return result
}

// queueEntry represents an entry in the priority queue for the whitespace algorithm.
type queueEntry struct {
	bound               core.PdfRectangle
	obstacles           []core.PdfRectangle
	quality             float64
	whitespaceFuzziness float64
}

func newQueueEntry(bound core.PdfRectangle, obstacles []core.PdfRectangle, whitespaceFuzziness float64) *queueEntry {
	return &queueEntry{
		bound:               bound,
		obstacles:           obstacles,
		quality:             scoringFunction(bound),
		whitespaceFuzziness: whitespaceFuzziness,
	}
}

func (e *queueEntry) getPivot(pageObstacles []core.PdfRectangle) core.PdfRectangle {
	if len(e.obstacles) == 0 {
		return e.bound
	}

	centroids := make([]core.PdfPoint, len(e.obstacles))
	for i, o := range e.obstacles {
		centroids[i] = o.Centroid()
	}

	idx, _ := FindIndexNearestPoint(
		e.bound.Centroid(),
		centroids,
		func(p core.PdfPoint) core.PdfPoint { return p },
		func(p core.PdfPoint) core.PdfPoint { return p },
		Euclidean,
	)

	if idx == -1 || idx >= len(e.obstacles) {
		return e.obstacles[0]
	}
	return e.obstacles[idx]
}

func (e *queueEntry) isEmptyEnough(pageObstacles []core.PdfRectangle) bool {
	if len(e.obstacles) == 0 {
		return true
	}

	var sum float64
	for _, obstacle := range pageObstacles {
		intersect, ok := rectIntersection(e.bound, obstacle)
		if !ok {
			return false
		}

		minArea := min(obstacle.Area(), e.bound.Area()) * e.whitespaceFuzziness

		if intersect.Area() > minArea {
			return false
		}
		sum += intersect.Area()
	}
	return sum < e.bound.Area()*e.whitespaceFuzziness
}

func (e *queueEntry) addWhitespace(rectangle core.PdfRectangle) {
	e.obstacles = append(e.obstacles, rectangle)
}

// scoringFunction computes the quality score Q(r) used to sort the priority queue.
// Tall rectangles are preferred while still allowing wide ones to be chosen.
func scoringFunction(rectangle core.PdfRectangle) float64 {
	return rectangle.Area() * (rectangle.Height / 4.0)
}

// sortedQueue implements a bounded priority queue sorted by entry quality (descending).
type sortedQueue struct {
	bound   int
	entries []*queueEntry
}

func newSortedQueue(maximumBound int) *sortedQueue {
	return &sortedQueue{bound: maximumBound}
}

func (q *sortedQueue) any() bool {
	return len(q.entries) > 0
}

func (q *sortedQueue) dequeue() *queueEntry {
	if len(q.entries) == 0 {
		panic("dequeue from empty sortedQueue")
	}

	sort.SliceStable(q.entries, func(i, j int) bool {
		return q.entries[i].quality > q.entries[j].quality
	})

	current := q.entries[0]
	q.entries = q.entries[1:]
	return current
}

func (q *sortedQueue) enqueue(entry *queueEntry) {
	if q.bound > 0 && len(q.entries) >= q.bound {
		sort.SliceStable(q.entries, func(i, j int) bool {
			return q.entries[i].quality < q.entries[j].quality
		})
		q.entries = q.entries[1:]
	}

	for _, e := range q.entries {
		if e == entry {
			return
		}
	}

	q.entries = append(q.entries, entry)
}

func (q *sortedQueue) all() []*queueEntry {
	return q.entries
}
