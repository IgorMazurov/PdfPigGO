// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/content"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// ReflectionGraphicsStateOperationFactory creates graphics state operations from operator tokens
// and operands using a reflection-based dispatch approach. It maps PDF content stream operators
// to their corresponding Go operation types.
type ReflectionGraphicsStateOperationFactory struct{}

// reflectionFactory is the package-level singleton instance of the factory.
var reflectionFactory = &ReflectionGraphicsStateOperationFactory{}

// GetReflectionFactory returns the singleton instance of ReflectionGraphicsStateOperationFactory.
func GetReflectionFactory() *ReflectionGraphicsStateOperationFactory {
	return reflectionFactory
}

// Create dispatches an operator token and its operands to the appropriate operation constructor,
// returning a fully constructed graphics state operation or an error if parsing fails.
func (f *ReflectionGraphicsStateOperationFactory) Create(op *tokens.OperatorToken, operands []tokens.Token) (content.GraphicsStateOperation, error) {
	if op == nil {
		return nil, fmt.Errorf("operator token cannot be nil")
	}

	switch op.Data() {
	case modifyClippingByEvenOddIntersectSymbol:
		return InstanceModifyClippingByEvenOddIntersect, nil
	case modifyClippingByNonZeroWindingIntersectSymbol:
		return InstanceModifyClippingByNonZeroWindingIntersect, nil
	case beginCompatibilitySectionSymbol:
		return InstanceBeginCompatibilitySection, nil
	case endCompatibilitySectionSymbol:
		return InstanceEndCompatibilitySection, nil
	case setColorRenderingIntentSymbol:
		name, ok := operands[0].(*tokens.NameToken)
		if !ok {
			return nil, fmt.Errorf("expected NameToken for ri operator, got %T", operands[0])
		}
		return NewSetColorRenderingIntent(name), nil
	case setFlatnessToleranceSymbol:
		if len(operands) == 0 {
			return nil, nil
		}
		return NewSetFlatnessTolerance(operandToDouble(operands[0])), nil
	case setLineCapSymbol:
		capStyle := graphiccore.LineCapStyle(operandToInt(operands[0]))
		return NewSetLineCap(capStyle)
	case setLineDashPatternSymbol:
		arr := tokensToDoubleArray(operands, true)
		phase := operandToInt(operands[len(operands)-1])
		return NewSetLineDashPattern(arr, phase)
	case setLineJoinSymbol:
		joinStyle := graphiccore.LineJoinStyle(operandToInt(operands[0]))
		return NewSetLineJoin(joinStyle)
	case setLineWidthSymbol:
		return NewSetLineWidth(operandToDouble(operands[0])), nil
	case setMiterLimitSymbol:
		return NewSetMiterLimit(operandToDouble(operands[0])), nil
	case appendDualControlPointBezierCurveSymbol:
		if len(operands) == 0 {
			return nil, nil
		}
		return NewAppendDualControlPointBezierCurve(
			operandToDouble(operands[0]), operandToDouble(operands[1]),
			operandToDouble(operands[2]), operandToDouble(operands[3]),
			operandToDouble(operands[4]), operandToDouble(operands[5]),
		), nil
	case appendEndControlPointBezierCurveSymbol:
		if len(operands) == 0 {
			return nil, nil
		}
		return NewAppendEndControlPointBezierCurve(
			operandToDouble(operands[0]), operandToDouble(operands[1]),
			operandToDouble(operands[2]), operandToDouble(operands[3]),
		), nil
	case appendRectangleSymbol:
		if len(operands) == 0 {
			return nil, nil
		}
		return NewAppendRectangle(
			operandToDouble(operands[0]), operandToDouble(operands[1]),
			operandToDouble(operands[2]), operandToDouble(operands[3]),
		), nil
	case appendStartControlPointBezierCurveSymbol:
		if len(operands) == 0 {
			return nil, nil
		}
		return NewAppendStartControlPointBezierCurve(
			operandToDouble(operands[0]), operandToDouble(operands[1]),
			operandToDouble(operands[2]), operandToDouble(operands[3]),
		), nil
	case appendStraightLineSegmentSymbol:
		return NewAppendStraightLineSegment(
			operandToDouble(operands[0]), operandToDouble(operands[1]),
		), nil
	case beginNewSubpathSymbol:
		return NewBeginNewSubpath(
			operandToDouble(operands[0]), operandToDouble(operands[1]),
		), nil
	case closeSubpathSymbol:
		return InstanceCloseSubpath, nil
	case modifyCurrentTransformationMatrixSymbol:
		arr := tokensToDoubleArray(operands, false)
		if len(arr) < 6 {
			return nil, fmt.Errorf("expected at least 6 values for cm operator, got %d", len(arr))
		}
		return NewModifyCurrentTransformationMatrix(
			arr[0], arr[1], arr[2], arr[3], arr[4], arr[5],
		), nil
	case popSymbol:
		return InstancePop, nil
	case pushSymbol:
		return InstancePush, nil
	case symbolSetGraphicsState:
		name, ok := operands[0].(*tokens.NameToken)
		if !ok {
			return nil, fmt.Errorf("expected NameToken for gs operator, got %T", operands[0])
		}
		return NewSetGraphicsStateParametersFromDictionary(name), nil
	case beginTextSymbol:
		return InstanceBeginText, nil
	case endTextSymbol:
		return InstanceEndText, nil
	case setCharacterSpacingSymbol:
		return NewSetCharacterSpacing(operandToDouble(operands[0])), nil
	case setFontAndSizeSymbol:
		font, ok := operands[0].(*tokens.NameToken)
		if !ok {
			return nil, fmt.Errorf("expected NameToken for Tf operator font name, got %T", operands[0])
		}
		size := operandToDouble(operands[1])
		return NewSetFontAndSize(font, size), nil
	case setHorizontalScalingSymbol:
		return NewSetHorizontalScaling(operandToDouble(operands[0])), nil
	case setTextLeadingSymbol:
		return NewSetTextLeading(operandToDouble(operands[0])), nil
	case setTextRenderingModeSymbol:
		return NewSetTextRenderingMode(operandToInt(operands[0])), nil
	case setTextRiseSymbol:
		return NewSetTextRise(operandToDouble(operands[0])), nil
	case setWordSpacingSymbol:
		return NewSetWordSpacing(operandToDouble(operands[0])), nil
	case closeAndStrokePathSymbol:
		return InstanceCloseAndStrokePath, nil
	case closeFillEvenOddStrokePathSymbol:
		return InstanceCloseFillEvenOddStrokePath, nil
	case closeFillNonZeroStrokePathSymbol:
		return InstanceCloseFillNonZeroStrokePath, nil
	case beginInlineImageSymbol:
		return InstanceBeginInlineImage, nil
	case beginMarkedContentSymbol:
		name, ok := operands[0].(*tokens.NameToken)
		if !ok {
			return nil, fmt.Errorf("expected NameToken for BMC operator, got %T", operands[0])
		}
		return NewBeginMarkedContent(name)
	case beginMarkedContentWithPropertiesSymbol:
		bdcName, ok := operands[0].(*tokens.NameToken)
		if !ok {
			return nil, fmt.Errorf("expected NameToken for BDC operator name, got %T", operands[0])
		}
		if dict, ok := operands[1].(*tokens.DictionaryToken); ok {
			return NewBeginMarkedContentWithPropertiesDict(bdcName, dict)
		} else if name, ok := operands[1].(*tokens.NameToken); ok {
			return NewBeginMarkedContentWithPropertiesName(bdcName, name)
		}
		return nil, fmt.Errorf("attempted to set a marked-content sequence with invalid parameters: [%s]", printOperands(operands))
	case designateMarkedContentPointSymbol:
		name, ok := operands[0].(*tokens.NameToken)
		if !ok {
			return nil, fmt.Errorf("expected NameToken for MP operator, got %T", operands[0])
		}
		return NewDesignateMarkedContentPoint(name)
	case designateMarkedContentPointWithPropertiesSymbol:
		dpName, ok := operands[0].(*tokens.NameToken)
		if !ok {
			return nil, fmt.Errorf("expected NameToken for DP operator name, got %T", operands[0])
		}
		if dict, ok := operands[1].(*tokens.DictionaryToken); ok {
			return NewDesignateMarkedContentPointWithPropertiesDict(dpName, dict)
		} else if name, ok := operands[1].(*tokens.NameToken); ok {
			return NewDesignateMarkedContentPointWithPropertiesName(dpName, name)
		}
		return nil, fmt.Errorf("attempted to set a marked-content point with invalid parameters: [%s]", printOperands(operands))
	case endMarkedContentSymbol:
		return InstanceEndMarkedContent, nil
	case endPathSymbol:
		return InstanceEndPath, nil
	case fillPathEvenOddRuleSymbol:
		return InstanceFillPathEvenOddRule, nil
	case fillPathEvenOddRuleAndStrokeSymbol:
		return InstanceFillPathEvenOddRuleAndStroke, nil
	case fillPathNonZeroWindingSymbol:
		return InstanceFillPathNonZeroWinding, nil
	case fillPathNonZeroWindingAndStrokeSymbol:
		return InstanceFillPathNonZeroWindingAndStroke, nil
	case fillPathNonZeroWindingCompatSymbol:
		return InstanceFillPathNonZeroWindingCompatibility, nil
	case invokeNamedXObjectSymbol:
		name, ok := operands[0].(*tokens.NameToken)
		if !ok {
			return nil, fmt.Errorf("expected NameToken for Do operator, got %T", operands[0])
		}
		return NewInvokeNamedXObject(name)
	case moveToNextLineSymbol:
		return InstanceMoveToNextLine, nil
	case moveToNextLineShowTextSymbol:
		if len(operands) != 1 {
			return nil, fmt.Errorf("attempted to create a move to next line and show text operation with %d operands", len(operands))
		}
		if st, ok := operands[0].(*tokens.StringToken); ok {
			return NewMoveToNextLineShowText(st.Data()), nil
		}
		if ht, ok := operands[0].(*tokens.HexToken); ok {
			return NewMoveToNextLineShowTextBytes(ht.Bytes()), nil
		}
		return nil, fmt.Errorf("tried to create a move to next line and show text operation with operand type: %T", operands[0])
	case moveToNextLineShowTextWithSpacingSymbol:
		wordSpacing, ok1 := operands[0].(*tokens.NumericToken)
		charSpacing, ok2 := operands[1].(*tokens.NumericToken)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("expected NumericTokens for word and character spacing in \" operator")
		}
		text := operands[2]
		if st, ok := text.(*tokens.StringToken); ok {
			return NewMoveToNextLineShowTextWithSpacing(wordSpacing.Data(), charSpacing.Data(), st.Data()), nil
		}
		if ht, ok := text.(*tokens.HexToken); ok {
			return NewMoveToNextLineShowTextWithSpacingBytes(wordSpacing.Data(), charSpacing.Data(), ht.Bytes()), nil
		}
		return nil, fmt.Errorf("tried to create a MoveToNextLineShowTextWithSpacing operation with operand type: %T", text)
	case moveToNextLineWithOffsetSymbol:
		return NewMoveToNextLineWithOffset(
			operandToDouble(operands[0]), operandToDouble(operands[1]),
		), nil
	case moveToNextLineWithOffsetSetLeadingSymbol:
		return NewMoveToNextLineWithOffsetSetLeading(
			operandToDouble(operands[0]), operandToDouble(operands[1]),
		), nil
	case paintShadingSymbol:
		name, ok := operands[0].(*tokens.NameToken)
		if !ok {
			return nil, fmt.Errorf("expected NameToken for sh operator, got %T", operands[0])
		}
		return NewPaintShading(name)
	case setNonStrokeColorSymbol:
		return NewSetNonStrokeColor(tokensToDoubleArray(operands, false)), nil
	case setNonStrokeColorAdvancedSymbol:
		lastIdx := len(operands) - 1
		if name, ok := operands[lastIdx].(*tokens.NameToken); ok {
			nums := make([]float64, 0, lastIdx)
			for i := 0; i < lastIdx; i++ {
				if nt, ok := operands[i].(*tokens.NumericToken); ok {
					nums = append(nums, nt.Data())
				}
			}
			return NewSetNonStrokeColorAdvancedWithPattern(nums, name), nil
		}
		allNumeric := true
		for _, t := range operands {
			if _, ok := t.(*tokens.NumericToken); !ok {
				allNumeric = false
				break
			}
		}
		if allNumeric {
			nums := make([]float64, len(operands))
			for i, t := range operands {
				nums[i] = t.(*tokens.NumericToken).Data()
			}
			return NewSetNonStrokeColorAdvanced(nums), nil
		}
		return nil, fmt.Errorf("attempted to set a non-stroke color space (scn) with invalid arguments: [%s]", printOperands(operands))
	case setNonStrokeColorDeviceCmykSymbol:
		return NewSetNonStrokeColorDeviceCmyk(
			operandToDouble(operands[0]), operandToDouble(operands[1]),
			operandToDouble(operands[2]), operandToDouble(operands[3]),
		), nil
	case setNonStrokeColorDeviceGraySymbol:
		return NewSetNonStrokeColorDeviceGray(operandToDouble(operands[0])), nil
	case setNonStrokeColorDeviceRgbSymbol:
		return NewSetNonStrokeColorDeviceRgb(
			operandToDouble(operands[0]), operandToDouble(operands[1]),
			operandToDouble(operands[2]),
		), nil
	case setNonStrokeColorSpaceSymbol:
		name, ok := operands[0].(*tokens.NameToken)
		if !ok {
			return nil, fmt.Errorf("expected NameToken for cs operator, got %T", operands[0])
		}
		return NewSetNonStrokeColorSpace(name), nil
	case setStrokeColorSymbol:
		return NewSetStrokeColor(tokensToDoubleArray(operands, false)), nil
	case setStrokeColorAdvancedSymbol:
		lastIdx := len(operands) - 1
		if name, ok := operands[lastIdx].(*tokens.NameToken); ok {
			nums := make([]float64, 0, lastIdx)
			for i := 0; i < lastIdx; i++ {
				if nt, ok := operands[i].(*tokens.NumericToken); ok {
					nums = append(nums, nt.Data())
				}
			}
			return NewSetStrokeColorAdvancedWithPattern(nums, name), nil
		} else if allNumericTokens(operands) {
			nums := make([]float64, len(operands))
			for i, t := range operands {
				nums[i] = t.(*tokens.NumericToken).Data()
			}
			return NewSetStrokeColorAdvanced(nums), nil
		}
		return nil, fmt.Errorf("attempted to set a stroke color space (SCN) with invalid arguments: [%s]", printOperands(operands))
	case setStrokeColorDeviceCmykSymbol:
		args := getExpectedDoubles(setStrokeColorDeviceCmykSymbol, operands, 4)
		return NewSetStrokeColorDeviceCmyk(args[0], args[1], args[2], args[3]), nil
	case setStrokeColorDeviceGraySymbol:
		return NewSetStrokeColorDeviceGray(operandToDouble(operands[0])), nil
	case setStrokeColorDeviceRgbSymbol:
		return NewSetStrokeColorDeviceRgb(
			operandToDouble(operands[0]), operandToDouble(operands[1]),
			operandToDouble(operands[2]),
		), nil
	case setStrokeColorSpaceSymbol:
		name, ok := operands[0].(*tokens.NameToken)
		if !ok {
			return nil, fmt.Errorf("expected NameToken for CS operator, got %T", operands[0])
		}
		return NewSetStrokeColorSpace(name), nil
	case setTextMatrixSymbol:
		arr := tokensToDoubleArray(operands, false)
		if len(arr) < 6 {
			return nil, fmt.Errorf("expected at least 6 values for Tm operator, got %d", len(arr))
		}
		return NewSetTextMatrix(arr[0], arr[1], arr[2], arr[3], arr[4], arr[5]), nil
	case strokePathSymbol:
		return InstanceStrokePath, nil
	case showTextSymbol:
		if len(operands) != 1 {
			return nil, fmt.Errorf("attempted to create a show text operation with %d operands", len(operands))
		}
		if st, ok := operands[0].(*tokens.StringToken); ok {
			return NewShowText(st.Data()), nil
		}
		if ht, ok := operands[0].(*tokens.HexToken); ok {
			return NewShowTextBytes(ht.Bytes()), nil
		}
		return nil, fmt.Errorf("tried to create a show text operation with operand type: %T", operands[0])
	case showTextsWithPositioningSymbol:
		if len(operands) == 0 {
			return nil, fmt.Errorf("cannot have 0 parameters for a TJ operator")
		}
		if len(operands) == 1 {
			if arr, ok := operands[0].(*tokens.ArrayToken); ok {
				return NewShowTextsWithPositioning(arr.Data())
			}
		}
		return NewShowTextsWithPositioning(operands)
	case beginInlineImageDataSymbol:
		return nil, nil
	case endInlineImageSymbol:
		return nil, nil
	case type3SetGlyphWidthSymbol:
		args := getExpectedDoubles(type3SetGlyphWidthSymbol, operands, 2)
		return NewType3SetGlyphWidth(args[0], args[1]), nil
	case type3SetGlyphWidthAndBoundingBoxSymbol:
		args := getExpectedDoubles(type3SetGlyphWidthAndBoundingBoxSymbol, operands, 6)
		return NewType3SetGlyphWidthAndBoundingBox(
			args[0], args[1], args[2], args[3], args[4], args[5],
		), nil
	}

	// Unknown operators are silently ignored, matching C# behavior where
	// ReflectionGraphicsStateOperationFactory.Create returns null for operators
	// not in the Operations dictionary. The page content parser will log a warning.
	return nil, nil
}

