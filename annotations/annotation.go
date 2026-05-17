package annotations

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/actions"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Annotation represents an annotation on a page in a PDF document.
type Annotation struct {
	annotationDictionary    *tokens.DictionaryToken
	typeField               AnnotationType
	rectangle               core.PdfRectangle
	content                 *string
	name                    *string
	modifiedDate            *string
	flags                   AnnotationFlags
	border                  *AnnotationBorder
	quadPoints              []*QuadPointsQuadrilateral
	action                  actions.Action
	normalAppearanceStream  *AppearanceStream
	rollOverAppearanceStream *AppearanceStream
	downAppearanceStream    *AppearanceStream
	appearanceState         *string
	inReplyTo               *Annotation
}

// NewAnnotation creates a new Annotation.
// Pass nil for optional parameters (content, name, modifiedDate, action, appearance streams, appearanceState, inReplyTo).
// Pass nil for quadPoints if no quad points are present (will default to empty slice).
func NewAnnotation(
	annotationDictionary *tokens.DictionaryToken,
	typeField AnnotationType,
	rectangle core.PdfRectangle,
	content *string,
	name *string,
	modifiedDate *string,
	flags AnnotationFlags,
	border *AnnotationBorder,
	quadPoints []*QuadPointsQuadrilateral,
	action actions.Action,
	normalAppearanceStream *AppearanceStream,
	rollOverAppearanceStream *AppearanceStream,
	downAppearanceStream *AppearanceStream,
	appearanceState *string,
	inReplyTo *Annotation,
) *Annotation {
	if quadPoints == nil {
		quadPoints = []*QuadPointsQuadrilateral{}
	}

	return &Annotation{
		annotationDictionary:     annotationDictionary,
		typeField:                typeField,
		rectangle:                rectangle,
		content:                  content,
		name:                     name,
		modifiedDate:             modifiedDate,
		flags:                    flags,
		border:                   border,
		quadPoints:               quadPoints,
		action:                   action,
		normalAppearanceStream:   normalAppearanceStream,
		rollOverAppearanceStream: rollOverAppearanceStream,
		downAppearanceStream:     downAppearanceStream,
		appearanceState:          appearanceState,
		inReplyTo:                inReplyTo,
	}
}

// AnnotationDictionary returns the underlying PDF dictionary which this annotation was created from.
func (a *Annotation) AnnotationDictionary() tokens.Token {
	return a.annotationDictionary
}

// Type returns the type of this annotation.
func (a *Annotation) Type() AnnotationType {
	return a.typeField
}

// Rectangle returns the rectangle in user space units specifying the location to place this annotation on the page.
func (a *Annotation) Rectangle() core.PdfRectangle {
	return a.rectangle
}

// Content returns the annotation text, or if the annotation does not display text, a description of the annotation's contents. Returns nil if not present.
func (a *Annotation) Content() *string {
	return a.content
}

// Name returns the name of this annotation which should be unique per page. Returns nil if not present.
func (a *Annotation) Name() *string {
	return a.name
}

// ModifiedDate returns the date and time the annotation was last modified, can be in any format. Returns nil if not present.
func (a *Annotation) ModifiedDate() *string {
	return a.modifiedDate
}

// Flags returns the flags defining the appearance and behaviour of this annotation.
func (a *Annotation) Flags() AnnotationFlags {
	return a.flags
}

// Border returns the annotation's border definition.
func (a *Annotation) Border() *AnnotationBorder {
	return a.border
}

// QuadPoints returns the rectangles defined using QuadPoints. For Link annotations these are the regions used to activate the link,
// for text markup annotations these are the text regions to apply the markup to.
func (a *Annotation) QuadPoints() []*QuadPointsQuadrilateral {
	return a.quadPoints
}

// Action returns the action for this annotation, or nil if none is set.
func (a *Annotation) Action() actions.Action {
	return a.action
}

// HasNormalAppearance indicates whether a normal appearance is present for this annotation.
func (a *Annotation) HasNormalAppearance() bool {
	return a.normalAppearanceStream != nil
}

// HasRollOverAppearance indicates whether a roll over appearance is present for this annotation (shown when hovering).
func (a *Annotation) HasRollOverAppearance() bool {
	return a.rollOverAppearanceStream != nil
}

// HasDownAppearance indicates whether a down appearance is present for this annotation (shown when clicked).
func (a *Annotation) HasDownAppearance() bool {
	return a.downAppearanceStream != nil
}

// InReplyTo returns the annotation this annotation was in reply to, or nil if none.
func (a *Annotation) InReplyTo() *Annotation {
	return a.inReplyTo
}

// String returns a string representation of the annotation.
func (a *Annotation) String() string {
	return fmt.Sprintf("%v - %v", a.typeField, a.content)
}
