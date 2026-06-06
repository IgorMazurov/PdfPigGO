package graphics

import (
	"fmt"
	"math"
	"time"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/util"
	"github.com/uglytoad/pdfpig/go/xobjects"
)

// ContentStreamProcessor processes PDF content streams to extract letters, paths, images,
// and marked content elements. It extends BaseStreamProcessor[PageContent].
type ContentStreamProcessor struct {
	*BaseStreamProcessor[*content.PageContent]

	letters            []*content.Letter
	paths              []geometry.PdfPath
	images             []core.Union[xobjects.XObjectContentRecord, *content.InlineImage]
	markedContents     []content.MarkedContentElement
	markedContentStack *MarkedContentStack
	filterProviderFP   filters.FilterProvider

	currentSubpath *core.PdfSubpath
	currentPath    *geometry.PdfPath
}

// NewContentStreamProcessor creates a new ContentStreamProcessor.
func NewContentStreamProcessor(
	pageNumber int,
	resourceStore content.ResourceStore,
	pdfScanner tokenization.PdfTokenScanner,
	pageContentParser content.PageContentParser,
	filterProvider content.LookupFilterProvider,
	cropBox *content.CropBox,
	userSpaceUnit geometry.UserSpaceUnit,
	rotation content.PageRotationDegrees,
	initialMatrix core.TransformationMatrix,
	parsingOptions *content.ParsingOptions,
) *ContentStreamProcessor {
	base := NewBaseStreamProcessor[*content.PageContent](
		pageNumber,
		resourceStore,
		pdfScanner,
		pageContentParser,
		filterProvider,
		cropBox,
		userSpaceUnit,
		rotation,
		initialMatrix,
		parsingOptions,
	)

	var fp filters.FilterProvider
	if x, ok := filterProvider.(filters.FilterProvider); ok {
		fp = x
	}

	csp := &ContentStreamProcessor{
		BaseStreamProcessor:  base,
		letters:              make([]*content.Letter, 0, 256),
		paths:                make([]geometry.PdfPath, 0, 64),
		images:               make([]core.Union[xobjects.XObjectContentRecord, *content.InlineImage], 0, 8),
		markedContents:       make([]content.MarkedContentElement, 0, 16),
		markedContentStack:   NewMarkedContentStack(),
		filterProviderFP:     fp,
	}
	csp.BaseStreamProcessor.renderGlyphHook = csp.RenderGlyph
	csp.BaseStreamProcessor.renderXObjectImageHook = csp.RenderXObjectImage
	csp.BaseStreamProcessor.renderInlineImageHook = csp.RenderInlineImage
	return csp
}

// CurrentSubpath returns the current subpath being built.
func (p *ContentStreamProcessor) CurrentSubpath() *core.PdfSubpath {
	return p.currentSubpath
}

// CurrentPath returns the current path being built.
func (p *ContentStreamProcessor) CurrentPath() *geometry.PdfPath {
	return p.currentPath
}

// Process processes the graphics operations and returns a PageContent result.
func (p *ContentStreamProcessor) Process(pageNumberCurrent int, operations []content.GraphicsStateOperation) *content.PageContent {
	p.PageNumber = pageNumberCurrent
	p.CloneAllStates()
	p.ProcessOperations(operations)

	pageContent, _ := content.NewPageContent(
		operations,
		p.letters,
		p.pathsToInterface(),
		p.images,
		p.markedContents,
		p.PdfScanner,
		p.filterProviderFP,
		p.ResourceStore,
	)
	return pageContent
}

