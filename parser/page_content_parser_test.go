package parser_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/graphics"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/parser"
	"github.com/uglytoad/pdfpig/go/tokens"
)

var testParser parser.PageContentParser

func init() {
	guard, err := core.NewStackDepthGuard(256)
	if err != nil {
		panic(err)
	}
	testParser = parser.NewPageContentParser(
		graphics.GetReflectionFactory(),
		guard,
		false,
	)
	lenientGuard, _ := core.NewStackDepthGuard(256)
	testLenientParser = parser.NewPageContentParser(
		graphics.GetReflectionFactory(),
		lenientGuard,
		true,
	)
}

var testLenientParser parser.PageContentParser

func getParser(useLenientParsing bool) parser.PageContentParser {
	if useLenientParsing {
		return testLenientParser
	}
	return testParser
}

func lineEndingsToWhiteSpace(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return s
}

func TestCorrectlyWritesSmallTextContent(t *testing.T) {
	const s = `BT
/F13 48 Tf
20 38 Td
1 Tr
2 w
(ABC) Tj
ET`
	input := core.NewMemoryInputBytes([]byte(s))
	parsed := getParser(false).Parse(1, input, logging.NoopLog)

	var buf bytes.Buffer
	for _, op := range parsed {
		if wr, ok := any(op).(interface{ Write(w io.Writer) error }); ok {
			if err := wr.Write(&buf); err != nil {
				t.Fatalf("Write failed: %v", err)
			}
		}
	}

	text := strings.TrimSpace(lineEndingsToWhiteSpace(buf.String()))
	expected := lineEndingsToWhiteSpace(s)

	if text != expected {
		t.Errorf("Round-trip mismatch.\nExpected:\n%s\nGot:\n%s", expected, text)
	}
}

func TestCorrectlyExtractsOptionsInTextContext(t *testing.T) {
	const s = `BT
/F13 48 Tf
20 38 Td
1 Tr
2 w
(ABC) Tj
ET`
	input := core.NewMemoryInputBytes([]byte(s))
	result := getParser(false).Parse(1, input, logging.NoopLog)

	if len(result) != 7 {
		t.Fatalf("Expected 7 operations, got %d", len(result))
	}

	if _, ok := result[0].(graphics.BeginText); !ok {
		t.Errorf("result[0] should be BeginText, got %T", result[0])
	}

	font, ok := result[1].(*graphics.SetFontAndSize)
	if !ok {
		t.Fatalf("result[1] should be SetFontAndSize, got %T", result[1])
	}
	if font.Font == nil || font.Font.Data() != "F13" {
		t.Errorf("Expected font F13, got %v", font.Font)
	}
	if font.Size != 48 {
		t.Errorf("Expected size 48, got %g", font.Size)
	}

	nextLine, ok := result[2].(*graphics.MoveToNextLineWithOffset)
	if !ok {
		t.Fatalf("result[2] should be MoveToNextLineWithOffset, got %T", result[2])
	}
	if nextLine.Tx != 20 {
		t.Errorf("Expected Tx=20, got %g", nextLine.Tx)
	}
	if nextLine.Ty != 38 {
		t.Errorf("Expected Ty=38, got %g", nextLine.Ty)
	}

	renderingMode, ok := result[3].(*graphics.SetTextRenderingMode)
	if !ok {
		t.Fatalf("result[3] should be SetTextRenderingMode, got %T", result[3])
	}
	if renderingMode.Mode != core.StrokeText {
		t.Errorf("Expected StrokeText mode, got %v", renderingMode.Mode)
	}

	lineWidth, ok := result[4].(*graphics.SetLineWidth)
	if !ok {
		t.Fatalf("result[4] should be SetLineWidth, got %T", result[4])
	}
	if lineWidth.Width != 2 {
		t.Errorf("Expected width=2, got %g", lineWidth.Width)
	}

	textOp, ok := result[5].(*graphics.ShowText)
	if !ok {
		t.Fatalf("result[5] should be ShowText, got %T", result[5])
	}
	if textOp.Text != "ABC" {
		t.Errorf("Expected text 'ABC', got '%s'", textOp.Text)
	}

	if _, ok := result[6].(graphics.EndText); !ok {
		t.Errorf("result[6] should be EndText, got %T", result[6])
	}
}

func TestSkipsComments(t *testing.T) {
	const s = `BT
21 32 Td %A comment here
0 Tr
ET`

	input := core.NewMemoryInputBytes([]byte(s))
	result := getParser(false).Parse(1, input, logging.NoopLog)

	if len(result) != 4 {
		t.Fatalf("Expected 4 operations, got %d", len(result))
	}

	if _, ok := result[0].(graphics.BeginText); !ok {
		t.Errorf("result[0] should be BeginText, got %T", result[0])
	}

	moveLine, ok := result[1].(*graphics.MoveToNextLineWithOffset)
	if !ok {
		t.Fatalf("result[1] should be MoveToNextLineWithOffset, got %T", result[1])
	}
	if moveLine.Tx != 21 {
		t.Errorf("Expected Tx=21, got %g", moveLine.Tx)
	}
	if moveLine.Ty != 32 {
		t.Errorf("Expected Ty=32, got %g", moveLine.Ty)
	}

	renderingMode, ok := result[2].(*graphics.SetTextRenderingMode)
	if !ok {
		t.Fatalf("result[2] should be SetTextRenderingMode, got %T", result[2])
	}
	if renderingMode.Mode != core.FillText {
		t.Errorf("Expected FillText mode, got %v", renderingMode.Mode)
	}

	if _, ok := result[3].(graphics.EndText); !ok {
		t.Errorf("result[3] should be EndText, got %T", result[3])
	}
}

