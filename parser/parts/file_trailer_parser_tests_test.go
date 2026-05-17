package parts_test

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/parser/filestructure"
	"github.com/uglytoad/pdfpig/go/testutil"
	"github.com/uglytoad/pdfpig/go/tokenization"
)

func TestFindsCompliantStartXref(t *testing.T) {
	input := testutil.ConvertStringToBytes(`sta455%r endstream
endobj

12 0 obj
1234  %eof
endobj

startxref
    456

%%EOF`, false)

	guard, _ := core.NewStackDepthGuard(256)
	scanner := tokenization.NewCoreTokenScanner(input.Bytes, true, guard, tokenization.ScannerScopeNone, nil, false, false)

	result := filestructure.GetFirstCrossReferenceOffset(input.Bytes, scanner, logging.NoopLog)

	if result.StartXRefDeclaredOffset == nil {
		t.Fatal("expected StartXRefDeclaredOffset to not be nil")
	}
	if *result.StartXRefDeclaredOffset != 456 {
		t.Errorf("expected StartXRefDeclaredOffset == 456, got %d", *result.StartXRefDeclaredOffset)
	}
}

func TestIncludesStartXrefFollowingEndOfFile(t *testing.T) {
	input := testutil.ConvertStringToBytes(`11 0 obj
<< /Type/Something /W[12 0 5 6] >>
endobj

12 0 obj
1234  %eof
endobj

startxref
    1384733

%%EOF

% I decided to put some nonsense here:
% because I could hahaha
startxref
17`, false)

	guard, _ := core.NewStackDepthGuard(256)
	scanner := tokenization.NewCoreTokenScanner(input.Bytes, true, guard, tokenization.ScannerScopeNone, nil, false, false)

	result := filestructure.GetFirstCrossReferenceOffset(input.Bytes, scanner, logging.NoopLog)

	if result.StartXRefDeclaredOffset == nil {
		t.Fatal("expected StartXRefDeclaredOffset to not be nil")
	}
	if *result.StartXRefDeclaredOffset != 17 {
		t.Errorf("expected StartXRefDeclaredOffset == 17, got %d", *result.StartXRefDeclaredOffset)
	}
}

func TestMissingStartXrefThrows(t *testing.T) {
	input := testutil.ConvertStringToBytes(`11 0 obj
<< /Type/Something /W[12 0 5 6] >>
endobj

12 0 obj
1234  %eof
endobj

startref
    1384733

%%EOF

% I decided to put some nonsense here:
% because I could hahaha
start_rexf
17`, false)

	guard, _ := core.NewStackDepthGuard(256)
	scanner := tokenization.NewCoreTokenScanner(input.Bytes, true, guard, tokenization.ScannerScopeNone, nil, false, false)

	result := filestructure.GetFirstCrossReferenceOffset(input.Bytes, scanner, logging.NoopLog)

	if result.StartXRefDeclaredOffset == nil {
		t.Fatal("expected StartXRefDeclaredOffset to not be nil")
	}
	if *result.StartXRefDeclaredOffset != 1384733 {
		t.Errorf("expected StartXRefDeclaredOffset == 1384733, got %d", *result.StartXRefDeclaredOffset)
	}
}

func TestBadInputBytesReturnsNull(t *testing.T) {
	input := testutil.ConvertStringToBytes("11 0 obj", false)

	guard, _ := core.NewStackDepthGuard(256)
	scanner := tokenization.NewCoreTokenScanner(input.Bytes, true, guard, tokenization.ScannerScopeNone, nil, false, false)

	result := filestructure.GetFirstCrossReferenceOffset(input.Bytes, scanner, logging.NoopLog)

	if result.StartXRefDeclaredOffset != nil {
		t.Errorf("expected StartXRefDeclaredOffset to be nil, got %d", *result.StartXRefDeclaredOffset)
	}
	if result.StartXRefOperatorToken != nil {
		t.Errorf("expected StartXRefOperatorToken to be nil, got %d", *result.StartXRefOperatorToken)
	}
}

