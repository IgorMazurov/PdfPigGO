package tokens

import (
	"github.com/uglytoad/pdfpig/go/core"
)

// IndirectReferenceToken represents a reference to an indirect object in a PDF document.
type IndirectReferenceToken struct {
	data core.IndirectReference
}

var _ Token = (*IndirectReferenceToken)(nil)
var _ DataToken[core.IndirectReference] = (*IndirectReferenceToken)(nil)

// NewIndirectReferenceToken creates a new IndirectReferenceToken with the given data.
func NewIndirectReferenceToken(data core.IndirectReference) *IndirectReferenceToken {
	return &IndirectReferenceToken{data: data}
}

// Data returns the indirect reference identifier for the object this token references.
func (t *IndirectReferenceToken) Data() core.IndirectReference {
	return t.data
}

// Equals reports whether other is an IndirectReferenceToken with the same Data value.
func (t *IndirectReferenceToken) Equals(other Token) bool {
	o, ok := other.(*IndirectReferenceToken)
	if !ok {
		return false
	}
	return t.data.Equals(o.data)
}

// String returns the string representation of the indirect reference token.
func (t *IndirectReferenceToken) String() string {
	return t.data.String()
}
