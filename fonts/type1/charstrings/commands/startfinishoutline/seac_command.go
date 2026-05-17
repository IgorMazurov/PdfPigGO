package startfinishoutline

import (
	"github.com/uglytoad/pdfpig/go/fonts/encodings"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
)

// SeacName is the PostScript operator name for seac.
const SeacName = "seac"

// SeacFirst byte of the seac opcode sequence.
const SeacFirst byte = 12

// SeacSecond byte of the seac opcode sequence.
var SeacSecond *byte

func init() {
	s := byte(6)
	SeacSecond = &s
}

// SeacTakeFromStackBottom indicates whether arguments are taken from the bottom of the stack.
var SeacTakeFromStackBottom = true

// SeacClearsOperandStack indicates whether this command clears the operand stack.
var SeacClearsOperandStack = true

// SeacLazy wraps the Run function for deferred execution within a Type 1 charstring interpreter.
var SeacLazy = commands.NewLazyType1Command(SeacName, SeacRun)

// SeacRun executes the seac (standard encoding accented character) command on the given context.
// It makes an accented character from two other characters in the font program:
// a base character and an accent character, both identified by their Adobe StandardEncoding code.
func SeacRun(context *commands.Type1BuildCharContext) {
	// Pop all 5 operands; first three unused pending full seac implementation.
	context.Stack().PopBottom() // accentLeftSidebearingX
	context.Stack().PopBottom() // accentOriginX
	context.Stack().PopBottom() // accentOriginY
	baseCharacterCode, _ := context.Stack().PopBottom()
	accentCharacterCode, _ := context.Stack().PopBottom()

	baseCharacterName := encodings.StandardEncodingValue.GetName(int(baseCharacterCode))
	accentCharacterName := encodings.StandardEncodingValue.GetName(int(accentCharacterCode))

	baseCharacter := context.GetCharacterByName(baseCharacterName)
	_ = context.GetCharacterByName(accentCharacterName)

	// TODO: full seac implementation.
	context.SetPath(baseCharacter)

	context.Stack().Clear()
}
