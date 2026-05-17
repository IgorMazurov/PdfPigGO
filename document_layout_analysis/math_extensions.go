package document_layout_analysis

import (
	"math"
)

// ModeFloat32 computes the mode of a slice of float32 values.
// Returns math.NaN if the slice is empty or has no unique mode.
func ModeFloat32(values []float32) float32 {
	if len(values) == 0 {
		return float32(math.NaN())
	}

	counts := make(map[float32]int)
	for _, v := range values {
		counts[v]++
	}

	bestCount := 0
	var bestKey float32
	for k, c := range counts {
		if c > bestCount {
			bestCount = c
			bestKey = k
		}
	}

	tieCount := 0
	for _, c := range counts {
		if c == bestCount {
			tieCount++
		}
	}

	if tieCount > 1 {
		return float32(math.NaN())
	}

	return bestKey
}

// ModeFloat64 computes the mode of a slice of float64 values.
// Returns math.NaN if the slice is empty or has no unique mode.
func ModeFloat64(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}

	counts := make(map[float64]int)
	for _, v := range values {
		counts[v]++
	}

	bestCount := 0
	var bestKey float64
	for k, c := range counts {
		if c > bestCount {
			bestCount = c
			bestKey = k
		}
	}

	tieCount := 0
	for _, c := range counts {
		if c == bestCount {
			tieCount++
		}
	}

	if tieCount > 1 {
		return math.NaN()
	}

	return bestKey
}

// AlmostEqualsToZero tests whether number is approximately equal to zero,
// within the given epsilon. Default epsilon is 1e-5.
func AlmostEqualsToZero(number float64, epsilon ...float64) bool {
	eps := 1e-5
	if len(epsilon) > 0 {
		eps = epsilon[0]
	}
	return number > -eps && number < eps
}

// AlmostEquals tests whether two numbers are approximately equal,
// within the given epsilon. Default epsilon is 1e-5.
func AlmostEquals(number, other float64, epsilon ...float64) bool {
	return AlmostEqualsToZero(number-other, epsilon...)
}
