package geometry

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/geometry/clipperlibrary"
)

const epsilon = 1e-5

// ccw returns true if the points are in counter-clockwise order.
func ccw(p1, p2, p3 core.PdfPoint) bool {
	return (p2.X-p1.X)*(p3.Y-p1.Y) > (p2.Y-p1.Y)*(p3.X-p1.X)
}

// --- PdfPoint extensions ---

// DotProduct returns the dot product of two points treated as vectors.
func DotProduct(p1, p2 core.PdfPoint) float64 {
	return p1.X*p2.X + p1.Y*p2.Y
}

// PointAdd returns a point with summed coordinates.
func PointAdd(p1, p2 core.PdfPoint) core.PdfPoint {
	return core.NewPdfPoint(p1.X+p2.X, p1.Y+p2.Y)
}

// PointSubtract returns a point with subtracted coordinates.
func PointSubtract(p1, p2 core.PdfPoint) core.PdfPoint {
	return core.NewPdfPoint(p1.X-p2.X, p1.Y-p2.Y)
}

// --- Minimum Area Rectangle & Convex Hull ---

// parametricPerpendicularProjection finds the minimal bounding rectangle of a convex polygon.
// The polygon must be simple and convex with vertices in strict cyclic sequential order.
func parametricPerpendicularProjection(polygon []core.PdfPoint) (core.PdfRectangle, error) {
	if len(polygon) == 0 {
		return core.PdfRectangle{}, fmt.Errorf("parametricPerpendicularProjection: polygon must contain at least one point")
	}

	if len(polygon) == 1 {
		return core.NewPdfRectangle(polygon[0], polygon[0]), nil
	}

	if len(polygon) == 2 {
		return core.NewPdfRectangle(polygon[0], polygon[1]), nil
	}

	mrb := make([]float64, 8)
	amin := math.Inf(1)
	j := 1
	k := 0

	var qx, qy, r0x, r0y, r1x, r1y float64
	qx, qy = math.NaN(), math.NaN()
	r0x, r0y = math.NaN(), math.NaN()
	r1x, r1y = math.NaN(), math.NaN()

	for {
		pk := polygon[k]
		pj := polygon[j]

		vx := pj.X - pk.X
		vy := pj.Y - pk.Y
		r := 1.0 / (vx*vx + vy*vy)

		tmin := 1.0
		tmax := 0.0
		smax := 0.0
		l := -1

		for j = 0; j < len(polygon); j++ {
			pj = polygon[j]
			ux := pj.X - pk.X
			uy := pj.Y - pk.Y
			t := (ux*vx + uy*vy) * r

			ptx := t*vx + pk.X
			pty := t*vy + pk.Y
			ux = ptx - pj.X
			uy = pty - pj.Y

			s := ux*ux + uy*uy

			if t < tmin {
				tmin = t
				r0x = ptx
				r0y = pty
			}

			if t > tmax {
				tmax = t
				r1x = ptx
				r1y = pty
			}

			if s > smax {
				smax = s
				qx = ptx
				qy = pty
				l = j
			}
		}

		if l != -1 {
			pl := polygon[l]
			plMinusQX := pl.X - qx
			plMinusQY := pl.Y - qy

			r2x := r1x + plMinusQX
			r2y := r1y + plMinusQY

			r3x := r0x + plMinusQX
			r3y := r0y + plMinusQY

			ux := r1x - r0x
			uy := r1y - r0y

			a := (ux*ux + uy*uy) * smax

			if a < amin {
				amin = a
				mrb[0] = r0x
				mrb[1] = r0y
				mrb[2] = r1x
				mrb[3] = r1y
				mrb[4] = r2x
				mrb[5] = r2y
				mrb[6] = r3x
				mrb[7] = r3y
			}
		}

		k++
		j = k + 1

		if j == len(polygon) {
			j = 0
		}
		if k == len(polygon) {
			break
		}
	}

	return core.NewPdfRectangleFromCorners(
		core.NewPdfPoint(mrb[4], mrb[5]),
		core.NewPdfPoint(mrb[6], mrb[7]),
		core.NewPdfPoint(mrb[2], mrb[3]),
		core.NewPdfPoint(mrb[0], mrb[1]),
	), nil
}

// MinimumAreaRectangle finds the oriented minimum area rectangle enclosing the given points
// by first computing their convex hull then finding its MAR.
func MinimumAreaRectangle(points []core.PdfPoint) (core.PdfRectangle, error) {
	if len(points) == 0 {
		return core.PdfRectangle{}, fmt.Errorf("minimumAreaRectangle: points cannot be empty")
	}

	distinct := make(map[core.PdfPoint]struct{})
	for _, p := range points {
		distinct[p] = struct{}{}
	}

	distinctPoints := make([]core.PdfPoint, 0, len(distinct))
	for p := range distinct {
		distinctPoints = append(distinctPoints, p)
	}

	slices.SortFunc(distinctPoints, func(a, b core.PdfPoint) int {
		if a.X != b.X {
			if a.X < b.X {
				return -1
			}
			return 1
		}
		if a.Y < b.Y {
			return -1
		}
		if a.Y > b.Y {
			return 1
		}
		return 0
	})

	hull, err := GrahamScan(distinctPoints)
	if err != nil {
		return core.PdfRectangle{}, fmt.Errorf("minimumAreaRectangle: graham scan failed: %w", err)
	}

	return parametricPerpendicularProjection(hull)
}

