package charstrings

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
)

// Type2BuildCharContext holds the state used while interpreting Type 2 CharString commands.
type Type2BuildCharContext struct {
	stack          *fonts.CharStringStack
	path           []*core.PdfSubpath
	currentLocation core.PdfPoint
	width          *float64
	transientArray map[int]float64
}

// NewType2BuildCharContext creates a new Type2BuildCharContext.
func NewType2BuildCharContext() *Type2BuildCharContext {
	return &Type2BuildCharContext{
		stack:          fonts.NewCharStringStack(),
		path:           make([]*core.PdfSubpath, 0),
		currentLocation: core.Origin,
		transientArray: make(map[int]float64),
	}
}

// Stack returns the operand stack for the charstring interpreter.
func (c *Type2BuildCharContext) Stack() *fonts.CharStringStack {
	return c.stack
}

// Path returns the list of subpaths built so far.
func (c *Type2BuildCharContext) Path() []*core.PdfSubpath {
	return c.path
}

// CurrentLocation returns the current pen position.
func (c *Type2BuildCharContext) CurrentLocation() core.PdfPoint {
	return c.currentLocation
}

// SetCurrentLocation sets the current pen position.
func (c *Type2BuildCharContext) SetCurrentLocation(p core.PdfPoint) {
	c.currentLocation = p
}

// Width returns the character width if one has been set, or nil otherwise.
func (c *Type2BuildCharContext) Width() *float64 {
	return c.width
}

// SetWidth sets the character width.
func (c *Type2BuildCharContext) SetWidth(w float64) {
	c.width = &w
}

// AddRelativeHorizontalLine adds a horizontal line segment relative to the current position.
func (c *Type2BuildCharContext) AddRelativeHorizontalLine(dx float64) {
	c.addRelativeLine(dx, 0)
}

// AddRelativeVerticalLine adds a vertical line segment relative to the current position.
func (c *Type2BuildCharContext) AddRelativeVerticalLine(dy float64) {
	c.addRelativeLine(0, dy)
}

// AddRelativeMoveTo moves to a new position relative to the current one, closing the
// previous subpath and starting a new one.
func (c *Type2BuildCharContext) AddRelativeMoveTo(dx, dy float64) {
	c.beforeMoveTo()
	newLocation := core.NewPdfPoint(c.currentLocation.X+dx, c.currentLocation.Y+dy)
	c.path[len(c.path)-1].MoveTo(newLocation.X, newLocation.Y)
	c.currentLocation = newLocation
}

// AddHorizontalMoveTo moves horizontally from the current position, closing the previous
// subpath and starting a new one.
func (c *Type2BuildCharContext) AddHorizontalMoveTo(dx float64) {
	c.beforeMoveTo()
	c.path[len(c.path)-1].MoveTo(c.currentLocation.X+dx, c.currentLocation.Y)
	c.currentLocation = c.currentLocation.MoveX(dx)
}

// AddVerticalMoveTo moves vertically from the current position, closing the previous
// subpath and starting a new one.
func (c *Type2BuildCharContext) AddVerticalMoveTo(dy float64) {
	c.beforeMoveTo()
	c.path[len(c.path)-1].MoveTo(c.currentLocation.X, c.currentLocation.Y+dy)
	c.currentLocation = c.currentLocation.MoveY(dy)
}

// AddRelativeBezierCurve adds a cubic Bezier curve using relative coordinates.
func (c *Type2BuildCharContext) AddRelativeBezierCurve(dx1, dy1, dx2, dy2, dx3, dy3 float64) {
	x1 := c.currentLocation.X + dx1
	y1 := c.currentLocation.Y + dy1

	x2 := x1 + dx2
	y2 := y1 + dy2

	x3 := x2 + dx3
	y3 := y2 + dy3

	c.path[len(c.path)-1].BezierCurveToCubic(x1, y1, x2, y2, x3, y3)
	c.currentLocation = core.NewPdfPoint(x3, y3)
}

// addRelativeLine adds a line segment relative to the current position.
func (c *Type2BuildCharContext) addRelativeLine(dx, dy float64) {
	dest := core.NewPdfPoint(c.currentLocation.X+dx, c.currentLocation.Y+dy)
	c.path[len(c.path)-1].LineTo(dest.X, dest.Y)
	c.currentLocation = dest
}

// AddVerticalStemHints records vertical stem hints. Currently a no-op placeholder.
func (c *Type2BuildCharContext) AddVerticalStemHints(hints []core.PdfRange) {}

// AddHorizontalStemHints records horizontal stem hints. Currently a no-op placeholder.
func (c *Type2BuildCharContext) AddHorizontalStemHints(hints []core.PdfRange) {}

// AddToTransientArray stores a value at the given location in the transient array.
func (c *Type2BuildCharContext) AddToTransientArray(value float64, location int) {
	c.transientArray[location] = value
}

// GetFromTransientArray retrieves and removes a value from the transient array.
func (c *Type2BuildCharContext) GetFromTransientArray(location int) (float64, bool) {
	val, ok := c.transientArray[location]
	if ok {
		delete(c.transientArray, location)
	}
	return val, ok
}

// beforeMoveTo closes the previous subpath (if any) and starts a new one. Per the
// Type 2 CharString Format spec, every character path must begin with a moveto operator.
func (c *Type2BuildCharContext) beforeMoveTo() {
	if len(c.path) > 0 {
		c.path[len(c.path)-1].CloseSubpath()
	}
	c.path = append(c.path, core.NewPdfSubpath())
}

// CountToBias computes the bias value for subroutine/call indexing given the number of
// subroutines. See Adobe Technical Note #5177.
func CountToBias(count int) int {
	switch {
	case count < 1240:
		return 107
	case count < 33900:
		return 1131
	default:
		return 32768
	}
}
