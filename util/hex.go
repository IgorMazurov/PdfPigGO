package util

import (
	"encoding/hex"
	"strings"
)

// hexChars holds the uppercase hexadecimal character lookup table.
var hexChars = []byte("0123456789ABCDEF")

// GetUtf8Chars writes the two-character hex representation of each byte in bytes
// into utf8Chars, which must have capacity for len(bytes)*2 elements.
func GetUtf8Chars(bytes, utf8Chars []byte) {
	position := 0
	for _, b := range bytes {
		utf8Chars[position] = hexChars[getHighNibble(b)]
		position++
		utf8Chars[position] = hexChars[getLowNibble(b)]
		position++
	}
}

// GetString returns the uppercase hexadecimal string representation of the given bytes.
func GetString(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

// getHighNibble extracts the upper 4 bits of a byte as an index (0-15).
func getHighNibble(b byte) int {
	return int((b & 0xF0) >> 4)
}

// getLowNibble extracts the lower 4 bits of a byte as an index (0-15).
func getLowNibble(b byte) int {
	return int(b & 0x0F)
}