// OrientedBoundingBox computes the oriented bounding box by fitting a line through the points,
// rotating them to get an AABB, then rotating back.
func OrientedBoundingBox(points []core.PdfPoint) (core.PdfRectangle, error) {
	if len(points) < 2 {
		return core.PdfRectangle{}, fmt.Errorf("orientedBoundingBox: must contain at least two points")
	}

	sumX := 0.0
	sumY := 0.0
	for i := 0; i < len(points); i++ {
		sumX += points[i].X
		sumY += points[i].Y
	}
	x0 := sumX / float64(len(points))
	y0 := sumY / float64(len(points))

	sumProduct := 0.0
	sumDiffSquaredX := 0.0

	for i := 0; i < len(points); i++ {
		p := points[i]
		xDiff := p.X - x0
		yDiff := p.Y - y0
		sumProduct += xDiff * yDiff
		sumDiffSquaredX += xDiff * xDiff
	}

	slope := sumProduct / sumDiffSquaredX
	angleRad := math.Atan(slope)
	cos := math.Cos(angleRad)
	sin := math.Sin(angleRad)

	inverseRotation := core.NewTransformationMatrix(
		cos, -sin, 0,
		sin, cos, 0,
		0, 0, 1,
	)

	first := inverseRotation.TransformPoint(points[0])
	minX, minY := first.X, first.Y
	maxX, maxY := first.X, first.Y

	for i := 1; i < len(points); i++ {
		tp := inverseRotation.TransformPoint(points[i])
		if tp.X < minX {
			minX = tp.X
		}
		if tp.Y < minY {
			minY = tp.Y
		}
		if tp.X > maxX {
			maxX = tp.X
		}
		if tp.Y > maxY {
			maxY = tp.Y
		}
	}

	aabb := core.NewPdfRectangleFloat(minX, minY, maxX, maxY)

	rotateBack := core.NewTransformationMatrix(
		cos, sin, 0,
		-sin, cos, 0,
		0, 0, 1,
	)
	return rotateBack.TransformRect(aabb), nil
}

// GrahamScan computes the convex hull of points in O(n log n) time.
func GrahamScan(points []core.PdfPoint) ([]core.PdfPoint, error) {
	if len(points) == 0 {
		return nil, fmt.Errorf("grahamScan: points cannot be empty")
	}

	if len(points) < 3 {
		return slices.Clone(points), nil
	}

	polarAngle := func(p1, p2 core.PdfPoint) float64 {
		angle := math.Atan2(p2.Y-p1.Y, p2.X-p1.X)
		return math.Mod(angle, math.Pi)
	}

	slices.SortFunc(points, func(a, b core.PdfPoint) int {
		if a.X < b.X {
			return -1
		}
		if a.X > b.X {
			return 1
		}
		if a.Y < b.Y {
			return -1
		}
		if a.Y > b.Y {
			return 1
		}
		return 0
	})

	p0 := points[0]

	type groupEntry struct {
		key float64
		val core.PdfPoint
	}

	groupsMap := make(map[float64][]core.PdfPoint)
	for i := 1; i < len(points); i++ {
		key := polarAngle(p0, points[i])
		groupsMap[key] = append(groupsMap[key], points[i])
	}

	keys := make([]float64, 0, len(groupsMap))
	for k := range groupsMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	sortedPoints := make([]core.PdfPoint, len(keys))
	for i, key := range keys {
		group := groupsMap[key]
		var farthest core.PdfPoint
		maxDistSq := -1.0
		for _, p := range group {
			dx := p.X - p0.X
			dy := p.Y - p0.Y
			distSq := dx*dx + dy*dy
			if distSq > maxDistSq {
				maxDistSq = distSq
				farthest = p
			}
		}
		sortedPoints[i] = farthest
	}

	if len(keys) < 2 {
		return []core.PdfPoint{p0, sortedPoints[0]}, nil
	}

	hull := make([]core.PdfPoint, 0, len(keys)+1)
	hull = append(hull, p0)
	hull = append(hull, sortedPoints[0])
	hull = append(hull, sortedPoints[1])

	for i := 2; i < len(sortedPoints); i++ {
		point := sortedPoints[i]
		for len(hull) > 1 && !ccw(hull[len(hull)-2], hull[len(hull)-1], point) {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, point)
	}

	slices.Reverse(hull)
	return hull, nil
}

// --- PdfRectangle extensions ---

// RectangleToPdfPath converts a rectangle to its PdfPath representation.
func RectangleToPdfPath(rectangle core.PdfRectangle) *PdfPath {
	subpath := core.NewPdfSubpath()
	subpath.Rectangle(rectangle.BottomLeft.X, rectangle.BottomLeft.Y, rectangle.Width, rectangle.Height)
	path := NewPdfPath()
	path.Add(subpath)
	return path
}

// RectangleContainsPoint reports whether the point lies within the rectangle.
// If includeBorder is true, points on the boundary are considered inside.
func RectangleContainsPoint(rectangle core.PdfRectangle, point core.PdfPoint, includeBorder bool) bool {
	if math.Abs(rectangle.Area()) < epsilon {
		return false
	}

	if math.Abs(rectangle.Rotation()) < epsilon {
		if includeBorder {
			return point.X >= rectangle.Left() &&
				point.X <= rectangle.Right() &&
				point.Y >= rectangle.Bottom() &&
				point.Y <= rectangle.Top()
		}
		return point.X > rectangle.Left() &&
			point.X < rectangle.Right() &&
			point.Y > rectangle.Bottom() &&
			point.Y < rectangle.Top()
	}

	area := func(p1, p2, p3 core.PdfPoint) float64 {
		return math.Abs((p2.X*p1.Y-p1.X*p2.Y)+(p3.X*p2.Y-p2.X*p3.Y)+(p1.X*p3.Y-p3.X*p1.Y)) / 2.0
	}

	area1 := area(rectangle.BottomLeft, point, rectangle.TopLeft)
	area2 := area(rectangle.TopLeft, point, rectangle.TopRight)
	area3 := area(rectangle.TopRight, point, rectangle.BottomRight)
	area4 := area(rectangle.BottomRight, point, rectangle.BottomLeft)

	sum := area1 + area2 + area3 + area4

	if sum-rectangle.Area() > epsilon {
		return false
	}

	if area1 < epsilon || area2 < epsilon || area3 < epsilon || area4 < epsilon {
		return includeBorder
	}

	return true
}

// RectangleContainsRect reports whether the other rectangle is inside this one.
func RectangleContainsRect(rectangle, other core.PdfRectangle, includeBorder bool) bool {
	if !RectangleContainsPoint(rectangle, other.BottomLeft, includeBorder) {
		return false
	}
	if !RectangleContainsPoint(rectangle, other.TopRight, includeBorder) {
		return false
	}
	if !RectangleContainsPoint(rectangle, other.BottomRight, includeBorder) {
		return false
	}
	if !RectangleContainsPoint(rectangle, other.TopLeft, includeBorder) {
		return false
	}
	return true
}

