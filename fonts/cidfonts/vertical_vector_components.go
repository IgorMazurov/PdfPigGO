package cidfonts

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/geometry"
)

// VerticalVectorComponents defines the default position and displacement vector
// vertical components for fonts which have vertical writing modes.
type VerticalVectorComponents struct {
	Position     float64
	Displacement float64
}

// DefaultVerticalVectorComponents is the default value if not defined by a font.
var DefaultVerticalVectorComponents = NewVerticalVectorComponents(800, -1000)

// NewVerticalVectorComponents creates a new VerticalVectorComponents.
func NewVerticalVectorComponents(position, displacement float64) VerticalVectorComponents {
	return VerticalVectorComponents{Position: position, Displacement: displacement}
}

// GetPositionVector returns the full position vector for a given glyph width.
// The full position vector unless overridden by the W2 array is (w0/2, Position),
// where w0 is the width of the given glyph.
func (v VerticalVectorComponents) GetPositionVector(glyphWidth float64) geometry.PdfVector {
	return geometry.NewPdfVector(glyphWidth/2.0, v.Position)
}

// GetDisplacementVector returns the full displacement vector (0, Displacement).
func (v VerticalVectorComponents) GetDisplacementVector() geometry.PdfVector {
	return geometry.NewPdfVector(0, v.Displacement)
}

// String returns a string representation of VerticalVectorComponents.
func (v VerticalVectorComponents) String() string {
	return fmt.Sprintf("Position: %g, Displacement: %g.", v.Position, v.Displacement)
}
