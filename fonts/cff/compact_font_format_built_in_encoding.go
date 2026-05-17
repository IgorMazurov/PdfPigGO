// Package cff provides types for parsing and working with Compact Font Format (CFF) data in PDF files.
package cff

// CFFBuiltInEncodingSupplement represents a single supplement entry in a built-in
// CFF encoding, mapping a character code to an SID and glyph name.
type CFFBuiltInEncodingSupplement struct {
	Code int
	Sid  int
	Name string
}

// NewCFFBuiltInEncodingSupplement creates a new supplement entry.
func NewCFFBuiltInEncodingSupplement(code, sid int, name string) CFFBuiltInEncodingSupplement {
	return CFFBuiltInEncodingSupplement{
		Code: code,
		Sid:  sid,
		Name: name,
	}
}

// CFFBuiltInEncoding is the base type for built-in CFF encodings.
// It holds a list of supplement entries that define character-to-glyph mappings.
type CFFBuiltInEncoding struct {
	*CompactFontFormatBaseEncoding
	Supplements []CFFBuiltInEncodingSupplement
}

// NewCFFBuiltInEncoding creates a new built-in CFF encoding from the given supplements.
func NewCFFBuiltInEncoding(supplements []CFFBuiltInEncodingSupplement) *CFFBuiltInEncoding {
	e := &CFFBuiltInEncoding{
		CompactFontFormatBaseEncoding: NewCompactFontFormatBaseEncoding(),
		Supplements:                   supplements,
	}
	for _, s := range supplements {
		e.AddWithName(s.Code, s.Sid, s.Name)
	}
	return e
}
