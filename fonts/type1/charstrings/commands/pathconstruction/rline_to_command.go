package pathconstruction

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// RLineToCommand implements the rlineto operator.
// Draws a line from the current position to a new position,
// relative in both x and y dimensions.
type RLineToCommand struct{}

const RLineToName = "rlineto"

const RLineToFirst byte = 5

// RLineToSecond is nil since rlineto uses only a single-byte code.
var RLineToSecond *byte = nil

var RLineToTakeFromStackBottom = true

var RLineToClearsOperandStack = true

var RLineToLazy = commands.NewLazyType1Command(RLineToName, RLineToRun)

// RLineToRun executes the rlineto command on the given context.
func RLineToRun(context *commands.Type1BuildCharContext) {
	deltaX, _ := context.Stack().PopBottom()
	deltaY, _ := context.Stack().PopBottom()

	x := context.CurrentPosition().X + deltaX
	y := context.CurrentPosition().Y + deltaY

	path := context.Path()
	if len(path) > 0 {
		path[len(path)-1].LineTo(x, y)
	}

	context.SetCurrentPosition(core.NewPdfPoint(x, y))
	context.Stack().Clear()
}
