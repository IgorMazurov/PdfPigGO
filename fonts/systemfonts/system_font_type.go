package systemfonts

// SystemFontType represents the type of a system font file format.
type SystemFontType int

const (
	// Unknown indicates an unrecognized or undetermined font type.
	Unknown SystemFontType = iota

	// TrueType represents a TrueType font (.ttf).
	TrueType

	// OpenType represents an OpenType font (.otf).
	OpenType

	// Type1 represents a PostScript Type 1 font.
	Type1

	// TrueTypeCollection represents a TrueType Collection file (.ttc).
	TrueTypeCollection

	// OpenTypeCollection represents an OpenType Collection file.
	OpenTypeCollection
)
