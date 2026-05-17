package hint

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// VStem3Name is the PostScript operator name for vstem3.
const VStem3Name = "vstem3"

// VStem3First byte of the vstem3 opcode sequence.
const VStem3First byte = 12

// VStem3Second byte of the vstem3 opcode sequence (two-byte opcode).
var VStem3Second *byte // initialized in init()

func init() {
	second := byte(1)
	VStem3Second = &second
}

// VStem3TakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var VStem3TakeFromStackBottom = true

// VStem3ClearsOperandStack indicates whether this command clears the operand stack.
var VStem3ClearsOperandStack = true

// VStem3Lazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var VStem3Lazy = commands.NewLazyType1Command(VStem3Name, VStem3Run)

// VStem3Run executes the vstem3 command on the given context.
// It declares the horizontal ranges of three vertical stem zones suited to letters
// with 3 vertical stems like 'm'. The six operands (x0, dx0, x1, dx1, x2, dx2)
// are relative to the x coordinate of the left sidebearing point and are currently ignored.
func VStem3Run(context *commands.Type1BuildCharContext) {
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()

	context.Stack().Clear()
}
