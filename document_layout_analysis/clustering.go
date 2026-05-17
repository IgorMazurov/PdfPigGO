package document_layout_analysis

import (
	"math"
	"runtime"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
)

// NearestNeighboursPoint clusters elements by finding their nearest neighbour in point space.
// Uses a KdTree for efficient spatial search. Each element is mapped to a pivot point and
// a candidate point; the distance measure compares two points. Elements whose nearest
// neighbour falls within maxDistanceFunction and passes both filters are grouped together
// via connected-component analysis on the neighbour graph.
func NearestNeighboursPoint[T comparable](
	elements []T,
	distMeasure func(core.PdfPoint, core.PdfPoint) float64,
	maxDistanceFunction func(T, T) float64,
	pivotPoint func(T) core.PdfPoint,
	candidatesPoint func(T) core.PdfPoint,
	filterPivot func(T) bool,
	filterFinal func(T, T) bool,
	maxDegreeOfParallelism int,
) [][]T {
	if len(elements) == 0 {
		return nil
	}

	indexes := make([]int, len(elements))
	for i := range indexes {
		indexes[i] = -1
	}

	kdTree, err := NewGenericKdTree(elements, candidatesPoint)
	if err != nil {
		return nil
	}

	sem := make(chan struct{}, resolveParallelism(maxDegreeOfParallelism))
	var wg sync.WaitGroup

	for i := 0; i < len(elements); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			pivot := elements[idx]
			if !filterPivot(pivot) {
				return
			}

			paired, pairIdx, dist := kdTree.FindNearestNeighbour(pivot, pivotPoint, distMeasure)
			if pairIdx != -1 && filterFinal(pivot, paired) && dist < maxDistanceFunction(pivot, paired) {
				indexes[idx] = pairIdx
			}
		}(i)
	}

	wg.Wait()

	grouped := GroupIndexes(indexes)
	result := make([][]T, len(grouped))
	for i, g := range grouped {
		group := make([]T, len(g))
		for j, idx := range g {
			group[j] = elements[idx]
		}
		result[i] = group
	}
	return result
}

// NearestNeighboursPointK clusters elements by considering the k-nearest neighbours as
// candidates. For each element it iterates through its k nearest neighbours (ordered by
// distance) and picks the first one that passes both filters and falls within the maximum
// distance. Uses a KdTree for efficient spatial search.
func NearestNeighboursPointK[T comparable](
	elements []T,
	k int,
	distMeasure func(core.PdfPoint, core.PdfPoint) float64,
	maxDistanceFunction func(T, T) float64,
	pivotPoint func(T) core.PdfPoint,
	candidatesPoint func(T) core.PdfPoint,
	filterPivot func(T) bool,
	filterFinal func(T, T) bool,
	maxDegreeOfParallelism int,
) [][]T {
	if len(elements) == 0 || k <= 0 {
		return nil
	}

	indexes := make([]int, len(elements))
	for i := range indexes {
		indexes[i] = -1
	}

	kdTree, err := NewGenericKdTree(elements, candidatesPoint)
	if err != nil {
		return nil
	}

	sem := make(chan struct{}, resolveParallelism(maxDegreeOfParallelism))
	var wg sync.WaitGroup

	for i := 0; i < len(elements); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			pivot := elements[idx]
			if !filterPivot(pivot) {
				return
			}

			neighbours := kdTree.FindNearestNeighbours(pivot, k, pivotPoint, distMeasure)
			for _, n := range neighbours {
				if filterFinal(pivot, n.Element) && n.Distance < maxDistanceFunction(pivot, n.Element) {
					indexes[idx] = n.Index
					break
				}
			}
		}(i)
	}

	wg.Wait()

	grouped := GroupIndexes(indexes)
	result := make([][]T, len(grouped))
	for i, g := range grouped {
		group := make([]T, len(g))
		for j, idx := range g {
			group[j] = elements[idx]
		}
		result[i] = group
	}
	return result
}

// NearestNeighboursLine clusters elements by finding their nearest neighbour using a line-
// based distance measure. Unlike the point-based variants this does not use a KdTree and
// instead performs a linear scan via FindIndexNearestLine. Each element is mapped to a
// pivot line and a candidate line; the distance measure compares two lines. Elements whose
// nearest neighbour falls within maxDistanceFunction and passes both filters are grouped
// together via connected-component analysis on the neighbour graph.
func NearestNeighboursLine[T comparable](
	elements []T,
	distMeasure func(core.PdfLine, core.PdfLine) float64,
	maxDistanceFunction func(T, T) float64,
	pivotLine func(T) core.PdfLine,
	candidatesLine func(T) core.PdfLine,
	filterPivot func(T) bool,
	filterFinal func(T, T) bool,
	maxDegreeOfParallelism int,
) [][]T {
	if len(elements) == 0 {
		return nil
	}

	indexes := make([]int, len(elements))
	for i := range indexes {
		indexes[i] = -1
	}

	sem := make(chan struct{}, resolveParallelism(maxDegreeOfParallelism))
	var wg sync.WaitGroup

	for i := 0; i < len(elements); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			pivot := elements[idx]
			if !filterPivot(pivot) {
				return
			}

			pairIdx, dist := FindIndexNearestLine(pivot, elements, pivotLine, candidatesLine, distMeasure)
			if pairIdx != -1 {
				paired := elements[pairIdx]
				if filterFinal(pivot, paired) && dist < maxDistanceFunction(pivot, paired) {
					indexes[idx] = pairIdx
				}
			}
		}(i)
	}

	wg.Wait()

	grouped := GroupIndexes(indexes)
	result := make([][]T, len(grouped))
	for i, g := range grouped {
		group := make([]T, len(g))
		for j, idx := range g {
			group[j] = elements[idx]
		}
		result[i] = group
	}
	return result
}

// GroupIndexes groups element indexes that share neighbours in common using a depth-first
// search on the undirected neighbour graph. edges[i] == j means elements i and j are
// connected; edges[i] == -1 means no connection. Returns a list of index groups where each
// group forms a connected component.
func GroupIndexes(edges []int) [][]int {
	adjacency := make([][]int, len(edges))
	for i := range adjacency {
		adjacency[i] = make([]int, 0)
	}

	for i := 0; i < len(edges); i++ {
		j := edges[i]
		if j != -1 {
			adjacency[i] = append(adjacency[i], j)
			adjacency[j] = append(adjacency[j], i)
		}
	}

	var grouped [][]int
	isDone := make([]bool, len(edges))

	for p := 0; p < len(edges); p++ {
		if isDone[p] {
			continue
		}
		grouped = append(grouped, dfsIterative(p, adjacency, isDone))
	}
	return grouped
}

// resolveParallelism returns the effective concurrency limit. A value of -1 or less means
// unlimited (uses all available CPU cores).
func resolveParallelism(maxDegreeOfParallelism int) int {
	if maxDegreeOfParallelism <= 0 || maxDegreeOfParallelism == math.MaxInt32 {
		return runtime.GOMAXPROCS(0)
	}
	return maxDegreeOfParallelism
}

// dfsIterative performs an iterative depth-first search from start node s on the given
// adjacency list. Returns all reachable nodes as a connected component.
func dfsIterative(s int, adj [][]int, isDone []bool) []int {
	var group []int
	stack := []int{s}
	isDone[s] = true

	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		group = append(group, u)

		for _, v := range adj[u] {
			if !isDone[v] {
				isDone[v] = true
				stack = append(stack, v)
			}
		}
	}
	return group
}
