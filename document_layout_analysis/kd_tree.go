package document_layout_analysis

import (
	"errors"
	"math"
	"sort"

	"github.com/uglytoad/pdfpig/go/core"
)

// kdTreeElement holds an element with its original index and associated point.
type kdTreeElement[T any] struct {
	index   int
	value   core.PdfPoint
	element T
}

// kdTreeNode represents a node in the KD-tree.
type kdTreeNode[T any] struct {
	leftChild  *kdTreeNode[T]
	rightChild *kdTreeNode[T]
	value      core.PdfPoint
	element    T
	depth      int
	isAxisCutX bool
	index      int
	isLeaf     bool
}

// L returns the split value (X or Y axis depending on depth).
func (n *kdTreeNode[T]) L() float64 {
	if n.isAxisCutX {
		return n.value.X
	}
	return n.value.Y
}

// KdTree is a KD-tree data structure specialized for PdfPoint.
// It provides efficient nearest neighbour search in 2D space.
type KdTree struct {
	root  *kdTreeNode[core.PdfPoint]
	count int
}

// NewKdTree creates a new KD-tree from the given points.
// Returns an error if points is nil or empty.
func NewKdTree(points []core.PdfPoint) (*KdTree, error) {
	if len(points) == 0 {
		return nil, errors.New("NewKdTree: points cannot be null or empty")
	}

	count := len(points)
	arr := make([]kdTreeElement[core.PdfPoint], count)
	for i, p := range points {
		arr[i] = kdTreeElement[core.PdfPoint]{index: i, value: p, element: p}
	}

	root := buildTree(arr, 0)

	return &KdTree{root: root, count: count}, nil
}

// Root returns the root node of the tree.
func (t *KdTree) Root() *kdTreeNode[core.PdfPoint] {
	return t.root
}

// Count returns the number of elements in the tree.
func (t *KdTree) Count() int {
	return t.count
}

// FindNearestNeighbour finds the nearest neighbour to the pivot point.
// Only returns 1 neighbour, even if equidistant points are found.
// Returns (zero-value, -1, NaN) if no neighbour is found.
func (t *KdTree) FindNearestNeighbour(pivot core.PdfPoint, distanceMeasure func(core.PdfPoint, core.PdfPoint) float64) (core.PdfPoint, int, float64) {
	resultNode, dist := findNearestOne[core.PdfPoint](t.root, pivot, pivot, distanceMeasure)
	if resultNode == nil || dist == nil {
		return core.PdfPoint{}, -1, math.NaN()
	}
	return resultNode.value, resultNode.index, *dist
}

// NearestResult holds a single nearest-neighbour search result.
type NearestResult struct {
	Point    core.PdfPoint
	Index    int
	Distance float64
}

// FindNearestNeighbours finds the k nearest neighbours to the pivot point.
// Might return more than k neighbours if points are equidistant.
func (t *KdTree) FindNearestNeighbours(pivot core.PdfPoint, k int, distanceMeasure func(core.PdfPoint, core.PdfPoint) float64) []NearestResult {
	queue := newKNearestQueue[core.PdfPoint](k)
	findNearestMany(t.root, pivot, k, pivot, distanceMeasure, queue)

	seen := make(map[int]bool)
	results := make([]NearestResult, 0, k)
	for i := 0; i < queue.count(); i++ {
		dist := queue.keys()[i]
		nodes := queue.values()[i]
		for _, n := range nodes {
			if !seen[n.index] {
				seen[n.index] = true
				results = append(results, NearestResult{
					Point:    n.value,
					Index:    n.index,
					Distance: dist,
				})
			}
		}
	}
	return results
}

// GenericKdTree is a generic KD-tree data structure.
type GenericKdTree[T comparable] struct {
	root  *kdTreeNode[T]
	count int
}

