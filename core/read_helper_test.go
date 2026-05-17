package core

import "testing"

func TestIsEndOfLineChar_LineFeed(t *testing.T) {
	if !IsEndOfLineChar('\n') {
		t.Error("expected line feed to be end of line")
	}
}

func TestIsEndOfLineChar_CarriageReturn(t *testing.T) {
	if !IsEndOfLineChar('\r') {
		t.Error("expected carriage return to be end of line")
	}
}

func TestIsEndOfLineChar_OtherCharacters(t *testing.T) {
	tests := []struct {
		name string
		ch   rune
	}{
		{"null", '\x00'},
		{"form feed", '\f'},
		{"backslash", '\\' },
		{"backspace", '\b'},
		{"letter a", 'a'},
		{"space", ' '},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if IsEndOfLineChar(tc.ch) {
				t.Errorf("expected %q to not be end of line", tc.ch)
			}
		})
	}
}
