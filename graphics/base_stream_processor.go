// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"math"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/xobjects"
)

// OperationContext represents the context in which graphics operations are executed.
type OperationContext interface {
	GetCurrentState() *CurrentGraphicsState
	CurrentPosition() core.PdfPoint
	SetCurrentPosition(pos core.PdfPoint)
	StackSize() int
	PopState()
	PushState()
	ShowText(bytes core.InputBytes)
	ShowPositionedText(tokenList []tokens.Token)
	ApplyXObject(xObjectName *tokens.NameToken)
	SetNamedGraphicsState(stateName *tokens.NameToken)
	BeginInlineImage()
	SetInlineImageProperties(properties map[*tokens.NameToken]tokens.Token)
	EndInlineImage(bytes []byte)
	BeginMarkedContent(name *tokens.NameToken, propertyDictionaryName *tokens.NameToken, properties *tokens.DictionaryToken)
	EndMarkedContent()
	SetFlatnessTolerance(tolerance float64)
SetLineCap(cap graphiccore.LineCapStyle)
	SetLineDashPattern(pattern graphiccore.LineDashPattern)
SetLineJoin(join graphiccore.LineJoinStyle)
	SetLineWidth(width float64)
	SetMiterLimit(limit float64)
	MoveToNextLineWithOffset()
	SetFontAndSize(font *tokens.NameToken, size float64)
	SetHorizontalScaling(scale float64)
	SetTextLeading(leading float64)
	SetTextRenderingMode(mode core.TextRenderingMode)
	SetTextRise(rise float64)
	SetWordSpacing(spacing float64)
	ModifyCurrentTransformationMatrix(value core.TransformationMatrix)
	SetCharacterSpacing(spacing float64)
	PaintShading(shadingName *tokens.NameToken)
	BeginSubpath()
	CloseSubpath() *core.PdfPoint
	StrokePath(close bool)
	FillPath(fillingRule core.FillingRule, close bool)
	FillStrokePath(fillingRule core.FillingRule, close bool)
	MoveTo(x, y float64)
	BezierCurveToQuadratic(x2, y2, x3, y3 float64)
	BezierCurveToCubic(x1, y1, x2, y2, x3, y3 float64)
	BezierCurveToStartCP(x2, y2, x3, y3 float64)
	LineTo(x, y float64)
	Rectangle(x, y, width, height float64)
	EndPath()
	ClosePath()
	ModifyClippingIntersect(clippingRule core.FillingRule)
	TextMatrices() *TextMatrices
	CurrentTransformationMatrix() core.TransformationMatrix
}

// formProcessingContext is the interface used by ProcessFormXObject to dispatch
// to concrete-type overrides (e.g., ContentStreamProcessor.ClipToRectangle).
// Go embedding does not provide virtual dispatch: when a method on
// BaseStreamProcessor[T] calls sp.SomeMethod(), it always resolves to the
// BaseStreamProcessor implementation even if an outer type defines an override.
// This interface bridges that gap by letting the concrete processor register itself.
type formProcessingContext interface {
	ClipToRectangle(rectangle core.PdfRectangle, clippingRule core.FillingRule)
}

// BaseStreamProcessor provides shared logic for processing PDF content streams.
// TPageContent is the type of page content produced by concrete implementations.
type BaseStreamProcessor[TPageContent any] struct {
	ResourceStore           content.ResourceStore
	UserSpaceUnit           geometry.UserSpaceUnit
	Rotation                content.PageRotationDegrees
	PdfScanner              tokenization.PdfTokenScanner
	PageContentParser       content.PageContentParser
	FilterProvider          content.LookupFilterProvider
	ParsingOptions          *content.ParsingOptions
	GraphicsStack           []*CurrentGraphicsState
	ActiveExtendedGraphicsStateFont fonts.Font
	InlineImageBuilder      *InlineImageBuilder
	PageNumber              int
	TextSequence            int
	textMatrices            *TextMatrices
	currentPosition         core.PdfPoint
	formCtx                 formProcessingContext
	renderGlyphHook         func(font fonts.Font, currentState *CurrentGraphicsState, fontSize float64, pointSize float64, code int, unicode string, currentOffset int64, renderingMatrix core.TransformationMatrix, textMatrix core.TransformationMatrix, transformationMatrix core.TransformationMatrix, characterBoundingBox *fonts.CharacterBoundingBox)
	renderXObjectImageHook  func(record *xobjects.XObjectContentRecord)
	renderInlineImageHook   func(image *content.InlineImage)
	xObjects                map[xobjects.XObjectType][]*xobjects.XObjectContentRecord
}

// SetFormProcessingContext sets the dispatch context for polymorphic method calls
// inside ProcessFormXObject. Concrete processors (e.g., ContentStreamProcessor)
// should call this with themselves so that overridden methods are reached.
func (sp *BaseStreamProcessor[TPageContent]) SetFormProcessingContext(ctx formProcessingContext) {
	sp.formCtx = ctx
}

