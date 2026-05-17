package util

import (
	"testing"
)

func TestStartsWithOffset(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		prefix   string
		offset   int
		expected bool
	}{
		{
			name:     "single_char_at_zero_offset",
			input:    "a",
			prefix:   "a",
			offset:   0,
			expected: true,
		},
		{
			name:     "single_char_past_end_offset",
			input:    "a",
			prefix:   "a",
			offset:   1,
			expected: false,
		},
		{
			name:     "empty_input_non_empty_prefix",
			input:    "",
			prefix:   "abc",
			offset:   0,
			expected: false,
		},
		{
			name:     "offset_one_full_string_prefix",
			input:    "abc",
			prefix:   "abc",
			offset:   1,
			expected: false,
		},
		{
			name:     "large_offset",
			input:    "abc",
			prefix:   "abc",
			offset:   10,
			expected: false,
		},
		{
			name:     "prefix_not_at_start_offset_zero",
			input:    "pineapple",
			prefix:   "apple",
			offset:   0,
			expected: false,
		},
		{
			name:     "prefix_at_wrong_offset",
			input:    "pineapple",
			prefix:   "apple",
			offset:   3,
			expected: false,
		},
		{
			name:     "prefix_at_correct_offset",
			input:    "pineapple",
			prefix:   "apple",
			offset:   4,
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := StartsWithOffset(tc.input, tc.prefix, tc.offset)
			if err != nil {
				t.Fatalf("StartsWithOffset(%q, %q, %d) returned unexpected error: %v", tc.input, tc.prefix, tc.offset, err)
			}
			if result != tc.expected {
				t.Errorf("StartsWithOffset(%q, %q, %d) = %v, want %v", tc.input, tc.prefix, tc.offset, result, tc.expected)
			}
		})
	}
}

func TestStartsWithOffset_NegativeOffset(t *testing.T) {
	_, err := StartsWithOffset("any", "x", -1)
	if err == nil {
		t.Fatal("StartsWithOffset with negative offset expected error, got nil")
	}
	if err != ErrNegativeOffset {
		t.Errorf("StartsWithOffset with negative offset = %v, want ErrNegativeOffset", err)
	}
}
