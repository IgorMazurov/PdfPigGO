package startfinishoutline

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// SbwName is the PostScript operator name for sbw.
const SbwName = "sbw"

// SbwFirst byte of the sbw opcode sequence.
const SbwFirst byte = 12

// SbwSecond byte of the sbw opcode sequence.
var SbwSecond *byte

func init() {
	s := byte(7)
	SbwSecond = &s
}

// SbwTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var SbwTakeFromStackBottom = true

// SbwClearsOperandStack indicates whether this command clears the operand stack.
var SbwClearsOperandStack = true

// SbwLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var SbwLazy = commands.NewLazyType1Command(SbwName, SbwRun)

// SbwRun executes the sbw (sidebearing and width) command on the given context.
// It sets the left sidebearing point at (sbx, sby) and the character width vector to (wx, wy),
// then updates the current position to (sbx, sby) without placing it in the character path.
func SbwRun(context *commands.Type1BuildCharContext) {
	leftSidebearingX, _ := context.Stack().PopBottom()
	leftSidebearingY, _ := context.Stack().PopBottom()
	characterWidthX, _ := context.Stack().PopBottom()
	characterWidthY, _ := context.Stack().PopBottom()

	context.SetLeftSideBearingX(leftSidebearingX)
	context.SetLeftSideBearingY(leftSidebearingY)

	context.SetWidthX(characterWidthX)
	context.SetWidthY(characterWidthY)

	context.SetCurrentPosition(core.NewPdfPoint(leftSidebearingX, leftSidebearingY))

	context.Stack().Clear()
}
