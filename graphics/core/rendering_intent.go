package core

// RenderingIntent specifies priorities regarding which properties to preserve
// and which to sacrifice for CIE colors.
type RenderingIntent byte

const (
	// AbsoluteColorimetric means no correction for the output medium's white
	// point. Colors only represented relative to the light source.
	AbsoluteColorimetric RenderingIntent = iota

	// RelativeColorimetric combines light source and output medium's white point.
	RelativeColorimetric

	// RenderSaturation emphasises saturation rather than colorimetric accuracy.
	RenderSaturation

	// Perceptual modifies from colorimetric values to provide a pleasing
	// perceptual appearance.
	Perceptual
)

var renderingIntentMap = map[string]RenderingIntent{
	"AbsoluteColorimetric": AbsoluteColorimetric,
	"RelativeColorimetric": RelativeColorimetric,
	"Saturation":           RenderSaturation,
	"Perceptual":           Perceptual,
}

// ParseRenderingIntent converts a string to the corresponding RenderingIntent.
// If the name is not recognized, RelativeColorimetric is returned as default per
// PDF specification.
func ParseRenderingIntent(s string) RenderingIntent {
	if intent, ok := renderingIntentMap[s]; ok {
		return intent
	}

	return RelativeColorimetric
}
