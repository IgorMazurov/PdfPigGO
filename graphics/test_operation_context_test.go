// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/geometry"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// TestOperationContext is a test implementation of OperationContext used for unit testing.
type TestOperationContext struct {
	StateStack          []*CurrentGraphicsState
	textMatrices      *TextMatrices
	CurrentSubpath      *core.PdfSubpath
	CurrentPath         *geometry.PdfPath
	currentPosition     core.PdfPoint
}

// NewTestOperationContext creates a new TestOperationContext with default state.
func NewTestOperationContext() *TestOperationContext {
	ctx := &TestOperationContext{
		StateStack:   make([]*CurrentGraphicsState, 0, 16),
		textMatrices:      NewTextMatrices(),
	}

	getCurrentState := func() *CurrentGraphicsState {
		return ctx.GetCurrentState()
	}

	initialState := &CurrentGraphicsState{
		CurrentTransformationMatrix: core.Identity,
		LineWidth:                   1,
		ColorSpaceContext:           NewColorSpaceContext(getCurrentState, nil),
	}
	ctx.StateStack = append(ctx.StateStack, initialState)
	ctx.CurrentSubpath = core.NewPdfSubpath()

	return ctx
}

// GetCurrentState returns the topmost graphics state on the stack.
func (ctx *TestOperationContext) GetCurrentState() *CurrentGraphicsState {
	if len(ctx.StateStack) == 0 {
		return nil
	}
	return ctx.StateStack[len(ctx.StateStack)-1]
}

// CurrentPosition returns the current drawing position.
func (ctx *TestOperationContext) CurrentPosition() core.PdfPoint {
	return ctx.currentPosition
}

// SetCurrentPosition sets the current drawing position.
func (ctx *TestOperationContext) SetCurrentPosition(pos core.PdfPoint) {
	ctx.currentPosition = pos
}

// StackSize returns the number of graphics states on the stack.
func (ctx *TestOperationContext) StackSize() int {
	return len(ctx.StateStack)
}

// PopState removes the topmost graphics state from the stack.
func (ctx *TestOperationContext) PopState() {
	if len(ctx.StateStack) == 0 {
		panic("cannot pop from empty graphics state stack")
	}
	ctx.StateStack = ctx.StateStack[:len(ctx.StateStack)-1]
}

// PushState pushes a deep clone of the current graphics state onto the stack.
func (ctx *TestOperationContext) PushState() {
	current := ctx.GetCurrentState()
	if current != nil {
		ctx.StateStack = append(ctx.StateStack, current.DeepClone())
	}
}

// ShowText is not implemented for test context.
func (ctx *TestOperationContext) ShowText(_ core.InputBytes) {}

// ShowPositionedText is not implemented for test context.
func (ctx *TestOperationContext) ShowPositionedText(_ []tokens.Token) {}

// ApplyXObject is not implemented for test context.
func (ctx *TestOperationContext) ApplyXObject(_ *tokens.NameToken) {}

// SetNamedGraphicsState is not implemented for test context.
func (ctx *TestOperationContext) SetNamedGraphicsState(_ *tokens.NameToken) {}

// BeginInlineImage is not implemented for test context.
func (ctx *TestOperationContext) BeginInlineImage() {}

// SetInlineImageProperties is not implemented for test context.
func (ctx *TestOperationContext) SetInlineImageProperties(_ map[*tokens.NameToken]tokens.Token) {}

// EndInlineImage is not implemented for test context.
func (ctx *TestOperationContext) EndInlineImage(_ []byte) {}

// BeginMarkedContent is not implemented for test context.
func (ctx *TestOperationContext) BeginMarkedContent(_ *tokens.NameToken, _ *tokens.NameToken, _ *tokens.DictionaryToken) {}

// EndMarkedContent is not implemented for test context.
func (ctx *TestOperationContext) EndMarkedContent() {}

// SetFlatnessTolerance sets the flatness tolerance on the current graphics state.
func (ctx *TestOperationContext) SetFlatnessTolerance(tolerance float64) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.Flatness = tolerance
	}
}