// tokensToDoubleArray converts a list of token operands into a float64 slice.
// If exceptLast is true, the last operand is excluded from conversion.
func tokensToDoubleArray(tokensList []tokens.Token, exceptLast bool) []float64 {
	end := len(tokensList)
	if exceptLast {
		end--
	}

	result := make([]float64, 0, end)

	for i := 0; i < end; i++ {
		tok := tokensList[i]

		if arr, ok := tok.(*tokens.ArrayToken); ok {
			for _, inner := range arr.Data() {
				if nt, ok := inner.(*tokens.NumericToken); ok {
					result = append(result, nt.Data())
				} else {
					return result
				}
			}
			continue
		}

		if nt, ok := tok.(*tokens.NumericToken); ok {
			result = append(result, nt.Data())
		} else {
			return result
		}
	}

	return result
}

// operandToInt converts a token to an integer value.
func operandToInt(tok tokens.Token) int {
	if nt, ok := tok.(*tokens.NumericToken); ok {
		return nt.IntVal()
	}
	panic(fmt.Sprintf("invalid operand token encountered when expecting numeric: %v", tok))
}

// operandToDouble converts a token to a float64 value.
func operandToDouble(tok tokens.Token) float64 {
	if nt, ok := tok.(*tokens.NumericToken); ok {
		return nt.Data()
	}
	panic(fmt.Sprintf("invalid operand token encountered when expecting numeric: %v", tok))
}

