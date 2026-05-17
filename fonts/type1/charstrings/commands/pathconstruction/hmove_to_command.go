package pathconstruction

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// HMoveToCommand implements the hmoveto operator.
// Moves to a new position for horizontal dimension only,
// relative to the current position.
type HMoveToCommand struct{}

const HMoveToName = "hmoveto"

const HMoveToFirst byte = 22

// HMoveToSecond is nil since hmoveto uses only a single-byte code.
var HMoveToSecond *byte = nil

var HMoveToTakeFromStackBottom = true

var HMoveToClearsOperandStack = true

var HMoveToLazy = commands.NewLazyType1Command(HMoveToName, HMoveToRun)

// HMoveToRun executes the hmoveto command on the given context.
func HMoveToRun(context *commands.Type1BuildCharContext) {
	deltaX, _ := context.Stack().PopBottom()

	if context.IsFlexing() {
		context.AddFlexPoint(core.NewPdfPoint(deltaX, 0))
	} else {
		x := context.CurrentPosition().X + deltaX
		y := context.CurrentPosition().Y
		context.SetCurrentPosition(core.NewPdfPoint(x, y))

		path := context.Path()
		path = append(path, core.NewPdfSubpath())
		path[len(path)-1].MoveTo(x, y)
		context.SetPath(path)
	}

	context.Stack().Clear()
}
