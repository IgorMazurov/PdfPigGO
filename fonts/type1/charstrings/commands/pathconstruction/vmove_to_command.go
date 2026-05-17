package pathconstruction

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// VMoveToCommand implements the vmoveto operator.
// Moves to a new position for vertical dimension only,
// relative to the current position.
type VMoveToCommand struct{}

const VMoveToName = "vmoveto"

const VMoveToFirst byte = 4

// VMoveToSecond is nil since vmoveto uses only a single-byte code.
var VMoveToSecond *byte = nil

var VMoveToTakeFromStackBottom = true

var VMoveToClearsOperandStack = true

var VMoveToLazy = commands.NewLazyType1Command(VMoveToName, VMoveToRun)

// VMoveToRun executes the vmoveto command on the given context.
func VMoveToRun(context *commands.Type1BuildCharContext) {
	deltaY, _ := context.Stack().PopBottom()

	if context.IsFlexing() {
		context.AddFlexPoint(core.NewPdfPoint(0, deltaY))
	} else {
		y := context.CurrentPosition().Y + deltaY
		x := context.CurrentPosition().X
		context.SetCurrentPosition(core.NewPdfPoint(x, y))

		path := context.Path()
		path = append(path, core.NewPdfSubpath())
		path[len(path)-1].MoveTo(x, y)
		context.SetPath(path)
	}

	context.Stack().Clear()
}
