package util

import (
	"fmt"
)

// Matrix3x3 represents a 3x3 matrix with the layout:
// | m11 m12 m13 |
// | m21 m22 m23 |
// | m31 m32 m33 |
type Matrix3x3 struct {
	m11, m12, m13 float64
	m21, m22, m23 float64
	m31, m32, m33 float64
}

// Identity is the 3x3 identity matrix. Multiplying any matrix by the
// identity matrix yields the original matrix.
var Identity = NewMatrix3x3(1, 0, 0, 0, 1, 0, 0, 0, 1)

// NewMatrix3x3 creates a new 3x3 matrix with the given elements.
func NewMatrix3x3(m11, m12, m13, m21, m22, m23, m31, m32, m33 float64) Matrix3x3 {
	return Matrix3x3{
		m11: m11, m12: m12, m13: m13,
		m21: m21, m22: m22, m23: m23,
		m31: m31, m32: m32, m33: m33,
	}
}

// Elements returns all nine matrix elements as a slice in row-major order.
func (m Matrix3x3) Elements() []float64 {
	return []float64{m.m11, m.m12, m.m13, m.m21, m.m22, m.m23, m.m31, m.m32, m.m33}
}

// Inverse returns the inverse of this matrix. If the determinant is zero,
// an error is returned since no inverse exists.
func (m Matrix3x3) Inverse() (Matrix3x3, error) {
	determinant := m.getDeterminant()

	if determinant == 0 {
		return Matrix3x3{}, fmt.Errorf("may not inverse a matrix with a determinant of 0")
	}

	transposed := m.Transpose()

	minorm11 := transposed.m22*transposed.m33 - transposed.m23*transposed.m32
	minorm12 := transposed.m21*transposed.m33 - transposed.m23*transposed.m31
	minorm13 := transposed.m21*transposed.m32 - transposed.m22*transposed.m31

	minorm21 := transposed.m12*transposed.m33 - transposed.m13*transposed.m32
	minorm22 := transposed.m11*transposed.m33 - transposed.m13*transposed.m31
	minorm23 := transposed.m11*transposed.m32 - transposed.m12*transposed.m31

	minorm31 := transposed.m12*transposed.m23 - transposed.m13*transposed.m22
	minorm32 := transposed.m11*transposed.m23 - transposed.m13*transposed.m21
	minorm33 := transposed.m11*transposed.m22 - transposed.m12*transposed.m21

	adjugate := NewMatrix3x3(
		minorm11, -minorm12, minorm13,
		-minorm21, minorm22, -minorm23,
		minorm31, -minorm32, minorm33)

	return adjugate.MultiplyScalar(1 / determinant), nil
}

// MultiplyScalar returns a new matrix with each element multiplied by the factor.
func (m Matrix3x3) MultiplyScalar(factor float64) Matrix3x3 {
	return NewMatrix3x3(
		m.m11*factor, m.m12*factor, m.m13*factor,
		m.m21*factor, m.m22*factor, m.m23*factor,
		m.m31*factor, m.m32*factor, m.m33*factor)
}

// Vector3 is a 3-element vector used for matrix-vector multiplication.
type Vector3 struct {
	X, Y, Z float64
}

// MultiplyVector multiplies this matrix by a 3-element vector and returns
// the resulting 3-element vector.
func (m Matrix3x3) MultiplyVector(v Vector3) Vector3 {
	return Vector3{
		m.m11*v.X + m.m12*v.Y + m.m13*v.Z,
		m.m21*v.X + m.m22*v.Y + m.m23*v.Z,
		m.m31*v.X + m.m32*v.Y + m.m33*v.Z,
	}
}

// Multiply returns the dot product (matrix multiplication) of this matrix
// and the supplied matrix.
func (m Matrix3x3) Multiply(other Matrix3x3) Matrix3x3 {
	return NewMatrix3x3(
		m.m11*other.m11+m.m12*other.m21+m.m13*other.m31,
		m.m11*other.m12+m.m12*other.m22+m.m13*other.m32,
		m.m11*other.m13+m.m12*other.m23+m.m13*other.m33,
		m.m21*other.m11+m.m22*other.m21+m.m23*other.m31,
		m.m21*other.m12+m.m22*other.m22+m.m23*other.m32,
		m.m21*other.m13+m.m22*other.m23+m.m23*other.m33,
		m.m31*other.m11+m.m32*other.m21+m.m33*other.m31,
		m.m31*other.m12+m.m32*other.m22+m.m33*other.m32,
		m.m31*other.m13+m.m32*other.m23+m.m33*other.m33)
}

// Transpose returns a new matrix that is the transpose of this matrix
// (rows and columns interchanged).
func (m Matrix3x3) Transpose() Matrix3x3 {
	return NewMatrix3x3(
		m.m11, m.m21, m.m31,
		m.m12, m.m22, m.m32,
		m.m13, m.m23, m.m33)
}

// Equals reports whether this matrix is equal to another.
func (m Matrix3x3) Equals(other Matrix3x3) bool {
	return m.m11 == other.m11 &&
		m.m12 == other.m12 &&
		m.m13 == other.m13 &&
		m.m21 == other.m21 &&
		m.m22 == other.m22 &&
		m.m23 == other.m23 &&
		m.m31 == other.m31 &&
		m.m32 == other.m32 &&
		m.m33 == other.m33
}

// getDeterminant computes the determinant of this 3x3 matrix.
func (m Matrix3x3) getDeterminant() float64 {
	minorM11 := m.m22*m.m33 - m.m23*m.m32
	minorM12 := m.m21*m.m33 - m.m23*m.m31
	minorM13 := m.m21*m.m32 - m.m22*m.m31

	return m.m11*minorM11 - m.m12*minorM12 + m.m13*minorM13
}
