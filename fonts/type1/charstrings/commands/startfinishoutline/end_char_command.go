package startfinishoutline

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// EndCharName is the PostScript operator name for endchar.
const EndCharName = "endchar"

// EndCharFirst byte of the endchar opcode sequence.
const EndCharFirst byte = 14

// EndCharSecond byte of the endchar opcode sequence (nil for single-byte opcode).
var EndCharSecond *byte // nil, no init() needed

// EndCharTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var EndCharTakeFromStackBottom = false

// EndCharClearsOperandStack indicates whether this command clears the operand stack.
var EndCharClearsOperandStack = true

// EndCharLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var EndCharLazy = commands.NewLazyType1Command(EndCharName, EndCharRun)

// EndCharRun executes the endchar command on the given context.
// It finishes a charstring outline definition and must be the last command
// in a character's outline (except for accented characters defined using seac).
func EndCharRun(context *commands.Type1BuildCharContext) {
	context.Stack().Clear()
}
