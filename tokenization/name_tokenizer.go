package tokenization

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// NameTokenizer tokenizes PDF name objects starting with '/'.
type NameTokenizer struct{}

var _ Tokenizer = (*NameTokenizer)(nil)

// NewNameTokenizer creates a new NameTokenizer.
func NewNameTokenizer() *NameTokenizer {
	return &NameTokenizer{}
}

// ReadsNextByte returns true because this tokenizer reads ahead past the name.
func (t *NameTokenizer) ReadsNextByte() bool {
	return true
}

// Tokenize attempts to tokenize a PDF name starting with '/'.
// Returns the NameToken and true on success, or nil and false if currentByte is not '/'.
func (t *NameTokenizer) Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool) {
	if currentByte != '/' {
		return nil, false
	}

	var buf []byte
	escapeActive := false
	postEscapeRead := 0
	escapedChars := [2]byte{}

	for input.MoveNext() {
		b := input.CurrentByte()

		if b == '#' {
			escapeActive = true
		} else if escapeActive {
			if core.IsHexByte(b) {
				escapedChars[postEscapeRead] = b
				postEscapeRead++

				if postEscapeRead == 2 {
					high := hexValue(escapedChars[0])
					low := hexValue(escapedChars[1])

					buf = append(buf, byte(high*16+low))

					escapeActive = false
					postEscapeRead = 0
				}
			} else {
				buf = append(buf, '#')

				if postEscapeRead == 1 {
					buf = append(buf, escapedChars[0])
				}

				if core.IsEndOfName(b) {
					break
				}

				if b == '#' {
					escapeActive = true
					postEscapeRead = 0
					continue
				}

				buf = append(buf, b)
				escapeActive = false
				postEscapeRead = 0
			}
		} else if core.IsEndOfName(b) {
			break
		} else {
			buf = append(buf, b)
		}
	}

	str := decodeNameBytes(buf)

	return tokens.Create(str), true
}

func hexValue(b byte) int {
	if b >= '0' && b <= '9' {
		return int(b - '0')
	}
	if b >= 'A' && b <= 'F' {
		return int(b - 'A') + 10
	}
	return int(b - 'a') + 10
}

func decodeNameBytes(data []byte) string {
	if utf8.Valid(data) {
		return string(data)
	}

	decoded, err := charmap.Windows1252.NewDecoder().String(string(data))
	if err != nil {
		return strings.ToValidUTF8(string(data), "")
	}
	return decoded
}
