// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// SetStrokeColorAdvanced sets the stroking color with support for Pattern,
// Separation, DeviceN, and ICCBased color spaces ("SCN" operator).
type SetStrokeColorAdvanced struct {
	Operands    []float64
	PatternName *tokens.NameToken
}

const setStrokeColorAdvancedSymbol = "SCN"

// NewSetStrokeColorAdvanced creates a new SetStrokeColorAdvanced with operands only.
func NewSetStrokeColorAdvanced(operands []float64) *SetStrokeColorAdvanced {
	return &SetStrokeColorAdvanced{Operands: operands}
}

// NewSetStrokeColorAdvancedWithPattern creates a new SetStrokeColorAdvanced with pattern name.
func NewSetStrokeColorAdvancedWithPattern(operands []float64, patternName *tokens.NameToken) *SetStrokeColorAdvanced {
	return &SetStrokeColorAdvanced{Operands: operands, PatternName: patternName}
}

// Operator returns the operator symbol for this operation.
func (o SetStrokeColorAdvanced) Operator() string {
	return setStrokeColorAdvancedSymbol
}

// Run executes the operation on the given context, setting the stroking color with advanced support.
func (o *SetStrokeColorAdvanced) Run(ctx OperationContext) {
	state := ctx.GetCurrentState()
	if state == nil || state.ColorSpaceContext == nil {
		return
	}
	state.ColorSpaceContext.SetStrokingColor(o.Operands, o.PatternName)
}

// Write writes the operator to the given writer.
func (o SetStrokeColorAdvanced) Write(w io.Writer) error {
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
	if _, err := fmt.Fprint(w, setStrokeColorAdvancedSymbol); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w)
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetStrokeColorAdvanced) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetStrokeColorAdvanced) String() string {
	parts := make([]string, len(o.Operands))
	for i, operand := range o.Operands {
		parts[i] = fmt.Sprintf("%g", operand)
	}
	if o.PatternName != nil {
		parts = append(parts, o.PatternName.String())
	}
	return strings.Join(parts, " ") + " " + setStrokeColorAdvancedSymbol
}
