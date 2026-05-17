package hint

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// HStemName is the PostScript operator name for hstem.
const HStemName = "hstem"

// HStemFirst byte of the hstem opcode sequence.
const HStemFirst byte = 1

// HStemSecond byte of the hstem opcode sequence (nil for single-byte opcode).
var HStemSecond *byte

// HStemTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var HStemTakeFromStackBottom = true

// HStemClearsOperandStack indicates whether this command clears the operand stack.
var HStemClearsOperandStack = true

// HStemLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var HStemLazy = commands.NewLazyType1Command(HStemName, HStemRun)

// HStemRun executes the hstem command on the given context.
// It declares the vertical range of a horizontal stem zone between the y coordinates
// y and y+dy, where y is relative to the y coordinate of the left sidebearing point.
func HStemRun(context *commands.Type1BuildCharContext) {
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()

	context.Stack().Clear()
}
