package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/type1"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings"
)

const (
	eexecEncryptionKey    = 55665
	eexecRandomBytes      = 4
	len4Bytes             = 4
	passwordConstant      = 5839
	charstringEncKey      = 4330
)

// Type1EncryptedPortionParser parses the encrypted portion of a Type 1 font program.
type Type1EncryptedPortionParser struct{}

// NewType1EncryptedPortionParser creates a new parser instance.
func NewType1EncryptedPortionParser() *Type1EncryptedPortionParser {
	return &Type1EncryptedPortionParser{}
}

// ParseResult holds the parsed private dictionary and charstrings from an encrypted Type 1 portion.
type ParseResult struct {
	PrivateDictionary *type1.Type1PrivateDictionary
	CharStrings       *charstrings.Type1CharStrings
}

// Parse decrypts and tokenizes the encrypted binary segment of a Type 1 font,
// extracting the Private dictionary entries and CharString definitions.
func (p *Type1EncryptedPortionParser) Parse(bytes []byte, isLenientParsing bool) (*ParseResult, error) {
	if !isBinary(bytes) {
		bytes = convertHexToBinary(bytes)
	}

	decrypted := decrypt(bytes, eexecEncryptionKey, eexecRandomBytes)

	if len(decrypted) == 0 {
		defaultPrivate := buildDefaultPrivate()
		defaultCharStrings := &charstrings.Type1CharStrings{}
		return &ParseResult{
			PrivateDictionary: defaultPrivate,
			CharStrings:       defaultCharStrings,
		}, nil
	}

	tokenizer, err := type1.NewType1Tokenizer(core.NewMemoryInputBytes(decrypted))
	if err != nil {
		return nil, fmt.Errorf("failed to create tokenizer: %w", err)
	}

	for {
		current := tokenizer.CurrentToken()
		if current == nil || current.IsPrivateDictionary() {
			break
		}
		tokenizer.GetNext()
	}

	if tokenizer.CurrentToken() == nil {
		return nil, fmt.Errorf("did not find the private dictionary start token")
	}

	next := tokenizer.GetNext()
	if next == nil || next.TokenType() != type1.TokenTypeInteger {
		return nil, fmt.Errorf("no length token was present in the stream following the private dictionary start, instead got %v", next)
	}

	length := next.(*type1.Type1Token).AsInt()

	readExpected(tokenizer, type1.TokenTypeNone, "dict")

	readExpectedAfterOptional(tokenizer, type1.TokenTypeNone, "dup", type1.TokenTypeNone, "begin")

	lenIv := len4Bytes
	builder := &type1.Type1PrivateDictionaryBuilder{}

	for i := 0; i < length; i++ {
		token := tokenizer.GetNext()

		if token == nil || token.TokenType() != type1.TokenTypeLiteral {
			break
		}

		key := token.(*type1.Type1Token).Text

		switch key {
		case type1.RdProcedure, type1.RdProcedureAlt:
			builder.Rd = readProcedure(tokenizer, false)
			readTillDef(tokenizer, false)

		case type1.NoAccessDef, type1.NoAccessDefAlt:
			builder.NoAccessDef = readProcedure(tokenizer, false)
			readTillDef(tokenizer, false)

		case type1.NoAccessPut, type1.NoAccessPutAlt:
			builder.NoAccessPut = readProcedure(tokenizer, false)
			readTillDef(tokenizer, false)

		case type1.BlueValues:
			builder.BlueValues = readArrayValues(tokenizer, func(t *type1.Type1Token) int { return t.AsInt() })

		case type1.OtherBlues:
			builder.OtherBlues = readArrayValues(tokenizer, func(t *type1.Type1Token) int { return t.AsInt() })

		case type1.StdHorizontalStemWidth:
			widths := readArrayValues(tokenizer, func(t *type1.Type1Token) float64 { return t.AsDouble() })
			if len(widths) > 0 {
				w := widths[0]
				builder.StandardHorizontalWidth = &w
			}

		case type1.StdVerticalStemWidth:
			widths := readArrayValues(tokenizer, func(t *type1.Type1Token) float64 { return t.AsDouble() })
			if len(widths) > 0 {
				w := widths[0]
				builder.StandardVerticalWidth = &w
			}

		case type1.StemSnapHorizontalWidths:
			builder.StemSnapHorizontalWidths = readArrayValues(tokenizer, func(t *type1.Type1Token) float64 { return t.AsDouble() })

		case type1.StemSnapVerticalWidths:
			builder.StemSnapVerticalWidths = readArrayValues(tokenizer, func(t *type1.Type1Token) float64 { return t.AsDouble() })

		case type1.BlueScale:
			v := readNumeric(tokenizer)
			builder.BlueScale = &v
			readTillDef(tokenizer, false)

		case type1.ForceBold:
			v := readBoolean(tokenizer)
			builder.ForceBold = &v
			readTillDef(tokenizer, false)

		case type1.SymMinFeature:
			startToken := tokenizer.GetNext()
			if startToken != nil && startToken.TokenType() == type1.TokenTypeStartArray {
				arrayTokens := readArrayValuesIntWithStart(tokenizer, true)
				if len(arrayTokens) >= 2 {
					builder.MinFeature = type1.NewMinFeature(arrayTokens[0], arrayTokens[1])
				}
			} else if startToken != nil && startToken.TokenType() == type1.TokenTypeStartProc {
				procTokens := readProcedure(tokenizer, true)
				if len(procTokens) >= 2 {
					t0 := procTokens[0].(*type1.Type1Token)
					t1 := procTokens[1].(*type1.Type1Token)
					builder.MinFeature = type1.NewMinFeature(t0.AsInt(), t1.AsInt())
				}
			}
			readTillDef(tokenizer, false)

		case type1.Password:
			pwVal := int(readNumeric(tokenizer))
			if pwVal != passwordConstant && !isLenientParsing {
				return nil, fmt.Errorf("type 1 font had the wrong password: %d", pwVal)
			}
			pwPtr := &pwVal
			builder.Password = pwPtr
			readTillDef(tokenizer, false)

		case type1.UniqueId:
			idVal := int(readNumeric(tokenizer))
			idPtr := &idVal
			builder.UniqueId = idPtr
			readTillDef(tokenizer, false)

		case type1.Len4:
			lenIv = int(readNumeric(tokenizer))
			readTillDef(tokenizer, false)

		case type1.BlueShift:
			v := int(readNumeric(tokenizer))
			builder.BlueShift = &v
			readTillDef(tokenizer, false)

		case type1.BlueFuzz:
			v := int(readNumeric(tokenizer))
			builder.BlueFuzz = &v
			readTillDef(tokenizer, false)

		case type1.FamilyBlues:
			builder.FamilyBlues = readArrayValues(tokenizer, func(t *type1.Type1Token) int { return t.AsInt() })

		case type1.FamilyOtherBlues:
			builder.FamilyOtherBlues = readArrayValues(tokenizer, func(t *type1.Type1Token) int { return t.AsInt() })

		case type1.LanguageGroup:
			v := int(readNumeric(tokenizer))
			builder.LanguageGroup = &v
			readTillDef(tokenizer, false)

		case type1.RndStemUp:
			v := readBoolean(tokenizer)
			builder.RoundStemUp = &v
			readTillDef(tokenizer, false)

		case type1.Subroutines:
			builder.Subroutines = readSubroutines(tokenizer, lenIv, isLenientParsing)

		case type1.OtherSubroutines:
			readOtherSubroutines(tokenizer, isLenientParsing)
			readTillDef(tokenizer, false)

		case type1.ExpansionFactor:
			v := readNumeric(tokenizer)
			builder.ExpansionFactor = &v
			readTillDef(tokenizer, false)

		case type1.Erode:
			readTillDef(tokenizer, true)

		default:
			readTillDef(tokenizer, true)
		}
	}

	currentToken := tokenizer.CurrentToken()

	var charStringList []*charstrings.Type1CharstringDecryptedBytes
	if currentToken != nil {
		for currentToken != nil && (currentToken.TokenType() != type1.TokenTypeLiteral || currentToken.(*type1.Type1Token).Text != "CharStrings") {
			currentToken = tokenizer.GetNext()
		}

		if currentToken != nil {
			charStringList = readCharStrings(tokenizer, lenIv, isLenientParsing)
		}
	}

	privateDict := builder.Build()

	parser := &charstrings.Type1CharStringParser{}
	instructions, err := parser.Parse(charStringList, builder.Subroutines)
	if err != nil {
		return nil, fmt.Errorf("failed to parse charstrings: %w", err)
	}

	return &ParseResult{
		PrivateDictionary: privateDict,
		CharStrings:       instructions,
	}, nil
}

