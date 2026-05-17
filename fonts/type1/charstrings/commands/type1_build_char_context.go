package commands

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
)

// Type1BuildCharContext holds the state used while interpreting Type 1 CharString commands.
type Type1BuildCharContext struct {
	widthX            float64
	widthY            float64
	leftSideBearingX  float64
	leftSideBearingY  float64
	isFlexing         bool
	path              []*core.PdfSubpath
	currentPosition   core.PdfPoint
	stack             *fonts.CharStringStack
	postscriptStack   *fonts.CharStringStack
	flexPoints        []core.PdfPoint
	subroutines       map[int]any
	characterByIndex  func(int) []*core.PdfSubpath
	characterByName   func(string) []*core.PdfSubpath
}

// NewType1BuildCharContext creates a new Type1BuildCharContext.
// subroutines maps subroutine index to command sequences (type erased to avoid import cycles).
// characterByIndexFactory retrieves character paths by glyph index.
// characterByNameFactory retrieves character paths by character name.
func NewType1BuildCharContext(
	subroutines map[int]any,
	characterByIndexFactory func(int) []*core.PdfSubpath,
	characterByNameFactory func(string) []*core.PdfSubpath,
) *Type1BuildCharContext {
	return &Type1BuildCharContext{
		currentPosition:   core.Origin,
		stack:             fonts.NewCharStringStack(),
		postscriptStack:   fonts.NewCharStringStack(),
		flexPoints:        make([]core.PdfPoint, 0),
		subroutines:       subroutines,
		characterByIndex:  characterByIndexFactory,
		characterByName:   characterByNameFactory,
	}
}

// WidthX returns the horizontal width of the character.
func (c *Type1BuildCharContext) WidthX() float64 {
	return c.widthX
}

// SetWidthX sets the horizontal width of the character.
func (c *Type1BuildCharContext) SetWidthX(w float64) {
	c.widthX = w
}

// WidthY returns the vertical width of the character.
func (c *Type1BuildCharContext) WidthY() float64 {
	return c.widthY
}

// SetWidthY sets the vertical width of the character.
func (c *Type1BuildCharContext) SetWidthY(w float64) {
	c.widthY = w
}

// LeftSideBearingX returns the horizontal left-side bearing.
func (c *Type1BuildCharContext) LeftSideBearingX() float64 {
	return c.leftSideBearingX
}

// SetLeftSideBearingX sets the horizontal left-side bearing.
func (c *Type1BuildCharContext) SetLeftSideBearingX(b float64) {
	c.leftSideBearingX = b
}

// LeftSideBearingY returns the vertical left-side bearing.
func (c *Type1BuildCharContext) LeftSideBearingY() float64 {
	return c.leftSideBearingY
}

// SetLeftSideBearingY sets the vertical left-side bearing.
func (c *Type1BuildCharContext) SetLeftSideBearingY(b float64) {
	c.leftSideBearingY = b
}

// IsFlexing reports whether a flex operation is currently in progress.
func (c *Type1BuildCharContext) IsFlexing() bool {
	return c.isFlexing
}

// SetIsFlexing sets the flexing state.
func (c *Type1BuildCharContext) SetIsFlexing(v bool) {
	c.isFlexing = v
}

// Path returns the list of subpaths built so far.
func (c *Type1BuildCharContext) Path() []*core.PdfSubpath {
	return c.path
}

// SetPath replaces the current path with the given subpaths.
func (c *Type1BuildCharContext) SetPath(path []*core.PdfSubpath) {
	c.path = path
}

// CurrentPosition returns the current pen position.
func (c *Type1BuildCharContext) CurrentPosition() core.PdfPoint {
	return c.currentPosition
}

// SetCurrentPosition sets the current pen position.
func (c *Type1BuildCharContext) SetCurrentPosition(p core.PdfPoint) {
	c.currentPosition = p
}

// Stack returns the operand stack for the charstring interpreter.
func (c *Type1BuildCharContext) Stack() *fonts.CharStringStack {
	return c.stack
}

// PostscriptStack returns the PostScript operand stack used by certain commands.
func (c *Type1BuildCharContext) PostscriptStack() *fonts.CharStringStack {
	return c.postscriptStack
}

// FlexPoints returns the list of flex points accumulated during a flex operation.
func (c *Type1BuildCharContext) FlexPoints() []core.PdfPoint {
	return c.flexPoints
}

// AddFlexPoint adds a point to the flex points list.
func (c *Type1BuildCharContext) AddFlexPoint(point core.PdfPoint) {
	c.flexPoints = append(c.flexPoints, point)
}

// ClearFlexPoints removes all accumulated flex points.
func (c *Type1BuildCharContext) ClearFlexPoints() {
	c.flexPoints = c.flexPoints[:0]
}

// Subroutines returns the subroutine lookup map.
func (c *Type1BuildCharContext) Subroutines() map[int]any {
	return c.subroutines
}

// GetCharacterByIndex retrieves character subpaths by glyph index.
func (c *Type1BuildCharContext) GetCharacterByIndex(characterCode int) []*core.PdfSubpath {
	return c.characterByIndex(characterCode)
}

// GetCharacterByName retrieves character subpaths by character name.
func (c *Type1BuildCharContext) GetCharacterByName(characterName string) []*core.PdfSubpath {
	return c.characterByName(characterName)
}