// NewBaseStreamProcessor creates a new BaseStreamProcessor with the given configuration.
func NewBaseStreamProcessor[TPageContent any](
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
) *BaseStreamProcessor[TPageContent] {
	if pdfScanner == nil {
		panic(fmt.Errorf("pdfScanner cannot be nil"))
	}
	if pageContentParser == nil {
		panic(fmt.Errorf("pageContentParser cannot be nil"))
	}
	if filterProvider == nil {
		panic(fmt.Errorf("filterProvider cannot be nil"))
	}

	initialClipping := GetInitialClipping(cropBox)

	sp := &BaseStreamProcessor[TPageContent]{
		PageNumber:        pageNumber,
		ResourceStore:     resourceStore,
		UserSpaceUnit:     userSpaceUnit,
		Rotation:          rotation,
		PdfScanner:        pdfScanner,
		PageContentParser: pageContentParser,
		FilterProvider:    filterProvider,
		ParsingOptions:    parsingOptions,
		textMatrices:      NewTextMatrices(),
		xObjects: map[xobjects.XObjectType][]*xobjects.XObjectContentRecord{
			xobjects.Image:      {},
			xobjects.PostScript: {},
		},
	}

	sp.GraphicsStack = []*CurrentGraphicsState{{
		CurrentTransformationMatrix: initialMatrix,
		CurrentClippingPath:         initialClipping,
		ColorSpaceContext:           NewColorSpaceContext(sp.getCurrentStateFunc(), resourceStore),
		FontState: FontState{
			HorizontalScaling: 100,
		},
	}}

	return sp
}

// getCurrentStateFunc returns a closure that retrieves the current graphics state.
func (sp *BaseStreamProcessor[TPageContent]) getCurrentStateFunc() func() *CurrentGraphicsState {
	return func() *CurrentGraphicsState {
		if len(sp.GraphicsStack) == 0 {
			return nil
		}
		return sp.GraphicsStack[len(sp.GraphicsStack)-1]
	}
}

// GetInitialClipping creates a clipping path from the crop box bounds.
func GetInitialClipping(cropBox *content.CropBox) *geometry.PdfPath {
	clippingPath := geometry.RectangleToPdfPath(cropBox.Bounds)
	clippingPath.SetClipping(core.FillingRuleEvenOdd)
	return clippingPath
}

// Process is the abstract method that concrete implementations must provide.
func (sp *BaseStreamProcessor[TPageContent]) Process(_ int, _ []content.GraphicsStateOperation) TPageContent {
	var zero TPageContent
	return zero
}

// runnableOperation is the interface implemented by all concrete graphics operations
// that can be executed against an OperationContext. This avoids a circular dependency
// between content.GraphicsStateOperation and graphics.OperationContext.
type runnableOperation interface {
	Run(ctx OperationContext)
}

// ProcessOperations runs each graphics state operation against this context.
func (sp *BaseStreamProcessor[TPageContent]) ProcessOperations(operations []content.GraphicsStateOperation) {
	for _, op := range operations {
		if r, ok := op.(runnableOperation); ok {
			r.Run(sp)
		}
	}
}

// CloneAllStates clones the entire graphics state stack and pushes a clone of the top state.
func (sp *BaseStreamProcessor[TPageContent]) CloneAllStates() []*CurrentGraphicsState {
	saved := sp.GraphicsStack
	sp.GraphicsStack = []*CurrentGraphicsState{}
	if len(saved) > 0 {
		sp.GraphicsStack = append(sp.GraphicsStack, saved[len(saved)-1].DeepClone())
	}
	return saved
}

// GetCurrentState returns the current (top-of-stack) graphics state.
func (sp *BaseStreamProcessor[TPageContent]) GetCurrentState() *CurrentGraphicsState {
	if len(sp.GraphicsStack) == 0 {
		return nil
	}
	return sp.GraphicsStack[len(sp.GraphicsStack)-1]
}

// PopState pops the top graphics state from the stack, keeping at least one state.
func (sp *BaseStreamProcessor[TPageContent]) PopState() {
	if len(sp.GraphicsStack) > 1 {
		sp.GraphicsStack = sp.GraphicsStack[:len(sp.GraphicsStack)-1]
	} else {
		errMsg := "Cannot execute a pop of the graphics state stack, it would leave the stack empty."
		sp.ParsingOptions.Logger.Error(errMsg)
		if !sp.ParsingOptions.UseLenientParsing {
			panic(fmt.Errorf("%s", errMsg))
		}
	}
	sp.ActiveExtendedGraphicsStateFont = nil
}

// PushState pushes a deep clone of the current graphics state onto the stack.
func (sp *BaseStreamProcessor[TPageContent]) PushState() {
	current := sp.GetCurrentState()
	if current != nil {
		sp.GraphicsStack = append(sp.GraphicsStack, current.DeepClone())
	}
}

