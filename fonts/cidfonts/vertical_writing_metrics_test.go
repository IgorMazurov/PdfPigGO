package cidfonts

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/geometry"
)

func TestVerticalWritingMetricsUsesDefaultWhenOverridesNull(t *testing.T) {
	defaults := NewVerticalVectorComponents(250, 600)
	data := NewVerticalWritingMetrics(defaults, nil, nil)

	if len(data.IndividualVerticalWritingDisplacements) != 0 {
		t.Errorf("expected empty IndividualVerticalWritingDisplacements, got %d entries", len(data.IndividualVerticalWritingDisplacements))
	}
	if len(data.IndividualVerticalWritingPositions) != 0 {
		t.Errorf("expected empty IndividualVerticalWritingPositions, got %d entries", len(data.IndividualVerticalWritingPositions))
	}

	position := data.GetPositionVector(60, 250)
	if position.Y != defaults.Position {
		t.Errorf("position.Y = %g, want %g", position.Y, defaults.Position)
	}

	displacement := data.GetDisplacementVector(32)
	if displacement.Y != defaults.Displacement {
		t.Errorf("displacement.Y = %g, want %g", displacement.Y, defaults.Displacement)
	}
}

func TestVerticalWritingMetricsDefaultXComponentsOfVectorsAreCorrect(t *testing.T) {
	defaults := NewVerticalVectorComponents(250, 600)
	data := NewVerticalWritingMetrics(defaults, nil, nil)

	position := data.GetPositionVector(9, 120)
	expectedX := 120 / 2.0
	if position.X != expectedX {
		t.Errorf("position.X = %g, want %g", position.X, expectedX)
	}

	displacement := data.GetDisplacementVector(10)
	if displacement.X != 0.0 {
		t.Errorf("displacement.X = %g, want 0", displacement.X)
	}
}

func TestVerticalWritingMetricsUsesVectorOverridesWhenPresent(t *testing.T) {
	defaults := NewVerticalVectorComponents(250, 600)
	data := NewVerticalWritingMetrics(
		defaults,
		map[int]float64{7: 120},
		map[int]geometry.PdfVector{7: geometry.NewPdfVector(25, 250)},
	)

	position := data.GetPositionVector(7, 360)
	if position.X != 25 {
		t.Errorf("position.X = %g, want 25", position.X)
	}
	if position.Y != 250 {
		t.Errorf("position.Y = %g, want 250", position.Y)
	}

	displacement := data.GetDisplacementVector(7)
	if displacement.X != 0 {
		t.Errorf("displacement.X = %g, want 0", displacement.X)
	}
	if displacement.Y != 120 {
		t.Errorf("displacement.Y = %g, want 120", displacement.Y)
	}

	defaultPosition := data.GetPositionVector(6, 100)
	if defaultPosition.X != 50 {
		t.Errorf("defaultPosition.X = %g, want 50", defaultPosition.X)
	}
	if defaultPosition.Y != defaults.Position {
		t.Errorf("defaultPosition.Y = %g, want %g", defaultPosition.Y, defaults.Position)
	}

	defaultDisplacement := data.GetDisplacementVector(6)
	if defaultDisplacement.X != 0 {
		t.Errorf("defaultDisplacement.X = %g, want 0", defaultDisplacement.X)
	}
	if defaultDisplacement.Y != defaults.Displacement {
		t.Errorf("defaultDisplacement.Y = %g, want %g", defaultDisplacement.Y, defaults.Displacement)
	}
}