// RenderGlyph renders a single glyph, handling diacritic combination and bounding box calculation.
func (p *ContentStreamProcessor) RenderGlyph(
	font fonts.Font,
	currentState *CurrentGraphicsState,
	fontSize float64,
	pointSize float64,
	code int,
	unicode string,
	currentOffset int64,
	renderingMatrix core.TransformationMatrix,
	textMatrix core.TransformationMatrix,
	transformationMatrix core.TransformationMatrix,
	characterBoundingBox *fonts.CharacterBoundingBox,
) {
	transformedGlyphBounds := Transform(
		renderingMatrix, textMatrix, transformationMatrix,
		characterBoundingBox.GlyphBounds,
	)

	if p.ParsingOptions.ClipPaths && currentState.CurrentClippingPath != nil {
		if !geometry.PathIntersectsWithRect(currentState.CurrentClippingPath, transformedGlyphBounds, false) {
			return
		}
	}

	var letter *content.Letter

	if util.IsInCombiningDiacriticRange(unicode) && currentOffset > 0 && len(p.letters) > 0 {
		attachTo := p.letters[len(p.letters)-1]
		if attachTo.TextSequence == p.TextSequence {
			if newLetter, ok := util.TryCombineDiacriticWithPreviousLetter(unicode, attachTo.Value); ok {
				p.letters = p.letters[:len(p.letters)-1]
				letter = content.NewLetter(
					newLetter,
					attachTo.BoundingBox,
					attachTo.GlyphRectangleLoose,
					attachTo.StartBaseLine,
					attachTo.EndBaseLine,
					attachTo.Width,
					attachTo.FontSize,
					attachTo.GetFont(),
					attachTo.RenderingMode,
					attachTo.StrokeColor,
					attachTo.FillColor,
					attachTo.PointSize,
					attachTo.TextSequence,
				)
				letter.Code = attachTo.Code
				letter.GlyphTransform = attachTo.GlyphTransform
			}
		}
	}

	isBboxValid := transformedGlyphBounds.Width > 1e-10 || transformedGlyphBounds.Height > 1e-10

	if letter == nil {
		transformedPdfBounds := Transform(
			renderingMatrix, textMatrix, transformationMatrix,
			core.NewPdfRectangleFloat(0, 0, characterBoundingBox.Width, float64(p.UserSpaceUnit.PointMultiples)),
		)

		looseBox := Transform(
			renderingMatrix, textMatrix, transformationMatrix,
			core.NewPdfRectangleFloat(0, font.GetDescent(), characterBoundingBox.Width, font.GetAscent()),
		)

		var bbox core.PdfRectangle
		if isBboxValid {
			bbox = transformedGlyphBounds
		} else {
			bbox = transformedPdfBounds
		}

		letter = content.NewLetter(
			unicode,
			bbox,
			looseBox,
			transformedPdfBounds.BottomLeft,
			transformedPdfBounds.BottomRight,
			transformedPdfBounds.Width,
			fontSize,
			font,
			currentState.FontState.TextRenderingMode,
			currentState.CurrentStrokingColor,
			currentState.CurrentNonStrokingColor,
			pointSize,
			p.TextSequence,
		)
		letter.Code = code
		letter.GlyphTransform = renderingMatrix.Multiply(textMatrix).Multiply(transformationMatrix)
	}

	p.letters = append(p.letters, letter)
	p.markedContentStack.AddLetter(letter)
}

// RenderXObjectImage records an XObject image reference and adds it to marked content.
func (p *ContentStreamProcessor) RenderXObjectImage(xObjectContentRecord *xobjects.XObjectContentRecord) {
	p.images = append(p.images, core.UnionOne[xobjects.XObjectContentRecord, *content.InlineImage](*xObjectContentRecord))
	p.markedContentStack.AddXObject(xObjectContentRecord, p.PdfScanner, p.filterProviderFP, p.ResourceStore)
}

// BeginSubpath starts a new subpath within the current path.
func (p *ContentStreamProcessor) BeginSubpath() {
	if p.currentPath == nil {
		p.currentPath = geometry.NewPdfPath()
	}
	p.addCurrentSubpath()
	p.currentSubpath = core.NewPdfSubpath()
}

