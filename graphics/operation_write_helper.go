package graphics

import (
	"io"
	"math"
	"strconv"
	"unicode/utf8"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/util"
)

// byteWriter combines io.Writer with WriteByte for efficient single-byte writes.
type byteWriter interface {
	io.Writer
	WriteByte(c byte) error
}

const (
	whitespaceByte = ' '
	newlineByte    = '\n'
	zeroByte       = '0'
	pointByte      = '.'
)

// WriteText writes the string to the writer using ASCII or Latin-1 encoding.
// If appendWhitespace is true, a space byte is appended after the text.
func WriteText(w byteWriter, text string, appendWhitespace bool) error {
	if utf8.ValidString(text) && isASCII(text) {
		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}
	} else {
		bytes := core.StringAsLatin1Bytes(text)
		if _, err := w.Write(bytes); err != nil {
			return err
		}
	}

	if appendWhitespace {
		if err := w.WriteByte(whitespaceByte); err != nil {
			return err
		}
	}

	return nil
}

// WriteTextBytes writes raw bytes to the writer.
// If appendWhitespace is true, a space byte is appended after the text.
func WriteTextBytes(w byteWriter, asciiBytes []byte, appendWhitespace bool) error {
	if _, err := w.Write(asciiBytes); err != nil {
		return err
	}

	if appendWhitespace {
		if err := w.WriteByte(whitespaceByte); err != nil {
			return err
		}
	}

	return nil
}

// WriteHex writes the bytes as a hex string wrapped in angle brackets (<...>).
func WriteHex(w byteWriter, bytes []byte) error {
	hex := make([]byte, len(bytes)*2)
	util.GetUtf8Chars(bytes, hex)

	if err := w.WriteByte('<'); err != nil {
		return err
	}
	if _, err := w.Write(hex); err != nil {
		return err
	}
	return w.WriteByte('>')
}

// WriteWhiteSpace writes a single space byte to the writer.
func WriteWhiteSpace(w byteWriter) error {
	return w.WriteByte(whitespaceByte)
}

// WriteNewLine writes a newline byte to the writer.
func WriteNewLine(w byteWriter) error {
	return w.WriteByte(newlineByte)
}

// WriteDouble formats a double with up to 9 decimal places, stripping trailing
// zeros and the trailing decimal point if no fractional digits remain.
func WriteDouble(w byteWriter, value float64) error {
	formatted := formatDouble(value)
	_, err := w.Write(formatted)
	return err
}

func formatDouble(value float64) []byte {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		return []byte(strconv.FormatFloat(value, 'f', -1, 64))
	}

	formatted := strconv.FormatFloat(value, 'f', 9, 64)

	buf := []byte(formatted)
	lastIndex := getLastSignificantDigitIndex(buf, len(buf))

	return buf[:lastIndex]
}

func getLastSignificantDigitIndex(buffer []byte, bytesWritten int) int {
	lastIndex := bytesWritten
	for i := bytesWritten - 1; i > 1; i-- {
		if buffer[i] != zeroByte {
			break
		}
		lastIndex--
	}

	if buffer[lastIndex-1] == pointByte {
		lastIndex--
	}

	return lastIndex
}

// WriteNumberTextInt writes a number as double, then whitespace, text, and newline.
func WriteNumberTextInt(w byteWriter, number int, text string) error {
	if err := WriteDouble(w, float64(number)); err != nil {
		return err
	}
	if err := w.WriteByte(whitespaceByte); err != nil {
		return err
	}
	if err := WriteText(w, text, false); err != nil {
		return err
	}
	return WriteNewLine(w)
}

// WriteNumberTextIntBytes writes a number as double, then whitespace, raw bytes, and newline.
func WriteNumberTextIntBytes(w byteWriter, number int, asciiBytes []byte) error {
	if err := WriteDouble(w, float64(number)); err != nil {
		return err
	}
	if err := w.WriteByte(whitespaceByte); err != nil {
		return err
	}
	if err := WriteTextBytes(w, asciiBytes, false); err != nil {
		return err
	}
	return WriteNewLine(w)
}

// WriteNumberTextFloat writes a number as double, then whitespace, text, and newline.
func WriteNumberTextFloat(w byteWriter, number float64, text string) error {
	if err := WriteDouble(w, number); err != nil {
		return err
	}
	if err := w.WriteByte(whitespaceByte); err != nil {
		return err
	}
	if err := WriteText(w, text, false); err != nil {
		return err
	}
	return WriteNewLine(w)
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}

// The following are convenience wrappers that match the C# extension method pattern
// where callers pass a stream and call methods on it. These exist for compatibility
// with code that was ported directly from C#.

// WriteTextNoWhitespace is an alias for WriteText with appendWhitespace=false.
func WriteTextNoWhitespace(w byteWriter, text string) error {
	return WriteText(w, text, false)
}

// FormatDouble returns the formatted double as a string (for cases where
// the caller needs the string rather than writing directly to a stream).
func FormatDouble(value float64) string {
	buf := formatDouble(value)
	return string(buf)
}
