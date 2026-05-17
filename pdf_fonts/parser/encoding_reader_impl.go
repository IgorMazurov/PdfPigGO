// Package parser provides types for parsing PDF font structures.
package parser

import (
	"fmt"
	"strings"

	"github.com/uglytoad/pdfpig/go/fonts"
	enc "github.com/uglytoad/pdfpig/go/fonts/encodings"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// encodingReader provides a full implementation of EncodingReader.
type encodingReader struct {
	scanner tokenization.PdfTokenScanner
}

// NewEncodingReader creates a new EncodingReader with the given scanner.
func NewEncodingReader(scanner tokenization.PdfTokenScanner) EncodingReader {
	return &encodingReader{scanner: scanner}
}

// Read extracts the encoding from the given font dictionary, using the optional
// font descriptor and fallback encoding when available.
func (r *encodingReader) Read(
	fontDictionary *tokens.DictionaryToken,
	descriptor *fonts.FontDescriptor,
	fontEncoding *enc.Encoding,
) *enc.Encoding {
	baseEncodingObject, ok := fontDictionary.TryGet(tokens.Encoding)
	if !ok {
		return nil
	}

	name, ok := parts.TryGet[*tokens.NameToken](baseEncodingObject, r.scanner)
	if ok {
		namedEncoding := tryGetNamedEncoding(descriptor, name)
		if namedEncoding != nil {
			return namedEncoding
		}

		baseFontName, ok := parts.TryGet[*tokens.NameToken](
			getDictEntry(fontDictionary, tokens.BaseFont), r.scanner)
		if ok {
			if strings.EqualFold(baseFontName.Data(), "ZapfDingbats") {
				return enc.ZapfDingbatsEncodingValue.Encoding
			}

			if strings.EqualFold(baseFontName.Data(), "Symbol") {
				return enc.SymbolEncodingValue.Encoding
			}

			return enc.WinAnsiEncodingValue.Encoding
		}
	}

	encodingDictionary, err := parts.GetByToken[*tokens.DictionaryToken](baseEncodingObject, r.scanner)
	if err != nil || encodingDictionary == nil {
		return nil
	}

	return readEncodingDictionary(r.scanner, encodingDictionary, fontEncoding)
}

func getDictEntry(dict *tokens.DictionaryToken, name *tokens.NameToken) tokens.Token {
	token, _ := dict.TryGet(name)
	return token
}

// readEncodingDictionary reads an encoding dictionary and returns the resulting encoding.
// Returns nil on error (matching C# throw behavior mapped to Go conventions).
func readEncodingDictionary(scanner tokenization.PdfTokenScanner, encodingDictionary *tokens.DictionaryToken, fontEncoding *enc.Encoding) *enc.Encoding {
	if encodingDictionary == nil {
		return nil
	}

	var baseEncoding *enc.Encoding

	baseEncodingToken, ok := encodingDictionary.TryGet(tokens.BaseEncoding)
	if ok {
		if baseEncodingName, ok := any(baseEncodingToken).(*tokens.NameToken); ok {
			var found bool
			baseEncoding, found = enc.TryGetNamedEncoding(baseEncodingName)
			if !found {
				return nil
			}
		} else {
			return nil
		}
	} else {
		if fontEncoding != nil {
			baseEncoding = fontEncoding
		} else {
			baseEncoding = enc.StandardEncodingValue.Encoding
		}
	}

	differencesBase, ok := encodingDictionary.TryGet(tokens.Differences)
	if !ok {
		return baseEncoding
	}

	differenceArray, err := parts.GetByToken[*tokens.ArrayToken](differencesBase, scanner)
	if err != nil || differenceArray == nil {
		return baseEncoding
	}

	differences := processDifferences(differenceArray)

	newEncoding, err := enc.NewDifferenceBasedEncoding(baseEncoding, differences)
	if err != nil {
		return baseEncoding
	}

	return newEncoding.Encoding
}

// processDifferences parses the Differences array into a list of (code, name) pairs.
func processDifferences(differenceArray *tokens.ArrayToken) []enc.Difference {
	differences := make([]enc.Difference, 0, differenceArray.Length())

	if differenceArray.Length() == 0 {
		return differences
	}

	currentIndex := -1
	for i := 0; i < differenceArray.Length(); i++ {
		entry := differenceArray.Get(i)

		if number, ok := any(entry).(*tokens.NumericToken); ok {
			currentIndex = int(number.IntVal())
		} else if name, ok := any(entry).(*tokens.NameToken); ok {
			differences = append(differences, enc.Difference{Code: currentIndex, Name: name.Data()})
			currentIndex++
		} else {
			panic(fmt.Errorf("unexpected entry in the differences array: %v", differenceArray))
		}
	}

	return differences
}

// tryGetNamedEncoding attempts to resolve a named encoding. For symbolic fonts,
// it defaults to StandardEncoding. Otherwise it looks up by name token.
func tryGetNamedEncoding(descriptor *fonts.FontDescriptor, encodingName *tokens.NameToken) *enc.Encoding {
	if descriptor != nil && descriptor.Flags.HasFlag(fonts.Symbolic) {
		return enc.StandardEncodingValue.Encoding
	}

	result, ok := enc.TryGetNamedEncoding(encodingName)
	if !ok {
		return nil
	}

	return result
}
