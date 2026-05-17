package reading_order_detector_test

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	rod "github.com/uglytoad/pdfpig/go/document_layout_analysis/reading_order_detector"
)

// originTopLeft returns the point at coordinate top-left of an 800-unit page.
func originTopLeft() core.PdfPoint {
	return core.NewPdfPoint(0, 800)
}

// moveLeft returns a new point shifted "left" by dist units.
// In C# semantics, MoveX adds to X (positive X = left on screen).
func moveLeft(p core.PdfPoint, dist float64) core.PdfPoint {
	if dist < 0 {
		panic("dist must be positive")
	}
	return p.MoveX(dist)
}

// moveDown returns a new point shifted down by dist units.
func moveDown(p core.PdfPoint, dist float64) core.PdfPoint {
	if dist < 0 {
		panic("dist must be positive")
	}
	return p.MoveY(-dist)
}

// boxAtTopLeft creates a square rectangle anchored at the top-left origin.
func boxAtTopLeft(length float64) core.PdfRectangle {
	o := originTopLeft()
	return core.NewPdfRectangle(o, moveDown(moveLeft(o, length), length))
}

// rectMoveLeft returns a new rectangle shifted "left" by dist units.
// In C# semantics, left means increasing X.
func rectMoveLeft(r core.PdfRectangle, dist float64) core.PdfRectangle {
	if dist < 0 {
		panic("dist must be positive")
	}
	return r.Translate(dist, 0)
}

// rectMoveDown returns a new rectangle shifted down by dist units.
func rectMoveDown(r core.PdfRectangle, dist float64) core.PdfRectangle {
	if dist < 0 {
		panic("dist must be positive")
	}
	return r.Translate(0, -dist)
}

func TestIntervalRelationEqualsX(t *testing.T) {
	a := core.NewPdfRectangle(core.Origin, core.NewPdfPoint(10, 10))

	res := rod.GetRelationX(a, a, 5)

	if res != rod.Equals {
		t.Errorf("expected Equals, got %v", res)
	}
}

func TestIntervalRelationEqualsY(t *testing.T) {
	a := core.NewPdfRectangle(core.Origin, core.NewPdfPoint(10, 10))

	res := rod.GetRelationY(a, a, 5)

	if res != rod.Equals {
		t.Errorf("expected Equals, got %v", res)
	}
}

func TestIntervalRelationPrecedesX(t *testing.T) {
	a := boxAtTopLeft(10)
	b := rectMoveLeft(boxAtTopLeft(10), 100)

	res := rod.GetRelationX(a, b, 5)
	resInverse := rod.GetRelationX(b, a, 5)

	if res != rod.Precedes {
		t.Errorf("expected Precedes, got %v", res)
	}
	if resInverse != rod.PrecedesI {
		t.Errorf("expected PrecedesI, got %v", resInverse)
	}
}

func TestIntervalRelationPrecedesY(t *testing.T) {
	a := boxAtTopLeft(10)
	b := rectMoveDown(a, 200)

	res := rod.GetRelationY(a, b, 5)
	resInverse := rod.GetRelationY(b, a, 5)

	if res != rod.Precedes {
		t.Errorf("expected Precedes, got %v", res)
	}
	if resInverse != rod.PrecedesI {
		t.Errorf("expected PrecedesI, got %v", resInverse)
	}
}

func TestIntervalRelationMeetsX(t *testing.T) {
	a := boxAtTopLeft(100)
	b := rectMoveLeft(a, 100)

	res := rod.GetRelationX(a, b, 5)
	resInverse := rod.GetRelationX(b, a, 5)

	if res != rod.Meets {
		t.Errorf("expected Meets, got %v", res)
	}
	if resInverse != rod.MeetsI {
		t.Errorf("expected MeetsI, got %v", resInverse)
	}
}

func TestIntervalRelationMeetsXWithinTolerance(t *testing.T) {
	a := boxAtTopLeft(100)
	b := rectMoveLeft(a, 110)

	res := rod.GetRelationX(a, b, 11)
	resInverse := rod.GetRelationX(b, a, 11)

	if res != rod.Meets {
		t.Errorf("expected Meets, got %v", res)
	}
	if resInverse != rod.MeetsI {
		t.Errorf("expected MeetsI, got %v", resInverse)
	}
}

func TestIntervalRelationMeetsY(t *testing.T) {
	a := boxAtTopLeft(100)
	b := rectMoveDown(a, 100)

	res := rod.GetRelationY(a, b, 5)
	resInverse := rod.GetRelationY(b, a, 5)

	if res != rod.Meets {
		t.Errorf("expected Meets, got %v", res)
	}
	if resInverse != rod.MeetsI {
		t.Errorf("expected MeetsI, got %v", resInverse)
	}
}

