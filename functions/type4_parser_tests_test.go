// Package functions provides tests for the Type 4 PostScript function parser.
package functions

import "testing"

func TestParserBasics(t *testing.T) {
	CreateType4Tester("3 4 add 2 sub").PopInt(t, 5).IsEmpty(t)
}

func TestNested(t *testing.T) {
	CreateType4Tester("true { 2 1 add } { 2 1 sub } ifelse").
		PopInt(t, 3).IsEmpty(t)
	CreateType4Tester("{ true }").
		PopBool(t, true).IsEmpty(t)
}

func TestJira804(t *testing.T) {
	// Tint to CMYK function from PDFBOX-804.
	// Problems were: no whitespace between "mul" and "}", line breaks cause endless loops.
	CreateType4Tester("1 {dup dup .72 mul exch 0 exch .38 mul}\n").
		PopRealWithDelta(t, 0.38, 0.001).PopReal(t, 0).PopRealWithDelta(t, 0.72, 0.001).Pop(t, 1.0).IsEmpty(t)
}
