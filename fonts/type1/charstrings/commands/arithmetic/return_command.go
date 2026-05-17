package arithmetic

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// ReturnName is the PostScript operator name for return.
const ReturnName = "return"

// ReturnFirst byte of the return opcode sequence.
const ReturnFirst byte = 11

// ReturnSecond byte of the return opcode sequence (null since it's a single-byte opcode).
var ReturnSecond *byte // nil, indicating no second byte

// ReturnTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var ReturnTakeFromStackBottom = false

// ReturnClearsOperandStack indicates whether this command clears the operand stack.
var ReturnClearsOperandStack = false

// ReturnLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var ReturnLazy = commands.NewLazyType1Command(ReturnName, ReturnRun)

// ReturnRun executes the return command on the given context.
// It returns from a charstring subroutine and continues execution in the calling charstring.
func ReturnRun(context *commands.Type1BuildCharContext) {
	// Do nothing; the interpreter handles the return stack.
}