// CloseSubpath closes the current subpath and returns the closing point.
func (p *ContentStreamProcessor) CloseSubpath() *core.PdfPoint {
	if p.currentSubpath == nil {
		return nil
	}

	cmds := p.currentSubpath.Commands()
	if len(cmds) == 0 {
		return nil
	}

	move, ok := cmds[0].(*core.Move)
	if !ok {
		p.ParsingOptions.Logger.Error("CloseSubpath(): first command not Move.")
		return nil
	}

	point := move.Location
	p.currentSubpath.CloseSubpath()
	p.addCurrentSubpath()
	return &point
}

// addCurrentSubpath adds the current subpath to the current path and resets it.
func (p *ContentStreamProcessor) addCurrentSubpath() {
	if p.currentSubpath == nil || p.currentPath == nil {
		return
	}
	p.currentPath.Add(p.currentSubpath)
	p.currentSubpath = nil
}

// StrokePath strokes the current path.
func (p *ContentStreamProcessor) StrokePath(close bool) {
	if p.currentPath == nil {
		return
	}

	p.currentPath.SetStroked()

	if close && p.currentSubpath != nil {
		p.currentSubpath.CloseSubpath()
	}

	p.ClosePath()
}

// FillPath fills the current path using the given filling rule.
func (p *ContentStreamProcessor) FillPath(fillingRule core.FillingRule, close bool) {
	if p.currentPath == nil {
		return
	}

	p.currentPath.SetFilled(fillingRule)

	if close && p.currentSubpath != nil {
		p.currentSubpath.CloseSubpath()
	}

	p.ClosePath()
}

// FillStrokePath fills then strokes the current path.
func (p *ContentStreamProcessor) FillStrokePath(fillingRule core.FillingRule, close bool) {
	if p.currentPath == nil {
		return
	}

	p.currentPath.SetFilled(fillingRule)
	p.currentPath.SetStroked()

	if close && p.currentSubpath != nil {
		p.currentSubpath.CloseSubpath()
	}

	p.ClosePath()
}

// MoveTo moves to the specified coordinates, starting a new subpath.
func (p *ContentStreamProcessor) MoveTo(x, y float64) {
	p.BeginSubpath()
	point := p.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x, y))
	p.SetCurrentPosition(point)
	p.currentSubpath.MoveTo(point.X, point.Y)
}

// BezierCurveToQuadratic adds a quadratic Bézier curve segment.
func (p *ContentStreamProcessor) BezierCurveToQuadratic(x2, y2, x3, y3 float64) {
	if p.currentSubpath == nil {
		return
	}

	controlPoint2 := p.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x2, y2))
	end := p.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x3, y3))

	p.currentSubpath.BezierCurveToQuadratic(
		p.CurrentPosition().X, p.CurrentPosition().Y,
		controlPoint2.X, controlPoint2.Y,
	)
	p.SetCurrentPosition(end)
}

// BezierCurveToCubic adds a cubic Bézier curve segment.
func (p *ContentStreamProcessor) BezierCurveToCubic(x1, y1, x2, y2, x3, y3 float64) {
	if p.currentSubpath == nil {
		return
	}

	controlPoint1 := p.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x1, y1))
	controlPoint2 := p.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x2, y2))
	end := p.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x3, y3))

	p.currentSubpath.BezierCurveToCubic(
		controlPoint1.X, controlPoint1.Y,
		controlPoint2.X, controlPoint2.Y,
		end.X, end.Y,
	)
	p.SetCurrentPosition(end)
}

// BezierCurveToStartCP adds a cubic Bézier curve segment using the current point as first control point ("v" operator).
func (p *ContentStreamProcessor) BezierCurveToStartCP(x2, y2, x3, y3 float64) {
	if p.currentSubpath == nil {
		return
	}

	current := p.CurrentTransformationMatrix().TransformPoint(p.CurrentPosition())
	controlPoint2 := p.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x2, y2))
	end := p.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x3, y3))

	p.currentSubpath.BezierCurveToCubic(
		current.X, current.Y,
		controlPoint2.X, controlPoint2.Y,
		end.X, end.Y,
	)
	p.SetCurrentPosition(end)
}

