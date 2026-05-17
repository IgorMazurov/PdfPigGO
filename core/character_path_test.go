package core

import (
	"strings"
	"testing"
)

func TestBezierCurveGeneratesCorrectBoundingBox(t *testing.T) {
	curve := NewCubicBezierCurve(
		PdfPointFromInt(60, 105),
		PdfPointFromInt(75, 30),
		PdfPointFromInt(215, 115),
		PdfPointFromInt(140, 160),
	)

	result := curve.GetBoundingRectangle()

	if result == nil {
		t.Fatal("expected non-nil bounding rectangle")
	}

	if result.Top() != 160 {
		t.Errorf("expected Top == 160, got %g", result.Top())
	}

	if !(result.Bottom() < 105 && result.Bottom() > 30) {
		t.Errorf("expected Bottom between 30 and 105 (exclusive), got %g", result.Bottom())
	}

	if !(result.Right() > 140 && result.Right() < 215) {
		t.Errorf("expected Right between 140 and 215 (exclusive), got %g", result.Right())
	}

	if result.Left() != 60 {
		t.Errorf("expected Left == 60, got %g", result.Left())
	}
}

func TestLoopBezierCurveGeneratesCorrectBoundingBox(t *testing.T) {
	curve := NewCubicBezierCurve(
		PdfPointFromInt(166, 142),
		PdfPointFromInt(75, 30),
		PdfPointFromInt(215, 115),
		PdfPointFromInt(140, 160),
	)

	result := curve.GetBoundingRectangle()

	if result == nil {
		t.Fatal("expected non-nil bounding rectangle")
	}

	if result.Top() != 160 {
		t.Errorf("expected Top == 160, got %g", result.Top())
	}

	if !(result.Bottom() < 142 && result.Bottom() > 30) {
		t.Errorf("expected Bottom between 30 and 142 (exclusive), got %g", result.Bottom())
	}

	if result.Right() != 166 {
		t.Errorf("expected Right == 166, got %g", result.Right())
	}

	if !(result.Left() < 140) {
		t.Errorf("expected Left < 140, got %g", result.Left())
	}
}

func TestBezierCurveAddsCorrectSvgCommand(t *testing.T) {
	curve := NewCubicBezierCurve(
		PdfPointFromInt(60, 105),
		PdfPointFromInt(75, 30),
		PdfPointFromInt(215, 115),
		PdfPointFromInt(140, 160),
	)

	var builder strings.Builder
	curve.WriteSvg(&builder, 0)

	expected := "C 75 -30, 215 -115, 140 -160 "
	actual := builder.String()

	if actual != expected {
		t.Errorf("expected SVG command %q, got %q", expected, actual)
	}
}