func buildDefaultPrivate() *type1.Type1PrivateDictionary {
	baseBuilder := &fonts.AdobeStylePrivateDictionaryBuilder{}
	base := fonts.NewAdobeStylePrivateDictionary(baseBuilder)
	pw := passwordConstant
	mf := type1.DefaultMinFeature
	return &type1.Type1PrivateDictionary{
		AdobeStylePrivateDictionary: base,
		Password:                    pw,
		MinFeature:                  mf,
	}
}

// isBinary distinguishes between binary and hex encoded data.
func isBinary(bytes []byte) bool {
	if len(bytes) < 4 {
		return true
	}

	if core.IsWhitespace(bytes[0]) {
		return true
	}

	for i := 1; i < 4; i++ {
		b := bytes[i]
		if !core.IsHexByte(b) {
			return true
		}
	}

	return false
}

// convertHexToBinary converts a hex-encoded byte slice to binary.
func convertHexToBinary(bytes []byte) []byte {
	result := make([]byte, len(bytes)/2)
	index := 0

	var last rune = 0
	offset := 0
	for i := 0; i < len(bytes); i++ {
		c := rune(bytes[i])
		if !core.IsHexChar(c) {
			continue
		}

		if offset == 1 {
			result[index] = convertHexPair(last, c)
			index++
			offset = 0
		} else {
			offset++
		}

		last = c
	}

	return result[:index]
}

