package pathconstruction

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// VLineToCommand implements the vlineto operator.
// Draws a vertical line from the current position to a new y-coordinate,
// keeping the x-coordinate unchanged.
type VLineToCommand struct{}

const VLineToName = "vlineto"

const VLineToFirst byte = 7

// VLineToSecond is nil since vlineto uses only a single-byte code.
var VLineToSecond *byte = nil

var VLineToTakeFromStackBottom = true

var VLineToClearsOperandStack = true

var VLineToLazy = commands.NewLazyType1Command(VLineToName, VLineToRun)

// VLineToRun executes the vlineto command on the given context.
func VLineToRun(context *commands.Type1BuildCharContext) {
	deltaY, _ := context.Stack().PopBottom()
	y := context.CurrentPosition().Y + deltaY

	path := context.Path()
	if len(path) > 0 {
		path[len(path)-1].LineTo(context.CurrentPosition().X, y)
	}

	context.SetCurrentPosition(core.NewPdfPoint(context.CurrentPosition().X, y))
	context.Stack().Clear()
}
