package annotations

import (
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// LinkAnnotationAdapter wraps *Annotation to implement content.LinkAnnotationIface.
// This adapter exists because Go does not support covariant return types for interface
// satisfaction — Annotation.Type() returns AnnotationType and AnnotationDictionary()
// returns tokens.Token, but LinkAnnotationIface requires exact signature matches.
type LinkAnnotationAdapter struct {
	a *Annotation
}

func (w LinkAnnotationAdapter) Type() any {
	// Return the actual AnnotationType enum to match C# behavior where annotation.Type
	// returns an AnnotationType enum. isLinkAnnotation handles this via fmt.Stringer case
	// since annotations.AnnotationType implements String().
	return w.a.Type()
}

func (w LinkAnnotationAdapter) AnnotationDictionary() tokens.Token {
	return w.a.AnnotationDictionary()
}

func (w LinkAnnotationAdapter) Rectangle() core.PdfRectangle {
	return w.a.Rectangle()
}

func (w LinkAnnotationAdapter) InReplyTo() any {
	if w.a.InReplyTo() == nil {
		return nil
	}
	return NewLinkAnnotationAdapter(w.a.InReplyTo())
}

func (w LinkAnnotationAdapter) Action() any {
	if w.a.Action() == nil {
		return nil
	}
	return w.a.Action()
}

func (w LinkAnnotationAdapter) HasNormalAppearance() bool {
	return w.a.HasNormalAppearance()
}

func (w LinkAnnotationAdapter) HasRollOverAppearance() bool {
	return w.a.HasRollOverAppearance()
}

func (w LinkAnnotationAdapter) HasDownAppearance() bool {
	return w.a.HasDownAppearance()
}

func (w LinkAnnotationAdapter) NormalAppearanceStream() any {
	if w.a.normalAppearanceStream == nil {
		return nil
	}
	return w.a.normalAppearanceStream
}

func (w LinkAnnotationAdapter) DownAppearanceStream() any {
	if w.a.downAppearanceStream == nil {
		return nil
	}
	return w.a.downAppearanceStream
}

func (w LinkAnnotationAdapter) AppearanceState() *string {
	return w.a.appearanceState
}

func (w LinkAnnotationAdapter) QuadPoints() any {
	if len(w.a.quadPoints) == 0 {
		return nil
	}
	return w.a.quadPoints
}

// NewLinkAnnotationAdapter creates a LinkAnnotationIface wrapper around *Annotation.
func NewLinkAnnotationAdapter(a *Annotation) content.LinkAnnotationIface {
	return LinkAnnotationAdapter{a: a}
}

// newLinkAnnotationAdapter is the unexported convenience for internal use.
func newLinkAnnotationAdapter(a *Annotation) content.LinkAnnotationIface {
	return LinkAnnotationAdapter{a: a}
}