// decrypt applies the Type 1 eexec decryption algorithm.
func decrypt(bytes []byte, key int, randomBytes int) []byte {
	if randomBytes == -1 {
		return bytes
	}

	if randomBytes > len(bytes) || len(bytes) == 0 {
		return nil
	}

	const c1 = 52845
	const c2 = 22719

	plainBytes := make([]byte, len(bytes)-randomBytes)

	for i := 0; i < len(bytes); i++ {
		cipher := int(bytes[i]) & 0xFF
		plain := cipher ^ (key >> 8)

		if i >= randomBytes {
			plainBytes[i-randomBytes] = byte(plain)
		}

		key = ((cipher + key)*c1 + c2) & 0xffff
	}

	return plainBytes
}

// convertHexPair converts two hex characters to a single byte.
func convertHexPair(high, low rune) byte {
	highByte := hexMap[high]
	lowByte := hexMap[low]
	return (highByte << 4) | lowByte
}

var hexMap = map[rune]byte{
	'0': 0x00, '1': 0x01, '2': 0x02, '3': 0x03, '4': 0x04,
	'5': 0x05, '6': 0x06, '7': 0x07, '8': 0x08, '9': 0x09,
	'A': 0x0A, 'a': 0x0A,
	'B': 0x0B, 'b': 0x0B,
	'C': 0x0C, 'c': 0x0C,
	'D': 0x0D, 'd': 0x0D,
	'E': 0x0E, 'e': 0x0E,
	'F': 0x0F, 'f': 0x0F,
}

func asType1Token(t type1.TokenLike) *type1.Type1Token {
	if tt, ok := t.(*type1.Type1Token); ok {
		return tt
	}
	return nil
}

// readExpected reads the next token and validates its type and optional text.
func readExpected(tokenizer *type1.Type1Tokenizer, typ type1.TokenType, texts ...string) {
	text := ""
	if len(texts) > 0 {
		text = texts[0]
	}

	token := tokenizer.GetNext()
	if token == nil {
		panic(fmt.Sprintf("type 1 encrypted portion ended when a token with text '%s' was expected", text))
	}

	tt, ok := token.(*type1.Type1Token)
	if !ok || tt.Type != typ || (text != "" && tt.Text != text) {
		panic(fmt.Sprintf("found invalid token %v when type %v with text %s was expected", token, typ, text))
	}
}

// readExpectedAfterOptional reads the next token, accepting either the primary or optional match.
func readExpectedAfterOptional(tokenizer *type1.Type1Tokenizer, optionalType type1.TokenType, optionalText string, typ type1.TokenType, text string) {
	token := tokenizer.GetNext()
	if token == nil {
		panic(fmt.Sprintf("type 1 encrypted portion ended when a token with text '%s' or '%s' was expected", optionalText, text))
	}

	tt := asType1Token(token)
	if tt != nil && tt.Type == typ && tt.Text == text {
		return
	}

	if tt != nil && tt.Type == optionalType && tt.Text == optionalText {
		readExpected(tokenizer, typ, text)
		return
	}

	panic(fmt.Sprintf("found invalid token %v when type %v with text %s was expected", token, typ, text))
}

// readProcedure reads tokens inside a PostScript procedure block.
func readProcedure(tokenizer *type1.Type1Tokenizer, hasReadStartProc bool) []type1.TokenLike {
	tokens := make([]type1.TokenLike, 0)
	depth := -1
	if hasReadStartProc {
		depth = 1
	}
	readProcedureRecursive(tokenizer, &tokens, &depth)
	return tokens
}

