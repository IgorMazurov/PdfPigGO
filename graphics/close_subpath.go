// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"io"
)

// CloseSubpath closes the current subpath by appending a straight line segment from
// the current point to the starting point of the subpath. If the current subpath is
// already closed, this does nothing.
type CloseSubpath struct{}

const closeSubpathSymbol = "h"

// InstanceCloseSubpath is the singleton instance of CloseSubpath.
var InstanceCloseSubpath = CloseSubpath{}

// Operator returns the operator symbol for this operation.
func (o CloseSubpath) Operator() string {
	return closeSubpathSymbol
}

// Run executes the operation on the given context, closing the current subpath.
func (o CloseSubpath) Run(ctx OperationContext) {
	point := ctx.CloseSubpath()
	if point != nil {
		ctx.SetCurrentPosition(*point)
	}
}

// Write writes the operator symbol to the given writer.
func (o CloseSubpath) Write(w io.Writer) error {
	if _, err := io.WriteString(w, closeSubpathSymbol); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (CloseSubpath) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o CloseSubpath) String() string {
	return closeSubpathSymbol
}