// RectangleIntersectsWithRect reports whether two rectangles overlap (not just share a border).
func RectangleIntersectsWithRect(rectangle, other core.PdfRectangle) bool {
	if math.Abs(rectangle.Rotation()) < epsilon && math.Abs(other.Rotation()) < epsilon {
		if rectangle.Left() > other.Right() || other.Left() > rectangle.Right() {
			return false
		}
		if rectangle.Top() < other.Bottom() || other.Top() < rectangle.Bottom() {
			return false
		}
		return true
	}

	r1 := RectangleNormalise(rectangle)
	r2 := RectangleNormalise(other)
	if math.Abs(r1.Rotation()) < epsilon && math.Abs(r2.Rotation()) < epsilon {
		if !RectangleIntersectsWithRect(r1, r2) {
			return false
		}
	}

	if RectangleContainsPoint(rectangle, other.BottomLeft, false) {
		return true
	}
	if RectangleContainsPoint(rectangle, other.TopRight, false) {
		return true
	}
	if RectangleContainsPoint(rectangle, other.TopLeft, false) {
		return true
	}
	if RectangleContainsPoint(rectangle, other.BottomRight, false) {
		return true
	}

	if RectangleContainsPoint(other, rectangle.BottomLeft, false) {
		return true
	}
	if RectangleContainsPoint(other, rectangle.TopRight, false) {
		return true
	}
	if RectangleContainsPoint(other, rectangle.TopLeft, false) {
		return true
	}
	if RectangleContainsPoint(other, rectangle.BottomRight, false) {
		return true
	}

	if lineIntersects(rectangle.BottomLeft, rectangle.BottomRight, other.BottomLeft, other.BottomRight) {
		return true
	}
	if lineIntersects(rectangle.BottomLeft, rectangle.BottomRight, other.BottomRight, other.TopRight) {
		return true
	}
	if lineIntersects(rectangle.BottomLeft, rectangle.BottomRight, other.TopRight, other.TopLeft) {
		return true
	}
	if lineIntersects(rectangle.BottomLeft, rectangle.BottomRight, other.TopLeft, other.BottomLeft) {
		return true
	}

	if lineIntersects(rectangle.BottomRight, rectangle.TopRight, other.BottomLeft, other.BottomRight) {
		return true
	}
	if lineIntersects(rectangle.BottomRight, rectangle.TopRight, other.BottomRight, other.TopRight) {
		return true
	}
	if lineIntersects(rectangle.BottomRight, rectangle.TopRight, other.TopRight, other.TopLeft) {
		return true
	}
	if lineIntersects(rectangle.BottomRight, rectangle.TopRight, other.TopLeft, other.BottomLeft) {
		return true
	}

	if lineIntersects(rectangle.TopRight, rectangle.TopLeft, other.BottomLeft, other.BottomRight) {
		return true
	}
	if lineIntersects(rectangle.TopRight, rectangle.TopLeft, other.BottomRight, other.TopRight) {
		return true
	}
	if lineIntersects(rectangle.TopRight, rectangle.TopLeft, other.TopRight, other.TopLeft) {
		return true
	}
	if lineIntersects(rectangle.TopRight, rectangle.TopLeft, other.TopLeft, other.BottomLeft) {
		return true
	}

	if lineIntersects(rectangle.TopLeft, rectangle.BottomLeft, other.BottomLeft, other.BottomRight) {
		return true
	}
	if lineIntersects(rectangle.TopLeft, rectangle.BottomLeft, other.BottomRight, other.TopRight) {
		return true
	}
	if lineIntersects(rectangle.TopLeft, rectangle.BottomLeft, other.TopRight, other.TopLeft) {
		return true
	}
	if lineIntersects(rectangle.TopLeft, rectangle.BottomLeft, other.TopLeft, other.BottomLeft) {
		return true
	}

	return false
}

// PathIntersectsWithRect reports whether any corner of the rectangle is inside the path.
func PathIntersectsWithRect(path *PdfPath, rectangle core.PdfRectangle, includeBorder bool) bool {
	clipperPaths := make([][]clipperlibrary.ClipperIntPoint, len(path.Subpaths()))
	for i, sp := range path.Subpaths() {
		polygons := SubpathToClipperPolygon(sp)
		clipperPaths[i] = polygons
	}

	fillType := clipperlibrary.ClipperEvenOdd
	if path.FillingRule() == core.FillingRuleNonZeroWinding {
		fillType = clipperlibrary.ClipperNonZero
	}

	for _, pt := range RectangleToClipperPolygon(rectangle) {
		if pointInPaths(pt, clipperPaths, fillType, includeBorder) {
			return true
		}
	}

	return false
}

// RectangleIntersect returns the intersection of two axis-aligned rectangles, or nil if they don't overlap.
func RectangleIntersect(rectangle, other core.PdfRectangle) *core.PdfRectangle {
	if !RectangleIntersectsWithRect(rectangle, other) {
		return nil
	}
	rect := core.NewPdfRectangleFloat(
		math.Max(rectangle.BottomLeft.X, other.BottomLeft.X),
		math.Max(rectangle.BottomLeft.Y, other.BottomLeft.Y),
		math.Min(rectangle.TopRight.X, other.TopRight.X),
		math.Min(rectangle.TopRight.Y, other.TopRight.Y),
	)
	return &rect
}

// RectangleNormalise returns the axis-aligned bounding box of the rectangle with no rotation.
func RectangleNormalise(rectangle core.PdfRectangle) core.PdfRectangle {
	minX := math.Min(math.Min(rectangle.BottomLeft.X, rectangle.BottomRight.X), math.Min(rectangle.TopLeft.X, rectangle.TopRight.X))
	minY := math.Min(math.Min(rectangle.BottomLeft.Y, rectangle.BottomRight.Y), math.Min(rectangle.TopLeft.Y, rectangle.TopRight.Y))
	maxX := math.Max(math.Max(rectangle.BottomLeft.X, rectangle.BottomRight.X), math.Max(rectangle.TopLeft.X, rectangle.TopRight.X))
	maxY := math.Max(math.Max(rectangle.BottomLeft.Y, rectangle.BottomRight.Y), math.Max(rectangle.TopLeft.Y, rectangle.TopRight.Y))

	return core.NewPdfRectangleFloat(minX, minY, maxX, maxY)
}

