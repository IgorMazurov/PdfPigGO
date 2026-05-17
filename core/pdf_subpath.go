package core

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// PdfSubpath is made up of a sequence of connected segments.
type PdfSubpath struct {
	commands           []PathCommand
	currentPosition    *PdfPoint
	shoeLaceSum        float64
	isDrawnAsRectangle bool
}

// NewPdfSubpath creates a new empty PdfSubpath.
func NewPdfSubpath() *PdfSubpath {
	return &PdfSubpath{
		commands: make([]PathCommand, 0),
	}
}

// Commands returns the sequence of commands which form this PdfSubpath.
func (s *PdfSubpath) Commands() []PathCommand {
	return s.commands
}

// IsDrawnAsRectangle reports true if the PdfSubpath was originally drawn using
// the rectangle ('re') operator. Always false if paths are clipped.
func (s *PdfSubpath) IsDrawnAsRectangle() bool {
	return s.isDrawnAsRectangle
}

// SetIsDrawnAsRectangle sets whether this subpath was drawn as a rectangle.
func (s *PdfSubpath) SetIsDrawnAsRectangle(v bool) {
	s.isDrawnAsRectangle = v
}

// IsClockwise reports true if points are organised in a clockwise order.
// Works only with closed paths.
func (s *PdfSubpath) IsClockwise() bool {
	return s.IsClosed() && s.shoeLaceSum > 0
}

// IsCounterClockwise reports true if points are organised in a counterclockwise order.
// Works only with closed paths.
func (s *PdfSubpath) IsCounterClockwise() bool {
	return s.IsClosed() && s.shoeLaceSum < 0
}

// GetCentroid returns the PdfSubpath's centroid point.
func (s *PdfSubpath) GetCentroid() PdfPoint {
	var filtered []PathCommand
	for _, c := range s.commands {
		if _, ok := c.(*Line); ok {
			filtered = append(filtered, c)
		} else if isBezier(c) {
			filtered = append(filtered, c)
		}
	}

	if len(filtered) == 0 {
		return PdfPoint{}
	}

	var sumX, sumY float64
	for _, cmd := range filtered {
		startPt := GetStartPoint(cmd)
		endPt := GetEndPoint(cmd)
		sumX += startPt.X + endPt.X
		sumY += startPt.Y + endPt.Y
	}

	count := float64(len(filtered)) * 2.0
	return NewPdfPoint(sumX/count, sumY/count)
}

// Simplify converts everything to PdfLines.
// n is the number of lines required (minimum is 1).
func (s *PdfSubpath) Simplify(n int) *PdfSubpath {
	if n < 1 {
		n = 4
	}

	simplifiedPath := NewPdfSubpath()
	startPoint := GetStartPoint(s.commands[0])
	simplifiedPath.MoveTo(startPoint.X, startPoint.Y)

	for _, command := range s.commands {
		switch cmd := command.(type) {
		case *Line:
			simplifiedPath.LineTo(cmd.To.X, cmd.To.Y)
		default:
			if bz, ok := toBezier(command); ok {
				lines := bz.ToLines(n)
				for _, lineB := range lines {
					simplifiedPath.LineTo(lineB.To.X, lineB.To.Y)
				}
			}
		}
	}

	if s.IsClosed() && len(simplifiedPath.commands) > 0 {
		first := GetStartPoint(simplifiedPath.commands[0])
		lastEnd := GetEndPoint(simplifiedPath.commands[len(simplifiedPath.commands)-1])
		if !first.Equals(lastEnd) {
			simplifiedPath.LineTo(first.X, first.Y)
		}
	}

	return simplifiedPath
}

// MoveTo adds a Move command to the path.
func (s *PdfSubpath) MoveTo(x, y float64) {
	pos := NewPdfPoint(x, y)
	s.currentPosition = &pos
	s.commands = append(s.commands, &Move{Location: pos})
}

// LineTo adds a Line command to the path.
func (s *PdfSubpath) LineTo(x, y float64) error {
	if s.currentPosition == nil {
		return errors.New("LineTo(): currentPosition is null")
	}

	s.shoeLaceSum += (x - s.currentPosition.X) * (y + s.currentPosition.Y)

	to := NewPdfPoint(x, y)
	s.commands = append(s.commands, &Line{From: *s.currentPosition, To: to})
	s.currentPosition = &to
	return nil
}

