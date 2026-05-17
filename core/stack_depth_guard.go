package core

import (
	"fmt"
	"math"
)

// StackDepthGuard tracks and limits the depth of nested stack operations,
// such as recursive calls or nested parsing, to prevent excessive stack usage.
type StackDepthGuard struct {
	maxStackDepth int
	depth         int
}

// Infinite is a StackDepthGuard with no effective limit on the allowed depth.
var Infinite = &StackDepthGuard{maxStackDepth: int(math.MaxInt)}

// NewStackDepthGuard creates a new StackDepthGuard with the specified maximum
// stack depth. The maxStackDepth must be positive.
func NewStackDepthGuard(maxStackDepth int) (*StackDepthGuard, error) {
	if maxStackDepth <= 0 {
		return nil, fmt.Errorf("maxStackDepth must be positive, got %d", maxStackDepth)
	}
	return &StackDepthGuard{maxStackDepth: maxStackDepth}, nil
}

// Enter increments the current nesting depth and checks against the maximum
// allowed stack depth. Returns a PdfDocumentFormatException if the limit is exceeded.
func (g *StackDepthGuard) Enter() error {
	g.depth++
	if g.depth > g.maxStackDepth {
		g.depth--
		return NewPdfDocumentFormatException(
			fmt.Sprintf("Exceeded maximum nesting depth of %d.", g.maxStackDepth),
		)
	}
	return nil
}

// Exit decreases the current depth level by one, clamping to zero if already at minimum.
func (g *StackDepthGuard) Exit() {
	g.depth--
	if g.depth < 0 {
		g.depth = 0
	}
}