// RectangleIntersectsWithLine reports whether the rectangle and the line intersect.
func RectangleIntersectsWithLine(rectangle core.PdfRectangle, line core.PdfLine) bool {
	return rectLineIntersects(rectangle, line.Point1, line.Point2)
}

// RectangleIntersectLine returns the intersection of the rectangle and the line as a PdfLine, or nil.
func RectangleIntersectLine(rectangle core.PdfRectangle, line core.PdfLine) *core.PdfLine {
	intersection := rectLineIntersectPoints(rectangle, line.Point1, line.Point2)
	if len(intersection) == 0 {
		return nil
	}
	l := core.NewPdfLine(intersection[0], intersection[1])
	return &l
}

// RectangleIntersectLines returns the list of PdfLines that are the intersection of the rectangle and the lines.
func RectangleIntersectLines(rectangle core.PdfRectangle, lines []core.PdfLine) []core.PdfLine {
	clipper := clipperlibrary.NewClipper(0)
	clipper.AddPath(RectangleToClipperPolygon(rectangle), clipperlibrary.ClipperClip, true)

	for _, line := range lines {
		clipper.AddPath(PdfLineToClipperInt(line), clipperlibrary.ClipperSubject, false)
	}

	solutions := clipperlibrary.NewClipperPolyTree()
	if !clipper.ExecutePolyTree(clipperlibrary.ClipperIntersection, solutions, clipperlibrary.ClipperEvenOdd, clipperlibrary.ClipperEvenOdd) {
		return []core.PdfLine{}
	}

	rv := make([]core.PdfLine, 0)
	for i := 0; i < solutions.ChildCount(); i++ {
		sol := solutions.Children[i]
		contour := sol.Contour()
		l := core.NewPdfLine(
			core.NewPdfPoint(float64(contour[0].X)/Factor, float64(contour[0].Y)/Factor),
			core.NewPdfPoint(float64(contour[1].X)/Factor, float64(contour[1].Y)/Factor),
		)
		rv = append(rv, l)
	}
	return rv
}

// --- PdfLine extensions ---

// LineContainsPoint reports whether the point lies on the line segment.
func LineContainsPoint(line core.PdfLine, point core.PdfPoint) bool {
	return containsOnLine(line.Point1, line.Point2, point)
}

// PdfLineIntersectsWithPdfLine reports whether two PdfLines intersect.
func PdfLineIntersectsWithPdfLine(line, other core.PdfLine) bool {
	return lineIntersects(line.Point1, line.Point2, other.Point1, other.Point2)
}

// PdfLineIntersectsWithCoreLine reports whether a PdfLine and a core.Line intersect.
func PdfLineIntersectsWithCoreLine(line core.PdfLine, other *core.Line) bool {
	return lineIntersects(line.Point1, line.Point2, other.From, other.To)
}

// PdfLineIntersectPdfLine returns the intersection point of two lines, or nil if they don't intersect.
func PdfLineIntersectPdfLine(line, other core.PdfLine) *core.PdfPoint {
	p := lineIntersectPoint(line.Point1, line.Point2, other.Point1, other.Point2)
	return p
}

// PdfLineIntersectCoreLine returns the intersection point of a PdfLine and a core.Line, or nil.
func PdfLineIntersectCoreLine(line core.PdfLine, other *core.Line) *core.PdfPoint {
	p := lineIntersectPoint(line.Point1, line.Point2, other.From, other.To)
	return p
}

// PdfLineParallelTo reports whether two PdfLines are parallel.
func PdfLineParallelTo(line, other core.PdfLine) bool {
	return linesParallel(line.Point1, line.Point2, other.Point1, other.Point2)
}

// PdfLineParallelToCoreLine reports whether a PdfLine and a core.Line are parallel.
func PdfLineParallelToCoreLine(line core.PdfLine, other *core.Line) bool {
	return linesParallel(line.Point1, line.Point2, other.From, other.To)
}

// LineIntersectRect returns the intersection of the line and the rectangle as a PdfLine, or nil.
func LineIntersectRect(line core.PdfLine, rectangle core.PdfRectangle) *core.PdfLine {
	return RectangleIntersectLine(rectangle, line)
}

// LineIntersectsWithRect reports whether the line and the rectangle intersect.
func LineIntersectsWithRect(line core.PdfLine, rectangle core.PdfRectangle) bool {
	return RectangleIntersectsWithLine(rectangle, line)
}

// --- Core.Line extensions ---

// CoreLineContainsPoint reports whether the point lies on the core.Line segment.
func CoreLineContainsPoint(line *core.Line, point core.PdfPoint) bool {
	return containsOnLine(line.From, line.To, point)
}

// CoreLineIntersectsWithCoreLine reports whether two core.Lines intersect.
func CoreLineIntersectsWithCoreLine(line, other *core.Line) bool {
	return lineIntersects(line.From, line.To, other.From, other.To)
}

// CoreLineIntersectsWithPdfLine reports whether a core.Line and a PdfLine intersect.
func CoreLineIntersectsWithPdfLine(line *core.Line, other core.PdfLine) bool {
	return lineIntersects(line.From, line.To, other.Point1, other.Point2)
}

// CoreLineIntersectCoreLine returns the intersection point of two core.Lines, or nil.
func CoreLineIntersectCoreLine(line, other *core.Line) *core.PdfPoint {
	p := lineIntersectPoint(line.From, line.To, other.From, other.To)
	return p
}

// CoreLineIntersectPdfLine returns the intersection point of a core.Line and PdfLine, or nil.
func CoreLineIntersectPdfLine(line *core.Line, other core.PdfLine) *core.PdfPoint {
	p := lineIntersectPoint(line.From, line.To, other.Point1, other.Point2)
	return p
}