// LineTo adds a line segment to the current point.
func (p *ContentStreamProcessor) LineTo(x, y float64) {
	if p.currentSubpath == nil {
		return
	}

	endPoint := p.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x, y))
	p.currentSubpath.LineTo(endPoint.X, endPoint.Y)
	p.SetCurrentPosition(endPoint)
}

// Rectangle adds a rectangle subpath.
func (p *ContentStreamProcessor) Rectangle(x, y, width, height float64) {
	p.BeginSubpath()
	lowerLeft := p.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x, y))
	upperRight := p.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x+width, y+height))

	p.currentSubpath.Rectangle(lowerLeft.X, lowerLeft.Y, upperRight.X-lowerLeft.X, upperRight.Y-lowerLeft.Y)
	p.addCurrentSubpath()
}

// EndPath ends the current path construction without filling or stroking.
func (p *ContentStreamProcessor) EndPath() {
	if p.currentPath == nil {
		return
	}

	p.addCurrentSubpath()

	if p.currentPath.IsClipping() {
		if !p.ParsingOptions.ClipPaths {
			pathCopy := *p.currentPath
			p.paths = append(p.paths, pathCopy)
			p.markedContentStack.AddPath(&pathCopy)
		}
		p.currentPath = nil
		return
	}

	pathCopy := *p.currentPath
	p.paths = append(p.paths, pathCopy)
	p.markedContentStack.AddPath(&pathCopy)
	p.currentPath = nil
}

// ClosePath closes and ends the current path.
func (p *ContentStreamProcessor) ClosePath() {
	p.addCurrentSubpath()

	if p.currentPath == nil {
		return
	}

	if p.currentPath.IsClipping() {
		p.EndPath()
		return
	}

	currentState := p.GetCurrentState()
	if p.currentPath.IsStroked() {
		p.currentPath.SetStrokeDetails(geometry.StrokeState{
			LineDashPattern:      currentState.LineDashPattern,
			CurrentStrokingColor: currentState.CurrentStrokingColor,
			LineWidth:            currentState.LineWidth,
			CapStyle:             currentState.CapStyle,
			JoinStyle:            currentState.JoinStyle,
		})
	}

	if p.currentPath.IsFilled() {
		p.currentPath.SetFillDetails(geometry.FillState{
			CurrentNonStrokingColor: currentState.CurrentNonStrokingColor,
		})
	}

	if p.ParsingOptions.ClipPaths && currentState != nil && currentState.CurrentClippingPath != nil {
		clippedPath := currentState.CurrentClippingPath.Clip(p.currentPath, p.ParsingOptions.Logger)
		if clippedPath != nil {
			p.paths = append(p.paths, *clippedPath)
			p.markedContentStack.AddPath(clippedPath)
		} else if p.currentPath.IsStroked() || p.currentPath.IsFilled() {
			// Fallback: if clipping returns empty for a stroked/filled path, add the original.
			// This matches C# behavior where paths within the crop box are preserved even when
			// the clipper library fails to compute the intersection (e.g., open/stroked paths).
			pathCopy := *p.currentPath
			p.paths = append(p.paths, pathCopy)
			p.markedContentStack.AddPath(&pathCopy)
		}
	} else {
		pathCopy := *p.currentPath
		p.paths = append(p.paths, pathCopy)
		p.markedContentStack.AddPath(&pathCopy)
	}

	p.currentPath = nil
}

// ModifyClippingIntersect modifies the clipping path by intersecting with the current path.
func (p *ContentStreamProcessor) ModifyClippingIntersect(clippingRule core.FillingRule) {
	if p.currentPath == nil {
		return
	}

	p.addCurrentSubpath()
	p.currentPath.SetClipping(clippingRule)

	if p.ParsingOptions.ClipPaths {
		graphicsState := p.GetCurrentState()
		if graphicsState != nil && graphicsState.CurrentClippingPath != nil {
			currentClipping := graphicsState.CurrentClippingPath
			currentClipping.SetClipping(clippingRule)

			newClippings := clipWithTimeout(p.currentPath, currentClipping, p.ParsingOptions.Logger)
			if newClippings == nil {
				p.ParsingOptions.Logger.Warn("Empty clipping path found. Clipping path not updated.")
			} else {
				graphicsState.CurrentClippingPath = newClippings
			}
		}
	}
}

