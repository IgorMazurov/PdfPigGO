package util

import (
	"bytes"
	"testing"
)

func TestCircularByteBufferCanExceedCapacity(t *testing.T) {
	buffer := NewCircularByteBuffer(3)

	input := []byte("123456")
	for i := 0; i < len(input); i++ {
		buffer.Add(input[i])
	}

	if !buffer.IsCurrentlyEqual("456") {
		t.Errorf("IsCurrentlyEqual(\"456\") = false, want true")
	}

	if !bytes.Equal([]byte("456"), buffer.AsBytes()) {
		t.Errorf("AsBytes() = %v, want %v", buffer.AsBytes(), []byte("456"))
	}

	tests := []struct {
		suffix string
		want   bool
	}{
		{"6", true},
		{"56", true},
		{"456", true},
		{"3456", false},
	}

	for _, tc := range tests {
		if got := buffer.EndsWith(tc.suffix); got != tc.want {
			t.Errorf("EndsWith(%q) = %v, want %v", tc.suffix, got, tc.want)
		}
	}
}

func TestCircularByteBufferCanUndershootCapacity(t *testing.T) {
	buffer := NewCircularByteBuffer(9)

	input := []byte("123456")
	for i := 0; i < len(input); i++ {
		buffer.Add(input[i])
	}

	if !buffer.IsCurrentlyEqual("123456") {
		t.Errorf("IsCurrentlyEqual(\"123456\") = false, want true")
	}

	tests := []struct {
		suffix string
		want   bool
	}{
		{"3456", true},
		{"123", false},
	}

	for _, tc := range tests {
		if got := buffer.EndsWith(tc.suffix); got != tc.want {
			t.Errorf("EndsWith(%q) = %v, want %v", tc.suffix, got, tc.want)
		}
	}

	if !bytes.Equal([]byte("123456"), buffer.AsBytes()) {
		t.Errorf("AsBytes() = %v, want %v", buffer.AsBytes(), []byte("123456"))
	}
}

func TestCircularByteBufferCanAddReverse(t *testing.T) {
	bufferLen := len("startxref")

	s := "wibbly bibble startxref 2024"

	buffer := NewCircularByteBuffer(bufferLen)

	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		buffer.AddReverse(c)

		if i <= len(s)-bufferLen {
			str := s[i : i+bufferLen]

			if !buffer.IsCurrentlyEqual(str) {
				t.Errorf("IsCurrentlyEqual(%q) = false, want true at index %d", str, i)
			}
		}
	}
}
