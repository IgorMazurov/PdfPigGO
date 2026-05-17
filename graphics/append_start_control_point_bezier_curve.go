// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"
)

// AppendStartControlPointBezierCurve appends a cubic Bezier curve to the current path.
// The curve extends from the current point to the point (x3, y3), using the current point
// and (x2, y2) as the Bezier control points.
type AppendStartControlPointBezierCurve struct {
	X2 float64
	Y2 float64
	X3 float64
	Y3 float64
}

const appendStartControlPointBezierCurveSymbol = "v"

// NewAppendStartControlPointBezierCurve creates a new AppendStartControlPointBezierCurve operation.
func NewAppendStartControlPointBezierCurve(x2, y2, x3, y3 float64) *AppendStartControlPointBezierCurve {
	return &AppendStartControlPointBezierCurve{X2: x2, Y2: y2, X3: x3, Y3: y3}
}

// Operator returns the operator symbol for this operation.
func (o AppendStartControlPointBezierCurve) Operator() string {
	return appendStartControlPointBezierCurveSymbol
}

// Run executes the operation on the given context, adding a cubic Bezier curve segment.
// The current point is used as the first control point, (x2, y2) as the second control point,
// and (x3, y3) as the end point.
func (o *AppendStartControlPointBezierCurve) Run(ctx OperationContext) {
	ctx.BezierCurveToStartCP(o.X2, o.Y2, o.X3, o.Y3)
}

// Write writes the operator to the given writer.
func (o AppendStartControlPointBezierCurve) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%g %g %g %g %s\n", o.X2, o.Y2, o.X3, o.Y3, appendStartControlPointBezierCurveSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (AppendStartControlPointBezierCurve) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o AppendStartControlPointBezierCurve) String() string {
	return fmt.Sprintf("%g %g %g %g %s", o.X2, o.Y2, o.X3, o.Y3, appendStartControlPointBezierCurveSymbol)
}
