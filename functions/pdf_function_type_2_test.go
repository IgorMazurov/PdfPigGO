package functions

import (
	"math"
	"testing"
)

func createType2Function(domain, rng, c0, c1 []float64, n float64) *PdfFunctionType2 {
	domainArr := makeArrayToken(domain...)
	rangeArr := makeArrayToken(rng...)
	c0Arr := makeArrayToken(c0...)
	c1Arr := makeArrayToken(c1...)

	return NewPdfFunctionType2FromDict(nil, domainArr, rangeArr, c0Arr, c1Arr, n)
}

func assertAlmostEqual4(t *testing.T, got, want float64) {
	t.Helper()
	assertAlmostEqual(t, got, want, 4)
}

func TestPdfFunctionType2Simple(t *testing.T) {
	fn := createType2Function(
		[]float64{-1.0, 1.0, -1.0, 1.0},
		[]float64{-1.0, 1.0},
		[]float64{0.0},
		[]float64{1.0},
		1,
	)

	if fn.FunctionType() != Exponential {
		t.Errorf("expected FunctionType == Exponential, got %v", fn.FunctionType())
	}

	cases := []struct {
		input float64
		want  float64
	}{
		{-0.7, -0.7},
		{0.7, 0.7},
		{-0.5, -0.5},
		{0.5, 0.5},
		{0, 0},
		{1, 1},
		{-1, -1},
	}

	for _, c := range cases {
		output := fn.Eval(c.input)
		if len(output) != 1 {
			t.Fatalf("Eval(%g): expected length 1, got %d", c.input, len(output))
		}
		assertAlmostEqual4(t, output[0], c.want)
	}
}

func TestPdfFunctionType2SimpleClip(t *testing.T) {
	fn := createType2Function(
		[]float64{-1.0, 1.0, -1.0, 1.0},
		[]float64{-1.0, 1.0},
		[]float64{0.0},
		[]float64{1.0},
		1,
	)

	output := fn.Eval(-15)
	if len(output) != 1 {
		t.Fatalf("Eval(-15): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], -1)

	output = fn.Eval(15)
	if len(output) != 1 {
		t.Fatalf("Eval(15): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], 1)
}

func TestPdfFunctionType2N2(t *testing.T) {
	fn := createType2Function(
		[]float64{-1.0, 1.0, -1.0, 1.0},
		[]float64{-10.0, 10.0},
		[]float64{0.0},
		[]float64{1.0},
		2,
	)

	output := fn.Eval(1.12)
	if len(output) != 1 {
		t.Fatalf("Eval(1.12): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], 1.2544)

	output = fn.Eval(-1.35)
	if len(output) != 1 {
		t.Fatalf("Eval(-1.35): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], 1.82250)

	output = fn.Eval(5)
	if len(output) != 1 {
		t.Fatalf("Eval(5): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], 10)
}

func TestPdfFunctionType2N3(t *testing.T) {
	fn := createType2Function(
		[]float64{-1.0, 1.0, -1.0, 1.0},
		[]float64{-10.0, 10.0},
		[]float64{4.0},
		[]float64{9.53},
		3,
	)

	output := fn.Eval(1.0)
	if len(output) != 1 {
		t.Fatalf("Eval(1.0): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], 9.53)

	output = fn.Eval(-1.236)
	if len(output) != 1 {
		t.Fatalf("Eval(-1.236): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], -6.44192)
}

func TestPdfFunctionType2NSqrt(t *testing.T) {
	fn := createType2Function(
		[]float64{-1.0, 1.0, -1.0, 1.0},
		[]float64{-10.0, 10.0},
		[]float64{2.589},
		[]float64{10.58},
		0.5,
	)

	output := fn.Eval(0.5)
	if len(output) != 1 {
		t.Fatalf("Eval(0.5): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], 8.23949)

	output = fn.Eval(0.78)
	if len(output) != 1 {
		t.Fatalf("Eval(0.78): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], 9.64646)

	output = fn.Eval(-0.78)
	if len(output) != 1 {
		t.Fatalf("Eval(-0.78): expected length 1, got %d", len(output))
	}
	if !math.IsNaN(output[0]) {
		t.Errorf("expected NaN for negative input with sqrt, got %g", output[0])
	}
}
