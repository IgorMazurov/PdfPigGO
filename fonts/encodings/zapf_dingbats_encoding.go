// Package encodings provides character encoding types for PDF fonts.
package encodings

// ZapfDingbatsEncoding is the Zapf Dingbats PostScript encoding vector.
type ZapfDingbatsEncoding struct {
	*Encoding
}

// ZapfDingbatsEncodingValue is the singleton instance of ZapfDingbatsEncoding.
var ZapfDingbatsEncodingValue = func() *ZapfDingbatsEncoding {
	e := &ZapfDingbatsEncoding{
		Encoding: NewEncoding(),
	}
	for _, entry := range zapfDingbatsEncodingTable {
		e.Add(entry.code, entry.name)
	}
	return e
}()

// EncodingName returns the name of this encoding.
func (e *ZapfDingbatsEncoding) EncodingName() string {
	return "ZapfDingbatsEncoding"
}