// SetLineCap sets the line cap style on the current graphics state.
func (ctx *TestOperationContext) SetLineCap(cap graphiccore.LineCapStyle) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.CapStyle = cap
	}
}

// SetLineDashPattern sets the line dash pattern on the current graphics state.
func (ctx *TestOperationContext) SetLineDashPattern(pattern graphiccore.LineDashPattern) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.LineDashPattern = pattern
	}
}

// SetLineJoin sets the line join style on the current graphics state.
func (ctx *TestOperationContext) SetLineJoin(join graphiccore.LineJoinStyle) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.JoinStyle = join
	}
}

// SetLineWidth sets the line width on the current graphics state.
func (ctx *TestOperationContext) SetLineWidth(width float64) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.LineWidth = width
	}
}

// SetMiterLimit sets the miter limit on the current graphics state.
func (ctx *TestOperationContext) SetMiterLimit(limit float64) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.MiterLimit = limit
	}
}

// MoveToNextLineWithOffset moves to the next line using the current leading value.
func (ctx *TestOperationContext) MoveToNextLineWithOffset() {
	state := ctx.GetCurrentState()
	if state == nil {
		return
	}
	leading := -1 * state.FontState.Leading
	matrix := core.GetTranslationMatrix(0, leading)
	ctx.textMatrices.TextMatrix = matrix.Multiply(ctx.textMatrices.TextMatrix)
}

// SetFontAndSize sets the font name and size on the current graphics state.
func (ctx *TestOperationContext) SetFontAndSize(font *tokens.NameToken, size float64) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.FontState.FontSize = size
		state.FontState.FontName = font
	}
}

// SetHorizontalScaling sets the horizontal scaling on the current graphics state.
func (ctx *TestOperationContext) SetHorizontalScaling(scale float64) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.FontState.HorizontalScaling = scale
	}
}

// SetTextLeading sets the text leading on the current graphics state.
func (ctx *TestOperationContext) SetTextLeading(leading float64) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.FontState.Leading = leading
	}
}

// SetTextRenderingMode sets the text rendering mode on the current graphics state.
func (ctx *TestOperationContext) SetTextRenderingMode(mode core.TextRenderingMode) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.FontState.TextRenderingMode = mode
	}
}

// SetTextRise sets the text rise on the current graphics state.
func (ctx *TestOperationContext) SetTextRise(rise float64) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.FontState.Rise = rise
	}
}

// SetWordSpacing sets the word spacing on the current graphics state.
func (ctx *TestOperationContext) SetWordSpacing(spacing float64) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.FontState.WordSpacing = spacing
	}
}

// ModifyCurrentTransformationMatrix multiplies the given matrix with the current transformation matrix.
func (ctx *TestOperationContext) ModifyCurrentTransformationMatrix(value core.TransformationMatrix) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.CurrentTransformationMatrix = value.Multiply(state.CurrentTransformationMatrix)
	}
}

// SetCharacterSpacing sets the character spacing on the current graphics state.
func (ctx *TestOperationContext) SetCharacterSpacing(spacing float64) {
	state := ctx.GetCurrentState()
	if state != nil {
		state.FontState.CharacterSpacing = spacing
	}
}

// PaintShading is not implemented for test context.
func (ctx *TestOperationContext) PaintShading(_ *tokens.NameToken) {}

// BeginSubpath initializes a new empty subpath.
func (ctx *TestOperationContext) BeginSubpath() {
	ctx.CurrentSubpath = core.NewPdfSubpath()
}

// CloseSubpath closes the current subpath and returns the closing point.
func (ctx *TestOperationContext) CloseSubpath() *core.PdfPoint {
	return &core.PdfPoint{}
}

// StrokePath is not implemented for test context.
func (ctx *TestOperationContext) StrokePath(_ bool) {}

// FillPath is not implemented for test context.
func (ctx *TestOperationContext) FillPath(_ core.FillingRule, _ bool) {}

// FillStrokePath is not implemented for test context.
func (ctx *TestOperationContext) FillStrokePath(_ core.FillingRule, _ bool) {}

