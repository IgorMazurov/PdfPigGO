package annotations

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/actions"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/util"
)

// AnnotationProvider extracts annotation objects from a page's /Annots array.
type AnnotationProvider struct {
	tokenScanner      tokenization.PdfTokenScanner
	pageDictionary    *tokens.DictionaryToken
	namedDestinations *destinations.NamedDestinations
	log               logging.Log
	matrix            core.TransformationMatrix
}

// NewAnnotationProvider creates a new AnnotationProvider.
func NewAnnotationProvider(
	tokenScanner tokenization.PdfTokenScanner,
	pageDictionary *tokens.DictionaryToken,
	matrix core.TransformationMatrix,
	namedDestinations *destinations.NamedDestinations,
	log logging.Log,
) (*AnnotationProvider, error) {
	if tokenScanner == nil {
		return nil, errors.New("tokenScanner cannot be null")
	}
	if pageDictionary == nil {
		return nil, errors.New("pageDictionary cannot be null")
	}

	return &AnnotationProvider{
		tokenScanner:      tokenScanner,
		pageDictionary:    pageDictionary,
		namedDestinations: namedDestinations,
		log:               log,
		matrix:            matrix,
	}, nil
}

// GetAnnotations returns all annotations defined on the page.
func (p *AnnotationProvider) GetAnnotations() []*Annotation {
	lookupAnnotations := make(map[core.IndirectReference]*Annotation)

	annotToken, ok := p.pageDictionary.TryGet(tokens.Annots)
	if !ok {
		return nil
	}

	annotationsArray, ok := parts.TryGet[*tokens.ArrayToken](annotToken, p.tokenScanner)
	if !ok {
		return nil
	}

	var result []*Annotation

	for _, tokenItem := range annotationsArray.Data() {
		annotationDict, ok := parts.TryGet[*tokens.DictionaryToken](tokenItem, p.tokenScanner)
		if !ok {
			continue
		}

		replyTo := p.resolveReplyTo(annotationDict, lookupAnnotations)

		typeName, _ := parts.TryGet[*tokens.NameToken](
			func() tokens.Token { t, _ := annotationDict.TryGet(tokens.Subtype); return t }(),
			p.tokenScanner,
		)
		annotationType := ToAnnotationType(typeName)

		action := p.getAction(annotationDict)

		rectToken, rectOk := parts.TryGet[*tokens.ArrayToken](
			func() tokens.Token { t, _ := annotationDict.TryGet(tokens.Rect); return t }(),
			p.tokenScanner,
		)
		var rectangle core.PdfRectangle
		if rectOk {
			if rect, err := util.ToRectangle(rectToken, p.tokenScanner); err == nil && rect != nil {
				rectangle = p.matrix.TransformRect(*rect)
			}
		}

		content := p.getNamedString(tokens.Contents, annotationDict)
		name := p.getNamedString(tokens.Nm, annotationDict)
		modifiedDate := p.getNamedString(tokens.M, annotationDict)

		flags := AnnotationFlags(0)
		if flagsToken, ok := annotationDict.TryGet(tokens.F); ok {
			if numeric, found := parts.TryGet[*tokens.NumericToken](flagsToken, p.tokenScanner); found {
				flags = AnnotationFlags(numeric.IntVal())
			}
		}

		border := Default
		if borderToken, ok := annotationDict.TryGet(tokens.Border); ok {
			if borderArray, found := parts.TryGet[*tokens.ArrayToken](borderToken, p.tokenScanner); found && borderArray.Length() >= 3 {
				horizontal := getNumericValue(borderArray.Get(0))
				vertical := getNumericValue(borderArray.Get(1))
				width := getNumericValue(borderArray.Get(2))

				var dashes []float64
				if borderArray.Length() == 4 {
					if dashArray, ok := borderArray.Data()[3].(*tokens.ArrayToken); ok {
						dashes = extractDashes(dashArray)
					}
				}

				border = NewAnnotationBorder(horizontal, vertical, width, dashes)
			}
		}

		quadPointRectangles := p.extractQuadPoints(annotationDict)

		var normalAppearanceStream *AppearanceStream
		var rollOverAppearanceStream *AppearanceStream
		var downAppearanceStream *AppearanceStream

		if apToken, ok := annotationDict.TryGet(tokens.Ap); ok {
			if appearanceDictionary, found := parts.TryGet[*tokens.DictionaryToken](apToken, p.tokenScanner); found {
				streamN, _ := TryCreateAppearanceStream(appearanceDictionary, tokens.N, p.tokenScanner)
				normalAppearanceStream = streamN

				streamR, _ := TryCreateAppearanceStream(appearanceDictionary, tokens.R, p.tokenScanner)
				rollOverAppearanceStream = streamR

				streamD, _ := TryCreateAppearanceStream(appearanceDictionary, tokens.D, p.tokenScanner)
				downAppearanceStream = streamD
			}
		}

		appearanceState := p.extractAppearanceState(annotationDict)

		annotation := NewAnnotation(
			annotationDict,
			annotationType,
			rectangle,
			content,
			name,
			modifiedDate,
			flags,
			border,
			quadPointRectangles,
			action,
			normalAppearanceStream,
			rollOverAppearanceStream,
			downAppearanceStream,
			appearanceState,
			replyTo,
		)

		if indirectRef, ok := tokenItem.(*tokens.IndirectReferenceToken); ok {
			lookupAnnotations[indirectRef.Data()] = annotation
		}

		result = append(result, annotation)
	}

	return result
}

