package hint

import (
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// DotSectionName is the PostScript operator name for dotsection.
const DotSectionName = "dotsection"

// DotSectionFirst byte of the dotsection opcode sequence.
const DotSectionFirst byte = 12

// DotSectionSecond byte of the dotsection opcode sequence (two-byte opcode).
var DotSectionSecond *byte // initialized in init()

func init() {
	second := byte(0)
	DotSectionSecond = &second
}

// DotSectionTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var DotSectionTakeFromStackBottom = false

// DotSectionClearsOperandStack indicates whether this command clears the operand stack.
var DotSectionClearsOperandStack = true

// DotSectionLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var DotSectionLazy = commands.NewLazyType1Command(DotSectionName, DotSectionRun)

// DotSectionRun executes the dotsection command on the given context.
// It brackets an outline section for the dots in letters such as "i", "j" and "!".
func DotSectionRun(context *commands.Type1BuildCharContext) {
	context.Stack().Clear()
}
