package tokenization

import (
	"math"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// NumericTokenizer tokenizes PDF numeric literals including integers,
// fractional numbers, and scientific notation.
type NumericTokenizer struct{}

var _ Tokenizer = (*NumericTokenizer)(nil)

const (
	byteZero          byte = 48
	byteNine          byte = 57
	byteNegative      byte = '-'
	bytePositive      byte = '+'
	bytePeriod        byte = '.'
	byteExponentLower byte = 'e'
	byteExponentUpper byte = 'E'
)

// NewNumericTokenizer creates a new NumericTokenizer.
func NewNumericTokenizer() *NumericTokenizer {
	return &NumericTokenizer{}
}

// ReadsNextByte returns true because this tokenizer reads ahead past the number.
func (t *NumericTokenizer) ReadsNextByte() bool {
	return true
}

// Tokenize attempts to tokenize a PDF numeric literal.
// Returns the NumericToken and true on success, or nil and false if parsing fails.
func (t *NumericTokenizer) Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool) {
	readBytes := 0

	isNegative := false
	integerPart := 0.0

	hasFraction := false
	fractionalPart := int64(0)
	fractionalCount := 0

	hasExponent := false
	isExponentNegative := false
	exponentPart := 0

	for {
		b := input.CurrentByte()

		if b >= byteZero && b <= byteNine {
			digit := int(b - byteZero)
			if hasExponent {
				exponentPart = exponentPart*10 + digit
			} else if hasFraction {
				fractionalPart = fractionalPart*10 + int64(digit)
				fractionalCount++
			} else {
				integerPart = integerPart*10 + float64(digit)
			}
		} else if b == bytePositive {
			// No impact
		} else if b == byteNegative {
			if hasExponent {
				isExponentNegative = true
			} else {
				isNegative = true
			}
		} else if b == bytePeriod {
			if hasExponent || hasFraction {
				return nil, false
			}
			hasFraction = true
		} else if b == byteExponentLower || b == byteExponentUpper {
			if readBytes == 0 {
				return nil, false
			}
			if hasExponent {
				return nil, false
			}
			hasExponent = true
		} else {
			if readBytes == 0 {
				return nil, false
			}
			break
		}

		readBytes++
		if !input.MoveNext() {
			break
		}
	}

	if hasExponent && !isExponentNegative {
		combined := integerPart*pow10(fractionalCount) + float64(fractionalPart)

		shift := exponentPart - fractionalCount

		if shift >= 0 {
			integerPart = combined * pow10(shift)
		} else {
			integerPart = combined / pow10(-shift)
		}

		hasFraction = false
		hasExponent = false
	}

	if hasFraction && fractionalCount > 0 {
		switch fractionalCount {
		case 1:
			integerPart += float64(fractionalPart) / 10.0
		case 2:
			integerPart += float64(fractionalPart) / 100.0
		case 3:
			integerPart += float64(fractionalPart) / 1000.0
		default:
			integerPart += float64(fractionalPart) / math.Pow(10, float64(fractionalCount))
		}
	}

	if hasExponent {
		signedExponent := exponentPart
		if isExponentNegative {
			signedExponent = -signedExponent
		}
		integerPart *= math.Pow(10, float64(signedExponent))
	}

	if isNegative {
		integerPart = -integerPart
	}

	if integerPart == 0 {
		return tokens.Zero, true
	}

	return tokens.NewNumericToken(integerPart), true
}

func pow10(exp int) float64 {
	switch exp {
	case 0:
		return 1
	case 1:
		return 10
	case 2:
		return 100
	case 3:
		return 1000
	case 4:
		return 10000
	case 5:
		return 100000
	case 6:
		return 1000000
	case 7:
		return 10000000
	case 8:
		return 100000000
	case 9:
		return 1000000000
	default:
		return math.Pow(10, float64(exp))
	}
}
