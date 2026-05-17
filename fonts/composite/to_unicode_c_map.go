package composite

import (
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/cmap"
)

// ToUnicodeCMap defines the information content (actual text) of the font
// as opposed to the display format.
type ToUnicodeCMap struct {
	cMap                         *cmap.CMap
	isUsingIdentityAsUnicodeMap  bool
}

// NewToUnicodeCMap creates a new ToUnicodeCMap instance.
func NewToUnicodeCMap(cmapVal *cmap.CMap) *ToUnicodeCMap {
	t := &ToUnicodeCMap{}

	if cmapVal != nil {
		t.cMap = cmapVal
		t.isUsingIdentityAsUnicodeMap = strings.HasPrefix(strings.ToLower(cmapVal.Name()), "identity-")
	}

	return t
}

// CanMapToUnicode reports whether the font provides a CMap to map CIDs to Unicode values.
func (t *ToUnicodeCMap) CanMapToUnicode() bool {
	return t.cMap != nil
}

// IsUsingIdentityAsUnicodeMap reports whether this document is unexpectedly using
// a predefined Identity-H/V CMap as its ToUnicode CMap.
func (t *ToUnicodeCMap) IsUsingIdentityAsUnicodeMap() bool {
	return t.isUsingIdentityAsUnicodeMap
}

// TryGet returns the Unicode string for the given character code.
// Returns true if a mapping exists, false otherwise.
func (t *ToUnicodeCMap) TryGet(code int) (string, bool) {
	if t.cMap == nil {
		return "", false
	}

	return t.cMap.TryConvertToUnicode(code)
}

// ReadCode reads a character code from the input bytes according to the CMap's codespace ranges.
// Returns an error if the internal CMap is nil or if reading fails.
func (t *ToUnicodeCMap) ReadCode(bytes core.InputBytes, useLenientParsing bool) (int, error) {
	if t.cMap == nil {
		return 0, nil
	}

	return t.cMap.ReadCode(bytes, useLenientParsing)
}
