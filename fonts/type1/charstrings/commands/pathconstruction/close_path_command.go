package pathconstruction

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// ClosePathRun executes the closepath command on the given context.
// Closes the current sub-path without repositioning the current point.
func ClosePathRun(context *commands.Type1BuildCharContext) {
	path := context.Path()
	if len(path) > 0 {
		path[len(path)-1].CloseSubpath()
	}
	context.Stack().Clear()
}

const ClosePathName = "closepath"

const ClosePathFirst byte = 9

// ClosePathSecond is nil since closepath uses only a single-byte code.
var ClosePathSecond *byte = nil

var ClosePathTakeFromStackBottom = false

var ClosePathClearsOperandStack = true

var ClosePathLazy = commands.NewLazyType1Command(ClosePathName, ClosePathRun)
