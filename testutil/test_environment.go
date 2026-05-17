package testutil

import "strings"

// IsSingleByteNewLine returns true if the string contains only LF newlines (no CR characters).
func IsSingleByteNewLine(s string) bool {
	return strings.IndexRune(s, '\r') < 0
}
