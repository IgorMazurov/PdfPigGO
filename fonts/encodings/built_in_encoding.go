// Package encodings provides character encoding types for PDF fonts.
package encodings

// BuiltInEncoding represents an encoding built into a TrueType font.
type BuiltInEncoding struct {
	*Encoding
}

// NewBuiltInEncoding creates a new BuiltInEncoding from the given code-to-name mapping.
func NewBuiltInEncoding(codeToName map[int]string) *BuiltInEncoding {
	e := &BuiltInEncoding{
		Encoding: NewEncoding(),
	}
	for code, name := range codeToName {
		e.Add(code, name)
	}
	return e
}

// EncodingName returns the name of this encoding.
func (e *BuiltInEncoding) EncodingName() string {
	return "built-in (TTF)"
}
