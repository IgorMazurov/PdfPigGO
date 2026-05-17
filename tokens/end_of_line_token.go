package tokens

// EndOfLineToken represents an end-of-line marker found in Adobe Type 1 font files and the cross-reference table.
type EndOfLineToken struct{}

var _ Token = (*EndOfLineToken)(nil)

// EOLToken is the singleton instance of the end-of-line token.
var EOLToken = &EndOfLineToken{}

// Equals reports whether other is an EndOfLineToken.
func (t *EndOfLineToken) Equals(other Token) bool {
	_, ok := other.(*EndOfLineToken)
	return ok
}