// ShowText renders text from the given byte stream using the current font.
func (sp *BaseStreamProcessor[TPageContent]) ShowText(bytes core.InputBytes) {
	sp.TextSequence++

	currentState := sp.GetCurrentState()
	if currentState == nil {
		return
	}

	var font fonts.Font
	if currentState.FontState.FromExtendedGraphicsState {
		font = sp.ActiveExtendedGraphicsStateFont
	} else {
		font = sp.ResourceStore.GetFont(currentState.FontState.FontName)
	}

	if font == nil {
		if sp.ParsingOptions.SkipMissingFonts {
			sp.ParsingOptions.Logger.Warn(
				fmt.Sprintf("Skipping a missing font with name %v since it is not present in the document and SkipMissingFonts is set to true.", currentState.FontState.FontName))
			return
		}
		panic(fmt.Sprintf("Could not find the font with name %v in the resource store. It has not been loaded yet.", currentState.FontState.FontName))
	}

	fontSize := currentState.FontState.FontSize
	horizontalScaling := currentState.FontState.HorizontalScaling / 100.0
	characterSpacing := currentState.FontState.CharacterSpacing
	rise := currentState.FontState.Rise

	transformationMatrix := currentState.CurrentTransformationMatrix

	renderingMatrix := core.FromValues(fontSize*horizontalScaling, 0, 0, fontSize, 0, rise)

	combined := transformationMatrix.Multiply(sp.textMatrices.TextMatrix)
	transformedRect := combined.TransformRect(core.NewPdfRectangleFloat(0, 0, 1, fontSize))
	pointSize := math.Round(transformedRect.Height*100) / 100

	for bytes.MoveNext() {
		savedTextMatrix := sp.textMatrices.TextMatrix
		if err := sp.processShowTextChar(font, bytes, currentState, fontSize, pointSize, characterSpacing, horizontalScaling, renderingMatrix, transformationMatrix); err != nil {
			sp.textMatrices.TextMatrix = savedTextMatrix
			sp.ParsingOptions.Logger.Error(fmt.Sprintf("Failed to render character: %v", err))
		}
	}
}

// processShowTextChar renders one character from ShowText operation.
// Returns nil on success, error if the character could not be rendered.
// No defer/recover — errors are returned explicitly for zero-overhead hot path.
func (sp *BaseStreamProcessor[TPageContent]) processShowTextChar(
	font fonts.Font, bytes core.InputBytes, currentState *CurrentGraphicsState,
	fontSize, pointSize, characterSpacing, horizontalScaling float64,
	renderingMatrix, transformationMatrix core.TransformationMatrix) error {
	code, codeLength := font.ReadCharacterCode(bytes)

	unicode, foundUnicode := font.TryGetUnicode(code)

	if !foundUnicode || unicode == "" {
		sp.ParsingOptions.Logger.Warn(
			fmt.Sprintf("We could not find the corresponding character with code %d in font %v.", code, font.Name()))
		unicode = string(rune(code))
	}

	wordSpacing := 0.0
	if code == ' ' && codeLength == 1 {
		wordSpacing += sp.GetCurrentState().FontState.WordSpacing
	}

	textMatrix := sp.textMatrices.TextMatrix

	if font.IsVertical() {
		if verticalFont, ok := font.(IVerticalWritingSupported); ok {
			positionVector := verticalFont.GetPositionVector(code)
			textMatrix = core.NewTransformationMatrixFromSlice([]float64{1, 0, 0, 1, positionVector.X, positionVector.Y}).Multiply(textMatrix)
		} else {
			return fmt.Errorf("font %v was in vertical writing mode but did not implement IVerticalWritingSupported", font.Name())
		}
	}

	boundingBox, err := font.GetBoundingBox(code)
	if err != nil {
		return err
	}

	if sp.renderGlyphHook != nil {
		sp.renderGlyphHook(font, currentState, fontSize, pointSize, code, unicode, bytes.CurrentOffset(), renderingMatrix, textMatrix, transformationMatrix, boundingBox)
	} else {
		sp.RenderGlyph(font, currentState, fontSize, pointSize, code, unicode, bytes.CurrentOffset(), renderingMatrix, textMatrix, transformationMatrix, boundingBox)
	}

	var tx, ty float64
	if font.IsVertical() {
		if verticalFont, ok := font.(IVerticalWritingSupported); ok {
			displacement := verticalFont.GetDisplacementVector(code)
			tx = 0
			ty = (displacement.Y*fontSize) + characterSpacing + wordSpacing
		} else {
			tx = 0
			ty = characterSpacing + wordSpacing
		}
	} else {
		tx = (boundingBox.Width*fontSize + characterSpacing + wordSpacing) * horizontalScaling
		ty = 0
	}

	translation := core.GetTranslationMatrix(tx, ty)
	sp.textMatrices.TextMatrix = translation.Multiply(sp.textMatrices.TextMatrix)

	return nil
}

// RenderGlyph is the abstract method that concrete implementations must provide to render a single glyph.
func (sp *BaseStreamProcessor[TPageContent]) RenderGlyph(
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
	_ = font
	_ = currentState
	_ = fontSize
	_ = pointSize
	_ = code
	_ = unicode
	_ = currentOffset
	_ = renderingMatrix
	_ = textMatrix
	_ = transformationMatrix
	_ = characterBoundingBox
}