func TestHandlesEscapedLineBreaks(t *testing.T) {
	const s = `q 1 0 0 1 48 434
cm BT 0.0001 Tc 19 0 0 19 0 0 Tm /Tc1 1 Tf (   \(sleep 1; printf ""QUIT\\r\\n""\) | )
            Tj ET Q`

	input := core.NewMemoryInputBytes([]byte(s))
	result := getParser(false).Parse(1, input, logging.NoopLog)

	if len(result) != 9 {
		t.Fatalf("Expected 9 operations, got %d", len(result))
	}

	if _, ok := result[0].(graphics.Push); !ok {
		t.Errorf("result[0] should be Push, got %T", result[0])
	}
	if _, ok := result[1].(*graphics.ModifyCurrentTransformationMatrix); !ok {
		t.Errorf("result[1] should be ModifyCurrentTransformationMatrix, got %T", result[1])
	}
	if _, ok := result[2].(graphics.BeginText); !ok {
		t.Errorf("result[2] should be BeginText, got %T", result[2])
	}
	if _, ok := result[3].(*graphics.SetCharacterSpacing); !ok {
		t.Errorf("result[3] should be SetCharacterSpacing, got %T", result[3])
	}
	if _, ok := result[4].(*graphics.SetTextMatrix); !ok {
		t.Errorf("result[4] should be SetTextMatrix, got %T", result[4])
	}
	if _, ok := result[5].(*graphics.SetFontAndSize); !ok {
		t.Errorf("result[5] should be SetFontAndSize, got %T", result[5])
	}

	textOp, ok := result[6].(*graphics.ShowText)
	if !ok {
		t.Fatalf("result[6] should be ShowText, got %T", result[6])
	}
	expectedText := `   (sleep 1; printf ""QUIT\r\n"") | `
	if textOp.Text != expectedText {
		t.Errorf("Expected text '%s', got '%s'", expectedText, textOp.Text)
	}

	if _, ok := result[7].(graphics.EndText); !ok {
		t.Errorf("result[7] should be EndText, got %T", result[7])
	}
	if _, ok := result[8].(graphics.Pop); !ok {
		t.Errorf("result[8] should be Pop, got %T", result[8])
	}
}

func TestHandlesWeirdNumber(t *testing.T) {
	const s = `/Alpha1
gs
0
0
0
rg
0.00-90
151555.0
m
302399.97
151555.0
l`

	input := core.NewMemoryInputBytes([]byte(s))
	result := getParser(false).Parse(1, input, logging.NoopLog)

	if len(result) != 4 {
		t.Fatalf("Expected 4 operations, got %d", len(result))
	}
}

func TestHandlesIssue953_IntOverflowContent(t *testing.T) {
	const s = `BT
/TT6 1 Tf
12.007 0 0 12.007 163.2j
-0.19950 Tc
0 Tw
(x)Tj
-0.1949 1.4142 TD
(H)Tj
/TT7 1 Tf
12.031 0 0 12.031 157.38 85.2 Tm
<0077>Tj
-0.1945 1.4114 TD
<0077>Tj
/TT4 1 Tf
12.007 0 0 12.007 174.42 94.5601 Tm
0.0004 Tc
-0.0005 Tw
( + )Tj
E9 478l)]T862.68E9 478E9 484.54 9 155l)]T862.6av9 478E9 15.2(
ET
154.386( i92 m
171.6 97.62 l
S
BT
/TT6 28 Tf
12.03128 T2002.0307 163.2j
-0.19950 DAc
0 Tw853Tj
0.1945 1.4142 om)873j
-0.574142 om)68.80
-0.5797 0 TD
(f)Tj
/TT( )7Tf
0.31945 1.5341 TD371.4j
2.82
8.2652 0 5.724 TD
0 Tc
-0.0001 2748.3( = 091ity )-27483
[(te27483
[(te27483
[(te27483
[(te27483
[(te27483
[(Eq.)52   \(2.1
( 
`

	input := core.NewMemoryInputBytes([]byte(s))
	result := getParser(true).Parse(1, input, logging.NoopLog)

	if len(result) == 0 {
		t.Error("Expected non-empty result from lenient parser")
	}
}

