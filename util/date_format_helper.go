package util

import (
	"strconv"
	"time"
)

// TryParseDateTimeOffset attempts to parse a PDF-formatted date string into a time.Time.
// PDF date values follow the format D:YYYYMMDDHHmmSSOHH'mm as defined in ISO/IEC 8824,
// where the timezone offset O is one of Z, +, or -, optionally followed by HH'mm components.
func TryParseDateTimeOffset(s string) (time.Time, bool) {
	defer func() { recover() }()

	if s == "" || len(s) < 4 {
		return time.Time{}, false
	}

	location := 0
	if s[0] == 'D' && s[1] == ':' {
		location = 2
	}

	hasRemaining := func(pos, length int) bool {
		return pos+length <= len(s)
	}

	isAtEnd := func(pos int) bool {
		return pos == len(s)
	}

	inRange := func(val, min, max int) bool {
		return val >= min && val <= max
	}

	if !hasRemaining(location, 4) {
		return time.Time{}, false
	}

	year, err := strconv.Atoi(s[location : location+4])
	if err != nil {
		return time.Time{}, false
	}

	location += 4

	if !hasRemaining(location, 2) {
		if !isAtEnd(location) {
			return time.Time{}, false
		}
		return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC), true
	}

	month, err := strconv.Atoi(s[location : location+2])
	if err != nil || !inRange(month, 1, 12) {
		return time.Time{}, false
	}

	location += 2

	if !hasRemaining(location, 2) {
		if !isAtEnd(location) {
			return time.Time{}, false
		}
		return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), true
	}

	day, err := strconv.Atoi(s[location : location+2])
	if err != nil || !inRange(day, 1, 31) {
		return time.Time{}, false
	}

	location += 2

	// Validate day is valid for the given year/month (time.Date normalizes overflow silently).
	testDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if testDate.Year() != year || int(testDate.Month()) != month || testDate.Day() != day {
		return time.Time{}, false
	}

	if !hasRemaining(location, 2) {
		if !isAtEnd(location) {
			return time.Time{}, false
		}
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), true
	}

	hour, err := strconv.Atoi(s[location : location+2])
	if err != nil || !inRange(hour, 0, 23) {
		return time.Time{}, false
	}

	location += 2

	if !hasRemaining(location, 2) {
		if !isAtEnd(location) {
			return time.Time{}, false
		}
		return time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC), true
	}

	minute, err := strconv.Atoi(s[location : location+2])
	if err != nil || !inRange(minute, 0, 59) {
		return time.Time{}, false
	}

	location += 2

	if !hasRemaining(location, 2) {
		if !isAtEnd(location) {
			return time.Time{}, false
		}
		return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC), true
	}

	second, err := strconv.Atoi(s[location : location+2])
	if err != nil || !inRange(second, 0, 59) {
		return time.Time{}, false
	}

	location += 2

	if !hasRemaining(location, 1) {
		if !isAtEnd(location) {
			return time.Time{}, false
		}
		return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC), true
	}

	o := s[location]
	location++

	var sign int
	switch o {
	case '-':
		sign = -1
	case '+':
		sign = 1
	case 'Z':
		sign = 0
	default:
		return time.Time{}, false
	}

	if isAtEnd(location) {
		return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC), true
	}

	if !hasRemaining(location, 3) {
		return time.Time{}, false
	}

	hoursOffset, err := strconv.Atoi(s[location : location+2])
	if err != nil || s[location+2] != '\'' || !inRange(hoursOffset, 0, 23) {
		return time.Time{}, false
	}

	location += 3

	if isAtEnd(location) {
		offset := time.Duration(hoursOffset*sign) * time.Hour
		return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.FixedZone("", int(offset.Seconds()))), true
	}

	if !hasRemaining(location, 3) {
		return time.Time{}, false
	}

	minutesOffset, err := strconv.Atoi(s[location : location+2])
	if err != nil || s[location+2] != '\'' || !inRange(minutesOffset, 0, 59) {
		return time.Time{}, false
	}

	location += 3

	if isAtEnd(location) {
		offset := time.Duration(hoursOffset*sign)*time.Hour + time.Duration(minutesOffset*sign)*time.Minute
		return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.FixedZone("", int(offset.Seconds()))), true
	}

	return time.Time{}, false
}