// Rectangle adds a rectangle following the PDF specification (m, l, l, l, c) path.
// A new subpath is created.
func (s *PdfSubpath) Rectangle(x, y, width, height float64) {
	s.MoveTo(x, y)
	s.LineTo(x+width, y)
	s.LineTo(x+width, y+height)
	s.LineTo(x, y+height)
	s.CloseSubpath()
	s.isDrawnAsRectangle = true
}

// BezierCurveToCubic adds a CubicBezierCurve to the path.
func (s *PdfSubpath) BezierCurveToCubic(x1, y1, x2, y2, x3, y3 float64) error {
	if s.currentPosition == nil {
		return errors.New("BezierCurveTo(): currentPosition is null")
	}

	s.shoeLaceSum += (x1 - s.currentPosition.X) * (y1 + s.currentPosition.Y)
	s.shoeLaceSum += (x2 - x1) * (y2 + y1)
	s.shoeLaceSum += (x3 - x2) * (y3 + y2)

	to := NewPdfPoint(x3, y3)
	s.commands = append(s.commands, &CubicBezierCurve{
		StartPoint:         *s.currentPosition,
		EndPoint:           to,
		FirstControlPoint:  NewPdfPoint(x1, y1),
		SecondControlPoint: NewPdfPoint(x2, y2),
	})
	s.currentPosition = &to
	return nil
}

// BezierCurveToQuadratic adds a QuadraticBezierCurve to the path.
// Only used in fonts.
func (s *PdfSubpath) BezierCurveToQuadratic(x1, y1, x2, y2 float64) error {
	if s.currentPosition == nil {
		return errors.New("BezierCurveTo(): currentPosition is null")
	}

	s.shoeLaceSum += (x1 - s.currentPosition.X) * (y1 + s.currentPosition.Y)
	s.shoeLaceSum += (x2 - x1) * (y2 + y1)

	to := NewPdfPoint(x2, y2)
	s.commands = append(s.commands, &QuadraticBezierCurve{
		StartPoint:   *s.currentPosition,
		EndPoint:     to,
		ControlPoint: NewPdfPoint(x1, y1),
	})
	s.currentPosition = &to
	return nil
}

// CloseSubpath closes the path.
func (s *PdfSubpath) CloseSubpath() {
	if s.currentPosition != nil && len(s.commands) > 0 {
		startPoint := GetStartPoint(s.commands[0])
		if !startPoint.Equals(*s.currentPosition) {
			s.shoeLaceSum += (startPoint.X - s.currentPosition.X) * (startPoint.Y + s.currentPosition.Y)
		}
	}
	s.commands = append(s.commands, &Close{})
}

// IsClosed determines if the path is currently closed.
func (s *PdfSubpath) IsClosed() bool {
	var lastCmd, firstCmd PathCommand
	filteredCount := 0

	for i := len(s.commands) - 1; i >= 0; i-- {
		cmd := s.commands[i]

		if _, ok := cmd.(*Close); ok {
			return true
		}

		if isLineOrMove(cmd) || isBezier(cmd) {
			if lastCmd == nil {
				lastCmd = cmd
			}
			firstCmd = cmd
			filteredCount++
		}
	}

	if filteredCount < 2 || lastCmd == nil || firstCmd == nil {
		return false
	}

	if !GetStartPoint(firstCmd).Equals(GetEndPoint(lastCmd)) {
		return false
	}

	return true
}

// GetBoundingRectangle gets a PdfRectangle which entirely contains the geometry
// of the defined subpath. For subpaths which don't define any geometry this returns nil.
func (s *PdfSubpath) GetBoundingRectangle() *PdfRectangle {
	if len(s.commands) == 0 {
		return nil
	}

	minX := math.MaxFloat64
	maxX := -math.MaxFloat64
	minY := math.MaxFloat64
	maxY := -math.MaxFloat64

	for _, command := range s.commands {
		rect := command.GetBoundingRectangle()
		if rect == nil {
			continue
		}

		if rect.Left() < minX {
			minX = rect.Left()
		}
		if rect.Right() > maxX {
			maxX = rect.Right()
		}
		if rect.Bottom() < minY {
			minY = rect.Bottom()
		}
		if rect.Top() > maxY {
			maxY = rect.Top()
		}
	}

	if minX == math.MaxFloat64 || maxX == -math.MaxFloat64 ||
		minY == math.MaxFloat64 || maxY == -math.MaxFloat64 {
		return nil
	}

	rect := NewPdfRectangleFloat(minX, minY, maxX, maxY)
	return &rect
}

