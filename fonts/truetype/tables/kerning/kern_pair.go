package kerning

import "fmt"

// KernPair holds a kerning adjustment value for a pair of glyphs in a TrueType font.
type KernPair struct {
	// LeftGlyphIndex is the index of the left-hand glyph in the pair.
	LeftGlyphIndex int

	// RightGlyphIndex is the index of the right-hand glyph in the pair.
	RightGlyphIndex int

	// Value is the kerning adjustment. Positive values move characters apart,
	// negative values move them closer together.
	Value int16
}

// NewKernPair creates a new KernPair for the given glyph indices and kerning value.
func NewKernPair(leftGlyphIndex, rightGlyphIndex int, value int16) KernPair {
	return KernPair{
		LeftGlyphIndex:  leftGlyphIndex,
		RightGlyphIndex: rightGlyphIndex,
		Value:           value,
	}
}

// String returns a human-readable representation of the kern pair.
func (k KernPair) String() string {
	return fmt.Sprintf("Left: %d, Right: %d, Value %d.", k.LeftGlyphIndex, k.RightGlyphIndex, k.Value)
}
