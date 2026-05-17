// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
	"strings"

	"github.com/uglytoad/pdfpig/go/graphics/core"
)

// SetLineDashPattern sets the line dash pattern in the graphics state.
type SetLineDashPattern struct {
	Pattern core.LineDashPattern
}

// Symbol is the operator symbol for this operation in a PDF content stream.
const setLineDashPatternSymbol = "d"

// NewSetLineDashPattern creates a new SetLineDashPattern operation.
func NewSetLineDashPattern(array []float64, phase int) (*SetLineDashPattern, error) {
	pattern, err := core.NewLineDashPattern(phase, array)
	if err != nil {
		return nil, fmt.Errorf("invalid line dash pattern: %w", err)
	}
	return &SetLineDashPattern{Pattern: pattern}, nil
}

// Operator returns the operator symbol for this operation.
func (o SetLineDashPattern) Operator() string {
	return setLineDashPatternSymbol
}

// Run executes the operation on the given context, setting the line dash pattern.
func (o *SetLineDashPattern) Run(ctx OperationContext) {
	ctx.SetLineDashPattern(o.Pattern)
}

// Write writes the operator to the given writer.
func (o SetLineDashPattern) Write(w io.Writer) error {
	if _, err := fmt.Fprint(w, "["); err != nil {
		return err
	}

	for i, v := range o.Pattern.Array {
		if i > 0 {
			if _, err := fmt.Fprint(w, " "); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "%g", v); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "] %d %s\n", o.Pattern.Phase, setLineDashPatternSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (SetLineDashPattern) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o SetLineDashPattern) String() string {
	parts := make([]string, len(o.Pattern.Array))
	for i, v := range o.Pattern.Array {
		parts[i] = fmt.Sprintf("%g", v)
	}
	return fmt.Sprintf("[%s] %d %s", strings.Join(parts, " "), o.Pattern.Phase, setLineDashPatternSymbol)
}