func readProcedureRecursive(tokenizer *type1.Type1Tokenizer, tokens *[]type1.TokenLike, depth *int) {
	if *depth == -1 {
		readExpected(tokenizer, type1.TokenTypeStartProc)
		*depth = 1
	}

	if *depth == 0 {
		return
	}

	for {
		token := tokenizer.GetNext()
		if token == nil {
			break
		}

		tt := asType1Token(token)
		if tt == nil {
			continue
		}

		if tt.Type == type1.TokenTypeStartProc {
			*depth++
			readProcedureRecursive(tokenizer, tokens, depth)
		} else if tt.Type == type1.TokenTypeEndProc {
			*depth--
			break
		} else {
			*tokens = append(*tokens, token)
		}
	}
}

// readTillDef reads tokens until a 'def', 'ND', or '|-' name is encountered.
func readTillDef(tokenizer *type1.Type1Tokenizer, skip bool) {
	current := tokenizer.CurrentToken()
	if current != nil {
		tt := asType1Token(current)
		if tt != nil && tt.Text == "def" {
			return
		}
	}

	for {
		token := tokenizer.GetNext()
		if token == nil {
			break
		}

		tt := asType1Token(token)
		if tt == nil {
			continue
		}

		if tt.Type == type1.TokenTypeNone {
			if tt.Text == "ND" || tt.Text == "|-" {
				return
			}
			if tt.Text == "def" {
				break
			}
			if tt.Text == "systemdict" {
				readExpected(tokenizer, type1.TokenTypeLiteral, "internaldict")
				readExpected(tokenizer, type1.TokenTypeNone, "known")
				readExpected(tokenizer, type1.TokenTypeStartProc)

				skipToken := tokenizer.GetNext()
				for skipToken != nil {
					st := asType1Token(skipToken)
					if st != nil && st.Type == type1.TokenTypeNone && (st.Text == "ND" || st.Text == "def") {
						return
					}
					skipToken = tokenizer.GetNext()
				}
			}
		} else if !skip {
			panic(fmt.Sprintf("encountered unexpected non-name token while reading till 'def' token: %v", token))
		}
	}
}

// readTillPut reads tokens until a 'put', 'NP', or '|' name is encountered.
func readTillPut(tokenizer *type1.Type1Tokenizer) {
	for {
		token := tokenizer.GetNext()
		if token == nil {
			break
		}

		tt := asType1Token(token)
		if tt == nil {
			continue
		}

		if tt.Text == "put" {
			break
		}

		switch tt.Text {
		case "NP", "|":
			return
		}
	}
}

// readArrayValues reads an array of values from the tokenizer using the given converter.
func readArrayValues[T any](tokenizer *type1.Type1Tokenizer, converter func(*type1.Type1Token) T) []T {
	readExpected(tokenizer, type1.TokenTypeStartArray)
	return readArrayValuesInternal(tokenizer, converter, true)
}

func readArrayValuesIntWithStart(tokenizer *type1.Type1Tokenizer, hasReadStart bool) []int {
	return readArrayValuesInternal(tokenizer, func(t *type1.Type1Token) int { return t.AsInt() }, hasReadStart)
}

func readArrayValuesInternal[T any](tokenizer *type1.Type1Tokenizer, converter func(*type1.Type1Token) T, includeDef bool) []T {
	results := make([]T, 0)

	for {
		token := tokenizer.GetNext()
		if token == nil {
			break
		}

		tt := asType1Token(token)
		if tt == nil {
			continue
		}

		if tt.Type == type1.TokenTypeEndArray {
			break
		}

		if tt.Type == type1.TokenTypeStartArray {
			nested := readArrayValuesInternal(tokenizer, converter, false)
			results = append(results, nested...)
			continue
		}

		result := converter(tt)
		results = append(results, result)
	}

	if includeDef {
		readTillDef(tokenizer, false)
	}

	return results
}

// readNumeric reads the next token and returns its numeric value.
func readNumeric(tokenizer *type1.Type1Tokenizer) float64 {
	token := tokenizer.GetNext()

	if token == nil {
		panic("expected to read a numeric token but reached end of input")
	}

	tt := asType1Token(token)
	if tt == nil || (tt.Type != type1.TokenTypeInteger && tt.Type != type1.TokenTypeReal) {
		panic(fmt.Sprintf("expected to read a numeric token, instead got: %v", token))
	}

	return tt.AsDouble()
}

// readBoolean reads the next token and returns its boolean value.
func readBoolean(tokenizer *type1.Type1Tokenizer) bool {
	token := tokenizer.GetNext()

	if token == nil {
		panic("expected to read a boolean token but reached end of input")
	}

	tt := asType1Token(token)
	if tt == nil || (tt.Text != "true" && tt.Text != "false") {
		panic(fmt.Sprintf("expected to read a boolean token, instead got: %v", token))
	}

	return tt.AsBool()
}