// ClipToRectangle clips to the given rectangle with the specified filling rule.
func (p *ContentStreamProcessor) ClipToRectangle(rectangle core.PdfRectangle, clippingRule core.FillingRule) {
	graphicsState := p.GetCurrentState()
	if graphicsState == nil {
		return
	}

	transformedRect := graphicsState.CurrentTransformationMatrix.TransformRect(rectangle)

	// Guard against degenerate transformed rectangles (NaN/Inf) that would hang the clipper.
	if !isFiniteRect(transformedRect) {
		p.ParsingOptions.Logger.Warn("Degenerate clipping rectangle encountered (transformed), skipping clip.")
		return
	}

	normalizedRect := geometry.RectangleNormalise(transformedRect)
	clipPath := geometry.RectangleToPdfPath(normalizedRect)
	clipPath.SetClipping(clippingRule)

	if graphicsState.CurrentClippingPath != nil {
		currentClipping := graphicsState.CurrentClippingPath
		currentClipping.SetClipping(clippingRule)
		newClippings := clipPath.Clip(currentClipping, p.ParsingOptions.Logger)

		if newClippings == nil {
			p.ParsingOptions.Logger.Warn("Empty clipping path found. Clipping path not updated.")
		} else {
			graphicsState.CurrentClippingPath = newClippings
		}
	}
}

// isFiniteRect checks whether all coordinates of a rectangle are finite numbers.
func isFiniteRect(r core.PdfRectangle) bool {
	for _, pt := range []core.PdfPoint{r.TopLeft, r.TopRight, r.BottomLeft, r.BottomRight} {
		if math.IsNaN(pt.X) || math.IsInf(pt.X, 0) || math.IsNaN(pt.Y) || math.IsInf(pt.Y, 0) {
			return false
		}
	}
	return true
}

// RenderInlineImage records an inline image and adds it to marked content.
func (p *ContentStreamProcessor) RenderInlineImage(image *content.InlineImage) {
	p.images = append(p.images, core.UnionTwo[xobjects.XObjectContentRecord, *content.InlineImage](image))
	p.markedContentStack.AddImage(image)
}

// BeginMarkedContent starts a marked content section.
func (p *ContentStreamProcessor) BeginMarkedContent(name *tokens.NameToken, propertyDictionaryName *tokens.NameToken, properties *tokens.DictionaryToken) {
	if propertyDictionaryName != nil {
		actual := p.ResourceStore.GetMarkedContentPropertiesDictionary(propertyDictionaryName)
		if actual != nil {
			properties = actual
		}
	}

	p.markedContentStack.Push(name, properties)
}

// EndMarkedContent ends the current marked content section.
func (p *ContentStreamProcessor) EndMarkedContent() {
	if p.markedContentStack.CanPop() {
		mc := p.markedContentStack.Pop(p.PdfScanner)
		if mc != nil {
			p.markedContents = append(p.markedContents, *mc)
		}
	}
}

// PaintShading paints a shading fill (no-op for now).
func (p *ContentStreamProcessor) PaintShading(shadingName *tokens.NameToken) {
	_ = shadingName
}

// pathsToInterface converts geometry.PdfPath values to content.PdfPath wrappers.
func (p *ContentStreamProcessor) pathsToInterface() []content.PdfPath {
	result := make([]content.PdfPath, len(p.paths))
	for i := range p.paths {
		result[i] = content.PdfPath{PdfPath: &p.paths[i]}
	}
	return result
}

