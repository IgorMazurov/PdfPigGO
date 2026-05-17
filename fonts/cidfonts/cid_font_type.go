package cidfonts

// CidFontType specifies the type of a CID font in a PDF document.
type CidFontType int

const (
	// Type0 represents glyph descriptions based on Adobe Type 1 format.
	Type0 CidFontType = iota

	// Type2 represents glyph descriptions based on TrueType format.
	Type2 CidFontType = 2
)
