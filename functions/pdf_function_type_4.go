package functions

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// PdfFunctionType4 represents a PostScript calculator function (PDF function type 4).
// The function stream contains a sequence of PostScript operators and operands that are
// executed at evaluation time using a stack-based interpreter.
type PdfFunctionType4 struct {
	PdfFunctionBase

	instructions *InstructionSequence
}

// NewPdfFunctionType4 creates a PdfFunctionType4 backed by a stream token.
// The function stream data is decoded as ISO 8859-1 and parsed into an instruction sequence.
func NewPdfFunctionType4(stream *tokens.StreamToken, domain, rangeVals *tokens.ArrayToken) (*PdfFunctionType4, error) {
	if stream == nil {
		return nil, fmt.Errorf("function stream must not be nil")
	}

	str := core.BytesAsLatin1String(stream.Data())
	instructions, err := ParseInstructionSequence(str)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type 4 function stream: %w", err)
	}

	return &PdfFunctionType4{
		PdfFunctionBase: *NewPdfFunctionBaseFromStream(stream, domain, rangeVals),
		instructions:    instructions,
	}, nil
}

// FunctionType returns PostScript for this function type.
func (f *PdfFunctionType4) FunctionType() FunctionTypes {
	return PostScript
}

// Eval evaluates the PostScript calculator function at the given input values.
// Input values are clipped to their domain ranges and pushed onto the operand stack.
// The instruction sequence is then executed, and output values are popped from the stack
// and clipped to their range bounds.
func (f *PdfFunctionType4) Eval(input ...float64) []float64 {
	operators := NewOperators()
	context := NewExecutionContext(operators)

	for i := 0; i < len(input); i++ {
		domain := f.GetDomainForInput(i)
		value := ClipToRange(input[i], domain.Min(), domain.Max())
		context.Push(value)
	}

	if err := f.instructions.Execute(context); err != nil {
		panic(err)
	}

	numOutputValues := f.NumberOfOutputParameters()
	actualOutputCount := context.Count()
	if actualOutputCount < numOutputValues {
		panic(fmt.Errorf("the type 4 function returned %d values but the Range entry indicates that %d values be returned",
			actualOutputCount, numOutputValues))
	}

	outputValues := make([]float64, numOutputValues)
	for i := numOutputValues - 1; i >= 0; i-- {
		rng := f.GetRangeForOutput(i)
		val := context.PopReal()
		outputValues[i] = ClipToRange(val, rng.Min(), rng.Max())
	}

	return outputValues
}
