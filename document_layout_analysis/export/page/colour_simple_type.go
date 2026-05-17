package page

// PageXmlColourSimpleType represents colour values in PAGE XML documents.
type PageXmlColourSimpleType byte

const (
	// Black indicates black colour.
	Black PageXmlColourSimpleType = iota

	// Blue indicates blue colour.
	Blue

	// Brown indicates brown colour.
	Brown

	// Cyan indicates cyan colour.
	Cyan

	// Green indicates green colour.
	Green

	// Grey indicates grey colour.
	Grey

	// Indigo indicates indigo colour.
	Indigo

	// Magenta indicates magenta colour.
	Magenta

	// Orange indicates orange colour.
	Orange

	// Pink indicates pink colour.
	Pink

	// Red indicates red colour.
	Red

	// Turquoise indicates turquoise colour.
	Turquoise

	// Violet indicates violet colour.
	Violet

	// White indicates white colour.
	White

	// Yellow indicates yellow colour.
	Yellow

	// ColourOther indicates an unrecognized colour type.
	ColourOther
)
