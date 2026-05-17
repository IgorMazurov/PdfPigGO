package parser

import (
	"fmt"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/encodings"
	"github.com/uglytoad/pdfpig/go/fonts/type1"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

const (
	clearToMark       = "cleartomark"
	pfbFileIndicator  = 0x80
)

var encryptedPortionParser = NewType1EncryptedPortionParser()

// previousTokenSet tracks the last three tokens seen during scanning.
type previousTokenSet struct {
	tokens [3]tokens.Token
}

func (s *previousTokenSet) Get(index int) tokens.Token {
	return s.tokens[2-index]
}

func (s *previousTokenSet) Add(token tokens.Token) {
	s.tokens[0] = s.tokens[1]
	s.tokens[1] = s.tokens[2]
	s.tokens[2] = token
}

// Parse parses an embedded Adobe Type 1 font file.
func Parse(inputBytes core.InputBytes, length1, length2 int, stackDepthGuard *core.StackDepthGuard) (*Type1Font, error) {
	isEntirePfbFile := false
	peeked, hasPeek := inputBytes.Peek()
	if hasPeek && peeked == pfbFileIndicator {
		isEntirePfbFile = true
	}

	var eexecPortion []byte

	if isEntirePfbFile {
		ascii, binary, err := readPfbHeader(inputBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to read PFB header: %w", err)
		}
		eexecPortion = binary
		inputBytes = core.NewMemoryInputBytes(ascii)
	}

	scanner := tokenization.NewCoreTokenScanner(inputBytes, false, stackDepthGuard, tokenization.ScannerScopeNone, nil, false, false)

	if !scanner.Advance() {
		return nil, fmt.Errorf("the Type1 program did not start with '%%!'")
	}

	comment, ok := scanner.Current().(*tokens.CommentToken)
	if !ok || !strings.HasPrefix(comment.Data(), "!") {
		return nil, fmt.Errorf("the Type1 program did not start with '%%!'")
	}

	name := "Unknown"
	parts := strings.Fields(comment.Data())
	if len(parts) >= 2 {
		name = parts[1]
	}

	for scanner.Advance() {
		if _, ok := scanner.Current().(*tokens.CommentToken); !ok {
			break
		}
	}

	dictionaries := make([]*tokens.DictionaryToken, 0)

	arrayTokenizer := type1.NewType1ArrayTokenizer()
	nameTokenizer := type1.NewType1NameTokenizer()
	scanner.RegisterCustomTokenizer('{', arrayTokenizer)
	scanner.RegisterCustomTokenizer('/', nameTokenizer)

	defer func() {
		scanner.DeregisterCustomTokenizer(arrayTokenizer)
		scanner.DeregisterCustomTokenizer(nameTokenizer)
	}()

	tempEexecPortion := core.NewArrayPoolBufferWriter()
	defer tempEexecPortion.Dispose()

	tokenSet := &previousTokenSet{}
	tokenSet.Add(scanner.Current())

	for scanner.Advance() {
		opToken, isOp := scanner.Current().(*tokens.OperatorToken)
		if isOp {
			if opToken == tokens.Eexec {
				eexecBytes := readEexecPortion(inputBytes, tempEexecPortion)
				if !isEntirePfbFile {
					eexecPortion = eexecBytes
				}
			} else {
				handleOperator(opToken, scanner, tokenSet, &dictionaries)
			}
		}

		tokenSet.Add(scanner.Current())
	}

	encoding := getEncoding(dictionaries)
	matrix := getFontMatrix(dictionaries)
	boundingBox := getBoundingBox(dictionaries)

	if boundingBox == nil {
		boundingBox = &core.PdfRectangle{}
	}

	parseResult, err := encryptedPortionParser.Parse(eexecPortion, false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse encrypted portion: %w", err)
	}

	return NewType1Font(name, encoding, matrix, *boundingBox, parseResult.PrivateDictionary, parseResult.CharStrings), nil
}

// readEexecPortion reads bytes from eexec operator until cleartomark is found.
func readEexecPortion(inputBytes core.InputBytes, writer *core.ArrayPoolBufferWriter) []byte {
	offset := 0

	for inputBytes.MoveNext() {
		currentByte := inputBytes.CurrentByte()

		if currentByte == clearToMark[offset] {
			offset++
		} else {
			if offset > 0 {
				for i := 0; i < offset; i++ {
					writer.WriteSingle(clearToMark[i])
				}
			}
			offset = 0
		}

		if offset == len(clearToMark) {
			break
		}

		if offset > 0 {
			continue
		}

		writer.WriteSingle(currentByte)
	}

	return writer.WrittenSpan()
}

// readPfbHeader reads the PFB file header and returns the ASCII and binary portions.
func readPfbHeader(bytes core.InputBytes) (ascii []byte, binary []byte, err error) {
	readSize := func(recordType byte) (int, error) {
		if !bytes.MoveNext() {
			return 0, fmt.Errorf("unexpected end of input reading PFB header")
		}

		if bytes.CurrentByte() != pfbFileIndicator {
			return 0, fmt.Errorf("file does not start with 0x80, which indicates a full PFB file. Instead got: %d", bytes.CurrentByte())
		}

		if !bytes.MoveNext() {
			return 0, fmt.Errorf("unexpected end of input reading PFB record type")
		}

		if bytes.CurrentByte() != recordType {
			return 0, fmt.Errorf("encountered unexpected header type in the PFB file: %d", bytes.CurrentByte())
		}

		size := 0

		if !bytes.MoveNext() {
			return 0, fmt.Errorf("unexpected end of input reading PFB size")
		}
		size = int(bytes.CurrentByte())

		if !bytes.MoveNext() {
			return 0, fmt.Errorf("unexpected end of input reading PFB size")
		}
		size += int(bytes.CurrentByte()) << 8

		if !bytes.MoveNext() {
			return 0, fmt.Errorf("unexpected end of input reading PFB size")
		}
		size += int(bytes.CurrentByte()) << 16

		if !bytes.MoveNext() {
			return 0, fmt.Errorf("unexpected end of input reading PFB size")
		}
		size += int(bytes.CurrentByte()) << 24

		return size, nil
	}

	asciiSize, err := readSize(0x01)
	if err != nil {
		return nil, nil, err
	}

	asciiPart := make([]byte, asciiSize)
	for i := 0; i < asciiSize; i++ {
		if !bytes.MoveNext() {
			return nil, nil, fmt.Errorf("unexpected end of input reading ASCII portion of PFB file")
		}
		asciiPart[i] = bytes.CurrentByte()
	}

	binarySize, err := readSize(0x02)
	if err != nil {
		return nil, nil, err
	}

	binaryPart := make([]byte, binarySize)
	for i := 0; i < binarySize; i++ {
		if !bytes.MoveNext() {
			return nil, nil, fmt.Errorf("unexpected end of input reading binary portion of PFB file")
		}
		binaryPart[i] = bytes.CurrentByte()
	}

	return asciiPart, binaryPart, nil
}

func handleOperator(token *tokens.OperatorToken, scanner tokenization.TokenScanner, set *previousTokenSet, dictionaries *[]*tokens.DictionaryToken) {
	if token.Data() == "dict" {
		numToken, ok := set.Get(0).(*tokens.NumericToken)
		if !ok {
			return
		}
		dict := readDictionary(numToken.IntVal(), scanner.(tokenization.SeekableTokenScanner))
		*dictionaries = append(*dictionaries, dict)
	}
}

func readDictionary(keys int, scanner tokenization.SeekableTokenScanner) *tokens.DictionaryToken {
	var previousToken tokens.Token

	dictionary := make(map[*tokens.NameToken]tokens.Token)

	for scanner.Advance() {
		opToken, isOp := scanner.Current().(*tokens.OperatorToken)
		if isOp && opToken.Data() == "begin" {
			break
		}
	}

	for i := 0; i < keys; i++ {
		scanner.Advance()
		key, ok := scanner.Current().(*tokens.NameToken)
		if !ok {
			return createDictionary(dictionary)
		}

		if key == tokens.Encoding {
			enc := readEncoding(scanner)
			if enc.encoding != nil {
				dictionary[key] = enc.encoding
			} else if enc.name != nil {
				dictionary[key] = enc.name
			}
			continue
		}

		for scanner.Advance() {
			current := scanner.Current()

			if current == tokens.Def {
				dictionary[key] = previousToken
				break
			}

			if current == tokens.Dict {
				num, ok := previousToken.(*tokens.NumericToken)
				if !ok {
					return createDictionary(dictionary)
				}
				inner := readDictionary(num.IntVal(), scanner)
				previousToken = inner
				continue
			}

			if current == tokens.Readonly {
				continue
			}

			if op, ok := current.(*tokens.OperatorToken); ok && op.Data() == "end" {
				continue
			}

			previousToken = current
		}
	}

	return createDictionary(dictionary)
}

func createDictionary(data map[*tokens.NameToken]tokens.Token) *tokens.DictionaryToken {
	if len(data) == 0 {
		empty := make(map[string]tokens.Token)
		dt, _ := tokens.WithMap(empty)
		return dt
	}
	dt, err := tokens.NewDictionary(data)
	if err != nil {
		empty := make(map[string]tokens.Token)
		fallback, _ := tokens.WithMap(empty)
		return fallback
	}
	return dt
}

type encodingResult struct {
	encoding *tokens.ArrayToken
	name     *tokens.NameToken
}

func readEncoding(scanner tokenization.SeekableTokenScanner) encodingResult {
	result := make([]tokens.Token, 0)

	scanner.Advance()
	if _, ok := scanner.Current().(*tokens.NumericToken); !ok {
		if op, ok := scanner.Current().(*tokens.OperatorToken); ok && op.Data() == tokens.StandardEncoding.Data() {
			return encodingResult{nil, tokens.StandardEncoding}
		}
		at := tokens.NewArrayToken(result)
		return encodingResult{at, nil}
	}

	scanner.Advance()
	arrayOp, ok := scanner.Current().(*tokens.OperatorToken)
	if !ok || arrayOp.Data() != "array" {
		at := tokens.NewArrayToken(result)
		return encodingResult{at, nil}
	}

	stoppedOnDup := false

	for scanner.Advance() {
		opToken, isOp := scanner.Current().(*tokens.OperatorToken)
		if isOp && opToken.Data() == "for" {
			break
		}

		if !isOp {
			continue
		}

		if strings.EqualFold(opToken.Data(), "for") {
			break
		}

		if strings.EqualFold(opToken.Data(), tokens.Dup.Data()) {
			stoppedOnDup = true
			break
		}
	}

	currentIsFor := false
	if op, ok := scanner.Current().(*tokens.OperatorToken); ok && op == tokens.For {
		currentIsFor = true
	}

	if !currentIsFor && !stoppedOnDup {
		at := tokens.NewArrayToken(result)
		return encodingResult{at, nil}
	}

	isDefOrReadonly := func() bool {
		return scanner.Current() == tokens.Def || scanner.Current() == tokens.Readonly
	}

	for stoppedOnDup || scanner.Advance() {
		if isDefOrReadonly() {
			break
		}
		stoppedOnDup = false

		if scanner.Current() != tokens.Dup {
			panic(fmt.Errorf("expected the array for encoding to begin with 'dup'"))
		}

		scanner.Advance()
		number := scanner.Current().(*tokens.NumericToken)
		scanner.Advance()
		name := scanner.Current().(*tokens.NameToken)

		scanner.Advance()
		put, ok := scanner.Current().(*tokens.OperatorToken)
		if !ok || put != tokens.Put {
			panic(fmt.Errorf("expected the array entry to end with 'put'"))
		}

		result = append(result, number, name)
	}

	for scanner.Current() != tokens.Def {
		if !scanner.Advance() {
			break
		}
	}

	at := tokens.NewArrayToken(result)
	return encodingResult{at, nil}
}

func getEncoding(dictionaries []*tokens.DictionaryToken) map[int]string {
	result := make(map[int]string)

	for _, dict := range dictionaries {
		token, found := dict.TryGet(tokens.Encoding)
		if !found {
			continue
		}

		if encArray, ok := token.(*tokens.ArrayToken); ok {
			data := encArray.Data()
			for i := 0; i < len(data); i += 2 {
				code, ok1 := data[i].(*tokens.NumericToken)
				name, ok2 := data[i+1].(*tokens.NameToken)
				if ok1 && ok2 {
					result[code.IntVal()] = name.Data()
				}
			}
			return result
		}

		if encName, ok := token.(*tokens.NameToken); ok && encName.Equals(tokens.StandardEncoding) {
			return encodings.StandardEncodingValue.CodeToNameMap()
		}
	}

	return result
}

func getFontMatrix(dictionaries []*tokens.DictionaryToken) *tokens.ArrayToken {
	for _, dict := range dictionaries {
		token, found := dict.TryGet(tokens.FontMatrix)
		if found {
			if array, ok := token.(*tokens.ArrayToken); ok {
				return array
			}
		}
	}
	return nil
}

func getBoundingBox(dictionaries []*tokens.DictionaryToken) *core.PdfRectangle {
	for _, dict := range dictionaries {
		token, found := dict.TryGet(tokens.FontBbox)
		if !found {
			continue
		}

		array, ok := token.(*tokens.ArrayToken)
		if !ok {
			continue
		}

		data := array.Data()
		if len(data) != 4 {
			continue
		}

		x1, ok1 := data[0].(*tokens.NumericToken)
		y1, ok2 := data[1].(*tokens.NumericToken)
		x2, ok3 := data[2].(*tokens.NumericToken)
		y2, ok4 := data[3].(*tokens.NumericToken)

		if ok1 && ok2 && ok3 && ok4 {
			rect := core.NewPdfRectangleFloat(x1.DoubleVal(), y1.DoubleVal(), x2.DoubleVal(), y2.DoubleVal())
			return &rect
		}
	}

	return nil
}
