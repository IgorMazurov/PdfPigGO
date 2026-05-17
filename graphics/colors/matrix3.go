package colors

// mat3 represents a 3x3 matrix stored in column-major order.
type mat3 [9]float64

// newMat3 creates a 3x3 matrix from row-major elements.
func newMat3(m11, m12, m13, m21, m22, m23, m31, m32, m33 float64) mat3 {
	return mat3{m11, m12, m13, m21, m22, m23, m31, m32, m33}
}

// identityMat3 is the 3x3 identity matrix.
var identityMat3 = newMat3(1, 0, 0, 0, 1, 0, 0, 0, 1)

// mat3Multiply returns the product of two 3x3 matrices.
func mat3Multiply(a, b mat3) mat3 {
	return newMat3(
		a[0]*b[0]+a[1]*b[3]+a[2]*b[6],
		a[0]*b[1]+a[1]*b[4]+a[2]*b[7],
		a[0]*b[2]+a[1]*b[5]+a[2]*b[8],
		a[3]*b[0]+a[4]*b[3]+a[5]*b[6],
		a[3]*b[1]+a[4]*b[4]+a[5]*b[7],
		a[3]*b[2]+a[4]*b[5]+a[5]*b[8],
		a[6]*b[0]+a[7]*b[3]+a[8]*b[6],
		a[6]*b[1]+a[7]*b[4]+a[8]*b[7],
		a[6]*b[2]+a[7]*b[5]+a[8]*b[8],
	)
}

// mat3MultiplyVec multiplies a 3x3 matrix by a 3-element vector.
func mat3MultiplyVec(m mat3, v xyzTriplet) xyzTriplet {
	return xyzTriplet{
		X: m[0]*v.X + m[1]*v.Y + m[2]*v.Z,
		Y: m[3]*v.X + m[4]*v.Y + m[5]*v.Z,
		Z: m[6]*v.X + m[7]*v.Y + m[8]*v.Z,
	}
}

// mat3Inverse returns the inverse of a 3x3 matrix.
func mat3Inverse(m mat3) mat3 {
	det := m[0]*(m[4]*m[8]-m[5]*m[7]) -
		m[1]*(m[3]*m[8]-m[5]*m[6]) +
		m[2]*(m[3]*m[7]-m[4]*m[6])

	inverseDet := 1.0 / det

	return newMat3(
		(m[4]*m[8]-m[5]*m[7])*inverseDet,
		(m[2]*m[7]-m[1]*m[8])*inverseDet,
		(m[1]*m[5]-m[2]*m[4])*inverseDet,
		(m[5]*m[6]-m[3]*m[8])*inverseDet,
		(m[0]*m[8]-m[2]*m[6])*inverseDet,
		(m[3]*m[2]-m[0]*m[5])*inverseDet,
		(m[3]*m[7]-m[4]*m[6])*inverseDet,
		(m[1]*m[6]-m[0]*m[7])*inverseDet,
		(m[0]*m[4]-m[1]*m[3])*inverseDet,
	)
}
