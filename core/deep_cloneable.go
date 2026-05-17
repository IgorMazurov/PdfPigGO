package core

// DeepCloneable indicates the type may be cloned, making an entirely independent copy.
type DeepCloneable[T any] interface {
	// Clone returns a deep clone of the value including all referenced data.
	Clone() T
}
