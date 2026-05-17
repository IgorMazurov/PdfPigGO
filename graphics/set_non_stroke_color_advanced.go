// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// SetNonStrokeColorAdvanced sets the non-stroking color with support for Pattern,
// Separation, DeviceN, and ICCBased color spaces ("scn" operator).
type SetNonStrokeColorAdvanced struct {
	Operands    []float64
	PatternName *tokens.NameToken
}

const setNonStrokeColorAdvancedSymbol = "scn"

// NewSetNonStrokeColorAdvanced creates a new SetNonStrokeColorAdvanced with operands only.
func NewSetNonStrokeColorAdvanced(operands []float64) *SetNonStrokeColorAdvanced {
	return &SetNonStrokeColorAdvanced{Operands: operands}
}

// NewSetNonStrokeColorAdvancedWithPattern creates a new SetNonStrokeColorAdvanced with pattern name.
func NewSetNonStrokeColorAdvancedWithPattern(operands []float64, patternName *tokens.NameToken) *SetNonStrokeColorAdvanced {
	return &SetNonStrokeColorAdvanced{Operands: operands, PatternName: patternName}
}

// Operator returns the operator symbol for this operation.
func (o SetNonStrokeColorAdvanced) Operator() string {
	return setNonStrokeColorAdvancedSymbol
}

// Run executes the operation on the given context, setting the non-stroking color with advanced support.
func (o *SetNonStrokeColorAdvanced) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetNonStrokingColor(o.Operands, o.PatternName)
}

// Write writes the operator to the given writer.
func (o SetNonStrokeColorAdvanced) Write(w io.Writer) error {
	for _, operand := range o.Operands {
		if _, err := fmt.Fprintf(w, "%g ", operand); err != nil {
			return err
		}
	}
	if o.PatternName != nil {
		if _, err := fmt.Fprint(w, o.PatternName.String()); err != nil {
			return err
		}
		_, _ = fmt.Fprint(w, " ")
	}
	if _, err := fmt.Fprint(w, setNonStrokeColorAdvancedSymbol); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetNonStrokeColorAdvanced) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetNonStrokeColorAdvanced) String() string {
	parts := make([]string, len(o.Operands))
	for i, operand := range o.Operands {
		parts[i] = fmt.Sprintf("%g", operand)
	}
	if o.PatternName != nil {
		parts = append(parts, o.PatternName.String())
	}
	return strings.Join(parts, " ") + " " + setNonStrokeColorAdvancedSymbol
}