// ShowPositionedText renders text with explicit positioning adjustments from the token list.
func (sp *BaseStreamProcessor[TPageContent]) ShowPositionedText(tokenList []tokens.Token) {
	sp.TextSequence++

	currentState := sp.GetCurrentState()
	if currentState == nil {
		return
	}

	textState := currentState.FontState
	fontSize := textState.FontSize
	horizontalScaling := textState.HorizontalScaling / 100.0
	font := sp.ResourceStore.GetFont(textState.FontName)

	if font == nil {
		if sp.ParsingOptions.SkipMissingFonts {
			sp.ParsingOptions.Logger.Warn(
				fmt.Sprintf("Skipping a missing font with name %v since it is not present in the document and SkipMissingFonts is set to true.", textState.FontName))
			return
		}
		panic(fmt.Sprintf("Could not find the font with name %v in the resource store. It has not been loaded yet.", textState.FontName))
	}

	isVertical := font.IsVertical()

	for _, token := range tokenList {
		if number, ok := token.(*tokens.NumericToken); ok {
			positionAdjustment := number.Data()

			var tx, ty float64
			if isVertical {
				tx = 0
				ty = -positionAdjustment / 1000 * fontSize
			} else {
				tx = -positionAdjustment / 1000 * fontSize * horizontalScaling
				ty = 0
			}

			sp.adjustTextMatrix(tx, ty)
		} else {
			var bytesData []byte
			switch t := token.(type) {
			case *tokens.HexToken:
				bytesData = make([]byte, len(t.Bytes()))
				copy(bytesData, t.Bytes())
			case *tokens.StringToken:
				bytesData = t.GetBytes()
			}

			sp.ShowText(core.NewMemoryInputBytes(bytesData))
		}
	}
}

// ApplyXObject invokes an XObject by name from the resource store.
func (sp *BaseStreamProcessor[TPageContent]) ApplyXObject(xObjectName *tokens.NameToken) {
	xObjectStream, found := sp.ResourceStore.TryGetXObject(xObjectName)
	if !found {
		if sp.ParsingOptions.SkipMissingFonts {
			return
		}
		panic(fmt.Errorf("No XObject with name %v found on page %d.", xObjectName, sp.PageNumber))
	}

	subTypeToken := xObjectStream.StreamDictionary.Data()[tokens.Subtype.Data()]
	subType, ok := subTypeToken.(*tokens.NameToken)
	if !ok {
		panic(fmt.Errorf("XObject Subtype is not a NameToken: %v", subTypeToken))
	}

	state := sp.GetCurrentState()
	matrix := state.CurrentTransformationMatrix

	// Use current stroking color space or fall back to DeviceRGB (matches C# behavior).
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
		sp.xObjects[xobjects.PostScript] = append(sp.xObjects[xobjects.PostScript], contentRecord)
	} else if subType.Data() == tokens.Image.Data() {
		contentRecord := xobjects.NewXObjectContentRecord(
			xObjectName,
			xObjectStream,
			xobjects.Image,
			matrix,
			state.RenderingIntent,
			colorSpace,
		)
		if sp.renderXObjectImageHook != nil {
			sp.renderXObjectImageHook(contentRecord)
		}
	} else if subType.Data() == tokens.Form.Data() {
		sp.ProcessFormXObject(xObjectStream, xObjectName)
	} else {
		panic(fmt.Errorf("XObject encountered with unexpected SubType %v. %v.", subType, xObjectStream.StreamDictionary))
	}
}

// RenderXObjectImage is the abstract method for rendering an XObject image.
func (sp *BaseStreamProcessor[TPageContent]) RenderXObjectImage(record *xobjects.XObjectContentRecord) {
	_ = record
}

// ProcessFormXObject processes a form XObject following PDF spec steps 1-5.
func (sp *BaseStreamProcessor[TPageContent]) ProcessFormXObject(formStream *tokens.StreamToken, xObjectName *tokens.NameToken) {
	var formResources *tokens.DictionaryToken
	if dictRaw, ok := formStream.StreamDictionary.Data()[tokens.Resources.Data()]; ok {
		formResources = resolveToDictionary(dictRaw, sp.PdfScanner)
	}

	if formResources != nil {
		sp.ResourceStore.LoadResourceDictionary(formResources)
	}

	// 1. Save current state.
	sp.PushState()

	startState := sp.GetCurrentState()

	// Transparency Group XObjects - check for /Group entry
	if groupToken, ok := formStream.StreamDictionary.Data()[tokens.Group.Data()]; ok {
		formGroupDict := resolveToDictionary(groupToken, sp.PdfScanner)

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
	sp.ModifyCurrentTransformationMatrix(formMatrix)

	contentStream := sp.FilterProvider.DecodeStream(formStream, sp.PdfScanner)

	operations := sp.PageContentParser.Parse(
		sp.PageNumber,
		core.NewMemoryInputBytes(contentStream),
		sp.ParsingOptions.Logger,
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
					bbox := core.NewPdfRectangleFloat(points[0], points[1], points[2], points[3])
					if sp.formCtx != nil {
						sp.formCtx.ClipToRectangle(bbox, core.FillingRuleEvenOdd)
					} else {
						sp.ClipToRectangle(bbox, core.FillingRuleEvenOdd)
					}
				}
			}
		}
	}

	// 4. Check for circular references and paint objects.
	hasCircularReference := sp.HasFormXObjectCircularReference(formStream, xObjectName, operations)
	if hasCircularReference {
		if sp.ParsingOptions.UseLenientParsing {
			filteredOps := make([]content.GraphicsStateOperation, 0, len(operations))
			for _, op := range operations {
				if invoke, ok := op.(*InvokeNamedXObject); ok && invoke.Name != nil && invoke.Name.Data() == xObjectName.Data() {
					continue
				}
				filteredOps = append(filteredOps, op)
			}
			operations = filteredOps
			sp.ParsingOptions.Logger.Warn(
				fmt.Sprintf("An XObject form named '%v' is referencing itself which can cause unexpected behaviour. The self reference was removed from the operations before further processing.", xObjectName))
		} else {
			panic(fmt.Errorf("An XObject form named '%v' is referencing itself which can cause unexpected behaviour.", xObjectName))
		}
	}

	sp.ProcessOperations(operations)

	// 5. Restore saved state.
	sp.PopState()

	if formResources != nil {
		sp.ResourceStore.UnloadResourceDictionary()
	}
}