// ApplyXObject overrides BaseStreamProcessor so that Form XObject processing
// dispatches through *ContentStreamProcessor (not the embedded base field).
func (p *ContentStreamProcessor) ApplyXObject(xObjectName *tokens.NameToken) {
	xObjectStream, found := p.ResourceStore.TryGetXObject(xObjectName)
	if !found {
		if p.ParsingOptions.SkipMissingFonts {
			return
		}
		panic(fmt.Errorf("No XObject with name %v found on page %d.", xObjectName, p.PageNumber))
	}

	subTypeToken := xObjectStream.StreamDictionary.Data()[tokens.Subtype.Data()]
	subType, ok := subTypeToken.(*tokens.NameToken)
	if !ok {
		panic(fmt.Errorf("XObject Subtype is not a NameToken: %v", subTypeToken))
	}

	state := p.GetCurrentState()
	matrix := state.CurrentTransformationMatrix

	colorSpace := colors.DeviceRgbColorSpaceDetails
	if state.ColorSpaceContext != nil && state.ColorSpaceContext.CurrentStrokingColorSpace != nil {
		colorSpace = state.ColorSpaceContext.CurrentStrokingColorSpace
	}

	if subType.Data() == tokens.Ps.Data() {
		contentRecord := xobjects.NewXObjectContentRecord(
			xObjectName,
			xObjectStream,
			xobjects.PostScript,
			matrix,
			state.RenderingIntent,
			colorSpace,
		)
		p.xObjects[xobjects.PostScript] = append(p.xObjects[xobjects.PostScript], contentRecord)
	} else if subType.Data() == tokens.Image.Data() {
		contentRecord := xobjects.NewXObjectContentRecord(
			xObjectName,
			xObjectStream,
			xobjects.Image,
			matrix,
			state.RenderingIntent,
			colorSpace,
		)
		if p.renderXObjectImageHook != nil {
			p.renderXObjectImageHook(contentRecord)
		}
	} else if subType.Data() == tokens.Form.Data() {
		p.ProcessFormXObject(xObjectStream, xObjectName)
	} else {
		panic(fmt.Errorf("XObject encountered with unexpected SubType %v. %v.", subType, xObjectStream.StreamDictionary))
	}
}

// ProcessOperations overrides BaseStreamProcessor to pass *ContentStreamProcessor
// as OperationContext so that overridden methods (StrokePath, FillPath, etc.) are reached.
func (p *ContentStreamProcessor) ProcessOperations(operations []content.GraphicsStateOperation) {
	for _, op := range operations {
		if r, ok := op.(runnableOperation); ok {
			r.Run(p)
		}
	}
}

