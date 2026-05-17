package pathconstruction

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// RMoveToCommand implements the rmoveto operator.
// Starts a new subpath relative to the current point,
// interpreting the number pair as a displacement rather than absolute coordinates.
type RMoveToCommand struct{}

const RMoveToName = "rmoveto"

const RMoveToFirst byte = 21

// RMoveToSecond is nil since rmoveto uses only a single-byte code.
var RMoveToSecond *byte = nil

var RMoveToTakeFromStackBottom = true

var RMoveToClearsOperandStack = true

var RMoveToLazy = commands.NewLazyType1Command(RMoveToName, RMoveToRun)

// RMoveToRun executes the rmoveto command on the given context.
func RMoveToRun(context *commands.Type1BuildCharContext) {
	deltaX, _ := context.Stack().PopBottom()
	deltaY, _ := context.Stack().PopBottom()

	if context.IsFlexing() {
		context.AddFlexPoint(core.NewPdfPoint(deltaX, deltaY))
	} else {
		x := context.CurrentPosition().X + deltaX
		y := context.CurrentPosition().Y + deltaY
		context.SetCurrentPosition(core.NewPdfPoint(x, y))

		path := context.Path()
		path = append(path, core.NewPdfSubpath())
		path[len(path)-1].MoveTo(x, y)
		context.SetPath(path)
	}

	context.Stack().Clear()
}
