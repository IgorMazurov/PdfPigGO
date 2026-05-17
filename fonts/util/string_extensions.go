package util

import "errors"

var (
	ErrStartOutOfRange   = errors.New("start index is out of range")
	ErrLengthOutOfRange  = errors.New("length is out of range or negative")
)

// SubstringFrom returns the substring of text starting at position start.
// This mirrors C# string.Substring(start) / ReadOnlySpan<char>.Slice(start).
func SubstringFrom(text string, start int) (string, error) {
	if start < 0 || start > len(text) {
		return "", ErrStartOutOfRange
	}
	return text[start:], nil
}

// SubstringRange returns a substring of text starting at position start with the given length.
// This mirrors C# string.Substring(start, length) / ReadOnlySpan<char>.Slice(start, length).
func SubstringRange(text string, start, length int) (string, error) {
	if start < 0 || start > len(text) {
		return "", ErrStartOutOfRange
	}
	if length < 0 || start+length > len(text) {
		return "", ErrLengthOutOfRange
	}
	return text[start : start+length], nil
}