func TestIntervalRelationMeetsYWhenMovedDownBecomesPrecedes(t *testing.T) {
	startPoint := core.NewPdfPoint(100, 600)
	a := core.NewPdfRectangle(startPoint, moveDown(startPoint, 100))
	meetsABox := rectMoveDown(a, 100)

	res := rod.GetRelationY(a, meetsABox, 5)
	resInverse := rod.GetRelationY(meetsABox, a, 5)

	if res != rod.Meets {
		t.Errorf("expected Meets, got %v", res)
	}
	if resInverse != rod.MeetsI {
		t.Errorf("expected MeetsI, got %v", resInverse)
	}

	preceededByABox := rectMoveDown(meetsABox, 100)

	moveRes := rod.GetRelationY(a, preceededByABox, 5)
	moveResInverse := rod.GetRelationY(preceededByABox, a, 5)

	if moveRes != rod.Precedes {
		t.Errorf("expected Precedes after further move down, got %v", moveRes)
	}
	if moveResInverse != rod.PrecedesI {
		t.Errorf("expected PrecedesI after further move down, got %v", moveResInverse)
	}
}

func TestIntervalRelationOverlapsX(t *testing.T) {
	a := boxAtTopLeft(100)
	b := rectMoveLeft(a, a.Width/2)

	res := rod.GetRelationX(a, b, 5)
	resInverse := rod.GetRelationX(b, a, 5)

	if res != rod.Overlaps {
		t.Errorf("expected Overlaps, got %v", res)
	}
	if resInverse != rod.OverlapsI {
		t.Errorf("expected OverlapsI, got %v", resInverse)
	}
}

func TestIntervalRelationOverlapsY(t *testing.T) {
	a := boxAtTopLeft(100)
	b := rectMoveDown(rectMoveLeft(a, 500), a.Height/2)

	res := rod.GetRelationY(a, b, 5)
	resInverse := rod.GetRelationY(b, a, 5)

	if res != rod.Overlaps {
		t.Errorf("expected Overlaps, got %v", res)
	}
	if resInverse != rod.OverlapsI {
		t.Errorf("expected OverlapsI, got %v", resInverse)
	}
}

func TestIntervalRelationStartsX(t *testing.T) {
	topLeft := originTopLeft()
	a := core.NewPdfRectangle(topLeft, moveDown(moveLeft(topLeft, 50), 10))
	b := core.NewPdfRectangle(topLeft, moveDown(moveLeft(topLeft, 100), 10))

	res := rod.GetRelationX(a, b, 5)
	resInverse := rod.GetRelationX(b, a, 5)

	if res != rod.Starts {
		t.Errorf("expected Starts, got %v", res)
	}
	if resInverse != rod.StartsI {
		t.Errorf("expected StartsI, got %v", resInverse)
	}
}

func TestIntervalRelationStartsY(t *testing.T) {
	topLeft := originTopLeft()
	a := core.NewPdfRectangle(topLeft, moveDown(moveLeft(topLeft, 100), 100))
	b := core.NewPdfRectangle(topLeft, moveDown(moveLeft(topLeft, 100), 200))

	res := rod.GetRelationY(a, b, 5)
	resInverse := rod.GetRelationY(b, a, 5)

	if res != rod.Starts {
		t.Errorf("expected Starts, got %v", res)
	}
	if resInverse != rod.StartsI {
		t.Errorf("expected StartsI, got %v", resInverse)
	}
}

func TestIntervalRelationDuringX(t *testing.T) {
	a := core.NewPdfRectangle(core.NewPdfPoint(20, 0), core.NewPdfPoint(80, 0))
	b := core.NewPdfRectangle(core.NewPdfPoint(0, 0), core.NewPdfPoint(100, 0))

	res := rod.GetRelationX(a, b, 5)
	resInverse := rod.GetRelationX(b, a, 5)

	if res != rod.During {
		t.Errorf("expected During, got %v", res)
	}
	if resInverse != rod.DuringI {
		t.Errorf("expected DuringI, got %v", resInverse)
	}
}

func TestIntervalRelationDuringY(t *testing.T) {
	a := core.NewPdfRectangle(core.NewPdfPoint(0, 20), core.NewPdfPoint(0, 80))
	b := core.NewPdfRectangle(core.NewPdfPoint(0, 0), core.NewPdfPoint(0, 100))

	res := rod.GetRelationY(a, b, 5)
	resInverse := rod.GetRelationY(b, a, 5)

	if res != rod.During {
		t.Errorf("expected During, got %v", res)
	}
	if resInverse != rod.DuringI {
		t.Errorf("expected DuringI, got %v", resInverse)
	}
}

func TestIntervalRelationFinishesX(t *testing.T) {
	topRight := moveLeft(originTopLeft(), 400)
	a := core.NewPdfRectangle(topRight.MoveX(-100), topRight)
	b := core.NewPdfRectangle(topRight.MoveX(-200), topRight)

	res := rod.GetRelationX(a, b, 5)
	resInverse := rod.GetRelationX(b, a, 5)

	if res != rod.Finishes {
		t.Errorf("expected Finishes, got %v", res)
	}
	if resInverse != rod.FinishesI {
		t.Errorf("expected FinishesI, got %v", resInverse)
	}
}

func TestIntervalRelationFinishesY(t *testing.T) {
	a := rectMoveDown(boxAtTopLeft(20), 20)
	b := boxAtTopLeft(40)

	res := rod.GetRelationY(a, b, 5)
	resInverse := rod.GetRelationY(b, a, 5)

	if res != rod.Finishes {
		t.Errorf("expected Finishes, got %v", res)
	}
	if resInverse != rod.FinishesI {
		t.Errorf("expected FinishesI, got %v", resInverse)
	}
}
