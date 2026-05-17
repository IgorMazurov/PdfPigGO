package startfinishoutline

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// HsbwName is the PostScript operator name for hsbw.
const HsbwName = "hsbw"

// HsbwFirst byte of the hsbw opcode sequence.
const HsbwFirst byte = 13

// HsbwSecond byte of the hsbw opcode sequence (nil for single-byte opcode).
var HsbwSecond *byte // nil, no init() needed

// HsbwTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var HsbwTakeFromStackBottom = true

// HsbwClearsOperandStack indicates whether this command clears the operand stack.
var HsbwClearsOperandStack = true

// HsbwLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var HsbwLazy = commands.NewLazyType1Command(HsbwName, HsbwRun)

// HsbwRun executes the hsbw (horizontal sidebearing and width) command on the given context.
// It sets the left sidebearing point at (sbx, 0) and the character width vector to (wx, 0),
// then updates the current position to (sbx, 0) without placing it in the character path.
func HsbwRun(context *commands.Type1BuildCharContext) {
	leftSidebearingPointX, _ := context.Stack().PopBottom()
	characterWidthVectorX, _ := context.Stack().PopBottom()

	context.SetLeftSideBearingX(leftSidebearingPointX)
	context.SetWidthX(characterWidthVectorX)

	context.SetCurrentPosition(core.NewPdfPoint(leftSidebearingPointX, 0))

	context.Stack().Clear()
}
