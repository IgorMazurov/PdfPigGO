package type1

import (
	"fmt"
	"strconv"
	"strings"
)

// TokenType represents the kind of a Type 1 PostScript token.
type TokenType int

const (
	TokenTypeNone       TokenType = iota // None
	TokenTypeString                    // String
	TokenTypeName                          // Name
	TokenTypeLiteral                   // Literal
	TokenTypeReal                        // Real
	TokenTypeInteger                     // Integer
	// TokenTypeStartArray marks the beginning of an array ('[' or '{').
	TokenTypeStartArray
	// TokenTypeEndArray marks the end of an array (']' or '}').
	TokenTypeEndArray
	TokenTypeStartProc // StartProc
	TokenTypeEndProc   // EndProc
	TokenTypeStartDict // StartDict
	TokenTypeEndDict   // EndDict
	TokenTypeCharstring // Charstring
)

func (t TokenType) String() string {
	switch t {
	case TokenTypeNone:
		return "None"
	case TokenTypeString:
		return "String"
	case TokenTypeName:
		return "Name"
	case TokenTypeLiteral:
		return "Literal"
	case TokenTypeReal:
		return "Real"
	case TokenTypeInteger:
		return "Integer"
	case TokenTypeStartArray:
		return "StartArray"
	case TokenTypeEndArray:
		return "EndArray"
	case TokenTypeStartProc:
		return "StartProc"
	case TokenTypeEndProc:
		return "EndProc"
	case TokenTypeStartDict:
		return "StartDict"
	case TokenTypeEndDict:
		return "EndDict"
	case TokenTypeCharstring:
		return "Charstring"
	default:
		return fmt.Sprintf("TokenType(%d)", t)
	}
}

// Type1Token represents a parsed token from a Type 1 font program.
type Type1Token struct {
	Type TokenType
	Text string
}

// NewType1Token creates a new Type1Token with the given text and type.
func NewType1Token(text string, typ TokenType) *Type1Token {
	return &Type1Token{Text: text, Type: typ}
}

// NewType1TokenFromChar creates a new Type1Token from a single character.
func NewType1TokenFromChar(c rune, typ TokenType) *Type1Token {
	return NewType1Token(string(c), typ)
}

// IsPrivateDictionary returns true if this token is the "Private" literal keyword.
func (t *Type1Token) IsPrivateDictionary() bool {
	return t.Type == TokenTypeLiteral && strings.EqualFold(t.Text, "Private")
}

// AsInt parses the token text as an integer.
func (t *Type1Token) AsInt() int {
	v, _ := strconv.ParseFloat(t.Text, 64)
	return int(v)
}

// AsDouble parses the token text as a floating-point number using invariant culture.
func (t *Type1Token) AsDouble() float64 {
	v, _ := strconv.ParseFloat(t.Text, 64)
	return v
}

// AsBool returns true if the token text equals "true" (case-insensitive).
func (t *Type1Token) AsBool() bool {
	return strings.EqualFold(t.Text, "true")
}

// String returns a string representation of the token.
func (t *Type1Token) String() string {
	return fmt.Sprintf("Token[type=%s, text=%s]", t.Type, t.Text)
}

// Type1DataToken represents a Type 1 token that carries raw byte data (charstring).
type Type1DataToken struct {
	Type TokenType
	Data []byte
}

// errInvalidCharstringType is returned when a Type1DataToken is created with a non-Charstring type.
var errInvalidCharstringType = fmt.Errorf("invalid token type for type 1 data token receiving bytes, expected Charstring")

// NewType1DataToken creates a new Type1DataToken. The type must be TokenTypeCharstring.
func NewType1DataToken(typ TokenType, data []byte) (*Type1DataToken, error) {
	if typ != TokenTypeCharstring {
		return nil, fmt.Errorf("%w, got %s", errInvalidCharstringType, typ)
	}
	return &Type1DataToken{Type: typ, Data: data}, nil
}

// IsPrivateDictionary always returns false for a data token.
func (t *Type1DataToken) IsPrivateDictionary() bool {
	return false
}

// String returns a string representation of the data token.
func (t *Type1DataToken) String() string {
	return fmt.Sprintf("Token[type=%s, data=%d bytes]", t.Type, len(t.Data))
}
