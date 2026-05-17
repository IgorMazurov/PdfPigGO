package tokens

// DataToken represents a token from a PDF document which contains data in some format.
type DataToken[T any] interface {
	Token
	Data() T
}
