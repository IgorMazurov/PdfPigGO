package colors

// ColorSpace represents the color space used to interpret color values in a PDF.
// Color spaces enable a PDF to specify abstract colors in a device independent way.
type ColorSpace int

const (
	// DeviceGray controls the intensity of achromatic light on a scale from black to white.
	DeviceGray ColorSpace = 0

	// DeviceRGB controls the intensities of red, green and blue light.
	DeviceRGB = 1

	// DeviceCMYK controls the concentrations of cyan, magenta, yellow and black (K) inks.
	DeviceCMYK = 2

	// CalGray is a special case of the CIE colorspace using a single channel (A) and a single transformation.
	// A represents the gray component of a calibrated gray space in the range 0 to 1.
	CalGray = 3

	// CalRGB is a CIE ABC color space with a single transformation.
	// A, B and C represent red, green and blue color values in the range 0 to 1.
	CalRGB = 4

	// Lab is a CIE ABC color space with two transforms. A, B and C represent the L*, a* and b*
	// components of a CIE 1976 L*a*b* space. The range of A (L*) is 0 to 100.
	// The range of B (a*) and C (b*) are defined by the Range of the color space.
	Lab = 5

	// ICCBased is a colorspace specified by a sequence of bytes which are interpreted according to
	// the ICC specification.
	ICCBased = 6

	// Indexed allows a PDF content stream to use small integers as indices into a color map or color table
	// of arbitrary colors in some other space. Each sample value is treated as an index into the color table.
	Indexed = 7

	// Pattern enables a PDF content stream to paint an area with a pattern rather than a single color.
	// The pattern may be either a tiling pattern (type 1) or a shading pattern (type 2).
	Pattern = 8

	// Separation provides a means for specifying the use of additional colorants or for isolating the control
	// of individual color components of a device color space for a subtractive device. When such a space is
	// the current color space, the current color is a single-component value called a tint that controls the
	// application of the given colorant or color components only.
	Separation = 9

// DeviceN can contain an arbitrary number of color components and provides greater flexibility than
	// standard device color spaces such as DeviceCMYK or individual Separation color spaces for a subtractive device.
	DeviceN = 10
)

// String returns the name of the color space.
func (cs ColorSpace) String() string {
	switch cs {
	case DeviceGray:
		return "DeviceGray"
	case DeviceRGB:
		return "DeviceRGB"
	case DeviceCMYK:
		return "DeviceCMYK"
	case CalGray:
		return "CalGray"
	case CalRGB:
		return "CalRGB"
	case Lab:
		return "Lab"
	case ICCBased:
		return "ICCBased"
	case Indexed:
		return "Indexed"
	case Pattern:
		return "Pattern"
	case Separation:
		return "Separation"
	case DeviceN:
		return "DeviceN"
	default:
		return "Unknown"
	}
}
