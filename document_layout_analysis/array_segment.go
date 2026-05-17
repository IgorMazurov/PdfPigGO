package document_layout_analysis

// ArraySegment represents a segment of an array, mirroring System.ArraySegment{T}.
type ArraySegment[T any] struct {
	Array  []T
	Offset int
	Count  int
}

// NewArraySegment creates a new ArraySegment wrapping the given slice.
func NewArraySegment[T any](arr []T, offset, count int) ArraySegment[T] {
	return ArraySegment[T]{
		Array:  arr,
		Offset: offset,
		Count:  count,
	}
}

// Slice returns a Go slice view of the segment's elements.
func (s ArraySegment[T]) Slice() []T {
	return s.Array[s.Offset : s.Offset+s.Count]
}
