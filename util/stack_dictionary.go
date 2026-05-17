package util

import "fmt"

// StackDictionary is a stack-based dictionary that searches from top to bottom
// across pushed scopes. The most recently pushed scope takes priority on lookup.
type StackDictionary[K comparable, V any] struct {
	values []map[K]V
}

// NewStackDictionary creates an empty StackDictionary.
func NewStackDictionary[K comparable, V any]() *StackDictionary[K, V] {
	return &StackDictionary[K, V]{values: make([]map[K]V, 0)}
}

// Get returns the value for key, searching from top scope downward.
// Returns an error if the stack is empty or the key is not found.
func (sd *StackDictionary[K, V]) Get(key K) (V, error) {
	if len(sd.values) == 0 {
		var zero V
		return zero, fmt.Errorf("cannot get item from empty stack, call Push before use")
	}

	val, ok := sd.TryGetValue(key)
	if !ok {
		var zero V
		return zero, fmt.Errorf("no item with key %v in stack", key)
	}

	return val, nil
}

// Set sets the value for key in the topmost scope.
func (sd *StackDictionary[K, V]) Set(key K, value V) error {
	if len(sd.values) == 0 {
		return fmt.Errorf("cannot set item in empty stack, call Push before use")
	}

	sd.values[len(sd.values)-1][key] = value
	return nil
}

// TryGetValue searches for key from top scope downward. Returns the value and true
// if found, or zero value and false otherwise.
func (sd *StackDictionary[K, V]) TryGetValue(key K) (V, bool) {
	if len(sd.values) == 0 {
		var zero V
		return zero, false
	}

	for i := len(sd.values) - 1; i >= 0; i-- {
		if val, ok := sd.values[i][key]; ok {
			return val, true
		}
	}

	var zero V
	return zero, false
}

// Push adds a new empty scope on top of the stack.
func (sd *StackDictionary[K, V]) Push() {
	sd.values = append(sd.values, make(map[K]V))
}

// Pop removes the topmost scope from the stack.
// Returns an error if the stack is empty.
func (sd *StackDictionary[K, V]) Pop() error {
	if len(sd.values) == 0 {
		return fmt.Errorf("cannot pop empty stacked dictionary")
	}

	sd.values = sd.values[:len(sd.values)-1]
	return nil
}
