package arithmetic

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// SetCurrentPointName is the PostScript operator name for setcurrentpoint.
const SetCurrentPointName = "setcurrentpoint"

// SetCurrentPointFirst byte of the setcurrentpoint opcode sequence.
const SetCurrentPointFirst byte = 12

// SetCurrentPointSecond byte of the setcurrentpoint opcode sequence (two-byte opcode).
var SetCurrentPointSecond *byte // initialized in init()

func init() {
	second := byte(33)
	SetCurrentPointSecond = &second
}

// SetCurrentPointTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var SetCurrentPointTakeFromStackBottom = true

// SetCurrentPointClearsOperandStack indicates whether this command clears the operand stack.
var SetCurrentPointClearsOperandStack = true

// SetCurrentPointLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var SetCurrentPointLazy = commands.NewLazyType1Command(SetCurrentPointName, SetCurrentPointRun)

// SetCurrentPointRun executes the setcurrentpoint command on the given context.
// It sets the current point to (x, y) in absolute character space coordinates without
// performing a charstring moveto command. This establishes the current point for a
// subsequent relative path building command. The setcurrentpoint command is used only
// in conjunction with results from OtherSubrs procedures.
func SetCurrentPointRun(context *commands.Type1BuildCharContext) {
	_, _ = context.Stack().PopBottom()
	_, _ = context.Stack().PopBottom()

	context.Stack().Clear()
}
