package testutil

import (
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/images/png"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// TestPdfImage is a test helper that implements content.PdfImage.
// It provides mutable fields for easy configuration in tests.
type TestPdfImage struct {
	boundingBox       core.PdfRectangle
	widthInSamples    int
	heightInSamples   int
	bitsPerComponent  int
	rawMemory         []byte
	renderingIntent   graphiccore.RenderingIntent
	isImageMask       bool
	decode            []float64
	interpolate       bool
	isInlineImage     bool
	imageDictionary   *tokens.DictionaryToken
	colorSpaceDetails *colors.ColorSpaceDetails
	decodedBytes      []byte
	maskImage         content.PdfImage
}

var _ content.PdfImage = (*TestPdfImage)(nil)

// NewTestPdfImage creates a new TestPdfImage with default values.
func NewTestPdfImage() *TestPdfImage {
	return &TestPdfImage{
		bitsPerComponent: 8,
		renderingIntent:  graphiccore.RelativeColorimetric,
	}
}

// BoundingBox returns the bounding rectangle of the image in PDF coordinates.
func (t *TestPdfImage) BoundingBox() core.PdfRectangle {
	return t.boundingBox
}

// SetBoundingBox sets the bounding rectangle of the image.
func (t *TestPdfImage) SetBoundingBox(rect core.PdfRectangle) {
	t.boundingBox = rect
}

// Bounds returns the placement rectangle of the image in PDF coordinates.
// Deprecated: use BoundingBox() instead.
func (t *TestPdfImage) Bounds() core.PdfRectangle {
	return t.boundingBox
}

// WidthInSamples returns the width of the image in samples.
func (t *TestPdfImage) WidthInSamples() int {
	return t.widthInSamples
}

// SetWidthInSamples sets the width of the image in samples.
func (t *TestPdfImage) SetWidthInSamples(w int) {
	t.widthInSamples = w
}

// HeightInSamples returns the height of the image in samples.
func (t *TestPdfImage) HeightInSamples() int {
	return t.heightInSamples
}

// SetHeightInSamples sets the height of the image in samples.
func (t *TestPdfImage) SetHeightInSamples(h int) {
	t.heightInSamples = h
}

// BitsPerComponent returns the number of bits used to represent each color component.
func (t *TestPdfImage) BitsPerComponent() int {
	return t.bitsPerComponent
}

// SetBitsPerComponent sets the number of bits per component.
func (t *TestPdfImage) SetBitsPerComponent(bpc int) {
	t.bitsPerComponent = bpc
}

// RawMemory returns the encoded bytes of the image with all filters still applied.
func (t *TestPdfImage) RawMemory() []byte {
	return t.rawMemory
}

// SetRawMemory sets the raw encoded bytes of the image.
func (t *TestPdfImage) SetRawMemory(b []byte) {
	t.rawMemory = b
}

// RawBytes returns the encoded byte slice of the image with all filters still applied.
func (t *TestPdfImage) RawBytes() []byte {
	return t.rawMemory
}

// RenderingIntent returns the color rendering intent for the image.
func (t *TestPdfImage) RenderingIntent() graphiccore.RenderingIntent {
	return t.renderingIntent
}

// SetRenderingIntent sets the color rendering intent for the image.
func (t *TestPdfImage) SetRenderingIntent(ri graphiccore.RenderingIntent) {
	t.renderingIntent = ri
}

// IsImageMask reports whether the image is to be treated as an image mask.
func (t *TestPdfImage) IsImageMask() bool {
	return t.isImageMask
}

// SetIsImageMask sets whether the image should be treated as an image mask.
func (t *TestPdfImage) SetIsImageMask(v bool) {
	t.isImageMask = v
}

// Decode returns the decode array that maps image samples into values
// appropriate for the color space.
func (t *TestPdfImage) Decode() []float64 {
	return t.decode
}

// SetDecode sets the decode array for the image.
func (t *TestPdfImage) SetDecode(d []float64) {
	t.decode = d
}

// Interpolate reports whether interpolation should be performed when rendering.
func (t *TestPdfImage) Interpolate() bool {
	return t.interpolate
}

// SetInterpolate sets whether interpolation should be performed.
func (t *TestPdfImage) SetInterpolate(v bool) {
	t.interpolate = v
}

// IsInlineImage reports whether this image is an inline image.
func (t *TestPdfImage) IsInlineImage() bool {
	return t.isInlineImage
}

// SetIsInlineImage sets whether this image is an inline image.
func (t *TestPdfImage) SetIsInlineImage(v bool) {
	t.isInlineImage = v
}

// ImageDictionary returns the full dictionary token for this image object.
func (t *TestPdfImage) ImageDictionary() *tokens.DictionaryToken {
	return t.imageDictionary
}

// SetImageDictionary sets the dictionary token for this image object.
func (t *TestPdfImage) SetImageDictionary(d *tokens.DictionaryToken) {
	t.imageDictionary = d
}

// ColorSpaceDetails returns the color space details used to interpret the image.
func (t *TestPdfImage) ColorSpaceDetails() *colors.ColorSpaceDetails {
	return t.colorSpaceDetails
}

// SetColorSpaceDetails sets the color space details for the image.
func (t *TestPdfImage) SetColorSpaceDetails(csd *colors.ColorSpaceDetails) {
	t.colorSpaceDetails = csd
}

// SetDecodedBytes sets the decoded bytes of the image.
func (t *TestPdfImage) SetDecodedBytes(b []byte) {
	t.decodedBytes = b
}

// MaskImage returns the image mask, either a soft-mask or a stencil mask.
func (t *TestPdfImage) MaskImage() content.PdfImage {
	return t.maskImage
}

// TryGetBytes attempts to return the decoded bytes of the image.
func (t *TestPdfImage) TryGetBytes() ([]byte, bool) {
	if len(t.decodedBytes) == 0 {
		return nil, false
	}
	return t.decodedBytes, true
}

// TryGetPng attempts to convert the image to PNG format.
func (t *TestPdfImage) TryGetPng() ([]byte, bool) {
	return png.TryGenerate(t, t.maskImage)
}
