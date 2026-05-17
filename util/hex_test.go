package util

import (
	"testing"
)

func TestGetString(t *testing.T) {
	cases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty_byte_slice",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "single_byte_37",
			input:    []byte{37},
			expected: "25",
		},
		{
			name:     "zero_byte",
			input:    []byte{0},
			expected: "00",
		},
		{
			name:     "max_byte",
			input:    []byte{255},
			expected: "FF",
		},
		{
			name:     "pdf_header_bytes",
			input:    []byte{37, 80, 68, 70, 45, 49, 46, 54, 13, 37},
			expected: "255044462D312E360D25",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetString(tc.input)
			if result != tc.expected {
				t.Errorf("GetString(%v) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}
