package document_layout_analysis

import (
	"math"

	"github.com/uglytoad/pdfpig/go/core"
)

// Euclidean returns the straight-line distance between two points.
func Euclidean(point1, point2 core.PdfPoint) float64 {
	dx := point1.X - point2.X
	dy := point1.Y - point2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// WeightedEuclidean returns the weighted Euclidean distance between two points.
func WeightedEuclidean(point1, point2 core.PdfPoint, wX, wY float64) float64 {
	dx := point1.X - point2.X
	dy := point1.Y - point2.Y
	return math.Sqrt(wX*dx*dx + wY*dy*dy)
}

// Manhattan returns the sum of absolute differences of Cartesian coordinates (L1 distance).
func Manhattan(point1, point2 core.PdfPoint) float64 {
	return math.Abs(point1.X-point2.X) + math.Abs(point1.Y-point2.Y)
}

// Angle returns the angle in degrees between the horizontal axis and the line
// connecting startPoint to endPoint. Result is in [-180, 180].
func Angle(startPoint, endPoint core.PdfPoint) float64 {
	return math.Atan2(endPoint.Y-startPoint.Y, endPoint.X-startPoint.X) * 180 / math.Pi
}

// Vertical returns the absolute difference of the Y coordinates of two points.
func Vertical(point1, point2 core.PdfPoint) float64 {
	return math.Abs(point2.Y - point1.Y)
}

// Horizontal returns the absolute difference of the X coordinates of two points.
func Horizontal(point1, point2 core.PdfPoint) float64 {
	return math.Abs(point2.X - point1.X)
}

// BoundAngle180 bounds angle so that -180 <= theta <= 180.
func BoundAngle180(angle float64) float64 {
	angle = fmod(angle+180, 360)
	if angle < 0 {
		angle += 360
	}
	return angle - 180
}

// BoundAngle0to360 bounds angle so that 0 <= theta <= 360.
func BoundAngle0to360(angle float64) float64 {
	angle = fmod(angle, 360)
	if angle < 0 {
		angle += 360
	}
	return angle
}

// MinimumEditDistance returns the Levenshtein distance between two strings.
func MinimumEditDistance(s1, s2 string) int {
	runes1 := []rune(s1)
	runes2 := []rune(s2)
	n := len(runes1)
	m := len(runes2)

	d := make([][]uint16, n+1)
	for i := range d {
		d[i] = make([]uint16, m+1)
	}

	for i := 1; i <= n; i++ {
		d[i][0] = uint16(i)
	}
	for j := 1; j <= m; j++ {
		d[0][j] = uint16(j)
	}

	for j := 1; j <= m; j++ {
		for i := 1; i <= n; i++ {
			cost := uint16(0)
			if runes1[i-1] != runes2[j-1] {
				cost = 1
			}
			d[i][j] = min3(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}

	return int(d[n][m])
}

// MinimumEditDistanceNormalised returns the normalised Levenshtein distance between
// two strings, in range [0, 1]. A value of 0 means the strings are identical.
func MinimumEditDistanceNormalised(s1, s2 string) float64 {
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}
	if maxLen == 0 {
		return 0
	}
	return float64(MinimumEditDistance(s1, s2)) / float64(maxLen)
}

// FindIndexNearestPoint finds the index of the nearest element to element in candidates,
// excluding itself. It returns -1 if no distinct candidate is found.
// distanceMeasure computes the distance between two points.
func FindIndexNearestPoint[T comparable](
	element T,
	candidates []T,
	pivotPoint func(T) core.PdfPoint,
	candidatePoint func(T) core.PdfPoint,
	distanceMeasure func(core.PdfPoint, core.PdfPoint) float64,
) (int, float64) {
	if len(candidates) == 0 || distanceMeasure == nil {
		panic("FindIndexNearestPoint: candidates must be non-empty and distanceMeasure must not be nil")
	}

	distance := math.MaxFloat64
	closestIdx := -1
	pivot := pivotPoint(element)

	for i, c := range candidates {
		currentDistance := distanceMeasure(pivot, candidatePoint(c))
		if currentDistance < distance && c != element {
			distance = currentDistance
			closestIdx = i
		}
	}

	return closestIdx, distance
}

// FindIndexNearestLine finds the index of the nearest element to element in candidates
// by line proximity, excluding itself. It returns -1 if no distinct candidate is found.
// distanceMeasure computes the distance between two lines.
func FindIndexNearestLine[T comparable](
	element T,
	candidates []T,
	pivotLine func(T) core.PdfLine,
	candidateLine func(T) core.PdfLine,
	distanceMeasure func(core.PdfLine, core.PdfLine) float64,
) (int, float64) {
	if len(candidates) == 0 || distanceMeasure == nil {
		panic("FindIndexNearestLine: candidates must be non-empty and distanceMeasure must not be nil")
	}

	distance := math.MaxFloat64
	closestIdx := -1
	pivot := pivotLine(element)

	for i, c := range candidates {
		currentDistance := distanceMeasure(pivot, candidateLine(c))
		if currentDistance < distance && c != element {
			distance = currentDistance
			closestIdx = i
		}
	}

	return closestIdx, distance
}

func min3(a, b, c uint16) uint16 {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

// fmod computes the floating-point remainder of x/y, truncating toward zero.
// Equivalent to C#'s % operator on doubles.
func fmod(x, y float64) float64 {
	return x - math.Trunc(x/y)*y
}
