package pathconstruction

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// HLineToCommand implements the hlineto operator.
// Draws a horizontal line from the current position to a new x-coordinate,
// keeping the y-coordinate unchanged.
type HLineToCommand struct{}

const HLineToName = "hlineto"

const HLineToFirst byte = 6

// HLineToSecond is nil since hlineto uses only a single-byte code.
var HLineToSecond *byte = nil

var HLineToTakeFromStackBottom = true

var HLineToClearsOperandStack = true

var HLineToLazy = commands.NewLazyType1Command(HLineToName, HLineToRun)

// HLineToRun executes the hlineto command on the given context.
func HLineToRun(context *commands.Type1BuildCharContext) {
	deltaX, _ := context.Stack().PopBottom()
	x := context.CurrentPosition().X + deltaX

	path := context.Path()
	if len(path) > 0 {
		path[len(path)-1].LineTo(x, context.CurrentPosition().Y)
	}

	context.SetCurrentPosition(core.NewPdfPoint(x, context.CurrentPosition().Y))
	context.Stack().Clear()
}