func TestInvalidTokensAfterStartXrefReturnsNull(t *testing.T) {
	input := testutil.ConvertStringToBytes(`11 0 obj
        << /Type/Font >>
endobj

startxref 
<< /Why (am i here?) >> 69
%EOF`, false)

	guard, _ := core.NewStackDepthGuard(256)
	scanner := tokenization.NewCoreTokenScanner(input.Bytes, true, guard, tokenization.ScannerScopeNone, nil, false, false)

	result := filestructure.GetFirstCrossReferenceOffset(input.Bytes, scanner, logging.NoopLog)

	if result.StartXRefDeclaredOffset != nil {
		t.Errorf("expected StartXRefDeclaredOffset to be nil, got %d", *result.StartXRefDeclaredOffset)
	}
	if result.StartXRefOperatorToken == nil {
		t.Error("expected StartXRefOperatorToken to not be nil")
	}
}

func TestMissingNumericAfterStartXrefReturnsNull(t *testing.T) {
	input := testutil.ConvertStringToBytes(`1 0 obj
<< /Type/Font >>
endobj

startxref 
`, false)

	guard, _ := core.NewStackDepthGuard(256)
	scanner := tokenization.NewCoreTokenScanner(input.Bytes, true, guard, tokenization.ScannerScopeNone, nil, false, false)

	result := filestructure.GetFirstCrossReferenceOffset(input.Bytes, scanner, logging.NoopLog)

	if result.StartXRefDeclaredOffset != nil {
		t.Errorf("expected StartXRefDeclaredOffset to be nil, got %d", *result.StartXRefDeclaredOffset)
	}
	if result.StartXRefOperatorToken == nil {
		t.Error("expected StartXRefOperatorToken to not be nil")
	}
}

func TestTakesLastStartXrefPrecedingEndOfFile(t *testing.T) {
	input := testutil.ConvertStringToBytes(`11 0 obj
<< /Type/Something /W[12 0 5 6] >>
endobj

12 0 obj
1234  %eof
endobj

startxref
    1384733

%actually I changed my mind

startxref
         1274665676543

%%EOF`, false)

	guard, _ := core.NewStackDepthGuard(256)
	scanner := tokenization.NewCoreTokenScanner(input.Bytes, true, guard, tokenization.ScannerScopeNone, nil, false, false)

	result := filestructure.GetFirstCrossReferenceOffset(input.Bytes, scanner, logging.NoopLog)

	if result.StartXRefDeclaredOffset == nil {
		t.Fatal("expected StartXRefDeclaredOffset to not be nil")
	}
	if *result.StartXRefDeclaredOffset != 1274665676543 {
		t.Errorf("expected StartXRefDeclaredOffset == 1274665676543, got %d", *result.StartXRefDeclaredOffset)
	}
	if result.StartXRefOperatorToken == nil {
		t.Error("expected StartXRefOperatorToken to not be nil")
	}
}

func TestCanReadStartXrefIfCommentsPresent(t *testing.T) {
	input := testutil.ConvertStringToBytes(`
startxref %Commented here
    57695

%%EOF`, false)

	guard, _ := core.NewStackDepthGuard(256)
	scanner := tokenization.NewCoreTokenScanner(input.Bytes, true, guard, tokenization.ScannerScopeNone, nil, false, false)

	result := filestructure.GetFirstCrossReferenceOffset(input.Bytes, scanner, logging.NoopLog)

	if result.StartXRefDeclaredOffset == nil {
		t.Fatal("expected StartXRefDeclaredOffset to not be nil")
	}
	if *result.StartXRefDeclaredOffset != 57695 {
		t.Errorf("expected StartXRefDeclaredOffset == 57695, got %d", *result.StartXRefDeclaredOffset)
	}
	if result.StartXRefOperatorToken == nil {
		t.Error("expected StartXRefOperatorToken to not be nil")
	}
}
