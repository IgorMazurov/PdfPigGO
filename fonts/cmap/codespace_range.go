package cmap

import (
	"errors"
	"fmt"
)

// CodespaceRange represents a codespace range specified by a pair of codes
// of some particular length giving the lower and upper bounds of that range.
type CodespaceRange struct {
	Start      []byte
	End        []byte
	StartInt   int
	EndInt     int
	CodeLength int
}

// NewCodespaceRange creates a new CodespaceRange from start and end byte slices.
func NewCodespaceRange(start, end []byte) (*CodespaceRange, error) {
	if len(start) == 0 || len(end) == 0 {
		return nil, errors.New("start and end must not be empty")
	}

	if len(start) != len(end) {
		return nil, fmt.Errorf("start and end must have the same length: start has %d bytes, end has %d bytes", len(start), len(end))
	}

	return &CodespaceRange{
		Start:      start,
		End:        end,
		StartInt:   ToInt(start),
		EndInt:     ToInt(end),
		CodeLength: len(start),
	}, nil
}

// Matches returns true if the given code bytes match this codespace range.
func (r CodespaceRange) Matches(code []byte) bool {
	if code == nil {
		return false
	}

	return r.IsFullMatch(code, len(code))
}

// IsFullMatch returns true if the given code bytes match this codespace range.
func (r CodespaceRange) IsFullMatch(code []byte, codeLength int) bool {
	if code == nil {
		return false
	}

	// The code must be the same length as the bounding codes.
	if codeLength != r.CodeLength {
		return false
	}

	value := ToInt(code[:codeLength])
	return value >= r.StartInt && value <= r.EndInt
}

// String returns a string representation of the range.
func (r CodespaceRange) String() string {
	return fmt.Sprintf("Length %d: %d -> %d", r.CodeLength, r.StartInt, r.EndInt)
}
