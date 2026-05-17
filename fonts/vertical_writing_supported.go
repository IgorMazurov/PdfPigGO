package fonts

import (
	"github.com/uglytoad/pdfpig/go/geometry"
)

// VerticalWritingSupported describes a font that supports vertical writing mode
// in addition to the default horizontal writing mode.
type VerticalWritingSupported interface {
	// GetPositionVector returns the position vector for the given character code.
	// In vertical fonts the glyph position is described by a position vector from
	// the origin used for horizontal writing. The position vector is applied to
	// the horizontal writing origin to give a new vertical writing origin.
	GetPositionVector(characterCode int) geometry.PdfVector

	// GetDisplacementVector returns the displacement vector for the given character code.
	GetDisplacementVector(characterCode int) geometry.PdfVector
}