// HasFormXObjectCircularReference checks if a form XObject references itself.
func (sp *BaseStreamProcessor[TPageContent]) HasFormXObjectCircularReference(
	formStream *tokens.StreamToken,
	xObjectName *tokens.NameToken,
	operations []content.GraphicsStateOperation,
) bool {
	if xObjectName == nil {
		return false
	}

	// Check if any operation invokes the same XObject name.
	hasSelfReference := false
	for _, op := range operations {
		if invoke, ok := op.(*InvokeNamedXObject); ok && invoke.Name != nil && invoke.Name.Data() == xObjectName.Data() {
			hasSelfReference = true
			break
		}
	}
	if !hasSelfReference {
		return false
	}

	// Get the XObject entry from formStream's own Resources/XObject dictionary (t1).
	t1 := getXObjectTokenFromStream(formStream, xObjectName)
	if t1 == nil {
		return false
	}

	// Get resourceStream from ResourceStore.
	resourceStream, found := sp.ResourceStore.TryGetXObject(xObjectName)
	if !found {
		return false
	}

	// Get the XObject entry from resourceStream's own Resources/XObject dictionary (t2).
	t2 := getXObjectTokenFromStream(resourceStream, xObjectName)
	if t2 == nil {
		return false
	}

	// Compare the two tokens — if they reference the same PDF object, it's circular.
	return tokensEqual(t1, t2)
}

// getXObjectTokenFromStream looks up xObjectName in streamToken's Resources/XObject dictionary.
func getXObjectTokenFromStream(streamToken *tokens.StreamToken, xObjectName *tokens.NameToken) tokens.Token {
	if streamToken == nil || streamToken.StreamDictionary == nil {
		return nil
	}

	dictData := streamToken.StreamDictionary.Data()
	resourcesVal, ok := dictData[tokens.Resources.Data()]
	if !ok {
		return nil
	}

	resourcesDict, ok := resourcesVal.(*tokens.DictionaryToken)
	if !ok {
		return nil
	}

	xobjectVal, ok := resourcesDict.Data()[tokens.Xobject.Data()]
	if !ok {
		return nil
	}

	xobjectDict, ok := xobjectVal.(*tokens.DictionaryToken)
	if !ok {
		return nil
	}

	token, ok := xobjectDict.Data()[xObjectName.Data()]
	if !ok {
		return nil
	}

	return token
}

// tokensEqual compares two tokens for equality, using the Equals method if available.
func tokensEqual(a, b tokens.Token) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a == b {
		return true
	}
	if eq, ok := a.(interface{ Equals(tokens.Token) bool }); ok {
		return eq.Equals(b)
	}
	return false
}