// resolveReplyTo resolves the /Irt (In Reply To) reference if present.
func (p *AnnotationProvider) resolveReplyTo(dict *tokens.DictionaryToken, lookup map[core.IndirectReference]*Annotation) *Annotation {
	if irtToken, ok := dict.TryGet(tokens.Irt); ok {
		refToken, found := parts.TryGet[*tokens.IndirectReferenceToken](irtToken, p.tokenScanner)
		if found {
			if linked, exists := lookup[refToken.Data()]; exists {
				return linked
			}
		}
	}
	return nil
}

// getAction retrieves the action for an annotation dictionary.
func (p *AnnotationProvider) getAction(annotationDictionary *tokens.DictionaryToken) actions.Action {
	destProvider := destinations.DestinationProvider{}

	if dest, ok := destProvider.TryGetDestination(
		annotationDictionary, tokens.Dest, p.namedDestinations, p.tokenScanner, p.log, false); ok {
		return actions.NewGoToAction(dest)
	}

	actionProvider := actions.ActionProvider{}

	action, ok, err := actionProvider.TryGetAction(
		annotationDictionary, p.namedDestinations, p.tokenScanner, p.log)
	if err != nil {
		return nil
	}
	if ok {
		return action
	}

	return nil
}

// getNamedString extracts a string value from the dictionary for the given key.
func (p *AnnotationProvider) getNamedString(name *tokens.NameToken, dictionary *tokens.DictionaryToken) *string {
	token, ok := dictionary.TryGet(name)
	if !ok {
		return nil
	}

	switch t := token.(type) {
	case *tokens.StringToken:
		s := t.Data()
		return &s
	case *tokens.HexToken:
		s := t.Data()
		return &s
	default:
		if str, found := parts.TryGet[*tokens.StringToken](token, p.tokenScanner); found {
			s := str.Data()
			return &s
		}
	}

	return nil
}

// extractQuadPoints parses the /Quadpoints array into quadrilaterals.
func (p *AnnotationProvider) extractQuadPoints(dict *tokens.DictionaryToken) []*QuadPointsQuadrilateral {
	qpToken, ok := dict.TryGet(tokens.Quadpoints)
	if !ok {
		return nil
	}

	quadPointsArray, found := parts.TryGet[*tokens.ArrayToken](qpToken, p.tokenScanner)
	if !found {
		return nil
	}

	var quadPointRectangles []*QuadPointsQuadrilateral
	values := make([]float64, 0, 8)

	for i := 0; i < quadPointsArray.Length(); i++ {
		numeric, ok := quadPointsArray.Get(i).(*tokens.NumericToken)
		if !ok {
			continue
		}

		values = append(values, numeric.Data())

		if len(values) == 8 {
			qp, err := NewQuadPointsQuadrilateral([]core.PdfPoint{
				p.matrix.TransformPoint(core.NewPdfPoint(values[0], values[1])),
				p.matrix.TransformPoint(core.NewPdfPoint(values[2], values[3])),
				p.matrix.TransformPoint(core.NewPdfPoint(values[4], values[5])),
				p.matrix.TransformPoint(core.NewPdfPoint(values[6], values[7])),
			})
			if err == nil {
				quadPointRectangles = append(quadPointRectangles, qp)
			}
			values = values[:0]
		}
	}

	return quadPointRectangles
}

// extractAppearanceState returns the current appearance state name if present.
func (p *AnnotationProvider) extractAppearanceState(dict *tokens.DictionaryToken) *string {
	asToken, ok := dict.TryGet(tokens.As)
	if !ok {
		return nil
	}

	nameToken, found := parts.TryGet[*tokens.NameToken](asToken, p.tokenScanner)
	if !found {
		return nil
	}

	s := nameToken.Data()
	return &s
}

// getNumericValue extracts the float64 data from a NumericToken, or returns 0 if not numeric.
func getNumericValue(token tokens.Token) float64 {
	if n, ok := token.(*tokens.NumericToken); ok {
		return n.Data()
	}
	return 0
}

// extractDashes extracts dash pattern values from an array of NumericTokens.
func extractDashes(arr *tokens.ArrayToken) []float64 {
	data := arr.Data()
	dashes := make([]float64, 0, len(data))
	for _, item := range data {
		if n, ok := item.(*tokens.NumericToken); ok {
			dashes = append(dashes, n.Data())
		}
	}
	return dashes
}

// unwrapIndirect resolves an IndirectReferenceToken through the scanner.
func unwrapIndirect(token tokens.Token, scanner tokenization.PdfTokenScanner) tokens.Token {
	if indRef, ok := token.(*tokens.IndirectReferenceToken); ok && scanner != nil {
		ref := indRef.Data()
		obj := scanner.Get(ref)
		if obj != nil {
			return obj.Data()
		}
	}
	return token
}
