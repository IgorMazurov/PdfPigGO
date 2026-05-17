package content

import (
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/images/png"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// XObjectImage represents an image XObject extracted from a PDF document.
type XObjectImage struct {
	boundingBox       core.PdfRectangle
	widthInSamples    int
	heightInSamples   int
	bitsPerComponent  int
	isJpxDecode       bool
	isImageMask       bool
	renderingIntent   graphiccore.RenderingIntent
	interpolate       bool
	decode            []float64
	imageDictionary   *tokens.DictionaryToken
	rawMemory         []byte
	colorSpaceDetails colors.ColorSpaceDetails
	maskImage         PdfImage

	filterProvider filters.FilterProvider
	stream         *tokens.StreamToken
	pdfScanner     tokenization.PdfTokenScanner

	decodedOnce  sync.Once
	decoded    []byte
	decodeOk   bool
	supported  bool
}

var _ PdfImage = (*XObjectImage)(nil)

// NewXObjectImage creates a new XObjectImage.
func NewXObjectImage(
	bounds core.PdfRectangle,
	widthInSamples int,
	heightInSamples int,
	bitsPerComponent int,
	isJpxDecode bool,
	isImageMask bool,
	renderingIntent graphiccore.RenderingIntent,
	interpolate bool,
	decode []float64,
	imageDictionary *tokens.DictionaryToken,
	rawMemory []byte,
	colorSpaceDetails colors.ColorSpaceDetails,
	maskImage PdfImage,
	filterProvider filters.FilterProvider,
	stream *tokens.StreamToken,
	supported bool,
) *XObjectImage {
	return &XObjectImage{
		boundingBox:       bounds,
		widthInSamples:    widthInSamples,
		heightInSamples:   heightInSamples,
		bitsPerComponent:  bitsPerComponent,
		isJpxDecode:       isJpxDecode,
		isImageMask:       isImageMask,
		renderingIntent:   renderingIntent,
		interpolate:       interpolate,
		decode:            decode,
		imageDictionary:   imageDictionary,
		rawMemory:         rawMemory,
		colorSpaceDetails: colorSpaceDetails,
		maskImage:         maskImage,
		filterProvider:    filterProvider,
		stream:            stream,
		supported:         supported,
	}
}

// SetPdfScanner sets the PDF token scanner for resolving indirect references during decoding.
func (i *XObjectImage) SetPdfScanner(scanner tokenization.PdfTokenScanner) {
	i.pdfScanner = scanner
}

// BoundingBox returns the bounding rectangle of the image in PDF coordinates.
func (i *XObjectImage) BoundingBox() core.PdfRectangle {
	return i.boundingBox
}

// Bounds returns the placement rectangle of the image in PDF coordinates.
// Deprecated: use BoundingBox() instead.
func (i *XObjectImage) Bounds() core.PdfRectangle {
	return i.boundingBox
}

// WidthInSamples returns the width of the image in samples.
func (i *XObjectImage) WidthInSamples() int {
	return i.widthInSamples
}

// HeightInSamples returns the height of the image in samples.
func (i *XObjectImage) HeightInSamples() int {
	return i.heightInSamples
}

// BitsPerComponent returns the number of bits used to represent each color component.
func (i *XObjectImage) BitsPerComponent() int {
	return i.bitsPerComponent
}

// IsImageMask reports whether the image is to be treated as an image mask.
func (i *XObjectImage) IsImageMask() bool {
	return i.isImageMask
}

// Decode returns the decode array that maps image samples into values
// appropriate for the color space.
func (i *XObjectImage) Decode() []float64 {
	return i.decode
}

// Interpolate reports whether interpolation should be performed when rendering.
func (i *XObjectImage) Interpolate() bool {
	return i.interpolate
}

// IsInlineImage reports whether this image is an inline image.
// XObject images are never inline, so this always returns false.
func (i *XObjectImage) IsInlineImage() bool {
	return false
}

// ImageDictionary returns the full dictionary token for this image object.
func (i *XObjectImage) ImageDictionary() *tokens.DictionaryToken {
	return i.imageDictionary
}

// RenderingIntent returns the color rendering intent for the image.
func (i *XObjectImage) RenderingIntent() graphiccore.RenderingIntent {
	return i.renderingIntent
}

// RawMemory returns the encoded bytes of the image with all filters still applied.
func (i *XObjectImage) RawMemory() []byte {
	return i.rawMemory
}

// RawBytes returns the encoded byte slice of the image with all filters still applied.
func (i *XObjectImage) RawBytes() []byte {
	return i.rawMemory
}

// ColorSpaceDetails returns the color space details used to interpret the image.
func (i *XObjectImage) ColorSpaceDetails() *colors.ColorSpaceDetails {
	return &i.colorSpaceDetails
}

// MaskImage returns the image mask, either a soft-mask or a stencil mask.
func (i *XObjectImage) MaskImage() PdfImage {
	return i.maskImage
}

// TryGetBytes attempts to return the decoded bytes of the image.
// Returns false if any filter in the chain is not supported.
func (i *XObjectImage) TryGetBytes() ([]byte, bool) {
	if !i.supported {
		return nil, false
	}

	i.decodedOnce.Do(func() {
		if i.stream == nil || i.filterProvider == nil {
			return
		}
		fl, err := getFiltersForStream(i.filterProvider, i.stream.StreamDictionary, i.pdfScanner)
		if err != nil {
			return
		}
		transform := i.stream.Data()
		for idx, f := range fl {
			transform, err = f.Decode(transform, i.stream.StreamDictionary, i.filterProvider, idx)
			if err != nil {
				return
			}
		}
		i.decoded = transform
		i.decodeOk = true
	})

	if !i.decodeOk {
		return nil, false
	}

	return i.decoded, true
}

func getFiltersForStream(fp filters.FilterProvider, dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) ([]filters.Filter, error) {
	if lp, ok := fp.(filters.LookupFilterProvider); ok && scanner != nil {
		return lp.GetFiltersWithScanner(dict, scanner)
	}
	return fp.GetFilters(dict)
}

// TryGetPng attempts to convert the image to PNG format.
func (i *XObjectImage) TryGetPng() ([]byte, bool) {
	return png.TryGenerate(i, i.maskImage)
}
