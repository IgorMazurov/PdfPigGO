package tokenization

import (
	"encoding/binary"
	"strings"
	"unicode/utf16"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// StringTokenizer tokenizes PDF parenthesized literal strings.
type StringTokenizer struct {
	usePdfDocEncoding bool
	stringBuilder     strings.Builder
}

var _ Tokenizer = (*StringTokenizer)(nil)

// NewStringTokenizer creates a new StringTokenizer.
func NewStringTokenizer(usePdfDocEncoding bool) *StringTokenizer {
	return &StringTokenizer{
		usePdfDocEncoding: usePdfDocEncoding,
	}
}

// ReadsNextByte returns false because this tokenizer does not read ahead past the token.
func (t *StringTokenizer) ReadsNextByte() bool {
	return false
}

// Tokenize attempts to tokenize a parenthesized string starting with '('.
func (t *StringTokenizer) Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool) {
	if currentByte != '(' {
		return nil, false
	}

	builder := &t.stringBuilder
	numberOfBrackets := 1
	isEscapeActive := false
	isLineBreaking := false

	octalModeActive := false
	octal := [3]int{0, 0, 0}
	octalsRead := 0

	for input.MoveNext() {
		b := input.CurrentByte()
		c := rune(b)

		if octalModeActive {
			nextCharacterOctal := c >= '0' && c <= '7'

			if nextCharacterOctal {
				leftShiftOctal(c, octalsRead, &octal)
				octalsRead++
			}

			if octalsRead == 3 || !nextCharacterOctal {
				characterCode := core.FromOctalDigits(octal[:])
				builder.WriteRune(rune(characterCode))

				octal[0] = 0
				octal[1] = 0
				octal[2] = 0
				octalsRead = 0
				octalModeActive = false
			}

			if nextCharacterOctal {
				continue
			}
		}

		switch c {
		case ')':
			isLineBreaking = false
			if !isEscapeActive {
				numberOfBrackets--
			}

			isEscapeActive = false
			if numberOfBrackets > 0 {
				builder.WriteRune(c)
			}

			numberOfBrackets = checkForEndOfString(numberOfBrackets, input)

		case '(':
			isLineBreaking = false

			if !isEscapeActive {
				numberOfBrackets++
			}

			isEscapeActive = false
			builder.WriteRune(c)

		case '\\':
			isLineBreaking = false
			if isEscapeActive {
				builder.WriteRune(c)
				isEscapeActive = false
			} else {
				isEscapeActive = true
			}

		default:
			if isLineBreaking {
				if core.IsEndOfLineChar(c) {
					continue
				}

				isLineBreaking = false
				builder.WriteRune(c)
			} else if isEscapeActive {
				processEscapedCharacter(c, builder, &octal, &octalModeActive, &octalsRead, &isLineBreaking)
				isEscapeActive = false
			} else {
				builder.WriteRune(c)
			}
		}

		if numberOfBrackets <= 0 {
			break
		}
	}

	tokenStr, encodedWith := t.decodeString(builder.String())
	builder.Reset()

	return tokens.NewStringToken(tokenStr, encodedWith), true
}

func (t *StringTokenizer) decodeString(builtStr string) (string, tokens.StringEncoding) {
	runes := []rune(builtStr)

	if len(runes) >= 2 {
		if runes[0] == 0xFE && runes[1] == 0xFF {
			rawBytes := core.StringAsLatin1Bytes(builtStr)
			codeUnits := bytesToUTF16BE(rawBytes)
			tokenStr := string(utf16.Decode(codeUnits[1:]))
			return tokenStr, tokens.Utf16BE
		}

		if runes[0] == 0xFF && runes[1] == 0xFE {
			rawBytes := core.StringAsLatin1Bytes(builtStr)
			codeUnits := bytesToUTF16LE(rawBytes)
			tokenStr := string(utf16.Decode(codeUnits[1:]))
			return tokenStr, tokens.Utf16
		}

		if t.usePdfDocEncoding {
			rawBytes := core.StringAsLatin1Bytes(builtStr)
			if str, ok := core.TryConvertBytesToString(rawBytes); ok {
				return str, tokens.StringEncodingPDFDoc
			}
			return builtStr, tokens.Iso88591
		}

		return builtStr, tokens.Iso88591
	}

	if t.usePdfDocEncoding {
		rawBytes := core.StringAsLatin1Bytes(builtStr)
		if str, ok := core.TryConvertBytesToString(rawBytes); ok {
			return str, tokens.StringEncodingPDFDoc
		}
		return builtStr, tokens.Iso88591
	}

	return builtStr, tokens.Iso88591
}

func bytesToUTF16BE(b []byte) []uint16 {
	units := make([]uint16, len(b)/2)
	for i := 0; i < len(units); i++ {
		units[i] = binary.BigEndian.Uint16(b[i*2 : i*2+2])
	}
	return units
}

func bytesToUTF16LE(b []byte) []uint16 {
	units := make([]uint16, len(b)/2)
	for i := 0; i < len(units); i++ {
		units[i] = binary.LittleEndian.Uint16(b[i*2 : i*2+2])
	}
	return units
}

func leftShiftOctal(nextOctalChar rune, octalsRead int, octals *[3]int) {
	for i := octalsRead; i > 0; i-- {
		octals[i] = octals[i-1]
	}

	value, _ := core.CharacterToShort(nextOctalChar)
	octals[0] = value
}

func processEscapedCharacter(c rune, builder *strings.Builder, octal *[3]int,
	isOctalActive *bool, octalsRead *int, isLineBreaking *bool) {

	switch c {
	case 'n':
		builder.WriteRune('\n')
	case 'r':
		builder.WriteRune('\r')
	case 't':
		builder.WriteRune('\t')
	case 'b':
		builder.WriteRune('\b')
	case 'f':
		builder.WriteRune('\f')
	case '0', '1', '2', '3', '4', '5', '6', '7':
		octal[0], _ = core.CharacterToShort(c)
		*isOctalActive = true
		*octalsRead = 1
	default:
		if c == rune(core.AsciiCarriageReturn) || c == rune(core.AsciiLineFeed) {
			*isLineBreaking = true
		} else {
			builder.WriteRune(c)
		}
	}
}

func checkForEndOfString(numberOfBrackets int, bytes core.InputBytes) int {
	lineFeed := byte(10)
	carriageReturn := byte(13)

	braces := numberOfBrackets
	nextThreeBytes := make([]byte, 3)

	startAt := bytes.CurrentOffset()

	amountRead, _ := bytes.Read(nextThreeBytes)

	if amountRead == 3 && nextThreeBytes[0] == carriageReturn {
		if (nextThreeBytes[1] == lineFeed && (nextThreeBytes[2] == '/' || nextThreeBytes[2] == '>')) ||
			nextThreeBytes[1] == '/' || nextThreeBytes[1] == '>' {
			braces = 0
		}
	}

	if amountRead > 0 {
		bytes.Seek(startAt, 0) // io.SeekStart
	}

	return braces
}