// CoreLineParallelTo reports whether two core.Lines are parallel.
func CoreLineParallelTo(line, other *core.Line) bool {
	return linesParallel(line.From, line.To, other.From, other.To)
}

// CoreLineParallelToPdfLine reports whether a core.Line and PdfLine are parallel.
func CoreLineParallelToPdfLine(line *core.Line, other core.PdfLine) bool {
	return linesParallel(line.From, line.To, other.Point1, other.Point2)
}

// --- Generic line helpers ---

func containsOnLine(p1, p2, point core.PdfPoint) bool {
	if math.Abs(p2.X-p1.X) < epsilon {
		if math.Abs(point.X-p2.X) < epsilon {
			signDiff := math.Signbit(point.Y-p2.Y) != math.Signbit(point.Y-p1.Y)
			if signDiff || (point.Y >= math.Min(p1.Y, p2.Y) && point.Y <= math.Max(p1.Y, p2.Y)) {
				return true
			}
		}
		return false
	}

	if math.Abs(p2.Y-p1.Y) < epsilon {
		if math.Abs(point.Y-p2.Y) < epsilon {
			signDiff := math.Signbit(point.X-p2.X) != math.Signbit(point.X-p1.X)
			if signDiff || (point.X >= math.Min(p1.X, p2.X) && point.X <= math.Max(p1.X, p2.X)) {
				return true
			}
		}
		return false
	}

	tx := (point.X - p1.X) / (p2.X - p1.X)
	ty := (point.Y - p1.Y) / (p2.Y - p1.Y)
	if math.Abs(tx-ty) > epsilon {
		return false
	}
	return tx >= 0 && (tx-1) <= epsilon
}

// LineIntersects reports whether the segment (p11,p12) intersects the segment (p21,p22).
func LineIntersects(p11, p12, p21, p22 core.PdfPoint) bool {
	return ccw(p11, p12, p21) != ccw(p11, p12, p22) &&
		ccw(p21, p22, p11) != ccw(p21, p22, p12)
}

func lineIntersects(p11, p12, p21, p22 core.PdfPoint) bool {
	return LineIntersects(p11, p12, p21, p22)
}

func lineIntersectPoint(p11, p12, p21, p22 core.PdfPoint) *core.PdfPoint {
	if !LineIntersects(p11, p12, p21, p22) {
		return nil
	}

	slope1, intercept1 := getSlopeIntercept(p11, p12)
	slope2, intercept2 := getSlopeIntercept(p21, p22)

	if math.IsNaN(slope1) {
		x := intercept1
		y := slope2*x + intercept2
		p := core.NewPdfPoint(x, y)
		return &p
	} else if math.IsNaN(slope2) {
		x := intercept2
		y := slope1*x + intercept1
		p := core.NewPdfPoint(x, y)
		return &p
	} else {
		x := (intercept2 - intercept1) / (slope1 - slope2)
		y := slope1*x + intercept1
		p := core.NewPdfPoint(x, y)
		return &p
	}
}

func rectLineIntersects(rectangle core.PdfRectangle, p1, p2 core.PdfPoint) bool {
	return len(rectLineIntersectPoints(rectangle, p1, p2)) > 0
}

func rectLineIntersectPoints(rectangle core.PdfRectangle, p1, p2 core.PdfPoint) []core.PdfPoint {
	clipper := clipperlibrary.NewClipper(0)
	clipper.AddPath(RectangleToClipperPolygon(rectangle), clipperlibrary.ClipperClip, true)
	clipper.AddPath([]clipperlibrary.ClipperIntPoint{
		PointToClipperInt(p1),
		PointToClipperInt(p2),
	}, clipperlibrary.ClipperSubject, false)

	solutions := clipperlibrary.NewClipperPolyTree()
	if !clipper.ExecutePolyTree(clipperlibrary.ClipperIntersection, solutions, clipperlibrary.ClipperEvenOdd, clipperlibrary.ClipperEvenOdd) {
		return nil
	}

	childCount := solutions.ChildCount()
	if childCount == 0 {
		return nil
	}

	if childCount == 1 {
		sol := solutions.Children[0]
		contour := sol.Contour()
		return []core.PdfPoint{
			core.NewPdfPoint(float64(contour[0].X)/Factor, float64(contour[0].Y)/Factor),
			core.NewPdfPoint(float64(contour[1].X)/Factor, float64(contour[1].Y)/Factor),
		}
	}

	return nil
}

// LineIntersectsRect reports whether the segment (p1,p2) intersects the rectangle.
func LineIntersectsRect(rectangle core.PdfRectangle, p1, p2 core.PdfPoint) bool {
	clipper := clipperlibrary.NewClipper(0)
	clipper.AddPath(RectangleToClipperPolygon(rectangle), clipperlibrary.ClipperClip, true)
	clipper.AddPath([]clipperlibrary.ClipperIntPoint{
		PointToClipperInt(p1),
		PointToClipperInt(p2),
	}, clipperlibrary.ClipperSubject, false)

	solutions := clipperlibrary.NewClipperPolyTree()
	if !clipper.ExecutePolyTree(clipperlibrary.ClipperIntersection, solutions, clipperlibrary.ClipperEvenOdd, clipperlibrary.ClipperEvenOdd) {
		return false
	}

	return solutions.ChildCount() > 0
}

func linesParallel(p11, p12, p21, p22 core.PdfPoint) bool {
	return math.Abs((p12.Y-p11.Y)*(p22.X-p21.X)-(p22.Y-p21.Y)*(p12.X-p11.X)) < epsilon
}

// --- CubicBezierCurve extensions ---

const oneThird = 0.333333333333333333333
const sqrtOfThree = 1.73205080756888

