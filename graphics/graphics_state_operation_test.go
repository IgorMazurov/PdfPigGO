package graphics

import (
	"bytes"
	"io"
	"strings"
	"testing"

	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// opWriter is the interface shared by all graphics operations for writing to an io.Writer.
type opWriter interface {
	Write(io.Writer) error
}

// operationTestCase describes a single graphics operation for testing.
type operationTestCase struct {
	name    string
	symbol  string
	instance opWriter
}

// getAllOperations returns test cases for every known GraphicsStateOperation type.
func getAllOperations() []operationTestCase {
	hogName := tokens.MustCreate("Hog")

	return []operationTestCase{
		{name: "Pop", symbol: popSymbol, instance: InstancePop},
		{name: "Push", symbol: pushSymbol, instance: InstancePush},
		{name: "BeginText", symbol: beginTextSymbol, instance: InstanceBeginText},
		{name: "EndText", symbol: endTextSymbol, instance: InstanceEndText},
		{name: "MoveToNextLine", symbol: moveToNextLineSymbol, instance: InstanceMoveToNextLine},
		{name: "CloseSubpath", symbol: closeSubpathSymbol, instance: InstanceCloseSubpath},
		{name: "EndPath", symbol: endPathSymbol, instance: InstanceEndPath},
		{name: "StrokePath", symbol: strokePathSymbol, instance: InstanceStrokePath},
		{name: "FillPathNonZeroWinding", symbol: fillPathNonZeroWindingSymbol, instance: InstanceFillPathNonZeroWinding},
		{name: "FillPathEvenOddRule", symbol: fillPathEvenOddRuleSymbol, instance: InstanceFillPathEvenOddRule},
		{name: "CloseAndStrokePath", symbol: closeAndStrokePathSymbol, instance: InstanceCloseAndStrokePath},
		{name: "FillPathNonZeroWindingAndStroke", symbol: fillPathNonZeroWindingAndStrokeSymbol, instance: InstanceFillPathNonZeroWindingAndStroke},
		{name: "FillPathEvenOddRuleAndStroke", symbol: fillPathEvenOddRuleAndStrokeSymbol, instance: InstanceFillPathEvenOddRuleAndStroke},
		{name: "CloseFillEvenOddStrokePath", symbol: closeFillEvenOddStrokePathSymbol, instance: InstanceCloseFillEvenOddStrokePath},
		{name: "CloseFillNonZeroStrokePath", symbol: closeFillNonZeroStrokePathSymbol, instance: InstanceCloseFillNonZeroStrokePath},
		{name: "BeginCompatibilitySection", symbol: beginCompatibilitySectionSymbol, instance: InstanceBeginCompatibilitySection},
		{name: "EndCompatibilitySection", symbol: endCompatibilitySectionSymbol, instance: InstanceEndCompatibilitySection},
		{name: "ModifyClippingByEvenOddIntersect", symbol: modifyClippingByEvenOddIntersectSymbol, instance: InstanceModifyClippingByEvenOddIntersect},
		{name: "ModifyClippingByNonZeroWindingIntersect", symbol: modifyClippingByNonZeroWindingIntersectSymbol, instance: InstanceModifyClippingByNonZeroWindingIntersect},
		{name: "EndMarkedContent", symbol: endMarkedContentSymbol, instance: InstanceEndMarkedContent},
		{name: "BeginInlineImage", symbol: beginInlineImageSymbol, instance: InstanceBeginInlineImage},
		{name: "FillPathNonZeroWindingCompatibility", symbol: fillPathNonZeroWindingCompatSymbol, instance: InstanceFillPathNonZeroWindingCompatibility},
		{name: "SetCharacterSpacing", symbol: setCharacterSpacingSymbol, instance: NewSetCharacterSpacing(0.5)},
		{name: "SetFontAndSize", symbol: setFontAndSizeSymbol, instance: NewSetFontAndSize(hogName, 12)},
		{name: "SetHorizontalScaling", symbol: setHorizontalScalingSymbol, instance: NewSetHorizontalScaling(100)},
		{name: "SetTextLeading", symbol: setTextLeadingSymbol, instance: NewSetTextLeading(14)},
		{name: "SetTextRenderingMode", symbol: setTextRenderingModeSymbol, instance: NewSetTextRenderingMode(0)},
		{name: "SetTextRise", symbol: setTextRiseSymbol, instance: NewSetTextRise(0)},
		{name: "SetWordSpacing", symbol: setWordSpacingSymbol, instance: NewSetWordSpacing(0)},
		{name: "SetTextMatrix", symbol: setTextMatrixSymbol, instance: NewSetTextMatrix(1, 0, 0, 1, 0, 0)},
		{name: "MoveToNextLineWithOffset", symbol: moveToNextLineWithOffsetSymbol, instance: NewMoveToNextLineWithOffset(0, 14)},
		{name: "MoveToNextLineWithOffsetSetLeading", symbol: moveToNextLineWithOffsetSetLeadingSymbol, instance: NewMoveToNextLineWithOffsetSetLeading(0, 14)},
		{name: "BeginNewSubpath", symbol: beginNewSubpathSymbol, instance: NewBeginNewSubpath(0, 0)},
		{name: "AppendStraightLineSegment", symbol: appendStraightLineSegmentSymbol, instance: NewAppendStraightLineSegment(10, 20)},
		{name: "ModifyCurrentTransformationMatrix", symbol: modifyCurrentTransformationMatrixSymbol, instance: NewModifyCurrentTransformationMatrix(1, 0, 0, 1, 0, 0)},
		{name: "SetLineWidth", symbol: setLineWidthSymbol, instance: NewSetLineWidth(1)},
		{name: "SetFlatnessTolerance", symbol: setFlatnessToleranceSymbol, instance: NewSetFlatnessTolerance(1)},
		{name: "SetMiterLimit", symbol: setMiterLimitSymbol, instance: NewSetMiterLimit(10)},
		{name: "SetColorRenderingIntent", symbol: setColorRenderingIntentSymbol, instance: NewSetColorRenderingIntent(hogName)},
		{name: "BeginMarkedContent", symbol: beginMarkedContentSymbol, instance: mustOpWriter(NewBeginMarkedContent(hogName))},
		{name: "DesignateMarkedContentPoint", symbol: designateMarkedContentPointSymbol, instance: mustOpWriter(NewDesignateMarkedContentPoint(hogName))},
		{name: "InvokeNamedXObject", symbol: invokeNamedXObjectSymbol, instance: mustOpWriter(NewInvokeNamedXObject(hogName))},
		{name: "PaintShading", symbol: paintShadingSymbol, instance: mustOpWriter(NewPaintShading(hogName))},
		{name: "SetNonStrokeColorSpace", symbol: setNonStrokeColorSpaceSymbol, instance: NewSetNonStrokeColorSpace(hogName)},
		{name: "SetStrokeColorSpace", symbol: setStrokeColorSpaceSymbol, instance: NewSetStrokeColorSpace(hogName)},
		{name: "SetGraphicsStateParametersFromDictionary", symbol: symbolSetGraphicsState, instance: NewSetGraphicsStateParametersFromDictionary(hogName)},
		{name: "ShowText", symbol: showTextSymbol, instance: NewShowText("Hello")},
		{name: "MoveToNextLineShowText", symbol: moveToNextLineShowTextSymbol, instance: NewMoveToNextLineShowText("Hello")},
		{name: "SetNonStrokeColorDeviceGray", symbol: setNonStrokeColorDeviceGraySymbol, instance: NewSetNonStrokeColorDeviceGray(0.5)},
		{name: "SetStrokeColorDeviceGray", symbol: setStrokeColorDeviceGraySymbol, instance: NewSetStrokeColorDeviceGray(0.5)},
		{name: "SetNonStrokeColorDeviceRgb", symbol: setNonStrokeColorDeviceRgbSymbol, instance: NewSetNonStrokeColorDeviceRgb(1, 0, 0)},
		{name: "SetStrokeColorDeviceRgb", symbol: setStrokeColorDeviceRgbSymbol, instance: NewSetStrokeColorDeviceRgb(1, 0, 0)},
		{name: "SetNonStrokeColorDeviceCmyk", symbol: setNonStrokeColorDeviceCmykSymbol, instance: NewSetNonStrokeColorDeviceCmyk(0, 0, 0, 1)},
		{name: "SetStrokeColorDeviceCmyk", symbol: setStrokeColorDeviceCmykSymbol, instance: NewSetStrokeColorDeviceCmyk(0, 0, 0, 1)},
		{name: "SetNonStrokeColor", symbol: setNonStrokeColorSymbol, instance: NewSetNonStrokeColor([]float64{0.5})},
		{name: "SetStrokeColor", symbol: setStrokeColorSymbol, instance: NewSetStrokeColor([]float64{0.5})},
		{name: "SetNonStrokeColorAdvanced", symbol: setNonStrokeColorAdvancedSymbol, instance: NewSetNonStrokeColorAdvanced([]float64{0.5})},
		{name: "SetStrokeColorAdvanced", symbol: setStrokeColorAdvancedSymbol, instance: NewSetStrokeColorAdvanced([]float64{0.5})},
		{name: "AppendRectangle", symbol: appendRectangleSymbol, instance: NewAppendRectangle(0, 0, 100, 100)},
		{name: "AppendStartControlPointBezierCurve", symbol: appendStartControlPointBezierCurveSymbol, instance: NewAppendStartControlPointBezierCurve(0, 0, 50, 50)},
		{name: "AppendEndControlPointBezierCurve", symbol: appendEndControlPointBezierCurveSymbol, instance: NewAppendEndControlPointBezierCurve(0, 0, 50, 50)},
		{name: "AppendDualControlPointBezierCurve", symbol: appendDualControlPointBezierCurveSymbol, instance: NewAppendDualControlPointBezierCurve(0, 0, 25, 25, 75, 75)},
		{name: "SetLineCap", symbol: setLineCapSymbol, instance: mustOpWriter(NewSetLineCap(graphiccore.LineCapStyle(0)))},
		{name: "SetLineJoin", symbol: setLineJoinSymbol, instance: mustOpWriter(NewSetLineJoin(graphiccore.LineJoinStyle(0)))},
		{name: "Type3SetGlyphWidth", symbol: type3SetGlyphWidthSymbol, instance: NewType3SetGlyphWidth(100, 0)},
		{name: "Type3SetGlyphWidthAndBoundingBox", symbol: type3SetGlyphWidthAndBoundingBoxSymbol, instance: NewType3SetGlyphWidthAndBoundingBox(100, 0, 0, 0, 100, 100)},
		{name: "EndInlineImage", symbol: endInlineImageSymbol, instance: NewEndInlineImage([]byte{})},
		{name: "BeginInlineImageData", symbol: beginInlineImageDataSymbol, instance: NewBeginInlineImageData(map[*tokens.NameToken]tokens.Token{})},
		{name: "SetLineDashPattern", symbol: setLineDashPatternSymbol, instance: mustOpWriter(NewSetLineDashPattern([]float64{1, 2, 3}, 0))},
		{name: "ShowTextsWithPositioning", symbol: showTextsWithPositioningSymbol, instance: mustOpWriter(NewShowTextsWithPositioning([]tokens.Token{tokens.NewStringToken("Hello")}))},
		{name: "BeginMarkedContentWithProperties", symbol: beginMarkedContentWithPropertiesSymbol, instance: mustOpWriter(NewBeginMarkedContentWithPropertiesName(tokens.MustCreate("Test"), tokens.MustCreate("Props")))},
		{name: "DesignateMarkedContentPointWithProperties", symbol: designateMarkedContentPointWithPropertiesSymbol, instance: mustOpWriter(NewDesignateMarkedContentPointWithPropertiesName(tokens.MustCreate("Test"), tokens.MustCreate("Props")))},
		{name: "MoveToNextLineShowTextWithSpacing", symbol: moveToNextLineShowTextWithSpacingSymbol, instance: NewMoveToNextLineShowTextWithSpacing(0, 0, "Hello")},
	}
}

// mustOpWriter unwraps a (opWriter, error) pair, panicking on error.
func mustOpWriter(v opWriter, err error) opWriter {
	if err != nil {
		panic(err)
	}
	return v
}

// TestAllOperationsCanBeWritten verifies every operation's Write method succeeds.
func TestAllOperationsCanBeWritten(t *testing.T) {
	ops := getAllOperations()
	for _, tc := range ops {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tc.instance.Write(&buf)
			if err != nil {
				t.Errorf("%s.Write() returned error: %v", tc.name, err)
			}
			if buf.Len() == 0 {
				t.Errorf("%s.Write() produced empty output", tc.name)
			}
		})
	}
}

