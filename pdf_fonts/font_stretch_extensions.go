package pdffonts

import (
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// ConvertToFontStretch converts a NameToken to the corresponding FontStretch value.
func ConvertToFontStretch(name *tokens.NameToken) fonts.FontStretch {
	switch name.Data() {
	case "UltraCondensed":
		return fonts.UltraCondensed
	case "ExtraCondensed":
		return fonts.ExtraCondensed
	case "Condensed":
		return fonts.Condensed
	case "Normal":
		return fonts.Normal
	case "SemiExpanded":
		return fonts.SemiExpanded
	case "Expanded":
		return fonts.Expanded
	case "ExtraExpanded":
		return fonts.ExtraExpanded
	case "UltraExpanded":
		return fonts.UltraExpanded
	default:
		return fonts.Unknown
	}
}
