package filestructure_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/parser/filestructure"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/testutil"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestFindsTwoXrefs(t *testing.T) {
	content := `%PDF-1.7
%âãÏÓ
5 0 obj
<</Filter/FlateDecode/Length 66>>stream
abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz
endstream
endobj
xref7 1
0000000000 65535 f 
0000000500 00000 n 
4 0 obj
<</Contents 5 0 R/MediaBox[0 0 595 842]/Parent 2 0 R/Resources<</Font<</F1 6 0 R>>>>/TrimBox[0 0 595 842]/Type/Page>>
endobj
xref
0 3
0000000000 65535 f 
0000000443 00000 n 
0000000576 00000 n
trailer
<< /Size 100 /Root 100 >>
startxref
9000
%%EOF`

	ib := testutil.ConvertStringToBytes(content, false)

	guard, _ := core.NewStackDepthGuard(256)
	scanner := tokenization.NewCoreTokenScanner(ib.Bytes, true, guard, tokenization.ScannerScopeNone, nil, false, false)

	results, err := filestructure.FirstPassParse(
		filestructure.NewFileHeaderOffset(0),
		ib.Bytes,
		scanner,
		logging.NoopLog,
	)
	if err != nil {
		t.Fatalf("FirstPassParse returned error: %v", err)
	}

	if len(results.Parts) != 2 {
		t.Errorf("expected 2 parts, got %d", len(results.Parts))
	}

	if results.Trailer == nil {
		t.Error("expected Trailer to not be nil")
	}

	ref, _ := core.NewIndirectReference(8, 0)
	loc, ok := results.XrefOffsets[ref]
	if !ok {
		t.Fatalf("XrefOffsets missing key for IndirectReference(8, 0)")
	}
	if loc.Value1 != 500 {
		t.Errorf("expected Value1 == 500, got %d", loc.Value1)
	}
}

func TestFindsXrefsFromRealFileTruncated(t *testing.T) {
	content := `%PDF-1.5
%âãÏÓ
5 0 obj <</Linearized 1/L 4631/O 8/E 1125/N 1/T 4485/H [ 436 129]>>
endobj
xref
5 7
0000000016 00000 n
0000000565 00000 n
0000000436 00000 n
0000000639 00000 n
0000000796 00000 n
0000001001 00000 n
0000001098 00000 n
trailer
<</Size 12/Prev 4475/Root 6 0 R/Info 4 0 R/ID[<2c9c3edf9641f1459e947e7f933f6da0><2c9c3edf9641f1459e947e7f933f6da0>]>>
startxref
0
%%EOF
7 0 obj<</Length 52/Filter/FlateDecode/O 67/S 38>>stream
3842973893927327893237832738923732923782987348
endstream
endobj
6 0 obj<</Pages 2 0 R/Outlines 1 0 R/Type/Catalog/Metadata 3 0 R>>
endobj
8 0 obj<</Contents 9 0 R/Type/Page/Parent 2 0 R/Rotate 0/MediaBox[0 0 612 792]/CropBox[0 0 612 792]/Resources<</Font<</F1 10 0 R>>/ProcSet 11 0 R>>>>
endobj
9 0 obj<</Length 137/Filter/FlateDecode>>stream
abajsgiwgbkeeuuehxh9x3oihx2h802chc280h2082x
endstream
endobj
10 0 obj<</Type/Font/Name/F1/Encoding/MacRomanEncoding/BaseFont/Helvetica/Subtype/Type1>>
endobj
11 0 obj[/PDF/Text]
endobj
1 0 obj<</Count 0/Type/Outlines>>
endobj
2 0 obj<</Count 1/Kids[8 0 R]/Type/Pages>>
endobj
4 0 obj<</ModDate(D:20070213222810-05'00')/CreationDate(D:20070213222810-05'00')>>
endobj
xref
0 5
0000000000 65535 f
0000001125 00000 n
0000001166 00000 n
0000001216 00000 n
0000004385 00000 n
trailer
<</Size 5>>
startxref
116
%%EOF`

	content = normalizeLineEndings(content)

	ib := testutil.ConvertStringToBytes(content, false)

	guard, _ := core.NewStackDepthGuard(256)
	scanner := tokenization.NewCoreTokenScanner(ib.Bytes, true, guard, tokenization.ScannerScopeNone, nil, false, false)

	results, err := filestructure.FirstPassParse(
		filestructure.NewFileHeaderOffset(0),
		ib.Bytes,
		scanner,
		logging.NoopLog,
	)
	if err != nil {
		t.Fatalf("FirstPassParse returned error: %v", err)
	}

	offsets := make([]int64, 0, len(results.Parts))
	for _, part := range results.Parts {
		offsets = append(offsets, part.Offset())
	}
	sort.Slice(offsets, func(i, j int) bool {
		return offsets[i] < offsets[j]
	})

	if len(offsets) == 0 || (len(offsets) >= 1 && offsets[0] != 98) {
		t.Errorf("expected first offset to be 98, got %v", offsets)
	}
	if len(offsets) < 2 || offsets[1] != 1186 {
		t.Errorf("expected second offset to be 1186, got %v", offsets)
	}

	if results.Trailer == nil {
		t.Error("expected Trailer to not be nil")
	}

	ib.Bytes.Seek(98, 0)
	scanner2 := tokenization.NewCoreTokenScanner(ib.Bytes, false, guard, tokenization.ScannerScopeNone, nil, false, false)
	scanner2.Advance()

	opToken, ok := scanner2.Current().(*tokens.OperatorToken)
	if !ok {
		t.Fatalf("expected OperatorToken at offset 98, got %T", scanner2.Current())
	}
	if opToken.Data() != "xref" {
		t.Errorf("expected token data \"xref\", got %q", opToken.Data())
	}
}

func normalizeLineEndings(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\n", "\r\n")
}
