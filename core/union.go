package core

import "fmt"

// Union defines a discriminated union (sum type) of two types A and B.
type Union[A any, B any] interface {
	// Match calls first if the union holds an A value, or second if it holds a B value.
	// The functions may return values; the result is returned as any.
	Match(first func(A) any, second func(B) any) any

	// TryGetFirst returns the inner value if the union holds type A, otherwise false.
	TryGetFirst() (A, bool)

	// TryGetSecond returns the inner value if the union holds type B, otherwise false.
	TryGetSecond() (B, bool)
}

// Case1 represents a Union[A, B] holding a value of type A.
type Case1[A any, B any] struct {
	Item A
}

func (c Case1[A, B]) Match(first func(A) any, _ func(B) any) any {
	return first(c.Item)
}

func (c Case1[A, B]) TryGetFirst() (A, bool) {
	return c.Item, true
}

func (c Case1[A, B]) TryGetSecond() (B, bool) {
	var zero B
	return zero, false
}

func (c Case1[A, B]) String() string {
	return fmt.Sprintf("%v", c.Item)
}

// Case2 represents a Union[A, B] holding a value of type B.
type Case2[A any, B any] struct {
	Item B
}

func (c Case2[A, B]) Match(_ func(A) any, second func(B) any) any {
	return second(c.Item)
}

func (c Case2[A, B]) TryGetFirst() (A, bool) {
	var zero A
	return zero, false
}

func (c Case2[A, B]) TryGetSecond() (B, bool) {
	return c.Item, true
}

func (c Case2[A, B]) String() string {
	return fmt.Sprintf("%v", c.Item)
}

// UnionOne creates a Union[A, B] containing an A value.
func UnionOne[A any, B any](item A) Union[A, B] {
	return Case1[A, B]{Item: item}
}

// UnionTwo creates a Union[A, B] containing a B value.
func UnionTwo[A any, B any](item B) Union[A, B] {
	return Case2[A, B]{Item: item}
}
