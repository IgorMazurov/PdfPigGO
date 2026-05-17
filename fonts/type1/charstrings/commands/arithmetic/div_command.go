package arithmetic

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// DivName is the PostScript operator name for div.
const DivName = "div"

// DivFirst byte of the div opcode sequence.
const DivFirst byte = 12

// DivSecond byte of the div opcode sequence (two-byte opcode).
var DivSecond *byte // initialized in init()

func init() {
	second := byte(12)
	DivSecond = &second
}

// DivTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var DivTakeFromStackBottom = false

// DivClearsOperandStack indicates whether this command clears the operand stack.
var DivClearsOperandStack = false

// DivLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var DivLazy = commands.NewLazyType1Command(DivName, DivRun)

// DivRun executes the div command on the given context.
// It pops two operands from the stack (num2 and num1), computes num1 / num2,
// and pushes the result back onto the stack. The result is always a real.
func DivRun(context *commands.Type1BuildCharContext) {
	first, _ := context.Stack().PopTop()
	second, _ := context.Stack().PopTop()

	result := second / first

	context.Stack().Push(result)
}