// readOtherSubroutines reads the /OtherSubrs dictionary entry.
func readOtherSubroutines(tokenizer *type1.Type1Tokenizer, isLenientParsing bool) {
	start := tokenizer.GetNext()

	if start == nil {
		return
	}

	st := asType1Token(start)

	if st != nil && st.Type == type1.TokenTypeStartArray {
		readArrayValuesInternal(tokenizer, func(t *type1.Type1Token) type1.TokenLike { return t }, false)
	} else if st != nil && (st.Type == type1.TokenTypeInteger || st.Type == type1.TokenTypeReal) {
		length := st.AsInt()
		readExpected(tokenizer, type1.TokenTypeNone, "array")

		for i := 0; i < length; i++ {
			readExpected(tokenizer, type1.TokenTypeNone, "dup")
			readNumeric(tokenizer)
			readTillPut(tokenizer)
		}
		readTillDef(tokenizer, false)
	} else if !isLenientParsing {
		panic(fmt.Sprintf("failed to read start of /OtherSubrs array. Got start token: %v", start))
	}
}

// readSubroutines reads the /Subrs dictionary entry and returns decrypted subroutine bytes.
func readSubroutines(tokenizer *type1.Type1Tokenizer, lenIv int, isLenientParsing bool) []*charstrings.Type1CharstringDecryptedBytes {
	length := int(readNumeric(tokenizer))
	subroutines := make([]*charstrings.Type1CharstringDecryptedBytes, 0, length)

	readExpected(tokenizer, type1.TokenTypeNone, "array")

	for i := 0; i < length; i++ {
		current := tokenizer.GetNext()
		ct := asType1Token(current)
		if ct == nil || ct.Type != type1.TokenTypeNone || ct.Text != "dup" {
			break
		}

		index := int(readNumeric(tokenizer))
		byteLength := int(readNumeric(tokenizer))

		charstring := tokenizer.GetNext()
		charstringToken, ok := charstring.(*type1.Type1DataToken)
		if !ok {
			panic(fmt.Sprintf("found an unexpected token instead of subroutine charstring: %v", charstring))
		}

		if !isLenientParsing && len(charstringToken.Data) != byteLength {
			panic(fmt.Sprintf("the subroutine charstring %s did not have the expected length of %d", charstringToken.String(), byteLength))
		}

		subroutine := decrypt(charstringToken.Data, charstringEncKey, lenIv)
		subroutines = append(subroutines, charstrings.NewType1CharstringDecryptedBytes(subroutine, index))
		readTillPut(tokenizer)
	}

	readTillDef(tokenizer, false)

	return subroutines
}

// readCharStrings reads the CharStrings dictionary and returns decrypted charstring bytes.
func readCharStrings(tokenizer *type1.Type1Tokenizer, lenIv int, isLenientParsing bool) []*charstrings.Type1CharstringDecryptedBytes {
	length := int(readNumeric(tokenizer))

	readExpected(tokenizer, type1.TokenTypeNone, "dict")

	readExpected(tokenizer, type1.TokenTypeNone)
	readExpected(tokenizer, type1.TokenTypeNone, "begin")

	results := make([]*charstrings.Type1CharstringDecryptedBytes, 0)

	for i := 0; i < length; i++ {
		token := tokenizer.GetNext()

		tt := asType1Token(token)
		if tt != nil && tt.Type == type1.TokenTypeNone && tt.Text == "end" {
			break
		}

		if tt == nil || tt.Type != type1.TokenTypeLiteral {
			panic(fmt.Sprintf("type 1 font error. Expected literal for charstring name, instead got: %v", token))
		}

		name := tt.Text

		charstringLength := readNumeric(tokenizer)
		charstring := tokenizer.GetNext()
		charstringToken, ok := charstring.(*type1.Type1DataToken)
		if !ok {
			panic(fmt.Sprintf("got wrong type of token, expected charstring, instead got: %v", charstring))
		}

		if !isLenientParsing && len(charstringToken.Data) != int(charstringLength) {
			panic(fmt.Sprintf("the charstring %s did not have the expected length of %g", charstringToken.String(), charstringLength))
		}

		data := decrypt(charstringToken.Data, charstringEncKey, lenIv)

		results = append(results, charstrings.NewType1CharstringWithName(name, data, i))

		readTillDef(tokenizer, false)
	}

	readExpected(tokenizer, type1.TokenTypeNone, "end")

	return results
}
