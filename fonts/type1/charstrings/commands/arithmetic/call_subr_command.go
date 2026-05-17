package arithmetic

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// CallSubrName is the PostScript operator name for callsubr.
const CallSubrName = "callsubr"

// CallSubrFirst byte of the callsubr opcode sequence.
const CallSubrFirst byte = 10

// CallSubrSecond byte of the callsubr opcode sequence (null for single-byte opcode).
var CallSubrSecond *byte = nil

// CallSubrTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var CallSubrTakeFromStackBottom = false

// CallSubrClearsOperandStack indicates whether this command clears the operand stack.
var CallSubrClearsOperandStack = false

// CallSubrLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var CallSubrLazy = commands.NewLazyType1Command(CallSubrName, CallSubrRun)

// CallSubrRun executes the callsubr command on the given context.
// It pops an index from the stack, looks up the corresponding subroutine,
// and executes each command in the subroutine's command sequence.
func CallSubrRun(context *commands.Type1BuildCharContext) {
	indexVal, _ := context.Stack().PopTop()
	index := int(indexVal)

	subroutineAny := context.Subroutines()[index]
	subroutine, ok := subroutineAny.([]core.Union[float64, *commands.LazyType1Command])
	if !ok {
		panic("invalid subroutine type")
	}

	for _, cmd := range subroutine {
		if num, ok := cmd.TryGetFirst(); ok {
			context.Stack().Push(num)
		} else if lazyCmd, ok := cmd.TryGetSecond(); ok {
			lazyCmd.Run(context)
		}
	}
}
