package util

// RecursiveChildFunc defines a function that returns the children of an item.
type RecursiveChildFunc[T any] func(T) []T

// ToRecursiveOrderList flattens a hierarchical collection into a single list,
// preserving recursive order by inserting children immediately after their parent.
// The childFunc callback retrieves the children for each item.
func ToRecursiveOrderList[T any](collection []T, childFunc RecursiveChildFunc[T]) []T {
	type queueItem struct {
		index int
		item  T
		depth int
	}

	resultList := make([]T, 0)
	currentItems := make([]queueItem, 0, len(collection))

	for _, item := range collection {
		currentItems = append(currentItems, queueItem{index: 0, item: item, depth: 0})
	}

	depthItemCounter := 0
	previousItemDepth := 0

	for len(currentItems) > 0 {
		currentItem := currentItems[0]
		currentItems = currentItems[1:]

		if currentItem.depth != previousItemDepth {
			depthItemCounter = 0
		}

		resultIndex := currentItem.index + depthItemCounter
		depthItemCounter++

		resultList = insertAt(resultList, resultIndex, currentItem.item)

		childItems := childFunc(currentItem.item)
		for _, childItem := range childItems {
			currentItems = append(currentItems, queueItem{
				index: resultIndex + 1,
				item:  childItem,
				depth: currentItem.depth + 1,
			})
		}

		previousItemDepth = currentItem.depth
	}

	return resultList
}

// insertAt inserts an element at the given index in a slice.
func insertAt[T any](slice []T, index int, item T) []T {
	if index >= len(slice) {
		return append(slice, item)
	}
	slice = append(slice, *new(T))
	copy(slice[index+1:], slice[index:])
	slice[index] = item
	return slice
}
