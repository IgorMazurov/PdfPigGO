package type1

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
)

// Type1Tokenizer tokenizes the binary portion of a Type 1 font program.
type Type1Tokenizer struct {
	bytes      core.InputBytes
	comments   []string
	openParens int
	prevToken  TokenLike
	current    TokenLike
}

// TokenLike is the common interface for Type1Token and Type1DataToken.
type TokenLike interface {
	TokenType() TokenType
	IsPrivateDictionary() bool
	String() string
}

// NewType1Tokenizer creates a new tokenizer for the given input bytes.
func NewType1Tokenizer(bytes core.InputBytes) (*Type1Tokenizer, error) {
	if bytes == nil {
		return nil, fmt.Errorf("bytes cannot be nil")
	}
	t := &Type1Tokenizer{
		bytes:    bytes,
		comments: make([]string, 0),
	}
	t.current = t.readNextToken()
	return t, nil
}

// CurrentToken returns the current token. May be nil at end of input.
func (t *Type1Tokenizer) CurrentToken() TokenLike {
	return t.current
}

// Comments returns the list of comments encountered during tokenization.
func (t *Type1Tokenizer) Comments() []string {
	return t.comments
}

// GetNext advances to and returns the next token. Returns nil at end of input.
func (t *Type1Tokenizer) GetNext() TokenLike {
	t.current = t.readNextToken()
	return t.current
}

func (t *Type1Tokenizer) readNextToken() TokenLike {
	t.prevToken = t.current
	var skip bool

	for {
		skip = false
		for t.bytes.MoveNext() {
			b := t.bytes.CurrentByte()
			c := rune(b)

			switch c {
			case '%':
				t.comments = append(t.comments, t.readComment())

			case '(':
				return t.readString()

			case ')':
				return NewType1TokenFromChar(')', TokenTypeNone)

			case '[':
				return NewType1TokenFromChar('[', TokenTypeStartArray)

			case ']':
				return NewType1TokenFromChar(']', TokenTypeEndArray)

			case '{':
				return NewType1TokenFromChar('{', TokenTypeStartProc)

			case '}':
				return NewType1TokenFromChar('}', TokenTypeEndProc)

			case '/':
				name := t.readLiteral(nil)
				if name == "" {
					return nil
				}
				return NewType1Token(name, TokenTypeLiteral)

			case '<':
				if following, ok := t.bytes.Peek(); ok && following == '<' {
					t.bytes.MoveNext()
					return NewType1Token("<<", TokenTypeStartDict)
				}
				return NewType1TokenFromChar('<', TokenTypeNone)

			case '>':
				if following, ok := t.bytes.Peek(); ok && following == '>' {
					t.bytes.MoveNext()
					return NewType1Token(">>", TokenTypeEndDict)
				}
				return NewType1TokenFromChar('>', TokenTypeNone)

			default:
				if core.IsWhitespace(b) || b == 0 {
					skip = true
					break
				}

				if number := t.tryReadNumber(c); number != nil {
					return number
				}

				name := t.readLiteral(&c)
				if name == "" {
					return NewType1TokenFromChar(c, TokenTypeNone)
				}

				if strings.EqualFold(name, RdProcedure) || name == RdProcedureAlt {
					prevInt := false
					if pt, ok := t.prevToken.(*Type1Token); ok && pt.Type == TokenTypeInteger {
						prevInt = true
					}
					if prevInt {
						return t.readCharString(t.prevTokenAsInt())
					}
					return NewType1TokenFromChar(c, TokenTypeNone)
				}

				return NewType1Token(name, TokenTypeNone)
			}
		}
		if !skip {
			break
		}
	}

	return nil
}

func (t *Type1Tokenizer) readString() TokenLike {
	var sb strings.Builder

	for t.bytes.MoveNext() {
		c := rune(t.bytes.CurrentByte())

		switch c {
		case '(':
			t.openParens++
			sb.WriteRune(c)

		case ')':
			if t.openParens == 0 {
				return NewType1Token(sb.String(), TokenTypeString)
			}
			sb.WriteRune(c)
			t.openParens--

		case '\\':
			t.bytes.MoveNext()
			c1 := rune(t.bytes.CurrentByte())
			switch c1 {
			case 'n', '\n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case 'b':
				sb.WriteByte('\b')
			case 'f':
				sb.WriteByte('\f')
			case '\\', '(', ')':
				sb.WriteRune(c1)
			default:
				if c1 >= '0' && c1 <= '9' {
					octal := string([]rune{c1, t.nextChar(), t.nextChar()})
					code, err := strconv.ParseInt(octal, 8, 32)
					if err == nil {
						sb.WriteRune(rune(code))
					} else {
						sb.WriteString(octal)
					}
				} else {
					sb.WriteRune(c1)
				}
			}

		case '\r', '\n':
			sb.WriteByte('\n')

		default:
			sb.WriteRune(c)
		}
	}

	return nil
}