// MoveTo moves to the given point transformed by the current transformation matrix.
func (ctx *TestOperationContext) MoveTo(x, y float64) {
	ctx.BeginSubpath()
	point := ctx.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x, y))
	ctx.currentPosition = point
	ctx.CurrentSubpath.MoveTo(point.X, point.Y)
}

// BezierCurveToQuadratic draws a quadratic Bézier curve from the current position to (x3, y3) with control point (x2, y2).
func (ctx *TestOperationContext) BezierCurveToQuadratic(x2, y2, x3, y3 float64) {
	if ctx.CurrentSubpath == nil {
		return
	}
	controlPoint2 := ctx.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x2, y2))
	end := ctx.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x3, y3))
	_ = ctx.CurrentSubpath.BezierCurveToQuadratic(ctx.currentPosition.X, ctx.currentPosition.Y, controlPoint2.X, controlPoint2.Y)
	ctx.currentPosition = end
}

// BezierCurveToCubic draws a cubic Bézier curve from the current position to (x3, y3).
func (ctx *TestOperationContext) BezierCurveToCubic(x1, y1, x2, y2, x3, y3 float64) {
	if ctx.CurrentSubpath == nil {
		return
	}
	controlPoint1 := ctx.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x1, y1))
	controlPoint2 := ctx.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x2, y2))
	end := ctx.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x3, y3))
	_ = ctx.CurrentSubpath.BezierCurveToCubic(controlPoint1.X, controlPoint1.Y, controlPoint2.X, controlPoint2.Y, end.X, end.Y)
	ctx.currentPosition = end
}

// BezierCurveToStartCP draws a Bézier curve using the current position as first control point.
func (ctx *TestOperationContext) BezierCurveToStartCP(x2, y2, x3, y3 float64) {
	if ctx.CurrentSubpath == nil {
		return
	}
	controlPoint2 := ctx.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x2, y2))
	end := ctx.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x3, y3))
	_ = ctx.CurrentSubpath.BezierCurveToCubic(ctx.currentPosition.X, ctx.currentPosition.Y, ctx.currentPosition.X, ctx.currentPosition.Y, controlPoint2.X, controlPoint2.Y)
	ctx.currentPosition = end
}

// LineTo draws a line from the current position to (x, y) transformed by the current matrix.
func (ctx *TestOperationContext) LineTo(x, y float64) {
	if ctx.CurrentSubpath == nil {
		return
	}
	endPoint := ctx.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x, y))
	_ = ctx.CurrentSubpath.LineTo(endPoint.X, endPoint.Y)
	ctx.currentPosition = endPoint
}

// Rectangle draws a rectangle at (x, y) with the given width and height.
func (ctx *TestOperationContext) Rectangle(x, y, width, height float64) {
	ctx.BeginSubpath()
	lowerLeft := ctx.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x, y))
	upperRight := ctx.CurrentTransformationMatrix().TransformPoint(core.NewPdfPoint(x+width, y+height))
	ctx.CurrentSubpath.Rectangle(lowerLeft.X, lowerLeft.Y, upperRight.X-lowerLeft.X, upperRight.Y-lowerLeft.Y)
}

// EndPath is not implemented for test context.
func (ctx *TestOperationContext) EndPath() {}

// ClosePath is not implemented for test context.
func (ctx *TestOperationContext) ClosePath() {}

// ModifyClippingIntersect is not implemented for test context.
func (ctx *TestOperationContext) ModifyClippingIntersect(_ core.FillingRule) {}

// TextMatrices returns the text matrices managed by this context.
func (ctx *TestOperationContext) TextMatrices() *TextMatrices {
	return ctx.textMatrices
}

// CurrentTransformationMatrix returns the current transformation matrix from the top graphics state.
func (ctx *TestOperationContext) CurrentTransformationMatrix() core.TransformationMatrix {
	state := ctx.GetCurrentState()
	if state != nil {
		return state.CurrentTransformationMatrix
	}
	return core.Identity
}

var _ OperationContext = (*TestOperationContext)(nil)

