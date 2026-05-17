package tokens

import "bytes"

// InlineImageDataToken represents inline image data embedded in a PDF content stream.
// The content is wrapped by ID and ED tags in a BI operation.
type InlineImageDataToken struct {
	data []byte
}

var _ Token = (*InlineImageDataToken)(nil)
var _ DataToken[[]byte] = (*InlineImageDataToken)(nil)

// NewInlineImageDataToken creates a new InlineImageDataToken with the given data.
func NewInlineImageDataToken(data []byte) *InlineImageDataToken {
	return &InlineImageDataToken{data: data}
}

// Data returns the raw image byte slice of this token.
func (t *InlineImageDataToken) Data() []byte {
	return t.data
}

// Equals reports whether other is an InlineImageDataToken with identical data bytes.
func (t *InlineImageDataToken) Equals(other Token) bool {
	if _, ok := other.(*InlineImageDataToken); !ok {
		return false
	}
	o := other.(*InlineImageDataToken)
	return bytes.Equal(t.data, o.data)
}