// SetNamedGraphicsState applies graphics state parameters from a named extended graphics state dictionary.
func (sp *BaseStreamProcessor[TPageContent]) SetNamedGraphicsState(stateName *tokens.NameToken) {
	currentGraphicsState := sp.GetCurrentState()
	if currentGraphicsState == nil {
		return
	}

	state := sp.ResourceStore.GetExtendedGraphicsStateDictionary(stateName)
	if state == nil {
		return
	}

	data := state.Data()

	if lwRaw, ok := data[tokens.Lw.Data()]; ok {
		if lwToken, isNum := lwRaw.(*tokens.NumericToken); isNum {
			currentGraphicsState.LineWidth = lwToken.Data()
		}
	}

	if lcRaw, ok := data[tokens.Lc.Data()]; ok {
		if lcToken, isNum := lcRaw.(*tokens.NumericToken); isNum {
			currentGraphicsState.CapStyle = graphiccore.LineCapStyle(int(lcToken.Data()))
		}
	}

	if ljRaw, ok := data[tokens.Lj.Data()]; ok {
		if ljToken, isNum := ljRaw.(*tokens.NumericToken); isNum {
			currentGraphicsState.JoinStyle = graphiccore.LineJoinStyle(int(ljToken.Data()))
		}
	}

	if fontArrayRaw, ok := data[tokens.Font.Data()]; ok {
		if fontArray, isArray := fontArrayRaw.(*tokens.ArrayToken); isArray {
			d := fontArray.Data()
			if len(d) == 2 {
				if fontRef, isRef := d[0].(*tokens.IndirectReferenceToken); isRef {
					if sizeToken, isNum := d[1].(*tokens.NumericToken); isNum {
						currentGraphicsState.FontState.FromExtendedGraphicsState = true
						currentGraphicsState.FontState.FontSize = sizeToken.Data()
						sp.ActiveExtendedGraphicsStateFont = sp.ResourceStore.GetFontDirectly(fontRef)
					}
				}
			}
		}
	}

	if aisRaw, ok := data[tokens.Ais.Data()]; ok {
		if aisToken, isBool := aisRaw.(*tokens.BooleanToken); isBool {
			currentGraphicsState.AlphaSource = aisToken.Data()
		}
	}

	if caRaw, ok := data[tokens.Ca.Data()]; ok {
		if caToken, isNum := caRaw.(*tokens.NumericToken); isNum {
			currentGraphicsState.AlphaConstantStroking = caToken.Data()
		}
	}

	if cansRaw, ok := data[tokens.CaNs.Data()]; ok {
		if cansToken, isNum := cansRaw.(*tokens.NumericToken); isNum {
			currentGraphicsState.AlphaConstantNonStroking = cansToken.Data()
		}
	}

	if opRaw, ok := data[tokens.Op.Data()]; ok {
		if opToken, isBool := opRaw.(*tokens.BooleanToken); isBool {
			currentGraphicsState.Overprint = opToken.Data()
		}
	}

	if opNsRaw, ok := data[tokens.OpNs.Data()]; ok {
		if opNsToken, isBool := opNsRaw.(*tokens.BooleanToken); isBool {
			currentGraphicsState.NonStrokingOverprint = opNsToken.Data()
		}
	}

	if opmRaw, ok := data[tokens.Opm.Data()]; ok {
		if opmToken, isNum := opmRaw.(*tokens.NumericToken); isNum {
			currentGraphicsState.OverprintMode = int(opmToken.Data())
		}
	}

	if saRaw, ok := data[tokens.Sa.Data()]; ok {
		if saToken, isBool := saRaw.(*tokens.BooleanToken); isBool {
			currentGraphicsState.StrokeAdjustment = saToken.Data()
		}
	}

	if bmRaw, ok := data[tokens.Bm.Data()]; ok {
		if bmToken, isName := bmRaw.(*tokens.NameToken); isName {
			if mode, found := graphiccore.ParseBlendMode(bmToken.Data()); found {
				currentGraphicsState.BlendMode = mode
			} else {
				currentGraphicsState.BlendMode = graphiccore.Normal
			}
		} else if bmArrayToken, isArray := bmRaw.(*tokens.ArrayToken); isArray {
			currentGraphicsState.BlendMode = graphiccore.Normal
			for _, item := range bmArrayToken.Data() {
				if nameTok, isName := item.(*tokens.NameToken); isName {
					if mode, found := graphiccore.ParseBlendMode(nameTok.Data()); found {
						currentGraphicsState.BlendMode = mode
						break
					}
				}
			}
		}
	}

	if smaskRaw, ok := data[tokens.Smask.Data()]; ok {
		if smToken, isName := smaskRaw.(*tokens.NameToken); isName {
			if smToken.Data() == tokens.None.Data() {
				currentGraphicsState.SoftMask = nil
			}
		} else if smDictToken, isDict := smaskRaw.(*tokens.DictionaryToken); isDict {
			lp, ok := sp.FilterProvider.(filters.LookupFilterProvider)
			if !ok {
				sp.ParsingOptions.Logger.Error("filter provider does not implement filters.LookupFilterProvider, skipping soft mask")
			} else {
				sm, err := ParseSoftMask(smDictToken, sp.PdfScanner, lp)
				if err != nil {
					sp.ParsingOptions.Logger.Error(fmt.Sprintf("Failed to parse soft mask: %v", err))
				} else {
					currentGraphicsState.SoftMask = sm
				}
			}
		}
	}
}

// BeginInlineImage starts the inline image collection process.
func (sp *BaseStreamProcessor[TPageContent]) BeginInlineImage() {
	if sp.InlineImageBuilder != nil {
		sp.ParsingOptions.Logger.Error("Begin inline image (BI) command encountered while another inline image was active.")
	}
	sp.InlineImageBuilder = NewInlineImageBuilder()
}

// SetInlineImageProperties sets the properties for the current inline image.
func (sp *BaseStreamProcessor[TPageContent]) SetInlineImageProperties(properties map[*tokens.NameToken]tokens.Token) {
	if sp.InlineImageBuilder == nil {
		sp.ParsingOptions.Logger.Error("Begin inline image data (ID) command encountered without a corresponding begin inline image (BI) command.")
		return
	}
	sp.InlineImageBuilder.Properties = properties
}