// CubicBezierSplit splits a cubic bezier curve into two at parameter tau using De Casteljau's algorithm.
func CubicBezierSplit(curve *core.CubicBezierCurve, tau float64) (*core.CubicBezierCurve, *core.CubicBezierCurve) {
	points := make([][]core.PdfPoint, 4)

	points[0] = []core.PdfPoint{
		curve.StartPoint,
		curve.FirstControlPoint,
		curve.SecondControlPoint,
		curve.EndPoint,
	}

	points[1] = make([]core.PdfPoint, 3)
	points[2] = make([]core.PdfPoint, 2)
	points[3] = make([]core.PdfPoint, 1)

	for j := 1; j <= 3; j++ {
		for i := 0; i <= 3-j; i++ {
			x := (1-tau)*points[j-1][i].X + tau*points[j-1][i+1].X
			y := (1-tau)*points[j-1][i].Y + tau*points[j-1][i+1].Y
			points[j][i] = core.NewPdfPoint(x, y)
		}
	}

	return core.NewCubicBezierCurve(points[0][0], points[1][0], points[2][0], points[3][0]),
		core.NewCubicBezierCurve(points[3][0], points[2][1], points[1][2], points[0][3])
}

// CubicBezierIntersectsWithPdfLine reports whether the bezier curve and PdfLine intersect.
func CubicBezierIntersectsWithPdfLine(curve *core.CubicBezierCurve, line core.PdfLine) bool {
	return cubicBezierIntersects(curve, line.Point1, line.Point2)
}

// CubicBezierIntersectsWithCoreLine reports whether the bezier curve and core.Line intersect.
func CubicBezierIntersectsWithCoreLine(curve *core.CubicBezierCurve, line *core.Line) bool {
	return cubicBezierIntersects(curve, line.From, line.To)
}

func cubicBezierIntersects(curve *core.CubicBezierCurve, p1, p2 core.PdfPoint) bool {
	ts := cubicBezierIntersectT(curve, p1, p2)
	return ts != nil && len(ts) > 0
}

// CubicBezierIntersectPdfLine returns the intersection points of the curve and PdfLine.
func CubicBezierIntersectPdfLine(curve *core.CubicBezierCurve, line core.PdfLine) []core.PdfPoint {
	return cubicBezierIntersect(curve, line.Point1, line.Point2)
}

// CubicBezierIntersectCoreLine returns the intersection points of the curve and core.Line.
func CubicBezierIntersectCoreLine(curve *core.CubicBezierCurve, line *core.Line) []core.PdfPoint {
	return cubicBezierIntersect(curve, line.From, line.To)
}

func cubicBezierIntersect(curve *core.CubicBezierCurve, p1, p2 core.PdfPoint) []core.PdfPoint {
	ts := cubicBezierIntersectT(curve, p1, p2)
	if ts == nil || len(ts) == 0 {
		return []core.PdfPoint{}
	}

	points := make([]core.PdfPoint, 0, len(ts))
	for _, t := range ts {
		point := core.NewPdfPoint(
			core.BezierValueCubic(curve.StartPoint.X, curve.FirstControlPoint.X, curve.SecondControlPoint.X, curve.EndPoint.X, t),
			core.BezierValueCubic(curve.StartPoint.Y, curve.FirstControlPoint.Y, curve.SecondControlPoint.Y, curve.EndPoint.Y, t),
		)
		if containsOnLine(p1, p2, point) {
			points = append(points, point)
		}
	}
	return points
}

// CubicBezierIntersectTPdfLine returns the t values where the curve and PdfLine intersect.
func CubicBezierIntersectTPdfLine(curve *core.CubicBezierCurve, line core.PdfLine) []float64 {
	return cubicBezierIntersectT(curve, line.Point1, line.Point2)
}

// CubicBezierIntersectTCoreLine returns the t values where the curve and core.Line intersect.
func CubicBezierIntersectTCoreLine(curve *core.CubicBezierCurve, line *core.Line) []float64 {
	return cubicBezierIntersectT(curve, line.From, line.To)
}

func cubicBezierIntersectT(curve *core.CubicBezierCurve, p1, p2 core.PdfPoint) []float64 {
	bezierBbox := curve.GetBoundingRectangle()
	if bezierBbox == nil {
		return nil
	}

	bb := *bezierBbox

	if bb.Left() > math.Max(p1.X, p2.X) || math.Min(p1.X, p2.X) > bb.Right() {
		return nil
	}

	if bb.Top() < math.Min(p1.Y, p2.Y) || math.Max(p1.Y, p2.Y) < bb.Bottom() {
		return nil
	}

	aCoeff := p2.Y - p1.Y
	bCoeff := p1.X - p2.X
	cCoeff := p1.X*(p1.Y-p2.Y) + p1.Y*(p2.X-p1.X)

	alpha := curve.StartPoint.X*aCoeff + curve.StartPoint.Y*bCoeff
	beta := 3.0 * (curve.FirstControlPoint.X*aCoeff + curve.FirstControlPoint.Y*bCoeff)
	gamma := 3.0 * (curve.SecondControlPoint.X*aCoeff + curve.SecondControlPoint.Y*bCoeff)
	delta := curve.EndPoint.X*aCoeff + curve.EndPoint.Y*bCoeff

	a := -alpha + beta - gamma + delta
	b := 3*alpha - 2*beta + gamma
	c := -3*alpha + beta
	d := alpha + cCoeff

	solution := solveCubicEquation(a, b, c, d)

	var result []float64
	for _, s := range solution {
		if !math.IsNaN(s) && s >= -epsilon && (s-1) <= epsilon {
			result = append(result, s)
		}
	}
	slices.Sort(result)
	return result
}

// --- Clipper helpers ---

func crossProduct(pt1, pt2, pt3 clipperlibrary.ClipperIntPoint) float64 {
	return float64((pt2.X-pt1.X)*(pt3.Y-pt2.Y)-(pt2.Y-pt1.Y)*(pt3.X-pt2.X))
}

