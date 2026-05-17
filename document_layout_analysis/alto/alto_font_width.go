package alto

// AltoFontWidth represents font width in ALTO format.
type AltoFontWidth string

const (
	// AltoFontWidthProportional indicates a proportional-width font.
	AltoFontWidthProportional AltoFontWidth = "proportional"
	// AltoFontWidthFixed indicates a fixed-width font.
	AltoFontWidthFixed AltoFontWidth = "fixed"
)

// String returns the XML string representation of this font width.
func (f AltoFontWidth) String() string {
	return string(f)
}
