package alto

import "strings"

// AltoFontStyles represents font styles in ALTO XML format as a bitmask.
type AltoFontStyles uint32

const (
	// AltoFontStylesBold indicates bold text.
	AltoFontStylesBold AltoFontStyles = 1 << iota
	// AltoFontStylesItalics indicates italicized text.
	AltoFontStylesItalics
	// AltoFontStylesSubscript indicates subscripted text.
	AltoFontStylesSubscript
	// AltoFontStylesSuperscript indicates superscripted text.
	AltoFontStylesSuperscript
	// AltoFontStylesSmallCaps indicates small-caps text.
	AltoFontStylesSmallCaps
	// AltoFontStylesUnderline indicates underlined text.
	AltoFontStylesUnderline
)

// String returns a comma-separated list of the XML enum values for each set flag.
func (f AltoFontStyles) String() string {
	if f == 0 {
		return ""
	}

	var parts []string
	if f&AltoFontStylesBold != 0 {
		parts = append(parts, "bold")
	}
	if f&AltoFontStylesItalics != 0 {
		parts = append(parts, "italics")
	}
	if f&AltoFontStylesSubscript != 0 {
		parts = append(parts, "subscript")
	}
	if f&AltoFontStylesSuperscript != 0 {
		parts = append(parts, "superscript")
	}
	if f&AltoFontStylesSmallCaps != 0 {
		parts = append(parts, "smallcaps")
	}
	if f&AltoFontStylesUnderline != 0 {
		parts = append(parts, "underline")
	}

	return strings.Join(parts, ",")
}
