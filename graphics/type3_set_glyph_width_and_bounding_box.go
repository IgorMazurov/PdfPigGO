// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// Type3SetGlyphWidthAndBoundingBox sets width information for a glyph and declares
// that the glyph description specifies both its shape and color for a Type 3 font.
// It also sets the glyph bounding box.
type Type3SetGlyphWidthAndBoundingBox struct {
	HorizontalDisplacement float64
	VerticalDisplacement   float64
	LowerLeftX             float64
	LowerLeftY             float64
	UpperRightX            float64
	UpperRightY            float64
}

// type3SetGlyphWidthAndBoundingBoxSymbol is the operator symbol for this operation in a PDF content stream.
const type3SetGlyphWidthAndBoundingBoxSymbol = "d1"

// NewType3SetGlyphWidthAndBoundingBox creates a new Type3SetGlyphWidthAndBoundingBox operation.
func NewType3SetGlyphWidthAndBoundingBox(horizontalDisplacement, verticalDisplacement, lowerLeftX, lowerLeftY, upperRightX, upperRightY float64) *Type3SetGlyphWidthAndBoundingBox {
	return &Type3SetGlyphWidthAndBoundingBox{
		HorizontalDisplacement: horizontalDisplacement,
		VerticalDisplacement:   verticalDisplacement,
		LowerLeftX:             lowerLeftX,
		LowerLeftY:             lowerLeftY,
		UpperRightX:            upperRightX,
		UpperRightY:            upperRightY,
	}
}

// Operator returns the operator symbol for this operation.
func (o Type3SetGlyphWidthAndBoundingBox) Operator() string {
	return type3SetGlyphWidthAndBoundingBoxSymbol
}

// Run executes the operation on the given context. This is a no-op as per spec.
func (o *Type3SetGlyphWidthAndBoundingBox) Run(ctx OperationContext) {
}

// Write writes the operator to the given writer.
func (o Type3SetGlyphWidthAndBoundingBox) Write(w io.Writer) error {
	if bw, ok := w.(byteWriter); ok {
		if err := WriteDouble(bw, o.HorizontalDisplacement); err != nil {
			return err
		}
		if err := WriteWhiteSpace(bw); err != nil {
			return err
		}
		if err := WriteDouble(bw, o.VerticalDisplacement); err != nil {
			return err
		}
		if err := WriteWhiteSpace(bw); err != nil {
			return err
		}
		if err := WriteDouble(bw, o.LowerLeftX); err != nil {
			return err
		}
		if err := WriteWhiteSpace(bw); err != nil {
			return err
		}
		if err := WriteDouble(bw, o.LowerLeftY); err != nil {
			return err
		}
		if err := WriteWhiteSpace(bw); err != nil {
			return err
		}
		if err := WriteDouble(bw, o.UpperRightX); err != nil {
			return err
		}
		if err := WriteWhiteSpace(bw); err != nil {
			return err
		}
		if err := WriteNumberTextFloat(bw, o.UpperRightY, type3SetGlyphWidthAndBoundingBoxSymbol); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(w, "%s %s %s %s %s %s %s\n", FormatDouble(o.HorizontalDisplacement), FormatDouble(o.VerticalDisplacement), FormatDouble(o.LowerLeftX), FormatDouble(o.LowerLeftY), FormatDouble(o.UpperRightX), FormatDouble(o.UpperRightY), type3SetGlyphWidthAndBoundingBoxSymbol); err != nil {
			return err
		}
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (Type3SetGlyphWidthAndBoundingBox) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o Type3SetGlyphWidthAndBoundingBox) String() string {
	return fmt.Sprintf("%g %g %g %g %g %g %s", o.HorizontalDisplacement, o.VerticalDisplacement, o.LowerLeftX, o.LowerLeftY, o.UpperRightX, o.UpperRightY, type3SetGlyphWidthAndBoundingBoxSymbol)
}
