package util

import (
	"errors"
	"strings"
)

var ErrNegativeOffset = errors.New("offset cannot be negative")

// StartsWithOffset reports whether value[offset:] begins with prefix.
// If offset is negative, it returns an error.
// If value is empty and prefix is also empty with offset 0, it returns true.
func StartsWithOffset(value, prefix string, offset int) (bool, error) {
	if offset < 0 {
		return false, ErrNegativeOffset
	}

	if value == "" {
		return len(prefix) == 0 && offset == 0, nil
	}

	if offset > len(value)-1 {
		return false, nil
	}

	return strings.HasPrefix(value[offset:], prefix), nil
}

// SubstringFrom returns the substring of text starting at position start.
// This mirrors C# string.AsSpan(start) / string.Substring(start).
func SubstringFrom(text string, start int) (string, error) {
	if start < 0 || start > len(text) {
		return "", errors.New("start index is out of range")
	}
	return text[start:], nil
}

// SubstringRange returns a substring of text starting at position start with the given length.
// This mirrors C# string.AsSpan(start, length) / string.Substring(start, length).
func SubstringRange(text string, start, length int) (string, error) {
	if start < 0 || start > len(text) {
		return "", errors.New("start index is out of range")
	}
	if length < 0 || start+length > len(text) {
		return "", errors.New("length is out of range or negative")
	}
	return text[start : start+length], nil
}