func (t *Type1Tokenizer) nextChar() rune {
	t.bytes.MoveNext()
	return rune(t.bytes.CurrentByte())
}

func (t *Type1Tokenizer) tryReadNumber(c rune) TokenLike {
	currentPos := t.bytes.CurrentOffset()
	var sb strings.Builder
	var radixStr strings.Builder
	hasDigit := false

	if c == '+' || c == '-' {
		sb.WriteRune(c)
		c = t.nextChar()
	}

	for c >= '0' && c <= '9' {
		sb.WriteRune(c)
		c = t.nextChar()
		hasDigit = true
	}

	if c == '.' {
		sb.WriteRune(c)
		c = t.nextChar()
	} else if c == '#' {
		radixStr.WriteString(sb.String())
		sb.Reset()
		c = t.nextChar()
	} else if sb.Len() == 0 || !hasDigit {
		t.bytes.Seek(currentPos, 0) // io.SeekStart
		return nil
	} else {
		t.bytes.Seek(t.bytes.CurrentOffset()-1, 0) // io.SeekStart
		return NewType1Token(sb.String(), TokenTypeInteger)
	}

	if c >= '0' && c <= '9' {
		sb.WriteRune(c)
		c = t.nextChar()
	} else {
		t.bytes.Seek(currentPos, 0) // io.SeekStart
		return nil
	}

	for c >= '0' && c <= '9' {
		sb.WriteRune(c)
		c = t.nextChar()
	}

	if c == 'E' {
		sb.WriteRune(c)
		c = t.nextChar()

		if c == '-' {
			sb.WriteRune(c)
			c = t.nextChar()
		}

		if c >= '0' && c <= '9' {
			sb.WriteRune(c)
			c = t.nextChar()
		} else {
			t.bytes.Seek(currentPos, 0) // io.SeekStart
			return nil
		}

		for c >= '0' && c <= '9' {
			sb.WriteRune(c)
			c = t.nextChar()
		}
	}

	t.bytes.Seek(t.bytes.CurrentOffset()-1, 0) // io.SeekStart

	if radixStr.Len() > 0 {
		radix, err := strconv.ParseInt(radixStr.String(), 10, 32)
		if err == nil {
			numStr := sb.String()
			num, err2 := strconv.ParseInt(numStr, int(radix), 64)
			if err2 == nil {
				return NewType1Token(strconv.FormatInt(num, 10), TokenTypeInteger)
			}
		}
		return NewType1Token(sb.String(), TokenTypeReal)
	}

	return NewType1Token(sb.String(), TokenTypeReal)
}

func (t *Type1Tokenizer) readLiteral(prevChar *rune) string {
	var sb strings.Builder
	if prevChar != nil {
		sb.WriteRune(*prevChar)
	}

	for {
		b, ok := t.bytes.Peek()
		if !ok {
			break
		}

		c := rune(b)
		if core.IsWhitespace(byte(c)) || c == '(' || c == ')' || c == '<' || c == '>' ||
			c == '[' || c == ']' || c == '{' || c == '}' || c == '/' || c == '%' {
			break
		}

		sb.WriteRune(c)
		t.bytes.MoveNext()
	}

	return sb.String()
}

func (t *Type1Tokenizer) readComment() string {
	var sb strings.Builder

	for t.bytes.MoveNext() {
		c := rune(t.bytes.CurrentByte())
		if core.IsEndOfLineChar(c) {
			break
		}
		sb.WriteRune(c)
	}

	return sb.String()
}

func (t *Type1Tokenizer) readCharString(length int) TokenLike {
	t.bytes.MoveNext()

	data := make([]byte, length)
	for i := 0; i < length; i++ {
		t.bytes.MoveNext()
		data[i] = t.bytes.CurrentByte()
	}

	token, _ := NewType1DataToken(TokenTypeCharstring, data)
	return token
}

func (t *Type1Tokenizer) prevTokenAsInt() int {
	if pt, ok := t.prevToken.(*Type1Token); ok {
		return pt.AsInt()
	}
	return 0
}

// TokenType returns the token type. Implements TokenLike.
func (t *Type1Token) TokenType() TokenType {
	return t.Type
}

// TokenType returns the token type. Implements TokenLike.
func (t *Type1DataToken) TokenType() TokenType {
	return t.Type
}
