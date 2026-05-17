package arithmetic

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// PopName is the PostScript operator name for pop.
const PopName = "pop"

// PopFirst byte of the pop opcode sequence.
const PopFirst byte = 12

// PopSecond byte of the pop opcode sequence (two-byte opcode).
var PopSecond *byte // initialized in init()

func init() {
	second := byte(17)
	PopSecond = &second
}

// PopTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var PopTakeFromStackBottom = false

// PopClearsOperandStack indicates whether this command clears the operand stack.
var PopClearsOperandStack = false

// PopLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var PopLazy = commands.NewLazyType1Command(PopName, PopRun)

// PopRun executes the pop command on the given context.
// It pops a number from the top of the PostScript operand stack and pushes
// that number onto the main operand stack. This command is used only to
// retrieve a result from an OtherSubrs procedure.
func PopRun(context *commands.Type1BuildCharContext) {
	num, _ := context.PostscriptStack().PopTop()
	context.Stack().Push(num)
}
