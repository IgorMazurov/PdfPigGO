package pathconstruction

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// RelativeRCurveToCommand implements the rrcurveto operator.
// Unlike rcurveto where all arguments are relative to the current point,
// rrcurveto arguments are relative to each other sequentially.
// Equivalent to: dx1 dy1 (dx1+dx2) (dy1+dy2) (dx1+dx2+dx3) (dy1+dy2+dy3) rcurveto.
type RelativeRCurveToCommand struct{}

const Name = "rrcurveto"

const First byte = 8

// Second is nil since rrcurveto uses only a single-byte code.
var Second *byte = nil

var TakeFromStackBottom = true

var ClearsOperandStack = true

var Lazy = commands.NewLazyType1Command(Name, Run)

// Run executes the rrcurveto command on the given context.
func Run(context *commands.Type1BuildCharContext) {
	dx1, _ := context.Stack().PopBottom()
	dy1, _ := context.Stack().PopBottom()
	dx2, _ := context.Stack().PopBottom()
	dy2, _ := context.Stack().PopBottom()
	dx3, _ := context.Stack().PopBottom()
	dy3, _ := context.Stack().PopBottom()

	cp := context.CurrentPosition()
	x1 := cp.X + dx1
	y1 := cp.Y + dy1

	x2 := x1 + dx2
	y2 := y1 + dy2

	x3 := x2 + dx3
	y3 := y2 + dy3

	path := context.Path()
	if len(path) > 0 {
		path[len(path)-1].BezierCurveToCubic(x1, y1, x2, y2, x3, y3)
	}

	context.SetCurrentPosition(core.NewPdfPoint(x3, y3))
	context.Stack().Clear()
}
