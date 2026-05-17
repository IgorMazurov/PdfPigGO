package document_layout_analysis

import (
	"sort"
)

// ArraySegmentTake returns a new ArraySegment containing the first count elements from source.
func ArraySegmentTake[T any](source ArraySegment[T], count int) ArraySegment[T] {
	return NewArraySegment(source.Array, source.Offset, count)
}

// ArraySegmentSkip returns a new ArraySegment with the first count elements skipped.
func ArraySegmentSkip[T any](source ArraySegment[T], count int) ArraySegment[T] {
	return NewArraySegment(source.Array, source.Offset+count, source.Count-count)
}

// ArraySegmentSort sorts the elements in source using the provided less function.
// The less function receives two elements and should return true if a < b.
func ArraySegmentSort[T any](source *ArraySegment[T], less func(a, b T) bool) {
	s := source.Slice()
	sort.SliceStable(s, func(i, j int) bool {
		return less(s[i], s[j])
	})
}

// ArraySegmentGetAt returns the element at the specified index within the segment.
func ArraySegmentGetAt[T any](source ArraySegment[T], index int) (T, error) {
	return source.Array[source.Offset+index], nil
}
