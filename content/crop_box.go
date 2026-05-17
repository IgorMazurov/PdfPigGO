package content

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
)

// CropBox defines the visible region of a page, contents expanding beyond the crop box should be clipped.
type CropBox struct {
	// Bounds defines the clipping of the content when the page is displayed or printed.
	// The page's contents are to be clipped (cropped) to this rectangle and then imposed on the output medium.
	Bounds core.PdfRectangle
}

// NewCropBox creates a new CropBox with the given bounds.
func NewCropBox(bounds core.PdfRectangle) *CropBox {
	return &CropBox{Bounds: bounds}
}

// String returns the string representation of the crop box bounds.
func (c *CropBox) String() string {
	return c.Bounds.String()
}

var _ fmt.Stringer = (*CropBox)(nil)
