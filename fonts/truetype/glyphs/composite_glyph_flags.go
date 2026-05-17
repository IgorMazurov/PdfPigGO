package glyphs

// CompositeGlyphFlags specifies the meaning of the transformation entries
// for a composite glyph definition in TrueType fonts.
type CompositeGlyphFlags uint16

const (
	// Args1And2AreWords indicates that arguments are words, otherwise they are bytes.
	Args1And2AreWords CompositeGlyphFlags = 1 << iota

	// ArgsAreXAndYValues indicates that arguments are x y offset values,
	// otherwise they are point indices.
	ArgsAreXAndYValues

	// RoundXAndYToGrid indicates that if arguments are x y offset values,
	// the values are rounded to the closest grid lines before addition to the glyph.
	RoundXAndYToGrid

	// WeHaveAScale indicates that the scale value is read in 2.14 format
	// (between -2 and < 2) and the glyph is scaled before grid-fitting.
	// Otherwise scale is 1.
	WeHaveAScale

	// Reserved is reserved for future use and should be set to 0.
	Reserved

	// MoreComponents indicates that there is a glyph component following the current one.
	MoreComponents

	// WeHaveAnXAndYScale indicates that X is scaled differently from Y.
	WeHaveAnXAndYScale

	// WeHaveATwoByTwo indicates that there is a 2x2 transformation matrix
	// used to scale the component.
	WeHaveATwoByTwo

	// WeHaveInstructions indicates that there are instructions for the composite
	// character following the last component.
	WeHaveInstructions

	// UseMyMetrics forces advance width and left side bearing for the composite
	// to be equal to those from the original glyph.
	UseMyMetrics
)
