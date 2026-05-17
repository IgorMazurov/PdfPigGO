package tokens

import (
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf16"
)

// HexToken is a token containing string data where the string is encoded as hexadecimal.
type HexToken struct {
	data  string
	bytes []byte
}

var _ Token = (*HexToken)(nil)
var _ DataToken[string] = (*HexToken)(nil)

var hexMap = func() map[rune]byte {
	m := make(map[rune]byte, 32)
	for i := byte(0); i <= 9; i++ {
		m[rune('0'+i)] = i
	}
	pairs := [][2]rune{{'A', 'a'}, {'B', 'b'}, {'C', 'c'}, {'D', 'd'}, {'E', 'e'}, {'F', 'f'}}
	for i, p := range pairs {
		m[p[0]] = byte(10 + i)
		m[p[1]] = byte(10 + i)
	}
	return m
}()

// NewHexToken creates a new HexToken from the provided hex characters.
// If the character count is odd, the final character is treated as having a trailing '0' per PDF spec 7.3.4.3.
func NewHexToken(characters []rune) *HexToken {
	bytes := make([]byte, (len(characters)+1)/2)
	index := 0

	for i := 0; i < len(characters); i += 2 {
		high := characters[i]
		var low rune
		if i == len(characters)-1 {
			low = '0'
		} else {
			low = characters[i+1]
		}

		bytes[index] = convertPair(high, low)
		index++
	}

	data := decodeHexBytes(bytes)

	return &HexToken{
		data:  data,
		bytes: bytes,
	}
}

// NewHexTokenFromRaw creates a new HexToken from raw bytes.
func NewHexTokenFromRaw(raw []byte) *HexToken {
	copied := make([]byte, len(raw))
	copy(copied, raw)
	return &HexToken{
		data:  string(copied),
		bytes: copied,
	}
}

// Data returns the string contained in the hex data.
func (t *HexToken) Data() string {
	return t.data
}

// Bytes returns the raw byte slice of the hex data.
func (t *HexToken) Bytes() []byte {
	return t.bytes
}

// GetHexString converts the binary data back to a hex string representation.
func (t *HexToken) GetHexString() string {
	return strings.ToUpper(hex.EncodeToString(t.bytes))
}

// Equals reports whether other is a HexToken with identical Data.
func (t *HexToken) Equals(other Token) bool {
	o, ok := other.(*HexToken)
	if !ok {
		return false
	}
	return t.data == o.data
}

// String returns the hex token formatted as a PDF hexadecimal string with angle brackets.
func (t *HexToken) String() string {
	return fmt.Sprintf("<%s>", t.GetHexString())
}

// ConvertPair converts two hex characters to a single byte.
func ConvertPair(high, low rune) byte {
	return convertPair(high, low)
}

func convertPair(high, low rune) byte {
	highByte := hexMap[high]
	lowByte := hexMap[low]
	return (highByte << 4) | lowByte
}

// ConvertHexBytesToInt converts the bytes in a HexToken to an integer.
// Only the first one or two bytes are used.
func ConvertHexBytesToInt(token *HexToken) int {
	b := token.bytes
	value := int(b[0])
	if len(b) >= 2 {
		value <<= 8
		value += int(b[1])
	}
	return value
}

// decodeHexBytes interprets the byte slice as either UTF-16BE (with FE FF BOM) or plain Latin-1,
// skipping null bytes for the latter.
func decodeHexBytes(bytes []byte) string {
	if len(bytes) >= 2 && bytes[0] == 0xFE && bytes[1] == 0xFF {
		return decodeUTF16BE(bytes[2:])
	}

	var builder strings.Builder
	for _, b := range bytes {
		if b != 0 {
			builder.WriteRune(rune(b))
		}
	}
	return builder.String()
}

// decodeUTF16BE decodes a UTF-16 big-endian byte slice (without BOM) to a Go string.
func decodeUTF16BE(data []byte) string {
	if len(data)%2 != 0 {
		data = append(data, 0)
	}

	codeUnits := make([]uint16, len(data)/2)
	for i := 0; i < len(codeUnits); i++ {
		codeUnits[i] = uint16(data[i*2])<<8 | uint16(data[i*2+1])
	}

	runes := utf16.Decode(codeUnits)
	return string(runes)
}
