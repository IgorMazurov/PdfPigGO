package tokens

// BooleanToken represents a boolean value token in a PDF document.
type BooleanToken struct {
	data bool
}

var _ Token = (*BooleanToken)(nil)
var _ DataToken[bool] = (*BooleanToken)(nil)

// True is the singleton instance representing the boolean true token.
var True = &BooleanToken{data: true}

// False is the singleton instance representing the boolean false token.
var False = &BooleanToken{data: false}

// NewBooleanToken creates a new BooleanToken with the given value.
func NewBooleanToken(data bool) *BooleanToken {
	return &BooleanToken{data: data}
}

// Data returns the boolean value of this token.
func (t *BooleanToken) Data() bool {
	return t.data
}

// Equals reports whether other is a BooleanToken with the same Data value.
func (t *BooleanToken) Equals(other Token) bool {
	o, ok := other.(*BooleanToken)
	if !ok {
		return false
	}
	return o.data == t.data
}

// HashCode returns a hash code for this token matching the data value.
func (t *BooleanToken) HashCode() int {
	if t.data {
		return 1
	}
	return 0
}

// String returns the string representation of the boolean token.
func (t *BooleanToken) String() string {
	if t.data {
		return "True"
	}
	return "False"
}
