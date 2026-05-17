package page

// PageXmlGraphicsSimpleType represents the type of graphic element in a PAGE XML document.
type PageXmlGraphicsSimpleType byte

const (
	// Logo indicates a logo graphic.
	Logo PageXmlGraphicsSimpleType = iota

	// Letterhead indicates a letterhead graphic.
	Letterhead

	// Decoration indicates a decorative graphic.
	Decoration

	// Frame indicates a frame graphic.
	Frame

	// HandwrittenAnnotation indicates a handwritten annotation graphic.
	HandwrittenAnnotation

	// Stamp indicates a stamp graphic.
	Stamp

	// Signature indicates a signature graphic.
	Signature

	// Barcode indicates a barcode graphic.
	Barcode

	// PaperGrow indicates a paper grow artifact.
	PaperGrow

	// PunchHole indicates a punch hole artifact.
	PunchHole

	// OtherGraphics indicates an unrecognized graphic type.
	OtherGraphics
)
