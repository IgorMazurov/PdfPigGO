package fonts

import (
	"bytes"
	"fmt"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// ToUnicodeCMapBuilder builds a ToUnicode CMap stream from a Unicode-to-character-code mapping.
type ToUnicodeCMapBuilder struct{}

// ConvertToCMapStream converts a map of Unicode characters to character codes into
// a CMap byte stream suitable for embedding as a ToUnicode stream in a PDF font dictionary.
func (b *ToUnicodeCMapBuilder) ConvertToCMapStream(unicodeToCharacterCode map[rune]byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	writeName(buf, tokens.CidInit)
	writeName(buf, tokens.ProcSet)
	writeTextBytes(buf, "findresource", true)
	writeTextBytes(buf, "begin", false)
	writeNewLine(buf)

	writeDouble(buf, 12)
	buf.WriteByte(' ')
	writeTextBytes(buf, "dict", true)
	writeTextBytes(buf, "begin", false)
	writeNewLine(buf)

	writeTextBytes(buf, "begincmap", false)
	writeNewLine(buf)

	writeName(buf, tokens.CidSystemInfo)

	dict, err := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.Registry:   tokens.NewStringToken("Adobe"),
		tokens.Ordering:   tokens.NewStringToken("UCS"),
		tokens.Supplement: tokens.NewNumericTokenFromInt(0),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CMap dictionary: %w", err)
	}

	writeDictionary(buf, dict)
	buf.WriteByte(' ')

	writeTextBytes(buf, "def", false)
	writeNewLine(buf)

	writeName(buf, tokens.Cmapname)
	writeName(buf, tokens.Create("Adobe-Identity-UCS"))
	writeTextBytes(buf, "def", false)
	writeNewLine(buf)

	writeName(buf, tokens.CmapType)
	writeNumberTextInt(buf, 2, "def")
	writeNumberTextInt(buf, 1, "begincodespacerange")
	writeHex(buf, []byte{0x00})
	writeHex(buf, []byte{0xFF})
	writeNewLine(buf)

	writeTextBytes(buf, "endcodespacerange", false)
	writeNewLine(buf)

	writeNumberTextInt(buf, len(unicodeToCharacterCode), "beginbfchar")

	for unicodeChar, charCode := range unicodeToCharacterCode {
		unicodeInt := uint16(unicodeChar)
		low := byte(unicodeInt >> 0)
		high := byte(unicodeInt >> 8)

		writeHex(buf, []byte{charCode})
		writeHex(buf, []byte{high, low})
		writeNewLine(buf)
	}

	writeTextBytes(buf, "endbfchar", false)
	writeNewLine(buf)

	writeTextBytes(buf, "endcmap", false)
	writeNewLine(buf)

	writeTextBytes(buf, "CMapName currentdict /CMap defineresource pop", false)
	writeNewLine(buf)

	writeTextBytes(buf, "end", false)
	writeNewLine(buf)

	writeTextBytes(buf, "end", false)
	writeNewLine(buf)

	return buf.Bytes(), nil
}

// writeName writes a PDF name token (e.g., "/CIDInit ") to the buffer with trailing whitespace.
func writeName(buf *bytes.Buffer, name *tokens.NameToken) {
	buf.WriteByte('/')
	buf.WriteString(name.Data())
	buf.WriteByte(' ')
}

// writeTextBytes writes raw ASCII bytes to the buffer.
// If appendWhitespace is true, a space byte is appended after the text.
func writeTextBytes(buf *bytes.Buffer, ascii string, appendWhitespace bool) {
	buf.WriteString(ascii)
	if appendWhitespace {
		buf.WriteByte(' ')
	}
}

// writeNewLine writes a newline byte to the buffer.
func writeNewLine(buf *bytes.Buffer) {
	buf.WriteByte('\n')
}

// writeDouble formats a float64 with up to 9 decimal places, stripping trailing
// zeros and the trailing decimal point if no fractional digits remain.
func writeDouble(buf *bytes.Buffer, value float64) {
	formatted := formatDouble(value)
	buf.Write(formatted)
}

// formatDouble returns the byte representation of a double with trailing zeros removed.
func formatDouble(value float64) []byte {
	formatted := fmt.Sprintf("%.9f", value)
	buf := []byte(formatted)
	lastIndex := len(buf)
	for i := len(buf) - 1; i > 1; i-- {
		if buf[i] != '0' {
			break
		}
		lastIndex--
	}
	if lastIndex > 0 && buf[lastIndex-1] == '.' {
		lastIndex--
	}
	return buf[:lastIndex]
}

// writeNumberTextInt writes a number as double, then whitespace, text, and newline.
func writeNumberTextInt(buf *bytes.Buffer, number int, text string) {
	writeDouble(buf, float64(number))
	buf.WriteByte(' ')
	writeTextBytes(buf, text, false)
	writeNewLine(buf)
}

// writeHex writes bytes as a hex string wrapped in angle brackets (<...>).
func writeHex(buf *bytes.Buffer, b []byte) {
	hexChars := "0123456789ABCDEF"
	buf.WriteByte('<')
	for _, byteVal := range b {
		buf.WriteByte(hexChars[(byteVal>>4)&0xF])
		buf.WriteByte(hexChars[byteVal&0xF])
	}
	buf.WriteByte('>')
}

// writeDictionary writes a PDF dictionary token (<< ... >>) to the buffer.
func writeDictionary(buf *bytes.Buffer, dict *tokens.DictionaryToken) {
	if dict == nil {
		return
	}
	buf.WriteString("<<")
	data := dict.Data()
	for key, value := range data {
		buf.WriteString(" /")
		buf.WriteString(key)
		buf.WriteString(" ")
		buf.WriteString(tokenToString(value))
	}
	buf.WriteString(" >>")
}

// tokenToString converts a Token to its PDF string representation.
func tokenToString(t tokens.Token) string {
	switch v := t.(type) {
	case *tokens.StringToken:
		return fmt.Sprintf("(%s)", v.Data())
	case *tokens.NumericToken:
		return v.String()
	case *tokens.NameToken:
		return "/" + v.Data()
	default:
		return fmt.Sprint(t)
	}
}
