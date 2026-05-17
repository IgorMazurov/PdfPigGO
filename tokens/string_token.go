package tokens

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"

	"github.com/uglytoad/pdfpig/go/core"
)

// StringEncoding represents the encoding used to convert the underlying file bytes to a string.
type StringEncoding int

const (
	// Iso88591 uses ISO 8859-1 (Latin-1) encoding.
	Iso88591 StringEncoding = iota
	// Utf16 uses UTF-16 little-endian encoding with BOM.
	Utf16
	// Utf16BE uses UTF-16 big-endian encoding with BOM prefix (0xFE 0xFF).
	Utf16BE
	// StringEncodingPDFDoc uses the PDF DocEncoding for strings in the body of a PDF file.
	StringEncodingPDFDoc
)

// StringToken represents a string of text contained in a PDF document.
type StringToken struct {
	data        string
	encodedWith StringEncoding
}

var _ Token = (*StringToken)(nil)
var _ DataToken[string] = (*StringToken)(nil)

// NewStringToken creates a new StringToken with the given data and encoding.
// If encodedWith is omitted, it defaults to Iso88591.
func NewStringToken(data string, encodedWith ...StringEncoding) *StringToken {
	enc := Iso88591
	if len(encodedWith) > 0 {
		enc = encodedWith[0]
	}
	return &StringToken{
		data:        data,
		encodedWith: enc,
	}
}

// Data returns the string data contained in the token.
func (t *StringToken) Data() string {
	return t.data
}

// EncodedWith returns the encoding used to generate the Data string from the file bytes.
func (t *StringToken) EncodedWith() StringEncoding {
	return t.encodedWith
}

// GetBytes converts the Data string back to bytes using the stored encoding.
func (t *StringToken) GetBytes() []byte {
	switch t.encodedWith {
	case Utf16BE:
		data := encodeUTF16BigEndian(t.data)
		result := make([]byte, len(data)+2)
		result[0] = 0xFE
		result[1] = 0xFF
		copy(result[2:], data)
		return result
	case Utf16:
		return encodeUTF16LE(t.data)
	case StringEncodingPDFDoc:
		return core.StringToBytes(t.data)
	default:
		return core.StringAsLatin1Bytes(t.data)
	}
}

// Equals reports whether other is a StringToken with the same Data and EncodedWith values.
func (t *StringToken) Equals(other Token) bool {
	o, ok := other.(*StringToken)
	if !ok {
		return false
	}
	return o.encodedWith == t.encodedWith && o.data == t.data
}

// String returns the string token formatted as a PDF literal string with parentheses.
func (t *StringToken) String() string {
	return fmt.Sprintf("(%s)", t.data)
}

func encodeUTF16BigEndian(s string) []byte {
	surrogates := utf16.Encode([]rune(s))
	result := make([]byte, len(surrogates)*2)
	for i, codeUnit := range surrogates {
		binary.BigEndian.PutUint16(result[i*2:], uint16(codeUnit))
	}
	return result
}

func encodeUTF16LE(s string) []byte {
	surrogates := utf16.Encode([]rune(s))
	result := make([]byte, 2+len(surrogates)*2)
	result[0] = 0xFF
	result[1] = 0xFE
	for i, codeUnit := range surrogates {
		binary.LittleEndian.PutUint16(result[2+i*2:], uint16(codeUnit))
	}
	return result
}
