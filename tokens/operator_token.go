package tokens

import "fmt"

// OperatorToken represents an operator token encountered in a page content or Adobe Type 1 font stream.
type OperatorToken struct {
	data string
}

var _ Token = (*OperatorToken)(nil)
var _ DataToken[string] = (*OperatorToken)(nil)

// Begin text.
var Bt = &OperatorToken{data: "BT"}

// Def.
var Def = &OperatorToken{data: "def"}

// Dict.
var Dict = &OperatorToken{data: "dict"}

// Dup.
var Dup = &OperatorToken{data: "dup"}

// Eexec.
var Eexec = &OperatorToken{data: "eexec"}

// End object.
var EndObject = &OperatorToken{data: "endobj"}

// End stream.
var EndStream = &OperatorToken{data: "endstream"}

// End text.
var Et = &OperatorToken{data: "ET"}

// For.
var For = &OperatorToken{data: "for"}

// OpN is the "n" operator token (ends path construction).
var OpN = &OperatorToken{data: "n"}

// Put.
var Put = &OperatorToken{data: "put"}

// Pop (Q operator).
var QPop = &OperatorToken{data: "Q"}

// Push (q operator).
var QPush = &OperatorToken{data: "q"}

// OpR is the "R" operator token (object reference).
var OpR = &OperatorToken{data: "R"}

// Rectangle.
var Re = &OperatorToken{data: "re"}

// Readonly.
var Readonly = &OperatorToken{data: "readonly"}

// Object.
var StartObject = &OperatorToken{data: "obj"}

// Stream.
var StartStream = &OperatorToken{data: "stream"}

// Set font and size.
var Tf = &OperatorToken{data: "Tf"}

// Modify clipping.
var WStar = &OperatorToken{data: "W*"}

// OpXref is the "xref" operator token (cross reference).
var OpXref = &OperatorToken{data: "xref"}

// Cross reference section offset.
var StartXref = &OperatorToken{data: "startxref"}

// NewOperatorToken creates a new OperatorToken with the given data.
func NewOperatorToken(data string) *OperatorToken {
	return &OperatorToken{data: data}
}

// Data returns the operator name of this token.
func (t *OperatorToken) Data() string {
	return t.data
}

// CreateOperatorToken creates an OperatorToken from the given data, returning a pooled singleton for known operators.
func CreateOperatorToken(data string) *OperatorToken {
	switch data {
	case "BT":
		return Bt
	case "eexec":
		return Eexec
	case "endobj":
		return EndObject
	case "endstream":
		return EndStream
	case "ET":
		return Et
	case "def":
		return Def
	case "dict":
		return Dict
	case "for":
		return For
	case "dup":
		return Dup
	case "n":
		return OpN
	case "obj":
		return StartObject
	case "put":
		return Put
	case "Q":
		return QPop
	case "q":
		return QPush
	case "R":
		return OpR
	case "re":
		return Re
	case "readonly":
		return Readonly
	case "stream":
		return StartStream
	case "Tf":
		return Tf
	case "W*":
		return WStar
	case "xref":
		return OpXref
	case "startxref":
		return StartXref
	default:
		return NewOperatorToken(data)
	}
}

// Equals reports whether other is an OperatorToken with the same Data value.
func (t *OperatorToken) Equals(other Token) bool {
	o, ok := other.(*OperatorToken)
	if !ok {
		return false
	}
	return o.data == t.data
}

// String returns the string representation of the operator token.
func (t *OperatorToken) String() string {
	return fmt.Sprintf("%s", t.data)
}
