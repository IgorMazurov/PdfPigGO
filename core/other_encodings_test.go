package core

import (
	"testing"
)

func TestBytesAsLatin1StringEmptyReturnsEmpty(t *testing.T) {
	result := BytesAsLatin1String([]byte{})
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestStringAsLatin1BytesNullReturnsNil(t *testing.T) {
	result := StringAsLatin1Bytes("")
	if result != nil {
		t.Errorf("expected nil bytes, got %v", result)
	}
}
