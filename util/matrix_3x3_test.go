package util

import (
	"testing"
)

func TestMatrix3x3CanCreate(t *testing.T) {
	matrix := NewMatrix3x3(1, 2, 3, 4, 5, 6, 7, 8, 9)
	expected := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}

	got := matrix.Elements()
	if len(got) != len(expected) {
		t.Fatalf("Elements() length = %d, want %d", len(got), len(expected))
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("Elements()[%d] = %v, want %v", i, got[i], expected[i])
		}
	}
}

func TestMatrix3x3ExposesIdentity(t *testing.T) {
	matrix := NewMatrix3x3(1, 0, 0, 0, 1, 0, 0, 0, 1)

	if !matrix.Equals(Identity) {
		t.Errorf("identity matrix not equal to Identity singleton: got %+v, want %+v", matrix, Identity)
	}
}

func TestMatrix3x3MultiplyScalar(t *testing.T) {
	matrix := NewMatrix3x3(1, 2, 3, 4, 5, 6, 7, 8, 9)

	result := matrix.MultiplyScalar(3)
	expected := []float64{3, 6, 9, 12, 15, 18, 21, 24, 27}

	got := result.Elements()
	if len(got) != len(expected) {
		t.Fatalf("Elements() length = %d, want %d", len(got), len(expected))
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("Elements()[%d] = %v, want %v", i, got[i], expected[i])
		}
	}
}

func TestMatrix3x3MultiplyVector(t *testing.T) {
	matrix := NewMatrix3x3(1, 2, 3, 4, 5, 6, 7, 8, 9)
	vector := Vector3{1, 2, 3}

	product := matrix.MultiplyVector(vector)

	wantX, wantY, wantZ := 14.0, 32.0, 50.0
	if product.X != wantX || product.Y != wantY || product.Z != wantZ {
		t.Errorf("MultiplyVector() = (%v, %v, %v), want (%v, %v, %v)", product.X, product.Y, product.Z, wantX, wantY, wantZ)
	}
}

func TestMatrix3x3MultiplyMatrix(t *testing.T) {
	matrix := NewMatrix3x3(1, 2, 3, 4, 5, 6, 7, 8, 9)

	product := matrix.Multiply(matrix)

	expected := NewMatrix3x3(30, 36, 42, 66, 81, 96, 102, 126, 150)

	if !expected.Equals(product) {
		t.Errorf("Multiply() = %+v, want %+v", product, expected)
	}
}

func TestMatrix3x3Transpose(t *testing.T) {
	matrix := NewMatrix3x3(1, 2, 3, 4, 5, 6, 7, 8, 9)

	transposed := matrix.Transpose()
	expected := []float64{1, 4, 7, 2, 5, 8, 3, 6, 9}

	got := transposed.Elements()
	if len(got) != len(expected) {
		t.Fatalf("Elements() length = %d, want %d", len(got), len(expected))
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("Elements()[%d] = %v, want %v", i, got[i], expected[i])
		}
	}
}

func TestMatrix3x3InverseReturnsMatrix(t *testing.T) {
	matrix := NewMatrix3x3(1, 2, 3, 0, 1, 4, 5, 6, 0)

	inverse, err := matrix.Inverse()
	if err != nil {
		t.Fatalf("Inverse() unexpected error: %v", err)
	}

	expected := []float64{-24, 18, 5, 20, -15, -4, -5, 4, 1}
	got := inverse.Elements()
	if len(got) != len(expected) {
		t.Fatalf("Elements() length = %d, want %d", len(got), len(expected))
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("Elements()[%d] = %v, want %v", i, got[i], expected[i])
		}
	}

	product := matrix.Multiply(inverse)
	if !product.Equals(Identity) {
		t.Errorf("matrix * inverse != Identity: got %+v, want %+v", product, Identity)
	}
}

func TestMatrix3x3InverseReturnsError(t *testing.T) {
	matrix := NewMatrix3x3(1, 2, 3, 4, 5, 6, 7, 8, 9)

	_, err := matrix.Inverse()
	if err == nil {
		t.Error("Inverse() expected error for singular matrix (determinant = 0), got nil")
	}
}
