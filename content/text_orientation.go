package content

// TextOrientation represents the orientation of text in a PDF document.
type TextOrientation byte

const (
	// OtherTextOrientation indicates an unrecognized or other text orientation.
	OtherTextOrientation TextOrientation = 0

	// HorizontalTextOrientation is the usual left-to-right text orientation.
	HorizontalTextOrientation TextOrientation = 1

	// Rotate180TextOrientation is horizontal text, upside down.
	Rotate180TextOrientation TextOrientation = 2

	// Rotate90TextOrientation is rotated text going down.
	Rotate90TextOrientation TextOrientation = 3

	// Rotate270TextOrientation is rotated text going up.
	Rotate270TextOrientation TextOrientation = 4
)
