package png

// ColorType describes which color channels are present in a PNG image.
type ColorType byte

const (
	// ColorTypeNone indicates no color information.
	ColorTypeNone ColorType = 0

	// ColorTypePaletteUsed indicates that the PLTE chunk is used for indexed colors.
	ColorTypePaletteUsed ColorType = 1

	// ColorTypeColorUsed indicates that red, green, and blue channels are present.
	ColorTypeColorUsed ColorType = 2

	// ColorTypeAlphaChannelUsed indicates that an alpha channel is present.
	ColorTypeAlphaChannelUsed ColorType = 4
)