// pointInPathsWindingCount computes the winding count of a point relative to paths.
func pointInPathsWindingCount(pt clipperlibrary.ClipperIntPoint, paths [][]clipperlibrary.ClipperIntPoint) int {
	result := 0
	for _, path := range paths {
		j := 0
		p := path
		length := len(p)

		if length < 3 {
			continue
		}
		prevPt := p[length-1]

		for j < length && p[j].Y == prevPt.Y {
			j++
		}
		if j == length {
			continue
		}

		isAbove := prevPt.Y < pt.Y

		for j < length {
			if isAbove {
				for j < length && p[j].Y < pt.Y {
					j++
				}
				if j == length {
					break
				}
				if j > 0 {
					prevPt = p[j-1]
				}

				crossProd := crossProduct(prevPt, p[j], pt)
				if crossProd == 0 {
					return math.MaxInt32
				} else if crossProd < 0 {
					result--
				}
			} else {
				for j < length && p[j].Y > pt.Y {
					j++
				}
				if j == length {
					break
				}
				if j > 0 {
					prevPt = p[j-1]
				}

				crossProd := crossProduct(prevPt, p[j], pt)
				if crossProd == 0 {
					return math.MaxInt32
				} else if crossProd > 0 {
					result++
				}
			}

			j++
			isAbove = !isAbove
		}
	}
	return result
}

func pointInPaths(pt clipperlibrary.ClipperIntPoint, paths [][]clipperlibrary.ClipperIntPoint, fillRule clipperlibrary.ClipperPolyFillType, includeBorder bool) bool {
	wc := pointInPathsWindingCount(pt, paths)
	if wc == math.MaxInt32 {
		return includeBorder
	}

	switch fillRule {
	default:
	case clipperlibrary.ClipperEvenOdd:
		return wc%2 != 0
	case clipperlibrary.ClipperNonZero:
		return wc != 0
	}
	return false
}

// --- PdfSubpath extensions ---

// SubpathContainsPoint reports whether the point is inside the subpath (ignoring winding rule).
func SubpathContainsPoint(subpath *core.PdfSubpath, point core.PdfPoint, includeBorder bool) bool {
	return pointInPaths(
		PointToClipperInt(point),
		[][]clipperlibrary.ClipperIntPoint{SubpathToClipperPolygon(subpath)},
		clipperlibrary.ClipperEvenOdd,
		includeBorder,
	)
}

// SubpathContainsRect reports whether the rectangle is inside the subpath (ignoring winding rule).
func SubpathContainsRect(subpath *core.PdfSubpath, rectangle core.PdfRectangle, includeBorder bool) bool {
	clipperPaths := [][]clipperlibrary.ClipperIntPoint{SubpathToClipperPolygon(subpath)}
	for _, pt := range RectangleToClipperPolygon(rectangle) {
		if !pointInPaths(pt, clipperPaths, clipperlibrary.ClipperEvenOdd, includeBorder) {
			return false
		}
	}
	return true
}

// SubpathContainsSubpath reports whether the other subpath is inside this one (ignoring winding rule).
func SubpathContainsSubpath(subpath, other *core.PdfSubpath, includeBorder bool) bool {
	clipperPaths := [][]clipperlibrary.ClipperIntPoint{SubpathToClipperPolygon(subpath)}
	for _, pt := range SubpathToClipperPolygon(other) {
		if !pointInPaths(pt, clipperPaths, clipperlibrary.ClipperEvenOdd, includeBorder) {
			return false
		}
	}
	return true
}

// --- PdfPath extensions ---

// PathGetArea returns the total area of the path.
func PathGetArea(path *PdfPath) float64 {
	clipperPaths := make([][]clipperlibrary.ClipperIntPoint, len(path.Subpaths()))
	for i, sp := range path.Subpaths() {
		clipperPaths[i] = SubpathToClipperPolygon(sp)
	}

	fillType := clipperlibrary.ClipperEvenOdd
	if path.FillingRule() == core.FillingRuleNonZeroWinding {
		fillType = clipperlibrary.ClipperNonZero
	}

	simplifieds := clipperlibrary.SimplifyPolygons(clipperPaths, fillType)
	sum := 0.0
	for _, simplified := range simplifieds {
		sum += clipperlibrary.AreaPoly(simplified)
	}
	return sum / (Factor * Factor)
}

// PathContainsPoint reports whether the point is inside the path.
func PathContainsPoint(path *PdfPath, point core.PdfPoint, includeBorder bool) bool {
	clipperPaths := make([][]clipperlibrary.ClipperIntPoint, len(path.Subpaths()))
	for i, sp := range path.Subpaths() {
		clipperPaths[i] = SubpathToClipperPolygon(sp)
	}

	fillType := clipperlibrary.ClipperEvenOdd
	if path.FillingRule() == core.FillingRuleNonZeroWinding {
		fillType = clipperlibrary.ClipperNonZero
	}

	return pointInPaths(PointToClipperInt(point), clipperPaths, fillType, includeBorder)
}

// PathContainsRect reports whether the rectangle is inside the path.
func PathContainsRect(path *PdfPath, rectangle core.PdfRectangle, includeBorder bool) bool {
	clipperPaths := make([][]clipperlibrary.ClipperIntPoint, len(path.Subpaths()))
	for i, sp := range path.Subpaths() {
		clipperPaths[i] = SubpathToClipperPolygon(sp)
	}

	fillType := clipperlibrary.ClipperEvenOdd
	if path.FillingRule() == core.FillingRuleNonZeroWinding {
		fillType = clipperlibrary.ClipperNonZero
	}

	for _, pt := range RectangleToClipperPolygon(rectangle) {
		if !pointInPaths(pt, clipperPaths, fillType, includeBorder) {
			return false
		}
	}
	return true
}

// PathContainsSubpath reports whether the subpath is inside the path.
func PathContainsSubpath(path *PdfPath, subpath *core.PdfSubpath, includeBorder bool) bool {
	clipperPaths := make([][]clipperlibrary.ClipperIntPoint, len(path.Subpaths()))
	for i, sp := range path.Subpaths() {
		clipperPaths[i] = SubpathToClipperPolygon(sp)
	}

	fillType := clipperlibrary.ClipperEvenOdd
	if path.FillingRule() == core.FillingRuleNonZeroWinding {
		fillType = clipperlibrary.ClipperNonZero
	}

	for _, pt := range SubpathToClipperPolygon(subpath) {
		if !pointInPaths(pt, clipperPaths, fillType, includeBorder) {
			return false
		}
	}
	return true
}

