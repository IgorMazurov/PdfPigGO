package tokens

// CommentToken represents a comment from a PDF document. Any occurrence of the
// percent sign character (%) outside a string or stream introduces a comment.
// The comment consists of all characters between the percent sign and the end of the line.
type CommentToken struct {
	data string
}

var _ Token = (*CommentToken)(nil)
var _ DataToken[string] = (*CommentToken)(nil)

// NewCommentToken creates a new CommentToken with the given data.
func NewCommentToken(data string) *CommentToken {
	if data == "" {
		return &CommentToken{data: ""}
	}
	return &CommentToken{data: data}
}

// Data returns the text of the comment excluding the initial percent '%' sign.
func (t *CommentToken) Data() string {
	return t.data
}

// Equals reports whether other is a CommentToken with the same Data value.
func (t *CommentToken) Equals(other Token) bool {
	o, ok := other.(*CommentToken)
	if !ok {
		return false
	}
	return o.data == t.data
}

// String returns the text of the comment.
func (t *CommentToken) String() string {
	return t.data
}
