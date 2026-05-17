package colors

// ColorSpaceFamily classifies ColorSpaces into families.
// ColorSpaces within the same family share general characteristics.
type ColorSpaceFamily int

const (
	// Device represents colorspaces that directly specify colors or shades of gray
	// that the output device should produce.
	Device ColorSpaceFamily = iota

	// CIEBased represents color spaces based on an international standard for color
	// specification created by the Commission Internationale de l'Eclairage
	// (International Commission on Illumination). These spaces specify colors in a way
	// that is independent of the characteristics of any particular output device.
	CIEBased

	// Special represents color spaces that add features or properties to an underlying
	// color space. They include facilities for patterns, color mapping, separations,
	// and high-fidelity and multitone color.
	Special
)
