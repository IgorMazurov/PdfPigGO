package pathconstruction

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// HvCurveToCommand implements the hvcurveto operator.
// Draws a Bézier curve when the first Bézier tangent is horizontal
// and the second Bézier tangent is vertical.
// Equivalent to dx1 0 dx2 dy2 0 dy3 rrcurveto.
type HvCurveToCommand struct{}

const HvCurveToName = "hvcurveto"

const HvCurveToFirst byte = 31

// HvCurveToSecond is nil since hvcurveto uses only a single-byte code.
var HvCurveToSecond *byte = nil

var HvCurveToTakeFromStackBottom = true

var HvCurveToClearsOperandStack = true

var HvCurveToLazy = commands.NewLazyType1Command(HvCurveToName, HvCurveToRun)

// HvCurveToRun executes the hvcurveto command on the given context.
func HvCurveToRun(context *commands.Type1BuildCharContext) {
	dx1, _ := context.Stack().PopBottom()
	dx2, _ := context.Stack().PopBottom()
	dy2, _ := context.Stack().PopBottom()
	dy3, _ := context.Stack().PopBottom()

	cp := context.CurrentPosition()
	x1 := cp.X + dx1
	y1 := cp.Y

	x2 := x1 + dx2
	y2 := y1 + dy2

	x3 := x2
	y3 := y2 + dy3

	path := context.Path()
	if len(path) > 0 {
		path[len(path)-1].BezierCurveToCubic(x1, y1, x2, y2, x3, y3)
	}

	context.SetCurrentPosition(core.NewPdfPoint(x3, y3))
	context.Stack().Clear()
}
