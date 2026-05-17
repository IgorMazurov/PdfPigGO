package core

import (
	"math"
	"testing"
)

func TestTransformationMatrixMultipliesCorrectly(t *testing.T) {
	tests := []struct {
		name     string
		a        []float64
		b        []float64
		expected []float64
	}{
		{
			name: "identity times identity",
			a:    []float64{1, 0, 0, 0, 1, 0, 0, 0, 1},
			b:    []float64{1, 0, 0, 0, 1, 0, 0, 0, 1},
			expected: []float64{1, 0, 0, 0, 1, 0, 0, 0, 1},
		},
		{
			name: "general multiplication 1",
			a:    []float64{65, 9, 3, 5, 2, 7, 11, 1, 6},
			b:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expected: []float64{122, 199, 276, 62, 76, 90, 57, 75, 93},
		},
		{
			name: "general multiplication with negatives",
			a:    []float64{3, 5, 7, 11, 13, -3, 17, -6, -9},
			b:    []float64{5, 4, 3, 3, 7, 12, 1, 0, 6},
			expected: []float64{37, 47, 111, 91, 135, 171, 58, 26, -75},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matrixA, err := FromArray(tc.a)
			if err != nil {
				t.Fatalf("FromArray(a) error: %v", err)
			}
			matrixB, err := FromArray(tc.b)
			if err != nil {
				t.Fatalf("FromArray(b) error: %v", err)
			}
			expectedMatrix, err := FromArray(tc.expected)
			if err != nil {
				t.Fatalf("FromArray(expected) error: %v", err)
			}

			result := matrixA.Multiply(matrixB)

			if !result.Equals(expectedMatrix) {
				t.Errorf("Multiply result mismatch:\nexpected: %s\n  actual: %s", expectedMatrix.String(), result.String())
			}
		})
	}
}

func TestTransformationMatrixInversesCorrectly(t *testing.T) {
	tests := []struct {
		name     string
		a        []float64
		expected []float64
	}{
		{
			name:     "inverse case 1",
			a:        []float64{0, 2, 2, 1, 1, 1, 0, 1, 2},
			expected: []float64{-0.5, 1, 0, 1, 0, -1, -0.5, 0, 1},
		},
		{
			name:     "inverse case 2",
			a:        []float64{1, 1, 0, 0, 1, 0, 2, 1, 1},
			expected: []float64{1, -1, 0, 0, 1, 0, -2, 1, 1},
		},
		{
			name:     "inverse case 3",
			a:        []float64{2, 0, 0, 0, 2, 1, 2, 0, 2},
			expected: []float64{0.5, 0, 0, 0.25, 0.5, -0.25, -0.5, 0, 0.5},
		},
		{
			name: "inverse case 4 decimal values",
			a:    []float64{-4.68, 2.47, 3.12, 5.00, -6.19, -0.58, 9.37, -7.11, 4.51},
			expected: []float64{
				-0.212368047525923, -0.220866560683805, 0.118511242369019,
				-0.18548392709254, -0.333665068307247, 0.0854066769202931,
				0.14880219150553, -0.0671483286158035, 0.110153244324962,
			},
		},
	}

	tolerance := 1e-6

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matrixA, err := FromArray(tc.a)
			if err != nil {
				t.Fatalf("FromArray(a) error: %v", err)
			}
			expectedMatrix, err := FromArray(tc.expected)
			if err != nil {
				t.Fatalf("FromArray(expected) error: %v", err)
			}

			result := matrixA.Inverse()

			for i := 0; i < Rows; i++ {
				for j := 0; j < Columns; j++ {
					expVal, err := expectedMatrix.GetAt(i, j)
					if err != nil {
						t.Fatalf("expectedMatrix.GetAt(%d, %d) error: %v", i, j, err)
					}
					resVal, err := result.GetAt(i, j)
					if err != nil {
						t.Fatalf("result.GetAt(%d, %d) error: %v", i, j, err)
					}

					if math.Abs(resVal-expVal) > tolerance {
						t.Errorf("element [%d][%d]: expected %g, got %g (diff=%g)", i, j, expVal, resVal, math.Abs(resVal-expVal))
					}
				}
			}
		})
	}
}
