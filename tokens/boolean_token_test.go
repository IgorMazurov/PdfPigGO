package tokens

import "testing"

func TestBooleanTokenObjectEquals(t *testing.T) {
	one := True
	two := interface{}(True)

	if !one.Equals(two.(Token)) {
		t.Error("expected BooleanToken.True to equal itself")
	}
}

func TestBooleanTokenObjectNotEqual(t *testing.T) {
	one := False
	two := interface{}(True)

	if one.Equals(two.(Token)) {
		t.Error("expected BooleanToken.False not to equal BooleanToken.True")
	}
}

func TestBooleanTokenHashCodeMatch(t *testing.T) {
	if True.HashCode() != True.HashCode() {
		t.Error("expected hash codes of True to match themselves")
	}
}

func TestBooleanTokenHashCodeNotMatch(t *testing.T) {
	if True.HashCode() == False.HashCode() {
		t.Error("expected hash codes of True and False to not match")
	}
}

func TestBooleanTokenStringRepresentationCorrect(t *testing.T) {
	tests := []struct {
		name     string
		token    *BooleanToken
		expected string
	}{
		{"True", True, "True"},
		{"False", False, "False"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.token.String(); got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}
