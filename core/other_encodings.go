package core

import (
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// Iso88591 is the ISO 8859-1 (Latin-1) character encoding.
var Iso88591 = charmap.ISO8859_1

// StringAsLatin1Bytes converts a string to bytes using the ISO 8859-1 encoding.
func StringAsLatin1Bytes(s string) []byte {
	if s == "" {
		return nil
	}

	result, err := Iso88591.NewEncoder().String(s)
	if err != nil {
		return nil
	}

	return []byte(result)
}

// BytesAsLatin1String converts bytes to a string using the ISO 8859-1 encoding.
func BytesAsLatin1String(bytes []byte) string {
	result, _, err := transform.String(Iso88591.NewDecoder(), string(bytes))
	if err != nil {
		return ""
	}

	return result
}
