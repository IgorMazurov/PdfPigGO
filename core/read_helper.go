package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// AsciiLineFeed is the line-feed '\n' character.
const AsciiLineFeed byte = 10

// AsciiCarriageReturn is the carriage return '\r' character.
const AsciiCarriageReturn byte = 13

var endOfNameCharacters = map[byte]bool{
	' ':              true,
	AsciiCarriageReturn: true,
	AsciiLineFeed:     true,
	9:                 true, // tab
	'>':               true,
	'<':               true,
	'[':               true,
	'/':               true,
	']':               true,
	')':               true,
	'(':               true,
	0:                 true,
	'\f':              true, // form feed
}

// ReadLine reads a string from the input until a newline.
func ReadLine(bytes InputBytes) (string, error) {
	if bytes == nil {
		return "", errors.New("bytes cannot be nil")
	}

	if bytes.IsAtEnd() {
		return "", errors.New("error: end-of-file, expected line")
	}

	var buf strings.Builder
	buf.Grow(11)

	var c byte = 0
	for bytes.MoveNext() {
		c = bytes.CurrentByte()

		if IsEndOfLineByte(c) {
			break
		}

		buf.WriteByte(c)
	}

	if c == AsciiCarriageReturn {
		if peek, ok := bytes.Peek(); ok && peek == AsciiLineFeed {
			bytes.MoveNext()
		}
	}

	return buf.String(), nil
}

// SkipSpaces skips any whitespace characters and PDF comments.
func SkipSpaces(bytes InputBytes) {
	const commentCharacter byte = 37 // '#'

	bytes.MoveNext()
	c := bytes.CurrentByte()

	for IsWhitespace(c) || c == commentCharacter {
		if c == commentCharacter {
			for bytes.MoveNext() {
				c = bytes.CurrentByte()
				if IsEndOfLineByte(c) {
					break
				}
			}
		}
		bytes.MoveNext()
		c = bytes.CurrentByte()
	}

	if !bytes.IsAtEnd() {
		bytes.Seek(bytes.CurrentOffset()-1, 0) // io.SeekStart
	}
}

// IsEndOfName reports whether the given byte value is the end of a PDF Name token.
func IsEndOfName(ch byte) bool {
	return endOfNameCharacters[ch]
}

// IsWhitespace determines if a byte is whitespace or not, including newlines.
// These values are specified in table 1 (page 12) of ISO 32000-1:2008.
func IsWhitespace(c byte) bool {
	return c == 0 || c == 32 || c == AsciiLineFeed || c == AsciiCarriageReturn || c == 9 || c == 12
}

// IsEndOfLineChar reports whether the character is an end of line character.
func IsEndOfLineChar(c rune) bool {
	return IsEndOfLineByte(byte(c))
}

// IsEndOfLineByte reports whether the byte is an end of line character.
func IsEndOfLineByte(b byte) bool {
	return IsLineFeed(b) || IsCarriageReturn(b)
}

// IsLineFeed reports whether the byte is a line feed '\n' character.
func IsLineFeed(c byte) bool {
	return c == AsciiLineFeed
}

// IsCarriageReturn reports whether the byte is a carriage return '\r' character.
func IsCarriageReturn(c byte) bool {
	return c == AsciiCarriageReturn
}

// IsStringStr reports whether the given string is at this position in the input.
// Resets to the current offset once read.
func IsStringStr(bytes InputBytes, s string) bool {
	found := true

	startOffset := bytes.CurrentOffset()

	for i := 0; i < len(s); i++ {
		if bytes.IsAtEnd() && i < len(s)-1 {
			found = false
			break
		}
		if bytes.CurrentByte() != s[i] {
			found = false
			break
		}
		bytes.MoveNext()
	}

	bytes.Seek(startOffset, 0) // io.SeekStart

	return found
}

// IsStringBytes reports whether the given byte slice is at this position in the input.
// Resets to the current offset once read.
func IsStringBytes(bytes InputBytes, s []byte) bool {
	found := true

	startOffset := bytes.CurrentOffset()

	for i := 0; i < len(s); i++ {
		if bytes.CurrentByte() != s[i] {
			found = false
			break
		}
		bytes.MoveNext()
	}

	bytes.Seek(startOffset, 0) // io.SeekStart

	return found
}

// ReadLong reads a long from the input.
func ReadLong(bytes InputBytes) (int64, error) {
	SkipSpaces(bytes)

	buffer := make([]byte, 19) // max formatted int64 length
	bytesRead := readNumberAsUtf8Bytes(bytes, buffer)

	if bytesRead < 0 {
		return 0, errors.New("number exceeded maximum length")
	}

	longStr := string(buffer[:bytesRead])
	result, err := strconv.ParseInt(longStr, 10, 64)
	if err != nil {
		bytes.Seek(bytes.CurrentOffset()-int64(bytesRead), 0) // io.SeekStart
		return 0, fmt.Errorf("error: expected a long type at offset %d, instead got '%s'",
			bytes.CurrentOffset(), BytesAsLatin1String(buffer[:bytesRead]))
	}

	return result, nil
}

// ReadInt reads an int from the input.
func ReadInt(bytes InputBytes) (int64, error) {
	if bytes == nil {
		return 0, errors.New("bytes cannot be nil")
	}

	SkipSpaces(bytes)

	buffer := make([]byte, 10) // max formatted int32 length
	bytesRead := readNumberAsUtf8Bytes(bytes, buffer)

	if bytesRead < 0 {
		return 0, errors.New("number exceeded maximum length")
	}

	intStr := string(buffer[:bytesRead])
	result, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		bytes.Seek(bytes.CurrentOffset()-int64(bytesRead), 0) // io.SeekStart
		return 0, NewPdfDocumentFormatException(
			fmt.Sprintf("error: expected an integer type at offset %d, instead got '%s'",
				bytes.CurrentOffset(), BytesAsLatin1String(buffer[:bytesRead])))
	}

	return result, nil
}

// IsDigit reports whether the given byte value is a digit.
func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// IsHexByte reports whether the given byte value is a valid hex digit.
func IsHexByte(b byte) bool {
	return IsHexChar(rune(b))
}

// IsHexChar reports whether the given character value is a valid hex digit.
func IsHexChar(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// IsValidUtf8 reports whether the given byte slice is valid UTF-8.
func IsValidUtf8(input []byte) bool {
	return strings.ToValidUTF8(string(input), "") == string(input)
}

// readNumberAsUtf8Bytes reads bytes until a delimiter and stores them in buffer.
// Returns the number of bytes read, or -1 if the buffer is too small.
func readNumberAsUtf8Bytes(reader InputBytes, buffer []byte) int {
	position := 0

	for reader.MoveNext() {
		lastByte := reader.CurrentByte()

		if lastByte == ' ' || lastByte == AsciiLineFeed || lastByte == AsciiCarriageReturn ||
			lastByte == '<' || // see sourceforge bug 1714707
			lastByte == '[' || // PDFBOX-1845
			lastByte == '(' || // PDFBOX-2579
			lastByte == 0 {
			break
		}

		if position >= len(buffer) {
			return -1
		}

		buffer[position] = lastByte
		position++
	}

	if !reader.IsAtEnd() {
		reader.Seek(reader.CurrentOffset()-1, 0) // io.SeekStart
	}

	return position
}
