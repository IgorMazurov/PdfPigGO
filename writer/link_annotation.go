package writer

import (
	"github.com/uglytoad/pdfpig/go/annotations"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// BorderStyle specifies the border style for a link annotation.
type BorderStyle int

const (
	// Solid is a solid border.
	Solid BorderStyle = iota
	// Dashed is a dashed border.
	Dashed
	// Beveled is a simulated embossed border that appears to be raised above the surface of the page.
	Beveled
	// Inset is a simulated engraved border that appears to be recessed below the surface of the page.
	Inset
	// Underline is an underline border drawn along the bottom of the annotation rectangle.
	Underline
)

// LinkAnnotation represents a link annotation that can be added to a PDF page.
// Link annotations provide clickable areas that can trigger actions such as navigating to another page or opening a URL.
type LinkAnnotation struct {
	annotationBorder *annotations.AnnotationBorder
	border           *BorderStyle
	borderWidth      *int
	rect             core.PdfRectangle
	quadPoints       []annotations.QuadPointsQuadrilateral
	action           any // stores full concrete action type (UriAction, GoToAction, etc.)
}

// AnnotationBorder returns the border style for the link annotation.
// This is overwritten by Border if both are provided.
func (la LinkAnnotation) AnnotationBorder() *annotations.AnnotationBorder {
	return la.annotationBorder
}

// Border returns the border style for the link annotation.
func (la LinkAnnotation) Border() *BorderStyle {
	return la.border
}

// BorderWidth returns the width of the border for the link annotation.
func (la LinkAnnotation) BorderWidth() *int {
	return la.borderWidth
}

// Rect returns the rectangle defining the location and size of the link annotation on the page.
func (la LinkAnnotation) Rect() core.PdfRectangle {
	return la.rect
}

// QuadPoints returns the quadrilaterals defining the clickable regions of the link.
// These are typically used to define precise clickable areas that may not be rectangular.
func (la LinkAnnotation) QuadPoints() []annotations.QuadPointsQuadrilateral {
	if la.quadPoints == nil {
		return []annotations.QuadPointsQuadrilateral{}
	}
	return la.quadPoints
}

// Action returns the action to be performed when the link is activated.
func (la LinkAnnotation) Action() any {
	return la.action
}

// NewLinkAnnotation creates a new LinkAnnotation instance.
// Pass nil for optional parameters: annotationBorder, borderStyle, borderWidth, quadPoints.
func NewLinkAnnotation(
	action any,
	rect core.PdfRectangle,
	annotationBorder *annotations.AnnotationBorder,
	borderStyle *BorderStyle,
	borderWidth *int,
	quadPoints []annotations.QuadPointsQuadrilateral,
) *LinkAnnotation {
	if quadPoints == nil {
		quadPoints = []annotations.QuadPointsQuadrilateral{}
	}

	return &LinkAnnotation{
		annotationBorder: annotationBorder,
		border:           borderStyle,
		borderWidth:      borderWidth,
		rect:             rect,
		quadPoints:       quadPoints,
		action:           action,
	}
}

// ToToken converts this link annotation to a PDF dictionary token representation.
func (la LinkAnnotation) ToToken() *tokens.DictionaryToken {
	dictEntries := map[*tokens.NameToken]tokens.Token{
		tokens.Type:     tokens.Annot,
		tokens.Subtype:  tokens.Link,
		tokens.Rect:     tokens.NewArrayToken([]tokens.Token{
			tokens.NewNumericToken(la.rect.BottomLeft.X),
			tokens.NewNumericToken(la.rect.BottomLeft.Y),
			tokens.NewNumericToken(la.rect.TopRight.X),
			tokens.NewNumericToken(la.rect.TopRight.Y),
		}),
	}

	if len(la.quadPoints) > 0 {
		var quadPointsArray []tokens.Token
		for _, quad := range la.quadPoints {
			for _, point := range quad.Points() {
				quadPointsArray = append(quadPointsArray, tokens.NewNumericToken(point.X))
				quadPointsArray = append(quadPointsArray, tokens.NewNumericToken(point.Y))
			}
		}
		dictEntries[tokens.Quadpoints] = tokens.NewArrayToken(quadPointsArray)
	}

	if la.annotationBorder != nil {
		borderArray := []tokens.Token{
			tokens.NewNumericToken(la.annotationBorder.HorizontalCornerRadius()),
			tokens.NewNumericToken(la.annotationBorder.VerticalCornerRadius()),
			tokens.NewNumericToken(la.annotationBorder.BorderWidth()),
		}

		dashPattern := la.annotationBorder.LineDashPattern()
		if len(dashPattern) > 0 {
			var dashArray []tokens.Token
			for _, dash := range dashPattern {
				dashArray = append(dashArray, tokens.NewNumericToken(dash))
			}
			borderArray = append(borderArray, tokens.NewArrayToken(dashArray))
		}
		dictEntries[tokens.Border] = tokens.NewArrayToken(borderArray)
	}

	if la.border != nil {
		var styleName *tokens.NameToken
		switch *la.border {
		case Solid:
			styleName = tokens.S
		case Dashed:
			styleName = tokens.D
		case Beveled:
			styleName = tokens.B
		case Inset:
			styleName = tokens.I
		case Underline:
			styleName = tokens.U
		default:
			styleName = tokens.S
		}

		width := 1
		if la.borderWidth != nil {
			width = *la.borderWidth
		}

		bsDict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
			tokens.S: styleName,
			tokens.W: tokens.NewNumericTokenFromInt(width),
		})
		dictEntries[tokens.Bs] = bsDict
	}

	result, _ := tokens.NewDictionary(dictEntries)
	return result
}
