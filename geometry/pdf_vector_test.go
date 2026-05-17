package geometry_test

import (
	"math"
	"testing"

	"github.com/uglytoad/pdfpig/go/geometry"
)

func TestConstructorSetsValues(t *testing.T) {
	vector := geometry.NewPdfVector(5.2, 6.9)

	if vector.X != 5.2 {
		t.Errorf("expected X = 5.2, got %g", vector.X)
	}
	if vector.Y != 6.9 {
		t.Errorf("expected Y = 6.9, got %g", vector.Y)
	}
}

func TestScaleMultipliesLeavesOriginalUnchanged(t *testing.T) {
	vector := geometry.NewPdfVector(5.2, 6.9)

	scaled := vector.Scale(0.7)

	if vector.X != 5.2 {
		t.Errorf("expected original X unchanged at 5.2, got %g", vector.X)
	}
	expectedScaledX := 5.2 * 0.7
	if !floatEqual(scaled.X, expectedScaledX) {
		t.Errorf("expected scaled X = %g, got %g", expectedScaledX, scaled.X)
	}

	if vector.Y != 6.9 {
		t.Errorf("expected original Y unchanged at 6.9, got %g", vector.Y)
	}
	expectedScaledY := 6.9 * 0.7
	if !floatEqual(scaled.Y, expectedScaledY) {
		t.Errorf("expected scaled Y = %g, got %g", expectedScaledY, scaled.Y)
	}
}

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}
