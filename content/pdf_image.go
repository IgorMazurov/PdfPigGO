package content

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// PdfImage represents an image extracted from a PDF document.
// It may be an inline image or a PostScript image XObject.
type PdfImage interface {
	BoundingBox

	// Bounds returns the placement rectangle of the image in PDF coordinates.
	// Deprecated: use BoundingBox() instead.
	Bounds() core.PdfRectangle

	// WidthInSamples returns the width of the image in samples.
	WidthInSamples() int

	// HeightInSamples returns the height of the image in samples.
	HeightInSamples() int

	// BitsPerComponent returns the number of bits used to represent each color component.
	BitsPerComponent() int

	// RawMemory returns the encoded bytes of the image with all filters still applied.
	RawMemory() []byte

	// RawBytes returns the encoded byte slice of the image with all filters still applied.
	RawBytes() []byte

	// RenderingIntent returns the color rendering intent for the image.
	RenderingIntent() graphiccore.RenderingIntent

	// IsImageMask reports whether the image is to be treated as an image mask.
	// If true, the image is a monochrome image where each sample is a single bit
	// and represents a stencil for marking or masking regions on the page.
	IsImageMask() bool

	// Decode returns the decode array that maps image samples into values
	// appropriate for the ColorSpace. Each pair of numbers in the array defines
	// the interpolation range for one component.
	Decode() []float64

	// Interpolate reports whether interpolation should be performed when rendering.
	Interpolate() bool

	// IsInlineImage reports whether this image is an inline image as opposed to
	// an XObject image.
	IsInlineImage() bool

	// ImageDictionary returns the full dictionary token for this image object.
	ImageDictionary() *tokens.DictionaryToken

	// ColorSpaceDetails returns the color space details used to interpret the image.
	// Returns nil when IsImageMask is true or when the image is JPX encoded.
	ColorSpaceDetails() *colors.ColorSpaceDetails

	// MaskImage returns the image mask, either a soft-mask or a stencil mask.
	MaskImage() PdfImage

	// TryGetBytes attempts to return the decoded bytes of the image.
	// For JPEG images and some other types, RawMemory should be used directly.
	TryGetBytes() ([]byte, bool)

	// TryGetPng attempts to convert the image to PNG format.
	// Does not support conversion of JPG to PNG.
	TryGetPng() ([]byte, bool)
}
