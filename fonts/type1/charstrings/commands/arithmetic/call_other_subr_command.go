package arithmetic

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands/pathconstruction"
)

const (
	flexEnd           = 0
	flexBegin         = 1
	flexMiddle        = 2
	hintReplacement   = 3
)

// Name is the PostScript operator name for callothersubr.
const Name = "callothersubr"

// First byte of the callothersubr opcode sequence.
const First byte = 12

// Second byte of the callothersubr opcode sequence.
var Second *byte = ptrByte(16)

// TakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var TakeFromStackBottom = false

// ClearsOperandStack indicates whether this command clears the operand stack.
var ClearsOperandStack = false

// Lazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var Lazy = commands.NewLazyType1Command(Name, Run)

// Run executes the callothersubr command on the given context.
// It pops an index and argument count from the stack, then dispatches to one of four
// built-in other subroutines: flex end (0), flex begin (1), flex middle (2), or hint replacement (3).
func Run(context *commands.Type1BuildCharContext) {
	indexVal, _ := context.Stack().PopTop()
	index := int(indexVal)

	numberOfArgumentsVal, _ := context.Stack().PopTop()
	numberOfArguments := int(numberOfArgumentsVal)

	otherSubroutineArguments := make([]float64, numberOfArguments)
	for j := 0; j < numberOfArguments; j++ {
		otherSubroutineArguments[j], _ = context.Stack().PopTop()
	}

	switch index {
	case flexEnd:
		context.SetIsFlexing(false)

		flexPoints := context.FlexPoints()
		if len(flexPoints) < 7 {
			panic("there must be at least 7 flex points defined by an other subroutine")
		}

		cp := context.CurrentPosition()

		reference := flexPoints[0]
		reference = reference.Translate(cp.X, cp.Y)

		first := flexPoints[1]
		first = first.Translate(reference.X, reference.Y)
		first = first.Translate(-cp.X, -cp.Y)

		context.Stack().Push(first.X)
		context.Stack().Push(first.Y)
		context.Stack().Push(flexPoints[2].X)
		context.Stack().Push(flexPoints[2].Y)
		context.Stack().Push(flexPoints[3].X)
		context.Stack().Push(flexPoints[3].Y)
		pathconstruction.Run(context)

		context.Stack().Push(flexPoints[4].X)
		context.Stack().Push(flexPoints[4].Y)
		context.Stack().Push(flexPoints[5].X)
		context.Stack().Push(flexPoints[5].Y)
		context.Stack().Push(flexPoints[6].X)
		context.Stack().Push(flexPoints[6].Y)
		pathconstruction.Run(context)

		context.ClearFlexPoints()

	case flexBegin:
		context.PostscriptStack().Clear()
		cp := context.CurrentPosition()
		context.PostscriptStack().Push(cp.X)
		context.PostscriptStack().Push(cp.Y)
		context.SetIsFlexing(true)

	case flexMiddle:
		cp := context.CurrentPosition()
		context.PostscriptStack().Push(cp.X)
		context.PostscriptStack().Push(cp.Y)

	case hintReplacement:
		if len(otherSubroutineArguments) != 1 {
			panic("the hint replacement subroutine only takes a single argument")
		}
		context.PostscriptStack().Clear()
		context.PostscriptStack().Push(otherSubroutineArguments[0])

	default:
		context.PostscriptStack().Clear()
		for i := 0; i < len(otherSubroutineArguments); i++ {
			context.PostscriptStack().Push(otherSubroutineArguments[i])
		}
	}
}

func ptrByte(v byte) *byte {
	return &v
}