// GetDrawnRectangle returns the rectangle dimensions if IsDrawnAsRectangle is true.
// Otherwise returns nil.
func (s *PdfSubpath) GetDrawnRectangle() *PdfRectangle {
	if !s.isDrawnAsRectangle || len(s.commands) != 5 {
		return nil
	}

	mv, ok0 := s.commands[0].(*Move)
	line1, ok1 := s.commands[1].(*Line)
	line2, ok2 := s.commands[2].(*Line)
	_, ok3 := s.commands[3].(*Line)
	_, ok4 := s.commands[4].(*Close)

	if !ok0 || !ok1 || !ok2 || !ok3 || !ok4 {
		return nil
	}

	if !line1.From.Equals(mv.Location) || line1.To.Y != mv.Location.Y {
		return nil
	}

	width := line1.To.X - mv.Location.X

	if !line2.From.Equals(line1.To) || line2.To.X != line1.To.X {
		return nil
	}

	height := line2.To.Y - line1.To.Y

	rect := NewPdfRectangle(mv.Location, NewPdfPoint(mv.Location.X+width, mv.Location.Y+height))
	return &rect
}

// GetBoundingRectangleForPath gets a PdfRectangle which entirely contains the geometry
// of the defined path. For paths which don't define any geometry this returns nil.
func GetBoundingRectangleForPath(path []*PdfSubpath) *PdfRectangle {
	if len(path) == 0 {
		return nil
	}

	var bboxes []PdfRectangle
	for _, subpath := range path {
		rect := subpath.GetBoundingRectangle()
		if rect != nil {
			bboxes = append(bboxes, *rect)
		}
	}

	if len(bboxes) == 0 {
		return nil
	}

	minX := bboxes[0].Left()
	minY := bboxes[0].Bottom()
	maxX := bboxes[0].Right()
	maxY := bboxes[0].Top()

	for _, box := range bboxes[1:] {
		if box.Left() < minX {
			minX = box.Left()
		}
		if box.Bottom() < minY {
			minY = box.Bottom()
		}
		if box.Right() > maxX {
			maxX = box.Right()
		}
		if box.Top() > maxY {
			maxY = box.Top()
		}
	}

	rect := NewPdfRectangleFloat(minX, minY, maxX, maxY)
	return &rect
}

// Equals compares two PdfSubpaths for equality. Paths will only be considered equal
// if the commands which construct the paths are in the same order.
func (s *PdfSubpath) Equals(other *PdfSubpath) bool {
	if other == nil || len(s.commands) != len(other.commands) {
		return false
	}

	for i := range s.commands {
		if !s.commands[i].Equals(other.commands[i]) {
			return false
		}
	}

	return true
}

// PathCommand is a command in a PdfSubpath.
type PathCommand interface {
	// GetBoundingRectangle returns the smallest rectangle which contains the path region given by this command.
	GetBoundingRectangle() *PdfRectangle

	// WriteSvg converts from the path command to an SVG string representing the path operation.
	WriteSvg(builder *strings.Builder, height float64)

	// Equals compares two PathCommands for equality.
	Equals(other PathCommand) bool
}

var _ PathCommand = (*Close)(nil)
var _ PathCommand = (*Move)(nil)
var _ PathCommand = (*Line)(nil)
var _ PathCommand = (*QuadraticBezierCurve)(nil)
var _ PathCommand = (*CubicBezierCurve)(nil)

// Close closes the current PdfSubpath.
type Close struct{}

func (c *Close) GetBoundingRectangle() *PdfRectangle {
	return nil
}

func (c *Close) WriteSvg(builder *strings.Builder, height float64) {
	builder.WriteString("Z ")
}