func TestCorrectlyExtractsOperations(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "SimpleGoogleDocPageContent.txt"))
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	content := strings.TrimPrefix(string(data), "\ufeff")

	input := core.NewMemoryInputBytes([]byte(content))
	result := getParser(false).Parse(1, input, logging.NoopLog)

	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}
}

func TestCorrectlyWritesOperations(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "SimpleGoogleDocPageContent.txt"))
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	content := strings.TrimPrefix(string(data), "\ufeff")

	input := core.NewMemoryInputBytes([]byte(content))
	result := getParser(false).Parse(1, input, logging.NoopLog)

	var buf bytes.Buffer
	for _, op := range result {
		if wr, ok := any(op).(interface{ Write(w io.Writer) error }); ok {
			if err := wr.Write(&buf); err != nil {
				t.Fatalf("Write failed: %v", err)
			}
		}
	}

	text := lineEndingsToWhiteSpace(buf.String())
	expected := lineEndingsToWhiteSpace(content)

	replacementRegex := regexp.MustCompile(`\s(\.\d+)\b`)
	expected = replacementRegex.ReplaceAllString(expected, " 0$1")

	if text != expected {
		t.Errorf("Round-trip mismatch.\nExpected:\n%s\nGot:\n%s", expected, text)
	}
}

func TestCorrectlyHandlesFile0007511CorruptInlineImage(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "0007511-page-2.txt"))
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	input := core.NewMemoryInputBytes(data)
	result := getParser(true).Parse(1, input, logging.NoopLog)

	if len(result) == 0 {
		t.Error("Expected non-empty result from lenient parser")
	}
}

func TestParseReturnsOperations(t *testing.T) {
	const s = `BT
/F1 12 Tf
(Hello) Tj
ET`

	input := core.NewMemoryInputBytes([]byte(s))
	result := getParser(false).Parse(1, input, logging.NoopLog)

	if len(result) != 4 {
		t.Fatalf("Expected 4 operations, got %d", len(result))
	}

	for i, op := range result {
		if op == nil {
			t.Errorf("Operation at index %d is nil", i)
		}
	}
}

func TestTypeAssertionsOnParsedOperations(t *testing.T) {
	const s = `q
1 0 0 1 48 750 cm
Q`

	input := core.NewMemoryInputBytes([]byte(s))
	result := getParser(false).Parse(1, input, logging.NoopLog)

	if len(result) != 3 {
		t.Fatalf("Expected 3 operations, got %d", len(result))
	}

	pushOp, ok := result[0].(graphics.Push)
	if !ok {
		t.Errorf("result[0] should be Push, got %T", result[0])
	}

	cmOp, ok := result[1].(*graphics.ModifyCurrentTransformationMatrix)
	if !ok {
		t.Errorf("result[1] should be ModifyCurrentTransformationMatrix, got %T", result[1])
	} else {
		if cmOp.Value[0] != 1 || cmOp.Value[5] != 750 {
			t.Errorf("Unexpected matrix values: %+v", cmOp.Value)
		}
	}

	popOp, ok := result[2].(graphics.Pop)
	if !ok {
		t.Errorf("result[2] should be Pop, got %T", result[2])
	}

	_ = pushOp
	_ = popOp
}

func TestShowTextWithHex(t *testing.T) {
	const s = `BT
/F1 12 Tf
<48656C6C6F> Tj
ET`

	input := core.NewMemoryInputBytes([]byte(s))
	result := getParser(false).Parse(1, input, logging.NoopLog)

	if len(result) < 4 {
		t.Fatalf("Expected at least 4 operations, got %d", len(result))
	}

	showTextOp, ok := result[2].(*graphics.ShowText)
	if !ok {
		t.Errorf("result[2] should be ShowText, got %T", result[2])
	} else if showTextOp.Text != "" {
		t.Logf("ShowText has Text field: '%s'", showTextOp.Text)
	}
}

func TestGraphicsStateOperationInterface(t *testing.T) {
	const s = `BT
/F1 12 Tf
(Hello) Tj
ET`

	input := core.NewMemoryInputBytes([]byte(s))
	result := getParser(false).Parse(1, input, logging.NoopLog)

	for i, op := range result {
		var _ content.GraphicsStateOperation = op
		if op == nil {
			t.Errorf("Operation at index %d is nil", i)
		}
	}
}

func TestNameTokenInterning(t *testing.T) {
	const s = `/F1 12 Tf`

	input := core.NewMemoryInputBytes([]byte(s))
	result := getParser(false).Parse(1, input, logging.NoopLog)

	if len(result) != 1 {
		t.Fatalf("Expected 1 operation, got %d", len(result))
	}

	fontOp, ok := result[0].(*graphics.SetFontAndSize)
	if !ok {
		t.Fatalf("result[0] should be SetFontAndSize, got %T", result[0])
	}

	if fontOp.Font == nil {
		t.Fatal("Expected non-nil Font")
	}

	nameTok, ok := any(fontOp.Font).(tokens.Token)
	if !ok {
		t.Error("Font should implement Token interface")
	}

	_ = nameTok
}
