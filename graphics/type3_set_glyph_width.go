// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// Type3SetGlyphWidth sets width information for a glyph and declares that the
// glyph description specifies both its shape and color for a Type 3 font.
type Type3SetGlyphWidth struct {
	HorizontalDisplacement float64
	VerticalDisplacement   float64
}

// type3SetGlyphWidthSymbol is the operator symbol for this operation in a PDF content stream.
const type3SetGlyphWidthSymbol = "d0"

// NewType3SetGlyphWidth creates a new Type3SetGlyphWidth operation.
func NewType3SetGlyphWidth(horizontalDisplacement, verticalDisplacement float64) *Type3SetGlyphWidth {
	return &Type3SetGlyphWidth{
		HorizontalDisplacement: horizontalDisplacement,
		VerticalDisplacement:   verticalDisplacement,
	}
}

// Operator returns the operator symbol for this operation.
func (o Type3SetGlyphWidth) Operator() string {
	return type3SetGlyphWidthSymbol
}

// Run executes the operation on the given context. This is a no-op as per spec.
func (o *Type3SetGlyphWidth) Run(ctx OperationContext) {
}

// Write writes the operator to the given writer.
func (o Type3SetGlyphWidth) Write(w io.Writer) error {
	if bw, ok := w.(byteWriter); ok {
		if err := WriteDouble(bw, o.HorizontalDisplacement); err != nil {
			return err
		}
		if err := WriteWhiteSpace(bw); err != nil {
			return err
		}
		if err := WriteNumberTextFloat(bw, o.VerticalDisplacement, type3SetGlyphWidthSymbol); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(w, "%s %g %s\n", FormatDouble(o.HorizontalDisplacement), o.VerticalDisplacement, type3SetGlyphWidthSymbol); err != nil {
			return err
		}
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (Type3SetGlyphWidth) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o Type3SetGlyphWidth) String() string {
	return fmt.Sprintf("%g %g %s", o.HorizontalDisplacement, o.VerticalDisplacement, type3SetGlyphWidthSymbol)
}
