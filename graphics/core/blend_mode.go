package core

// BlendMode defines how overlapping graphical objects are composited.
type BlendMode byte

const (
	// Normal is the default blend mode, equivalent to Compatible.
	Normal BlendMode = iota

	// Multiply multiplies the top and bottom pixel colors and produces a darker
	// result.
	Multiply

	// Screen inverts the top and bottom pixel colors, multiplies them, then
	// inverts the result, producing a lighter image.
	Screen

	// Darken examines the color information in all the layers and selects the
	// darkest color as the result.
	Darken

	// Lighten examines the color information in all the layers and selects the
	// lightest color as the result.
	Lighten

	// ColorDodge brightens the bottom color to reflect the top color.
	ColorDodge

	// ColorBurn darkens the bottom color to reflect the top color.
	ColorBurn

	// HardLight combines Multiply and Screen depending on the top color.
	HardLight

	// SoftLight similar to HardLight but softer.
	SoftLight

	// Overlay combines Multiply and Screen depending on the bottom color.
	Overlay

	// Difference subtracts the darker from the lighter of the two colors.
	Difference

	// Exclusion similar to Difference but with lower contrast.
	Exclusion

	// Hue preserves the luminance and saturation of the bottom color while
	// adopting the hue of the top color (non-separable).
	Hue

	// Saturation preserves the luminance and hue of the bottom color while
	// adopting the saturation of the top color (non-separable).
	Saturation

	// Color preserves the luminance of the bottom color while adopting the hue
	// and saturation of the top color (non-separable).
	Color

	// Luminosity preserves the hue and saturation of the bottom color while
	// adopting the luminance of the top color (non-separable).
	Luminosity
)

var blendModeMap = map[string]BlendMode{
	"Normal":     Normal,
	"Compatible": Normal,
	"Multiply":   Multiply,
	"Screen":     Screen,
	"Darken":     Darken,
	"Lighten":    Lighten,
	"ColorDodge": ColorDodge,
	"ColorBurn":  ColorBurn,
	"HardLight":  HardLight,
	"SoftLight":  SoftLight,
	"Overlay":    Overlay,
	"Difference": Difference,
	"Exclusion":  Exclusion,
	"Hue":        Hue,
	"Saturation": Saturation,
	"Color":      Color,
	"Luminosity": Luminosity,
}

// ParseBlendMode converts a string to the corresponding BlendMode.
// Returns (mode, true) if recognized, or (0, false) otherwise.
func ParseBlendMode(s string) (BlendMode, bool) {
	mode, ok := blendModeMap[s]
	return mode, ok
}
