package tokens

import (
	"fmt"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
)

// ArrayToken represents a PDF array, a one-dimensional collection of objects arranged sequentially.
// PDF arrays may be heterogeneous; elements can be any combination of numbers, strings,
// dictionaries, or other objects, including nested arrays.
type ArrayToken struct {
	data []Token
}

var _ Token = (*ArrayToken)(nil)
var _ DataToken[[]Token] = (*ArrayToken)(nil)

// NewArrayToken creates a new ArrayToken from the given slice of tokens.
// During construction, sequences matching [NumericToken][NumericToken] OpR are
// collapsed into a single IndirectReferenceToken.
func NewArrayToken(data []Token) *ArrayToken {
	if data == nil {
		panic("data cannot be nil")
	}

	result := make([]Token, 0, len(data))
	for i := 0; i < len(data); i++ {
		token := data[i]

		if i >= 2 && token == OpR {
			if gen, ok := data[i-1].(*NumericToken); ok {
				if objNum, ok := data[i-2].(*NumericToken); ok {
					result = result[:len(result)-2]

					ref, err := core.NewIndirectReference(objNum.LongVal(), gen.IntVal())
					if err == nil {
						result = append(result, NewIndirectReferenceToken(ref))
					} else {
						result = append(result, token)
					}
					continue
				}
			}
		}

		result = append(result, token)
	}

	return &ArrayToken{data: result}
}

// Data returns the tokens contained in this array.
func (a *ArrayToken) Data() []Token {
	return a.data
}

// Length returns the number of tokens in this array.
func (a *ArrayToken) Length() int {
	return len(a.data)
}

// Get returns the token at the given index.
func (a *ArrayToken) Get(i int) Token {
	return a.data[i]
}

// Equals reports whether other is an ArrayToken with equal elements.
func (a *ArrayToken) Equals(other Token) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*ArrayToken)
	if !ok {
		return false
	}
	if len(a.data) != len(o.data) {
		return false
	}
	for i := range a.data {
		if !tokensEqual(a.data[i], o.data[i]) {
			return false
		}
	}
	return true
}

// String returns the string representation of the array token.
func (a *ArrayToken) String() string {
	var b strings.Builder
	b.WriteString("[ ")

	for i, token := range a.data {
		if s, ok := any(token).(fmt.Stringer); ok {
			b.WriteString(s.String())
		} else {
			b.WriteString(fmt.Sprintf("%v", token))
		}

		if i < len(a.data)-1 {
			b.WriteByte(',')
		}
		b.WriteByte(' ')
	}

	b.WriteByte(']')
	return b.String()
}

func tokensEqual(a, b Token) bool {
	if e, ok := a.(interface{ Equals(Token) bool }); ok {
		return e.Equals(b)
	}
	return fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b)
}
