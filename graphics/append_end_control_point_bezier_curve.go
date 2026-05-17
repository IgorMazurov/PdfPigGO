// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// AppendEndControlPointBezierCurve appends a cubic Bezier curve to the current path.
// The curve extends from the current point to the point (x3, y3), using (x1, y1) and
// (x3, y3) as the Bezier control points.
type AppendEndControlPointBezierCurve struct {
	X1 float64
	Y1 float64
	X3 float64
	Y3 float64
}

const appendEndControlPointBezierCurveSymbol = "y"

// NewAppendEndControlPointBezierCurve creates a new AppendEndControlPointBezierCurve operation.
func NewAppendEndControlPointBezierCurve(x1, y1, x3, y3 float64) *AppendEndControlPointBezierCurve {
	return &AppendEndControlPointBezierCurve{X1: x1, Y1: y1, X3: x3, Y3: y3}
}

// Operator returns the operator symbol for this operation.
func (o AppendEndControlPointBezierCurve) Operator() string {
	return appendEndControlPointBezierCurveSymbol
}

// Run executes the operation on the given context, adding a cubic Bezier curve segment.
func (o *AppendEndControlPointBezierCurve) Run(ctx OperationContext) {
	ctx.BezierCurveToCubic(o.X1, o.Y1, o.X3, o.Y3, o.X3, o.Y3)
}

// Write writes the operator to the given writer.
func (o AppendEndControlPointBezierCurve) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %g %g %s\n", o.X1, o.Y1, o.X3, o.Y3, appendEndControlPointBezierCurveSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (AppendEndControlPointBezierCurve) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o AppendEndControlPointBezierCurve) String() string {
	return fmt.Sprintf("%g %g %g %g %s", o.X1, o.Y1, o.X3, o.Y3, appendEndControlPointBezierCurveSymbol)
}
