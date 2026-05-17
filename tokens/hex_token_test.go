package tokens

import (
	"testing"
)

func TestHexTokenMapsCorrectlyToString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"AE", "AE", "®"},
		{"61", "61", "a"},
		{"0061", "0061", "a"},
		{"7465787420736f", "7465787420736f", "text so"},
		{"6170", "6170", "ap"},
		{"617 odd length", "617", "ap"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token := NewHexToken([]rune(tc.input))

			if got := token.Data(); got != tc.expected {
				t.Errorf("expected Data %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestHexTokenMapsCorrectlyToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"0003", "0003", 3},
		{"0011", "0011", 17},
		{"0024", "0024", 36},
		{"0037", "0037", 55},
		{"0044", "0044", 68},
		{"005B", "005B", 91},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token := NewHexToken([]rune(tc.input))

			got := ConvertHexBytesToInt(token)
			if got != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}
