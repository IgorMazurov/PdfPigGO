package util

import (
	"testing"
	"time"
)

func TestTryParseDateTimeOffset_InvalidInput(t *testing.T) {
	cases := []string{
		"",
		"D:",
		"D:FEHTR$54",
		"D:49454",
		"9454AE",
		"20190107121634!",
		"D:19990209153925+11A",
		"D:19990209153925+11",
		"D:19990209153925E11",
		"D:19993209",
		"D:19990750",
		"D:20100231",
	}

	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			_, ok := TryParseDateTimeOffset(input)
			if ok {
				t.Errorf("expected false for input %q, got true", input)
			}
		})
	}

	t.Run("nil_string", func(t *testing.T) {
		var input string
		_, ok := TryParseDateTimeOffset(input)
		if ok {
			t.Error("expected false for nil string, got true")
		}
	})
}

func TestTryParseDateTimeOffset_ValidDate(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:  "full_date_with_positive_offset",
			input: "D:20190710205447+01'00'",
			expected: time.Date(2019, 7, 10, 20, 54, 47, 0, time.FixedZone("", int(time.Hour.Seconds()))),
		},
		{
			name:  "year_only_with_prefix",
			input: "D:2017",
			expected: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "year_only_no_prefix",
			input: "2017",
			expected: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "year_month_no_prefix",
			input: "196712",
			expected: time.Date(1967, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "year_month_with_prefix",
			input: "D:196712",
			expected: time.Date(1967, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "full_date_no_time",
			input: "D:20100520",
			expected: time.Date(2010, 5, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "full_date_no_prefix",
			input: "20121106",
			expected: time.Date(2012, 11, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "date_with_hour",
			input: "D:2012110623",
			expected: time.Date(2012, 11, 6, 23, 0, 0, 0, time.UTC),
		},
		{
			name:  "date_with_hour_minute",
			input: "D:201211061655",
			expected: time.Date(2012, 11, 6, 16, 55, 0, 0, time.UTC),
		},
		{
			name:  "date_with_full_time",
			input: "D:20121106005512",
			expected: time.Date(2012, 11, 6, 0, 55, 12, 0, time.UTC),
		},
		{
			name:  "date_with_zulu",
			input: "D:20121106165512Z",
			expected: time.Date(2012, 11, 6, 16, 55, 12, 0, time.UTC),
		},
		{
			name:  "date_with_zulu_no_prefix",
			input: "20121106165512Z",
			expected: time.Date(2012, 11, 6, 16, 55, 12, 0, time.UTC),
		},
		{
			name:  "negative_offset_hours_minutes",
			input: "D:19970915110347-07'30'",
			expected: time.Date(1997, 9, 15, 11, 3, 47, 0, time.FixedZone("", int((-7*time.Hour-30*time.Minute).Seconds()))),
		},
		{
			name:  "positive_hourly_offset",
			input: "D:19990209153925+11'",
			expected: time.Date(1999, 2, 9, 15, 39, 25, 0, time.FixedZone("", int((11*time.Hour).Seconds()))),
		},
		{
			name:  "negative_hourly_offset",
			input: "D:19990209153925-03'",
			expected: time.Date(1999, 2, 9, 15, 39, 25, 0, time.FixedZone("", int((-3*time.Hour).Seconds()))),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, ok := TryParseDateTimeOffset(tc.input)
			if !ok {
				t.Fatalf("expected true for input %q, got false", tc.input)
			}
			if !result.Equal(tc.expected) {
				t.Errorf("input %q:\n  expected %+v\n  got    %+v", tc.input, tc.expected, result)
			}
		})
	}
}
