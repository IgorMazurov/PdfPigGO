// Package writer provides types for writing PDF documents and their content streams.
package writer

import (
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/graphics"
)

// GetGlobalTransform iterates through the graphics state operations and returns
// the active transformation matrix at stack depth 0 (i.e., the global transform).
// Only ModifyCurrentTransformationMatrix operations at stackDepth == 0 are considered.
// Push increases the stack depth, Pop decreases it.
func GetGlobalTransform(operations []content.GraphicsStateOperation) *core.TransformationMatrix {
	var activeMatrix *core.TransformationMatrix
	stackDepth := 0

	for _, op := range operations {
		switch o := op.(type) {
		case *graphics.ModifyCurrentTransformationMatrix:
			if stackDepth == 0 && len(o.Value) == 6 {
				matrix, _ := core.FromArray(o.Value[:])
				activeMatrix = &matrix
			}
		case graphics.Push:
			stackDepth++
		case graphics.Pop:
			stackDepth--
		}
	}

	return activeMatrix
}