// NewGenericKdTree creates a new generic KD-tree from the given elements.
// elementToPoint converts each element to its associated PdfPoint for tree construction.
// Returns an error if elements is nil or empty.
func NewGenericKdTree[T comparable](elements []T, elementToPoint func(T) core.PdfPoint) (*GenericKdTree[T], error) {
	if len(elements) == 0 {
		return nil, errors.New("NewGenericKdTree: elements cannot be null or empty")
	}

	count := len(elements)
	arr := make([]kdTreeElement[T], count)
	for i, el := range elements {
		arr[i] = kdTreeElement[T]{index: i, value: elementToPoint(el), element: el}
	}

	root := buildTree(arr, 0)

	return &GenericKdTree[T]{root: root, count: count}, nil
}

// Root returns the root node of the tree.
func (t *GenericKdTree[T]) Root() *kdTreeNode[T] {
	return t.root
}

// Count returns the number of elements in the tree.
func (t *GenericKdTree[T]) Count() int {
	return t.count
}

// FindNearestNeighbour finds the nearest neighbour to the pivot element.
// pivotToPoint converts the pivot to its associated PdfPoint.
// Only returns 1 neighbour, even if equidistant points are found.
// Returns (zero-value, -1, NaN) if no neighbour is found.
func (t *GenericKdTree[T]) FindNearestNeighbour(pivot T, pivotToPoint func(T) core.PdfPoint, distanceMeasure func(core.PdfPoint, core.PdfPoint) float64) (T, int, float64) {
	pivotPoint := pivotToPoint(pivot)
	resultNode, dist := findNearestOne(t.root, pivot, pivotPoint, distanceMeasure)
	if resultNode == nil || dist == nil {
		var zero T
		return zero, -1, math.NaN()
	}
	return resultNode.element, resultNode.index, *dist
}

// GenericNearestResult holds a single nearest-neighbour search result for generic elements.
type GenericNearestResult[T comparable] struct {
	Element  T
	Index    int
	Distance float64
}

// FindNearestNeighbours finds the k nearest neighbours to the pivot element.
// Might return more than k neighbours if points are equidistant.
func (t *GenericKdTree[T]) FindNearestNeighbours(pivot T, k int, pivotToPoint func(T) core.PdfPoint, distanceMeasure func(core.PdfPoint, core.PdfPoint) float64) []GenericNearestResult[T] {
	pivotPoint := pivotToPoint(pivot)
	queue := newKNearestQueue[T](k)
	findNearestMany(t.root, pivot, k, pivotPoint, distanceMeasure, queue)

	seen := make(map[int]bool)
	results := make([]GenericNearestResult[T], 0, k)
	for i := 0; i < queue.count(); i++ {
		dist := queue.keys()[i]
		nodes := queue.values()[i]
		for _, n := range nodes {
			if !seen[n.index] {
				seen[n.index] = true
				results = append(results, GenericNearestResult[T]{
					Element:  n.element,
					Index:    n.index,
					Distance: dist,
				})
			}
		}
	}
	return results
}

// buildTree recursively builds a KD-tree from the given elements.
func buildTree[T any](elements []kdTreeElement[T], depth int) *kdTreeNode[T] {
	if len(elements) == 0 {
		return nil
	}

	if len(elements) == 1 {
		return &kdTreeNode[T]{
			value:      elements[0].value,
			element:    elements[0].element,
			index:      elements[0].index,
			depth:      depth,
			isAxisCutX: depth%2 == 0,
			isLeaf:     true,
		}
	}

	if depth%2 == 0 {
		sort.Slice(elements, func(i, j int) bool {
			if elements[i].value.X != elements[j].value.X {
				return elements[i].value.X < elements[j].value.X
			}
			return elements[i].index < elements[j].index
		})
	} else {
		sort.Slice(elements, func(i, j int) bool {
			if elements[i].value.Y != elements[j].value.Y {
				return elements[i].value.Y < elements[j].value.Y
			}
			return elements[i].index < elements[j].index
		})
	}

	if len(elements) == 2 {
		left := &kdTreeNode[T]{
			value:      elements[0].value,
			element:    elements[0].element,
			index:      elements[0].index,
			depth:      depth + 1,
			isAxisCutX: (depth+1)%2 == 0,
			isLeaf:     true,
		}
		return &kdTreeNode[T]{
			leftChild:  left,
			rightChild: nil,
			value:      elements[1].value,
			element:    elements[1].element,
			index:      elements[1].index,
			depth:      depth,
			isAxisCutX: depth%2 == 0,
			isLeaf:     false,
		}
	}

	median := len(elements) / 2

	vLeft := buildTree(elements[:median], depth+1)
	vRight := buildTree(elements[median+1:], depth+1)

	return &kdTreeNode[T]{
		leftChild:  vLeft,
		rightChild: vRight,
		value:      elements[median].value,
		element:    elements[median].element,
		index:      elements[median].index,
		depth:      depth,
		isAxisCutX: depth%2 == 0,
		isLeaf:     false,
	}
}

