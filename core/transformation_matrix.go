package core

import (
	"errors"
	"fmt"
	"math"
)

// TransformationMatrix specifies the conversion from the transformed coordinate space
// to the original untransformed coordinate space.
type TransformationMatrix struct {
	A    float64 // The value at (0, 0) - The scale for the X dimension.
	B    float64 // The value at (0, 1).
	row1 float64
	C    float64 // The value at (1, 0).
	D    float64 // The value at (1, 1) - The scale for the Y dimension.
	row2 float64
	E    float64 // The value at (2, 0) - translation in X.
	F    float64 // The value at (2, 1) - translation in Y.
	row3 float64
}

// Rows is the number of rows in the matrix.
const Rows = 3

// Columns is the number of columns in the matrix.
const Columns = 3

// Identity is the default identity transformation matrix.
var Identity = TransformationMatrix{1, 0, 0, 0, 1, 0, 0, 0, 1}

// NewTransformationMatrix creates a new TransformationMatrix from 9 values.
func NewTransformationMatrix(a, b, r1, c, d, r2, e, f, r3 float64) TransformationMatrix {
	return TransformationMatrix{A: a, B: b, row1: r1, C: c, D: d, row2: r2, E: e, F: f, row3: r3}
}

// NewTransformationMatrixFromSlice creates a new TransformationMatrix from the values.
// Accepts either all 9 values of the matrix, 6 values in the default PDF order, or the 4 values of the top left square.
func NewTransformationMatrixFromSlice(values []float64) TransformationMatrix {
	if len(values) == 9 {
		return NewTransformationMatrix(values[0], values[1], values[2], values[3], values[4], values[5], values[6], values[7], values[8])
	}

	if len(values) == 6 {
		return NewTransformationMatrix(values[0], values[1], 0, values[2], values[3], 0, values[4], values[5], 1)
	}

	if len(values) == 4 {
		return NewTransformationMatrix(values[0], values[1], 0, values[2], values[3], 0, 0, 0, 1)
	}

	panic(fmt.Sprintf("The array must either define all 9 elements of the matrix or all 6 key elements. Instead array was: %v", values))
}

// GetTranslationMatrix creates a new TransformationMatrix with the X and Y translation values set.
func GetTranslationMatrix(x, y float64) TransformationMatrix {
	return NewTransformationMatrix(1, 0, 0, 0, 1, 0, x, y, 1)
}

// GetScaleMatrix creates a new TransformationMatrix with the X and Y scaling values set.
func GetScaleMatrix(scaleX, scaleY float64) TransformationMatrix {
	return NewTransformationMatrix(scaleX, 0, 0, 0, scaleY, 0, 0, 0, 1)
}

// GetRotationMatrix creates a new TransformationMatrix with the rotation value set.
func GetRotationMatrix(degreesCounterclockwise float64) TransformationMatrix {
	var cosVal, sinVal float64

	deg := degreesCounterclockwise - 360*math.Floor(degreesCounterclockwise/360)
	if deg < 0 {
		deg += 360
	}

	switch int(deg) {
	case 0:
		cosVal = 1
		sinVal = 0
	case 90:
		cosVal = 0
		sinVal = 1
	case 180:
		cosVal = -1
		sinVal = 0
	case 270:
		cosVal = 0
		sinVal = -1
	default:
		rad := degreesCounterclockwise * (math.Pi / 180)
		cosVal = math.Cos(rad)
		sinVal = math.Sin(rad)
	}

	return NewTransformationMatrix(cosVal, sinVal, 0, -sinVal, cosVal, 0, 0, 0, 1)
}

// GetAt gets the value at the specific row and column.
func (m TransformationMatrix) GetAt(row, col int) (float64, error) {
	if row < 0 {
		return 0, fmt.Errorf("cannot access negative rows in a matrix")
	}
	if row >= Rows {
		return 0, fmt.Errorf("the transformation matrix only contains %d rows and is zero indexed, you tried to access row %d", Rows, row)
	}
	if col < 0 {
		return 0, fmt.Errorf("cannot access negative columns in a matrix")
	}
	if col >= Columns {
		return 0, fmt.Errorf("the transformation matrix only contains %d columns and is zero indexed, you tried to access column %d", Columns, col)
	}

	switch row {
	case 0:
		switch col {
		case 0:
			return m.A, nil
		case 1:
			return m.B, nil
		case 2:
			return m.row1, nil
		}
	case 1:
		switch col {
		case 0:
			return m.C, nil
		case 1:
			return m.D, nil
		case 2:
			return m.row2, nil
		}
	case 2:
		switch col {
		case 0:
			return m.E, nil
		case 1:
			return m.F, nil
		case 2:
			return m.row3, nil
		}
	}

	return 0, fmt.Errorf("trying to access %d, %d which was not in the value array", row, col)
}

