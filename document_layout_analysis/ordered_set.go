package document_layout_analysis

// OrderedSet maintains O(1) membership tests while preserving insertion order,
// backed by both a map for lookups and a slice for ordering.
type OrderedSet[T comparable] struct {
	set  map[T]struct{}
	list []T
}

// NewOrderedSet creates an empty OrderedSet.
func NewOrderedSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		set:  make(map[T]struct{}),
		list: make([]T, 0),
	}
}

// Count returns the number of elements in the set.
func (o *OrderedSet[T]) Count() int {
	return len(o.set)
}

// TryAdd adds an item if it is not already present. Returns true when the item was added.
func (o *OrderedSet[T]) TryAdd(item T) bool {
	if _, exists := o.set[item]; exists {
		return false
	}
	o.list = append(o.list, item)
	o.set[item] = struct{}{}
	return true
}

// Clear removes all elements from the set.
func (o *OrderedSet[T]) Clear() {
	clear(o.set)
	o.list = o.list[:0]
}

// Contains reports whether the set contains the given item.
func (o *OrderedSet[T]) Contains(item T) bool {
	_, exists := o.set[item]
	return exists
}

// CopyTo copies the elements to the given slice starting at arrayIndex.
func (o *OrderedSet[T]) CopyTo(arr []T, arrayIndex int) {
	copy(arr[arrayIndex:], o.list)
}

// GetList returns the internal ordered list of elements.
func (o *OrderedSet[T]) GetList() []T {
	return o.list
}
