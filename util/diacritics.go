package util

import (
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// nonCombiningDiacritics is a set of characters that are not combining diacritics
// but can appear as standalone diacritical marks.
var nonCombiningDiacritics = map[string]bool{
	"´": true,
	"^": true,
	"ˆ": true,
	"¨": true,
	"©": true,
	"™": true,
	"®": true,
	"`": true,
	"˜": true,
	"∼": true,
	"¸": true,
}

// isPotentialStandaloneDiacritic reports whether value is a known non-combining
// diacritical mark that could appear as a standalone character in PDF text.
func isPotentialStandaloneDiacritic(value string) bool {
	return nonCombiningDiacritics[value]
}

// IsInCombiningDiacriticRange reports whether the single-character string value
// falls within the Unicode combining diacritic range (U+0300 to U+036F).
// It returns false if value does not contain exactly one character.
func IsInCombiningDiacriticRange(value string) bool {
	r, size := utf8.DecodeRuneInString(value)
	if size != len(value) {
		return false
	}

	code := int(r)
	return code >= 768 && code <= 879
}

// TryCombineDiacriticWithPreviousLetter attempts to combine a diacritic character
// with the preceding letter to form a single combined grapheme cluster.
// It returns the combined string and true if successful (grapheme count unchanged),
// or an empty string and false otherwise.
func TryCombineDiacriticWithPreviousLetter(diacritic, previous string) (string, bool) {
	if previous == "" {
		return "", false
	}

	result := previous + diacritic

	beforeLen := measureDiacriticAwareLength(previous)
	afterLen := measureDiacriticAwareLength(result)

	if beforeLen == afterLen {
		return result, true
	}

	return "", false
}

// measureDiacriticAwareLength counts the number of grapheme clusters in input,
// matching the behavior of C#'s StringInfo.GetTextElementEnumerator.
func measureDiacriticAwareLength(input string) int {
	nfd := norm.NFD.String(input)
	count := 0
	for _, r := range nfd {
		if !unicode.Is(unicode.Mn, r) && !unicode.Is(unicode.Me, r) && !unicode.Is(unicode.Mc, r) {
			count++
		}
	}
	return count
}


