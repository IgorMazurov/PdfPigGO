package export

import (
	"fmt"
	"strings"
)

// InvalidCharHandler transforms a string by handling invalid XML characters.
type InvalidCharHandler func(string) string

// GetXmlInvalidCharHandler returns an InvalidCharHandler for the given strategy.
func GetXmlInvalidCharHandler(strategy InvalidCharStrategy) InvalidCharHandler {
	switch strategy {
	case DoNotCheck:
		return func(s string) string { return s }

	case Remove:
		return func(s string) string {
			if s == "" {
				return s
			}
			var b strings.Builder
			b.Grow(len(s))
			for _, r := range s {
				if isValidXmlChar(r) {
					b.WriteRune(r)
				}
			}
			return b.String()
		}

	case ConvertToHexadecimal:
		return func(s string) string {
			if s == "" {
				return s
			}
			var b strings.Builder
			b.Grow(len(s))
			for _, r := range s {
				if isValidXmlChar(r) {
					b.WriteRune(r)
				} else {
					bytes := []byte(string(r))
					hexParts := make([]string, len(bytes))
					for i, byt := range bytes {
						hexParts[i] = fmt.Sprintf("%02X", byt)
					}
					b.WriteString("0x")
					b.WriteString(strings.Join(hexParts, "-"))
				}
			}
			return b.String()
		}

	default:
		panic(fmt.Sprintf("unsupported InvalidCharStrategy: %d", strategy))
	}
}

// RestoreEscapedInvalidChars reverses the placeholder escaping done by escapeInvalidXmlCharsPlaceholder.
func RestoreEscapedInvalidChars(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\x90' && i+1 < len(s) && s[i+1] == '\x91' {
			// This is a placeholder for an invalid char - but we can't recover the original byte
			// since it was already lost. Instead, skip these markers.
			i += 1 // skip \x91
			b.WriteRune('\ufffd') // fallback to replacement char
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// isValidXmlChar reports whether r is a valid XML 1.0 character per the W3C spec.
// Valid chars: #x9 | #xA | #xD | [#x20-#xD7FF] | [#xE000-#xFFFD] | [#x10000-#x10FFFF]
func isValidXmlChar(r rune) bool {
	return r == '\t' || r == '\n' || r == '\r' ||
		(r >= ' ' && r <= '\uD7FF') ||
		(r >= '\uE000' && r <= '\uFFFD') ||
		(r >= '\U00010000' && r <= '\U0010FFFF')
}