// findNearestOne recursively finds the single nearest neighbour to pivotPoint.
func findNearestOne[T comparable](
	node *kdTreeNode[T],
	pivot T,
	pivotPoint core.PdfPoint,
	distance func(core.PdfPoint, core.PdfPoint) float64,
) (*kdTreeNode[T], *float64) {
	if node == nil {
		return nil, nil
	}

	if node.isLeaf {
		if node.element == pivot {
			return nil, nil
		}
		d := distance(node.value, pivotPoint)
		return node, &d
	}

	currentNearestNode := node
	currentDistance := distance(node.value, pivotPoint)

	var newNode *kdTreeNode[T]
	var newDist *float64

	pointValue := pivotPoint.X
	if !node.isAxisCutX {
		pointValue = pivotPoint.Y
	}

	if pointValue < node.L() {
		newNode, newDist = findNearestOne(node.leftChild, pivot, pivotPoint, distance)

		if newDist != nil && *newDist <= currentDistance && newNode.element != pivot {
			currentDistance = *newDist
			currentNearestNode = newNode
		}

		if node.rightChild != nil && pointValue+currentDistance >= node.L() {
			newNode, newDist = findNearestOne(node.rightChild, pivot, pivotPoint, distance)
		}
	} else {
		newNode, newDist = findNearestOne(node.rightChild, pivot, pivotPoint, distance)

		if newDist != nil && *newDist <= currentDistance && newNode.element != pivot {
			currentDistance = *newDist
			currentNearestNode = newNode
		}

		if node.leftChild != nil && pointValue-currentDistance <= node.L() {
			newNode, newDist = findNearestOne(node.leftChild, pivot, pivotPoint, distance)
		}
	}

	if newDist != nil && *newDist <= currentDistance && newNode.element != pivot {
		currentDistance = *newDist
		currentNearestNode = newNode
	}

	return currentNearestNode, &currentDistance
}

// kNearestQueue maintains a sorted collection of distances to nodes for k-NN search.
type kNearestQueue[T comparable] struct {
	distances    []float64
	nodeSets     [][]*kdTreeNode[T]
	k            int
	lastDistance float64
	lastElement  *kdTreeNode[T]
}

func newKNearestQueue[T comparable](k int) *kNearestQueue[T] {
	return &kNearestQueue[T]{
		distances:    make([]float64, 0),
		nodeSets:     make([][]*kdTreeNode[T], 0),
		k:            k,
		lastDistance: math.Inf(1),
	}
}

func (q *kNearestQueue[T]) count() int {
	return len(q.distances)
}

func (q *kNearestQueue[T]) isFull() bool {
	return q.count() >= q.k
}

func (q *kNearestQueue[T]) keys() []float64 {
	return q.distances
}

func (q *kNearestQueue[T]) values() [][]*kdTreeNode[T] {
	return q.nodeSets
}

