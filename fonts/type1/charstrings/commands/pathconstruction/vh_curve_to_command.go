package pathconstruction

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// VhCurveToCommand implements the vhcurveto operator.
// Draws a Bézier curve when the first Bézier tangent is vertical
// and the second Bézier tangent is horizontal.
// Equivalent to 0 dy1 dx2 dy2 dx3 0 rrcurveto.
type VhCurveToCommand struct{}

const VhCurveToName = "vhcurveto"

const VhCurveToFirst byte = 30

// VhCurveToSecond is nil since vhcurveto uses only a single-byte code.
var VhCurveToSecond *byte = nil

var VhCurveToTakeFromStackBottom = true

var VhCurveToClearsOperandStack = true

var VhCurveToLazy = commands.NewLazyType1Command(VhCurveToName, VhCurveToRun)

// VhCurveToRun executes the vhcurveto command on the given context.
func VhCurveToRun(context *commands.Type1BuildCharContext) {
	dy1, _ := context.Stack().PopBottom()
	dx2, _ := context.Stack().PopBottom()
	dy2, _ := context.Stack().PopBottom()
	dx3, _ := context.Stack().PopBottom()

	cp := context.CurrentPosition()
	x1 := cp.X
	y1 := cp.Y + dy1

	x2 := x1 + dx2
	y2 := y1 + dy2

	x3 := x2 + dx3
	y3 := y2

	path := context.Path()
	if len(path) > 0 {
		path[len(path)-1].BezierCurveToCubic(x1, y1, x2, y2, x3, y3)
	}

	context.SetCurrentPosition(core.NewPdfPoint(x3, y3))
	context.Stack().Clear()
}
