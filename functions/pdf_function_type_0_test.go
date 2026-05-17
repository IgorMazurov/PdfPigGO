package functions

import (
	"encoding/binary"
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

func makeArrayToken(vals ...float64) *tokens.ArrayToken {
	toks := make([]tokens.Token, len(vals))
	for i, v := range vals {
		toks[i] = tokens.NewNumericToken(v)
	}
	arr := tokens.NewArrayToken(toks)
	return arr
}

func createType0Function(domain, rng, size *tokens.ArrayToken, bitsPerSample int, encode, decode *tokens.ArrayToken, data []byte) (*PdfFunctionType0, error) {
	dictEntries := map[*tokens.NameToken]tokens.Token{
		tokens.FunctionType:    tokens.NewNumericToken(0),
		tokens.Domain:          domain,
		tokens.Range:           rng,
		tokens.BitsPerSample:   tokens.NewNumericToken(float64(bitsPerSample)),
		tokens.Size:            size,
	}
	if encode != nil {
		dictEntries[tokens.Encode] = encode
	}
	if decode != nil {
		dictEntries[tokens.Decode] = decode
	}

	dict, err := tokens.NewDictionary(dictEntries)
	if err != nil {
		return nil, err
	}

	stream, err := tokens.NewStreamToken(dict, data)
	if err != nil {
		return nil, err
	}

	order := 1
	return NewPdfFunctionType0FromStream(stream, domain, rng, size, bitsPerSample, order, encode, decode), nil
}

func assertAlmostEqual(t *testing.T, got, want float64, precision int) {
	t.Helper()
	diff := got - want
	if diff < 0 {
		diff = -diff
	}
	tolerance := float64(1)
	for i := 0; i < precision; i++ {
		tolerance /= 10
	}
	if diff > tolerance {
		t.Errorf("expected %.6f, got %.6f (tolerance %.6f)", want, got, tolerance)
	}
}

func assertEqualSlice(t *testing.T, got, want []float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("expected length %d, got %d", len(want), len(got))
		return
	}
	for i := range got {
		assertAlmostEqual(t, got[i], want[i], 4)
	}
}

func TestTIKA_1228_0(t *testing.T) {
	domain := makeArrayToken(0, 1)
	rng := makeArrayToken(0, 1, 0, 1, 0, 1, 0, 1)
	bitsPerSample := 8
	decodeArr := makeArrayToken(0, 1, 0, 1, 0, 1, 0, 1)
	encodeArr := makeArrayToken(0, 254)
	sizeArr := makeArrayToken(255)

	data := make([]byte, 255*4)
	for i := 0; i < 255; i++ {
		data[i*4+3] = byte(i)
	}

	fn, err := createType0Function(domain, rng, sizeArr, bitsPerSample, encodeArr, decodeArr, data)
	if err != nil {
		t.Fatalf("createType0Function error: %v", err)
	}

	if fn.FunctionType() != Sampled {
		t.Errorf("expected FunctionType == Sampled, got %v", fn.FunctionType())
	}

	result := fn.Eval(0)
	if len(result) != 4 {
		t.Errorf("Eval(0): expected length 4, got %d", len(result))
	}

	result = fn.Eval(0.5)
	if len(result) != 4 {
		t.Errorf("Eval(0.5): expected length 4, got %d", len(result))
	}

	result = fn.Eval(1)
	if len(result) != 4 {
		t.Errorf("Eval(1): expected length 4, got %d", len(result))
	}

	result = fn.Eval(0.2)
	if len(result) != 4 {
		t.Errorf("Eval(0.2): expected length 4, got %d", len(result))
	}
}

func TestSimple16(t *testing.T) {
	domain := makeArrayToken(0, 1)
	rng := makeArrayToken(0, 1)
	sizeArr := makeArrayToken(5)

	values := []uint16{0, 8192, 16384, 32768, 65535}
	data := make([]byte, 0, len(values)*2)
	for _, v := range values {
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, v)
		data = append(data, b...)
	}

	fn, err := createType0Function(domain, rng, sizeArr, 16, nil, nil, data)
	if err != nil {
		t.Fatalf("createType0Function error: %v", err)
	}

	if fn.FunctionType() != Sampled {
		t.Errorf("expected FunctionType == Sampled, got %v", fn.FunctionType())
	}

	result := fn.Eval(0.00)
	if len(result) != 1 {
		t.Fatalf("Eval(0.00): expected length 1, got %d", len(result))
	}
	assertAlmostEqual(t, result[0], 0.0, 3)

	result = fn.Eval(0.25)
	if len(result) != 1 {
		t.Fatalf("Eval(0.25): expected length 1, got %d", len(result))
	}
	assertAlmostEqual(t, result[0], 0.125, 3)

	result = fn.Eval(0.50)
	if len(result) != 1 {
		t.Fatalf("Eval(0.50): expected length 1, got %d", len(result))
	}
	assertAlmostEqual(t, result[0], 0.25, 2)

	result = fn.Eval(0.75)
	if len(result) != 1 {
		t.Fatalf("Eval(0.75): expected length 1, got %d", len(result))
	}
	assertAlmostEqual(t, result[0], 0.50, 2)

	result = fn.Eval(1.0)
	if len(result) != 1 {
		t.Fatalf("Eval(1.0): expected length 1, got %d", len(result))
	}
	assertAlmostEqual(t, result[0], 1.00, 2)
}

