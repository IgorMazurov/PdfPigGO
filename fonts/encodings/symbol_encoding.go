// Package encodings provides character encoding types for PDF fonts.
package encodings

// SymbolEncoding is the Symbol PostScript encoding vector.
type SymbolEncoding struct {
	*Encoding
}

// SymbolEncodingValue is the singleton instance of SymbolEncoding.
var SymbolEncodingValue = func() *SymbolEncoding {
	e := &SymbolEncoding{
		Encoding: NewEncoding(),
	}
	for _, entry := range symbolEncodingTable {
		e.Add(entry.code, entry.name)
	}
	return e
}()

// EncodingName returns the name of this encoding.
func (e *SymbolEncoding) EncodingName() string {
	return "SymbolEncoding"
}