func (c *Close) Equals(other PathCommand) bool {
	_, ok := other.(*Close)
	return ok
}

// Move moves drawing of the current PdfSubpath to the specified location.
type Move struct {
	Location PdfPoint
}

// NewMove creates a new Move path command.
func NewMove(location PdfPoint) *Move {
	return &Move{Location: location}
}

func (m *Move) GetBoundingRectangle() *PdfRectangle {
	return nil
}

func (m *Move) WriteSvg(builder *strings.Builder, height float64) {
	fmt.Fprintf(builder, "M %g %g ", m.Location.X, height-m.Location.Y)
}

func (m *Move) Equals(other PathCommand) bool {
	if o, ok := other.(*Move); ok {
		return m.Location.Equals(o.Location)
	}
	return false
}

// Line draws a straight line between two points.
type Line struct {
	From PdfPoint
	To   PdfPoint
}

// NewLine creates a new Line.
func NewLine(from, to PdfPoint) *Line {
	return &Line{From: from, To: to}
}

// Length returns the length of the line.
func (l *Line) Length() float64 {
	dx := l.From.X - l.To.X
	dy := l.From.Y - l.To.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (l *Line) GetBoundingRectangle() *PdfRectangle {
	rect := NewPdfRectangle(l.From, l.To)
	return &rect
}

func (l *Line) WriteSvg(builder *strings.Builder, height float64) {
	fmt.Fprintf(builder, "L %g %g ", l.To.X, height-l.To.Y)
}

func (l *Line) Equals(other PathCommand) bool {
	if o, ok := other.(*Line); ok {
		return l.From.Equals(o.From) && l.To.Equals(o.To)
	}
	return false
}

// isLineOrMove reports whether cmd is a Line or Move.
func isLineOrMove(cmd PathCommand) bool {
	_, ok1 := cmd.(*Line)
	_, ok2 := cmd.(*Move)
	return ok1 || ok2
}

// isBezier reports whether cmd is any Bezier curve type.
func isBezier(cmd PathCommand) bool {
	_, ok1 := cmd.(*QuadraticBezierCurve)
	_, ok2 := cmd.(*CubicBezierCurve)
	return ok1 || ok2
}

// toBezier extracts a bezierLines interface from a PathCommand if it's a Bezier curve.
type bezierLines interface {
	ToLines(n int) []*Line
}

func toBezier(cmd PathCommand) (bezierLines, bool) {
	if q, ok := cmd.(*QuadraticBezierCurve); ok {
		return q, true
	}
	if c, ok := cmd.(*CubicBezierCurve); ok {
		return c, true
	}
	return nil, false
}

// QuadraticBezierCurve draws a quadratic Bezier-curve given by the start, control and end points.
// Only used in fonts.
type QuadraticBezierCurve struct {
	StartPoint   PdfPoint
	EndPoint     PdfPoint
	ControlPoint PdfPoint
}

// NewQuadraticBezierCurve creates a quadratic Bezier-curve at the provided points.
func NewQuadraticBezierCurve(startPoint, controlPoint, endPoint PdfPoint) *QuadraticBezierCurve {
	return &QuadraticBezierCurve{
		StartPoint:   startPoint,
		EndPoint:     endPoint,
		ControlPoint: controlPoint,
	}
}

func (q *QuadraticBezierCurve) GetBoundingRectangle() *PdfRectangle {
	minX := q.StartPoint.X
	maxX := q.EndPoint.X
	if minX > maxX {
		minX, maxX = maxX, minX
	}

	minY := q.StartPoint.Y
	maxY := q.EndPoint.Y
	if minY > maxY {
		minY, maxY = maxY, minY
	}

	solved, xsMin, xsMax := q.trySolve(true, minX, maxX)
	if solved {
		minX = xsMin
		maxX = xsMax
	}

	solved, ysMin, ysMax := q.trySolve(false, minY, maxY)
	if solved {
		minY = ysMin
		maxY = ysMax
	}

	rect := NewPdfRectangleFloat(minX, minY, maxX, maxY)
	return &rect
}

func (q *QuadraticBezierCurve) WriteSvg(builder *strings.Builder, height float64) {
	fmt.Fprintf(builder, "C %g %g, %g %g ", q.ControlPoint.X, height-q.ControlPoint.Y, q.EndPoint.X, height-q.EndPoint.Y)
}

func (q *QuadraticBezierCurve) trySolve(isX bool, currentMin, currentMax float64) (bool, float64, float64) {
	p1 := q.StartPoint.X
	p2 := q.ControlPoint.X
	p3 := q.EndPoint.X
	if !isX {
		p1 = q.StartPoint.Y
		p2 = q.ControlPoint.Y
		p3 = q.EndPoint.Y
	}

	t := (p1 - p2) / (p1 - 2.0*p2 + p3)

	if t >= 0 && t <= 1 {
		sol := BezierValueQuad(p1, p2, p3, t)
		if sol < currentMin {
			currentMin = sol
		}
		if sol > currentMax {
			currentMax = sol
		}
	}

	return true, currentMin, currentMax
}

// ToLines converts the quadratic bezier curve into approximated lines.
func (q *QuadraticBezierCurve) ToLines(n int) []*Line {
	lines := make([]*Line, n)
	previousPoint := q.StartPoint

	for p := 1; p <= n; p++ {
		t := float64(p) / float64(n)
		currentPoint := NewPdfPoint(
			BezierValueQuad(q.StartPoint.X, q.ControlPoint.X, q.EndPoint.X, t),
			BezierValueQuad(q.StartPoint.Y, q.ControlPoint.Y, q.EndPoint.Y, t),
		)
		lines[p-1] = NewLine(previousPoint, currentPoint)
		previousPoint = currentPoint
	}

	return lines
}

func (q *QuadraticBezierCurve) Equals(other PathCommand) bool {
	if o, ok := other.(*QuadraticBezierCurve); ok {
		return q.StartPoint.Equals(o.StartPoint) &&
			q.ControlPoint.Equals(o.ControlPoint) &&
			q.EndPoint.Equals(o.EndPoint)
	}
	return false
}

// CubicBezierCurve draws a cubic Bezier-curve given by the start, control and end points.
type CubicBezierCurve struct {
	StartPoint         PdfPoint
	EndPoint           PdfPoint
	FirstControlPoint  PdfPoint
	SecondControlPoint PdfPoint
}

// NewCubicBezierCurve creates a cubic Bezier-curve at the provided points.
func NewCubicBezierCurve(startPoint, firstControlPoint, secondControlPoint, endPoint PdfPoint) *CubicBezierCurve {
	return &CubicBezierCurve{
		StartPoint:         startPoint,
		EndPoint:           endPoint,
		FirstControlPoint:  firstControlPoint,
		SecondControlPoint: secondControlPoint,
	}
}

func (c *CubicBezierCurve) GetBoundingRectangle() *PdfRectangle {
	minX := c.StartPoint.X
	maxX := c.EndPoint.X
	if minX > maxX {
		minX, maxX = maxX, minX
	}

	minY := c.StartPoint.Y
	maxY := c.EndPoint.Y
	if minY > maxY {
		minY, maxY = maxY, minY
	}

	solved, xsMin, xsMax := c.trySolve(true, minX, maxX)
	if solved {
		minX = xsMin
		maxX = xsMax
	}

	solved, ysMin, ysMax := c.trySolve(false, minY, maxY)
	if solved {
		minY = ysMin
		maxY = ysMax
	}

	rect := NewPdfRectangleFloat(minX, minY, maxX, maxY)
	return &rect
}

func (c *CubicBezierCurve) WriteSvg(builder *strings.Builder, height float64) {
	fmt.Fprintf(builder, "C %g %g, %g %g, %g %g ",
		c.FirstControlPoint.X, height-c.FirstControlPoint.Y,
		c.SecondControlPoint.X, height-c.SecondControlPoint.Y,
		c.EndPoint.X, height-c.EndPoint.Y)
}

func (c *CubicBezierCurve) trySolve(isX bool, currentMin, currentMax float64) (bool, float64, float64) {
	p1 := c.StartPoint.X
	p2 := c.FirstControlPoint.X
	p3 := c.SecondControlPoint.X
	p4 := c.EndPoint.X
	if !isX {
		p1 = c.StartPoint.Y
		p2 = c.FirstControlPoint.Y
		p3 = c.SecondControlPoint.Y
		p4 = c.EndPoint.Y
	}

	threeda := 3 * (p2 - p1)
	sixdb := 6 * (p3 - p2)
	threedc := 3 * (p4 - p3)

	a := threeda - sixdb + threedc
	bCoeff := sixdb - threeda - threeda
	ccoeff := threeda

	sqrtable := bCoeff*bCoeff - 4*a*ccoeff

	if sqrtable < 0 {
		return false, currentMin, currentMax
	}

	sqrt := math.Sqrt(sqrtable)
	divisor := 2 * a

	t1 := (-bCoeff + sqrt) / divisor
	t2 := (-bCoeff - sqrt) / divisor

	if t1 >= 0 && t1 <= 1 {
		sol1 := BezierValueCubic(p1, p2, p3, p4, t1)
		if sol1 < currentMin {
			currentMin = sol1
		}
		if sol1 > currentMax {
			currentMax = sol1
		}
	}

	if t2 >= 0 && t2 <= 1 {
		sol2 := BezierValueCubic(p1, p2, p3, p4, t2)
		if sol2 < currentMin {
			currentMin = sol2
		}
		if sol2 > currentMax {
			currentMax = sol2
		}
	}

	return true, currentMin, currentMax
}

// ToLines converts the cubic bezier curve into approximated lines.
func (c *CubicBezierCurve) ToLines(n int) []*Line {
	lines := make([]*Line, n)
	previousPoint := c.StartPoint

	for p := 1; p <= n; p++ {
		t := float64(p) / float64(n)
		currentPoint := NewPdfPoint(
			BezierValueCubic(c.StartPoint.X, c.FirstControlPoint.X, c.SecondControlPoint.X, c.EndPoint.X, t),
			BezierValueCubic(c.StartPoint.Y, c.FirstControlPoint.Y, c.SecondControlPoint.Y, c.EndPoint.Y, t),
		)
		lines[p-1] = NewLine(previousPoint, currentPoint)
		previousPoint = currentPoint
	}

	return lines
}

func (c *CubicBezierCurve) Equals(other PathCommand) bool {
	if o, ok := other.(*CubicBezierCurve); ok {
		return c.StartPoint.Equals(o.StartPoint) &&
			c.FirstControlPoint.Equals(o.FirstControlPoint) &&
			c.SecondControlPoint.Equals(o.SecondControlPoint) &&
			c.EndPoint.Equals(o.EndPoint)
	}
	return false
}

// BezierValueQuad calculates the value of the Quadratic Bezier-curve at t.
func BezierValueQuad(p1, p2, p3, t float64) float64 {
	oneMinusT := 1 - t
	p := (oneMinusT*oneMinusT)*p1 + (2*(oneMinusT)*t*p2) + (t*t)*p3
	return p
}

// BezierValueCubic calculates the value of the Cubic Bezier-curve at t.
func BezierValueCubic(p1, p2, p3, p4, t float64) float64 {
	oneMinusT := 1 - t
	p := (oneMinusT*oneMinusT*oneMinusT)*p1 +
		(3*(oneMinusT*oneMinusT)*t*p2) +
		(3*oneMinusT*(t*t)*p3) +
		(t*t*t)*p4
	return p
}

// GetStartPoint returns the start point of a path command.
func GetStartPoint(command PathCommand) PdfPoint {
	switch cmd := command.(type) {
	case *Line:
		return cmd.From
	case *QuadraticBezierCurve:
		return cmd.StartPoint
	case *CubicBezierCurve:
		return cmd.StartPoint
	case *Move:
		return cmd.Location
	default:
		return PdfPoint{}
	}
}

// GetEndPoint returns the end point of a path command.
func GetEndPoint(command PathCommand) PdfPoint {
	switch cmd := command.(type) {
	case *Line:
		return cmd.To
	case *QuadraticBezierCurve:
		return cmd.EndPoint
	case *CubicBezierCurve:
		return cmd.EndPoint
	case *Move:
		return cmd.Location
	default:
		return PdfPoint{}
	}
}