// TransformPoint transforms a point using this transformation matrix.
func (m TransformationMatrix) TransformPoint(original PdfPoint) PdfPoint {
	x := m.A*original.X + m.C*original.Y + m.E
	y := m.B*original.X + m.D*original.Y + m.F
	return NewPdfPoint(x, y)
}

// Transform transforms x and y coordinates using this transformation matrix.
func (m TransformationMatrix) Transform(x, y float64) (float64, float64) {
	newX := m.A*x + m.C*y + m.E
	newY := m.B*x + m.D*y + m.F
	return newX, newY
}

// TransformX transforms an X coordinate using this transformation matrix.
func (m TransformationMatrix) TransformX(x float64) float64 {
	return m.A*x + m.E
}

// TransformY transforms a Y coordinate using this transformation matrix.
func (m TransformationMatrix) TransformY(y float64) float64 {
	return m.D*y + m.F
}

// TransformRect transforms a rectangle using this transformation matrix.
func (m TransformationMatrix) TransformRect(original PdfRectangle) PdfRectangle {
	return NewPdfRectangleFromCorners(
		m.TransformPoint(original.TopLeft),
		m.TransformPoint(original.TopRight),
		m.TransformPoint(original.BottomLeft),
		m.TransformPoint(original.BottomRight),
	)
}

// TransformSubpath transforms a subpath using this transformation matrix.
func (m TransformationMatrix) TransformSubpath(subpath *PdfSubpath) (*PdfSubpath, error) {
	trSubpath := NewPdfSubpath()

	for _, c := range subpath.Commands() {
		switch cmd := c.(type) {
		case *Move:
			loc := m.TransformPoint(cmd.Location)
			trSubpath.MoveTo(loc.X, loc.Y)
		case *Line:
			to := m.TransformPoint(cmd.To)
			if err := trSubpath.LineTo(to.X, to.Y); err != nil {
				return nil, err
			}
		case *CubicBezierCurve:
			first := m.TransformPoint(cmd.FirstControlPoint)
			second := m.TransformPoint(cmd.SecondControlPoint)
			end := m.TransformPoint(cmd.EndPoint)
			if err := trSubpath.BezierCurveToCubic(first.X, first.Y, second.X, second.Y, end.X, end.Y); err != nil {
				return nil, err
			}
		case *QuadraticBezierCurve:
			control := m.TransformPoint(cmd.ControlPoint)
			end := m.TransformPoint(cmd.EndPoint)
			if err := trSubpath.BezierCurveToQuadratic(control.X, control.Y, end.X, end.Y); err != nil {
				return nil, err
			}
		case *Close:
			trSubpath.CloseSubpath()
		default:
			return nil, errors.New("unknown PdfSubpath command type")
		}
	}

	return trSubpath, nil
}

// TransformPath transforms a path (list of subpaths) using this transformation matrix.
func (m TransformationMatrix) TransformPath(path []*PdfSubpath) ([]*PdfSubpath, error) {
	result := make([]*PdfSubpath, 0, len(path))
	for _, subpath := range path {
		transformed, err := m.TransformSubpath(subpath)
		if err != nil {
			return nil, err
		}
		result = append(result, transformed)
	}
	return result, nil
}

// Translate generates a TransformationMatrix translated by the specified amount.
func (m TransformationMatrix) Translate(x, y float64) TransformationMatrix {
	e := x*m.A + y*m.C + m.E
	f := x*m.B + y*m.D + m.F
	r3 := x*m.row1 + y*m.row2 + m.row3

	return NewTransformationMatrix(
		m.A, m.B, m.row1,
		m.C, m.D, m.row2,
		e, f, r3,
	)
}

// FromValues creates a new TransformationMatrix from 6 values in the default PDF order.
func FromValues(a, b, c, d, e, f float64) TransformationMatrix {
	return NewTransformationMatrix(a, b, 0, c, d, 0, e, f, 1)
}

// FromValues4 creates a new TransformationMatrix from 4 values in the default PDF order.
func FromValues4(a, b, c, d float64) TransformationMatrix {
	return NewTransformationMatrix(a, b, 0, c, d, 0, 0, 0, 1)
}

