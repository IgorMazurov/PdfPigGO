package glyphs

// SimpleGlyphFlags specifies the meaning of each coordinate in the simple glyph definition.
type SimpleGlyphFlags uint8

const (
	// OnCurve indicates that the point is on the curve.
	OnCurve SimpleGlyphFlags = 1 << iota

	// XSingleByte indicates that the x-coordinate is 1 byte long instead of 2.
	XSingleByte

	// YSingleByte indicates that the y-coordinate is 1 byte long instead of 2.
	YSingleByte

	// Repeat indicates that the next byte specifies the number of times to repeat this set of flags.
	Repeat

	// ThisXIsTheSame means: if XSingleByte is set, the sign of the x-coordinate is positive;
	// if XSingleByte is not set, the current x-coordinate is the same as the previous.
	ThisXIsTheSame

	// ThisYIsTheSame means: if YSingleByte is set, the sign of the y-coordinate is positive;
	// if YSingleByte is not set, the current y-coordinate is the same as the previous.
	ThisYIsTheSame
)
