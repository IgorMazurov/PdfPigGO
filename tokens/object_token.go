package tokens

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
)

// ObjectToken represents an indirect object in a PDF file.
// Indirect objects have a unique identifier allowing other objects to reference them,
// and contain inner data of any token type.
type ObjectToken struct {
	position core.XrefLocation
	number   core.IndirectReference
	data     Token
}

var _ Token = (*ObjectToken)(nil)
var _ DataToken[Token] = (*ObjectToken)(nil)

// NewObjectToken creates a new ObjectToken with the given position, identifier, and inner data.
func NewObjectToken(position core.XrefLocation, number core.IndirectReference, data Token) *ObjectToken {
	return &ObjectToken{
		position: position,
		number:   number,
		data:     data,
	}
}

// Position returns the offset to the start of the object number from the start of the file in bytes.
func (t *ObjectToken) Position() core.XrefLocation {
	return t.position
}

// Number returns the object and generation number of the indirect object.
func (t *ObjectToken) Number() core.IndirectReference {
	return t.number
}

// Data returns the inner data contained in this object token.
func (t *ObjectToken) Data() Token {
	return t.data
}

// Equals reports whether other is an ObjectToken with the same Number and Data.
func (t *ObjectToken) Equals(other Token) bool {
	o, ok := other.(*ObjectToken)
	if !ok {
		return false
	}
	return t.number.Equals(o.number) && equalTokens(t.data, o.data)
}

// String returns the string representation of the object token.
func (t *ObjectToken) String() string {
	return fmt.Sprintf("Number: %s, Position: %+v, Type: %T", t.number, t.position, t.data)
}

// equalTokens compares two tokens for equality.
func equalTokens(a, b Token) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if eq, ok := a.(interface{ Equals(Token) bool }); ok {
		return eq.Equals(b)
	}
	return false
}