func TestNewTestOperationContextInitialState(t *testing.T) {
	ctx := NewTestOperationContext()

	if ctx.StackSize() != 1 {
		t.Errorf("expected stack size 1, got %d", ctx.StackSize())
	}

	if ctx.GetCurrentState() == nil {
		t.Error("expected initial graphics state to be non-nil")
	}

	if ctx.TextMatrices() == nil {
		t.Error("expected TextMatrices to be initialized")
	}

	if ctx.CurrentSubpath == nil {
		t.Error("expected CurrentSubpath to be initialized")
	}
}

func TestTestOperationContextPushPopState(t *testing.T) {
	ctx := NewTestOperationContext()

	initialSize := ctx.StackSize()
	ctx.PushState()

	if got := ctx.StackSize(); got != initialSize+1 {
		t.Errorf("expected stack size %d after push, got %d", initialSize+1, got)
	}

	ctx.PopState()

	if got := ctx.StackSize(); got != initialSize {
		t.Errorf("expected stack size %d after pop, got %d", initialSize, got)
	}
}

func TestTestOperationContextSetLineWidth(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.SetLineWidth(2.5)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if got := state.LineWidth; got != 2.5 {
		t.Errorf("expected line width 2.5, got %f", got)
	}
}

func TestTestOperationContextSetFlatnessTolerance(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.SetFlatnessTolerance(0.1)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if got := state.Flatness; got != 0.1 {
		t.Errorf("expected flatness 0.1, got %f", got)
	}
}

func TestTestOperationContextSetFontAndSize(t *testing.T) {
	ctx := NewTestOperationContext()
	fontName := tokens.Create("Helvetica")

	ctx.SetFontAndSize(fontName, 12.0)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if state.FontState.FontName != fontName {
		t.Errorf("expected font name %v, got %v", fontName, state.FontState.FontName)
	}

	if got := state.FontState.FontSize; got != 12.0 {
		t.Errorf("expected font size 12.0, got %f", got)
	}
}

func TestTestOperationContextSetLineCap(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.SetLineCap(graphiccore.LineCapRound)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if got := state.CapStyle; got != graphiccore.LineCapRound {
		t.Errorf("expected line cap Round, got %v", got)
	}
}

func TestTestOperationContextSetLineJoin(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.SetLineJoin(graphiccore.LineJoinRound)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if got := state.JoinStyle; got != graphiccore.LineJoinRound {
		t.Errorf("expected line join Round, got %v", got)
	}
}

func TestTestOperationContextSetMiterLimit(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.SetMiterLimit(10.0)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if got := state.MiterLimit; got != 10.0 {
		t.Errorf("expected miter limit 10.0, got %f", got)
	}
}

func TestTestOperationContextMoveTo(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.MoveTo(100, 200)

	if got := ctx.currentPosition.X; got != 100 {
		t.Errorf("expected X=100, got %f", got)
	}
	if got := ctx.currentPosition.Y; got != 200 {
		t.Errorf("expected Y=200, got %f", got)
	}

	if len(ctx.CurrentSubpath.Commands()) != 1 {
		t.Errorf("expected 1 command in subpath after MoveTo, got %d", len(ctx.CurrentSubpath.Commands()))
	}
}

func TestTestOperationContextLineTo(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.MoveTo(0, 0)
	ctx.LineTo(50, 60)

	if got := ctx.currentPosition.X; got != 50 {
		t.Errorf("expected X=50 after LineTo, got %f", got)
	}
	if got := ctx.currentPosition.Y; got != 60 {
		t.Errorf("expected Y=60 after LineTo, got %f", got)
	}

	if len(ctx.CurrentSubpath.Commands()) < 2 {
		t.Errorf("expected at least 2 commands in subpath after MoveTo+LineTo, got %d", len(ctx.CurrentSubpath.Commands()))
	}
}

func TestTestOperationContextLineToNilSubpath(t *testing.T) {
	ctx := NewTestOperationContext()
	ctx.CurrentSubpath = nil

	ctx.LineTo(10, 20)
}

func TestTestOperationContextBezierCurveToQuadratic(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.MoveTo(0, 0)
	ctx.BezierCurveToQuadratic(10, 20, 30, 40)

	if got := ctx.currentPosition.X; got != 30 {
		t.Errorf("expected X=30 after quadratic BezierCurveTo, got %f", got)
	}
	if got := ctx.currentPosition.Y; got != 40 {
		t.Errorf("expected Y=40 after quadratic BezierCurveTo, got %f", got)
	}
}

