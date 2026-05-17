package hint

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// VStemName is the PostScript operator name for vstem.
const VStemName = "vstem"

// VStemFirst byte of the vstem opcode sequence.
const VStemFirst byte = 3

// VStemSecond byte of the vstem opcode sequence (nil for single-byte opcode).
var VStemSecond *byte

// VStemTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var VStemTakeFromStackBottom = true

// VStemClearsOperandStack indicates whether this command clears the operand stack.
var VStemClearsOperandStack = true

// VStemLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var VStemLazy = commands.NewLazyType1Command(VStemName, VStemRun)

// VStemRun executes the vstem command on the given context.
// It declares the horizontal range of a vertical stem zone between the x coordinates
// x and x+dx, where x is relative to the x coordinate of the left sidebearing point.
func VStemRun(context *commands.Type1BuildCharContext) {
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()

	context.Stack().Clear()
}
