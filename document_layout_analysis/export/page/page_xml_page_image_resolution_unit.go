package page

// PageXmlPageImageResolutionUnit specifies the unit of resolution information
// referring to a standardised unit of measurement (pixels per inch, pixels per
// centimeter or other).
type PageXmlPageImageResolutionUnit byte

const (
	// PPI indicates pixels per inch.
	PPI PageXmlPageImageResolutionUnit = iota

	// PPCM indicates pixels per centimeter.
	PPCM

	// ResolutionUnitOther indicates a non-standard resolution unit.
	ResolutionUnitOther
)