// EndInlineImage finalizes and renders the inline image.
func (sp *BaseStreamProcessor[TPageContent]) EndInlineImage(bytes []byte) {
	if sp.InlineImageBuilder == nil {
		sp.ParsingOptions.Logger.Error("End inline image (EI) command encountered without a corresponding begin inline image (BI) command.")
		return
	}

	sp.InlineImageBuilder.Bytes = bytes

	image := sp.InlineImageBuilder.CreateInlineImage(
		sp.GetCurrentState().CurrentTransformationMatrix,
		sp.FilterProvider,
		sp.PdfScanner,
		sp.GetCurrentState().RenderingIntent,
		sp.ResourceStore,
	)

	if sp.renderInlineImageHook != nil {
		sp.renderInlineImageHook(image)
	}
	sp.InlineImageBuilder = nil
}

// RenderInlineImage is the abstract method for rendering an inline image.
func (sp *BaseStreamProcessor[TPageContent]) RenderInlineImage(image *content.InlineImage) {
	_ = image
}

// BeginMarkedContent starts a marked content section.
func (sp *BaseStreamProcessor[TPageContent]) BeginMarkedContent(name *tokens.NameToken, propertyDictionaryName *tokens.NameToken, properties *tokens.DictionaryToken) {
	_ = name
	_ = propertyDictionaryName
	_ = properties
}

// EndMarkedContent ends the current marked content section.
func (sp *BaseStreamProcessor[TPageContent]) EndMarkedContent() {
}

// SetFlatnessTolerance sets the flatness tolerance for curve rendering.
func (sp *BaseStreamProcessor[TPageContent]) SetFlatnessTolerance(tolerance float64) {
	state := sp.GetCurrentState()
	if state != nil {
		state.Flatness = tolerance
	}
}

// SetLineCap sets the line cap style.
func (sp *BaseStreamProcessor[TPageContent]) SetLineCap(cap graphiccore.LineCapStyle) {
	state := sp.GetCurrentState()
	if state != nil {
		state.CapStyle = cap
	}
}

// SetLineDashPattern sets the line dash pattern.
func (sp *BaseStreamProcessor[TPageContent]) SetLineDashPattern(pattern graphiccore.LineDashPattern) {
	state := sp.GetCurrentState()
	if state != nil {
		state.LineDashPattern = pattern
	}
}

// SetLineJoin sets the line join style.
func (sp *BaseStreamProcessor[TPageContent]) SetLineJoin(join graphiccore.LineJoinStyle) {
	state := sp.GetCurrentState()
	if state != nil {
		state.JoinStyle = join
	}
}

// SetLineWidth sets the stroke width.
func (sp *BaseStreamProcessor[TPageContent]) SetLineWidth(width float64) {
	state := sp.GetCurrentState()
	if state != nil {
		state.LineWidth = width
	}
}

// SetMiterLimit sets the miter limit for miter joins.
func (sp *BaseStreamProcessor[TPageContent]) SetMiterLimit(limit float64) {
	state := sp.GetCurrentState()
	if state != nil {
		state.MiterLimit = limit
	}
}

// MoveToNextLineWithOffset moves to the next line using current leading.
func (sp *BaseStreamProcessor[TPageContent]) MoveToNextLineWithOffset() {
	state := sp.GetCurrentState()
	if state == nil {
		return
	}
	td := NewMoveToNextLineWithOffset(0, -1*state.FontState.Leading)
	td.Run(sp)
}

// SetFontAndSize selects a font and sets its size.
func (sp *BaseStreamProcessor[TPageContent]) SetFontAndSize(font *tokens.NameToken, size float64) {
	state := sp.GetCurrentState()
	if state != nil {
		state.FontState.FontSize = size
		state.FontState.FontName = font
	}
}

// SetHorizontalScaling sets the horizontal text scaling percentage.
func (sp *BaseStreamProcessor[TPageContent]) SetHorizontalScaling(scale float64) {
	state := sp.GetCurrentState()
	if state != nil {
		state.FontState.HorizontalScaling = scale
	}
}

// SetTextLeading sets the text leading (line spacing).
func (sp *BaseStreamProcessor[TPageContent]) SetTextLeading(leading float64) {
	state := sp.GetCurrentState()
	if state != nil {
		state.FontState.Leading = leading
	}
}

// SetTextRenderingMode sets how text is rendered.
func (sp *BaseStreamProcessor[TPageContent]) SetTextRenderingMode(mode core.TextRenderingMode) {
	state := sp.GetCurrentState()
	if state != nil {
		state.FontState.TextRenderingMode = mode
	}
}

// SetTextRise sets the vertical displacement for text drawing.
func (sp *BaseStreamProcessor[TPageContent]) SetTextRise(rise float64) {
	state := sp.GetCurrentState()
	if state != nil {
		state.FontState.Rise = rise
	}
}

// SetWordSpacing sets additional space between words.
func (sp *BaseStreamProcessor[TPageContent]) SetWordSpacing(spacing float64) {
	state := sp.GetCurrentState()
	if state != nil {
		state.FontState.WordSpacing = spacing
	}
}

// ModifyCurrentTransformationMatrix concatenates a matrix to the current transformation matrix.
func (sp *BaseStreamProcessor[TPageContent]) ModifyCurrentTransformationMatrix(value core.TransformationMatrix) {
	state := sp.GetCurrentState()
	if state != nil {
		state.CurrentTransformationMatrix = value.Multiply(state.CurrentTransformationMatrix)
	}
}

