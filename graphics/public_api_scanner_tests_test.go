package graphics

import (
	"testing"
)

// TestAllSpecificationOperatorsArePresent verifies that every PDF specification
// graphics state operator has a corresponding implementation in this package.
// This is the Go equivalent of C# PublicApiScannerTests.AllSpecificationOperatorsArePresent().
func TestAllSpecificationOperatorsArePresent(t *testing.T) {
	expectedSymbols := []string{
		"b",
		"B",
		"b*",
		"B*",
		"BDC",
		"BI",
		"BMC",
		"BT",
		"BX",
		"c",
		"cm",
		"CS",
		"cs",
		"d",
		"d0",
		"d1",
		"Do",
		"DP",
		"EI",
		"EMC",
		"ET",
		"EX",
		"f",
		"F",
		"f*",
		"G",
		"g",
		"gs",
		"h",
		"i",
		"ID",
		"j",
		"J",
		"K",
		"k",
		"l",
		"m",
		"M",
		"MP",
		"n",
		"q",
		"Q",
		"re",
		"RG",
		"rg",
		"ri",
		"s",
		"S",
		"SC",
		"sc",
		"SCN",
		"scn",
		"sh",
		"T*",
		"Tc",
		"Td",
		"TD",
		"Tf",
		"Tj",
		"TJ",
		"TL",
		"Tm",
		"Tr",
		"Ts",
		"Tw",
		"Tz",
		"v",
		"w",
		"W",
		"W*",
		"y",
		"'",
		"\"",
	}

	actualSymbols := getAllSymbols()

	for _, expectedSymbol := range expectedSymbols {
		t.Run(expectedSymbol, func(t *testing.T) {
			if _, ok := actualSymbols[expectedSymbol]; !ok {
				t.Errorf("there is no operation defined with the symbol: %s", expectedSymbol)
			}
		})
	}
}
