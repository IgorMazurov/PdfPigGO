package functions

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

func createType4Function(functionText string, domain, rng []float64) *PdfFunctionType4 {
	domainArr := makeArrayToken(domain...)
	rangeArr := makeArrayToken(rng...)

	dictEntries := map[*tokens.NameToken]tokens.Token{
		tokens.FunctionType: tokens.NewNumericToken(4),
		tokens.Domain:       domainArr,
		tokens.Range:        rangeArr,
	}

	dict, err := tokens.NewDictionary(dictEntries)
	if err != nil {
		panic(err)
	}

	streamData := []byte(functionText)
	stream, err := tokens.NewStreamToken(dict, streamData)
	if err != nil {
		panic(err)
	}

	fn, err := NewPdfFunctionType4(stream, domainArr, rangeArr)
	if err != nil {
		panic(err)
	}

	return fn
}

func TestPdfFunctionType4Simple(t *testing.T) {
	fn := createType4Function("{ add }",
		[]float64{-1.0, 1.0, -1.0, 1.0},
		[]float64{-1.0, 1.0},
	)

	if fn.FunctionType() != PostScript {
		t.Errorf("expected FunctionType == PostScript, got %v", fn.FunctionType())
	}

	output := fn.Eval(0.8, 0.1)
	if len(output) != 1 {
		t.Fatalf("Eval(0.8, 0.1): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], 0.9)

	output = fn.Eval(0.8, 0.3)
	if len(output) != 1 {
		t.Fatalf("Eval(0.8, 0.3): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], 1)

	output = fn.Eval(0.8, 1.2)
	if len(output) != 1 {
		t.Fatalf("Eval(0.8, 1.2): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], 1)
}

func TestPdfFunctionType4ArgumentOrder(t *testing.T) {
	fn := createType4Function("{ pop }",
		[]float64{-1.0, 1.0, -1.0, 1.0},
		[]float64{-1.0, 1.0},
	)

	output := fn.Eval(-0.7, 0.0)
	if len(output) != 1 {
		t.Fatalf("Eval(-0.7, 0.0): expected length 1, got %d", len(output))
	}
	assertAlmostEqual4(t, output[0], -0.7)
}

func TestPdfFunctionType4Advanced(t *testing.T) {
	functionText := "{ dup 0.0 mul 1 exch sub 2 index 1.0 mul 1 exch sub mul  1 exch sub 3 1 roll dup 0.75 mul 1 exch sub 2 index 0.723 mul 1 exch sub mul  1 exch sub 3 1 roll dup 0.9 mul 1 exch sub 2 index 0.0 mul 1 exch sub mul  1 exch sub 3 1 roll dup 0.0 mul 1 exch sub 2 index 0.02 mul 1 exch sub mul  1 exch sub 3 1 roll pop pop }"

	fn := createType4Function(functionText,
		[]float64{0, 1, 0, 1},
		[]float64{0, 1, 0, 1, 0, 1, 0, 1},
	)

	output := fn.Eval(1.0, 1.0)

	if len(output) != 4 {
		t.Errorf("Eval(1.0, 1.0): expected length 4, got %d", len(output))
	}
}
