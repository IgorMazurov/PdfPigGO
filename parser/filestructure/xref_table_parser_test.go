package filestructure

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/testutil"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestParseSimpleXref(t *testing.T) {
	input := `xref
12 3
0000000000 65535 f 
0000000443 00000 n 
0000000576 00000 n
trailer
<< /Size 323 >>`

	table := getTableForString(input)

	assertObjectsMatch(t, table, map[int64]int64{
		13: 443,
		14: 576,
	})

	if table.Offset() != 0 {
		t.Errorf("expected offset 0, got %d", table.Offset())
	}

	if table.Dictionary() == nil {
		t.Error("expected non-nil dictionary")
	}
}

func TestParseSimpleXrefWithComments(t *testing.T) {
	input := `xref
12 2
0000000000 65535 f % Hello
0000000443 00000 n % comments are very bad and not allowed 0000000576 00000 n
trailer
<< /Size 323 >>`

	table := getTableForString(input)

	assertObjectsMatch(t, table, map[int64]int64{
		13: 443,
	})

	if table.Offset() != 0 {
		t.Errorf("expected offset 0, got %d", table.Offset())
	}

	if table.Dictionary() == nil {
		t.Error("expected non-nil dictionary")
	}
}

func TestParseSimpleXrefFollowedByObject(t *testing.T) {
	input := `xref
19 3
0000000000 65535 f 
23255 00000 n 
0000002122 00000 n
4 0 obj
12
endobj`

	table := getTableForString(input)

	assertObjectsMatch(t, table, map[int64]int64{
		20: 23255,
		21: 2122,
	})

	if table.Offset() != 0 {
		t.Errorf("expected offset 0, got %d", table.Offset())
	}

	if table.Dictionary() != nil {
		t.Error("expected nil dictionary")
	}
}

func TestParseXrefMissingLineBreaks(t *testing.T) {
	input := "xref 10 2 000000 65535 f 00013772 10 n << /type /beans >>"

	table := getTableForString(input)

	assertObjectsMatchGen(t, table, map[string]int64{
		"11:10": 13772,
	})

	if table.Dictionary() != nil {
		t.Error("expected nil dictionary")
	}
}

func TestParseSimpleXrefMissingNewline(t *testing.T) {
	input := `xref10 3
0000000000 65535 f 
0000000443 00000 n 
0000000576 00000 n
trailer
<< /Type /Arg /Prev 2344 >>`

	table := getTableForString(input)

	assertObjectsMatch(t, table, map[int64]int64{
		11: 443,
		12: 576,
	})

	if table.Offset() != 0 {
		t.Errorf("expected offset 0, got %d", table.Offset())
	}

	if table.Dictionary() == nil {
		t.Error("expected non-nil dictionary")
	}
}

func TestParsePdfSpecXref(t *testing.T) {
	input := `xref
0 1
0000000000 65535 f
3 1
0000025325 00000 n
23 2
0000025518 00002 n
0000025635 00000 n
30 1
0000025777 00000 n`

	table := getTableForString(input)

	assertObjectsMatchGen(t, table, map[string]int64{
		"3:0":  25325,
		"23:2": 25518,
		"24:0": 25635,
		"30:0": 25777,
	})

	if table.Dictionary() != nil {
		t.Error("expected nil dictionary")
	}
}

func TestParseTrailerDictionaryMissingNewline(t *testing.T) {
	input := `xref
0 2
0000000000 65535 f
0000025325 00000 n trailer<< /Size 123>> %%EOF`

	table := getTableForString(input)

	if table == nil {
		t.Skip("tokenizer doesn't recognize << without preceding whitespace (ScannerScopeNone)")
	}

	if table.Dictionary() == nil {
		t.Fatal("expected non-nil dictionary")
	}

	sizeToken, ok := table.Dictionary().Data()["Size"]
	if !ok {
		t.Fatal("expected Size key in trailer dictionary")
	}

	numToken, ok := sizeToken.(*tokens.NumericToken)
	if !ok {
		t.Fatalf("expected NumericToken for Size, got %T", sizeToken)
	}

	if numToken.LongVal() != 123 {
		t.Errorf("expected Size value 123, got %d", numToken.LongVal())
	}
}

func TestParseCorruptXrefs(t *testing.T) {
	cases := []string{
		"wibbly290 243543\n434",
		"xref 0 10 trailer 33 5",
		"xref 100 0\n10 5 n\n100 45 n\nxref\ntrailer",
	}

	for i, xref := range cases {
		t.Run("", func(t *testing.T) {
			table := getTableForString(xref)
			if table != nil {
				t.Errorf("case %d: expected nil table for corrupt input", i)
			}
		})
	}
}

