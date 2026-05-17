package tokens

// NullToken represents a null object in a PDF document.
// The null object has a type and value unequal to those of any other object.
// There is only one instance of the null token, accessible via the Null var.
type NullToken struct{}

var _ Token = (*NullToken)(nil)
var _ DataToken[any] = (*NullToken)(nil)

// Null is the singleton instance of the null token.
var Null = &NullToken{}

// Data returns nil, representing the PDF null value.
func (t *NullToken) Data() any {
	return nil
}

// Equals reports whether other is a NullToken.
func (t *NullToken) Equals(other Token) bool {
	_, ok := other.(*NullToken)
	return ok
}

// String returns the string representation of the null token.
func (t *NullToken) String() string {
	return "null"
}