// PathContainsPath reports whether the other path is inside this one.
func PathContainsPath(path, other *PdfPath, includeBorder bool) bool {
	clipperPaths := make([][]clipperlibrary.ClipperIntPoint, len(path.Subpaths()))
	for i, sp := range path.Subpaths() {
		clipperPaths[i] = SubpathToClipperPolygon(sp)
	}

	fillType := clipperlibrary.ClipperEvenOdd
	if path.FillingRule() == core.FillingRuleNonZeroWinding {
		fillType = clipperlibrary.ClipperNonZero
	}

	for _, subpath := range other.Subpaths() {
		for _, pt := range SubpathToClipperPolygon(subpath) {
			if !pointInPaths(pt, clipperPaths, fillType, includeBorder) {
				return false
			}
		}
	}
	return true
}

// --- SVG exports ---

// SubpathToSvg converts a subpath to an SVG path data string.
func SubpathToSvg(subpath *core.PdfSubpath, height float64) string {
	var builder strings.Builder
	for _, cmd := range subpath.Commands() {
		switch c := cmd.(type) {
		case *core.Move:
			fmt.Fprintf(&builder, "M %g %g ", c.Location.X, height-c.Location.Y)
		case *core.Line:
			fmt.Fprintf(&builder, "L %g %g ", c.To.X, height-c.To.Y)
		case *core.QuadraticBezierCurve:
			fmt.Fprintf(&builder, "Q %g %g, %g %g ",
				c.ControlPoint.X, height-c.ControlPoint.Y,
				c.EndPoint.X, height-c.EndPoint.Y)
		case *core.CubicBezierCurve:
			fmt.Fprintf(&builder, "C %g %g, %g %g, %g %g ",
				c.FirstControlPoint.X, height-c.FirstControlPoint.Y,
				c.SecondControlPoint.X, height-c.SecondControlPoint.Y,
				c.EndPoint.X, height-c.EndPoint.Y)
		case *core.Close:
			builder.WriteString("Z ")
		}
	}

	s := builder.String()
	if len(s) == 0 {
		return ""
	}
	if s[len(s)-1] == ' ' {
		s = s[:len(s)-1]
	}
	return s
}

// SubpathToFullSvg converts a subpath to a full SVG document with bounding boxes.
func SubpathToFullSvg(subpath *core.PdfSubpath, height float64) string {
	bboxToRect := func(box core.PdfRectangle, stroke string) string {
		return fmt.Sprintf("<rect x='%g' y='%g' width='%g' height='%g' stroke-width='2' fill='none' stroke='%s'></rect>",
			box.Left(), box.Bottom(), box.Width, box.Height, stroke)
	}

	glyph := SubpathToSvg(subpath, height)
	bbox := subpath.GetBoundingRectangle()
	var bboxes []core.PdfRectangle

	for _, cmd := range subpath.Commands() {
		segBbox := cmd.GetBoundingRectangle()
		if segBbox != nil {
			bboxes = append(bboxes, *segBbox)
		}
	}

	path := fmt.Sprintf("<path d='%s' stroke='cyan' stroke-width='3'></path>", glyph)
	var bboxRect string
	if bbox != nil {
		bboxRect = bboxToRect(*bbox, "yellow")
	}

	othersParts := make([]string, len(bboxes))
	for i, b := range bboxes {
		othersParts[i] = bboxToRect(b, "gray")
	}
	others := strings.Join(othersParts, " ")

	return fmt.Sprintf("<svg width='500' height='500'><g transform=\"scale(0.2, -0.2) translate(100, -700)\">%s %s %s</g></svg>",
		path, bboxRect, others)
}

// --- Private helpers ---

func getSlopeIntercept(p1, p2 core.PdfPoint) (float64, float64) {
	if math.Abs(p1.X-p2.X) > epsilon {
		slope := (p2.Y - p1.Y) / (p2.X - p1.X)
		intercept := p2.Y - slope*p2.X
		return slope, intercept
	}
	return math.NaN(), p1.X
}

func cubicRoot(d float64) float64 {
	if d < 0.0 {
		return -math.Pow(-d, oneThird)
	}
	return math.Pow(d, oneThird)
}

// solveCubicEquation solves ax^3 + bx^2 + cx + d = 0 and returns up to 3 real roots (NaN for missing).
func solveCubicEquation(a, b, c, d float64) []float64 {
	if math.Abs(a) <= epsilon {
		detQ := c*c - 4*b*d
		if detQ >= 0 {
			sqrtDetQ := math.Sqrt(detQ)
			oneOverTwiceB := 1 / (2.0 * b)
			x := (-c + sqrtDetQ) * oneOverTwiceB
			x0 := (-c - sqrtDetQ) * oneOverTwiceB
			return []float64{x, x0, math.NaN()}
		}
		return []float64{math.NaN(), math.NaN(), math.NaN()}
	}

	aSquared := a * a
	aCubed := aSquared * a
	bCubed := b * b * b
	abc := a * b * c
	bOver3a := b / (3.0 * a)

	Q := (3.0*a*c - b*b) / (9.0 * aSquared)
	R := (9.0*abc - 27.0*aSquared*d - 2.0*bCubed) / (54.0 * aCubed)

	det := Q*Q*Q + R*R

	x1 := math.NaN()
	x2 := math.NaN()
	x3 := math.NaN()

	if det >= 0 {
		sqrtDet := math.Sqrt(det)

		S := cubicRoot(R + sqrtDet)
		T := cubicRoot(R - sqrtDet)
		sPlusT := S + T

		x1 = sPlusT - bOver3a

		complexPart := sqrtOfThree / 2.0 * (S - T)
		if math.Abs(complexPart) <= epsilon {
			x2 = -sPlusT/2 - bOver3a
			x3 = x2
		}
	} else {
		viet := func(p_, q_ float64, k int) float64 {
			return 2.0 * math.Sqrt(-p_/3.0) *
				math.Cos(oneThird*math.Acos((3.0*q_)/(2.0*p_)*math.Sqrt(-3.0/p_))-(2.0*math.Pi*float64(k))/3.0)
		}

		p := Q * 3.0
		q := -R * 2.0
		x1 = viet(p, q, 0) - bOver3a
		x2 = viet(p, q, 1) - bOver3a
		x3 = viet(p, q, 2) - bOver3a
	}

	return []float64{x1, x2, x3}
}