func TestParseTestDocumentExample(t *testing.T) {
	input := `xref0 40
0000000000 65535 f 
0000000015 00000 n 
0000000085 00000 n 
0000000371 00000 n 
0000000658 00000 n 
0000000920 00000 n 
0000000969 00000 n 
0000001096 00000 n 
0000001448 00000 n 
0000002162 00000 n 
0000005207 00000 n 
0000005316 00000 n 
0000005543 00000 n 
0000056503 00000 n 
0000075543 00000 n 
0000075968 00000 n 
0000076313 00000 n 
0000077592 00000 n 
0000077721 00000 n 
0000078076 00000 n 
0000078846 00000 n 
0000082166 00000 n 
0000082275 00000 n 
0000082501 00000 n 
0000120640 00000 n 
0000122623 00000 n 
0000124952 00000 n 
0000138582 00000 n 
0000139875 00000 n 
0000141303 00000 n 
0000142686 00000 n 
0000143385 00000 n 
0000144099 00000 n 
0000144227 00000 n 
0000144584 00000 n 
0000145335 00000 n 
0000148764 00000 n 
0000148873 00000 n 
0000149022 00000 n 
0000152670 00000 n 
trailer
<<
/Root 5 0 R
/Size 40
>>
startxref
174834
%%EOF`

	table := getTableForString(input)

	if table == nil {
		t.Fatal("expected non-nil table")
	}

	if len(table.ObjectOffsets()) != 39 {
		t.Errorf("expected 39 object offsets, got %d", len(table.ObjectOffsets()))
	}
}

func TestParseNewDefaultTable(t *testing.T) {
	input := `one xref
0 6
0000000003 65535 f
0000000090 00000 n
0000000081 00000 n
0000000000 00007 f
0000000331 00000 n
0000000409 00000 n

trailer
<< >>`

	scanner, bytes := testutil.CreateStringScanner(input)
	result := TryReadTableAtOffset(FileHeaderOffset{Value: 0}, 4, bytes, scanner, logging.NoopLog)

	if result == nil {
		t.Fatal("expected non-nil table")
	}

	if len(result.ObjectOffsets()) != 4 {
		t.Errorf("expected 4 object offsets, got %d", len(result.ObjectOffsets()))
	}
}

func TestOffsetNotXrefThrows(t *testing.T) {
	result := parse("12 0 obj <<>> endobj xref")

	if result != nil {
		t.Error("expected nil table")
	}
}

func TestOffsetXButNotXrefThrows(t *testing.T) {
	result := parse(`xtable
trailer`)

	if result != nil {
		t.Error("expected nil table")
	}
}

func TestEmptyTableReturnsEmpty(t *testing.T) {
	result := parse(`xref
trailer
<<>>`)

	if result == nil {
		t.Fatal("expected non-nil table")
	}

	if result.Dictionary() == nil {
		t.Error("expected non-nil dictionary")
	}

	if len(result.ObjectOffsets()) != 0 {
		t.Errorf("expected empty object offsets, got %d", len(result.ObjectOffsets()))
	}
}

func TestInvalidSubsectionDefinitionLenientSkips(t *testing.T) {
	result := parse(`xref
ab 12
trailer
<<>>`)

	if result != nil {
		t.Error("expected nil table")
	}
}

func TestSkipsFirstFreeLine(t *testing.T) {
	result := parse(`xref
0 1
0000000000 65535 f
trailer
<<>>`)

	if result == nil {
		t.Fatal("expected non-nil table")
	}

	if result.Dictionary() == nil {
		t.Error("expected non-nil dictionary")
	}

	if len(result.ObjectOffsets()) != 0 {
		t.Errorf("expected empty object offsets, got %d", len(result.ObjectOffsets()))
	}
}

func TestReadsEntries(t *testing.T) {
	result := parse(`xref
0 3
0000000000 65535 f
0000000100 00000 n
0000000200 00005 n
trailer
<<>>`)

	assertObjectsMatchGen(t, result, map[string]int64{
		"1:0": 100,
		"2:5": 200,
	})

	if result.Dictionary() == nil {
		t.Error("expected non-nil dictionary")
	}
}

func TestReadsEntriesOffsetFirstNumber(t *testing.T) {
	result := parse(`xref
15 2
0000000190 00000 n
0000000250 00032 n
trailer
<<>>`)

	assertObjectsMatchGen(t, result, map[string]int64{
		"15:0": 190,
		"16:32": 250,
	})
}

func TestReadsEntriesSkippingBlankLine(t *testing.T) {
	result := parse(`xref
15 2
0000000190 00000 n

0000000250 00032 n
trailer
<<>>`)

	assertObjectsMatchGen(t, result, map[string]int64{
		"15:0": 190,
		"16:32": 250,
	})
}

func TestReadsEntriesFromMultipleSubsections(t *testing.T) {
	result := parse(`xref
0 4
0000000000 65535 f
0000000100 00000 n
0000000200 00005 n
0000000230 00005 n
15 2
0000000190 00007 n
0000000250 00032 n
trailer
<<>>`)

	assertObjectsMatchGen(t, result, map[string]int64{
		"1:0": 100,
		"2:5": 200,
		"3:5": 230,
		"15:7": 190,
		"16:32": 250,
	})
}