func TestTestOperationContextBezierCurveToCubic(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.MoveTo(0, 0)
	ctx.BezierCurveToCubic(10, 20, 30, 40, 50, 60)

	if got := ctx.currentPosition.X; got != 50 {
		t.Errorf("expected X=50 after cubic BezierCurveTo, got %f", got)
	}
	if got := ctx.currentPosition.Y; got != 60 {
		t.Errorf("expected Y=60 after cubic BezierCurveTo, got %f", got)
	}
}

func TestTestOperationContextBezierNilSubpath(t *testing.T) {
	ctx := NewTestOperationContext()
	ctx.CurrentSubpath = nil

	ctx.BezierCurveToQuadratic(1, 2, 3, 4)
	ctx.BezierCurveToCubic(1, 2, 3, 4, 5, 6)
}

func TestTestOperationContextRectangle(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.Rectangle(10, 20, 100, 50)

	if len(ctx.CurrentSubpath.Commands()) == 0 {
		t.Error("expected commands in subpath after Rectangle")
	}
}

func TestTestOperationContextSetTextLeading(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.SetTextLeading(14.0)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if got := state.FontState.Leading; got != 14.0 {
		t.Errorf("expected leading 14.0, got %f", got)
	}
}

func TestTestOperationContextSetCharacterSpacing(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.SetCharacterSpacing(2.0)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if got := state.FontState.CharacterSpacing; got != 2.0 {
		t.Errorf("expected character spacing 2.0, got %f", got)
	}
}

func TestTestOperationContextSetWordSpacing(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.SetWordSpacing(5.0)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if got := state.FontState.WordSpacing; got != 5.0 {
		t.Errorf("expected word spacing 5.0, got %f", got)
	}
}

func TestTestOperationContextSetTextRise(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.SetTextRise(3.0)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if got := state.FontState.Rise; got != 3.0 {
		t.Errorf("expected rise 3.0, got %f", got)
	}
}

func TestTestOperationContextSetHorizontalScaling(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.SetHorizontalScaling(150.0)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if got := state.FontState.HorizontalScaling; got != 150.0 {
		t.Errorf("expected horizontal scaling 150.0, got %f", got)
	}
}

func TestTestOperationContextSetTextRenderingMode(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.SetTextRenderingMode(core.FillThenStrokeText)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	if got := state.FontState.TextRenderingMode; got != core.FillThenStrokeText {
		t.Errorf("expected TextRenderingMode FillThenStroke, got %v", got)
	}
}

func TestTestOperationContextModifyCurrentTransformationMatrix(t *testing.T) {
	ctx := NewTestOperationContext()

	translation := core.GetTranslationMatrix(10, 20)
	ctx.ModifyCurrentTransformationMatrix(translation)

	state := ctx.GetCurrentState()
	if state == nil {
		t.Fatal("expected non-nil graphics state")
	}

	transformed := state.CurrentTransformationMatrix.TransformPoint(core.NewPdfPoint(0, 0))
	if got := transformed.X; got != 10 {
		t.Errorf("expected X=10 after translation, got %f", got)
	}
	if got := transformed.Y; got != 20 {
		t.Errorf("expected Y=20 after translation, got %f", got)
	}
}

func TestTestOperationContextMoveToNextLineWithOffset(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.textMatrices.TextMatrix = core.Identity
	ctx.GetCurrentState().FontState.Leading = 14.0

	ctx.MoveToNextLineWithOffset()

	trm := ctx.textMatrices.TextMatrix
	expectedY := -14.0
	if got := trm.F; got != expectedY {
		t.Errorf("expected text matrix F=%f after MoveToNextLineWithOffset, got %f", expectedY, got)
	}
}

func TestTestOperationContextCloseSubpath(t *testing.T) {
	ctx := NewTestOperationContext()

	point := ctx.CloseSubpath()
	if point == nil {
		t.Error("expected CloseSubpath to return a non-nil PdfPoint")
	}
}
