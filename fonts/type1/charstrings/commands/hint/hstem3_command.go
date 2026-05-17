package hint

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// HStem3Name is the PostScript operator name for hstem3.
const HStem3Name = "hstem3"

// HStem3First byte of the hstem3 opcode sequence.
const HStem3First byte = 12

// HStem3Second byte of the hstem3 opcode sequence (two-byte opcode).
var HStem3Second *byte // initialized in init()

func init() {
	second := byte(2)
	HStem3Second = &second
}

// HStem3TakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var HStem3TakeFromStackBottom = true

// HStem3ClearsOperandStack indicates whether this command clears the operand stack.
var HStem3ClearsOperandStack = true

// HStem3Lazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var HStem3Lazy = commands.NewLazyType1Command(HStem3Name, HStem3Run)

// HStem3Run executes the hstem3 command on the given context.
// It declares the vertical ranges of three horizontal stem zones suited to letters
// with 3 horizontal stems like 'E'. The six operands (y0, dy0, y1, dy1, y2, dy2)
// are relative to the y coordinate of the left sidebearing point and are currently ignored.
func HStem3Run(context *commands.Type1BuildCharContext) {
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()

	context.Stack().Clear()
}
