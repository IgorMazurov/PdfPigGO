// Package encodings provides character encoding types for PDF fonts.
package encodings

// MacOsRomanEncoding is similar to MacRomanEncoding with 15 additional entries.
type MacOsRomanEncoding struct {
	*MacRomanEncoding
}

// MacOsRomanEncodingValue is the singleton instance of MacOsRomanEncoding.
var MacOsRomanEncodingValue = func() *MacOsRomanEncoding {
	e := &MacOsRomanEncoding{
		MacRomanEncoding: &MacRomanEncoding{
			Encoding: NewEncoding(),
		},
	}
	for _, entry := range macOsRomanEncodingTable {
		e.Add(entry.code, entry.name)
	}
	return e
}()

// EncodingName returns the name of this encoding.
func (e *MacOsRomanEncoding) EncodingName() string {
	return "MacOsRomanEncoding"
}

// macOsRomanEncodingTable maps octal character codes to glyph names for the Mac OS Roman encoding supplement.
var macOsRomanEncodingTable = []struct {
	code int
	name string
}{
	{0o255, "notequal"},
	{0o260, "infinity"},
	{0o262, "lessequal"},
	{0o263, "greaterequal"},
	{0o266, "partialdiff"},
	{0o267, "summation"},
	{0o270, "product"},
	{0o271, "pi"},
	{0o272, "integral"},
	{0o275, "Omega"},
	{0o303, "radical"},
	{0o305, "approxequal"},
	{0o306, "Delta"},
	{0o327, "lozenge"},
	{0o333, "Euro"},
	{0o360, "apple"},
}
