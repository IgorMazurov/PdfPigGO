// Package fonts provides types for PDF font handling.
package fonts

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/uglytoad/pdfpig/go/fonts/encodings"
)

// GlyphList maps PostScript glyph names to Unicode values.
type GlyphList struct {
	nameToUnicode          map[string]string
	unicodeToName          map[string]string
	oddNameToUnicodeCache  map[string]string
	oddNameCacheMu         sync.RWMutex
}

// NotDefined is the standard glyph name for undefined characters.
const NotDefined = ".notdef"

var (
	adobeGlyphListOnce     sync.Once
	adobeGlyphListInstance *GlyphList
	adobeGlyphListErr      error

	zapfDingbatsOnce       sync.Once
	zapfDingbatsInstance   *GlyphList
	zapfDingbatsErr        error
)

// AdobeGlyphList returns the Adobe Glyph List singleton (includes extension).
func AdobeGlyphList() (*GlyphList, error) {
	adobeGlyphListOnce.Do(func() {
		list, err := glyphListFromFactory("glyphlist", "additional")
		if err != nil {
			adobeGlyphListErr = err
			return
		}
		adobeGlyphListInstance = list
	})
	return adobeGlyphListInstance, adobeGlyphListErr
}

// ZapfDingbats returns the Zapf Dingbats glyph list singleton.
func ZapfDingbats() (*GlyphList, error) {
	zapfDingbatsOnce.Do(func() {
		list, err := glyphListFromFactory("zapfdingbats")
		if err != nil {
			zapfDingbatsErr = err
			return
		}
		zapfDingbatsInstance = list
	})
	return zapfDingbatsInstance, zapfDingbatsErr
}

// NewGlyphList creates a GlyphList from a name-to-Unicode mapping.
func NewGlyphList(namesToUnicode map[string]string) *GlyphList {
	unicodeToNameTemp := make(map[string]string, len(namesToUnicode))

	for name, unicodeVal := range namesToUnicode {
		forceOverride := encodings.WinAnsiEncodingValue.ContainsName(name) ||
			encodings.MacRomanEncodingValue.ContainsName(name) ||
			encodings.MacExpertEncodingValue.ContainsName(name) ||
			encodings.SymbolEncodingValue.ContainsName(name) ||
			encodings.ZapfDingbatsEncodingValue.ContainsName(name)

		if _, exists := unicodeToNameTemp[unicodeVal]; !exists || forceOverride {
			unicodeToNameTemp[unicodeVal] = name
		}
	}

	return &GlyphList{
		nameToUnicode:         namesToUnicode,
		unicodeToName:         unicodeToNameTemp,
		oddNameToUnicodeCache: make(map[string]string),
	}
}

// UnicodeCodePointToName returns the glyph name for a Unicode code point value.
func (g *GlyphList) UnicodeCodePointToName(unicodeValue int) string {
	value := string(rune(unicodeValue))

	if result, ok := g.unicodeToName[value]; ok {
		return result
	}

	return NotDefined
}

