package core

import "fmt"

// CharacterToShort converts an ASCII digit character to its integer value.
// Supports '0' through '9'. Returns an error for non-digit characters.
func CharacterToShort(c rune) (int, error) {
	if c >= '0' && c <= '9' {
		return int(c - '0'), nil
	}
	return 0, fmt.Errorf("could not convert the character %c to a short", c)
}

// FromOctalDigits converts a slice of octal digit values (0-7) to their integer representation.
// Each element in the octal slice must be in range [0, 7].
func FromOctalDigits(octal []int) int {
	sum := 0
	for i := len(octal) - 1; i >= 0; i-- {
		sum += octal[i] * quickPower(8, i)
	}
	return sum
}

// quickPower computes x raised to the power of pow using exponentiation by squaring.
func quickPower(x, pow int) int {
	ret := 1
	for pow != 0 {
		if (pow & 1) == 1 {
			ret *= x
		}
		x *= x
		pow >>= 1
	}
	return ret
}
