package tokens

import (
	"fmt"
	"math"
)

// NumericToken represents a numeric value token in a PDF document.
// PDF supports integer and real numbers; this token represents both types
// as they are used interchangeably in the specification.
type NumericToken struct {
	data float64
}

var _ Token = (*NumericToken)(nil)
var _ DataToken[float64] = (*NumericToken)(nil)

// MinusOne is a singleton instance of numeric token for -1.
var MinusOne = &NumericToken{data: -1}

// Zero is a singleton instance of numeric token for 0.
var Zero = &NumericToken{data: 0}

// One is a singleton instance of numeric token for 1.
var One = &NumericToken{data: 1}

// Two is a singleton instance of numeric token for 2.
var Two = &NumericToken{data: 2}

// Three is a singleton instance of numeric token for 3.
var Three = &NumericToken{data: 3}

// Four is a singleton instance of numeric token for 4.
var Four = &NumericToken{data: 4}

// Five is a singleton instance of numeric token for 5.
var Five = &NumericToken{data: 5}

// Six is a singleton instance of numeric token for 6.
var Six = &NumericToken{data: 6}

// Seven is a singleton instance of numeric token for 7.
var Seven = &NumericToken{data: 7}

// Eight is a singleton instance of numeric token for 8.
var Eight = &NumericToken{data: 8}

// Nine is a singleton instance of numeric token for 9.
var Nine = &NumericToken{data: 9}

// Ten is a singleton instance of numeric token for 10.
var Ten = &NumericToken{data: 10}

// Eleven is a singleton instance of numeric token for 11.
var Eleven = &NumericToken{data: 11}

// Twelve is a singleton instance of numeric token for 12.
var Twelve = &NumericToken{data: 12}

// Thirteen is a singleton instance of numeric token for 13.
var Thirteen = &NumericToken{data: 13}

// Fourteen is a singleton instance of numeric token for 14.
var Fourteen = &NumericToken{data: 14}

// Fifteen is a singleton instance of numeric token for 15.
var Fifteen = &NumericToken{data: 15}

// Sixteen is a singleton instance of numeric token for 16.
var Sixteen = &NumericToken{data: 16}

// Seventeen is a singleton instance of numeric token for 17.
var Seventeen = &NumericToken{data: 17}

// Eighteen is a singleton instance of numeric token for 18.
var Eighteen = &NumericToken{data: 18}

// Nineteen is a singleton instance of numeric token for 19.
var Nineteen = &NumericToken{data: 19}

// Twenty is a singleton instance of numeric token for 20.
var Twenty = &NumericToken{data: 20}

// OneHundred is a singleton instance of numeric token for 100.
var OneHundred = &NumericToken{data: 100}

// FiveHundred is a singleton instance of numeric token for 500.
var FiveHundred = &NumericToken{data: 500}

// OneThousand is a singleton instance of numeric token for 1000.
var OneThousand = &NumericToken{data: 1000}

// NewNumericToken creates a new NumericToken with the given float64 value.
func NewNumericToken(data float64) *NumericToken {
	return &NumericToken{data: data}
}

// NewNumericTokenFromInt creates a new NumericToken from an integer value.
func NewNumericTokenFromInt(value int) *NumericToken {
	return &NumericToken{data: float64(value)}
}

// Data returns the numeric value of this token.
func (t *NumericToken) Data() float64 {
	return t.data
}

// HasDecimalPlaces reports whether the number has a non-zero decimal part.
func (t *NumericToken) HasDecimalPlaces() bool {
	return math.Floor(t.data) != t.data
}

// IntVal returns the value of this number as an int.
func (t *NumericToken) IntVal() int {
	return int(t.data)
}

// LongVal returns the value of this number as a long.
func (t *NumericToken) LongVal() int64 {
	return int64(t.data)
}

// DoubleVal returns the value of this number as a double.
func (t *NumericToken) DoubleVal() float64 {
	return t.data
}

// Equals reports whether other is a NumericToken with the same Data value.
func (t *NumericToken) Equals(other Token) bool {
	o, ok := other.(*NumericToken)
	if !ok {
		return false
	}
	return t.data == o.data
}

// String returns the string representation of the numeric token using invariant culture formatting.
func (t *NumericToken) String() string {
	return fmt.Sprintf("%g", t.data)
}
