package alto

// AltoFontType represents font type (Serif or Sans-Serif) in ALTO format.
type AltoFontType string

const (
	// AltoFontTypeSerif indicates a serif font.
	AltoFontTypeSerif AltoFontType = "serif"
	// AltoFontTypeSansSerif indicates a sans-serif font.
	AltoFontTypeSansSerif AltoFontType = "sans-serif"
)

// String returns the XML string representation of this font type.
func (f AltoFontType) String() string {
	return string(f)
}