func (q *kNearestQueue[T]) add(key float64, value *kdTreeNode[T]) {
	if key > q.lastDistance && q.isFull() {
		return
	}

	idx := sort.SearchFloat64s(q.distances, key)
	hasKey := idx < len(q.distances) && q.distances[idx] == key

	if !hasKey {
		q.distances = append(q.distances, 0)
		q.nodeSets = append(q.nodeSets, nil)
		copy(q.distances[idx+1:], q.distances[idx:])
		copy(q.nodeSets[idx+1:], q.nodeSets[idx:])
		q.distances[idx] = key
		q.nodeSets[idx] = []*kdTreeNode[T]{value}

		for len(q.distances) > q.k {
			q.distances = q.distances[:len(q.distances)-1]
			q.nodeSets = q.nodeSets[:len(q.nodeSets)-1]
		}
	} else {
		q.nodeSets[idx] = append(q.nodeSets[idx], value)
	}

	if len(q.distances) > 0 {
		q.lastDistance = q.distances[len(q.distances)-1]
		lastSet := q.nodeSets[len(q.nodeSets)-1]
		q.lastElement = lastSet[len(lastSet)-1]
	}
}

// findNearestMany recursively finds the k nearest neighbours to pivotPoint.
func findNearestMany[T comparable](
	node *kdTreeNode[T],
	pivot T,
	k int,
	pivotPoint core.PdfPoint,
	distance func(core.PdfPoint, core.PdfPoint) float64,
	queue *kNearestQueue[T],
) (*kdTreeNode[T], float64) {
	if node == nil {
		return nil, math.NaN()
	}

	if node.isLeaf {
		if node.element == pivot {
			return nil, math.NaN()
		}

		currentDistance := distance(node.value, pivotPoint)
		currentNearestNode := node

		if !queue.isFull() || currentDistance <= queue.lastDistance {
			queue.add(currentDistance, currentNearestNode)
			currentDistance = queue.lastDistance
			currentNearestNode = queue.lastElement
		}

		return currentNearestNode, currentDistance
	}

	currentNearestNode := node
	currentDistance := distance(node.value, pivotPoint)
	if (!queue.isFull() || currentDistance <= queue.lastDistance) && node.element != pivot {
		queue.add(currentDistance, currentNearestNode)
		currentDistance = queue.lastDistance
		currentNearestNode = queue.lastElement
	}

	var newNode *kdTreeNode[T]
	newDist := math.NaN()

	pointValue := pivotPoint.X
	if !node.isAxisCutX {
		pointValue = pivotPoint.Y
	}

	if pointValue < node.L() {
		newNode, newDist = findNearestMany(node.leftChild, pivot, k, pivotPoint, distance, queue)

		if !math.IsNaN(newDist) && newDist <= currentDistance && newNode.element != pivot {
			queue.add(newDist, newNode)
			currentDistance = queue.lastDistance
			currentNearestNode = queue.lastElement
		}

		if node.rightChild != nil && pointValue+currentDistance >= node.L() {
			newNode, newDist = findNearestMany(node.rightChild, pivot, k, pivotPoint, distance, queue)
		}
	} else {
		newNode, newDist = findNearestMany(node.rightChild, pivot, k, pivotPoint, distance, queue)

		if !math.IsNaN(newDist) && newDist <= currentDistance && newNode.element != pivot {
			queue.add(newDist, newNode)
			currentDistance = queue.lastDistance
			currentNearestNode = queue.lastElement
		}

		if node.leftChild != nil && pointValue-currentDistance <= node.L() {
			newNode, newDist = findNearestMany(node.leftChild, pivot, k, pivotPoint, distance, queue)
		}
	}

	if !math.IsNaN(newDist) && newDist <= currentDistance && newNode.element != pivot {
		queue.add(newDist, newNode)
		currentDistance = queue.lastDistance
		currentNearestNode = queue.lastElement
	}

	return currentNearestNode, currentDistance
}

// GetLeaves returns all leaf nodes in the subtree rooted at this node.
func (n *kdTreeNode[T]) GetLeaves() []*kdTreeNode[T] {
	var leaves []*kdTreeNode[T]
	getLeaves(n.leftChild, &leaves)
	getLeaves(n.rightChild, &leaves)
	return leaves
}

func getLeaves[T any](node *kdTreeNode[T], leaves *[]*kdTreeNode[T]) {
	if node == nil {
		return
	}
	if node.isLeaf {
		*leaves = append(*leaves, node)
	} else {
		getLeaves(node.leftChild, leaves)
		getLeaves(node.rightChild, leaves)
	}
}