// getExpectedDoubles validates that operands contain at least resultCount numeric values
// and returns them as a float64 slice.
func getExpectedDoubles(operatorSymbol string, operands []tokens.Token, resultCount int) []float64 {
	results := make([]float64, resultCount)
	if len(operands) < resultCount {
		panic(fmt.Sprintf("invalid operands for %s, needed %d numbers, got: %s", operatorSymbol, resultCount, printOperands(operands)))
	}
	for i := 0; i < resultCount; i++ {
		tok := operands[i]
		if nt, ok := tok.(*tokens.NumericToken); ok {
			results[i] = nt.Data()
		} else {
			panic(fmt.Sprintf("invalid operands for %s, needed %d numbers, got: %s", operatorSymbol, resultCount, printOperands(operands)))
		}
	}
	return results
}

// printOperands returns a string representation of the operand list.
func printOperands(operands []tokens.Token) string {
	parts := make([]string, len(operands))
	for i, t := range operands {
		parts[i] = fmt.Sprint(t)
	}
	return "[" + joinStrings(parts, ", ") + "]"
}

// allNumericTokens checks whether every token in the slice is a NumericToken.
func allNumericTokens(tokensList []tokens.Token) bool {
	for _, t := range tokensList {
		if _, ok := t.(*tokens.NumericToken); !ok {
			return false
		}
	}
	return true
}