// NameToUnicode returns the Unicode value for a glyph name.
// It follows the Adobe Glyph List specification for name decomposition.
func (g *GlyphList) NameToUnicode(name string) (string, error) {
	if name == "" {
		return "", nil
	}

	if unicodeVal, ok := g.nameToUnicode[name]; ok {
		return unicodeVal, nil
	}

	g.oddNameCacheMu.RLock()
	if result, ok := g.oddNameToUnicodeCache[name]; ok {
		g.oddNameCacheMu.RUnlock()
		return result, nil
	}
	g.oddNameCacheMu.RUnlock()

	var unicode string
	var err error

	if idx := strings.Index(name, "."); idx > 0 {
		unicode, err = g.NameToUnicode(name[:idx])
	} else if strings.Contains(name, "_") {
		parts := strings.Split(name, "_")
		var sb strings.Builder
		for _, part := range parts {
			sub, subErr := g.NameToUnicode(part)
			if subErr != nil {
				return "", subErr
			}
			sb.WriteString(sub)
		}
		unicode = sb.String()
		err = nil
	} else if strings.HasPrefix(name, "uni") && (len(name)-3)%4 == 0 {
		unicode, err = g.decodeUniName(name)
	} else if strings.HasPrefix(name, "u") && len(name) >= 5 && len(name) <= 7 {
		unicode, err = g.decodeUName(name)
	} else if strings.HasPrefix(strings.ToLower(name), "c") && len(name) >= 3 && len(name) <= 4 {
		unicode, err = g.decodeCName(name)
	} else {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	g.oddNameCacheMu.Lock()
	g.oddNameToUnicodeCache[name] = unicode
	g.oddNameCacheMu.Unlock()

	return unicode, nil
}

func (g *GlyphList) decodeUniName(name string) (string, error) {
	nameLen := len(name)
	var sb strings.Builder

	for pos := 3; pos+4 <= nameLen; pos += 4 {
		hexStr := name[pos : pos+4]
		codePoint, parseErr := parseHex(hexStr)
		if parseErr != nil {
			return "", parseErr
		}

		if codePoint > 0xD7FF && codePoint < 0xE000 {
			return "", NewInvalidFontFormatException(fmt.Sprintf("Unicode character name with disallowed code area: %s", name))
		}

		sb.WriteRune(rune(codePoint))
	}

	return sb.String(), nil
}

func (g *GlyphList) decodeUName(name string) (string, error) {
	hexStr := name[1:]
	codePoint, parseErr := parseHex(hexStr)
	if parseErr != nil {
		return "", parseErr
	}

	if codePoint > 0xD7FF && codePoint < 0xE000 {
		return "", NewInvalidFontFormatException(fmt.Sprintf("Unicode character name with disallowed code area: %s", name))
	}

	return string(rune(codePoint)), nil
}

func (g *GlyphList) decodeCName(name string) (string, error) {
	numStr := name[1:]
	codePoint, parseErr := parseIntDec(numStr)
	if parseErr != nil {
		return "", parseErr
	}

	return string(rune(codePoint)), nil
}

func parseHex(s string) (int, error) {
	v, err := strconv.ParseInt(s, 16, 32)
	return int(v), err
}

func parseIntDec(s string) (int, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	return int(v), err
}

// AdobeGlyphListNotDefined is the standard glyph name for undefined characters,
// exposed as a package-level constant for use by callers that don't want to hold
// a GlyphList instance. Equivalent to C# GlyphList.NotDefined.
const AdobeGlyphListNotDefined = NotDefined

// AdobeGlyphListUnicodeCodePointToName returns the glyph name for a Unicode code point
// using the Adobe Glyph List singleton. This is a convenience wrapper equivalent to
// C#'s GlyphList.AdobeGlyphList.UnicodeCodePointToName(code). Returns ".notdef" if
// the code point has no mapping or if the glyph list fails to load.
func AdobeGlyphListUnicodeCodePointToName(unicodeValue int) string {
	gl, err := AdobeGlyphList()
	if err != nil || gl == nil {
		return NotDefined
	}
	return gl.UnicodeCodePointToName(unicodeValue)
}

var glyphFactory = GlyphListFactory{}

func glyphListFromFactory(listNames ...string) (*GlyphList, error) {
	return glyphFactory.Get(listNames...)
}

// GlyphListNameToUnicode returns the Unicode value for a glyph name using the
// Adobe Glyph List singleton. Returns ("", false) if the name has no mapping
// or if the glyph list fails to load. Equivalent to C#'s
// GlyphList.AdobeGlyphList.NameToUnicode(name).
func GlyphListNameToUnicode(name string) (string, bool) {
	if name == "" {
		return "", false
	}
	gl, err := AdobeGlyphList()
	if err != nil || gl == nil {
		return "", false
	}
	unicode, err := gl.NameToUnicode(name)
	if err != nil {
		return "", false
	}
	if unicode == "" {
		return "", false
	}
	return unicode, true
}
