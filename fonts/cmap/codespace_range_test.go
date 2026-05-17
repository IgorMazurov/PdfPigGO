package cmap

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

func hexBytes(chars ...rune) []byte {
	token := tokens.NewHexToken(chars)
	result := make([]byte, len(token.Bytes()))
	copy(result, token.Bytes())
	return result
}

func TestCreatesCorrectly(t *testing.T) {
	tests := []struct {
		name       string
		startHex   []rune
		endHex     []rune
		startInt   int
		endInt     int
		codeLength int
	}{
		{
			name:       "single byte range 00-80",
			startHex:   []rune{'0', '0'},
			endHex:     []rune{'8', '0'},
			startInt:   0,
			endInt:     128,
			codeLength: 1,
		},
		{
			name:       "two byte range 8140-9ffc",
			startHex:   []rune{'8', '1', '4', '0'},
			endHex:     []rune{'9', 'f', 'f', 'c'},
			startInt:   33088,
			endInt:     40956,
			codeLength: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rng, err := NewCodespaceRange(hexBytes(tc.startHex...), hexBytes(tc.endHex...))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if rng.StartInt != tc.startInt {
				t.Errorf("expected StartInt %d, got %d", tc.startInt, rng.StartInt)
			}
			if rng.EndInt != tc.endInt {
				t.Errorf("expected EndInt %d, got %d", tc.endInt, rng.EndInt)
			}
			if rng.CodeLength != tc.codeLength {
				t.Errorf("expected CodeLength %d, got %d", tc.codeLength, rng.CodeLength)
			}
		})
	}
}

func TestMatchesNilCode(t *testing.T) {
	start := hexBytes('0', 'A')
	end := hexBytes('8', '0')

	rng, err := NewCodespaceRange(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rng.Matches(nil) {
		t.Error("expected Matches(nil) to be false")
	}
}

func TestIsFullMatchNilCode(t *testing.T) {
	start := hexBytes('0', 'A')
	end := hexBytes('8', '0')

	rng, err := NewCodespaceRange(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rng.IsFullMatch(nil, 2) {
		t.Error("expected IsFullMatch(nil, 2) to be false")
	}
}

func TestMatchesWrongLengthFalse(t *testing.T) {
	start := hexBytes('0', 'A')
	end := hexBytes('8', '0')

	rng, err := NewCodespaceRange(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	matches := rng.Matches(hexBytes('6', '9', '0', '1'))
	if matches {
		t.Error("expected Matches with wrong length code to be false")
	}
}

func TestMatchesLowerThanStartFalse(t *testing.T) {
	start := hexBytes('0', 'A')
	end := hexBytes('8', '0')

	rng, err := NewCodespaceRange(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	matches := rng.Matches(hexBytes('0', '1'))
	if matches {
		t.Error("expected Matches for code lower than start to be false")
	}
}

func TestMatchesHigherThanEndFalse(t *testing.T) {
	start := hexBytes('0', 'A')
	end := hexBytes('8', '0')

	rng, err := NewCodespaceRange(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	matches := rng.Matches(hexBytes('9', '6'))
	if matches {
		t.Error("expected Matches for code higher than end to be false")
	}
}

func TestMatchesInRangeTrue(t *testing.T) {
	start := hexBytes('0', 'A')
	end := hexBytes('8', '0')

	rng, err := NewCodespaceRange(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	matches := rng.Matches(hexBytes('5', 'A'))
	if !matches {
		t.Error("expected Matches for code in range to be true")
	}
}


