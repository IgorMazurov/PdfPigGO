package page

// ColourDepthSimpleType represents colour depth values in PAGE XML documents.
type ColourDepthSimpleType byte

const (
	// BiLevel indicates bilevel (black and white) image data.
	BiLevel ColourDepthSimpleType = iota

	// GreyScale indicates greyscale image data.
	GreyScale

	// Colour indicates full colour image data.
	Colour

	// ColourDepthOther indicates an unrecognized colour depth type.
	ColourDepthOther
)
