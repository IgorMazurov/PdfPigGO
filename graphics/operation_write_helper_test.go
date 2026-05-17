package graphics

import (
	"bytes"
	"math"
	"strconv"
	"testing"
)

func TestWriteDouble(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{"Zero", 0, "0"},
		{"Five", 5, "5"},
		{"MinusFive", -5, "-5"},
		{"Ten", 10, "10"},
		{"MinusTen", -10, "-10"},
		{"SmallPositive", 0.00000001, "0.00000001"},
		{"SmallNegative", -0.00000001, "-0.00000001"},
		{"TinyTrailingZeros", 0.00000005100, "0.000000051"},
		{"TinyNegTrailingZeros", -0.0000000510, "-0.000000051"},
		{"RoundNumber", 10000.000, "10000"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteDouble(&buf, tc.value)
			if err != nil {
				t.Fatalf("WriteDouble(%g) returned error: %v", tc.value, err)
			}
			got := buf.String()
			if got != tc.expected {
				t.Errorf("WriteDouble(%g) = %q, want %q", tc.value, got, tc.expected)
			}
		})
	}
}

func TestWriteDecimalRoundTrip(t *testing.T) {
	value := 15001.98
	var buf bytes.Buffer
	err := WriteDouble(&buf, value)
	if err != nil {
		t.Fatalf("WriteDouble(%g) returned error: %v", value, err)
	}

	got, err := strconv.ParseFloat(buf.String(), 64)
	if err != nil {
		t.Fatalf("ParseFloat(%q) returned error: %v", buf.String(), err)
	}
	if got != value {
		t.Errorf("round-trip: parsed %g, want %g", got, value)
	}
}

func TestWriteDoubleExtremes(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{
			name:     "LargeNegative",
			value:    -340282346638528859811704183484516925440,
			expected: "-340282346638528859811704183484516925440",
		},
		{
			name:     "LargePositive",
			value:    340282346638528859811704183484516925440,
			expected: "340282346638528859811704183484516925440",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteDouble(&buf, tc.value)
			if err != nil {
				t.Fatalf("WriteDouble(%g) returned error: %v", tc.value, err)
			}
			got := buf.String()
			if got != tc.expected {
				t.Errorf("WriteDouble(%g) = %q, want %q", tc.value, got, tc.expected)
			}
		})
	}
}

func TestWriteDoubleSpecial(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"Infinity", math.Inf(1)},
		{"NegInfinity", math.Inf(-1)},
		{"NaN", math.NaN()},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteDouble(&buf, tc.value)
			if err != nil {
				t.Fatalf("WriteDouble(%g) returned error: %v", tc.value, err)
			}
			if buf.Len() == 0 {
				t.Errorf("WriteDouble(%g) produced empty output", tc.value)
			}
		})
	}
}