// ProcessFormXObject overrides BaseStreamProcessor to ensure form XObject content
// is processed through *ContentStreamProcessor (not the embedded base field), so that
// overridden path methods are invoked correctly. This matches C# virtual dispatch behavior.
func (p *ContentStreamProcessor) ProcessFormXObject(formStream *tokens.StreamToken, xObjectName *tokens.NameToken) {
	var formResources *tokens.DictionaryToken
	if dictRaw, ok := formStream.StreamDictionary.Data()[tokens.Resources.Data()]; ok {
		formResources = resolveToDictionary(dictRaw, p.PdfScanner)
	}

	if formResources != nil {
		p.ResourceStore.LoadResourceDictionary(formResources)
	}

	// 1. Save current state.
	p.PushState()

	startState := p.GetCurrentState()

	// Transparency Group XObjects - check for /Group entry
	if groupToken, ok := formStream.StreamDictionary.Data()[tokens.Group.Data()]; ok {
		formGroupDict := resolveToDictionary(groupToken, p.PdfScanner)

		if formGroupDict != nil {
			sTokenRaw, hasS := formGroupDict.Data()[tokens.S.Data()]
			if !hasS {
				panic(fmt.Errorf("Invalid Transparency Group XObject, '/S' token is not set."))
			}
			sToken, ok := sTokenRaw.(*tokens.NameToken)
			if !ok || sToken.Data() != tokens.Transparency.Data() {
				panic(fmt.Errorf("Invalid Transparency Group XObject, '/S' token is not equal to '/Transparency'."))
			}

			startState.BlendMode = graphiccore.Normal
			startState.SoftMask = nil
			startState.AlphaConstantNonStroking = 1.0
			startState.AlphaConstantStroking = 1.0
		}
	}

	formMatrix := core.Identity
	if matrixToken, ok := formStream.StreamDictionary.Data()[tokens.Matrix.Data()]; ok {
		if arr, isArray := matrixToken.(*tokens.ArrayToken); isArray {
			data := arr.Data()
			if len(data) >= 6 {
				vals := make([]float64, 0, len(data))
				for _, item := range data {
					if num, isNum := item.(*tokens.NumericToken); isNum {
						vals = append(vals, num.Data())
					}
				}
				if len(vals) >= 6 {
					if m, err := core.FromArray(vals[:6]); err == nil {
						formMatrix = m
					}
				}
			}
		}
	}

	// 2. Update current transformation matrix.
	p.ModifyCurrentTransformationMatrix(formMatrix)

	contentStream := p.FilterProvider.DecodeStream(formStream, p.PdfScanner)

	operations := p.PageContentParser.Parse(
		p.PageNumber,
		core.NewMemoryInputBytes(contentStream),
		p.ParsingOptions.Logger,
	)

	// 3. Clip according to the form dictionary's BBox entry.
	if bboxToken, ok := formStream.StreamDictionary.Data()[tokens.Bbox.Data()]; ok {
		if arr, isArray := bboxToken.(*tokens.ArrayToken); isArray {
			data := arr.Data()
			if len(data) >= 4 {
				points := make([]float64, 0, len(data))
				for _, item := range data {
					if num, isNum := item.(*tokens.NumericToken); isNum {
						points = append(points, num.Data())
					}
				}
				if len(points) >= 4 {
					bboxRect := core.NewPdfRectangleFloat(points[0], points[1], points[2], points[3])
					normalizedBbox := geometry.RectangleNormalise(bboxRect)
					p.ClipToRectangle(normalizedBbox, core.FillingRuleEvenOdd)
				}
			}
		}
	}

	// 4. Check for circular references and paint objects.
	hasCircularReference := p.HasFormXObjectCircularReference(formStream, xObjectName, operations)
	if hasCircularReference {
		if p.ParsingOptions.UseLenientParsing {
			filteredOps := make([]content.GraphicsStateOperation, 0, len(operations))
			for _, op := range operations {
				if invoke, ok := op.(*InvokeNamedXObject); ok && invoke.Name != nil && invoke.Name.Data() == xObjectName.Data() {
					continue
				}
				filteredOps = append(filteredOps, op)
			}
			operations = filteredOps
			p.ParsingOptions.Logger.Warn(
				fmt.Sprintf("An XObject form named '%v' is referencing itself which can cause unexpected behaviour. The self reference was removed from the operations before further processing.", xObjectName))
		} else {
			panic(fmt.Errorf("An XObject form named '%v' is referencing itself which can cause unexpected behaviour.", xObjectName))
		}
	}

	p.ProcessOperations(operations)

	// 5. Restore saved state.
	p.PopState()

	if formResources != nil {
		p.ResourceStore.UnloadResourceDictionary()
	}
}

// clipWithTimeout runs PdfPath.Clip with a timeout to prevent hangs from the
// clipper library on certain path combinations. Returns nil if clipping times out.
func clipWithTimeout(subject, clipping *geometry.PdfPath, logger logging.Log) *geometry.PdfPath {
	resultCh := make(chan *geometry.PdfPath, 1)
	go func() {
		resultCh <- subject.Clip(clipping, logger)
	}()

	select {
	case result := <-resultCh:
		return result
	case <-time.After(60 * time.Second):
		logger.Warn("Clipping operation timed out (60s), returning original subject path.")
		return subject
	}
}
