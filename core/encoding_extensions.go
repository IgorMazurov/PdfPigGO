package core

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// GetString converts a byte slice to a string using the provided encoding.
// This is a polyfill equivalent for C#'s Encoding.GetString(ReadOnlySpan<byte>).
// If decoding fails, an empty string is returned.
func GetString(enc encoding.Encoding, bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	}

	result, _, err := transform.String(enc.NewDecoder(), string(bytes))
	if err != nil {
		return ""
	}

	return result
}
