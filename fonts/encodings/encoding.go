// Package encodings provides character encoding types for PDF fonts.
package encodings

import (
	"github.com/uglytoad/pdfpig/go/tokens"
)

// NotDefined is the glyph name used when a code has no mapping.
const NotDefined = ".notdef"

// Encoding maps character codes to glyph names from a PostScript encoding.
type Encoding struct {
	CodeToName map[int]string
	NameToCode map[string]int
}

// NewEncoding creates a new Encoding with pre-allocated maps.
func NewEncoding() *Encoding {
	return &Encoding{
		CodeToName: make(map[int]string, 250),
		NameToCode: make(map[string]int, 250),
	}
}

// CodeToNameMap returns the code-to-name mapping as a read-only view.
func (e *Encoding) CodeToNameMap() map[int]string {
	return e.CodeToName
}

// NameToCodeMap returns the name-to-code mapping as a read-only view.
func (e *Encoding) NameToCodeMap() map[string]int {
	return e.NameToCode
}

// EncodingName returns the name of this encoding.
func (e *Encoding) EncodingName() string {
	return ""
}

// ContainsName reports whether the encoding contains a code for the given name.
func (e *Encoding) ContainsName(name string) bool {
	_, ok := e.NameToCode[name]
	return ok
}

// ContainsCode reports whether the encoding contains a name for the given code.
func (e *Encoding) ContainsCode(code int) bool {
	_, ok := e.CodeToName[code]
	return ok
}

// GetName returns the character name corresponding to the given code,
// or NotDefined if not found.
func (e *Encoding) GetName(code int) string {
	name, ok := e.CodeToName[code]
	if !ok {
		return NotDefined
	}
	return name
}

// GetCode returns the character code for the given name, or -1 if not found.
func (e *Encoding) GetCode(name string) int {
	code, ok := e.NameToCode[name]
	if !ok {
		return -1
	}
	return code
}

// Add registers a character code and name pair in both direction maps.
func (e *Encoding) Add(code int, name string) {
	e.CodeToName[code] = name
	if _, exists := e.NameToCode[name]; !exists {
		e.NameToCode[name] = code
	}
}

// StandardEncoding is the standard PostScript encoding vector.
type StandardEncoding struct {
	*Encoding
}

// StandardEncodingValue is the singleton instance of StandardEncoding.
var StandardEncodingValue = func() *StandardEncoding {
	e := &StandardEncoding{
		Encoding: NewEncoding(),
	}
	for _, entry := range standardEncodingTable {
		e.Add(entry.code, entry.name)
	}
	return e
}()

// EncodingName returns the name of this encoding.
func (e *StandardEncoding) EncodingName() string {
	return "StandardEncoding"
}

// WinAnsiEncoding is the Windows ANSI (cp1252) encoding vector.
type WinAnsiEncoding struct {
	*Encoding
}

// WinAnsiEncodingValue is the singleton instance of WinAnsiEncoding.
var WinAnsiEncodingValue = func() *WinAnsiEncoding {
	e := &WinAnsiEncoding{
		Encoding: NewEncoding(),
	}
	for _, entry := range winAnsiEncodingTable {
		e.Add(entry.code, entry.name)
	}
	// In WinAnsiEncoding, all unused codes greater than 0o41 map to bullet.
	for i := 0o41; i <= 255; i++ {
		if !e.ContainsCode(i) {
			e.Add(i, "bullet")
		}
	}
	return e
}()

// EncodingName returns the name of this encoding.
func (e *WinAnsiEncoding) EncodingName() string {
	return "WinAnsiEncoding"
}

// MacExpertEncoding is the Mac OS expert encoding vector.
type MacExpertEncoding struct {
	*Encoding
}

// MacExpertEncodingValue is the singleton instance of MacExpertEncoding.
var MacExpertEncodingValue = func() *MacExpertEncoding {
	e := &MacExpertEncoding{
		Encoding: NewEncoding(),
	}
	for _, entry := range macExpertEncodingTable {
		e.Add(entry.code, entry.name)
	}
	return e
}()

// EncodingName returns the name of this encoding.
func (e *MacExpertEncoding) EncodingName() string {
	return "MacExpertEncoding"
}

// MacRomanEncoding is the Mac OS Roman encoding vector.
type MacRomanEncoding struct {
	*Encoding
}

// MacRomanEncodingValue is the singleton instance of MacRomanEncoding.
var MacRomanEncodingValue = func() *MacRomanEncoding {
	e := &MacRomanEncoding{
		Encoding: NewEncoding(),
	}
	for _, entry := range macRomanEncodingTable {
		e.Add(entry.code, entry.name)
	}
	return e
}()

// EncodingName returns the name of this encoding.
func (e *MacRomanEncoding) EncodingName() string {
	return "MacRomanEncoding"
}

// TryGetNamedEncoding returns a known encoding instance for the given NameToken.
// It returns false if the name does not match any built-in PDF encoding.
func TryGetNamedEncoding(name *tokens.NameToken) (*Encoding, bool) {
	if name == nil {
		return nil, false
	}

	switch {
	case name.Equals(tokens.StandardEncoding):
		return StandardEncodingValue.Encoding, true
	case name.Equals(tokens.WinAnsiEncoding):
		return WinAnsiEncodingValue.Encoding, true
	case name.Equals(tokens.MacExpertEncoding):
		return MacExpertEncodingValue.Encoding, true
	case name.Equals(tokens.MacRomanEncoding):
		return MacRomanEncodingValue.Encoding, true
	case name.Equals(tokens.SymbolEncoding):
		return SymbolEncodingValue.Encoding, true
	case name.Equals(tokens.ZapfDingbatsEncoding):
		return ZapfDingbatsEncodingValue.Encoding, true
	}

	return nil, false
}