func TestSimple8(t *testing.T) {
	domain := makeArrayToken(0, 1)
	rng := makeArrayToken(0, 1)
	sizeArr := makeArrayToken(5)

	data := []byte{0, 32, 64, 128, 255}

	fn, err := createType0Function(domain, rng, sizeArr, 8, nil, nil, data)
	if err != nil {
		t.Fatalf("createType0Function error: %v", err)
	}

	if fn.FunctionType() != Sampled {
		t.Errorf("expected FunctionType == Sampled, got %v", fn.FunctionType())
	}

	result := fn.Eval(0.00)
	if len(result) != 1 {
		t.Fatalf("Eval(0.00): expected length 1, got %d", len(result))
	}
	assertAlmostEqual(t, result[0], 0.0, 3)

	result = fn.Eval(0.25)
	if len(result) != 1 {
		t.Fatalf("Eval(0.25): expected length 1, got %d", len(result))
	}
	assertAlmostEqual(t, result[0], 0.125, 3)

	result = fn.Eval(0.50)
	if len(result) != 1 {
		t.Fatalf("Eval(0.50): expected length 1, got %d", len(result))
	}
	assertAlmostEqual(t, result[0], 0.25, 2)

	result = fn.Eval(0.75)
	if len(result) != 1 {
		t.Fatalf("Eval(0.75): expected length 1, got %d", len(result))
	}
	assertAlmostEqual(t, result[0], 0.50, 2)

	result = fn.Eval(1.0)
	if len(result) != 1 {
		t.Fatalf("Eval(1.0): expected length 1, got %d", len(result))
	}
	assertAlmostEqual(t, result[0], 1.00, 2)
}

func TestRgbColorSpace(t *testing.T) {
	domain := makeArrayToken(0, 1, 0, 1)
	rng := makeArrayToken(0, 1, 0, 1, 0, 1)
	sizeArr := makeArrayToken(2, 2)

	data := []byte{255, 255, 0, 0, 0, 0, 255, 0, 0, 0, 0, 255}

	fn, err := createType0Function(domain, rng, sizeArr, 8, nil, nil, data)
	if err != nil {
		t.Fatalf("createType0Function error: %v", err)
	}

	if fn.FunctionType() != Sampled {
		t.Errorf("expected FunctionType == Sampled, got %v", fn.FunctionType())
	}

	result := fn.Eval(0, 0)
	if len(result) != 3 {
		t.Fatalf("Eval(0,0): expected length 3, got %d", len(result))
	}
	assertEqualSlice(t, result, []float64{1, 1, 0})

	result = fn.Eval(1, 0)
	if len(result) != 3 {
		t.Fatalf("Eval(1,0): expected length 3, got %d", len(result))
	}
	assertEqualSlice(t, result, []float64{0, 0, 0})

	result = fn.Eval(0, 1)
	if len(result) != 3 {
		t.Fatalf("Eval(0,1): expected length 3, got %d", len(result))
	}
	assertEqualSlice(t, result, []float64{1, 0, 0})

	result = fn.Eval(1, 1)
	if len(result) != 3 {
		t.Fatalf("Eval(1,1): expected length 3, got %d", len(result))
	}
	assertEqualSlice(t, result, []float64{0, 0, 1})

	result = fn.Eval(0.5, 0.5)
	if len(result) != 3 {
		t.Fatalf("Eval(0.5,0.5): expected length 3, got %d", len(result))
	}
	assertEqualSlice(t, result, []float64{0.5, 0.25, 0.25})
}

func TestRedBlueGradient(t *testing.T) {
	domain := makeArrayToken(0, 1)
	rng := makeArrayToken(0, 1, 0, 1, 0, 1)
	sizeArr := makeArrayToken(2)

	data := []byte{255, 0, 0, 0, 0, 255}

	fn, err := createType0Function(domain, rng, sizeArr, 8, nil, nil, data)
	if err != nil {
		t.Fatalf("createType0Function error: %v", err)
	}

	if fn.FunctionType() != Sampled {
		t.Errorf("expected FunctionType == Sampled, got %v", fn.FunctionType())
	}

	result := fn.Eval(0)
	if len(result) != 3 {
		t.Fatalf("Eval(0): expected length 3, got %d", len(result))
	}
	assertEqualSlice(t, result, []float64{1, 0, 0})

	result = fn.Eval(1)
	if len(result) != 3 {
		t.Fatalf("Eval(1): expected length 3, got %d", len(result))
	}
	assertEqualSlice(t, result, []float64{0, 0, 1})

	result = fn.Eval(0.5)
	if len(result) != 3 {
		t.Fatalf("Eval(0.5): expected length 3, got %d", len(result))
	}
	assertEqualSlice(t, result, []float64{0.5, 0.0, 0.5})

	result = fn.Eval(0.3333)
	if len(result) != 3 {
		t.Fatalf("Eval(0.3333): expected length 3, got %d", len(result))
	}
	assertEqualSlice(t, result, []float64{0.6667, 0.0, 0.3333})
}
