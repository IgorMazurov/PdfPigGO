package content

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AnnotationType represents the type of a PDF annotation.
type AnnotationType int

const (
	AnnotationTypeUnknown AnnotationType = iota
	AnnotationTypeLink
)

// LinkAnnotationIface represents the minimal interface needed from an annotation
// for hyperlink extraction. This avoids an import cycle between content and annotations packages.
type LinkAnnotationIface interface {
	Type() any
	AnnotationDictionary() tokens.Token
	Rectangle() core.PdfRectangle
	InReplyTo() any
	Action() any
	HasNormalAppearance() bool
	HasRollOverAppearance() bool
	HasDownAppearance() bool
	NormalAppearanceStream() any
	DownAppearanceStream() any
	AppearanceState() *string
	QuadPoints() any
}

// Hyperlink represents full details for a link annotation which references an external resource.
type Hyperlink struct {
	// Bounds is the area on the page which when clicked will open the hyperlink.
	Bounds core.PdfRectangle
	// Text is the text in the link region (if any).
	Text string
	// Letters are the letters in the link region.
	Letters []*Letter
	// Uri is the URI the link directs to.
	Uri string
	// Annotation is the underlying link annotation.
	Annotation LinkAnnotationIface
}

// NewHyperlink creates a new Hyperlink.
func NewHyperlink(bounds core.PdfRectangle, letters []*Letter, text, uri string, annotation LinkAnnotationIface) *Hyperlink {
	return &Hyperlink{
		Bounds:     bounds,
		Text:       text,
		Letters:    letters,
		Uri:        uri,
		Annotation: annotation,
	}
}

// String returns a string representation of the hyperlink.
func (h *Hyperlink) String() string {
	return fmt.Sprintf("Link: %s (%s)", h.Text, h.Uri)
}