// SetCharacterSpacing sets additional space between characters.
func (sp *BaseStreamProcessor[TPageContent]) SetCharacterSpacing(spacing float64) {
	state := sp.GetCurrentState()
	if state != nil {
		state.FontState.CharacterSpacing = spacing
	}
}

// PaintShading paints a shading fill by name.
func (sp *BaseStreamProcessor[TPageContent]) PaintShading(shadingName *tokens.NameToken) {
	_ = shadingName
}

// TextMatrices returns the text matrix state.
func (sp *BaseStreamProcessor[TPageContent]) TextMatrices() *TextMatrices {
	return sp.textMatrices
}

// CurrentTransformationMatrix returns the current transformation matrix.
func (sp *BaseStreamProcessor[TPageContent]) CurrentTransformationMatrix() core.TransformationMatrix {
	state := sp.GetCurrentState()
	if state != nil {
		return state.CurrentTransformationMatrix
	}
	return core.Identity
}

// StackSize returns the number of graphics states on the stack.
func (sp *BaseStreamProcessor[TPageContent]) StackSize() int {
	return len(sp.GraphicsStack)
}

// CurrentPosition returns the current drawing position.
func (sp *BaseStreamProcessor[TPageContent]) CurrentPosition() core.PdfPoint {
	return sp.currentPosition
}

// SetCurrentPosition sets the current drawing position.
func (sp *BaseStreamProcessor[TPageContent]) SetCurrentPosition(pos core.PdfPoint) {
	sp.currentPosition = pos
}

// ClipToRectangle clips to the given rectangle with the specified filling rule.
func (sp *BaseStreamProcessor[TPageContent]) ClipToRectangle(rectangle core.PdfRectangle, clippingRule core.FillingRule) {
	_ = rectangle
	_ = clippingRule
}

// adjustTextMatrix applies a translation to the text matrix.
func (sp *BaseStreamProcessor[TPageContent]) adjustTextMatrix(tx, ty float64) {
	matrix := core.GetTranslationMatrix(tx, ty)
	sp.textMatrices.TextMatrix = matrix.Multiply(sp.textMatrices.TextMatrix)
}

// --- Path construction methods (abstract stubs for concrete implementations) ---

// BeginSubpath starts a new subpath at the current point.
func (sp *BaseStreamProcessor[TPageContent]) BeginSubpath() {
}

// CloseSubpath closes the current subpath and returns the closing point.
func (sp *BaseStreamProcessor[TPageContent]) CloseSubpath() *core.PdfPoint {
	return nil
}

// StrokePath strokes the current path.
func (sp *BaseStreamProcessor[TPageContent]) StrokePath(close bool) {
	_ = close
}

// FillPath fills the current path using the given filling rule.
func (sp *BaseStreamProcessor[TPageContent]) FillPath(fillingRule core.FillingRule, close bool) {
	_ = fillingRule
	_ = close
}

// FillStrokePath fills then strokes the current path.
func (sp *BaseStreamProcessor[TPageContent]) FillStrokePath(fillingRule core.FillingRule, close bool) {
	_ = fillingRule
	_ = close
}

// MoveTo moves to the specified coordinates.
func (sp *BaseStreamProcessor[TPageContent]) MoveTo(x, y float64) {
	_ = x
	_ = y
}

// BezierCurveToQuadratic adds a quadratic Bézier curve segment.
func (sp *BaseStreamProcessor[TPageContent]) BezierCurveToQuadratic(x2, y2, x3, y3 float64) {
	_ = x2
	_ = y2
	_ = x3
	_ = y3
}

// BezierCurveToCubic adds a cubic Bézier curve segment.
func (sp *BaseStreamProcessor[TPageContent]) BezierCurveToCubic(x1, y1, x2, y2, x3, y3 float64) {
	_ = x1
	_ = y1
	_ = x2
	_ = y2
	_ = x3
	_ = y3
}

// BezierCurveToStartCP adds a cubic Bézier curve segment using the current point as first control point ("v" operator).
func (sp *BaseStreamProcessor[TPageContent]) BezierCurveToStartCP(x2, y2, x3, y3 float64) {
	_ = x2
	_ = y2
	_ = x3
	_ = y3
}

// LineTo adds a line segment to the current point.
func (sp *BaseStreamProcessor[TPageContent]) LineTo(x, y float64) {
	_ = x
	_ = y
}

// Rectangle adds a rectangle subpath.
func (sp *BaseStreamProcessor[TPageContent]) Rectangle(x, y, width, height float64) {
	_ = x
	_ = y
	_ = width
	_ = height
}

// EndPath ends the current path construction without filling or stroking.
func (sp *BaseStreamProcessor[TPageContent]) EndPath() {
}

// ClosePath closes and ends the current path.
func (sp *BaseStreamProcessor[TPageContent]) ClosePath() {
}

// ModifyClippingIntersect modifies the clipping path by intersecting with the current path.
func (sp *BaseStreamProcessor[TPageContent]) ModifyClippingIntersect(clippingRule core.FillingRule) {
	_ = clippingRule
}