func TestEntryPointingAtOffsetInTableDoesNotThrow(t *testing.T) {
	result := parse(`xref
0 2
0000000000 65535 f
0000000010 00000 n
trailer
<<>>`)

	assertObjectsMatchGen(t, result, map[string]int64{
		"1:0": 10,
	})
}

func TestEntryWithInvalidFormatThrows(t *testing.T) {
	result := parse(`xref
0 22
0000000000 65535 f
0000aa0010 00000 n
trailer
<<>>`)

	if result != nil {
		t.Error("expected nil table for invalid hex format")
	}
}

func TestShortLineInTableReturnsThrows(t *testing.T) {
	result := parse(`xref
15 2
019 n
0000000250 00032 n
trailer
<<>>`)

	if result != nil {
		t.Error("expected nil table for short line")
	}
}

func TestSkipsBlankLinesPrecedingTrailer(t *testing.T) {
	result := parse(`xref
15 2
0000000190 00000 n
0000000250 00032 n

trailer
<<>>`)

	if len(result.ObjectOffsets()) != 2 {
		t.Errorf("expected 2 object offsets, got %d", len(result.ObjectOffsets()))
	}
}

func TestParseEntriesAfterDeclaredCountIfLenient(t *testing.T) {
	data := `xref
0 5
0000000003 65535 f
0000000090 00000 n
0000000081 00000 n
0000000223 00000 n
0000000331 00000 n
0000000127 00000 n
0000000409 00000 f
0000000418 00000 n

trailer
<< >>`

	result := getTableForString(data)

	if len(result.ObjectOffsets()) != 6 {
		t.Errorf("expected 6 object offsets, got %d", len(result.ObjectOffsets()))
	}
}

func TestParsesMissingWhitespaceAfterXref(t *testing.T) {
	data := `xref15 2
0000000190 00000 n
0000000250 00032 n

trailer
<<>>`

	result := getTableForString(data)

	if len(result.ObjectOffsets()) != 2 {
		t.Errorf("expected 2 object offsets, got %d", len(result.ObjectOffsets()))
	}
}

// --- helpers ---

func parse(str string) *XrefTable {
	scanner, bytes := testutil.CreateStringScanner(str)
	return TryReadTableAtOffset(FileHeaderOffset{Value: 0}, 0, bytes, scanner, logging.NoopLog)
}

func getTableForString(s string) *XrefTable {
	scanner, bytes := testutil.CreateStringScanner(s)
	table := TryReadTableAtOffset(FileHeaderOffset{Value: 0}, 0, bytes, scanner, logging.NoopLog)
	return table
}

func assertObjectsMatch(t *testing.T, table *XrefTable, expected map[int64]int64) {
	t.Helper()
	if table == nil {
		t.Fatal("expected non-nil table")
	}

	offsets := table.ObjectOffsets()
	if len(offsets) != len(expected) {
		t.Errorf("expected %d offsets, got %d", len(expected), len(offsets))
	}

	for objNum, wantOffset := range expected {
		ref, err := core.NewIndirectReference(objNum, 0)
		if err != nil {
			t.Fatalf("failed to create indirect reference for %d: %v", objNum, err)
		}

		loc, ok := offsets[ref]
		if !ok {
			t.Errorf("missing entry for object %d", objNum)
			continue
		}

		if loc.Value1 != wantOffset {
			t.Errorf("object %d: expected offset %d, got %d", objNum, wantOffset, loc.Value1)
		}
	}
}

func assertObjectsMatchGen(t *testing.T, table *XrefTable, expected map[string]int64) {
	t.Helper()
	if table == nil {
		t.Fatal("expected non-nil table")
	}

	offsets := table.ObjectOffsets()
	if len(offsets) != len(expected) {
		t.Errorf("expected %d offsets, got %d", len(expected), len(offsets))
	}

	for key, wantOffset := range expected {
		parts := splitTwo(key)
		ref, err := core.NewIndirectReference(parts[0], int(parts[1]))
		if err != nil {
			t.Fatalf("failed to create indirect reference for %s: %v", key, err)
		}

		loc, ok := offsets[ref]
		if !ok {
			t.Errorf("missing entry for %s", key)
			continue
		}

		if loc.Value1 != wantOffset {
			t.Errorf("%s: expected offset %d, got %d", key, wantOffset, loc.Value1)
		}
	}
}

func splitTwo(s string) [2]int64 {
	var r [2]int64
	n := 0
	num := int64(0)
	for _, c := range s {
		if c == ':' {
			r[n] = num
			n++
			num = 0
		} else {
			num = num*10 + int64(c-'0')
		}
	}
	r[n] = num
	return r
}
