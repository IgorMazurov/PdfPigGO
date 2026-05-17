package clipperlibrary

// Int128 represents a signed 128-bit integer as two 64-bit parts.
// Enables safe multiplication of two int64 values without overflow.
type Int128 struct {
	hi int64
	lo uint64
}

// NewInt128 creates an Int128 from a single int64 value with sign extension.
func NewInt128(lo int64) Int128 {
	var hi int64
	if lo < 0 {
		hi = -1
	}
	return Int128{hi: hi, lo: uint64(lo)}
}

// NewInt128FromParts creates an Int128 from explicit high and low parts.
func NewInt128FromParts(hi int64, lo uint64) Int128 {
	return Int128{hi: hi, lo: lo}
}

// IsNegative returns true if the value is negative.
func (v Int128) IsNegative() bool {
	return v.hi < 0
}

// Negate returns the negation of the value.
func (v Int128) Negate() Int128 {
	if v.lo == 0 {
		return NewInt128FromParts(-v.hi, 0)
	}
	return NewInt128FromParts(^v.hi, ^v.lo+1)
}

// ToFloat64 converts the value to float64.
func (v Int128) ToFloat64() float64 {
	const shift64 = 18446744073709551616.0 // 2^64
	if v.hi < 0 {
		if v.lo == 0 {
			return float64(v.hi) * shift64
		}
		return -(float64(^v.lo) + float64(^v.hi)*shift64)
	}
	return float64(v.lo) + float64(v.hi)*shift64
}

// Int128Mul multiplies two int64 values producing a precise 128-bit result.
func Int128Mul(lhs, rhs int64) Int128 {
	negate := (lhs < 0) != (rhs < 0)
	if lhs < 0 {
		lhs = -lhs
	}
	if rhs < 0 {
		rhs = -rhs
	}
	int1Hi := uint64(lhs) >> 32
	int1Lo := uint64(lhs) & 0xFFFFFFFF
	int2Hi := uint64(rhs) >> 32
	int2Lo := uint64(rhs) & 0xFFFFFFFF

	a := int1Hi * int2Hi
	b := int1Lo * int2Lo
	c := int1Hi*int2Lo + int1Lo*int2Hi

	hi := int64(a + (c >> 32))
	lo := (c << 32) + b
	if lo < b {
		hi++
	}
	result := NewInt128FromParts(hi, lo)
	if negate {
		return result.Negate()
	}
	return result
}
