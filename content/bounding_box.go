package content

import "github.com/uglytoad/pdfpig/go/core"

// BoundingBox is the interface for types that provide a bounding rectangle.
type BoundingBox interface {
	// BoundingBox returns the rectangle completely containing this object.
	BoundingBox() core.PdfRectangle
}