// TestAllOperationsHaveSymbols verifies every operation type has an associated symbol constant.
func TestAllOperationsHaveSymbols(t *testing.T) {
	ops := getAllOperations()
	for _, tc := range ops {
		t.Run(tc.name, func(t *testing.T) {
			if tc.symbol == "" {
				t.Errorf("%s has no symbol defined", tc.name)
			}
		})
	}
}

// TestReflectionGraphicsStateOperationFactoryKnowsAllOperations verifies that
// the factory's Create method can handle every known operation symbol.
func TestReflectionGraphicsStateOperationFactoryKnowsAllOperations(t *testing.T) {
	allOps := getAllSymbols()

	factory := GetReflectionFactory()

	for symbol, opName := range allOps {
		t.Run(opName, func(t *testing.T) {
			opToken := tokens.CreateOperatorToken(symbol)
			if opToken == nil {
				t.Fatalf("could not create operator token for symbol %q", symbol)
			}

			unrecognized := false
			func() {
				defer func() {
					_ = recover() // ignore panics from operand access — means symbol was recognized
				}()
				_, err := factory.Create(opToken, []tokens.Token{})
				if err != nil && containsString(err.Error(), "no support implemented") {
					unrecognized = true
				}
			}()

			if unrecognized {
				t.Errorf("factory does not handle symbol %q for operation %s", symbol, opName)
			}
		})
	}
}

// getAllSymbols returns a map of every operation symbol to its name.
func getAllSymbols() map[string]string {
	result := make(map[string]string)

	ops := getAllOperations()
	for _, tc := range ops {
		result[tc.symbol] = tc.name
	}

	return result
}

// containsString checks if s contains substr.
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}