// FromArray creates a new TransformationMatrix from the given slice of values.
// The slice must contain either 9, 6, or 4 elements.
func FromArray(values []float64) (TransformationMatrix, error) {
	switch len(values) {
	case 9:
		return NewTransformationMatrixFromSlice(values), nil
	case 6:
		return NewTransformationMatrix(values[0], values[1], 0, values[2], values[3], 0, values[4], values[5], 1), nil
	case 4:
		return NewTransformationMatrix(values[0], values[1], 0, values[2], values[3], 0, 0, 0, 1), nil
	default:
		return TransformationMatrix{}, fmt.Errorf("the array must either define all 9 elements of the matrix or all 6 key elements. Instead array was: %v", values)
	}
}

// Multiply multiplies one transformation matrix by another without modifying either matrix.
// Order is: (this * matrix).
func (m TransformationMatrix) Multiply(matrix TransformationMatrix) TransformationMatrix {
	a := m.A*matrix.A + m.B*matrix.C + m.row1*matrix.E
	b := m.A*matrix.B + m.B*matrix.D + m.row1*matrix.F
	r1 := m.A*matrix.row1 + m.B*matrix.row2 + m.row1*matrix.row3

	c := m.C*matrix.A + m.D*matrix.C + m.row2*matrix.E
	d := m.C*matrix.B + m.D*matrix.D + m.row2*matrix.F
	r2 := m.C*matrix.row1 + m.D*matrix.row2 + m.row2*matrix.row3

	e := m.E*matrix.A + m.F*matrix.C + m.row3*matrix.E
	f := m.E*matrix.B + m.F*matrix.D + m.row3*matrix.F
	r3 := m.E*matrix.row1 + m.F*matrix.row2 + m.row3*matrix.row3

	return NewTransformationMatrix(a, b, r1, c, d, r2, e, f, r3)
}

// MultiplyScalar multiplies the matrix by a scalar value without modifying this matrix.
func (m TransformationMatrix) MultiplyScalar(scalar float64) TransformationMatrix {
	return NewTransformationMatrix(
		m.A*scalar, m.B*scalar, m.row1*scalar,
		m.C*scalar, m.D*scalar, m.row2*scalar,
		m.E*scalar, m.F*scalar, m.row3*scalar,
	)
}

// Inverse returns the inverse of the current matrix.
func (m TransformationMatrix) Inverse() TransformationMatrix {
	a := m.D*m.row3 - m.row2*m.F
	c := -(m.C*m.row3 - m.row2*m.E)
	e := m.C*m.F - m.D*m.E

	b := -(m.B*m.row3 - m.row1*m.F)
	d := m.A*m.row3 - m.row1*m.E
	f := -(m.A*m.F - m.B*m.E)

	r1 := m.B*m.row2 - m.row1*m.D
	r2 := -(m.A*m.row2 - m.row1*m.C)
	r3 := m.A*m.D - m.B*m.C
	det := m.A*a + m.B*c + m.row1*e

	return NewTransformationMatrix(
		a/det, b/det, r1/det,
		c/det, d/det, r2/det,
		e/det, f/det, r3/det,
	)
}

// GetScalingFactorX returns the X scaling component of the current matrix.
func (m TransformationMatrix) GetScalingFactorX() float64 {
	xScale := m.A

	if m.B != 0 || m.C != 0 {
		xScale = math.Sqrt(m.A*m.A + m.B*m.B)
	}

	return xScale
}

// Equals reports whether m and other have the same matrix elements.
func (m TransformationMatrix) Equals(other TransformationMatrix) bool {
	return m.row1 == other.row1 &&
		m.row2 == other.row2 &&
		m.row3 == other.row3 &&
		m.A == other.A &&
		m.B == other.B &&
		m.C == other.C &&
		m.D == other.D &&
		m.E == other.E &&
		m.F == other.F
}

// EqualsMatrices determines whether two transformation matrices are equal.
func EqualsMatrices(a, b TransformationMatrix) bool {
	return a.Equals(b)
}

// String returns a string representation of the matrix in row-major form.
func (m TransformationMatrix) String() string {
	return fmt.Sprintf("%g, %g, %g\n%g, %g, %g\n%g, %g, %g",
		m.A, m.B, m.row1,
		m.C, m.D, m.row2,
		m.E, m.F, m.row3)
}
